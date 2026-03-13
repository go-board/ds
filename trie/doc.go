// Package trie implements a generic trie (prefix tree) data structure.
//
// A trie (also known as a prefix tree or digital tree) is a tree data structure
// used for efficient retrieval of sequences. It's particularly optimized for
// string-based operations like autocomplete, spell checking, and IP routing.
//
// # Time Complexity
//
//   - Insert: O(k) where k is the key length
//   - Delete: O(k)
//   - Get: O(k)
//   - PrefixSearch: O(p + m) where p is prefix length, m is number of matches
//   - KeysWithPrefix: O(p + m) where p is prefix length, m is number of matches
//   - LongestPrefix: O(k)
//
// # Features
//
//   - Efficient prefix-based operations
//   - O(k) time complexity for basic operations where k is key length
//   - Support for storing key-value pairs where keys are sequences
//   - Prefix search and autocomplete support
//   - Longest prefix matching
//   - Generic key element type support
//
// # Usage
//
// Trie excels at string/sequence operations:
//
//	// Keys are sequences of elements (e.g., []rune, []byte, []string)
//	m := trie.NewOrdered[[]rune, int]()
//
//	// Insert words
//	m.Insert([]rune("apple"), 1)
//	m.Insert([]rune("app"), 2)
//	m.Insert([]rune("application"), 3)
//	m.Insert([]rune("banana"), 4)
//
//	// Get exact match
//	val, found := m.Get([]rune("apple"))
//	fmt.Println(val, found) // 1 true
//
//	// Prefix search (autocomplete)
//	for key, val := range m.PrefixSearch([]rune("app")) {
//	    fmt.Println(string(key), val)
//	    // Output: app 2, apple 1, application 3
//	}
//
//	// Longest prefix match
//	prefix, val, found := m.LongestPrefix([]rune("applet"))
//	fmt.Println(string(prefix), val) // "apple" 1
//
// # Use Cases
//
//   - Autocomplete and type-ahead suggestions
//   - Spell checking
//   - IP routing (longest prefix match)
//   - Word frequency analysis
//   - T9 predictive text
//
// # Key Type
//
// The key type is a slice of elements. Common choices:
//
//	// For strings: []rune or []byte
//	trie.New[[]rune, V]()
//	trie.New[[]byte, V]()
//
//	// For custom sequences
//	trie.New[[]int, V]()
//	trie.New[[]CustomType, V]()
package trie
