# Either Benchmarks

This document describes the benchmark suite for the `either` package and how to interpret the results to identify performance bottlenecks.

## Running Benchmarks

To run all benchmarks:
```bash
cd either
go test -bench=. -benchmem
```

To run specific benchmarks:
```bash
go test -bench=BenchmarkMap -benchmem
go test -bench=BenchmarkChain -benchmem
```

To run with more iterations for stable results:
```bash
go test -bench=. -benchmem -benchtime=1000000x
```

## Benchmark Categories

### 1. Core Constructors
- `BenchmarkLeft` - Creating Left values
- `BenchmarkRight` - Creating Right values  
- `BenchmarkOf` - Creating Right values via Of

**Key Insights:**
- Right construction is ~2x faster than Left (~0.45ns vs ~0.87ns)
- All constructors have zero allocations
- These are the fastest operations in the package

### 2. Predicates
- `BenchmarkIsLeft` - Testing if Either is Left
- `BenchmarkIsRight` - Testing if Either is Right

**Key Insights:**
- Both predicates are extremely fast (~0.3ns)
- Zero allocations
- No performance difference between Left and Right checks

### 3. Fold Operations
- `BenchmarkMonadFold_Right/Left` - Direct fold with handlers
- `BenchmarkFold_Right/Left` - Curried fold

**Key Insights:**
- Right path is ~10x faster than Left path (0.5ns vs 5-7ns)
- Curried version adds ~2ns overhead on Left path
- Zero allocations for both paths
- **Bottleneck:** Left path has higher overhead due to type assertions

### 4. Unwrap Operations
- `BenchmarkUnwrap_Right/Left` - Converting to (value, error) tuple
- `BenchmarkUnwrapError_Right/Left` - Specialized for error types

**Key Insights:**
- Right path is ~10x faster than Left path
- Zero allocations
- Similar performance characteristics to Fold

### 5. Functor Operations (Map)
- `BenchmarkMonadMap_Right/Left` - Direct map
- `BenchmarkMap_Right/Left` - Curried map
- `BenchmarkMapLeft_Right/Left` - Mapping over Left channel
- `BenchmarkBiMap_Right/Left` - Mapping both channels

**Key Insights:**
- Map operations: ~11-14ns, zero allocations
- MapLeft on Left: ~34ns with 1 allocation (16B)
- BiMap: ~29-39ns with 1 allocation (16B)
- **Bottleneck:** BiMap and MapLeft allocate closures

### 6. Monad Operations (Chain)
- `BenchmarkMonadChain_Right/Left` - Direct chain
- `BenchmarkChain_Right/Left` - Curried chain
- `BenchmarkChainFirst_Right/Left` - Chain preserving original value
- `BenchmarkFlatten_Right/Left` - Removing nesting

**Key Insights:**
- Chain is very fast: 2-7ns, zero allocations
- **Bottleneck:** ChainFirst on Right: ~168ns with 5 allocations (120B)
- ChainFirst on Left short-circuits efficiently (~9ns)
- Curried Chain is faster than MonadChain

### 7. Applicative Operations (Ap)
- `BenchmarkMonadAp_RightRight/RightLeft/LeftRight` - Direct apply
- `BenchmarkAp_RightRight` - Curried apply

**Key Insights:**
- Ap operations: 5-12ns, zero allocations
- Short-circuits efficiently when either operand is Left
- MonadAp slightly slower than curried Ap

### 8. Alternative Operations
- `BenchmarkAlt_RightRight/LeftRight` - Providing alternatives
- `BenchmarkOrElse_Right/Left` - Recovery from Left

**Key Insights:**
- Very fast: 1.6-4.5ns, zero allocations
- Right path short-circuits without evaluating alternative
- Efficient error recovery mechanism

### 9. Conversion Operations
- `BenchmarkTryCatch_Success/Error` - Converting (value, error) tuples
- `BenchmarkTryCatchError_Success/Error` - Specialized for errors
- `BenchmarkSwap_Right/Left` - Swapping Left/Right
- `BenchmarkGetOrElse_Right/Left` - Extracting with default

**Key Insights:**
- All operations: 0.3-6ns, zero allocations
- Very efficient conversions
- GetOrElse on Right is extremely fast (~0.3ns)

### 10. Pipeline Operations
- `BenchmarkPipeline_Map_Right/Left` - Single Map in pipeline
- `BenchmarkPipeline_Chain_Right/Left` - Single Chain in pipeline
- `BenchmarkPipeline_Complex_Right/Left` - Multiple operations

**Key Insights:**
- **Major Bottleneck:** Pipeline operations allocate heavily
- Single Map: ~100ns with 4 allocations (96B)
- Complex pipeline: ~200-250ns with 8 allocations (192B)
- Chain in pipeline is much faster: 2-4.5ns, zero allocations
- **Recommendation:** Use direct function calls instead of F.Pipe for hot paths

### 11. Sequence Operations
- `BenchmarkMonadSequence2_RightRight/LeftRight` - Combining 2 Eithers
- `BenchmarkMonadSequence3_RightRightRight` - Combining 3 Eithers

**Key Insights:**
- Sequence2: ~6ns, zero allocations
- Sequence3: ~9ns, zero allocations
- Efficient short-circuiting on Left
- Linear scaling with number of Eithers

### 12. Do-Notation Operations
- `BenchmarkDo` - Creating empty context
- `BenchmarkBind_Right` - Binding values to context
- `BenchmarkLet_Right` - Pure computations in context

**Key Insights:**
- Do is extremely fast: ~0.4ns, zero allocations
- **Bottleneck:** Bind: ~130ns with 6 allocations (144B)
- Let is more efficient: ~23ns with 1 allocation (16B)
- **Recommendation:** Prefer Let over Bind when possible

### 13. String Formatting
- `BenchmarkString_Right/Left` - Converting to string

**Key Insights:**
- Right: ~69ns with 1 allocation (16B)
- Left: ~111ns with 1 allocation (48B)
- Only use for debugging, not in hot paths

## Performance Bottlenecks Summary

### Critical Bottlenecks (>100ns or multiple allocations)

1. **Pipeline operations with F.Pipe** (~100-250ns, 4-8 allocations)
   - **Impact:** High - commonly used pattern
   - **Mitigation:** Use direct function calls in hot paths
   - **Example:**
     ```go
     // Slow (100ns, 4 allocs)
     result := F.Pipe1(either, Map[error](transform))
     
     // Fast (12ns, 0 allocs)
     result := Map[error](transform)(either)
     ```

2. **ChainFirst on Right path** (~168ns, 5 allocations)
   - **Impact:** Medium - used for side effects
   - **Mitigation:** Avoid in hot paths, use direct Chain if side effect result not needed
   
3. **Bind in do-notation** (~130ns, 6 allocations)
   - **Impact:** Medium - used in complex workflows
   - **Mitigation:** Use Let for pure computations, minimize Bind calls

### Minor Bottlenecks (20-50ns or 1 allocation)

4. **BiMap operations** (~29-39ns, 1 allocation)
   - **Impact:** Low - less commonly used
   - **Mitigation:** Use Map or MapLeft separately if only one channel needs transformation

5. **MapLeft on Left path** (~34ns, 1 allocation)
   - **Impact:** Low - error path typically not hot
   - **Mitigation:** None needed unless in critical error handling

## Optimization Recommendations

### For Hot Paths

1. **Prefer direct function calls over pipelines:**
   ```go
   // Instead of:
   result := F.Pipe3(either, Map(f1), Chain(f2), Map(f3))
   
   // Use:
   result := Map(f3)(Chain(f2)(Map(f1)(either)))
   ```

2. **Use Chain instead of ChainFirst when possible:**
   ```go
   // If you don't need the original value:
   result := Chain(f)(either)  // Fast
   
   // Instead of:
   result := ChainFirst(func(a A) Either[E, B] { 
       f(a)
       return Right[error](a) 
   })(either)  // Slow
   ```

3. **Prefer Let over Bind in do-notation:**
   ```go
   // For pure computations:
   result := Let(setter, pureFunc)(either)  // 23ns
   
   // Instead of:
   result := Bind(setter, func(s S) Either[E, T] { 
       return Right[error](pureFunc(s)) 
   })(either)  // 130ns
   ```

### For Error Paths

- Left path operations are generally slower but acceptable since errors are exceptional
- Focus optimization on Right (success) path
- Use MapLeft and BiMap freely in error handling code

### Memory Considerations

- Most operations have zero allocations
- Avoid string formatting (String()) in production code
- Pipeline operations allocate - use sparingly in hot paths

## Comparative Analysis

### Fastest Operations (<5ns)
- Constructors (Left, Right, Of)
- Predicates (IsLeft, IsRight)
- GetOrElse on Right
- Chain (curried version)
- Alt/OrElse on Right
- Do

### Medium Speed (5-20ns)
- Fold operations
- Unwrap operations
- Map operations
- MonadChain
- Ap operations
- Sequence operations
- Let

### Slower Operations (>20ns)
- BiMap
- MapLeft on Left
- String formatting
- Pipeline operations
- Bind
- ChainFirst on Right

## Conclusion

The `either` package is highly optimized with most operations completing in nanoseconds with zero allocations. The main performance considerations are:

1. Avoid F.Pipe in hot paths - use direct function calls
2. Prefer Let over Bind in do-notation
3. Use Chain instead of ChainFirst when the original value isn't needed
4. The Right (success) path is consistently faster than Left (error) path
5. Most operations have zero allocations, making them GC-friendly

For typical use cases, the performance is excellent. Only in extremely hot paths (millions of operations per second) should you consider the micro-optimizations suggested above.