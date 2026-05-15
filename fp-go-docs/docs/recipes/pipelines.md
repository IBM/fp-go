---
title: Function Pipelines
hide_title: true
description: Build composable data processing pipelines using fp-go's Pipe, Flow, and composition utilities.
sidebar_position: 14
---

<PageHeader
  eyebrow="Recipes · 14 / 17"
  title="Function"
  titleAccent="Pipelines"
  lede="Build composable data processing pipelines using fp-go's Pipe, Flow, and composition utilities for readable, left-to-right data flow."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '8' },
    { label: 'Use Cases', value: 'Data transformation, ETL, request processing' }
  ]}
/>

<TLDR>
  <TLDRCard title="Left-to-Right Flow" icon="arrow-right">
    Pipe and Flow enable readable, sequential operations that mirror natural data flow direction.
  </TLDRCard>
  <TLDRCard title="Reusable Pipelines" icon="recycle">
    Flow creates reusable transformation functions that can be applied to different inputs.
  </TLDRCard>
  <TLDRCard title="Type-Safe Composition" icon="shield">
    Compile-time guarantees ensure each step's output matches the next step's input type.
  </TLDRCard>
</TLDR>

<Section id="basic-pipelines" number="01" title="Basic" titleAccent="Pipelines">

Function pipelines chain operations in a readable, left-to-right flow using **Pipe** and **Flow**.

<Compare>
<CompareCol kind="bad">
<CodeCard file="nested.go">
{`// ❌ Nested function calls (right-to-left)
result := replaceSpaces(
    toUpper(
        trim("  hello world  ")
    )
)
// Hard to read, inside-out
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="pipeline.go">
{`// ✅ Pipeline (left-to-right)
result := F.Pipe3(
    "  hello world  ",
    strings.TrimSpace,
    strings.ToUpper,
    func(s string) string {
        return strings.ReplaceAll(s, " ", "_")
    },
)
// HELLO_WORLD
// Easy to read, natural flow
`}
</CodeCard>
</CompareCol>
</Compare>

<CodeCard file="flow_reusable.go">
{`package main

import (
    "fmt"
    "strings"
    F "github.com/IBM/fp-go/v2/function"
)

// Create reusable pipeline with Flow
var normalizeString = F.Flow3(
    strings.TrimSpace,
    strings.ToLower,
    func(s string) string {
        return strings.ReplaceAll(s, " ", "-")
    },
)

func main() {
    result1 := normalizeString("  Hello World  ")
    result2 := normalizeString("  Functional Programming  ")
    
    fmt.Println(result1) // hello-world
    fmt.Println(result2) // functional-programming
}
`}
</CodeCard>

<Callout type="info">
**Pipe vs Flow**: Use `Pipe` to process a value immediately, `Flow` to create a reusable transformation function.
</Callout>

</Section>

<Section id="data-transformation" number="02" title="Data Transformation" titleAccent="Pipelines">

Process collections with array operations in a pipeline.

<CodeCard file="array_pipeline.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

type Product struct {
    Name  string
    Price float64
}

func main() {
    products := []Product{
        {Name: "Laptop", Price: 999.99},
        {Name: "Mouse", Price: 29.99},
        {Name: "Keyboard", Price: 79.99},
        {Name: "Monitor", Price: 299.99},
    }
    
    // Pipeline: filter expensive items, extract prices, format
    result := F.Pipe3(
        products,
        A.Filter(func(p Product) bool {
            return p.Price > 50
        }),
        A.Map(func(p Product) string {
            return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
        }),
        A.Reduce(func(acc, item string) string {
            if acc == "" {
                return item
            }
            return acc + "\\n" + item
        })(""),
    )
    
    fmt.Println(result)
    // Laptop: $999.99
    // Keyboard: $79.99
    // Monitor: $299.99
}
`}
</CodeCard>

</Section>

<Section id="error-handling" number="03" title="Error Handling" titleAccent="Pipelines">

Build pipelines with error handling using Either and Option.

<CodeCard file="either_pipeline.go">
{`package main

import (
    "fmt"
    "strconv"
    E "github.com/IBM/fp-go/v2/either"
    F "github.com/IBM/fp-go/v2/function"
)

func parseInt(s string) E.Either[error, int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return E.Left[int](err)
    }
    return E.Right[error](n)
}

func validatePositive(n int) E.Either[error, int] {
    if n <= 0 {
        return E.Left[int](fmt.Errorf("must be positive, got %d", n))
    }
    return E.Right[error](n)
}

func double(n int) int {
    return n * 2
}

func processNumber(s string) E.Either[error, int] {
    return F.Pipe2(
        parseInt(s),
        E.Chain(func(n int) E.Either[error, int] {
            return validatePositive(n)
        }),
        E.Map(double),
    )
}

func main() {
    // Valid input
    result1 := processNumber("42")
    E.Match(
        func(err error) { fmt.Println("Error:", err) },
        func(n int) { fmt.Println("Result:", n) },
    )(result1)
    // Result: 84
    
    // Invalid input
    result2 := processNumber("-5")
    E.Match(
        func(err error) { fmt.Println("Error:", err) },
        func(n int) { fmt.Println("Result:", n) },
    )(result2)
    // Error: must be positive, got -5
}
`}
</CodeCard>

<CodeCard file="option_pipeline.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

type User struct {
    ID    int
    Name  string
    Email string
}

func findUserByID(users []User, id int) O.Option[User] {
    return A.FindFirst(func(u User) bool {
        return u.ID == id
    })(users)
}

func getUserEmail(users []User, id int) O.Option[string] {
    return F.Pipe2(
        findUserByID(users, id),
        O.Map(func(u User) string {
            return u.Email
        }),
    )
}

func main() {
    users := []User{
        {ID: 1, Name: "Alice", Email: "alice@example.com"},
        {ID: 2, Name: "Bob", Email: "bob@example.com"},
    }
    
    // Found
    email1 := getUserEmail(users, 1)
    O.Match(
        func() { fmt.Println("User not found") },
        func(email string) { fmt.Println("Email:", email) },
    )(email1)
    // Email: alice@example.com
    
    // Not found
    email2 := getUserEmail(users, 999)
    O.Match(
        func() { fmt.Println("User not found") },
        func(email string) { fmt.Println("Email:", email) },
    )(email2)
    // User not found
}
`}
</CodeCard>

</Section>

<Section id="io-pipelines" number="04" title="IO" titleAccent="Pipelines">

Chain IO operations for file processing and HTTP requests.

<CodeCard file="file_pipeline.go">
{`package main

import (
    "fmt"
    "os"
    "strings"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

func readFile(path string) IOE.IOEither[error, string] {
    return IOE.TryCatch(func() (string, error) {
        data, err := os.ReadFile(path)
        return string(data), err
    })
}

func writeFile(path string, content string) IOE.IOEither[error, int] {
    return IOE.TryCatch(func() (int, error) {
        err := os.WriteFile(path, []byte(content), 0644)
        return len(content), err
    })
}

func processLines(lines []string) []string {
    return F.Pipe2(
        lines,
        A.Filter(func(line string) bool {
            return strings.TrimSpace(line) != ""
        }),
        A.Map(strings.ToUpper),
    )
}

func transformFile(input, output string) IOE.IOEither[error, int] {
    return F.Pipe3(
        readFile(input),
        IOE.Map(func(content string) []string {
            return strings.Split(content, "\\n")
        }),
        IOE.Map(processLines),
        IOE.Chain(func(lines []string) IOE.IOEither[error, int] {
            content := strings.Join(lines, "\\n")
            return writeFile(output, content)
        }),
    )
}

func main() {
    result := transformFile("input.txt", "output.txt")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Wrote %d bytes\\n", result.Right())
    }
}
`}
</CodeCard>

</Section>

<Section id="complex-pipelines" number="05" title="Multi-Stage" titleAccent="Pipelines">

Build complex pipelines with multiple transformation stages.

<CodeCard file="multistage_pipeline.go">
{`package main

import (
    "fmt"
    "strconv"
    "strings"
    A "github.com/IBM/fp-go/v2/array"
    E "github.com/IBM/fp-go/v2/either"
    F "github.com/IBM/fp-go/v2/function"
)

type RawData struct {
    ID    string
    Value string
}

type ParsedData struct {
    ID    int
    Value float64
}

type ValidatedData struct {
    ID    int
    Value float64
}

type ProcessedData struct {
    ID     int
    Result string
}

// Stage 1: Parse
func parseData(raw RawData) E.Either[error, ParsedData] {
    id, err := strconv.Atoi(raw.ID)
    if err != nil {
        return E.Left[ParsedData](fmt.Errorf("invalid ID: %w", err))
    }
    
    value, err := strconv.ParseFloat(raw.Value, 64)
    if err != nil {
        return E.Left[ParsedData](fmt.Errorf("invalid value: %w", err))
    }
    
    return E.Right[error](ParsedData{ID: id, Value: value})
}

// Stage 2: Validate
func validateData(parsed ParsedData) E.Either[error, ValidatedData] {
    if parsed.Value < 0 {
        return E.Left[ValidatedData](fmt.Errorf("value must be non-negative"))
    }
    return E.Right[error](ValidatedData{ID: parsed.ID, Value: parsed.Value})
}

// Stage 3: Process
func processData(validated ValidatedData) ProcessedData {
    return ProcessedData{
        ID:     validated.ID,
        Result: fmt.Sprintf("Processed: %.2f", validated.Value*2),
    }
}

// Complete pipeline
func processPipeline(raw []RawData) E.Either[error, []ProcessedData] {
    return F.Pipe3(
        raw,
        A.Traverse[RawData](E.Applicative[error, ParsedData]())(parseData),
        E.Chain(func(parsed []ParsedData) E.Either[error, []ValidatedData] {
            return A.Traverse[ParsedData](E.Applicative[error, ValidatedData]())(
                validateData,
            )(parsed)
        }),
        E.Map(func(validated []ValidatedData) []ProcessedData {
            return A.Map(processData)(validated)
        }),
    )
}

func main() {
    rawData := []RawData{
        {ID: "1", Value: "10.5"},
        {ID: "2", Value: "20.3"},
        {ID: "3", Value: "15.7"},
    }
    
    result := processPipeline(rawData)
    
    E.Match(
        func(err error) { fmt.Println("Error:", err) },
        func(processed []ProcessedData) {
            fmt.Println("Processed data:")
            for _, p := range processed {
                fmt.Printf("  ID %d: %s\\n", p.ID, p.Result)
            }
        },
    )(result)
}
`}
</CodeCard>

</Section>

<Section id="branching" number="06" title="Branching" titleAccent="Pipelines">

Handle different processing paths based on input type or conditions.

<CodeCard file="branching_pipeline.go">
{`package main

import (
    "fmt"
    E "github.com/IBM/fp-go/v2/either"
)

type Input struct {
    Type  string
    Value int
}

type Output struct {
    Result string
}

func processTypeA(value int) E.Either[error, Output] {
    return E.Right[error](Output{
        Result: fmt.Sprintf("Type A: %d", value*2),
    })
}

func processTypeB(value int) E.Either[error, Output] {
    return E.Right[error](Output{
        Result: fmt.Sprintf("Type B: %d", value*3),
    })
}

func processInput(input Input) E.Either[error, Output] {
    switch input.Type {
    case "A":
        return processTypeA(input.Value)
    case "B":
        return processTypeB(input.Value)
    default:
        return E.Left[Output](fmt.Errorf("unknown type: %s", input.Type))
    }
}

func main() {
    inputs := []Input{
        {Type: "A", Value: 10},
        {Type: "B", Value: 10},
        {Type: "C", Value: 10},
    }
    
    for _, input := range inputs {
        result := processInput(input)
        E.Match(
            func(err error) { fmt.Printf("Error processing %+v: %v\\n", input, err) },
            func(output Output) { fmt.Println(output.Result) },
        )(result)
    }
    // Type A: 20
    // Type B: 30
    // Error processing {Type:C Value:10}: unknown type: C
}
`}
</CodeCard>

</Section>

<Section id="lazy-evaluation" number="07" title="Lazy Evaluation" titleAccent="Pipelines">

Defer computation until needed with lazy pipelines.

<CodeCard file="lazy_pipeline.go">
{`package main

import (
    "fmt"
    L "github.com/IBM/fp-go/v2/lazy"
    F "github.com/IBM/fp-go/v2/function"
)

func expensiveOperation(n int) L.Lazy[int] {
    return L.MakeLazy(func() int {
        fmt.Printf("Computing %d...\\n", n)
        return n * 2
    })
}

func main() {
    // Pipeline is not executed until needed
    pipeline := F.Pipe2(
        expensiveOperation(10),
        L.Map(func(n int) int {
            fmt.Println("Doubling again...")
            return n * 2
        }),
    )
    
    fmt.Println("Pipeline created, not executed yet")
    
    // Execute pipeline
    result := pipeline()
    fmt.Println("Result:", result)
    // Computing 10...
    // Doubling again...
    // Result: 40
}
`}
</CodeCard>

</Section>

<Section id="best-practices" number="08" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Keep pipelines focused** — Each pipeline should have a single, clear responsibility
  </ChecklistItem>
  <ChecklistItem status="required">
    **Name intermediate steps** — Use variables for clarity in complex pipelines
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Handle errors early** — Validate inputs at the start of the pipeline
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Use Flow for reusability** — Create reusable transformation functions
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Consider lazy evaluation** — Defer expensive computations when appropriate
  </ChecklistItem>
</Checklist>

<Compare>
<CompareCol kind="good">
<CodeCard file="good_pipeline.go">
{`// ✅ Good: Single responsibility
var normalizeEmail = F.Flow2(
    strings.TrimSpace,
    strings.ToLower,
)

var validateEmail = func(email string) E.Either[error, string] {
    if !strings.Contains(email, "@") {
        return E.Left[string](fmt.Errorf("invalid email"))
    }
    return E.Right[error](email)
}

// ✅ Good: Clear intermediate steps
func processOrder(order Order) E.Either[error, Receipt] {
    validated := validateOrder(order)
    priced := E.Map(calculatePrice)(validated)
    charged := E.Chain(chargeCustomer)(priced)
    return E.Map(generateReceipt)(charged)
}

// ✅ Good: Validate early
func processData(input string) E.Either[error, Result] {
    return F.Pipe2(
        validateInput(input),
        E.Chain(parseInput),
        E.Chain(transformData),
    )
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="bad">
<CodeCard file="bad_pipeline.go">
{`// ❌ Avoid: Doing too much in one pipeline
var processEmail = func(email string) E.Either[error, string] {
    // Normalize, validate, send, log, update DB...
}

// ❌ Avoid: Unclear nested pipes
func processOrder(order Order) E.Either[error, Receipt] {
    return F.Pipe3(
        validateOrder(order),
        E.Chain(func(o Order) E.Either[error, Order] {
            return E.Map(func(p float64) Order {
                o.Total = p
                return o
            })(calculatePrice(o))
        }),
        // ... more nesting
    )
}

// ❌ Avoid: Late validation
func processData(input string) E.Either[error, Result] {
    return F.Pipe2(
        parseInput(input), // Might fail on invalid input
        E.Chain(validateParsed),
        E.Chain(transformData),
    )
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
