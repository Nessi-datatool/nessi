package delta

import (
	"fmt"
	"os"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/parquet"
	"github.com/apache/arrow/go/v15/parquet/compress"
	"github.com/apache/arrow/go/v15/parquet/pqarrow"
)

// WriteRecordToParquet writes an Arrow record to a Parquet file.
func WriteRecordToParquet(path string, record arrow.Record) error {
	osFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer osFile.Close()

	// Ensure correct compression codec is used.
	pqWriterProps := parquet.NewWriterProperties(parquet.WithCompression(compress.Codecs.Snappy))
	arrowWriterProps := pqarrow.DefaultWriterProps() // This is ArrowWriterProperties (struct)

	// pqarrow.NewFileWriter handles writing Arrow data in Parquet format to the sink (osFile).
	arrowFw, err := pqarrow.NewFileWriter(record.Schema(), osFile, pqWriterProps, arrowWriterProps)
	if err != nil {
		return fmt.Errorf("failed to create arrow file writer: %w", err)
	}

	if err := arrowFw.Write(record); err != nil {
		_ = arrowFw.Close()
		return fmt.Errorf("failed to write record: %w", err)
	}

	if err := arrowFw.Close(); err != nil {
		return fmt.Errorf("failed to close arrow writer: %w", err)
	}

	return nil
}
