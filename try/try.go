// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package try provides helpers for often deferred actions.
package try

import (
	"errors"
	"fmt"
	"io"
)

// Close will close the given [io.Closer] and join any error
// with the existing value referenced by the given *error.
func Close(err *error, c io.Closer) {
	if c == nil {
		return
	}

	cerr := c.Close()
	if cerr == nil {
		return
	}
	if *err == nil {
		*err = cerr
		return
	}
	*err = errors.Join(*err, cerr)
}

// Flusher represents anything which can flush itself.
type Flusher interface {
	Flush() error
}

// Flush will flush the given [Flusher] and join any error
// with the existing value referenced by the given *error.
func Flush(err *error, f Flusher) {
	if f == nil {
		return
	}

	ferr := f.Flush()
	if ferr == nil {
		return
	}
	if *err == nil {
		*err = ferr
		return
	}
	*err = errors.Join(*err, ferr)

}

// PanicError represents recovering from a panic and contains any value
// recovered from the panic.
type PanicError struct {
	Value any
}

// Error implements the [error] interface.
func (e PanicError) Error() string {
	return fmt.Sprintf("recovered from panic: %v", e.Value)
}

// Unwrap implements the implicit interface used by [errors.Is] and [errors.As].
func (e PanicError) Unwrap() error {
	return e.Value.(error)
}

// Recover with call [recover] and wrap and recovered any value
// into a [PanicError]. This [PanicError] will then be joined
// with the existing value reference by the given *error.
func Recover(err *error) {
	r := recover()
	if r == nil {
		return
	}

	perr := PanicError{
		Value: r,
	}
	if *err == nil {
		*err = perr
		return
	}
	*err = errors.Join(*err, perr)
}
