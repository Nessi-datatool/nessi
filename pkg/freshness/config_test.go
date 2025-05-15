package freshness

import (
	"testing"
	"time"
)

func TestSLAConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  SLAConfig
		wantErr bool
	}{
		{
			name: "Valid config",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			wantErr: false,
		},
		{
			name: "Empty table name",
			config: SLAConfig{
				TableName:         "",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			wantErr: true,
		},
		{
			name: "Empty table path",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			wantErr: true,
		},
		{
			name: "Zero expected frequency",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: 0,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			wantErr: true,
		},
		{
			name: "Negative expected frequency",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: -time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			wantErr: true,
		},
		{
			name: "Warning threshold less than 100",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  50,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			wantErr: true,
		},
		{
			name: "Critical threshold less than warning threshold",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  200,
				CriticalThreshold: 150,
				Enabled:           true,
			},
			wantErr: true,
		},
		{
			name: "Critical threshold equal to warning threshold",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 150,
				Enabled:           true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SLAConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "Valid hourly",
			input:   "1h",
			want:    time.Hour,
			wantErr: false,
		},
		{
			name:    "Valid daily",
			input:   "1d",
			want:    24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "Valid weekly",
			input:   "1w",
			want:    7 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "Valid monthly",
			input:   "1mo",
			want:    30 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "Valid minutes",
			input:   "30m",
			want:    30 * time.Minute,
			wantErr: false,
		},
		{
			name:    "Valid seconds",
			input:   "45s",
			want:    45 * time.Second,
			wantErr: false,
		},
		{
			name:    "Valid complex",
			input:   "1h30m",
			want:    90 * time.Minute,
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "Seconds",
			duration: 45 * time.Second,
			want:     "45 seconds",
		},
		{
			name:     "Minutes",
			duration: 30 * time.Minute,
			want:     "30 minutes",
		},
		{
			name:     "Hours",
			duration: 2 * time.Hour,
			want:     "2 hours",
		},
		{
			name:     "Days",
			duration: 3 * 24 * time.Hour,
			want:     "3 days",
		},
		{
			name:     "Weeks",
			duration: 2 * 7 * 24 * time.Hour,
			want:     "2 weeks",
		},
		{
			name:     "Months",
			duration: 2 * 30 * 24 * time.Hour,
			want:     "2 months",
		},
		{
			name:     "Zero",
			duration: 0,
			want:     "0 seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuration(tt.duration); got != tt.want {
				t.Errorf("FormatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
