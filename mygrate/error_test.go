package mygrate_test

import (
	"errors"
	"testing"

	"github.com/lanz-dev/go-mygrate/mygrate"
)

func TestError_Error(t *testing.T) {
	t.Parallel() // marks TLog as capable of running in parallel with other tests

	e := mygrate.Error{
		Err:         errUnitTest,
		InternalErr: mygrate.ErrInitFn,
	}

	expected := "store returned an error on init: unittest"
	if e.Error() != expected {
		t.Fatalf("expected error to be '%s', got '%s'", expected, e.Error())
	}
}

func TestError_Is(t *testing.T) {
	t.Parallel()

	e := mygrate.Error{
		Err:         errUnitTest,
		InternalErr: mygrate.ErrInitFn,
	}

	if !errors.Is(e, errUnitTest) {
		t.Fatalf("expected mygrate.Error to be '%s'", errUnitTest)
	}
	if !errors.Is(e, mygrate.ErrInitFn) {
		t.Fatalf("expected mygrate.Error to be '%s'", mygrate.ErrInitFn)
	}
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	e := mygrate.Error{
		Err:         errUnitTest,
		InternalErr: mygrate.ErrInitFn,
	}

	unwrapped := e.Unwrap()
	if !errors.Is(unwrapped, errUnitTest) {
		t.Fatalf("expected unwrapped error to be '%s'", errUnitTest)
	}
}
