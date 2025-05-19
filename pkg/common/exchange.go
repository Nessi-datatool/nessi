package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/ipc"
	"github.com/apache/arrow/go/v14/arrow/memory"
)

// ExchangeStrategy defines the interface for data exchange methods
type ExchangeStrategy interface {
	Serialize(data interface{}) ([]byte, error)
	Deserialize(data []byte, target interface{}) error
	SupportsType(data interface{}) bool
}

// ExchangeConfig holds configuration for data exchange
type ExchangeConfig struct {
	ArrowThreshold int64 `json:"arrow_threshold"` // Size threshold in bytes
	MaxMemoryUsage int64 `json:"max_memory_usage"`
	UseCompression bool  `json:"use_compression"`
}

// ExchangeStats tracks performance metrics
type ExchangeStats struct {
	RowCount      int64
	SizeBytes     int64
	MemoryUsage   int64
	SerializeMs   int64
	DeserializeMs int64
	CompressionMs int64
	LastOperation time.Time
}

// ExchangeManager orchestrates data exchange
type ExchangeManager struct {
	strategies []ExchangeStrategy
	config     ExchangeConfig
	stats      map[string]*ExchangeStats
	pool       memory.Allocator
}

// NewExchangeManager creates a new exchange manager
func NewExchangeManager(config ExchangeConfig) *ExchangeManager {
	return &ExchangeManager{
		strategies: make([]ExchangeStrategy, 0),
		config:     config,
		stats:      make(map[string]*ExchangeStats),
		pool:       memory.NewGoAllocator(),
	}
}

// RegisterStrategy adds a new exchange strategy
func (m *ExchangeManager) RegisterStrategy(strategy ExchangeStrategy) {
	m.strategies = append(m.strategies, strategy)
}

// SelectStrategy chooses the appropriate strategy based on data
func (m *ExchangeManager) SelectStrategy(data interface{}) (ExchangeStrategy, error) {
	// Check memory usage before proceeding
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	if memStats.Alloc > uint64(m.config.MaxMemoryUsage) {
		return nil, fmt.Errorf("memory usage exceeds limit: %d > %d", memStats.Alloc, m.config.MaxMemoryUsage)
	}

	// Prefer ArrowStrategy for arrow.Record
	for _, strategy := range m.strategies {
		if _, ok := data.(arrow.Record); ok {
			if _, isArrow := strategy.(*ArrowStrategy); isArrow {
				return strategy, nil
			}
		}
	}
	// Fallback: choose first strategy that supports the type
	for _, strategy := range m.strategies {
		if strategy.SupportsType(data) {
			return strategy, nil
		}
	}
	return nil, fmt.Errorf("no suitable strategy found for data type")
}

// ArrowStrategy implements ExchangeStrategy for Arrow data
type ArrowStrategy struct {
	pool memory.Allocator
}

func NewArrowStrategy() *ArrowStrategy {
	return &ArrowStrategy{
		pool: memory.NewGoAllocator(),
	}
}

func (s *ArrowStrategy) Serialize(data interface{}) ([]byte, error) {
	record, ok := data.(arrow.Record)
	if !ok {
		return nil, fmt.Errorf("data is not an arrow.Record")
	}

	buf := new(bytes.Buffer)
	writer := ipc.NewWriter(buf, ipc.WithSchema(record.Schema()))
	defer writer.Close()

	if err := writer.Write(record); err != nil {
		return nil, fmt.Errorf("failed to write record: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *ArrowStrategy) Deserialize(data []byte, target interface{}) error {
	reader, err := ipc.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create Arrow reader: %w", err)
	}

	record, ok := target.(*arrow.Record)
	if !ok {
		return fmt.Errorf("target is not a *arrow.Record")
	}

	if !reader.Next() {
		return fmt.Errorf("no record found in data")
	}

	*record = reader.Record()
	return nil
}

func (s *ArrowStrategy) SupportsType(data interface{}) bool {
	_, ok := data.(arrow.Record)
	return ok
}

// JSONStrategy implements ExchangeStrategy for JSON data
type JSONStrategy struct{}

func NewJSONStrategy() *JSONStrategy {
	return &JSONStrategy{}
}

func (s *JSONStrategy) Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (s *JSONStrategy) Deserialize(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

func (s *JSONStrategy) SupportsType(data interface{}) bool {
	// JSON strategy supports any data that can be marshaled
	_, err := json.Marshal(data)
	return err == nil
}
