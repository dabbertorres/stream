# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`stream` is a single-package Go generics library (module `github.com/dabbertorres/stream/v2`) implementing lazily-evaluated streams, similar in spirit to Java Streams or Rust iterators. The public API is built around `Seq[T]`, a defined type over the standard library's `iter.Seq[T]`. (An earlier `Stream[T]` API, built on a hand-rolled `streamer[T]` interface, has been removed in favor of `Seq[T]` now that Go supports native range-over-func iterators and generic methods — if you see references to `Stream[T]`/`From*` constructors/`Transform`/`ForEach` in older docs or history, they no longer exist.) This repo requires Go 1.27+ (per `go.mod`), which added support for generic methods (methods with their own type parameters, e.g. `func (s Seq[T]) Map[O any](...)`) — most of `Seq[T]`'s operations depend on this.

## Commands

```sh
go build ./...     # build
go test ./...      # run all tests
go test -run TestName ./...       # run a single test
go test -run TestName/subtest ./...  # run a table-driven subtest
go vet ./...
go tool gofumpt -l .    # list files needing formatting
go tool gofumpt -w .    # format in place
```

Use `gofumpt` (via `go tool gofumpt`), not plain `gofmt`, for formatting — it's registered as a tool dependency in `go.mod`.

CI (`.github/workflows/go.yml`) runs on push/PR to `master`, using Go 1.27 (matching `go.mod`): a `fmt` job runs `go tool gofumpt -l .` and fails if it lists any file, and a `build` job runs `go build -v ./...` and `go test -v ./...` across a matrix of build tags (`""` and `noseqdistinctunsafe`) so both configurations are checked.

Tests use `testify` (`assert`) and Go's built-in example tests (`Example*` functions with `// Output:` comments, e.g. in `README.md`/`example_test.go`).

The `Seq[T].DistinctUnsafe()` method (in `distinct_unsafe.go`) can be excluded from a build with `-tags noseqdistinctunsafe`, for callers that don't want `unsafe` linked in — `go build -tags noseqdistinctunsafe ./...` / `go test -tags noseqdistinctunsafe ./...` should be checked alongside the default build when touching that file or its free-function counterpart in `distinct.go`.

## Architecture

- `Seq[T]` (`type Seq[T any] iter.Seq[T]`, in `seq.go`) is a direct defined type over the standard library range-over-func iterator type. Because it *is* an `iter.Seq[T]`, a value can be ranged over directly (`for v := range seq { ... }`) — there's no `ForEach`/`Range` equivalent; that's just what ranging over the value already gives you.
- Most operations are **generic methods** on `Seq[T]` (e.g. `func (s Seq[T]) Map[O any](mapper func(T) O) Seq[O]`), each in its own file named after the operation (or a closely related pair): `map.go` (`Map`), `filter.go` (`Filter`), `limit.go` (`Limit`), `skip.go` (`Skip`), `reduce.go` (`Reduce`), `associate.go` (`Associate`, plus `ToMap`, see below), `while.go` (`DropWhile`/`TakeWhile`), `flat.go` (`FlatMap`), `first.go` (`First`, `FirstWhere`), `terminal.go` (`ToSlice`, `Append`, `All`, `Any`, `None`, `ToChan`, `Max`, `Min`, `MinMax` — see below). `seq.go` itself is left holding only the `Seq[T]`/`Pair[K, V]` type definitions and the three most basic constructors (`Of`, `OfChan`, `OfKeyed`) — everything else gets pulled out into its own file. When adding a new intermediate or terminal operation, follow this one-file-per-op, fluent-generic-method pattern.
  - **Exception — `Join`** (`join.go`): a free function (`func Join[T any](seqs ...Seq[T]) Seq[T]`), not a method, matching how the standard library always shapes `Join` (`strings.Join`, `path.Join`, `filepath.Join`) and concatenation (`slices.Concat`) as free functions rather than methods — concatenating N sequences is symmetric, with no single one of them more "primary" than the others. (It used to be a method solely to avoid colliding with the old `Stream[T]`'s free-function `Join`; now that `Stream[T]` is gone, there's no reason not to use the more idiomatic shape.)
  - **Exception — `ToMap`** (`associate.go`, alongside the `Associate` method): a free function (`func ToMap[K comparable, V any](Seq[Pair[K, V]]) map[K]V`), because it needs `K comparable`, a stricter constraint than `Pair[K, V]`'s own `K any` — a generic method can't add a constraint to a receiver type parameter that's already bound by the type it's instantiated with, so this has to live outside the method set. Duplicate keys are last-write-wins.
  - **Exception — `Distinct`**: exists in two forms, in separate files, with deliberately different names so the safe one is the one that reads naturally:
    - `distinct.go`: free function `Distinct[T comparable](Seq[T]) Seq[T]`, using a real `map[T]struct{}` — the correct, recommended version, and the only thing named plain `Distinct`.
    - `distinct_unsafe.go`: fluent method `func (s Seq[T]) DistinctUnsafe() Seq[T]`, using `unsafe`/`hash/maphash` to hash each element's raw bytes so it works without `T: comparable` — unsound for types whose equality isn't determined by their bit pattern (strings/slices with different backing arrays, types with padding, interfaces, maps). It was originally named `Distinct` too (method vs. free function, same as the old `Stream[T]`/`Seq[T]` split elsewhere), but that made the unsound version the one a caller reaches for by default when chaining — the `Unsafe` suffix is deliberate, not decorative. Gated behind `//go:build !noseqdistinctunsafe`; its test (`distinct_unsafe_test.go`) carries the same build tag.
- Constructors, all named `Of*`, live in `seq.go` (`Of` wraps an `iter.Seq[T]`, `OfChan` a receive-only channel, `OfKeyed` an `iter.Seq2[K, V]` into `Seq[Pair[K, V]]`) and `source.go` (`OfSlice`, `OfMap`, `OfOptional`, `OfFunc`, `OfDecoder`, plus the `Decoder` interface `OfDecoder` reads from). Several are thin wrappers around stdlib `slices`/`maps` iterator helpers (`OfSlice` → `slices.Values`, `OfMap` → `maps.All` + `OfKeyed`) rather than reimplementing iteration.
- `Pair[K, V]` (`seq.go`) is the package's 2-tuple type, with no constraint on `K` — `ToMap` is where a `comparable` requirement actually gets applied, not the type itself.
- Terminal ops lean on the standard library wherever `Seq[T]`'s `iter.Seq[T]`-compatibility makes that possible instead of reimplementing: `ToSlice` (`slices.Collect`), `Append` (`slices.AppendSeq`).
- There is deliberately no `Sorted` on `Seq[T]`: sorting requires fully draining the source before producing anything, which can't be done inside a lazily-returned `Seq` without either violating the laziness contract every other op honors (eagerly running at the point `.Sorted(...)` is called, not when the result is later ranged over) or buffering internally, which is exactly what `Optional[T]`-returning ops like `Max`/`Min` avoid by staying single-pass. Sorting is a post-collection concern: call `.ToSlice()` and use the standard library's `slices.SortFunc`/`slices.SortedFunc` directly on the result.
- `Optional[T]` (`optional.go`) is the result type for operations that may not produce a value (`First`, `FirstWhere`, `Max`, `Min`, `MinMax`): `Some`/`None`/`OptionalFromPointer` construct one, `Get`/`MustGet`/`GetOrDefault(Func)`/`IfSome`/`IfNone` consume one, and it implements `json.Marshaler`/`Unmarshaler` (`null` for `None`). `optional.go` holds only `Optional[T]` itself, not the `Seq[T]` methods that return one — those live with the rest of their kind (`first.go`, `terminal.go`), per the one-file-per-op rule above. `MinMax` returns two named `Optional[T]` results (`min, max Optional[T]`) rather than a single `Optional` wrapping a struct — there's no `MinMax[T]` type.
- `Max`/`Min`/`MinMax` take their ordering function as a plain `func(lhs, rhs T) bool` literal rather than a named type — there's no `LessFunc[T]`. For any `cmp.Ordered` type, the standard library's `cmp.Less` already satisfies it directly (see tests) — there's no package-local `OrderedLess`/numeric-constraint hierarchy reimplementing that either.
- Every op that can short-circuit (`Limit`, `TakeWhile`, `First`, `FirstWhere`, `Any`, etc.) has a companion test asserting it stops *pulling* from its parent at the right point, not just that its output is correct — see `TestSeqLimitStopsPulling` and `TestSeqTakeWhileStopsPulling` for the pattern. This matters because a subtly wrong index/order in a hand-written `for v := range s { ... }` loop can produce correct output while still over- or under-consuming the parent (this is exactly how a couple of real `Limit`/`Skip` bugs were caught during development — see git history).

## Conventions when adding new operations

- New `Seq[T]` operations are fluent generic methods (see above), not free functions, unless the operation needs a type-parameter constraint stricter than what the receiver's own type already provides (as with `ToMap`) or Go convention itself favors a free function for that shape of operation (as with `Join`).
- Test naming mirrors this: a method's test is `TestSeq<Op>` (e.g. `TestSeqMap`), a free function's test is just `Test<Op>` (e.g. `TestDistinct`, `TestJoin`).
- Prefer a standard-library type/function over a hand-rolled one when one exists and fits (e.g. `cmp.Less`/`cmp.Ordered` used directly at call sites instead of a package-local `OrderedLess`, and no reimplemented `Number`/`Integer`/`Ordered` constraint hierarchy) — don't reintroduce a general-purpose `util.go` grab-bag.
- Prefer a plain function-type literal (e.g. `func(lhs, rhs T) bool`) over introducing a named type for it, unless the named type earns its keep (e.g. by carrying its own methods) — there's no `LessFunc[T]`.
- Keep one operation (or a closely related pair, like `DropWhile`/`TakeWhile` in `while.go`, or the paired `distinct.go`/`distinct_unsafe.go`) per file, named after the operation — no `seq_` prefix, even though every file in this package is part of the `Seq[T]` API (there's nothing else to disambiguate from anymore).
- `Seq[T]`'s zero value is a nil func; ranging over it panics, same as calling any nil func. This is accepted as normal Go behavior for a func-shaped type — no nil guard is added in any method.
