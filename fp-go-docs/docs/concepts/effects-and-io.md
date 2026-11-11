---
sidebar_position: 14
title: Effects and IO
hide_title: true
description: Learn to manage side effects explicitly and safely with fp-go's IO types.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Concepts · 04 / 06"
  title="Effects and"
  titleAccent="IO."
  lede="Manage side effects explicitly and safely. Learn how IO types let you describe effects as pure values, then execute them when ready."
  meta={[
    {label: '// Difficulty', value: 'Intermediate'},
    {label: '// Reading time', value: '15 min · 9 sections'},
    {label: '// Prereqs', value: 'Pure functions, composition'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Key insight" prose value={<>Effects are <em>descriptions</em>, not executions.</>} variant="up" />
  <TLDRCard label="// Recommended type" prose value={<><code>IOResult[A]</code> for most effects.</>} />
  <TLDRCard label="// Benefits" prose value={<>Lazy, composable, <em>testable</em> effects.</>} />
</TLDR>

<Section id="what-are-effects" number="01" title="What are" titleAccent="effects?">

An **effect** is anything that interacts with the outside world:

<CodeCard file="effects.go">
{`// ❌ Reading files
data, err := os.ReadFile("config.json")

// ❌ Writing files
os.WriteFile("output.txt", data, 0644)

// ❌ Network calls
resp, err := http.Get("https://api.example.com")

// ❌ Database queries
rows, err := db.Query("SELECT * FROM users")

// ❌ Printing to console
fmt.Println("Hello, World!")

// ❌ Getting current time
now := time.Now()

// ❌ Random numbers
rand.Int()

// ❌ Modifying global state
globalCounter++`}
</CodeCard>

**Why are these effects?**
- They depend on external state
- They modify external state
- They're not deterministic
- They can fail

</Section>

<Section id="problem" number="02" title="The problem with" titleAccent="effects.">

### Effects Make Functions Impure

<Compare>
<CompareCol kind="bad">

<CodeCard file="impure.go">
{`// Impure: executes immediately, can't control when
func loadConfig() Config {
    data, _ := os.ReadFile("config.json")  // Effect happens NOW
    var config Config
    json.Unmarshal(data, &config)
    return config
}

// Problems:
// 1. Executes immediately (can't delay)
// 2. Hard to test (needs real file)
// 3. Can't compose easily
// 4. Errors ignored`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="pure.go">
{`// Pure: returns a DESCRIPTION of the effect
func loadConfig() ioresult.IOResult[Config] {
    return func() result.Result[Config] {
        data, err := os.ReadFile("config.json")
        if err != nil {
            return result.Err[Config](err)
        }
        
        var config Config
        err = json.Unmarshal(data, &config)
        return result.FromGoError(config, err)
    }
}

// Benefits:
// 1. Doesn't execute until called
// 2. Easy to test (mock the function)
// 3. Composable
// 4. Proper error handling`}
</CodeCard>

</CompareCol>
</Compare>

</Section>

<Section id="io-type" number="03" title="The IO" titleAccent="type.">

### What is IO?

**IO** is a type that represents a lazy computation that performs side effects.

<CodeCard file="io-basic.go">
{`// IO[A] is just a function that returns A
type IO[A any] func() A

// Example
var readFile IO[[]byte] = func() []byte {
    data, _ := os.ReadFile("file.txt")
    return data
}

// Nothing happens yet!
// The function is just stored

// Execute it
data := readFile()  // NOW the file is read`}
</CodeCard>

### Key Properties

1. **Lazy** - Doesn't execute until you call it
2. **Composable** - Can be chained with other IOs
3. **Testable** - Can be mocked or replaced
4. **Explicit** - Effect is visible in the type

</Section>

<Section id="io-variants" number="04" title="IO variants in" titleAccent="fp-go.">

### IO - Simple Effects

<CodeCard file="io-simple.go">
{`import "github.com/IBM/fp-go/v2/io"

// IO[A] - returns A
type IO[A any] func() A

// Example
printHello := io.Of(func() string {
    fmt.Println("Hello!")
    return "done"
})

// Execute
result := printHello()  // Prints "Hello!", returns "done"`}
</CodeCard>

### IOOption - Effects that Might Fail

<CodeCard file="io-option.go">
{`import "github.com/IBM/fp-go/v2/iooption"

// IOOption[A] - returns Option[A]
type IOOption[A any] func() option.Option[A]

// Example
findUser := func(id string) iooption.IOOption[User] {
    return func() option.Option[User] {
        user := db.FindByID(id)
        if user == nil {
            return option.None[User]()
        }
        return option.Some(*user)
    }
}

// Execute
userOpt := findUser("123")()`}
</CodeCard>

### IOResult - Effects with Error Handling ⭐

<CodeCard file="io-result.go">
{`import "github.com/IBM/fp-go/v2/ioresult"

// IOResult[A] - returns Result[A]
type IOResult[A any] func() result.Result[A]

// Example
readConfig := func(path string) ioresult.IOResult[Config] {
    return func() result.Result[Config] {
        data, err := os.ReadFile(path)
        if err != nil {
            return result.Err[Config](err)
        }
        
        var config Config
        err = json.Unmarshal(data, &config)
        return result.FromGoError(config, err)
    }
}

// Execute
configResult := readConfig("config.json")()`}
</CodeCard>

<Callout type="success">
**Recommendation:** Use IOResult for most effects with error handling.
</Callout>

</Section>

<Section id="working-with-io" number="05" title="Working with" titleAccent="IO.">

### Creating IO

<Tabs groupId="io">
<TabItem value="of" label="Of - Wrap Value">

<CodeCard file="io-of.go">
{`// Wrap a pure value in IO
io := io.Of(func() int {
    return 42
})

result := io()  // 42`}
</CodeCard>

</TabItem>
<TabItem value="effect" label="Effect - Side Effect">

<CodeCard file="io-effect.go">
{`// Create IO from side effect
printIO := io.Of(func() string {
    fmt.Println("Hello!")
    return "printed"
})

printIO()  // Prints and returns "printed"`}
</CodeCard>

</TabItem>
<TabItem value="ioresult" label="IOResult - With Errors">

<CodeCard file="io-result-create.go">
{`// Create IOResult
readFile := func(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}

result := readFile("file.txt")()`}
</CodeCard>

</TabItem>
</Tabs>

### Transforming IO

#### Map - Transform the Result

<CodeCard file="io-map.go">
{`// Map: transform the value inside
readNumber := io.Of(func() int { return 42 })

doubled := io.Map(func(x int) int {
    return x * 2
})(readNumber)

result := doubled()  // 84`}
</CodeCard>

#### Chain - Sequence Effects

<CodeCard file="io-chain.go">
{`// Chain: sequence two IOs
readFile := func(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}

parseJSON := func(data []byte) ioresult.IOResult[Config] {
    return func() result.Result[Config] {
        var config Config
        err := json.Unmarshal(data, &config)
        return result.FromGoError(config, err)
    }
}

// Chain them
loadConfig := ioresult.Chain(parseJSON)(readFile("config.json"))

// Execute
config := loadConfig()`}
</CodeCard>

### Composing IO

<CodeCard file="io-compose.go">
{`import "github.com/IBM/fp-go/v2/function"

// Build pipeline
loadAndValidate := function.Pipe3(
    readFile("config.json"),
    ioresult.Chain(parseJSON),
    ioresult.Chain(validateConfig),
)

// Execute when ready
result := loadAndValidate()`}
</CodeCard>

</Section>

<Section id="examples" number="06" title="Real-world" titleAccent="examples.">

### Example 1: File Operations

<Compare>
<CompareCol kind="bad">

<CodeCard file="file-without-io.go">
{`func processFile(input, output string) error {
    // Executes immediately
    data, err := os.ReadFile(input)
    if err != nil {
        return err
    }
    
    // Transform
    processed := transform(data)
    
    // Write
    return os.WriteFile(output, processed, 0644)
}

// Hard to test
// Can't delay execution
// Can't compose`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="file-with-io.go">
{`import (
    "github.com/IBM/fp-go/v2/ioresult"
    "github.com/IBM/fp-go/v2/function"
)

func readFile(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}

func writeFile(path string) func([]byte) ioresult.IOResult[[]byte] {
    return func(data []byte) ioresult.IOResult[[]byte] {
        return func() result.Result[[]byte] {
            err := os.WriteFile(path, data, 0644)
            return result.FromGoError(data, err)
        }
    }
}

func transform(data []byte) []byte {
    // Pure transformation
    return processed
}

func processFile(input, output string) ioresult.IOResult[[]byte] {
    return function.Pipe3(
        readFile(input),
        ioresult.Map(transform),
        ioresult.Chain(writeFile(output)),
    )
}

// Easy to test (mock readFile/writeFile)
// Lazy execution
// Composable`}
</CodeCard>

</CompareCol>
</Compare>

### Example 2: HTTP API Call

<Compare>
<CompareCol kind="bad">

<CodeCard file="http-without-io.go">
{`func fetchUser(id string) (User, error) {
    // Executes immediately
    resp, err := http.Get("https://api.example.com/users/" + id)
    if err != nil {
        return User{}, err
    }
    defer resp.Body.Close()
    
    var user User
    err = json.NewDecoder(resp.Body).Decode(&user)
    return user, err
}

// Executes on call
// Hard to test
// Can't retry easily`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="http-with-io.go">
{`func fetchUser(id string) ioresult.IOResult[User] {
    return func() result.Result[User] {
        resp, err := http.Get("https://api.example.com/users/" + id)
        if err != nil {
            return result.Err[User](err)
        }
        defer resp.Body.Close()
        
        var user User
        err = json.NewDecoder(resp.Body).Decode(&user)
        return result.FromGoError(user, err)
    }
}

// Lazy - doesn't execute until called
// Easy to test (mock the function)
// Can retry with ioresult.Retry`}
</CodeCard>

</CompareCol>
</Compare>

</Section>

<Section id="lazy-evaluation" number="07" title="Lazy" titleAccent="evaluation.">

### What is Lazy Evaluation?

**Lazy** means the computation doesn't happen until you explicitly execute it.

<CodeCard file="lazy.go">
{`// Create IO (doesn't execute)
io := readFile("large-file.txt")

// Still hasn't executed
time.Sleep(1 * time.Second)

// NOW it executes
data := io()`}
</CodeCard>

### Benefits

#### 1. Control When Effects Happen

<CodeCard file="control-effects.go">
{`// Build pipeline
pipeline := function.Pipe3(
    readFile("input.txt"),
    ioresult.Map(process),
    ioresult.Chain(writeFile("output.txt")),
)

// Nothing has happened yet!

// Execute when ready
if shouldProcess {
    result := pipeline()
}`}
</CodeCard>

#### 2. Compose Before Executing

<CodeCard file="compose-before-execute.go">
{`// Build complex pipeline
step1 := readFile("file1.txt")
step2 := ioresult.Chain(parseJSON)(step1)
step3 := ioresult.Chain(validate)(step2)
step4 := ioresult.Chain(transform)(step3)
step5 := ioresult.Chain(writeFile("output.txt"))(step4)

// Execute entire pipeline
result := step5()`}
</CodeCard>

#### 3. Retry Logic

<CodeCard file="retry.go">
{`import "github.com/IBM/fp-go/v2/ioresult"

// Create effect
fetchData := func() ioresult.IOResult[Data] {
    return func() result.Result[Data] {
        // Network call
        return result.Ok(data)
    }
}

// Retry on failure
withRetry := ioresult.Retry(
    3,                    // Max attempts
    100*time.Millisecond, // Delay
)(fetchData())

// Execute with retry
result := withRetry()`}
</CodeCard>

#### 4. Testing

<CodeCard file="testing-io.go">
{`// Production
var readFile = func(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}

// Test
func TestProcessFile(t *testing.T) {
    // Mock readFile
    readFile = func(path string) ioresult.IOResult[[]byte] {
        return func() result.Result[[]byte] {
            return result.Ok([]byte("test data"))
        }
    }
    
    // Test
    result := processFile("input.txt", "output.txt")()
    assert.True(t, result.IsOk())
}`}
</CodeCard>

</Section>

<Section id="separation" number="08" title="Separating description from" titleAccent="execution.">

### The Key Insight

<CodeCard file="separation.go">
{`// Description (pure)
var loadConfig = func() ioresult.IOResult[Config] {
    return function.Pipe3(
        readFile("config.json"),
        ioresult.Chain(parseJSON),
        ioresult.Chain(validate),
    )
}

// Execution (impure)
func main() {
    // Build description
    io := loadConfig()
    
    // Execute
    result := io()
    
    // Handle result
    result.Fold(
        func(err error) {
            log.Fatal(err)
        },
        func(config Config) {
            // Use config
        },
    )
}`}
</CodeCard>

**Benefits:**
- Description is pure and testable
- Execution is isolated
- Clear boundary between pure/impure

</Section>

<Section id="patterns" number="09" title="Common" titleAccent="patterns.">

### Pattern 1: Sequential Effects

<CodeCard file="sequential.go">
{`// Execute effects in sequence
pipeline := function.Pipe4(
    effect1(),
    ioresult.Chain(effect2),
    ioresult.Chain(effect3),
    ioresult.Chain(effect4),
)

result := pipeline()`}
</CodeCard>

### Pattern 2: Parallel Effects

<CodeCard file="parallel.go">
{`// Execute effects in parallel
var wg sync.WaitGroup
results := make(chan result.Result[Data], 3)

effects := []ioresult.IOResult[Data]{
    fetchFromAPI1(),
    fetchFromAPI2(),
    fetchFromAPI3(),
}

for _, effect := range effects {
    wg.Add(1)
    go func(e ioresult.IOResult[Data]) {
        defer wg.Done()
        results <- e()
    }(effect)
}

wg.Wait()
close(results)`}
</CodeCard>

### Pattern 3: Conditional Effects

<CodeCard file="conditional.go">
{`// Execute effect conditionally
loadConfig := func(env string) ioresult.IOResult[Config] {
    if env == "production" {
        return readFile("/etc/app/config.json")
    }
    return readFile("./config.dev.json")
}`}
</CodeCard>

### Pattern 4: Effect with Fallback

<CodeCard file="fallback.go">
{`// Try effect, fallback on error
loadConfig := ioresult.OrElse(func(err error) ioresult.IOResult[Config] {
    log.Printf("Failed to load config: %v, using defaults", err)
    return ioresult.Of(result.Ok(defaultConfig))
})(readFile("config.json"))`}
</CodeCard>

</Section>
