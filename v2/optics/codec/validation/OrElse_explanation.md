# OrElse is Equivalent to ChainLeft

## Overview

In [`optics/codec/validation/monad.go`](monad.go:474-476), the [`OrElse`](monad.go:474) function is defined as a simple alias for [`ChainLeft`](monad.go:304):

```go
//go:inline
func OrElse[A any](f Kleisli[Errors, A]) Operator[A, A] {
	return ChainLeft(f)
}
```

This means **`OrElse` and `ChainLeft` are functionally identical** - they produce exactly the same results for all inputs.

## Why Have Both?

While they are technically the same, they serve different **semantic purposes**:

### ChainLeft - Technical Perspective
[`ChainLeft`](monad.go:304-309) emphasizes the **technical operation**: it chains a computation on the Left (failure) channel of the Either/Validation monad. This name comes from category theory and functional programming terminology.

### OrElse - Semantic Perspective  
[`OrElse`](monad.go:474-476) emphasizes the **intent**: it provides an alternative or fallback when validation fails. The name reads naturally in code: "try this validation, **or else** try this alternative."

## Key Behavior

Both functions share the same critical behavior that distinguishes them from standard Either operations:

### Error Aggregation
When the transformation function returns a failure, **both the original errors AND the new errors are combined** using the Errors monoid. This ensures no validation errors are lost.

```go
// Example: Error aggregation
result := OrElse(func(errs Errors) Validation[string] {
    return Failures[string](Errors{
        &ValidationError{Messsage: "additional error"},
    })
})(Failures[string](Errors{
    &ValidationError{Messsage: "original error"},
}))

// Result contains BOTH errors: ["original error", "additional error"]
```

### Success Pass-Through
Success values pass through unchanged - the function is never called:

```go
result := OrElse(func(errs Errors) Validation[int] {
    return Failures[int](Errors{
        &ValidationError{Messsage: "never called"},
    })
})(Success(42))

// Result: Success(42) - unchanged
```

### Error Recovery
The function can recover from failures by returning a Success:

```go
recoverFromNotFound := OrElse(func(errs Errors) Validation[int] {
    for _, err := range errs {
        if err.Messsage == "not found" {
            return Success(0) // recover with default
        }
    }
    return Failures[int](errs)
})

result := recoverFromNotFound(Failures[int](Errors{
    &ValidationError{Messsage: "not found"},
}))

// Result: Success(0) - recovered from failure
```

## Use Cases

### 1. Fallback Validation (OrElse reads better)
```go
validatePositive := func(x int) Validation[int] {
    if x > 0 {
        return Success(x)
    }
    return Failures[int](Errors{
        &ValidationError{Messsage: "must be positive"},
    })
}

// Use OrElse for semantic clarity
withDefault := OrElse(func(errs Errors) Validation[int] {
    return Success(1) // default to 1 if validation fails
})

result := F.Pipe1(validatePositive(-5), withDefault)
// Result: Success(1)
```

### 2. Error Context Addition (ChainLeft reads better)
```go
addContext := ChainLeft(func(errs Errors) Validation[string] {
    return Failures[string](Errors{
        &ValidationError{
            Messsage: "validation failed in user.email field",
        },
    })
})

result := F.Pipe1(
    Failures[string](Errors{
        &ValidationError{Messsage: "invalid format"},
    }),
    addContext,
)
// Result contains: ["invalid format", "validation failed in user.email field"]
```

### 3. Pipeline Composition
Both can be used in pipelines, with errors accumulating at each step:

```go
result := F.Pipe2(
    Failures[int](Errors{
        &ValidationError{Messsage: "database error"},
    }),
    OrElse(func(errs Errors) Validation[int] {
        return Failures[int](Errors{
            &ValidationError{Messsage: "context added"},
        })
    }),
    OrElse(func(errs Errors) Validation[int] {
        return Failures[int](errs) // propagate
    }),
)
// Errors accumulate at each step in the pipeline
```

## Verification

The test suite in [`monad_test.go`](monad_test.go:1698) includes comprehensive tests proving that `OrElse` and `ChainLeft` are equivalent:

- ✅ Identical behavior for Success values
- ✅ Identical behavior for error recovery
- ✅ Identical behavior for error aggregation
- ✅ Identical behavior in pipeline composition
- ✅ Identical behavior for multiple error scenarios

Run the tests:
```bash
go test -v -run TestOrElse ./optics/codec/validation
```

## Conclusion

**`OrElse` is exactly the same as `ChainLeft`** - they are aliases with identical implementations and behavior. The choice between them is purely about **code readability and semantic intent**:

- Use **`OrElse`** when emphasizing fallback/alternative validation logic
- Use **`ChainLeft`** when emphasizing technical error channel transformation

Both maintain the critical validation property of **error aggregation**, ensuring all validation failures are preserved and reported together.