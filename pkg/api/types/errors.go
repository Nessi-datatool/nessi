package types

import "errors"

// Common errors
var (
	ErrNotImplemented          = errors.New("not implemented")
	ErrUnsupportedCatalogType = errors.New("unsupported catalog type")
)
