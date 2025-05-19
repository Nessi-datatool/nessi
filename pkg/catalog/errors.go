package catalog

import "errors"

// Common errors
var (
	// ErrNotImplemented is returned when a feature is not yet implemented
	ErrNotImplemented = errors.New("feature not implemented")

	// ErrUnsupportedCatalogType is returned when an unsupported catalog type is requested
	ErrUnsupportedCatalogType = errors.New("unsupported catalog type")

	// ErrNotConnected is returned when an operation is attempted on a catalog that is not connected
	ErrNotConnected = errors.New("not connected to data catalog")

	// ErrAlreadyConnected is returned when attempting to connect to an already connected catalog
	ErrAlreadyConnected = errors.New("already connected to data catalog")

	// ErrInvalidConfiguration is returned when invalid configuration is provided
	ErrInvalidConfiguration = errors.New("invalid catalog configuration")

	// ErrResourceNotFound is returned when a requested resource is not found
	ErrResourceNotFound = errors.New("resource not found in catalog")

	// ErrOperationNotSupported is returned when an operation is not supported by a catalog
	ErrOperationNotSupported = errors.New("operation not supported by catalog")
)
