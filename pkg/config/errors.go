package config

import "emperror.dev/errors"

const (
	ErrNotValidValue = errors.Sentinel("value is not valid")
	ErrNotStructure  = errors.Sentinel("value is not a structure")
	ErrNotPointer    = errors.Sentinel("value must be a pointer")
	ErrDuplicateKey  = errors.Sentinel("key duplication found")
)
