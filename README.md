# ds

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

`github.com/go-board/ds` is a generic data-structures library for modern Go.
It provides consistent APIs for maps, sets, trees, linear containers, and iterators,
with explicit ordered traversal and composable range bounds.

---

## Why this project

- **Generic-first design** (Go 1.24+)
- **Consistent API style** across related data structures
- **Ordered & unordered variants** for map/set workloads
- **Explicit iteration direction** (`IterAsc` / `IterDesc`)
- **Composable range model** via `bound.RangeBounds`
- **Rust-inspired Entry workflows** for map updates

---

## Installation

```bash
go get github.com/go-board/ds
```

## Requirements

- Go **1.24** or newer

---

## Package overview

| Category | Implementations |
|---|---|
| Linear containers | `ArrayDeque`, `ArrayStack`, `LinkedList` |
| Heap/queue | `PriorityQueue` |
| Unordered map/set | `HashMap`, `HashSet` |
| Ordered tree map/set | `BTree`, `BTreeMap`, `BTreeSet` |
| Ordered skip-list map/set | `SkipMap`, `SkipSet` |
| Prefix map | `TrieMap` |
| Range abstraction | `bound` (`Bound`, `RangeBounds`) |

---

## Quick start

```go
package main

import (
	"fmt"

	"github.com/go-board/ds"
)

func main() {
	// HashMap (unordered)
	hm := ds.NewComparableHashMap[string, int]()
	hm.Insert("apple", 3)
	hm.Insert("banana", 5)

	if v, ok := hm.Get("banana"); ok {
		fmt.Println("banana:", v)
	}

	// BTreeMap (ordered)
	om := ds.NewOrderedBTreeMap[string, int]()
	om.Insert("a", 1)
	om.Insert("b", 2)
	om.Insert("c", 3)

	for k, v := range om.IterAsc() {
		fmt.Println(k, v)
	}
}
```

---

## Core API concepts

### 1) Directional iteration

Ordered structures use explicit direction in API names:

- `IterAsc()` / `IterDesc()`
- `RangeAsc(bounds)` / `RangeDesc(bounds)`

This avoids ambiguity and improves call-site readability.

### 2) Bounds-based range queries

Use `RangeBounds` to model inclusive/exclusive/unbounded boundaries.

```go
m := ds.NewOrderedSkipMap[int, string]()
m.Insert(10, "a")
m.Insert(20, "b")
m.Insert(30, "c")

bounds := ds.NewRangeBounds(
	ds.NewIncluded(10),
	ds.NewExcluded(30),
)

for k, v := range m.RangeAsc(bounds) {
	fmt.Println(k, v) // 10 a, 20 b
}
```

You can also create fully unbounded ranges:

```go
all := ds.NewRangeBounds(ds.NewUnbounded[int](), ds.NewUnbounded[int]())
_ = all
```

### 3) Entry workflow (map-like types)

Map-like structures expose `Entry` helpers for upsert/delete flows:

```go
m := ds.NewComparableHashMap[string, int]()

m.Entry("x").AndModify(func(v *int) { *v += 1 }).OrInsert(1)

old, existed := m.Entry("x").Insert(10)
fmt.Println(old, existed) // 1 true

removed := m.Entry("x").Delete()
fmt.Println(removed) // true
```

---

## Focused examples

### BTreeSet + range bounds

```go
s := ds.NewOrderedBTreeSet[int]()
for _, v := range []int{1, 2, 3, 4, 5} {
	s.Insert(v)
}

bounds := ds.NewRangeBounds(ds.NewIncluded(2), ds.NewIncluded(4))
for v := range s.RangeAsc(bounds) {
	fmt.Println(v) // 2, 3, 4
}
```

### SkipSet set algebra

```go
a := ds.NewOrderedSkipSet[int]()
b := ds.NewOrderedSkipSet[int]()
for _, v := range []int{1, 2, 3} { a.Insert(v) }
for _, v := range []int{3, 4, 5} { b.Insert(v) }

u := a.Union(b)
for v := range u.IterAsc() {
	fmt.Println(v) // 1 2 3 4 5
}
```

### TrieMap prefix iteration

```go
m := ds.NewOrderedTrieMap[string, int]()
m.Insert([]string{"app"}, 1)
m.Insert([]string{"app", "le"}, 2)
m.Insert([]string{"app", "store"}, 3)

for key, val := range m.IterByPrefix([]string{"app"}) {
	fmt.Println(key, val)
}
```

### PriorityQueue

```go
pq := ds.NewMinOrderedPriorityQueue[int]()
pq.Push(4)
pq.Push(1)
pq.Push(3)

for !pq.IsEmpty() {
	v, _ := pq.Pop()
	fmt.Println(v) // 1, 3, 4
}
```

---

## API consistency notes

- Ordered map APIs (`BTreeMap`, `SkipMap`) are intentionally aligned:
  - directional iterators and ranges (`Asc` / `Desc`)
  - mutable iterator variants for values
  - `Entry` support
- Ordered set APIs (`BTreeSet`, `SkipSet`) expose matching set-algebra operations.
- `HashMap`/`HashSet` keep unordered semantics but follow the same return-style conventions where possible.


## Development style

The repository follows a unified style guide for naming, range APIs, Entry semantics,
and documentation conventions. See [CONTRIBUTING.md](./CONTRIBUTING.md).

---

## Running tests

```bash
go test ./...
```

---

## Contributing

Issues and pull requests are welcome.

Please keep contributions aligned with repository conventions:

1. Add/update tests for behavior changes.
2. Keep naming consistent across equivalent structures.
3. Update docs/examples when public APIs change.
4. Run `gofmt` and `go test ./...` before submitting.

---

## License

MIT. See [LICENSE](./LICENSE).
