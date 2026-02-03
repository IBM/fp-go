# ChainLeft and OrElse in the Decode Package

## Overview

In [`optics/codec/decode/monad.go`](monad.go:53-62), the [`ChainLeft`](monad.go:53) and [`OrElse`](monad.go:60) functions work with decoders that may fail during decoding operations.

```go
func ChainLeft[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return readert.Chain[Decode[I, A]](
		validation.ChainLeft,
		f,
	)
}

func OrElse[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return ChainLeft(f)
}
```

## Key Insight: OrElse is ChainLeft

**`OrElse` is exactly the same as `ChainLeft`** - they are aliases with identical implementations and behavior. The choice between them is purely about **code readability and semantic intent**.

## Understanding the Types

### Decode[I, A]
A decoder that takes input of type `I` and produces a `Validation[A]`:
```go
type Decode[I, A any] = func(I) Validation[A]
```

### Kleisli[I, Errors, A]
A function that takes `Errors` and produces a `Decode[I, A]`:
```go
type Kleisli[I, Errors, A] = func(Errors) Decode[I, A]
```

This allows error handlers to:
1. Access the validation errors that occurred
2. Access the original input (via the returned Decode function)
3. Either recover with a success value or produce new errors

### Operator[I, A, A]
A function that transforms one decoder into another:
```go
type Operator[I, A, A] = func(Decode[I, A]) Decode[I, A]
```

## Core Behavior

Both [`ChainLeft`](monad.go:53) and [`OrElse`](monad.go:60) delegate to [`validation.ChainLeft`](../validation/monad.go:304), which provides:

### 1. Error Aggregation
When the transformation function returns a failure, **both the original errors AND the new errors are combined** using the Errors monoid:

```go
failingDecoder := func(input string) Validation[int] {
    return either.Left[int](validation.Errors{
        {Value: input, Messsage: "original error"},
    })
}

handler := ChainLeft(func(errs Errors) Decode[string, int] {
    return func(input string) Validation[int] {
        return either.Left[int](validation.Errors{
            {Messsage: "additional error"},
        })
    }
})

decoder := handler(failingDecoder)
result := decoder("input")
// Result contains BOTH errors: ["original error", "additional error"]
```

### 2. Success Pass-Through
Success values pass through unchanged - the handler is never called:

```go
successDecoder := Of[string](42)

handler := ChainLeft(func(errs Errors) Decode[string, int] {
    return func(input string) Validation[int] {
        return either.Left[int](validation.Errors{
            {Messsage: "never called"},
        })
    }
})

decoder := handler(successDecoder)
result := decoder("input")
// Result: Success(42) - unchanged
```

### 3. Error Recovery
The handler can recover from failures by returning a successful decoder:

```go
failingDecoder := func(input string) Validation[int] {
    return either.Left[int](validation.Errors{
        {Value: input, Messsage: "not found"},
    })
}

recoverFromNotFound := ChainLeft(func(errs Errors) Decode[string, int] {
    for _, err := range errs {
        if err.Messsage == "not found" {
            return Of[string](0) // recover with default
        }
    }
    return func(input string) Validation[int] {
        return either.Left[int](errs)
    }
})

decoder := recoverFromNotFound(failingDecoder)
result := decoder("input")
// Result: Success(0) - recovered from failure
```

### 4. Access to Original Input
The handler returns a `Decode[I, A]` function, giving it access to the original input:

```go
handler := ChainLeft(func(errs Errors) Decode[string, int] {
    return func(input string) Validation[int] {
        // Can access both errs and input here
        if input == "special" {
            return validation.Of(999)
        }
        return either.Left[int](errs)
    }
})
```

## Use Cases

### 1. Fallback Decoding (OrElse reads better)

```go
// Primary decoder that may fail
primaryDecoder := func(input string) Validation[int] {
    n, err := strconv.Atoi(input)
    if err != nil {
        return either.Left[int](validation.Errors{
            {Value: input, Messsage: "not a valid integer"},
        })
    }
    return validation.Of(n)
}

// Use OrElse for semantic clarity - "try primary, or else use default"
withDefault := OrElse(func(errs Errors) Decode[string, int] {
    return Of[string](0) // default to 0 if decoding fails
})

decoder := withDefault(primaryDecoder)

result1 := decoder("42")    // Success(42)
result2 := decoder("abc")   // Success(0) - fallback
```

### 2. Error Context Addition (ChainLeft reads better)

```go
decodeUserAge := func(data map[string]any) Validation[int] {
    age, ok := data["age"].(int)
    if !ok {
        return either.Left[int](validation.Errors{
            {Value: data["age"], Messsage: "invalid type"},
        })
    }
    return validation.Of(age)
}

// Use ChainLeft when emphasizing error transformation
addContext := ChainLeft(func(errs Errors) Decode[map[string]any, int] {
    return func(data map[string]any) Validation[int] {
        return either.Left[int](validation.Errors{
            {
                Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
                Messsage: "failed to decode user age",
            },
        })
    }
})

decoder := addContext(decodeUserAge)
// Errors will include both original error and context
```

### 3. Conditional Recovery Based on Input

```go
decodePort := func(input string) Validation[int] {
    port, err := strconv.Atoi(input)
    if err != nil {
        return either.Left[int](validation.Errors{
            {Value: input, Messsage: "invalid port"},
        })
    }
    return validation.Of(port)
}

// Recover with different defaults based on input
smartDefault := OrElse(func(errs Errors) Decode[string, int] {
    return func(input string) Validation[int] {
        // Check input to determine appropriate default
        if strings.Contains(input, "http") {
            return validation.Of(80)
        }
        if strings.Contains(input, "https") {
            return validation.Of(443)
        }
        return validation.Of(8080)
    }
})

decoder := smartDefault(decodePort)
result1 := decoder("http-server")   // Success(80)
result2 := decoder("https-server")  // Success(443)
result3 := decoder("other")         // Success(8080)
```

### 4. Pipeline Composition

```go
type Config struct {
    DatabaseURL string
}

decodeConfig := func(data map[string]any) Validation[Config] {
    url, ok := data["db_url"].(string)
    if !ok {
        return either.Left[Config](validation.Errors{
            {Messsage: "missing db_url"},
        })
    }
    return validation.Of(Config{DatabaseURL: url})
}

// Build a pipeline with multiple error handlers
decoder := F.Pipe2(
    decodeConfig,
    OrElse(func(errs Errors) Decode[map[string]any, Config] {
        // Try environment variable as fallback
        return func(data map[string]any) Validation[Config] {
            if url := os.Getenv("DATABASE_URL"); url != "" {
                return validation.Of(Config{DatabaseURL: url})
            }
            return either.Left[Config](errs)
        }
    }),
    OrElse(func(errs Errors) Decode[map[string]any, Config] {
        // Final fallback to default
        return Of[map[string]any](Config{
            DatabaseURL: "localhost:5432",
        })
    }),
)
```

## Comparison with validation.ChainLeft

The decode package's [`ChainLeft`](monad.go:53) wraps [`validation.ChainLeft`](../validation/monad.go:304) using the Reader transformer pattern:

| Aspect | validation.ChainLeft | decode.ChainLeft |
|--------|---------------------|------------------|
| **Input** | `Validation[A]` | `Decode[I, A]` (function) |
| **Handler** | `func(Errors) Validation[A]` | `func(Errors) Decode[I, A]` |
| **Output** | `Validation[A]` | `Decode[I, A]` (function) |
| **Context** | No input access | Access to original input `I` |
| **Use Case** | Pure validation logic | Decoding with input-dependent recovery |

The key difference is that decode's version gives handlers access to the original input through the returned `Decode[I, A]` function.

## When to Use Which Name

### Use **OrElse** when:
- Emphasizing fallback/alternative decoding logic
- Providing default values on decode failure
- The intent is "try this, or else try that"
- Code reads more naturally with "or else"

### Use **ChainLeft** when:
- Emphasizing technical error channel transformation
- Adding context or enriching error information
- The focus is on error handling mechanics
- Working with other functional programming concepts

## Verification

The test suite in [`monad_test.go`](monad_test.go:385) includes comprehensive tests proving that `OrElse` and `ChainLeft` are equivalent:

- ✅ Identical behavior for Success values
- ✅ Identical behavior for error recovery
- ✅ Identical behavior for error aggregation
- ✅ Identical behavior in pipeline composition
- ✅ Identical behavior for multiple error scenarios
- ✅ Both provide access to original input

Run the tests:
```bash
go test -v -run "TestChainLeft|TestOrElse" ./optics/codec/decode
```

## Conclusion

**`OrElse` is exactly the same as `ChainLeft`** in the decode package - they are aliases with identical implementations and behavior. Both:

1. **Delegate to validation.ChainLeft** for error handling logic
2. **Aggregate errors** when transformations fail
3. **Preserve successes** unchanged
4. **Enable recovery** from decode failures
5. **Provide access** to the original input

The choice between them is purely about **code readability and semantic intent**:
- Use **`OrElse`** when emphasizing fallback/alternative decoding
- Use **`ChainLeft`** when emphasizing error transformation

Both maintain the critical property of **error aggregation**, ensuring all validation failures are preserved and reported together.