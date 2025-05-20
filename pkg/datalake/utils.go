package datalake

import (
	"os"
)

// createDirIfNotExists creates a directory if it doesn't exist
func createDirIfNotExists(path string) error {
	return os.MkdirAll(path, 0755)
}

// writeFile writes data to a file
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// statFile gets file info
func statFile(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
