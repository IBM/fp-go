---
title: Parallel Tasks
hide_title: true
description: Execute tasks concurrently using functional patterns with fp-go's ApplicativePar, Traverse, and Race combinators.
sidebar_position: 12
---

<PageHeader
  eyebrow="Recipes · 12 / 17"
  title="Parallel"
  titleAccent="Tasks"
  lede="Execute tasks concurrently using functional patterns with fp-go's ApplicativePar, Traverse, and Race combinators."
  meta={[
    { label: 'Difficulty', value: 'Advanced' },
    { label: 'Patterns', value: '7' },
    { label: 'Use Cases', value: 'Concurrent processing, HTTP requests, worker pools' }
  ]}
/>

<TLDR>
  <TLDRCard title="Type-Safe Parallelism" icon="zap">
    Use ApplicativePar for compile-time safe concurrent operations without manual goroutine management.
  </TLDRCard>
  <TLDRCard title="Controlled Concurrency" icon="sliders">
    Worker pools and rate limiters prevent resource exhaustion while maximizing throughput.
  </TLDRCard>
  <TLDRCard title="Race & Timeout" icon="clock">
    Race multiple operations and implement timeouts to handle slow or failing services gracefully.
  </TLDRCard>
</TLDR>

<Section id="basic-parallel" number="01" title="Basic Parallel" titleAccent="Execution">

Parallel execution allows multiple operations to run concurrently. fp-go provides **ApplicativePar** for safe parallel operations.

<CodeCard file="parallel_basic.go">
{`package main

import (
    "fmt"
    "time"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func slowOperation(n int) IOE.IOEither[error, int] {
    return IOE.TryCatch(func() (int, error) {
        time.Sleep(time.Duration(n) * time.Second)
        return n * 2, nil
    })
}

func processSequential(numbers []int) IOE.IOEither[error, []int] {
    return A.Traverse[int](IOE.Applicative[error, int]())(
        slowOperation,
    )(numbers)
}

func processParallel(numbers []int) IOE.IOEither[error, []int] {
    return A.Traverse[int](IOE.ApplicativePar[error, int]())(
        slowOperation,
    )(numbers)
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Sequential: ~15 seconds (1+2+3+4+5)
    start := time.Now()
    result1 := processSequential(numbers)()
    fmt.Printf("Sequential: %v (took %v)\\n", result1.Right(), time.Since(start))
    
    // Parallel: ~5 seconds (max of all)
    start = time.Now()
    result2 := processParallel(numbers)()
    fmt.Printf("Parallel: %v (took %v)\\n", result2.Right(), time.Since(start))
}
`}
</CodeCard>

<Callout type="info">
**ApplicativePar vs Applicative**: Use `ApplicativePar` for concurrent execution, `Applicative` for sequential. The API is identical—only the execution strategy differs.
</Callout>

</Section>

<Section id="http-requests" number="02" title="Concurrent HTTP" titleAccent="Requests">

Fetch multiple URLs in parallel to reduce total request time.

<CodeCard file="parallel_http.go">
{`package main

import (
    "fmt"
    "io"
    "net/http"
    "time"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type Response struct {
    URL    string
    Status int
    Body   string
}

func fetchURL(url string) IOE.IOEither[error, Response] {
    return IOE.TryCatch(func() (Response, error) {
        resp, err := http.Get(url)
        if err != nil {
            return Response{}, err
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return Response{}, err
        }
        
        return Response{
            URL:    url,
            Status: resp.StatusCode,
            Body:   string(body),
        }, nil
    })
}

func fetchAllParallel(urls []string) IOE.IOEither[error, []Response] {
    return A.Traverse[string](IOE.ApplicativePar[error, Response]())(
        fetchURL,
    )(urls)
}

func main() {
    urls := []string{
        "https://api.github.com/users/octocat",
        "https://api.github.com/users/torvalds",
        "https://api.github.com/users/gvanrossum",
    }
    
    start := time.Now()
    result := fetchAllParallel(urls)()
    duration := time.Since(start)
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        responses := result.Right()
        fmt.Printf("Fetched %d URLs in %v\\n", len(responses), duration)
        for _, resp := range responses {
            fmt.Printf("  %s: %d (%d bytes)\\n", resp.URL, resp.Status, len(resp.Body))
        }
    }
}
`}
</CodeCard>

</Section>

<Section id="worker-pool" number="03" title="Worker Pool" titleAccent="Pattern">

Control concurrency level with a fixed worker pool to prevent resource exhaustion.

<CodeCard file="worker_pool.go">
{`package main

import (
    "fmt"
    "sync"
    "time"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type WorkerPool[A, B any] struct {
    workers int
    work    func(A) IOE.IOEither[error, B]
}

func NewWorkerPool[A, B any](
    workers int,
    work func(A) IOE.IOEither[error, B],
) *WorkerPool[A, B] {
    return &WorkerPool[A, B]{
        workers: workers,
        work:    work,
    }
}

func (wp *WorkerPool[A, B]) Execute(items []A) IOE.IOEither[error, []B] {
    return func() IOE.Either[error, []B] {
        jobs := make(chan A, len(items))
        results := make(chan IOE.Either[error, B], len(items))
        
        // Start workers
        var wg sync.WaitGroup
        for i := 0; i < wp.workers; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for item := range jobs {
                    results <- wp.work(item)()
                }
            }()
        }
        
        // Send jobs
        for _, item := range items {
            jobs <- item
        }
        close(jobs)
        
        // Wait for completion
        go func() {
            wg.Wait()
            close(results)
        }()
        
        // Collect results
        output := make([]B, 0, len(items))
        for result := range results {
            if result.IsLeft() {
                return IOE.Left[[]B](result.Left())()
            }
            output = append(output, result.Right())
        }
        
        return IOE.Right[error](output)()
    }
}

func processItem(n int) IOE.IOEither[error, int] {
    return IOE.TryCatch(func() (int, error) {
        time.Sleep(100 * time.Millisecond)
        return n * 2, nil
    })
}

func main() {
    items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    pool := NewWorkerPool(3, processItem)
    
    start := time.Now()
    result := pool.Execute(items)()
    duration := time.Since(start)
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Processed %d items in %v with 3 workers\\n", len(result.Right()), duration)
        fmt.Println("Results:", result.Right())
    }
}
`}
</CodeCard>

</Section>

<Section id="race-timeout" number="04" title="Race &" titleAccent="Timeout">

Race multiple operations and implement timeouts for resilient systems.

<CodeCard file="race_timeout.go">
{`package main

import (
    "fmt"
    "sync"
    "time"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func race[A any](operations []IOE.IOEither[error, A]) IOE.IOEither[error, A] {
    return func() IOE.Either[error, A] {
        resultChan := make(chan IOE.Either[error, A], len(operations))
        var wg sync.WaitGroup
        
        for _, op := range operations {
            wg.Add(1)
            go func(operation IOE.IOEither[error, A]) {
                defer wg.Done()
                resultChan <- operation()
            }(op)
        }
        
        go func() {
            wg.Wait()
            close(resultChan)
        }()
        
        // Return first successful result
        var lastErr error
        for result := range resultChan {
            if result.IsRight() {
                return result
            }
            lastErr = result.Left()
        }
        
        return IOE.Left[A](fmt.Errorf("all operations failed: %w", lastErr))()
    }
}

func withTimeout[A any](
    timeout time.Duration,
    operation IOE.IOEither[error, A],
) IOE.IOEither[error, A] {
    timeoutOp := IOE.TryCatch(func() (A, error) {
        time.Sleep(timeout)
        var zero A
        return zero, fmt.Errorf("operation timed out after %v", timeout)
    })
    
    return race([]IOE.IOEither[error, A]{operation, timeoutOp})
}

func slowOperation() IOE.IOEither[error, string] {
    return IOE.TryCatch(func() (string, error) {
        time.Sleep(2 * time.Second)
        return "Completed", nil
    })
}

func main() {
    // Will timeout
    result1 := withTimeout(1*time.Second, slowOperation())()
    if result1.IsLeft() {
        fmt.Println("Result 1:", result1.Left())
    }
    
    // Will succeed
    result2 := withTimeout(3*time.Second, slowOperation())()
    if result2.IsRight() {
        fmt.Println("Result 2:", result2.Right())
    }
}
`}
</CodeCard>

</Section>

<Section id="batch-processing" number="05" title="Batch" titleAccent="Processing">

Process large datasets in parallel batches for optimal throughput.

<CodeCard file="batch_processing.go">
{`package main

import (
    "fmt"
    "time"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func processInBatches[A, B any](
    batchSize int,
    process func([]A) IOE.IOEither[error, []B],
) func([]A) IOE.IOEither[error, []B] {
    return func(items []A) IOE.IOEither[error, []B] {
        chunks := A.Chunksof(batchSize)(items)
        
        return F.Pipe2(
            A.Traverse[[]A](IOE.ApplicativePar[error, []B]())(process)(chunks),
            IOE.Map(func(results [][]B) []B {
                return A.Flatten(results)
            }),
        )
    }
}

func processBatch(items []int) IOE.IOEither[error, []int] {
    return IOE.TryCatch(func() ([]int, error) {
        time.Sleep(100 * time.Millisecond)
        return A.Map(func(n int) int { return n * 2 })(items), nil
    })
}

func main() {
    items := A.MakeBy(100)(func(i int) int { return i + 1 })
    
    processor := processInBatches(10, processBatch)
    
    start := time.Now()
    result := processor(items)()
    duration := time.Since(start)
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Processed %d items in %v\\n", len(result.Right()), duration)
    }
}
`}
</CodeCard>

</Section>

<Section id="error-handling" number="06" title="Error Handling in" titleAccent="Parallel">

Handle errors gracefully in parallel operations with fail-fast or collect-all strategies.

<Compare>
<CompareCol kind="bad">
<CodeCard file="fail_fast.go">
{`// Fail fast: stop on first error
func processWithFailFast(items []int) IOE.IOEither[error, []int] {
    return A.Traverse[int](IOE.ApplicativePar[error, int]())(
        func(n int) IOE.IOEither[error, int] {
            if n == 5 {
                return IOE.Left[int](fmt.Errorf("failed on item %d", n))
            }
            return IOE.Right[error](n * 2)
        },
    )(items)
}

func main() {
    items := []int{1, 2, 3, 4, 5, 6, 7, 8}
    result := processWithFailFast(items)()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
        // Error: failed on item 5
    }
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="collect_errors.go">
{`// Collect all: process everything, track errors
type Result[A any] struct {
    Value A
    Error error
}

func processAllWithErrors[A, B any](
    items []A,
    process func(A) IOE.IOEither[error, B],
) []Result[B] {
    results := make([]Result[B], len(items))
    var wg sync.WaitGroup
    
    for i, item := range items {
        wg.Add(1)
        go func(idx int, val A) {
            defer wg.Done()
            result := process(val)()
            if result.IsLeft() {
                results[idx] = Result[B]{Error: result.Left()}
            } else {
                results[idx] = Result[B]{Value: result.Right()}
            }
        }(i, item)
    }
    
    wg.Wait()
    return results
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>

<Section id="best-practices" number="07" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Control concurrency level** — Use worker pools to limit goroutines and prevent resource exhaustion
  </ChecklistItem>
  <ChecklistItem status="required">
    **Set timeouts** — Always timeout external operations to prevent hanging forever
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Handle partial failures** — Collect all results when appropriate instead of failing fast
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Batch small operations** — Reduce overhead by batching tiny operations together
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Use context for cancellation** — Pass context.Context for graceful shutdown
  </ChecklistItem>
</Checklist>

<CodeCard file="best_practices.go">
{`// ✅ Good: Controlled concurrency with timeout
pool := NewWorkerPool(10, processItem)
result := withTimeout(5*time.Second, pool.Execute(items))

// ❌ Avoid: Unbounded parallelism without timeout
result := A.Traverse[int](IOE.ApplicativePar[error, int]())(
    processItem,
)(thousandsOfItems) // Can overwhelm system
`}
</CodeCard>

</Section>
