[![Go](https://github.com/hadv/ringcache/actions/workflows/go.yml/badge.svg)](https://github.com/hadv/ringcache/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A non-thread-safe ring cache, also known as a circular buffer, is a fixed-size data structure designed to efficiently handle a continuous stream of data by overwriting the oldest entries when new data is added and the buffer is full. This implementation ensures constant memory usage and provides O(1) time complexity for both adding and retrieving elements. However, it is not safe for concurrent use and should be used in single-threaded environments where predictable performance and efficient memory management are required.

Here's a complete usage example:

```go
package main

import (
	"fmt"
	"ringcache"
)

// evictionHandler is a callback function that gets called when an item is evicted from the cache.
func evictionHandler(key interface{}, value interface{}) {
	fmt.Printf("Evicted key: %v, value: %v\n", key, value)
}

func main() {
	// Create a new ring cache with a maximum size of 3 and an eviction callback
	cache, err := ringcache.NewWithEvict(3, evictionHandler)
	if err != nil {
		fmt.Println("Error creating cache:", err)
		return
	}

	// Add some items to the cache
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	// Print the current cache contents
	fmt.Println("Cache after adding 3 items:")
	for _, key := range []string{"key1", "key2", "key3"} {
		if value, ok := cache.Get(key); ok {
			fmt.Printf("%s: %v\n", key, value)
		} else {
			fmt.Printf("%s not found\n", key)
		}
	}

	// Add another item, causing the first item to be evicted
	cache.Add("key4", "value4")

	// Print the current cache contents again
	fmt.Println("\nCache after adding 4th item (key1 should be evicted):")
	for _, key := range []string{"key1", "key2", "key3", "key4"} {
		if value, ok := cache.Get(key); ok {
			fmt.Printf("%s: %v\n", key, value)
		} else {
			fmt.Printf("%s not found\n", key)
		}
	}

	// Remove an item from the cache
	cache.Remove("key3")
	fmt.Println("\nCache after removing key3:")
	for _, key := range []string{"key2", "key3", "key4"} {
		if value, ok := cache.Get(key); ok {
			fmt.Printf("%s: %v\n", key, value)
		} else {
			fmt.Printf("%s not found\n", key)
		}
	}

	// Clear the cache
	cache.Purge()
	fmt.Println("\nCache after purge:")
	for _, key := range []string{"key2", "key3", "key4"} {
		if value, ok := cache.Get(key); ok {
			fmt.Printf("%s: %v\n", key, value)
		} else {
			fmt.Printf("%s not found\n", key)
		}
	}
}
```

Output:
```
Cache after adding 3 items:
key1: value1
key2: value2
key3: value3

Evicted key: key1, value: value1
Cache after adding 4th item (key1 should be evicted):
key1 not found
key2: value2
key3: value3
key4: value4

Evicted key: key3, value: value3
Cache after removing key3:
key2: value2
key3 not found
key4: value4

Evicted key: key2, value: value2
Evicted key: key4, value: value4
Cache after purge:
key2 not found
key3 not found
key4 not found
```

License
-------
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Contributing
------------
Contributions are welcome! Please feel free to submit a pull request or open an issue if you encounter any problems or have any suggestions for improvements.
