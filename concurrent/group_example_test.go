// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package concurrent

import (
	"context"
	"fmt"
)

func ExampleLazyGroup() {
	var computeGroup LazyGroup
	partialSumCh := make(chan int)
	for i := range 10 {
		computeGroup.Go(func(ctx context.Context) error {
			var total int
			for j := range 100000 {
				total += (j + 1) + i*100000
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case partialSumCh <- total:
			}
			return nil
		})
	}

	var sumGroup LazyGroup
	sumGroup.Go(func(ctx context.Context) error {
		defer close(partialSumCh)
		return computeGroup.Wait(ctx)
	})
	sumGroup.Go(func(ctx context.Context) error {
		var total int
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case partialSum, ok := <-partialSumCh:
				if !ok {
					fmt.Println(total)
					return nil
				}

				total += partialSum
			}
		}
	})

	err := sumGroup.Wait(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output: 500000500000
}
