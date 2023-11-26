package constant

import "github.com/pkg/errors"

var (
	NullPointerError = errors.New("null pointer error")
	NotFoundError    = errors.New("not found error")
	UnknownError     = errors.New("unknown error")
)
