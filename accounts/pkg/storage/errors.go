package storage

import (
	"fmt"
)

type notFoundErr struct {
	typ, id string
}

func (e *notFoundErr) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.typ, e.id)
}

// IsNotFoundErr can be returned by repo Load and Delete operations
func IsNotFoundErr(e error) bool {
	_, ok := e.(*notFoundErr)
	return ok
}

type unsupportedErr struct{}

func (e *unsupportedErr) Error() string {
	return "unsupported"
}

// IsUnsupportedErr can be returned by repo Write operations
func IsUnsupportedErr(e error) bool {
	_, ok := e.(*unsupportedErr)
	return ok
}
