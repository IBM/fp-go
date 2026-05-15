---
sidebar_position: 2
title: 5-Minute Quickstart
hide_title: true
description: Install fp-go, write your first program, and learn the core composition patterns in 5 minutes.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Getting started · Section 03 / 03"
  title="5-minute"
  titleAccent="quickstart."
  lede="Install fp-go, write your first program, and learn the core composition patterns — Pipe, Map, Chain, GetOrElse — in one sitting."
  meta={[
    {label: '// Version', value: <>v2.2.82 <MetaPill>LATEST</MetaPill></>},
    {label: '// Go required', value: '1.24+ (v2)'},
    {label: '// Reading time', value: '5 min · 5 steps'},
  ]}
/>

<TLDR>
  <TLDRCard label="// You'll learn" value="5" unit="patterns" description="Pipe, Map, Chain, GetOrElse, automatic error propagation." />
  <TLDRCard label="// Prereqs" prose value={<>Go <em>1.24+</em> and basic Go familiarity.</>} />
  <TLDRCard label="// You'll build" prose value={<>A safe divider and a <em>composable</em> calculator.</>} />
</TLDR>

<Section
  id="prereqs"
  number="01"
  title="Prerequisites"
  tag="Difficulty · Beginner"
>

<Callout title="Before you start.">
  <ul>
    <li><strong>Go 1.24+</strong> for v2 (recommended)</li>
    <li><strong>Go 1.18+</strong> for v1 (legacy)</li>
    <li>Basic understanding of Go</li>
  </ul>
</Callout>

</Section>

<Section
  id="install"
  number="02"
  title="Install"
  titleAccent="fp-go."
  lede="Choose your version. v2 is recommended for any new project."
>

<Tabs>
  <TabItem value="v2" label="v2 (Recommended)" default>

<CodeCard file="shell" status="tested">
{`# Initialize your Go module (if not already done)
go mod init myapp

# Install fp-go v2
go get github.com/IBM/fp-go/v2`}
</CodeCard>

<Callout type="success" title="Why v2.">
  <ul>
    <li>Latest features (Result, Effect, idiomatic packages)</li>
    <li>Better type inference</li>
    <li>Actively maintained</li>
    <li>Requires Go 1.24+</li>
  </ul>
</Callout>

  </TabItem>
  <TabItem value="v1" label="v1 (Legacy)">

<CodeCard file="shell">
{`# Initialize your Go module (if not already done)
go mod init myapp

# Install fp-go v1
go get github.com/IBM/fp-go`}
</CodeCard>

<Callout type="warn" title="When to choose v1.">
  <ul>
    <li>Stuck on Go 1.18–1.23</li>
    <li>Existing v1 codebase</li>
    <li>Need Writer monad (v1 only)</li>
  </ul>
</Callout>

  </TabItem>
</Tabs>

</Section>

<Section
  id="first-program"
  number="03"
  title="Your first"
  titleAccent="program."
  lede="A safe divider. Compare the three approaches side-by-side."
>

<Tabs>
  <TabItem value="without" label="Without fp-go">

<CodeCard file="main.go">
{`package main

import (
    "errors"
    "fmt"
)

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    // Manual error handling
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    doubled := result * 2
    fmt.Println("Result:", doubled) // Output: Result: 10

    // What if we want to chain more operations?
    // More if err != nil checks...
}`}
</CodeCard>

<Compare>
  <CompareCol kind="bad" pill="problems">
    <p>Repetitive error checking.</p>
    <p>Hard to compose operations.</p>
    <p>Error handling mixed with business logic.</p>
  </CompareCol>
  <CompareCol kind="good" pill="see right tab" title="fp-go">
    <p>Automatic error propagation.</p>
    <p>Composable pipelines.</p>
    <p>Business logic stays clear.</p>
  </CompareCol>
</Compare>

  </TabItem>
  <TabItem value="v2" label="With v2" default>

<CodeCard file="main.go" status="tested">
{`package main

import (
    "errors"
    "fmt"
    "github.com/IBM/fp-go/v2/function"
    "github.com/IBM/fp-go/v2/result"
)

// Pure function that returns Result instead of (value, error)
func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}

func main() {
    // Compose operations with automatic error handling
    finalResult := function.Pipe2(
        divide(10, 2),
        result.Map(func(x int) int { return x * 2 }),
        result.GetOrElse(func() int { return 0 }),
    )

    fmt.Println("Result:", finalResult) // Output: Result: 10

    // Error case is handled automatically
    errorResult := function.Pipe2(
        divide(10, 0), // Returns Err
        result.Map(func(x int) int { return x * 2 }), // Skipped!
        result.GetOrElse(func() int { return 0 }), // Returns default
    )

    fmt.Println("Error result:", errorResult) // Output: Error result: 0
}`}
</CodeCard>

<Callout type="success" title="Benefits.">
  No repetitive error checking. Easy composition. Errors propagate automatically. Business logic stays prominent.
</Callout>

  </TabItem>
  <TabItem value="v1" label="With v1">

<CodeCard file="main.go">
{`package main

import (
    "errors"
    "fmt"
    "github.com/IBM/fp-go/either"
    "github.com/IBM/fp-go/function"
)

// Pure function that returns Either instead of (value, error)
func divide(a, b int) either.Either[error, int] {
    if b == 0 {
        return either.Left[int](errors.New("division by zero"))
    }
    return either.Right[error](a / b)
}

func main() {
    // Compose operations with automatic error handling
    finalResult := function.Pipe2(
        divide(10, 2),
        either.Map(func(x int) int { return x * 2 }),
        either.GetOrElse(func() int { return 0 }),
    )

    fmt.Println("Result:", finalResult) // Output: Result: 10

    // Error case is handled automatically
    errorResult := function.Pipe2(
        divide(10, 0), // Returns Left
        either.Map(func(x int) int { return x * 2 }), // Skipped!
        either.GetOrElse(func() int { return 0 }), // Returns default
    )

    fmt.Println("Error result:", errorResult) // Output: Error result: 0
}`}
</CodeCard>

<Callout type="warn" title="Note.">
  v1 uses <code>Either[error, A]</code> instead of v2's <code>Result[A]</code>.
</Callout>

  </TabItem>
</Tabs>

### Run it

<CodeCard file="shell">
{`go run main.go`}
</CodeCard>

<CodeCard file="output">
{`Result: 10
Error result: 0`}
</CodeCard>

</Section>

<Section
  id="pattern"
  number="04"
  title="Understanding the"
  titleAccent="pattern."
  lede="Four building blocks: Result/Either return types, Pipe to compose, Map to transform, GetOrElse to extract."
>

### 1. Pure functions return results

<Tabs>
  <TabItem value="v2" label="v2" default>

<CodeCard file="signature.go">
{`// Instead of this:
func divide(a, b int) (int, error)

// We write this:
func divide(a, b int) result.Result[int]`}
</CodeCard>

  </TabItem>
  <TabItem value="v1" label="v1">

<CodeCard file="signature.go">
{`// Instead of this:
func divide(a, b int) (int, error)

// We write this:
func divide(a, b int) either.Either[error, int]`}
</CodeCard>

  </TabItem>
</Tabs>

### 2. Compose with Pipe

<CodeCard file="pipe.go">
{`result := function.Pipe2(
    operation1(),      // Returns Result[A]
    operation2,        // A -> Result[B]
    operation3,        // B -> C
)`}
</CodeCard>

<p><strong>Pipe</strong> feeds the output of one function into the next, automatically handling errors.</p>

### 3. Transform with Map

<Tabs>
  <TabItem value="v2" label="v2" default>

<CodeCard file="map.go">
{`result.Map(func(x int) int {
    return x * 2
})`}
</CodeCard>

  </TabItem>
  <TabItem value="v1" label="v1">

<CodeCard file="map.go">
{`either.Map(func(x int) int {
    return x * 2
})`}
</CodeCard>

  </TabItem>
</Tabs>

<p><strong>Map</strong> transforms the value inside a Result/Either — but only if it's successful. Errors pass through unchanged.</p>

### 4. Extract with GetOrElse

<Tabs>
  <TabItem value="v2" label="v2" default>

<CodeCard file="getorelse.go">
{`result.GetOrElse(func() int {
    return 0
})`}
</CodeCard>

  </TabItem>
  <TabItem value="v1" label="v1">

<CodeCard file="getorelse.go">
{`either.GetOrElse(func() int {
    return 0
})`}
</CodeCard>

  </TabItem>
</Tabs>

<p><strong>GetOrElse</strong> extracts the value or provides a default if there was an error.</p>

</Section>

<Section
  id="complex"
  number="05"
  title="A more complex"
  titleAccent="example."
  lede="Chain multiple operations into one pipeline. Failure in any step short-circuits the rest."
>

<Tabs>
  <TabItem value="v2" label="v2" default>

<CodeCard file="calculator.go" status="tested">
{`package main

import (
    "errors"
    "fmt"
    "github.com/IBM/fp-go/v2/function"
    "github.com/IBM/fp-go/v2/result"
)

func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}

func sqrt(n int) result.Result[int] {
    if n < 0 {
        return result.Err[int](errors.New("cannot sqrt negative number"))
    }
    // Simplified integer sqrt
    result := 0
    for result*result <= n {
        result++
    }
    return result.Ok(result - 1)
}

func calculate(a, b, c int) result.Result[string] {
    return function.Pipe3(
        divide(a, b),                                    // 100 / 4 = 25
        result.Chain(func(x int) result.Result[int] {   // Chain another Result operation
            return divide(x, c)                          // 25 / 5 = 5
        }),
        result.Chain(sqrt),                              // sqrt(5) ≈ 2
        result.Map(func(x int) string {                  // Convert to string
            return fmt.Sprintf("Final result: %d", x)
        }),
    )
}

func main() {
    // Success case
    success := calculate(100, 4, 5)
    fmt.Println(result.GetOrElse(func() string {
        return "Error occurred"
    })(success))
    // Output: Final result: 2

    // Error case (division by zero)
    failure := calculate(100, 0, 5)
    fmt.Println(result.GetOrElse(func() string {
        return "Error occurred"
    })(failure))
    // Output: Error occurred
}`}
</CodeCard>

  </TabItem>
  <TabItem value="v1" label="v1">

<CodeCard file="calculator.go">
{`package main

import (
    "errors"
    "fmt"
    "github.com/IBM/fp-go/either"
    "github.com/IBM/fp-go/function"
)

func divide(a, b int) either.Either[error, int] {
    if b == 0 {
        return either.Left[int](errors.New("division by zero"))
    }
    return either.Right[error](a / b)
}

func sqrt(n int) either.Either[error, int] {
    if n < 0 {
        return either.Left[int](errors.New("cannot sqrt negative number"))
    }
    // Simplified integer sqrt
    result := 0
    for result*result <= n {
        result++
    }
    return either.Right[error](result - 1)
}

func calculate(a, b, c int) either.Either[error, string] {
    return function.Pipe3(
        divide(a, b),                                           // 100 / 4 = 25
        either.Chain(func(x int) either.Either[error, int] {   // Chain another Either operation
            return divide(x, c)                                 // 25 / 5 = 5
        }),
        either.Chain(sqrt),                                     // sqrt(5) ≈ 2
        either.Map(func(x int) string {                         // Convert to string
            return fmt.Sprintf("Final result: %d", x)
        }),
    )
}

func main() {
    // Success case
    success := calculate(100, 4, 5)
    fmt.Println(either.GetOrElse(func() string {
        return "Error occurred"
    })(success))
    // Output: Final result: 2

    // Error case (division by zero)
    failure := calculate(100, 0, 5)
    fmt.Println(either.GetOrElse(func() string {
        return "Error occurred"
    })(failure))
    // Output: Error occurred
}`}
</CodeCard>

  </TabItem>
</Tabs>

<ApiTable
  columns={['Concept', 'When to use', 'Notes']}
  rows={[
    {symbol: 'Chain', signature: 'a.k.a. FlatMap', description: 'Use when your transformation returns another Result/Either.'},
    {symbol: 'Map', signature: 'Functor map', description: 'Use when your transformation returns a plain value.'},
    {symbol: 'Pipe', signature: 'function.PipeN', description: 'Compose multiple operations into a pipeline.'},
    {symbol: 'Auto error', signature: 'Short-circuit', description: 'If any step fails, subsequent steps are skipped.'},
  ]}
/>

</Section>

<Section
  id="next"
  number="06"
  title="What's"
  titleAccent="next?"
  lede="Pick a thread to keep learning."
>

<ApiTable
  columns={['Topic', 'Page', 'Why']}
  rows={[
    {symbol: 'Why', signature: <a href="./why-fp-go">Why fp-go?</a>, description: 'Understand the benefits and when to reach for it.'},
    {symbol: 'Pure functions', signature: <a href="./concepts/pure-functions">Pure functions</a>, description: 'The foundation of functional programming.'},
    {symbol: 'Monads', signature: <a href="./concepts/monads">Monads</a>, description: 'The pattern behind Result, Either, IO, etc.'},
    {symbol: 'Result', signature: <a href="./v2/result">Result</a>, description: 'Recommended type for v2 error handling.'},
    {symbol: 'Either', signature: <a href="./v2/either">Either</a>, description: 'Generic sum type.'},
    {symbol: 'Option', signature: <a href="./v2/option">Option</a>, description: 'Handle optional values safely.'},
    {symbol: 'IO', signature: <a href="./v2/io">IO</a>, description: 'Manage side effects.'},
    {symbol: 'Recipes · errors', signature: <a href="./recipes/error-handling">Error handling</a>, description: 'Production-style patterns.'},
    {symbol: 'Recipes · HTTP', signature: <a href="./recipes/http-requests">HTTP requests</a>, description: 'Effectful pipelines.'},
  ]}
/>

### Common questions

<Callout title="When should I use fp-go?">
  <p><strong>Good fit:</strong> complex error handling, data transformation pipelines, composable business logic, testing-heavy codebases.</p>
  <p><strong>Not ideal:</strong> simple CRUD, performance-critical hot paths (use idiomatic packages), teams unfamiliar with FP.</p>
</Callout>

<Callout type="success" title="Is fp-go production-ready?">
  Yes. fp-go is used in production at IBM and elsewhere. v2 is actively maintained and recommended for new projects.
</Callout>

<Callout type="info" title="What's the performance impact?">
  Standard packages: minimal overhead. Idiomatic packages (v2): 2–32× faster, near-native performance. See the <a href="./advanced/performance">performance guide</a> for details.
</Callout>

<Callout type="info" title="Migrating from v1 to v2?">
  See the <a href="./migration/v1-to-v2">migration guide</a> for a complete walkthrough of the 5 breaking changes and how to handle them.
</Callout>

</Section>

<Section
  id="summary"
  number="07"
  title="Summary"
>

<Checklist
  title="What you learned"
  items={[
    {label: 'How to install fp-go (v1 or v2)', done: true},
    {label: 'How to write pure functions with Result/Either', done: true},
    {label: 'How to compose operations with Pipe', done: true},
    {label: 'How to transform values with Map and Chain', done: true},
    {label: 'How to handle errors automatically', done: true},
    {label: 'How to build complex pipelines', done: true},
  ]}
/>

<Callout type="success" title="Next step.">
  Read <a href="./why-fp-go">Why fp-go?</a> to understand when and why to use functional programming in Go.
</Callout>

### Need help?

<ul>
  <li><a href="./intro">Full documentation</a></li>
  <li><a href="https://github.com/IBM/fp-go/discussions">GitHub Discussions</a></li>
  <li><a href="https://github.com/IBM/fp-go/issues">Report issues</a></li>
  <li><a href="https://pkg.go.dev/github.com/IBM/fp-go/v2">API reference</a></li>
</ul>

</Section>
