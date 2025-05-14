package profile

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// DistributionAnalysis contains detailed distribution information
type DistributionAnalysis struct {
	Histogram      map[string]int    `json:"histogram"`
	Frequencies    map[string]float64 `json:"frequencies"`
	TopValues      []ValueCount      `json:"top_values"`
	BottomValues   []ValueCount      `json:"bottom_values"`
	Skewness       float64           `json:"skewness"`
	Kurtosis       float64           `json:"kurtosis"`
	IsNormal       bool              `json:"is_normal"`
	Entropy        float64           `json:"entropy"`
}

// ValueCount is defined in shared_types.go

// DataQualityScore is defined in shared_types.go

// PatternInfo is defined in shared_types.go

// EnhancedProfile is defined in shared_types.go

// AdvancedProfiler extends the basic Profiler with advanced analytics
type AdvancedProfiler struct {
	*Profiler
	patternLibrary map[string]*regexp.Regexp
}

// NewAdvancedProfiler creates a new AdvancedProfiler
func NewAdvancedProfiler(tablePath string) *AdvancedProfiler {
	return &AdvancedProfiler{
		Profiler:       NewProfiler(tablePath),
		patternLibrary: initializePatternLibrary(),
	}
}

// initializePatternLibrary creates a library of common data patterns
func initializePatternLibrary() map[string]*regexp.Regexp {
	patterns := map[string]string{
		"email":        `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		"url":          `^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`,
		"phone_number": `^\+?[0-9]{10,15}$`,
		"ip_address":   `^(\d{1,3}\.){3}\d{1,3}$`,
		"date_iso":     `^\d{4}-\d{2}-\d{2}$`,
		"date_us":      `^\d{1,2}/\d{1,2}/\d{4}$`,
		"date_eu":      `^\d{1,2}\.\d{1,2}\.\d{4}$`,
		"time":         `^\d{1,2}:\d{2}(:\d{2})?$`,
		"credit_card":  `^\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}$`,
		"zip_code":     `^\d{5}(-\d{4})?$`,
		"currency":     `^[$€£¥]\d+(\.\d{2})?$`,
		"uuid":         `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	}

	compiled := make(map[string]*regexp.Regexp)
	for name, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			compiled[name] = re
		}
	}

	return compiled
}

// EnhanceProfile adds advanced analytics to a basic profile
func (p *AdvancedProfiler) EnhanceProfile(profile *Profile, values []interface{}) *EnhancedProfile {
	enhanced := &EnhancedProfile{
		Profile: profile,
	}

	// Add distribution analysis
	if profile.Type == "integer" || profile.Type == "float" {
		enhanced.Distribution = p.analyzeNumericDistribution(values)
	} else if profile.Type == "string" {
		enhanced.Distribution = p.analyzeStringDistribution(values)
	}

	// Add quality score
	enhanced.QualityScore = p.calculateQualityScore(profile, values)

	// Add detailed pattern analysis
	if profile.Type == "string" {
		enhanced.DetailedPatterns = p.analyzeStringPatterns(values)
		enhanced.FormatConsistency = p.calculateFormatConsistency(values)
	}

	// Add type consistency
	enhanced.TypeConsistency = p.calculateTypeConsistency(values, profile.Type)

	return enhanced
}

// analyzeNumericDistribution performs detailed distribution analysis for numeric columns
func (p *AdvancedProfiler) analyzeNumericDistribution(values []interface{}) *DistributionAnalysis {
	if len(values) == 0 {
		return nil
	}

	// Convert to float64 for analysis
	floatValues := make([]float64, 0, len(values))
	for _, v := range values {
		if v == nil {
			continue
		}

		var fv float64
		switch vt := v.(type) {
		case int:
			fv = float64(vt)
		case int32:
			fv = float64(vt)
		case int64:
			fv = float64(vt)
		case float32:
			fv = float64(vt)
		case float64:
			fv = vt
		default:
			continue
		}

		floatValues = append(floatValues, fv)
	}

	if len(floatValues) == 0 {
		return nil
	}

	// Calculate basic statistics
	stats := calculateStatistics(floatValues)

	// Create histogram
	numBins := int(math.Min(10, float64(len(floatValues))))
	if numBins < 2 {
		numBins = 2
	}

	binWidth := (stats.Max - stats.Min) / float64(numBins)
	histogram := make(map[string]int)
	
	for _, val := range floatValues {
		binIdx := int((val - stats.Min) / binWidth)
		if binIdx >= numBins {
			binIdx = numBins - 1
		}
		
		binLabel := fmt.Sprintf("%.2f-%.2f", stats.Min+float64(binIdx)*binWidth, stats.Min+float64(binIdx+1)*binWidth)
		histogram[binLabel]++
	}

	// Calculate frequencies
	frequencies := make(map[string]float64)
	for bin, count := range histogram {
		frequencies[bin] = float64(count) / float64(len(floatValues))
	}

	// Calculate skewness and kurtosis
	skewness := calculateSkewness(floatValues, stats.Mean, stats.StdDev)
	kurtosis := calculateKurtosis(floatValues, stats.Mean, stats.StdDev)
	
	// Determine if distribution is approximately normal
	isNormal := math.Abs(skewness) < 0.5 && math.Abs(kurtosis-3) < 1.0

	// Calculate entropy
	entropy := calculateEntropy(frequencies)

	// Create value counts for top/bottom values
	valueCountMap := make(map[float64]int)
	for _, val := range floatValues {
		valueCountMap[val]++
	}

	var valueCounts []ValueCount
	for val, count := range valueCountMap {
		valueCounts = append(valueCounts, ValueCount{
			Value: val,
			Count: count,
		})
	}

	// Sort by count
	sort.Slice(valueCounts, func(i, j int) bool {
		return valueCounts[i].Count > valueCounts[j].Count
	})

	topValues := valueCounts
	if len(topValues) > 5 {
		topValues = topValues[:5]
	}

	// Sort by count (ascending)
	sort.Slice(valueCounts, func(i, j int) bool {
		return valueCounts[i].Count < valueCounts[j].Count
	})

	bottomValues := valueCounts
	if len(bottomValues) > 5 {
		bottomValues = bottomValues[:5]
	}

	return &DistributionAnalysis{
		Histogram:    histogram,
		Frequencies:  frequencies,
		TopValues:    topValues,
		BottomValues: bottomValues,
		Skewness:     skewness,
		Kurtosis:     kurtosis,
		IsNormal:     isNormal,
		Entropy:      entropy,
	}
}

// analyzeStringDistribution performs detailed distribution analysis for string columns
func (p *AdvancedProfiler) analyzeStringDistribution(values []interface{}) *DistributionAnalysis {
	if len(values) == 0 {
		return nil
	}

	// Count string lengths
	lengthHistogram := make(map[string]int)
	valueCountMap := make(map[string]int)
	
	for _, v := range values {
		if v == nil {
			continue
		}

		s, ok := v.(string)
		if !ok {
			continue
		}

		// Group by length
		lengthBin := fmt.Sprintf("%d", len(s))
		lengthHistogram[lengthBin]++
		
		// Count values
		valueCountMap[s]++
	}

	// Calculate frequencies
	frequencies := make(map[string]float64)
	for length, count := range lengthHistogram {
		frequencies[length] = float64(count) / float64(len(values))
	}

	// Calculate entropy
	entropy := calculateEntropy(frequencies)

	// Create value counts for top/bottom values
	var valueCounts []ValueCount
	for val, count := range valueCountMap {
		valueCounts = append(valueCounts, ValueCount{
			Value: val,
			Count: count,
		})
	}

	// Sort by count
	sort.Slice(valueCounts, func(i, j int) bool {
		return valueCounts[i].Count > valueCounts[j].Count
	})

	topValues := valueCounts
	if len(topValues) > 5 {
		topValues = topValues[:5]
	}

	// Sort by count (ascending)
	sort.Slice(valueCounts, func(i, j int) bool {
		return valueCounts[i].Count < valueCounts[j].Count
	})

	bottomValues := valueCounts
	if len(bottomValues) > 5 {
		bottomValues = bottomValues[:5]
	}

	return &DistributionAnalysis{
		Histogram:    lengthHistogram,
		Frequencies:  frequencies,
		TopValues:    topValues,
		BottomValues: bottomValues,
		Entropy:      entropy,
	}
}

// calculateSkewness calculates the skewness of a distribution
func calculateSkewness(values []float64, mean, stdDev float64) float64 {
	if len(values) == 0 || stdDev == 0 {
		return 0
	}

	var sum float64
	for _, val := range values {
		sum += math.Pow((val-mean)/stdDev, 3)
	}

	return sum / float64(len(values))
}

// calculateKurtosis calculates the kurtosis of a distribution
func calculateKurtosis(values []float64, mean, stdDev float64) float64 {
	if len(values) == 0 || stdDev == 0 {
		return 0
	}

	var sum float64
	for _, val := range values {
		sum += math.Pow((val-mean)/stdDev, 4)
	}

	return sum / float64(len(values))
}

// calculateEntropy calculates the Shannon entropy of a distribution
func calculateEntropy(frequencies map[string]float64) float64 {
	var entropy float64
	for _, freq := range frequencies {
		if freq > 0 {
			entropy -= freq * math.Log2(freq)
		}
	}
	return entropy
}

// calculateQualityScore calculates quality metrics for a column
func (p *AdvancedProfiler) calculateQualityScore(profile *Profile, values []interface{}) *DataQualityScore {
	// Calculate completeness (1 - null ratio)
	completeness := 1.0 - (float64(profile.NullCount) / float64(profile.RowCount))

	// Calculate uniqueness
	uniqueness := float64(profile.Distinct) / float64(profile.RowCount-profile.NullCount)
	if profile.RowCount == profile.NullCount {
		uniqueness = 0
	}

	// Calculate consistency based on patterns
	consistency := p.calculateConsistency(profile, values)

	// Calculate accuracy based on anomalies
	accuracy := 1.0 - (float64(len(profile.Anomalies)) / float64(profile.RowCount-profile.NullCount))
	if profile.RowCount == profile.NullCount {
		accuracy = 0
	}

	// Calculate overall score
	overall := (completeness + consistency + accuracy + uniqueness) / 4.0

	// Determine recommended type
	recommendedType := profile.Type
	if profile.Type == "string" {
		recommendedType = p.inferType(values)
	}

	return &DataQualityScore{
		Completeness:    completeness,
		Consistency:     consistency,
		Accuracy:        accuracy,
		Uniqueness:      uniqueness,
		Overall:         overall,
		RecommendedType: recommendedType,
	}
}

// calculateConsistency calculates the consistency score based on patterns
func (p *AdvancedProfiler) calculateConsistency(profile *Profile, values []interface{}) float64 {
	if len(values) == 0 {
		return 0
	}

	if profile.Type == "string" {
		// For string columns, check pattern consistency
		patternCounts := make(map[string]int)
		totalNonNull := 0

		for _, v := range values {
			if v == nil {
				continue
			}

			s, ok := v.(string)
			if !ok {
				continue
			}

			totalNonNull++
			patternFound := false

			// Check against pattern library
			for name, re := range p.patternLibrary {
				if re.MatchString(s) {
					patternCounts[name]++
					patternFound = true
					break
				}
			}

			if !patternFound {
				patternCounts["unknown"]++
			}
		}

		// Find dominant pattern
		var maxCount int
		for _, count := range patternCounts {
			if count > maxCount {
				maxCount = count
			}
		}

		return float64(maxCount) / float64(totalNonNull)
	} else if profile.Type == "integer" || profile.Type == "float" {
		// For numeric columns, check for outliers
		outlierCount := 0
		for _, anomaly := range profile.Anomalies {
			if anomaly.Type == "outlier" {
				outlierCount++
			}
		}

		return 1.0 - (float64(outlierCount) / float64(len(values)))
	}

	return 1.0
}

// inferType tries to infer a more specific type from string values
func (p *AdvancedProfiler) inferType(values []interface{}) string {
	if len(values) == 0 {
		return "string"
	}

	// Count matches for each pattern
	patternCounts := make(map[string]int)
	totalNonNull := 0

	for _, v := range values {
		if v == nil {
			continue
		}

		s, ok := v.(string)
		if !ok {
			continue
		}

		totalNonNull++

		// Check against pattern library
		for name, re := range p.patternLibrary {
			if re.MatchString(s) {
				patternCounts[name]++
				break
			}
		}

		// Check if it could be a number
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			patternCounts["numeric"]++
		}

		// Check if it could be a boolean
		if strings.ToLower(s) == "true" || strings.ToLower(s) == "false" {
			patternCounts["boolean"]++
		}
	}

	// Find dominant pattern
	var maxPattern string
	var maxCount int
	for pattern, count := range patternCounts {
		if count > maxCount {
			maxCount = count
			maxPattern = pattern
		}
	}

	// Only recommend a type if it's a strong match (>80%)
	threshold := int(float64(totalNonNull) * 0.8)
	if maxCount >= threshold {
		switch maxPattern {
		case "email":
			return "email"
		case "url":
			return "url"
		case "date_iso", "date_us", "date_eu":
			return "date"
		case "time":
			return "time"
		case "ip_address":
			return "ip"
		case "numeric":
			return "numeric"
		case "boolean":
			return "boolean"
		case "uuid":
			return "uuid"
		}
	}

	return "string"
}

// analyzeStringPatterns performs detailed pattern analysis for string columns
func (p *AdvancedProfiler) analyzeStringPatterns(values []interface{}) []PatternInfo {
	if len(values) == 0 {
		return nil
	}

	// Count matches for each pattern
	patternCounts := make(map[string]int)
	patternExamples := make(map[string][]string)
	totalNonNull := 0

	for _, v := range values {
		if v == nil {
			continue
		}

		s, ok := v.(string)
		if !ok {
			continue
		}

		totalNonNull++

		// Check against pattern library
		for name, re := range p.patternLibrary {
			if re.MatchString(s) {
				patternCounts[name]++
				
				// Store examples (up to 3 per pattern)
				if len(patternExamples[name]) < 3 {
					patternExamples[name] = append(patternExamples[name], s)
				}
				break
			}
		}
	}

	// Create pattern info
	var patterns []PatternInfo
	for name, count := range patternCounts {
		confidence := float64(count) / float64(totalNonNull)
		
		// Only include patterns with at least 10% confidence
		if confidence >= 0.1 {
			pattern := PatternInfo{
				Name:        name,
				Regex:       p.patternLibrary[name].String(),
				Description: getPatternDescription(name),
				Confidence:  confidence,
				Examples:    patternExamples[name],
			}
			patterns = append(patterns, pattern)
		}
	}

	// Sort by confidence
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Confidence > patterns[j].Confidence
	})

	return patterns
}

// getPatternDescription returns a human-readable description of a pattern
func getPatternDescription(pattern string) string {
	descriptions := map[string]string{
		"email":        "Email address",
		"url":          "Web URL",
		"phone_number": "Phone number",
		"ip_address":   "IP address",
		"date_iso":     "ISO date (YYYY-MM-DD)",
		"date_us":      "US date (MM/DD/YYYY)",
		"date_eu":      "European date (DD.MM.YYYY)",
		"time":         "Time (HH:MM[:SS])",
		"credit_card":  "Credit card number",
		"zip_code":     "ZIP code",
		"currency":     "Currency amount",
		"uuid":         "UUID",
	}

	if desc, ok := descriptions[pattern]; ok {
		return desc
	}
	return "Unknown pattern"
}

// calculateFormatConsistency measures the consistency of string formats
func (p *AdvancedProfiler) calculateFormatConsistency(values []interface{}) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create format signatures
	signatures := make(map[string]int)
	totalNonNull := 0

	for _, v := range values {
		if v == nil {
			continue
		}

		s, ok := v.(string)
		if !ok {
			continue
		}

		totalNonNull++
		signature := createFormatSignature(s)
		signatures[signature]++
	}

	// Find dominant signature
	var maxCount int
	for _, count := range signatures {
		if count > maxCount {
			maxCount = count
		}
	}

	return float64(maxCount) / float64(totalNonNull)
}

// createFormatSignature creates a format signature for a string
func createFormatSignature(s string) string {
	var signature strings.Builder
	for _, c := range s {
		if unicode.IsDigit(c) {
			signature.WriteRune('d')
		} else if unicode.IsLetter(c) {
			if unicode.IsUpper(c) {
				signature.WriteRune('U')
			} else {
				signature.WriteRune('l')
			}
		} else if unicode.IsSpace(c) {
			signature.WriteRune('s')
		} else {
			signature.WriteRune(c)
		}
	}
	return signature.String()
}

// calculateTypeConsistency measures the consistency of value types
func (p *AdvancedProfiler) calculateTypeConsistency(values []interface{}, expectedType string) float64 {
	if len(values) == 0 {
		return 0
	}

	matchCount := 0
	totalNonNull := 0

	for _, v := range values {
		if v == nil {
			continue
		}

		totalNonNull++
		
		switch expectedType {
		case "integer":
			_, ok := v.(int)
			if ok || isIntType(v) {
				matchCount++
			}
		case "float":
			_, ok := v.(float64)
			if ok || isFloatType(v) {
				matchCount++
			}
		case "string":
			_, ok := v.(string)
			if ok {
				matchCount++
			}
		case "boolean":
			_, ok := v.(bool)
			if ok {
				matchCount++
			}
		}
	}

	if totalNonNull == 0 {
		return 0
	}
	return float64(matchCount) / float64(totalNonNull)
}

// isIntType checks if a value is any integer type
func isIntType(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

// isFloatType checks if a value is any float type
func isFloatType(v interface{}) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}
	return false
}
