// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package try

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type closeFunc func() error

func (f closeFunc) Close() error {
	return f()
}

func TestClose(t *testing.T) {
	t.Run("will update the error ref value", func(t *testing.T) {
		t.Run("if the close fails and the ref value is nil", func(t *testing.T) {
			closeErr := errors.New("close failed")
			c := closeFunc(func() error {
				return closeErr
			})

			f := func() (err error) {
				defer Close(&err, c)
				return nil
			}

			err := f()
			require.ErrorIs(t, err, closeErr)
		})

		t.Run("if the close fails and the ref value is non-nil", func(t *testing.T) {
			closeErr := errors.New("close failed")
			c := closeFunc(func() error {
				return closeErr
			})

			funcErr := errors.New("func error")
			f := func() (err error) {
				defer Close(&err, c)
				return funcErr
			}

			err := f()
			require.ErrorIs(t, err, funcErr)
			require.ErrorIs(t, err, closeErr)
		})
	})

	t.Run("will change the error ref value", func(t *testing.T) {
		t.Run("if the value is not an io.Closer", func(t *testing.T) {
			funcErr := errors.New("func error")
			f := func() (err error) {
				defer Close(&err, nil)
				return funcErr
			}

			err := f()
			require.ErrorIs(t, err, funcErr)
		})

		t.Run("if Close succeeds", func(t *testing.T) {
			c := closeFunc(func() error {
				return nil
			})

			funcErr := errors.New("func error")
			f := func() (err error) {
				defer Close(&err, c)
				return funcErr
			}

			err := f()
			require.ErrorIs(t, err, funcErr)
		})
	})
}

type flushFunc func() error

func (f flushFunc) Flush() error {
	return f()
}

func TestFlush(t *testing.T) {
	t.Run("will update the error ref value", func(t *testing.T) {
		t.Run("if the flush fails and the ref value is nil", func(t *testing.T) {
			flushErr := errors.New("flush failed")
			flusher := flushFunc(func() error {
				return flushErr
			})

			f := func() (err error) {
				defer Flush(&err, flusher)
				return nil
			}

			err := f()
			require.ErrorIs(t, err, flushErr)
		})

		t.Run("if the close fails and the ref value is non-nil", func(t *testing.T) {
			flushErr := errors.New("flush failed")
			flusher := flushFunc(func() error {
				return flushErr
			})

			funcErr := errors.New("func error")
			f := func() (err error) {
				defer Flush(&err, flusher)
				return funcErr
			}

			err := f()
			require.ErrorIs(t, err, funcErr)
			require.ErrorIs(t, err, flushErr)
		})
	})

	t.Run("will change the error ref value", func(t *testing.T) {
		t.Run("if the value is not an io.Closer", func(t *testing.T) {
			funcErr := errors.New("func error")
			f := func() (err error) {
				defer Flush(&err, nil)
				return funcErr
			}

			err := f()
			require.ErrorIs(t, err, funcErr)
		})

		t.Run("if Close succeeds", func(t *testing.T) {
			flusher := flushFunc(func() error {
				return nil
			})

			funcErr := errors.New("func error")
			f := func() (err error) {
				defer Flush(&err, flusher)
				return funcErr
			}

			err := f()
			require.ErrorIs(t, err, funcErr)
		})
	})
}

func TestRecover(t *testing.T) {
	t.Run("will update the error ref value", func(t *testing.T) {
		t.Run("if a panic is successfully recovered from and the ref is set to nil", func(t *testing.T) {
			f := func() (err error) {
				defer Recover(&err)
				panic("hello world")
			}

			err := f()

			var perr PanicError
			require.ErrorAs(t, err, &perr)
			require.NotEmpty(t, perr.Error())
			require.Equal(t, "hello world", perr.Value)
		})

		t.Run("if a panic is successfully recovered from and the ref is set to a non-nil value", func(t *testing.T) {
			funcErr := errors.New("error value")
			panicErr := errors.New("panic error")
			f := func() (err error) {
				defer Recover(&err)
				err = funcErr
				panic(panicErr)
			}

			err := f()

			require.ErrorIs(t, err, funcErr)

			var perr PanicError
			require.ErrorAs(t, err, &perr)
			require.NotEmpty(t, perr.Error())
			require.ErrorIs(t, perr, panicErr)
		})
	})

	t.Run("will not update the error ref value", func(t *testing.T) {
		t.Run("if no panic is occurred", func(t *testing.T) {
			f := func() (err error) {
				defer Recover(&err)
				return nil
			}

			err := f()
			require.Nil(t, err)
		})
	})
}
