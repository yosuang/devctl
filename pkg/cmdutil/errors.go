package cmdutil

import (
	"errors"
	"fmt"
)

func FlagErrorf(format string, args ...interface{}) error {
	return FlagErrorWrap(fmt.Errorf(format, args...))
}

func FlagErrorWrap(err error) error { return &FlagError{err} }

type FlagError struct {
	// Note: not struct{error}: only *FlagError should satisfy error.
	err error
}

func (fe *FlagError) Error() string {
	return fe.err.Error()
}

var ErrSilent = errors.New("SilentError")
