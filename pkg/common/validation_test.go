package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathValidator(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nessi-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp(tempDir, "test-file-*")
	assert.NoError(t, err)
	tempFile.Close()

	// Create a non-existent path
	nonExistentPath := filepath.Join(tempDir, "non-existent")

	tests := []struct {
		name      string
		validator *PathValidator
		path      string
		wantErr   bool
		errCode   ErrorCode
	}{
		{
			name:      "Empty path",
			validator: NewPathValidator(),
			path:      "",
			wantErr:   true,
			errCode:   ErrInvalidPath,
		},
		{
			name:      "Non-existent path not allowed",
			validator: NewPathValidator(),
			path:      nonExistentPath,
			wantErr:   true,
			errCode:   ErrInvalidPath,
		},
		{
			name: "Non-existent path allowed",
			validator: &PathValidator{
				AllowNonExistent: true,
			},
			path:    nonExistentPath,
			wantErr: false,
		},
		{
			name: "Directory required but file provided",
			validator: &PathValidator{
				RequireDirectory: true,
			},
			path:    tempFile.Name(),
			wantErr: true,
			errCode: ErrNotADirectory,
		},
		{
			name: "File required but directory provided",
			validator: &PathValidator{
				RequireFile: true,
			},
			path:    tempDir,
			wantErr: true,
			errCode: ErrNotAFile,
		},
		{
			name:      "Valid directory",
			validator: NewPathValidator(),
			path:      tempDir,
			wantErr:   false,
		},
		{
			name:      "Valid file",
			validator: NewPathValidator(),
			path:      tempFile.Name(),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := tt.validator.ValidatePath(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				// Check if it's a NessiError with the expected code
				var nessiErr *NessiError
				if assert.ErrorAs(t, err, &nessiErr) {
					assert.Equal(t, tt.errCode, nessiErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, gotPath)
			}
		})
	}
}

func TestValidateTablePath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nessi-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a fake Delta table
	deltaTablePath := filepath.Join(tempDir, "delta-table")
	err = os.Mkdir(deltaTablePath, 0755)
	assert.NoError(t, err)

	// Create _delta_log directory
	deltaLogPath := filepath.Join(deltaTablePath, "_delta_log")
	err = os.Mkdir(deltaLogPath, 0755)
	assert.NoError(t, err)

	// Test cases
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errCode ErrorCode
	}{
		{
			name:    "Valid Delta table",
			path:    deltaTablePath,
			wantErr: false,
		},
		{
			name:    "Not a Delta table",
			path:    tempDir,
			wantErr: true,
			errCode: ErrInvalidDeltaTable,
		},
		{
			name:    "Non-existent path",
			path:    filepath.Join(tempDir, "non-existent"),
			wantErr: true,
			errCode: ErrInvalidPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := ValidateTablePath(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				// Check if it's a NessiError with the expected code
				var nessiErr *NessiError
				if assert.ErrorAs(t, err, &nessiErr) {
					assert.Equal(t, tt.errCode, nessiErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, deltaTablePath, gotPath)
			}
		})
	}
}

func TestConfigValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator *ConfigValidator
		key       string
		value     string
		wantErr   bool
		errCode   ErrorCode
	}{
		{
			name: "Required config missing",
			validator: &ConfigValidator{
				Required: true,
			},
			key:     "host",
			value:   "",
			wantErr: true,
			errCode: ErrMissingConfig,
		},
		{
			name: "Required config present",
			validator: &ConfigValidator{
				Required: true,
			},
			key:     "host",
			value:   "localhost",
			wantErr: false,
		},
		{
			name: "Optional config missing",
			validator: &ConfigValidator{
				Required: false,
			},
			key:     "host",
			value:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.ValidateConfig(tt.key, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				// Check if it's a NessiError with the expected code
				var nessiErr *NessiError
				if assert.ErrorAs(t, err, &nessiErr) {
					assert.Equal(t, tt.errCode, nessiErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOutputFormat(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		wantFormat string
		wantErr   bool
		errCode   ErrorCode
	}{
		{
			name:       "Valid format - html",
			format:     "html",
			wantFormat: "html",
			wantErr:    false,
		},
		{
			name:       "Valid format - case insensitive",
			format:     "JSON",
			wantFormat: "json",
			wantErr:    false,
		},
		{
			name:     "Invalid format",
			format:   "invalid",
			wantErr:  true,
			errCode:  ErrInvalidFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFormat, err := ValidateOutputFormat(tt.format)
			if tt.wantErr {
				assert.Error(t, err)
				// Check if it's a NessiError with the expected code
				var nessiErr *NessiError
				if assert.ErrorAs(t, err, &nessiErr) {
					assert.Equal(t, tt.errCode, nessiErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFormat, gotFormat)
			}
		})
	}
}
