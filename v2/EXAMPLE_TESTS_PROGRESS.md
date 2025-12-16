# Example Tests Progress

This document tracks the progress of converting documentation examples into executable example test files.

## Overview

The codebase has 300+ documentation examples across many packages. This document tracks which packages have been completed and which still need work.

## Completed Packages

### Core Packages
- [x] **result** - Created `examples_bind_test.go`, `examples_curry_test.go`, `examples_apply_test.go`
  - Files: `bind.go` (10 examples), `curry.go` (5 examples), `apply.go` (2 examples)
  - Status: ✅ 17 tests passing

### Utility Packages
- [x] **pair** - Created `examples_test.go`
  - Files: `pair.go` (14 examples)
  - Status: ✅ 14 tests passing

- [x] **tuple** - Created `examples_test.go`
  - Files: `tuple.go` (6 examples)
  - Status: ✅ 6 tests passing

### Type Class Packages
- [x] **semigroup** - Created `examples_test.go`
  - Files: `semigroup.go` (7 examples)
  - Status: ✅ 7 tests passing

### Utility Packages (continued)
- [x] **predicate** - Created `examples_test.go`
  - Files: `bool.go` (3 examples), `contramap.go` (1 example)
  - Status: ✅ 4 tests passing

### Context Reader Packages
- [x] **idiomatic/context/readerresult** - Created `examples_reader_test.go`, `examples_bind_test.go`
  - Files: `reader.go` (8 examples), `bind.go` (14 examples)
  - Status: ✅ 22 tests passing

## Summary Statistics
- **Total Example Tests Created**: 74
- **Total Packages Completed**: 7 (result, pair, tuple, semigroup, predicate, idiomatic/context/readerresult)
- **All Tests Status**: ✅ PASSING

### Breakdown by Package
- **result**: 21 tests (bind: 10, curry: 5, apply: 2, array: 4)
- **pair**: 14 tests
- **tuple**: 6 tests
- **semigroup**: 7 tests
- **predicate**: 4 tests
- **idiomatic/context/readerresult**: 22 tests (reader: 8, bind: 14)

## Packages with Existing Examples

These packages already have some example test files:
- result (has `examples_create_test.go`, `examples_extract_test.go`)
- option (has `examples_create_test.go`, `examples_extract_test.go`)
- either (has `examples_create_test.go`, `examples_extract_test.go`)
- ioeither (has `examples_create_test.go`, `examples_do_test.go`, `examples_extract_test.go`)
- ioresult (has `examples_create_test.go`, `examples_do_test.go`, `examples_extract_test.go`)
- lazy (has `example_lazy_test.go`)
- array (has `examples_basic_test.go`, `examples_sort_test.go`, `example_any_test.go`, `example_find_test.go`)
- readerioeither (has `traverse_example_test.go`)
- context/readerioresult (has `flip_example_test.go`)

## Packages Needing Example Tests

### Core Packages (High Priority)
- [ ] **result** - Additional files need examples:
  - `apply.go` (2 examples)
  - `array.go` (7 examples)
  - `core.go` (6 examples)
  - `either.go` (26 examples)
  - `eq.go` (2 examples)
  - `functor.go` (1 example)
  
- [ ] **option** - Additional files need examples
- [ ] **either** - Additional files need examples

### Reader Packages (High Priority)
- [ ] **reader** - Many examples in:
  - `array.go` (12 examples)
  - `bind.go` (10 examples)
  - `curry.go` (8 examples)
  - `flip.go` (2 examples)
  - `reader.go` (21 examples)

- [ ] **readeroption** - Examples in:
  - `array.go` (3 examples)
  - `bind.go` (7 examples)
  - `curry.go` (5 examples)
  - `flip.go` (2 examples)
  - `from.go` (4 examples)
  - `reader.go` (18 examples)
  - `sequence.go` (4 examples)

- [ ] **readerresult** - Examples in:
  - `array.go` (3 examples)
  - `bind.go` (24 examples)
  - `curry.go` (7 examples)
  - `flip.go` (2 examples)
  - `from.go` (4 examples)
  - `monoid.go` (3 examples)

- [ ] **readereither** - Examples in:
  - `array.go` (3 examples)
  - `bind.go` (7 examples)
  - `flip.go` (3 examples)

- [ ] **readerio** - Examples in:
  - `array.go` (3 examples)
  - `bind.go` (7 examples)
  - `flip.go` (2 examples)
  - `logging.go` (4 examples)
  - `reader.go` (30 examples)

- [ ] **readerioeither** - Examples in:
  - `bind.go` (7 examples)
  - `flip.go` (1 example)

- [ ] **readerioresult** - Examples in:
  - `array.go` (8 examples)
  - `bind.go` (24 examples)

### State Packages
- [ ] **statereaderioeither** - Examples in:
  - `bind.go` (5 examples)
  - `resource.go` (1 example)
  - `state.go` (13 examples)

### Utility Packages
- [ ] **lazy** - Additional examples in:
  - `apply.go` (2 examples)
  - `bind.go` (7 examples)
  - `lazy.go` (10 examples)
  - `sequence.go` (4 examples)
  - `traverse.go` (2 examples)

- [ ] **pair** - Additional examples in:
  - `monad.go` (12 examples)
  - `pair.go` (remaining ~20 examples)

- [ ] **tuple** - Examples in:
  - `tuple.go` (6 examples)

- [ ] **predicate** - Examples in:
  - `bool.go` (3 examples)
  - `contramap.go` (1 example)
  - `monoid.go` (4 examples)

- [ ] **retry** - Examples in:
  - `retry.go` (7 examples)

- [ ] **logging** - Examples in:
  - `logger.go` (5 examples)

### Collection Packages
- [ ] **record** - Examples in:
  - `bind.go` (3 examples)

### Type Class Packages
- [ ] **semigroup** - Examples in:
  - `alt.go` (1 example)
  - `apply.go` (1 example)
  - `array.go` (4 examples)
  - `semigroup.go` (7 examples)

- [ ] **ord** - Examples in:
  - `ord.go` (1 example)

## Strategy for Completion

1. **Prioritize by usage**: Focus on core packages (result, option, either) first
2. **Group by package**: Complete all examples for one package before moving to next
3. **Test incrementally**: Run tests after each file to catch errors early
4. **Follow patterns**: Use existing example test files as templates
5. **Document as you go**: Update this file with progress

## Example Test File Template

```go
// Copyright header...

package packagename_test

import (
	"fmt"
	PKG "github.com/IBM/fp-go/v2/packagename"
)

func ExampleFunctionName() {
	// Copy example from doc comment
	// Ensure it compiles and produces correct output
	fmt.Println(result)
	// Output:
	// expected output
}
```

## Notes

- Use `F.Constant1[error](defaultValue)` for GetOrElse in result package
- Use `F.Pipe1` instead of `F.Pipe2` when only one transformation
- Check function signatures carefully for type parameters
- Some functions like `BiMap` are capitalized differently than in docs
- **Prefer `R.Eitherize1(func)` over manual error handling** - converts `func(T) (R, error)` to `func(T) Result[R]`
  - Example: Use `R.Eitherize1(strconv.Atoi)` instead of manual if/else error checking
- **Add Go documentation comments to all example functions** - Each example should have a comment explaining what it demonstrates
- **Idiomatic vs Non-Idiomatic packages**:
  - Non-idiomatic (e.g., `result`): Uses `Result[A]` type (Either monad)
  - Idiomatic (e.g., `idiomatic/result`): Uses `(A, error)` tuples (Go-style)
  - Context readers use non-idiomatic `Result[A]` internally