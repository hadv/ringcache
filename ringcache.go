package ringcache

import (
	"errors"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback func(key interface{}, value interface{})

// RingCache, often known as a circular buffer or ring buffer, is a data
// structure that uses a single, fixed-size buffer as if it were connected
// end-to-end. It is particularly useful for applications that require a buffer
// with a consistent and predictable size, such as in real-time data processing
// systems or network packet buffering.
type RingCache struct {
	maxSize int
	next    int
	keys    []interface{}
	items   map[interface{}]interface{}
	onEvict EvictCallback
}

// New creates a ring cache of the given size.
func New(maxSize int) (*RingCache, error) {
	return NewWithEvict(maxSize, nil)
}

// NewWithEvict constructs ring cache of the given size with callback
func NewWithEvict(maxSize int, onEvict EvictCallback) (*RingCache, error) {
	if maxSize <= 0 {
		return nil, errors.New("cache size should be greater than zero")
	}
	cache := &RingCache{
		maxSize: maxSize,
		next:    0,
		keys:    make([]interface{}, maxSize),
		items:   make(map[interface{}]interface{}),
		onEvict: onEvict,
	}

	return cache, nil
}

// Purge is used to completely clear the cache.
func (c *RingCache) Purge() {
	// evict all items
	if c.onEvict != nil {
		for _, k := range c.keys {
			if k != nil {
				c.onEvict(k, c.items[k])
			}
		}
	}

	// re-initialize
	c.items = make(map[interface{}]interface{})
	c.keys = make([]interface{}, c.maxSize)
	c.next = 0
}

// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *RingCache) Add(key, value interface{}) (evicted bool) {
	evicted = false

	// Do nothing if key or value is nil
	if key == nil || value == nil {
		return
	}

	// Check for existing item
	if k := c.keys[c.next]; k != nil {
		if c.onEvict != nil {
			c.onEvict(k, c.items[k])
			evicted = true
		}
		delete(c.items, k)
	}

	c.items[key] = value
	c.keys[c.next] = key
	c.next = (c.next + 1) % c.maxSize

	return
}

// Get looks up a key's value from the cache.
func (c *RingCache) Get(key interface{}) (interface{}, bool) {
	value, ok := c.items[key]

	return value, ok
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *RingCache) Contains(key interface{}) bool {
	_, ok := c.items[key]

	return ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *RingCache) Remove(key interface{}) bool {
	if val, ok := c.items[key]; ok {
		delete(c.items, key)
		for i, k := range c.keys {
			if k == key {
				c.keys[i] = nil
				if c.onEvict != nil {
					c.onEvict(key, val)
				}

				return true
			}
		}
	}

	return false
}

// Len returns the number of items in the cache.
func (c *RingCache) Len() int {
	return len(c.items)
}

// Cap returns the capacity of the cache.
func (c *RingCache) Cap() int {
	return c.maxSize
}
