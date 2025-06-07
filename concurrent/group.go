// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package concurrent

import (
	"context"
	"sync"

	"github.com/z5labs/sdk-go/try"
)

// LazyGroup is a collection of goroutines. The goroutines are not created until
// [LazyGroup.Wait] is called. A zero LazyGroup is valid, has no limit on the number
// of active goroutines.
type LazyGroup struct {
	funcs []func(context.Context) error
}

// Go registers the given func to be ran on its own goroutine
// once [LazyGroup.Wait] is called.
func (g *LazyGroup) Go(f func(context.Context) error) {
	g.funcs = append(g.funcs, func(ctx context.Context) (err error) {
		defer try.Recover(&err)

		return f(ctx)
	})
}

// Wait runs all registered funcs in their own goroutines and waits for all
// of them to complete. Wait will returns the first error to be returned by
// any of the funcs which returned a non-nil error.
func (g *LazyGroup) Wait(ctx context.Context) error {
	groupCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	var wg sync.WaitGroup
	for _, f := range g.funcs {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := f(groupCtx)
			if err == nil {
				return
			}

			cancel(err)
		}()
	}

	wg.Wait()
	return context.Cause(groupCtx)
}
