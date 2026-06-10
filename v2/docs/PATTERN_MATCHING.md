# Pattern Matching with FindFirstMap and AltAllArray

This document explains how to use `array.FindFirstMap` combined with `option.AltAllArray` to implement multi-branch pattern matching in Go, providing a functional alternative to switch/case statements.

## Table of Contents

- [Overview](#overview)
- [Basic Concept](#basic-concept)
- [Why Use This Approach?](#why-use-this-approach)
- [Pattern Matching Patterns](#pattern-matching-patterns)
- [Best Practices](#best-practices)
- [Comparison with Traditional Approaches](#comparison-with-traditional-approaches)

## Overview

Pattern matching in fp-go combines two key functions:

1. **`option.AltAllArray`**: Tries multiple Option values in sequence, returning the first Some
2. **`array.FindFirstMap`**: Finds the first array element for which a selector returns Some

Together, they enable functional pattern matching similar to switch/case statements.

```go
// Signature of key functions
func AltAllArray[A any](startWith Option[A]) func([]Option[A]) Option[A]
func FindFirstMap[A, B any](sel Kleisli[A, B]) Kleisli[[]A, B]
```

## Basic Concept

Pattern matching works by:

1. **Defining matchers**: Functions that return `Some[B]` on match, `None[B]()` otherwise
2. **Combining matchers**: Apply all matchers to input and collect results in an array
3. **Finding first match**: Use `AltAllArray` to return the first Some result

```go
// Define matchers
matchCase1 := func(x T) Option[R] {
    if condition1(x) { return Some(transform1(x)) }
    return None[R]()
}

matchCase2 := func(x T) Option[R] {
    if condition2(x) { return Some(transform2(x)) }
    return None[R]()
}

defaultCase := func(x T) Option[R] {
    return Some(defaultTransform(x)) // Always matches
}

// Combine into pattern matcher
matcher := func(x T) Option[R] {
    matchers := []Option[R]{
        matchCase1(x),
        matchCase2(x),
        defaultCase(x),
    }
    return AltAllArray(None[R]())(matchers)
}

// Use
result := matcher(input) // Returns first matching branch
```

## Why Use This Approach?

### 1. Type Safety

Each branch must return the same type, enforced at compile time:

```go
// Compile error: inconsistent types
matcher := func(x int) Option[string] {
    matchers := []Option[string]{
        Some("text"),
        Some(42), // Error: cannot use int as string
    }
    return AltAllArray(None[string]())(matchers)
}
```

### 2. Composability

Matchers are first-class functions that can be composed, reused, and tested independently:

```go
// Define reusable matchers
var (
    isPositive = func(n int) Option[string] {
        if n > 0 { return Some("positive") }
        return None[string]()
    }
    
    isNegative = func(n int) Option[string] {
        if n < 0 { return Some("negative") }
        return None[string]()
    }
    
    isZero = func(n int) Option[string] {
        if n == 0 { return Some("zero") }
        return None[string]()
    }
)

// Compose into different matchers
signMatcher := func(n int) Option[string] {
    return AltAllArray(None[string]())([]Option[string]{
        isPositive(n), isNegative(n), isZero(n),
    })
}

nonZeroMatcher := func(n int) Option[string] {
    return AltAllArray(None[string]())([]Option[string]{
        isPositive(n), isNegative(n),
    })
}
```

### 3. Explicit Priority

The order of matchers in the array defines priority explicitly:

```go
// More specific matchers first
matcher := func(e Event) Option[string] {
    return AltAllArray(None[string]())([]Option[string]{
        matchCriticalError(e),  // Priority 1: Most specific
        matchError(e),          // Priority 2: Less specific
        matchWarning(e),        // Priority 3: Even less specific
        matchDefault(e),        // Priority 4: Catch-all
    })
}
```

### 4. Exhaustiveness

Add a default matcher that always returns `Some` to ensure exhaustive matching:

```go
defaultCase := func(x T) Option[R] {
    return Some(defaultValue) // Always matches
}

matcher := func(x T) Option[R] {
    return AltAllArray(None[R]())([]Option[R]{
        case1(x), case2(x), case3(x), defaultCase(x),
    })
}
// Result is always Some[R]
```

## Pattern Matching Patterns

### 1. Simple Value Matching

Match on specific values or ranges:

```go
matchZero := func(n int) Option[string] {
    if n == 0 { return Some("zero") }
    return None[string]()
}

matchPositive := func(n int) Option[string] {
    if n > 0 { return Some("positive") }
    return None[string]()
}

matchNegative := func(n int) Option[string] {
    if n < 0 { return Some("negative") }
    return None[string]()
}

classify := func(n int) Option[string] {
    return AltAllArray(None[string]())([]Option[string]{
        matchZero(n), matchPositive(n), matchNegative(n),
    })
}
```

### 2. Type-Based Matching

Match on type properties or fields:

```go
type Request struct {
    Method string
    Path   string
}

matchGET := func(r Request) Option[Handler] {
    if r.Method == "GET" { return Some(getHandler) }
    return None[Handler]()
}

matchPOST := func(r Request) Option[Handler] {
    if r.Method == "POST" { return Some(postHandler) }
    return None[Handler]()
}

routeRequest := func(r Request) Option[Handler] {
    return AltAllArray(None[Handler]())([]Option[Handler]{
        matchGET(r), matchPOST(r), matchDefault(r),
    })
}
```

### 3. Guard Patterns

Use additional conditions (guards) for precise matching:

```go
matchAdminUser := func(u User) Option[string] {
    if u.Role == "admin" && u.Active {
        return Some("full access")
    }
    return None[string]()
}

matchActiveUser := func(u User) Option[string] {
    if u.Active {
        return Some("standard access")
    }
    return None[string]()
}

matchInactiveUser := func(u User) Option[string] {
    return Some("no access") // Default
}

checkAccess := func(u User) Option[string] {
    return AltAllArray(None[string]())([]Option[string]{
        matchAdminUser(u),    // Most specific
        matchActiveUser(u),   // Less specific
        matchInactiveUser(u), // Catch-all
    })
}
```

### 4. Transformation Patterns

Each branch can transform the input differently:

```go
parseDecimal := func(s string) Option[int] {
    n, err := strconv.Atoi(s)
    if err == nil { return Some(n) }
    return None[int]()
}

parseHex := func(s string) Option[int] {
    if strings.HasPrefix(s, "0x") {
        n, err := strconv.ParseInt(s[2:], 16, 64)
        if err == nil { return Some(int(n)) }
    }
    return None[int]()
}

parseNumber := func(s string) Option[int] {
    return AltAllArray(None[int]())([]Option[int]{
        parseDecimal(s), parseHex(s),
    })
}
```

### 5. Array Pattern Matching

Use `FindFirstMap` to find the first matching element in an array:

```go
// Define a combined parser
parseNumber := func(s string) Option[int] {
    return AltAllArray(None[int]())([]Option[int]{
        parseDecimal(s), parseHex(s), parseOctal(s),
    })
}

// Find first parseable string in array
inputs := []string{"invalid", "also bad", "42", "0x2A"}
result := FindFirstMap(parseNumber)(inputs) // Some(42)
```

## Best Practices

### 1. Order Matchers by Specificity

Place more specific matchers before less specific ones:

```go
// Good: Specific to general
matcher := func(x T) Option[R] {
    return AltAllArray(None[R]())([]Option[R]{
        matchSpecificCase(x),
        matchGeneralCase(x),
        matchDefault(x),
    })
}

// Bad: General matcher shadows specific ones
matcher := func(x T) Option[R] {
    return AltAllArray(None[R]())([]Option[R]{
        matchDefault(x),      // Matches everything!
        matchSpecificCase(x), // Never reached
    })
}
```

### 2. Use Descriptive Names

Name matchers to clearly indicate what they match:

```go
// Good: Clear intent
matchAdminUser := func(u User) Option[string] { ... }
matchGuestUser := func(u User) Option[string] { ... }

// Bad: Unclear
match1 := func(u User) Option[string] { ... }
match2 := func(u User) Option[string] { ... }
```

### 3. Keep Matchers Pure

Matchers should be pure functions without side effects:

```go
// Good: Pure function
matchValid := func(x int) Option[int] {
    if x > 0 { return Some(x * 2) }
    return None[int]()
}

// Bad: Side effects
matchValid := func(x int) Option[int] {
    log.Println("Checking", x) // Side effect!
    if x > 0 { return Some(x * 2) }
    return None[int]()
}
```

### 4. Provide Default Cases When Appropriate

For exhaustive matching, add a default case:

```go
// Exhaustive: Always returns Some
matcher := func(x T) Option[R] {
    return AltAllArray(None[R]())([]Option[R]{
        matchCase1(x), matchCase2(x),
        func(T) Option[R] { return Some(defaultValue) }(x),
    })
}

// Non-exhaustive: May return None
matcher := func(x T) Option[R] {
    return AltAllArray(None[R]())([]Option[R]{
        matchCase1(x), matchCase2(x),
    })
}
```

### 5. Test Matchers Independently

Each matcher is a function that can be tested in isolation:

```go
func TestMatchPositive(t *testing.T) {
    result := matchPositive(5)
    assert.Equal(t, Some("positive"), result)
    
    result = matchPositive(-5)
    assert.Equal(t, None[string](), result)
}
```

## Comparison with Traditional Approaches

### Switch Statement

```go
// Traditional switch
func classify(n int) string {
    switch {
    case n == 0:
        return "zero"
    case n > 0:
        return "positive"
    case n < 0:
        return "negative"
    default:
        return "unknown"
    }
}
```

### Pattern Matching with AltAllArray

```go
// Functional pattern matching
var classify = func(n int) Option[string] {
    return AltAllArray(None[string]())([]Option[string]{
        func() Option[string] {
            if n == 0 { return Some("zero") }
            return None[string]()
        }(),
        func() Option[string] {
            if n > 0 { return Some("positive") }
            return None[string]()
        }(),
        func() Option[string] {
            if n < 0 { return Some("negative") }
            return None[string]()
        }(),
        Some("unknown"),
    })
}
```

### Advantages of Pattern Matching

1. **Composability**: Matchers can be reused and combined
2. **Testability**: Each matcher is independently testable
3. **Type Safety**: Return type consistency enforced at compile time
4. **Explicit Priority**: Order of matchers is clear and intentional
5. **Functional Style**: Fits naturally into functional pipelines

### When to Use Each

**Use switch/case when:**
- Simple, straightforward branching
- Performance is critical (switch may be slightly faster)
- Team prefers imperative style

**Use pattern matching when:**
- Matchers need to be reused or composed
- Building complex pattern matching logic
- Working in a functional codebase
- Need independent testing of branches
- Want explicit type safety for all branches

## See Also

- `array.FindFirstMap` - Find first matching element in array
- `option.AltAllArray` - Combine multiple Options, returning first Some
- `option.Alt` - Alternative combinator for Option values
- `array/example_pattern_matching_test.go` - Comprehensive examples