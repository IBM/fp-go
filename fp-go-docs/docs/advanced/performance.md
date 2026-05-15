---
title: Performance Optimization
hide_title: true
description: Optimize functional code with lazy evaluation, memoization, tail call optimization, parallel processing, and memory management techniques.
sidebar_position: 3
---

<PageHeader
  eyebrow="Advanced · 03 / 04"
  title="Performance"
  titleAccent="Optimization"
  lede="Optimize functional code with lazy evaluation, memoization, tail call optimization, parallel processing, structure sharing, and memory management."
  meta={[
    { label: 'Difficulty', value: 'Advanced' },
    { label: 'Topics', value: '8' },
    { label: 'Prerequisites', value: 'Profiling, Benchmarking' }
  ]}
/>

<TLDR>
  <TLDRCard title="Lazy Evaluation" icon="clock">
    Defer computation until needed—cache results and process infinite sequences efficiently.
  </TLDRCard>
  <TLDRCard title="Memoization" icon="database">
    Cache expensive function results with LRU eviction to balance memory and speed.
  </TLDRCard>
  <TLDRCard title="Parallel Processing" icon="zap">
    Leverage concurrency with worker pools and parallel map for CPU-bound operations.
  </TLDRCard>
</TLDR>

<Section id="lazy-evaluation" number="01" title="Lazy" titleAccent="Evaluation">

Defer computation until the value is actually needed.

<CodeCard file="lazy_values.go">
{`package main

import (
    "fmt"
    "time"
)

type Lazy[A any] struct {
    compute func() A
    cached  *A
}

func NewLazy[A any](f func() A) *Lazy[A] {
    return &Lazy[A]{compute: f}
}

func (l *Lazy[A]) Force() A {
    if l.cached == nil {
        value := l.compute()
        l.cached = &value
    }
    return *l.cached
}

func main() {
    // Expensive computation
    expensive := NewLazy(func() int {
        fmt.Println("Computing...")
        time.Sleep(100 * time.Millisecond)
        return 42
    })
    
    fmt.Println("Lazy value created")
    
    // First force - computes
    result1 := expensive.Force()
    fmt.Println("Result:", result1)
    
    // Second force - uses cache
    result2 := expensive.Force()
    fmt.Println("Result:", result2)
    
    // Output:
    // Lazy value created
    // Computing...
    // Result: 42
    // Result: 42
}
`}
</CodeCard>

<CodeCard file="lazy_sequences.go">
{`package main

import (
    "fmt"
)

type LazySeq[A any] struct {
    head func() A
    tail func() *LazySeq[A]
}

func Take[A any](n int, seq *LazySeq[A]) []A {
    result := make([]A, 0, n)
    current := seq
    
    for i := 0; i < n && current != nil; i++ {
        result = append(result, current.head())
        current = current.tail()
    }
    
    return result
}

func Naturals() *LazySeq[int] {
    var generate func(int) *LazySeq[int]
    generate = func(n int) *LazySeq[int] {
        return &LazySeq[int]{
            head: func() int { return n },
            tail: func() *LazySeq[int] { return generate(n + 1) },
        }
    }
    return generate(0)
}

func main() {
    // Infinite sequence of natural numbers
    nats := Naturals()
    
    // Take only what we need
    first10 := Take(10, nats)
    fmt.Println(first10) // [0 1 2 3 4 5 6 7 8 9]
}
`}
</CodeCard>

</Section>

<Section id="memoization" number="02" title="Memoization" titleAccent="& Caching">

Cache function results to avoid recomputation.

<CodeCard file="simple_memoization.go">
{`package main

import (
    "fmt"
    "sync"
)

func Memoize[K comparable, V any](f func(K) V) func(K) V {
    cache := make(map[K]V)
    var mu sync.RWMutex
    
    return func(k K) V {
        mu.RLock()
        if v, ok := cache[k]; ok {
            mu.RUnlock()
            return v
        }
        mu.RUnlock()
        
        mu.Lock()
        defer mu.Unlock()
        
        // Double-check after acquiring write lock
        if v, ok := cache[k]; ok {
            return v
        }
        
        v := f(k)
        cache[k] = v
        return v
    }
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    // Memoized fibonacci
    memoFib := Memoize(fibonacci)
    
    fmt.Println(memoFib(10)) // Computed
    fmt.Println(memoFib(10)) // Cached
    fmt.Println(memoFib(15)) // Computed (uses cached 10)
}
`}
</CodeCard>

<CodeCard file="lru_cache.go">
{`package main

import (
    "container/list"
    "fmt"
    "sync"
)

type LRUCache[K comparable, V any] struct {
    capacity int
    cache    map[K]*list.Element
    lru      *list.List
    mu       sync.Mutex
}

type entry[K comparable, V any] struct {
    key   K
    value V
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
    return &LRUCache[K, V]{
        capacity: capacity,
        cache:    make(map[K]*list.Element),
        lru:      list.New(),
    }
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.cache[key]; ok {
        c.lru.MoveToFront(elem)
        return elem.Value.(*entry[K, V]).value, true
    }
    
    var zero V
    return zero, false
}

func (c *LRUCache[K, V]) Put(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.cache[key]; ok {
        c.lru.MoveToFront(elem)
        elem.Value.(*entry[K, V]).value = value
        return
    }
    
    if c.lru.Len() >= c.capacity {
        oldest := c.lru.Back()
        if oldest != nil {
            c.lru.Remove(oldest)
            delete(c.cache, oldest.Value.(*entry[K, V]).key)
        }
    }
    
    elem := c.lru.PushFront(&entry[K, V]{key: key, value: value})
    c.cache[key] = elem
}

func main() {
    cache := NewLRUCache[string, int](2)
    
    cache.Put("a", 1)
    cache.Put("b", 2)
    
    if v, ok := cache.Get("a"); ok {
        fmt.Println("a:", v) // 1
    }
    
    cache.Put("c", 3) // Evicts "b"
    
    if _, ok := cache.Get("b"); !ok {
        fmt.Println("b evicted")
    }
}
`}
</CodeCard>

</Section>

<Section id="tail-call" number="03" title="Tail Call" titleAccent="Optimization">

Go doesn't have TCO, but trampolines enable stack-safe recursion.

<CodeCard file="trampoline.go">
{`package main

import (
    "fmt"
)

type Trampoline[A any] interface {
    isTrampoline()
}

type Done[A any] struct {
    Value A
}

func (Done[A]) isTrampoline() {}

type More[A any] struct {
    Next func() Trampoline[A]
}

func (More[A]) isTrampoline() {}

func Run[A any](t Trampoline[A]) A {
    for {
        switch v := t.(type) {
        case Done[A]:
            return v.Value
        case More[A]:
            t = v.Next()
        }
    }
}

// Tail-recursive factorial using trampoline
func factorialTR(n, acc int) Trampoline[int] {
    if n <= 1 {
        return Done[int]{Value: acc}
    }
    return More[int]{
        Next: func() Trampoline[int] {
            return factorialTR(n-1, n*acc)
        },
    }
}

func factorial(n int) int {
    return Run(factorialTR(n, 1))
}

func main() {
    result := factorial(10)
    fmt.Println(result) // 3628800
}
`}
</CodeCard>

</Section>

<Section id="parallel-processing" number="04" title="Parallel" titleAccent="Processing">

Leverage concurrency for CPU-bound operations.

<CodeCard file="parallel_map.go">
{`package main

import (
    "fmt"
    "runtime"
    "sync"
)

func ParallelMap[A, B any](f func(A) B, items []A) []B {
    numWorkers := runtime.NumCPU()
    results := make([]B, len(items))
    
    var wg sync.WaitGroup
    chunkSize := (len(items) + numWorkers - 1) / numWorkers
    
    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > len(items) {
            end = len(items)
        }
        if start >= len(items) {
            break
        }
        
        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()
            for j := start; j < end; j++ {
                results[j] = f(items[j])
            }
        }(start, end)
    }
    
    wg.Wait()
    return results
}

func main() {
    numbers := make([]int, 1000)
    for i := range numbers {
        numbers[i] = i
    }
    
    // Parallel map
    squared := ParallelMap(func(n int) int {
        return n * n
    }, numbers)
    
    fmt.Println("First 10:", squared[:10])
}
`}
</CodeCard>

<CodeCard file="worker_pool.go">
{`package main

import (
    "fmt"
    "sync"
)

type WorkerPool[T, R any] struct {
    workers int
    jobs    chan T
    results chan R
    wg      sync.WaitGroup
}

func NewWorkerPool[T, R any](workers int, process func(T) R) *WorkerPool[T, R] {
    pool := &WorkerPool[T, R]{
        workers: workers,
        jobs:    make(chan T, workers*2),
        results: make(chan R, workers*2),
    }
    
    for i := 0; i < workers; i++ {
        pool.wg.Add(1)
        go func() {
            defer pool.wg.Done()
            for job := range pool.jobs {
                pool.results <- process(job)
            }
        }()
    }
    
    return pool
}

func (p *WorkerPool[T, R]) Submit(job T) {
    p.jobs <- job
}

func (p *WorkerPool[T, R]) Close() {
    close(p.jobs)
    p.wg.Wait()
    close(p.results)
}

func (p *WorkerPool[T, R]) Results() <-chan R {
    return p.results
}

func main() {
    pool := NewWorkerPool(4, func(n int) int {
        return n * n
    })
    
    // Submit jobs
    go func() {
        for i := 0; i < 10; i++ {
            pool.Submit(i)
        }
        pool.Close()
    }()
    
    // Collect results
    for result := range pool.Results() {
        fmt.Println(result)
    }
}
`}
</CodeCard>

</Section>

<Section id="structure-sharing" number="05" title="Structure" titleAccent="Sharing">

Reuse data structures instead of copying for memory efficiency.

<CodeCard file="persistent_list.go">
{`package main

import (
    "fmt"
)

// Persistent list (shares structure)
type PList[A any] interface {
    isPList()
}

type Nil[A any] struct{}

func (Nil[A]) isPList() {}

type Cons[A any] struct {
    Head A
    Tail PList[A]
}

func (Cons[A]) isPList() {}

// Prepend shares the tail
func Prepend[A any](head A, tail PList[A]) PList[A] {
    return Cons[A]{Head: head, Tail: tail}
}

func main() {
    // Original list: [1, 2, 3]
    list1 := Prepend(1, Prepend(2, Prepend(3, Nil[int]{})))
    
    // New list shares structure: [0, 1, 2, 3]
    list2 := Prepend(0, list1)
    
    // Both lists exist, sharing [1, 2, 3]
    fmt.Println("Lists created with structure sharing")
    _ = list2
}
`}
</CodeCard>

</Section>

<Section id="fusion" number="06" title="Fusion" titleAccent="Optimization">

Combine multiple operations into one pass to reduce allocations.

<Compare>
<CompareCol kind="bad">
<CodeCard file="naive_map.go">
{`// ❌ Naive: two passes, intermediate allocation
func naiveMapMap[A, B, C any](f func(A) B, g func(B) C, items []A) []C {
    temp := make([]B, len(items))
    for i, item := range items {
        temp[i] = f(item)
    }
    
    result := make([]C, len(temp))
    for i, item := range temp {
        result[i] = g(item)
    }
    
    return result
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="fused_map.go">
{`// ✅ Fused: one pass, no intermediate allocation
func fusedMapMap[A, B, C any](f func(A) B, g func(B) C, items []A) []C {
    result := make([]C, len(items))
    for i, item := range items {
        result[i] = g(f(item))
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    double := func(n int) int { return n * 2 }
    addTen := func(n int) int { return n + 10 }
    
    // Fused: [1,2,3,4,5] -> [12,14,16,18,20]
    result := fusedMapMap(double, addTen, numbers)
    fmt.Println(result) // [12 14 16 18 20]
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>

<Section id="memory-optimization" number="07" title="Memory" titleAccent="Optimization">

Manage memory efficiently with pooling and preallocation.

<CodeCard file="object_pooling.go">
{`package main

import (
    "fmt"
    "sync"
)

type Pool[T any] struct {
    pool sync.Pool
    new  func() T
}

func NewPool[T any](new func() T) *Pool[T] {
    return &Pool[T]{
        pool: sync.Pool{
            New: func() any {
                return new()
            },
        },
        new: new,
    }
}

func (p *Pool[T]) Get() T {
    return p.pool.Get().(T)
}

func (p *Pool[T]) Put(item T) {
    p.pool.Put(item)
}

func main() {
    // Pool of byte slices
    bufferPool := NewPool(func() []byte {
        return make([]byte, 1024)
    })
    
    // Get buffer from pool
    buf := bufferPool.Get()
    fmt.Println("Buffer size:", len(buf))
    
    // Return to pool
    bufferPool.Put(buf)
}
`}
</CodeCard>

<Compare>
<CompareCol kind="bad">
<CodeCard file="bad_append.go">
{`// ❌ Bad: grows dynamically, multiple allocations
func badAppend(n int) []int {
    var result []int
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="good_append.go">
{`// ✅ Good: preallocated, single allocation
func goodAppend(n int) []int {
    result := make([]int, 0, n)
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>

<Section id="best-practices" number="08" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Profile before optimizing** — Use pprof to identify actual bottlenecks, not assumed ones
  </ChecklistItem>
  <ChecklistItem status="required">
    **Benchmark changes** — Measure performance impact with `go test -bench` before and after
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Use appropriate data structures** — Maps for lookups, slices for sequential access, channels for communication
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Preallocate when size is known** — Avoid dynamic growth with `make([]T, 0, capacity)`
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Avoid premature optimization** — Start simple, optimize only proven bottlenecks
  </ChecklistItem>
</Checklist>

<CodeCard file="profiling.go">
{`import _ "net/http/pprof"
import "net/http"

func main() {
    // Enable profiling endpoint
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // Your application code
    // Access profiles at http://localhost:6060/debug/pprof/
}
`}
</CodeCard>

</Section>
