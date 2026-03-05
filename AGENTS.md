# AGENTS.md

This file defines repository-wide instructions for coding agents working in this project.

## Scope

These instructions apply to the entire repository.

## Language and comments

- Use English for all public-facing comments, docs, and test descriptions.
- Keep comments concise and behavior-focused.
- Remove stale comments during refactors.

## API style

- Keep naming consistent across equivalent data structures.
- For ordered traversal APIs, use explicit suffixes:
  - `IterAsc` / `IterDesc`
  - `RangeAsc` / `RangeDesc`
- For mutable traversal variants, keep `Mut` placement consistent:
  - `IterMutAsc`, `RangeMutDesc`, etc.

## Range APIs

- Prefer `bound.RangeBounds[T]` for ordered range queries.
- Avoid introducing ad-hoc public signatures like `lower/upper *T` for new range APIs.

## Entry APIs (map-like structures)

Keep `Entry` behavior aligned across implementations where applicable:

- `Insert(value) (old V, existed bool)`
- `Get() (V, bool)`
- `Delete() bool`

## Testing and validation

- Add or update tests for behavior changes.
- Run the following before finalizing:
  - `gofmt -w <changed_files>`
  - `go test ./...`

## Documentation

- Keep `README.md` and `CONTRIBUTING.md` aligned with public API changes.
- Prefer examples that compile against current APIs.
