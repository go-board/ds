# Contributing Guide

Thanks for contributing to `go-board/ds`.

## Development style (repository-wide)

### Language and formatting
- Use English for all public comments, docs, and test descriptions.
- Keep comments concise and behavior-focused.
- Run `gofmt` on all changed Go files before submitting.

### API naming conventions
- Ordered traversal APIs must use explicit direction suffixes:
  - `IterAsc` / `IterDesc`
  - `RangeAsc` / `RangeDesc`
- Mutable iterator variants should use `Mut` in the same position consistently:
  - `IterMutAsc`, `RangeMutDesc`, etc.

### Range APIs
- Prefer `bound.RangeBounds[T]` for all ordered range queries.
- Do not introduce ad-hoc `lower/upper *T` public signatures in new APIs.

### Entry APIs (map-like structures)
- Keep Entry behavior aligned across implementations:
  - `Insert(value) (old V, existed bool)`
  - `Get() (V, bool)`
  - `Delete() bool`

### Tests
- Add/adjust tests for any public API change.
- Prefer table-driven tests when multiple edge-cases share setup.
- Keep assertion messages explicit and action-oriented.

## Validation checklist
- `gofmt -w <changed_files>`
- `go test ./...`
