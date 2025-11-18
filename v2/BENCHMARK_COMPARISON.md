# Benchmark Comparison: Idiomatic vs Standard Either/Result

**Date:** 2025-11-18
**System:** AMD Ryzen 7 PRO 7840U w/ Radeon 780M Graphics (16 cores)
**Go Version:** go1.23+

This document provides a detailed performance comparison between the optimized `either` package and the `idiomatic/result` package after recent optimizations to the either package.

## Executive Summary

After optimizations to the `either` package, the performance characteristics have changed significantly:

### Key Findings

1. **Constructors & Predicates**: Both packages now perform comparably (~1-2 ns/op) with **zero heap allocations**
2. **Zero-allocation insight**: The `Either` struct (24 bytes) does NOT escape to heap - Go returns it by value on the stack
3. **Core Operations**: Idiomatic package has a **consistent advantage** of 1.2x - 2.3x for most operations
4. **Complex Operations**: Idiomatic package shows **massive advantages**:
   - ChainFirst (Right): **32.4x faster** (87.6 ns → 2.7 ns, 72 B → 0 B)
   - Pipeline operations: **2-3x faster** with lower allocations
5. **All simple operations**: Both maintain **zero heap allocations** (0 B/op, 0 allocs/op)

### Winner by Category

| Category | Winner | Reason |
|----------|--------|--------|
| Constructors | **TIE** | Both ~1.3-1.8 ns/op |
| Predicates | **TIE** | Both ~1.2-1.5 ns/op |
| Simple Transformations | **Idiomatic** | 1.2-2x faster |
| Monadic Operations | **Idiomatic** | 1.2-2.3x faster |
| Complex Chains | **Idiomatic** | 32x faster, zero allocs |
| Pipelines | **Idiomatic** | 2-2.4x faster, fewer allocs |
| Extraction | **Idiomatic** | 6x faster (GetOrElse) |

## Detailed Benchmark Results

### Constructor Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| Left      | 1.76           | **1.35**          | **1.3x** ✓ | 0 B/op | 0 B/op |
| Right     | 1.38           | 1.43              | 1.0x    | 0 B/op | 0 B/op |
| Of        | 1.68           | **1.22**          | **1.4x** ✓ | 0 B/op | 0 B/op |

**Analysis:** Both packages perform extremely well with **zero heap allocations**. Idiomatic has a slight edge on Left and Of.

**Important Clarification: Neither Package Escapes to Heap**

A common misconception is that struct-based Either escapes to heap while tuples stay on stack. The benchmarks prove this is FALSE:

```go
// Either package - NO heap allocation
type Either[E, A any] struct {
    r      A           // 8 bytes
    l      E           // 8 bytes
    isLeft bool        // 1 byte + 7 padding
}                      // Total: 24 bytes

func Of[E, A any](value A) Either[E, A] {
    return Right[E](value)  // Returns 24-byte struct BY VALUE
}

// Benchmark result: 0 B/op, 0 allocs/op ✓
```

**Why Either doesn't escape:**
1. **Small struct** - At 24 bytes, it's below Go's escape threshold (~64 bytes)
2. **Return by value** - Go returns small structs on the stack
3. **Inlining** - The `//go:inline` directive eliminates function overhead
4. **No pointers** - No pointer escapes in normal usage

**Idiomatic package:**
```go
// Returns native tuple - always stack allocated
func Right[A any](a A) (A, error) {
    return a, nil  // 16 bytes total (8 + 8)
}

// Benchmark result: 0 B/op, 0 allocs/op ✓
```

**Both achieve zero allocations** - the performance difference comes from other factors like function composition overhead, not from constructor allocations.

### Predicate Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| IsLeft    | 1.45           | **1.35**          | **1.1x** ✓ | 0 B/op | 0 B/op |
| IsRight   | 1.47           | 1.51              | 1.0x    | 0 B/op | 0 B/op |

**Analysis:** Virtually identical performance. The optimizations brought them to parity.

### Fold Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| MonadFold (Right) | 2.71   | -                 | -       | 0 B/op | - |
| MonadFold (Left)  | 2.26   | -                 | -       | 0 B/op | - |
| Fold (Right)      | 4.03   | **2.75**          | **1.5x** ✓ | 0 B/op | 0 B/op |
| Fold (Left)       | 3.69   | **2.40**          | **1.5x** ✓ | 0 B/op | 0 B/op |

**Analysis:** Idiomatic package is 1.5x faster for curried Fold operations.

### Unwrap Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Note |
|-----------|----------------|-------------------|------|
| Unwrap (Right)      | 1.27  | N/A | Either-specific |
| Unwrap (Left)       | 1.24  | N/A | Either-specific |
| UnwrapError (Right) | 1.27  | N/A | Either-specific |
| UnwrapError (Left)  | 1.27  | N/A | Either-specific |
| ToError (Right)     | N/A   | 1.40 | Idiomatic-specific |
| ToError (Left)      | N/A   | 1.84 | Idiomatic-specific |

**Analysis:** Both provide fast unwrapping. Idiomatic's tuple return is naturally unwrapped.

### Map Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| MonadMap (Right)  | 2.96   | -                | -       | 0 B/op | - |
| MonadMap (Left)   | 1.99   | -                | -       | 0 B/op | - |
| Map (Right)       | 5.13   | **4.34**         | **1.2x** ✓ | 0 B/op | 0 B/op |
| Map (Left)        | 4.19   | **2.48**         | **1.7x** ✓ | 0 B/op | 0 B/op |
| MapLeft (Right)   | 3.93   | **2.22**         | **1.8x** ✓ | 0 B/op | 0 B/op |
| MapLeft (Left)    | 7.22   | **3.51**         | **2.1x** ✓ | 0 B/op | 0 B/op |

**Analysis:** Idiomatic is consistently faster across all Map variants, especially for error path (Left).

### BiMap Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| BiMap (Right)  | 16.79  | **3.82**         | **4.4x** ✓ | 0 B/op | 0 B/op |
| BiMap (Left)   | 11.47  | **3.47**         | **3.3x** ✓ | 0 B/op | 0 B/op |

**Analysis:** Idiomatic package shows significant advantage for BiMap operations (3-4x faster).

### Chain (Monadic Bind) Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| MonadChain (Right) | 2.89  | -                | -       | 0 B/op | - |
| MonadChain (Left)  | 2.03  | -                | -       | 0 B/op | - |
| Chain (Right)      | 5.44  | **2.34**         | **2.3x** ✓ | 0 B/op | 0 B/op |
| Chain (Left)       | 4.44  | **2.53**         | **1.8x** ✓ | 0 B/op | 0 B/op |
| ChainFirst (Right) | 87.62 | **2.71**         | **32.4x** ✓✓✓ | 72 B, 3 allocs | 0 B, 0 allocs |
| ChainFirst (Left)  | 3.94  | **2.48**         | **1.6x** ✓ | 0 B/op | 0 B/op |

**Analysis:**
- Idiomatic is 2x faster for standard Chain operations
- **ChainFirst shows the most dramatic difference**: 32.4x faster with zero allocations vs 72 bytes!

### Flatten Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Note |
|-----------|----------------|-------------------|------|
| Flatten (Right) | 8.73  | N/A | Either-specific nested structure |
| Flatten (Left)  | 8.86  | N/A | Either-specific nested structure |

**Analysis:** Flatten is specific to Either's nested structure handling.

### Applicative Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| MonadAp (RR)  | 3.81  | -                | -       | 0 B/op | - |
| MonadAp (RL)  | 3.07  | -                | -       | 0 B/op | - |
| MonadAp (LR)  | 3.08  | -                | -       | 0 B/op | - |
| Ap (RR)       | 6.99  | -                | -       | 0 B/op | - |

**Analysis:** MonadAp is fast in Either. Idiomatic package doesn't expose direct Ap benchmarks.

### Alternative Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| Alt (RR)       | 5.72  | **2.40**         | **2.4x** ✓ | 0 B/op | 0 B/op |
| Alt (LR)       | 4.89  | **2.39**         | **2.0x** ✓ | 0 B/op | 0 B/op |
| OrElse (Right) | 5.28  | **2.40**         | **2.2x** ✓ | 0 B/op | 0 B/op |
| OrElse (Left)  | 3.99  | **2.42**         | **1.6x** ✓ | 0 B/op | 0 B/op |

**Analysis:** Idiomatic package is consistently 2x faster for alternative operations.

### GetOrElse Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| GetOrElse (Right) | 9.01  | **1.49**         | **6.1x** ✓✓ | 0 B/op | 0 B/op |
| GetOrElse (Left)  | 6.35  | **2.08**         | **3.1x** ✓✓ | 0 B/op | 0 B/op |

**Analysis:** Idiomatic package shows dramatic advantage for value extraction (3-6x faster).

### TryCatch Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Note |
|-----------|----------------|-------------------|------|
| TryCatch (Success)       | 2.39  | N/A | Either-specific |
| TryCatch (Error)         | 3.40  | N/A | Either-specific |
| TryCatchError (Success)  | 3.32  | N/A | Either-specific |
| TryCatchError (Error)    | 6.44  | N/A | Either-specific |

**Analysis:** TryCatch/TryCatchError are Either-specific for wrapping (value, error) tuples.

### Other Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| Swap (Right)      | 2.30  | -                | -       | 0 B/op | - |
| Swap (Left)       | 3.05  | -                | -       | 0 B/op | - |
| MapTo (Right)     | -     | 1.60             | -       | -      | 0 B/op |
| MapTo (Left)      | -     | 1.73             | -       | -      | 0 B/op |
| ChainTo (Right)   | -     | 2.66             | -       | -      | 0 B/op |
| ChainTo (Left)    | -     | 2.85             | -       | -      | 0 B/op |
| Reduce (Right)    | -     | 2.34             | -       | -      | 0 B/op |
| Reduce (Left)     | -     | 1.40             | -       | -      | 0 B/op |
| Flap (Right)      | -     | 3.86             | -       | -      | 0 B/op |
| Flap (Left)       | -     | 2.58             | -       | -      | 0 B/op |

### FromPredicate Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| FromPredicate (Pass) | -     | 3.38  | -       | -      | 0 B/op |
| FromPredicate (Fail) | -     | 5.03  | -       | -      | 0 B/op |

**Analysis:** FromPredicate in idiomatic shows good performance for validation patterns.

### Option Conversion

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| ToOption (Right)   | -     | 1.17  | -       | -      | 0 B/op |
| ToOption (Left)    | -     | 1.21  | -       | -      | 0 B/op |
| FromOption (Some)  | -     | 2.68  | -       | -      | 0 B/op |
| FromOption (None)  | -     | 3.72  | -       | -      | 0 B/op |

**Analysis:** Very fast conversion between Result and Option in idiomatic package.

## Pipeline Benchmarks

These benchmarks measure realistic composition scenarios using F.Pipe.

### Simple Map Pipeline

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| Pipeline Map (Right) | 112.7  | **46.5**  | **2.4x** ✓ | 72 B, 3 allocs | 48 B, 2 allocs |
| Pipeline Map (Left)  | 116.8  | **47.2**  | **2.5x** ✓ | 72 B, 3 allocs | 48 B, 2 allocs |

### Chain Pipeline

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| Pipeline Chain (Right) | 74.4  | **26.1**  | **2.9x** ✓ | 48 B, 2 allocs | 24 B, 1 allocs |
| Pipeline Chain (Left)  | 86.4  | **25.7**  | **3.4x** ✓ | 48 B, 2 allocs | 24 B, 1 allocs |

### Complex Pipeline (Map → Chain → Map)

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| Complex (Right) | 279.8  | **116.3**  | **2.4x** ✓ | 192 B, 8 allocs | 120 B, 5 allocs |
| Complex (Left)  | 288.1  | **115.8**  | **2.5x** ✓ | 192 B, 8 allocs | 120 B, 5 allocs |

**Analysis:**
- Idiomatic package shows **2-3.4x speedup** for realistic pipelines
- Significantly fewer allocations in all pipeline scenarios
- The gap widens as pipelines become more complex

## Array/Collection Operations

### TraverseArray

| Operation | Either (ns/op) | Idiomatic (ns/op) | Note |
|-----------|----------------|-------------------|------|
| TraverseArray (Success) | -     | 32.3  | 48 B, 1 alloc |
| TraverseArray (Error)   | -     | 28.3  | 48 B, 1 alloc |

**Analysis:** Idiomatic package provides efficient array traversal with minimal allocations.

## Validation (ApV)

### ApV Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| ApV (BothRight) | -     | 1.17   | -       | -      | 0 B/op |
| ApV (BothLeft)  | -     | 141.5  | -       | -      | 48 B, 2 allocs |

**Analysis:** Idiomatic's validation applicative shows fast success path, with allocations only when accumulating errors.

## String Formatting

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| String/ToString (Right) | 139.9  | **81.8**  | **1.7x** ✓ | 16 B, 1 alloc | 16 B, 1 alloc |
| String/ToString (Left)  | 161.6  | **72.7**  | **2.2x** ✓ | 48 B, 1 alloc | 24 B, 1 alloc |

**Analysis:** Idiomatic package formats strings faster with fewer allocations for Left values.

## Do-Notation

| Operation | Either (ns/op) | Idiomatic (ns/op) | Note |
|-----------|----------------|-------------------|------|
| Do        | 2.03  | -                | Either-specific |
| Bind      | 153.4 | -                | 96 B, 4 allocs |
| Let       | 33.5  | -                | 16 B, 1 alloc |

**Analysis:** Do-notation is specific to Either package for monadic composition patterns.

## Summary Statistics

### Simple Operations (< 10 ns/op)

**Either Package:**
- Count: 24 operations
- Average: 3.2 ns/op
- Range: 1.24 - 9.01 ns/op

**Idiomatic Package:**
- Count: 36 operations
- Average: 2.1 ns/op
- Range: 1.17 - 5.03 ns/op

**Winner:** Idiomatic (1.5x faster average)

### Complex Operations (Pipelines, allocations)

**Either Package:**
- Pipeline Map: 112.7 ns/op (72 B, 3 allocs)
- Pipeline Chain: 74.4 ns/op (48 B, 2 allocs)
- Complex: 279.8 ns/op (192 B, 8 allocs)
- ChainFirst: 87.6 ns/op (72 B, 3 allocs)

**Idiomatic Package:**
- Pipeline Map: 46.5 ns/op (48 B, 2 allocs)
- Pipeline Chain: 26.1 ns/op (24 B, 1 allocs)
- Complex: 116.3 ns/op (120 B, 5 allocs)
- ChainFirst: 2.71 ns/op (0 B, 0 allocs)

**Winner:** Idiomatic (2-32x faster, significantly fewer allocations)

### Allocation Analysis

**Either Package:**
- Zero-allocation operations: Most simple operations
- Operations with allocations: Pipelines, Bind, Do-notation, ChainFirst

**Idiomatic Package:**
- Zero-allocation operations: Almost all operations except pipelines and validation
- Significantly fewer allocations in pipeline scenarios
- ChainFirst: **Zero allocations** (vs 72 B in Either)

## Performance Characteristics

### Where Either Package Excels

1. **Comparable to Idiomatic**: After optimizations, Either matches Idiomatic for constructors and predicates
2. **Feature Richness**: More operations (Do-notation, Bind, Let, Flatten, Swap)
3. **Type Flexibility**: Full Either[E, A] with custom error types

### Where Idiomatic Package Excels

1. **Core Operations**: 1.2-2.3x faster for Map, Chain, Fold
2. **Complex Operations**: 32x faster for ChainFirst
3. **Pipelines**: 2-3.4x faster with fewer allocations
4. **Extraction**: 3-6x faster for GetOrElse
5. **Alternative**: 2x faster for Alt/OrElse
6. **BiMap**: 3-4x faster
7. **Consistency**: More predictable performance profile

## Real-World Performance Impact

### Hot Path Example (1 million operations)

```go
// Map operation (very common)
// Either:    5.13 ns/op × 1M = 5.13 ms
// Idiomatic: 4.34 ns/op × 1M = 4.34 ms
// Savings: 0.79 ms per million operations

// Chain operation (common in pipelines)
// Either:    5.44 ns/op × 1M = 5.44 ms
// Idiomatic: 2.34 ns/op × 1M = 2.34 ms
// Savings: 3.10 ms per million operations

// Pipeline Complex (realistic composition)
// Either:    279.8 ns/op × 1M = 279.8 ms
// Idiomatic: 116.3 ns/op × 1M = 116.3 ms
// Savings: 163.5 ms per million operations
```

### Memory Impact

For 1 million ChainFirst operations:
- Either: 72 MB allocated
- Idiomatic: 0 MB allocated
- **Savings: 72 MB + reduced GC pressure**

## Recommendations

### Use Idiomatic Package When:

1. **Performance is Critical**
   - Hot paths in your application
   - High-throughput services (>10k req/s)
   - Complex operation chains
   - Memory-constrained environments

2. **Natural Go Integration**
   - Working with stdlib (value, error) patterns
   - Team familiar with Go idioms
   - Simple migration from existing code
   - Want zero-cost abstractions

3. **Pipeline-Heavy Code**
   - 2-3.4x faster pipelines
   - Significantly fewer allocations
   - Better CPU cache utilization

### Use Either Package When:

1. **Feature Requirements**
   - Need custom error types (Either[E, A])
   - Using Do-notation for complex compositions
   - Need Flatten, Swap, or other Either-specific operations
   - Porting from FP languages (Scala, Haskell)

2. **Type Safety Over Performance**
   - Explicit Either semantics
   - Algebraic data type guarantees
   - Teaching/learning FP concepts

3. **Moderate Performance Needs**
   - After optimizations, Either is quite fast
   - Difference matters only at high scale
   - Code clarity > micro-optimizations

### Hybrid Approach

```go
// Use Either for complex type safety
import "github.com/IBM/fp-go/v2/either"
type ValidationError struct { Field, Message string }
validated := either.Either[ValidationError, Input]{...}

// Convert to Idiomatic for hot path
import "github.com/IBM/fp-go/v2/idiomatic/result"
value, err := either.UnwrapError(either.MapLeft(toError)(validated))
processed, err := result.Chain(hotPathProcessing)(value, err)
```

## Conclusion

After optimizations to the Either package:

1. **Both packages achieve zero heap allocations for constructors** - The Either struct (24 bytes) does NOT escape to heap
2. **Simple operations** are now **comparable** between both packages (~1-2 ns/op, 0 B/op)
3. **Core transformations** favor Idiomatic by **1.2-2.3x**
4. **Complex operations** heavily favor Idiomatic by **2-32x**
5. **Memory efficiency** strongly favors Idiomatic (especially ChainFirst: 72 B → 0 B)
6. **Real-world pipelines** show **2-3.4x speedup** with Idiomatic

### Key Insight: No Heap Escape Myth

A critical finding: **Both packages avoid heap allocations for simple operations.** The Either struct is small enough (24 bytes) that Go returns it by value on the stack, not the heap. The `0 B/op, 0 allocs/op` benchmarks confirm this.

The performance differences come from:
- **Function composition overhead** in complex operations
- **Currying and closure creation** in pipelines
- **Tuple simplicity** vs struct field access

Not from constructor allocations—both are equally efficient there.

### Final Verdict

The idiomatic package provides a compelling performance advantage for production workloads while maintaining zero-cost functional programming abstractions. The Either package remains excellent for type safety, feature richness, and scenarios where explicit Either[E, A] semantics are valuable.

**Bottom Line:**
- For **high-performance Go services**: idiomatic package is the clear winner (1.2-32x faster)
- For **type-safe, feature-rich FP**: Either package is excellent (comparable simple ops, more features)
- **Both avoid heap allocations** for constructors—choose based on your performance vs features trade-off
