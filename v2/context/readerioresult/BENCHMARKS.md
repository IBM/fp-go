# ReaderIOResult Benchmarks

This document describes the benchmark suite for the `context/readerioeither` package and how to interpret the results to identify performance bottlenecks.

## Running Benchmarks

To run all benchmarks:
```bash
cd context/readerioeither
go test -bench=. -benchmem
```

To run specific benchmarks:
```bash
go test -bench=BenchmarkMap -benchmem
go test -bench=BenchmarkChain -benchmem
go test -bench=BenchmarkApPar -benchmem
```

To run with more iterations for stable results:
```bash
go test -bench=. -benchmem -benchtime=100000x
```

## Benchmark Categories

### 1. Core Constructors
- `BenchmarkLeft` - Creating Left (error) values (~64ns, 2 allocs)
- `BenchmarkRight` - Creating Right (success) values (~64ns, 2 allocs)
- `BenchmarkOf` - Creating Right values via Of (~47ns, 2 allocs)

**Key Insights:**
- All constructors allocate 2 times (64B total)
- `Of` is slightly faster than `Right` due to inlining
- Construction is very fast, suitable for hot paths

### 2. Conversion Operations
- `BenchmarkFromEither_Right/Left` - Converting Either to ReaderIOResult (~70ns, 2 allocs)
- `BenchmarkFromIO` - Converting IO to ReaderIOResult (~78ns, 3 allocs)
- `BenchmarkFromIOEither_Right/Left` - Converting IOEither (~23ns, 1 alloc)

**Key Insights:**
- FromIOEither is the fastest conversion (~23ns)
- FromIO has an extra allocation due to wrapping
- All conversions are lightweight

### 3. Execution Operations
- `BenchmarkExecute_Right` - Executing Right computation (~37ns, 1 alloc)
- `BenchmarkExecute_Left` - Executing Left computation (~48ns, 1 alloc)
- `BenchmarkExecute_WithContext` - Executing with context (~42ns, 1 alloc)

**Key Insights:**
- Execution is very fast with minimal allocations
- Left path is slightly slower due to error handling
- Context overhead is minimal (~5ns)

### 4. Functor Operations (Map)
- `BenchmarkMonadMap_Right/Left` - Direct map (~135ns, 5 allocs)
- `BenchmarkMap_Right/Left` - Curried map (~24ns, 1 alloc)
- `BenchmarkMapTo_Right` - Replacing with constant (~69ns, 1 alloc)

**Key Insights:**
- **Bottleneck:** MonadMap has 5 allocations (128B)
- Curried Map is ~5x faster with fewer allocations
- **Recommendation:** Use curried `Map` instead of `MonadMap`

### 5. Monad Operations (Chain)
- `BenchmarkMonadChain_Right/Left` - Direct chain (~190ns, 6 allocs)
- `BenchmarkChain_Right/Left` - Curried chain (~28ns, 1 alloc)
- `BenchmarkChainFirst_Right/Left` - Chain preserving original (~27ns, 1 alloc)
- `BenchmarkFlatten_Right/Left` - Removing nesting (~147ns, 7 allocs)

**Key Insights:**
- **Bottleneck:** MonadChain has 6 allocations (160B)
- Curried Chain is ~7x faster
- ChainFirst is as fast as Chain
- **Bottleneck:** Flatten has 7 allocations
- **Recommendation:** Use curried `Chain` instead of `MonadChain`

### 6. Applicative Operations (Ap)
- `BenchmarkMonadApSeq_RightRight` - Sequential apply (~281ns, 8 allocs)
- `BenchmarkMonadApPar_RightRight` - Parallel apply (~49ns, 3 allocs)
- `BenchmarkExecuteApSeq_RightRight` - Executing sequential (~1403ns, 8 allocs)
- `BenchmarkExecuteApPar_RightRight` - Executing parallel (~5606ns, 61 allocs)

**Key Insights:**
- **Major Bottleneck:** Parallel execution has 61 allocations (1896B)
- Construction of ApPar is fast (~49ns), but execution is expensive
- Sequential execution is faster for simple operations (~1.4μs vs ~5.6μs)
- **Recommendation:** Use ApSeq for simple operations, ApPar only for truly independent, expensive computations
- Parallel overhead includes context management and goroutine coordination

### 7. Alternative Operations
- `BenchmarkAlt_RightRight/LeftRight` - Providing alternatives (~210-344ns, 6 allocs)
- `BenchmarkOrElse_Right/Left` - Recovery from Left (~40-52ns, 2 allocs)

**Key Insights:**
- Alt has significant overhead (6 allocations)
- OrElse is much more efficient for error recovery
- **Recommendation:** Prefer OrElse over Alt when possible

### 8. Chain Operations with Different Types
- `BenchmarkChainEitherK_Right/Left` - Chaining Either (~25ns, 1 alloc)
- `BenchmarkChainIOK_Right/Left` - Chaining IO (~55ns, 1 alloc)
- `BenchmarkChainIOEitherK_Right/Left` - Chaining IOEither (~53ns, 1 alloc)

**Key Insights:**
- All chain-K operations are efficient
- ChainEitherK is fastest (pure transformation)
- ChainIOK and ChainIOEitherK have similar performance

### 9. Context Operations
- `BenchmarkAsk` - Accessing context (~52ns, 3 allocs)
- `BenchmarkDefer` - Lazy generation (~34ns, 1 alloc)
- `BenchmarkMemoize` - Caching results (~82ns, 4 allocs)

**Key Insights:**
- Ask has 3 allocations for context wrapping
- Defer is lightweight
- Memoize has overhead but pays off for expensive computations

### 10. Delay Operations
- `BenchmarkDelay_Construction` - Creating delayed computation (~19ns, 1 alloc)
- `BenchmarkTimer_Construction` - Creating timer (~92ns, 3 allocs)

**Key Insights:**
- Delay construction is very cheap
- Timer has additional overhead for time operations

### 11. TryCatch Operations
- `BenchmarkTryCatch_Success/Error` - Creating TryCatch (~33ns, 1 alloc)
- `BenchmarkExecuteTryCatch_Success/Error` - Executing TryCatch (~3ns, 0 allocs)

**Key Insights:**
- TryCatch construction is cheap
- Execution is extremely fast with zero allocations
- Excellent for wrapping Go error-returning functions

### 12. Pipeline Operations
- `BenchmarkPipeline_Map_Right/Left` - Single Map in pipeline (~200-306ns, 9 allocs)
- `BenchmarkPipeline_Chain_Right/Left` - Single Chain in pipeline (~155-217ns, 7 allocs)
- `BenchmarkPipeline_Complex_Right/Left` - Multiple operations (~777-1039ns, 25 allocs)
- `BenchmarkExecutePipeline_Complex_Right` - Executing complex pipeline (~533ns, 10 allocs)

**Key Insights:**
- **Major Bottleneck:** Pipeline operations allocate heavily
- Single Map: ~200ns with 9 allocations (224B)
- Complex pipeline: ~900ns with 25 allocations (640B)
- **Recommendation:** Avoid F.Pipe in hot paths, use direct function calls

### 13. Do-Notation Operations
- `BenchmarkDo` - Creating empty context (~45ns, 2 allocs)
- `BenchmarkBind_Right` - Binding values (~25ns, 1 alloc)
- `BenchmarkLet_Right` - Pure computations (~23ns, 1 alloc)
- `BenchmarkApS_Right` - Applicative binding (~99ns, 4 allocs)

**Key Insights:**
- Do-notation operations are efficient
- Bind and Let are very fast
- ApS has more overhead (4 allocations)
- Much better than either package's Bind (~130ns vs ~25ns)

### 14. Traverse Operations
- `BenchmarkTraverseArray_Empty` - Empty array (~689ns, 13 allocs)
- `BenchmarkTraverseArray_Small` - 3 elements (~1971ns, 37 allocs)
- `BenchmarkTraverseArray_Medium` - 10 elements (~4386ns, 93 allocs)
- `BenchmarkTraverseArraySeq_Small` - Sequential (~1885ns, 52 allocs)
- `BenchmarkTraverseArrayPar_Small` - Parallel (~1362ns, 37 allocs)
- `BenchmarkExecuteTraverseArraySeq_Small` - Executing sequential (~1080ns, 34 allocs)
- `BenchmarkExecuteTraverseArrayPar_Small` - Executing parallel (~18560ns, 202 allocs)

**Key Insights:**
- **Bottleneck:** Traverse operations allocate per element
- Empty array still has 13 allocations (overhead)
- Parallel construction is faster but execution is much slower
- **Major Bottleneck:** Parallel execution: ~18.5μs with 202 allocations
- Sequential execution is ~17x faster for small arrays
- **Recommendation:** Use sequential traverse for small collections, parallel only for large, expensive operations

### 15. Record Operations
- `BenchmarkTraverseRecord_Small` - 3 entries (~1444ns, 55 allocs)
- `BenchmarkSequenceRecord_Small` - 3 entries (~1073ns, 54 allocs)

**Key Insights:**
- Record operations have high allocation overhead
- Similar performance to array traversal
- Allocations scale with map size

### 16. Resource Management
- `BenchmarkWithResource_Success` - Creating resource wrapper (~193ns, 8 allocs)
- `BenchmarkExecuteWithResource_Success` - Executing with resource (varies)
- `BenchmarkExecuteWithResource_ErrorInBody` - Error handling (varies)

**Key Insights:**
- Resource management has 8 allocations for safety
- Ensures proper cleanup even on errors
- Overhead is acceptable for resource safety guarantees

### 17. Context Cancellation
- `BenchmarkExecute_CanceledContext` - Executing with canceled context
- `BenchmarkExecuteApPar_CanceledContext` - Parallel with canceled context

**Key Insights:**
- Cancellation is handled efficiently
- Minimal overhead for checking cancellation
- ApPar respects cancellation properly

## Performance Bottlenecks Summary

### Critical Bottlenecks (>100ns or >5 allocations)

1. **Pipeline operations with F.Pipe** (~200-1000ns, 9-25 allocations)
   - **Impact:** High - commonly used pattern
   - **Mitigation:** Use direct function calls in hot paths
   - **Example:**
     ```go
     // Slow (200ns, 9 allocs)
     result := F.Pipe1(rioe, Map[int](transform))
     
     // Fast (24ns, 1 alloc)
     result := Map[int](transform)(rioe)
     ```

2. **MonadMap and MonadChain** (~135-207ns, 5-6 allocations)
   - **Impact:** High - fundamental operations
   - **Mitigation:** Use curried versions (Map, Chain)
   - **Speedup:** 5-7x faster

3. **Parallel applicative execution** (~5.6μs, 61 allocations)
   - **Impact:** High when used
   - **Mitigation:** Use ApSeq for simple operations
   - **Note:** Only use ApPar for truly independent, expensive computations

4. **Parallel traverse execution** (~18.5μs, 202 allocations)
   - **Impact:** High for collections
   - **Mitigation:** Use sequential traverse for small collections
   - **Threshold:** Consider parallel only for >100 elements with expensive operations

5. **Alt operations** (~210-344ns, 6 allocations)
   - **Impact:** Medium
   - **Mitigation:** Use OrElse for error recovery (40-52ns, 2 allocs)

### Minor Bottlenecks (50-100ns or 3-4 allocations)

6. **Flatten operations** (~147ns, 7 allocations)
   - **Impact:** Low - less commonly used
   - **Mitigation:** Avoid unnecessary nesting

7. **Memoize** (~82ns, 4 allocations)
   - **Impact:** Low - overhead pays off for expensive computations
   - **Mitigation:** Only use for computations >1μs

8. **ApS in do-notation** (~99ns, 4 allocations)
   - **Impact:** Low
   - **Mitigation:** Use Let or Bind when possible

## Optimization Recommendations

### For Hot Paths

1. **Use curried functions over Monad* versions:**
   ```go
   // Instead of:
   result := MonadMap(rioe, transform)  // 135ns, 5 allocs
   
   // Use:
   result := Map[int](transform)(rioe)  // 24ns, 1 alloc
   ```

2. **Avoid F.Pipe in performance-critical code:**
   ```go
   // Instead of:
   result := F.Pipe3(rioe, Map(f1), Chain(f2), Map(f3))  // 1000ns, 25 allocs
   
   // Use:
   result := Map(f3)(Chain(f2)(Map(f1)(rioe)))  // Much faster
   ```

3. **Use sequential operations for small collections:**
   ```go
   // For arrays with <10 elements:
   result := TraverseArraySeq(f)(arr)  // 1.9μs, 52 allocs
   
   // Instead of:
   result := TraverseArrayPar(f)(arr)  // 18.5μs, 202 allocs (when executed)
   ```

4. **Prefer OrElse over Alt for error recovery:**
   ```go
   // Instead of:
   result := Alt(alternative)(rioe)  // 210-344ns, 6 allocs
   
   // Use:
   result := OrElse(recover)(rioe)  // 40-52ns, 2 allocs
   ```

### For Context Operations

- Context operations are generally efficient
- Ask has 3 allocations but is necessary for context access
- Cancellation checking is fast and should be used liberally

### For Resource Management

- WithResource overhead (8 allocations) is acceptable for safety
- Always use for resources that need cleanup
- The RAII pattern prevents resource leaks

### Memory Considerations

- Most operations have 1-2 allocations
- Monad* versions have 5-8 allocations
- Pipeline operations allocate heavily
- Parallel operations have significant allocation overhead
- Traverse operations allocate per element

## Comparative Analysis

### Fastest Operations (<50ns, <2 allocations)
- Constructors (Left, Right, Of)
- Execution (Execute_Right/Left)
- Curried operations (Map, Chain, ChainFirst)
- TryCatch execution
- Chain-K operations
- Let and Bind in do-notation

### Medium Speed (50-200ns, 2-8 allocations)
- Conversion operations
- MonadMap, MonadChain
- Flatten
- Memoize
- Resource management construction
- Traverse construction

### Slower Operations (>200ns or >8 allocations)
- Pipeline operations
- Alt operations
- Traverse operations (especially parallel)
- Applicative operations (especially parallel execution)

## Parallel vs Sequential Trade-offs

### When to Use Sequential (ApSeq, TraverseArraySeq)
- Small collections (<10 elements)
- Fast operations (<1μs per element)
- When allocations matter
- Default choice for most use cases

### When to Use Parallel (ApPar, TraverseArrayPar)
- Large collections (>100 elements)
- Expensive operations (>10μs per element)
- Independent computations
- When latency matters more than throughput
- I/O-bound operations

### Parallel Overhead
- Construction: ~50ns, 3 allocs
- Execution: +4-17μs, +40-160 allocs
- Context management and goroutine coordination
- Only worthwhile for truly expensive operations

## Conclusion

The `context/readerioeither` package is well-optimized with most operations completing in nanoseconds with minimal allocations. Key recommendations:

1. Use curried functions (Map, Chain) instead of Monad* versions (5-7x faster)
2. Avoid F.Pipe in hot paths (5-10x slower)
3. Use sequential operations by default; parallel only for expensive, independent computations
4. Prefer OrElse over Alt for error recovery (5x faster)
5. TryCatch is excellent for wrapping Go functions (near-zero execution cost)
6. Context operations are efficient; use liberally
7. Resource management overhead is acceptable for safety guarantees

For typical use cases, the performance is excellent. Only in extremely hot paths (millions of operations per second) should you consider the micro-optimizations suggested above.