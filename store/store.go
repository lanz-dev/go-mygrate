// Package store provides some default store implementations.
package store

import (
	"errors"
)

var (
	// ErrIDNotFound will be returned if ID is not found.
	ErrIDNotFound = errors.New("id not found")
)
