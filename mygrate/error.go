package mygrate

import (
	"errors"
)

var (
	// ErrInitFn will be returned if the stores init return an error.
	ErrInitFn = errors.New("store returned an error on init")
	// ErrDownFn will be returned if the down func will return an error.
	ErrDownFn = errors.New("migrations down returned an error")
	// ErrStore will be returned if the store has an error.
	ErrStore = errors.New("store returned an error")
	// ErrUpFn will be returned if the up func will return an error.
	ErrUpFn = errors.New("migrations up returned an error")
)

// Error is a custom mygrate error type.
type Error struct {
	ID          string // ID of the migration.
	Err         error  // The underlying error from the store e.g. sql.ErrNoRows.
	InternalErr error
}

// Error makes this struct an error.
func (e Error) Error() string {
	return e.InternalErr.Error() + ": " + e.Err.Error()
}

// Is implements errors.Is.
func (e Error) Is(t error) bool {
	return (e.InternalErr.Error()) == t.Error() || (e.Err.Error() == t.Error())
}

// Unwrap implements errors.Unwrap.
func (e Error) Unwrap() error {
	return e.Err
}

func errInit(err error) error {
	return &Error{Err: err, InternalErr: ErrInitFn}
}

func errStore(id string, err error) error {
	return &Error{ID: id, Err: err, InternalErr: ErrStore}
}

func errUp(id string, err error) error {
	return &Error{ID: id, Err: err, InternalErr: ErrUpFn}
}

func errDown(id string, err error) error {
	return &Error{ID: id, Err: err, InternalErr: ErrDownFn}
}
