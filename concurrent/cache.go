// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package concurrent

import "sync"

// Cache provides a very simple in-memory cache which
// is based simply on a map and [sync.Mutex].
type Cache[K comparable, V any] struct {
	initOnce sync.Once
	mu       sync.Mutex
	data     map[K]V
}

func (c *Cache[K, V]) init() {
	c.initOnce.Do(func() {
		c.data = make(map[K]V)
	})
}

// Get retrieves the value for the given key. If the key does not exist
// in the cache, then the zero value for V will be returned along with
// false.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.init()

	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.data[key]
	return v, ok
}

// GetOrNew retrieves the value for the given key. If the key does not exist
// in the cache, then the given function will be called to get the value.
// If the given function succeeds, then the returned value will be placed
// in the cache before being returned.
func (c *Cache[K, V]) GetOrNew(key K, f func() (V, error)) (V, error) {
	c.init()

	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.data[key]
	if ok {
		return v, nil
	}

	v, err := f()
	if err != nil {
		return v, err
	}

	c.data[key] = v
	return v, nil
}

// Put places the key value pair into the cache. It will override any previously
// cached value.
func (c *Cache[K, V]) Put(key K, value V) {
	c.init()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}
