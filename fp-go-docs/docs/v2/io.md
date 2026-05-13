---
title: IO
hide_title: true
description: Lazy synchronous side effects with referential transparency - control when effects execute.
sidebar_position: 12
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="IO"
  lede="Lazy, synchronous computation that produces a value. IO[A] encapsulates side effects while maintaining referential transparency and composability."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/io' },
    { label: 'Type', value: 'Monad (func() A)' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

IO is simply a function that takes no arguments and returns a value:

<CodeCard file="type_definition.go">
{`package io

// IO is a lazy computation
type IO[A any] = func() A
`}
</CodeCard>

### Why IO?

<ApiTable>
| Benefit | Description |
|---------|-------------|
| **Lazy evaluation** | Computations don't execute until explicitly called |
| **Referential transparency** | Same IO value always describes the same computation |
| **Composability** | Build complex operations from simple ones without executing |
| **Testability** | Mock side effects by providing different IO values |
| **Explicit effects** | Type system tracks which functions have side effects |
</ApiTable>

<Compare>
<CompareCol kind="bad">
<CodeCard file="eager.go">
{`// ❌ Eager - executes immediately
func getTimestamp() time.Time {
    return time.Now()  // Runs NOW
}

// Hard to test
func processData() Result {
    timestamp := getTimestamp()  // Can't control
    return process(timestamp)
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="lazy.go">
{`// ✅ Lazy with IO - describes computation
func getTimestamp() io.IO[time.Time] {
    return io.Now  // Returns description
}

// Easy to test
func processData() io.IO[Result] {
    return io.Chain(func(t time.Time) io.IO[Result] {
        return io.Of(process(t))
    })(getTimestamp())
}

result := processData()()  // Execute when ready
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[A any](value A) IO[A]` | Wrap pure value in IO |
| `FromImpure` | `func FromImpure(f func()) IO[unit.Unit]` | Wrap side effect |
| `Now` | `IO[time.Time]` | Current time |
| `UnixTime` | `IO[int64]` | Current Unix timestamp |
| `MonotonicTime` | `IO[int64]` | Monotonic time in nanoseconds |
| `Random` | `IO[int]` | Random integer |
| `RandomRange` | `func RandomRange(min, max int) IO[int]` | Random in range |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func(IO[A]) IO[B]` | Transform result |
| `Chain` | `func Chain[A, B any](f func(A) IO[B]) func(IO[A]) IO[B]` | FlatMap - sequence operations |
| `Flatten` | `func Flatten[A any](IO[IO[A]]) IO[A]` | Unwrap nested IO |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[A, B any](fa IO[A]) func(IO[func(A) B]) IO[B]` | Apply function (parallel) |
| `ApSeq` | `func ApSeq[A, B any](fa IO[A]) func(IO[func(A) B]) IO[B]` | Apply function (sequential) |
| `SequenceArray` | `func SequenceArray[A any]([]IO[A]) IO[[]A]` | All-or-nothing (parallel) |
| `SequenceArraySeq` | `func SequenceArraySeq[A any]([]IO[A]) IO[[]A]` | All-or-nothing (sequential) |
| `TraverseArray` | `func TraverseArray[A, B any](f func(A) IO[B]) func([]A) IO[[]B]` | Map and sequence |
</ApiTable>

### Time Operations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Delay` | `func Delay(d time.Duration) func(IO[A]) IO[A]` | Defer execution by duration |
| `After` | `func After(t time.Time) func(IO[A]) IO[A]` | Execute after specific time |
| `WithDuration` | `func WithDuration[A any](IO[A]) IO[pair.Pair[A, time.Duration]]` | Measure execution time |
| `WithTime` | `func WithTime[A any](IO[A]) IO[tuple.Tuple3[A, time.Time, time.Time]]` | Include start/end times |
</ApiTable>

### Resource Management

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Bracket` | `func Bracket[R, A any](acquire IO[R], use func(R) IO[A], release func(R, IO[A]) IO[unit.Unit]) IO[A]` | Safe resource handling |
| `WithResource` | `func WithResource[R, A any](acquire func(...) IO[R], release func(R) IO[unit.Unit]) func(func(R) IO[A]) IO[A]` | Resource pattern |
</ApiTable>

### Utilities

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `ChainFirst` | `func ChainFirst[A, B any](f IO[B]) func(IO[A]) IO[A]` | Side effect, keep original value |
| `Logger` | `func Logger() func(string) IO[unit.Unit]` | Log message |
| `Printf` | `func Printf(format string) func(...any) IO[unit.Unit]` | Printf-style logging |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    "time"
    IO "github.com/IBM/fp-go/v2/io"
)

func main() {
    // Wrap pure value
    greeting := IO.Of("Hello, World!")
    result := greeting()  // "Hello, World!"
    
    // Current time (lazy)
    now := IO.Now
    timestamp := now()  // time.Time
    
    // Random number
    randomNum := IO.Random
    n := randomNum()  // int
    
    // Side effect
    printHello := IO.FromImpure(func() {
        fmt.Println("Hello!")
    })
    printHello()  // Prints "Hello!"
}
`}
</CodeCard>

### Transformations

<CodeCard file="transformations.go">
{`package main

import (
    "fmt"
    "time"
    IO "github.com/IBM/fp-go/v2/io"
    F "github.com/IBM/fp-go/v2/function"
)

func main() {
    // Map: transform result
    doubled := F.Pipe1(
        IO.Of(21),
        IO.Map(func(n int) int { return n * 2 }),
    )
    result := doubled()  // 42
    
    // Chain: sequence operations
    formatted := F.Pipe2(
        IO.Now,
        IO.Chain(func(t time.Time) IO.IO[int64] {
            return IO.Of(t.Unix())
        }),
        IO.Chain(func(unix int64) IO.IO[string] {
            return IO.Of(fmt.Sprintf("Timestamp: %d", unix))
        }),
    )
    output := formatted()  // "Timestamp: 1234567890"
}
`}
</CodeCard>

### Parallel vs Sequential

<CodeCard file="parallel.go">
{`package main

import (
    "time"
    IO "github.com/IBM/fp-go/v2/io"
)

func main() {
    operations := []IO.IO[int]{
        IO.Delay(100*time.Millisecond)(IO.Of(1)),
        IO.Delay(100*time.Millisecond)(IO.Of(2)),
        IO.Delay(100*time.Millisecond)(IO.Of(3)),
    }
    
    // Parallel execution (~100ms total)
    parallel := IO.SequenceArray(operations)
    results := parallel()  // [1, 2, 3]
    
    // Sequential execution (~300ms total)
    sequential := IO.SequenceArraySeq(operations)
    results = sequential()  // [1, 2, 3]
}
`}
</CodeCard>

### Resource Management

<CodeCard file="resource.go">
{`package main

import (
    "os"
    IO "github.com/IBM/fp-go/v2/io"
)

func processFile(path string) IO.IO[[]byte] {
    return IO.Bracket(
        // Acquire resource
        func() IO.IO[*os.File] {
            return IO.Of(func() *os.File {
                f, _ := os.Open(path)
                return f
            }())
        },
        // Use resource
        func(f *os.File) IO.IO[[]byte] {
            return IO.Of(func() []byte {
                data, _ := io.ReadAll(f)
                return data
            }())
        },
        // Release resource (always runs)
        func(f *os.File, _ IO.IO[[]byte]) IO.IO[unit.Unit] {
            return IO.FromImpure(func() {
                f.Close()
            })
        },
    )
}

func main() {
    data := processFile("config.json")()
    // File is guaranteed to be closed
}
`}
</CodeCard>

### Time-Based Operations

<CodeCard file="time_ops.go">
{`package main

import (
    "fmt"
    "time"
    IO "github.com/IBM/fp-go/v2/io"
)

func main() {
    // Delay execution
    delayed := IO.Delay(time.Second)(IO.Of(42))
    result := delayed()  // Waits 1 second, returns 42
    
    // Measure execution time
    operation := IO.Delay(100 * time.Millisecond)(IO.Of(42))
    withTime := IO.WithDuration(operation)
    value, duration := withTime()
    fmt.Printf("Value: %d, Duration: %v\n", value, duration)
    // Value: 42, Duration: ~100ms
}
`}
</CodeCard>

### Logging and Debugging

<CodeCard file="logging.go">
{`package main

import (
    IO "github.com/IBM/fp-go/v2/io"
    F "github.com/IBM/fp-go/v2/function"
)

func fetchUser(id string) IO.IO[User] {
    return F.Pipe2(
        IO.ChainFirst(IO.Logger()("Fetching user...")),
        fetchUserFromDB(id),
        IO.ChainFirst(IO.Printf("Fetched user: %+v")),
    )
}

func main() {
    user := fetchUser("123")()
    // Logs: "Fetching user..."
    // Logs: "Fetched user: {ID:123 Name:Alice}"
    // Returns: User{ID: "123", Name: "Alice"}
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: API Calls

<CodeCard file="api_calls.go">
{`package main

import (
    IO "github.com/IBM/fp-go/v2/io"
    F "github.com/IBM/fp-go/v2/function"
)

func fetchUserData(id string) IO.IO[UserData] {
    return F.Pipe2(
        fetchUser(id),  // IO.IO[User]
        IO.Chain(func(user User) IO.IO[UserData] {
            return IO.Map(func(posts []Post) UserData {
                return UserData{User: user, Posts: posts}
            })(fetchPosts(user.ID))
        }),
    )
}

// Execute when ready
data := fetchUserData("123")()
`}
</CodeCard>

### Pattern 2: Caching

<CodeCard file="caching.go">
{`package main

import (
    "sync"
    IO "github.com/IBM/fp-go/v2/io"
)

var cachedData IO.IO[Data]
var once sync.Once

func getCachedData() IO.IO[Data] {
    return func() Data {
        once.Do(func() {
            cachedData = expensiveComputation()
        })
        return cachedData()
    }
}

// First call computes, subsequent calls use cache
data1 := getCachedData()()
data2 := getCachedData()()  // Uses cached value
`}
</CodeCard>

### Pattern 3: Testing with Mocks

<CodeCard file="testing.go">
{`package main

import (
    "testing"
    "time"
    IO "github.com/IBM/fp-go/v2/io"
)

type Dependencies struct {
    GetTime   func() IO.IO[time.Time]
    FetchUser func(string) IO.IO[User]
}

func processUser(deps Dependencies, id string) IO.IO[Result] {
    return F.Pipe2(
        deps.GetTime(),
        IO.Chain(func(t time.Time) IO.IO[Result] {
            return IO.Map(func(u User) Result {
                return Result{User: u, Timestamp: t}
            })(deps.FetchUser(id))
        }),
    )
}

func TestProcessUser(t *testing.T) {
    // Mock dependencies
    deps := Dependencies{
        GetTime: func() IO.IO[time.Time] {
            return IO.Of(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
        },
        FetchUser: func(id string) IO.IO[User] {
            return IO.Of(User{ID: id, Name: "Test User"})
        },
    }
    
    result := processUser(deps, "123")()
    
    assert.Equal(t, "Test User", result.User.Name)
    assert.Equal(t, 2024, result.Timestamp.Year())
}
`}
</CodeCard>

### When to Use IO vs IOResult

<ApiTable>
| Use IO When | Use IOResult When |
|-------------|-------------------|
| Side effects that cannot fail | File operations that may fail |
| Time-based operations | HTTP requests |
| Random number generation | Database queries |
| Logging and debugging | Any operation that can error |
| Need lazy evaluation | Need error handling + lazy evaluation |
</ApiTable>

</Section>
