// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package concurrent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/z5labs/sdk-go/try"
)

func TestLazyGroup_Wait(t *testing.T) {
	t.Run("will return an error", func(t *testing.T) {
		t.Run("if the context given to Wait is cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			var lg LazyGroup

			for range 5 {
				lg.Go(func(ctx context.Context) error {
					defer cancel()

					<-time.After(100 * time.Millisecond)

					return nil
				})
			}

			err := lg.Wait(ctx)
			require.ErrorIs(t, err, context.Canceled)
		})

		t.Run("if one of the functions returns a error", func(t *testing.T) {
			var lg LazyGroup

			errFuncFailed := errors.New("failed")
			lg.Go(func(ctx context.Context) error {
				return errFuncFailed
			})

			err := lg.Wait(t.Context())
			require.ErrorIs(t, err, errFuncFailed)
		})

		t.Run("if one of the functions panics", func(t *testing.T) {
			var lg LazyGroup

			lg.Go(func(ctx context.Context) error {
				panic("hello")
			})

			err := lg.Wait(t.Context())

			var perr try.PanicError
			require.ErrorAs(t, err, &perr)
			require.Equal(t, "hello", perr.Value)
		})

		t.Run("if multiple functions return an error", func(t *testing.T) {
			var lg LazyGroup

			errFirst := errors.New("first goroutine failed")
			lg.Go(func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(100 * time.Millisecond):
				}
				return errFirst
			})

			lg.Go(func(ctx context.Context) error {
				<-ctx.Done()
				return context.Canceled
			})

			err := lg.Wait(t.Context())
			require.ErrorIs(t, err, errFirst)
		})
	})
}
