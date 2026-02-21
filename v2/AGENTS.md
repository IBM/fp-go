# Agent Guidelines for fp-go/v2

This document provides guidelines for AI agents working on the fp-go/v2 project.

## Documentation Standards

### Go Doc Comments

1. **Use Standard Go Doc Format**
   - Do NOT use markdown-style links like `[text](url)`
   - Use simple type references: `ReaderResult`, `Validate[I, A]`, `validation.Success`
   - Go's documentation system will automatically create links

2. **Structure**
   ```go
   // FunctionName does something useful.
   //
   // Longer description explaining the purpose and behavior.
   //
   // # Type Parameters
   //
   //   - T: Description of type parameter
   //
   // # Parameters
   //
   //   - param: Description of parameter
   //
   // # Returns
   //
   //   - ReturnType: Description of return value
   //
   // # Example Usage
   //
   //   code example here
   //
   // # See Also
   //
   //   - RelatedFunction: Brief description
   func FunctionName[T any](param T) ReturnType {
   ```

3. **Code Examples**
   - Use idiomatic Go patterns
   - Prefer `result.Eitherize1(strconv.Atoi)` over manual error handling
   - Show realistic, runnable examples

### File Headers

Always include the Apache 2.0 license header:

```go
// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
```

## Testing Standards

### Test Structure

1. **Organize Tests by Category**
   ```go
   func TestFunctionName_Success(t *testing.T) {
       t.Run("specific success case", func(t *testing.T) {
           // test code
       })
   }
   
   func TestFunctionName_Failure(t *testing.T) {
       t.Run("specific failure case", func(t *testing.T) {
           // test code
       })
   }
   
   func TestFunctionName_EdgeCases(t *testing.T) {
       // edge case tests
   }
   
   func TestFunctionName_Integration(t *testing.T) {
       // integration tests
   }
   ```

2. **Use Direct Assertions**
   - Prefer: `assert.Equal(t, validation.Success(expected), actual)`
   - Avoid: Verbose `either.MonadFold` patterns unless necessary
   - Exception: When you need to verify pointer is not nil or extract specific fields

3. **Use Idiomatic Patterns**
   - Use `result.Eitherize1` for converting `(T, error)` functions
   - Use `result.Of` for success values
   - Use `result.Left` for error values

### Test Coverage

Include tests for:
- **Success cases**: Normal operation with various input types
- **Failure cases**: Error handling and error preservation
- **Edge cases**: Nil, empty, zero values, boundary conditions
- **Integration**: Composition with other functions
- **Type safety**: Verify type parameters work correctly
- **Benchmarks**: Performance-critical paths

### Example Test Pattern

```go
func TestFromReaderResult_Success(t *testing.T) {
    t.Run("converts successful ReaderResult", func(t *testing.T) {
        // Arrange
        parseIntRR := result.Eitherize1(strconv.Atoi)
        validator := FromReaderResult[string, int](parseIntRR)
        
        // Act
        result := validator("42")(nil)
        
        // Assert
        assert.Equal(t, validation.Success(42), result)
    })
}
```

## Code Style

### Functional Patterns

1. **Prefer Composition**
   ```go
   validator := F.Pipe1(
       FromReaderResult[string, int](parseIntRR),
       Chain(validatePositive),
   )
   ```

2. **Use Type-Safe Helpers**
   - `result.Eitherize1` for `func(T) (R, error)`
   - `result.Of` for wrapping success values
   - `result.Left` for wrapping errors

3. **Avoid Verbose Patterns**
   - Don't manually handle `(value, error)` tuples when helpers exist
   - Don't use `either.MonadFold` in tests unless necessary

### Error Handling

1. **In Production Code**
   - Use `validation.Success` for successful validations
   - Use `validation.FailureWithMessage` for simple failures
   - Use `validation.FailureWithError` to preserve error causes

2. **In Tests**
   - Verify error messages and causes
   - Check error context is preserved
   - Test error accumulation when applicable

## Common Patterns

### Converting Error-Based Functions

```go
// Good: Use Eitherize1
parseIntRR := result.Eitherize1(strconv.Atoi)

// Avoid: Manual error handling
parseIntRR := func(input string) result.Result[int] {
    val, err := strconv.Atoi(input)
    if err != nil {
        return result.Left[int](err)
    }
    return result.Of(val)
}
```

### Testing Validation Results

```go
// Good: Direct comparison
assert.Equal(t, validation.Success(42), result)

// Avoid: Verbose extraction (unless you need to verify specific fields)
assert.True(t, either.IsRight(result))
value := either.MonadFold(result,
    func(Errors) int { return 0 },
    F.Identity[int],
)
assert.Equal(t, 42, value)
```

### Documentation Examples

```go
// Good: Concise and idiomatic
//  parseIntRR := result.Eitherize1(strconv.Atoi)
//  validator := FromReaderResult[string, int](parseIntRR)

// Avoid: Verbose manual patterns
//  parseIntRR := func(input string) result.Result[int] {
//      val, err := strconv.Atoi(input)
//      if err != nil {
//          return result.Left[int](err)
//      }
//      return result.Of(val)
//  }
```

## Checklist for New Code

- [ ] Apache 2.0 license header included
- [ ] Go doc comments use standard format (no markdown links)
- [ ] Code examples are idiomatic and concise
- [ ] Tests cover success, failure, edge cases, and integration
- [ ] Tests use direct assertions where possible
- [ ] Benchmarks included for performance-critical code
- [ ] All tests pass
- [ ] Code uses functional composition patterns
- [ ] Error handling preserves context and causes