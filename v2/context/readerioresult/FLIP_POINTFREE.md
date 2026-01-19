# Sequence Functions and Point-Free Style Programming

This document explains how the `Sequence*` functions in the `context/readerioresult` package enable point-free style programming and improve code composition.

## Table of Contents

1. [What is Point-Free Style?](#what-is-point-free-style)
2. [The Problem: Nested Function Application](#the-problem-nested-function-application)
3. [The Solution: Sequence Functions](#the-solution-sequence-functions)
4. [How Sequence Enables Point-Free Style](#how-sequence-enables-point-free-style)
5. [TraverseReader: Introducing Dependencies](#traversereader-introducing-dependencies)
6. [Practical Benefits](#practical-benefits)
7. [Examples](#examples)
8. [Comparison: With and Without Sequence](#comparison-with-and-without-sequence)

## What is Point-Free Style?

Point-free style (also called tacit programming) is a programming paradigm where function definitions don't explicitly mention their arguments. Instead, functions are composed using combinators and higher-order functions.

**Traditional style (with points):**
```go
func double(x int) int {
    return x * 2
}
```

**Point-free style (without points):**
```go
var double = N.Mul(2)
```

The key benefit is that point-free style emphasizes **what** the function does (its transformation) rather than **how** it manipulates data.

## The Problem: Nested Function Application

In functional programming with monadic types like `ReaderIOResult`, we often have nested structures where we need to apply parameters in a specific order. Consider:

```go
type ReaderIOResult[A any] = func(context.Context) func() Either[error, A]
type Reader[R, A any] = func(R) A

// A computation that produces a Reader
type Computation = ReaderIOResult[Reader[Config, int]]
// Expands to: func(context.Context) func() Either[error, func(Config) int]
```

To use this, we must apply parameters in this order:
1. First, provide `context.Context`
2. Then, execute the IO effect (call the function)
3. Then, unwrap the `Either` to get the `Reader`
4. Finally, provide the `Config`

This creates several problems:

### Problem 1: Awkward Parameter Order

```go
computation := getComputation()
ctx := t.Context()
cfg := Config{Value: 42}

// Must apply in this specific order
result := computation(ctx)()  // Get Either[error, Reader[Config, int]]
if reader, err := either.Unwrap(result); err == nil {
    value := reader(cfg)  // Finally apply Config
    // use value
}
```

The `Config` parameter, which is often known early and stable, must be provided last. This prevents partial application and reuse.

### Problem 2: Cannot Partially Apply Dependencies

```go
// Want to do this: create a reusable computation with Config baked in
// But can't because Config comes last!
withConfig := computation(cfg)  // ❌ Doesn't work - cfg comes last, not first
```

### Problem 3: Breaks Point-Free Composition

```go
// Want to compose like this:
var pipeline = F.Flow3(
    getComputation,
    applyConfig(cfg),  // ❌ Can't do this - Config comes last
    processResult,
)
```

## The Solution: Sequence Functions

The `Sequence*` functions solve this by "flipping" or "sequencing" the nested structure, changing the order in which parameters are applied.

### SequenceReader

```go
func SequenceReader[R, A any](
    ma ReaderIOResult[Reader[R, A]]
) Kleisli[R, A]
```

**Type transformation:**
```
From: func(context.Context) func() Either[error, func(R) A]
To:   func(R) func(context.Context) func() Either[error, A]
```

Now `R` (the Reader's environment) comes **first**, before `context.Context`!

### SequenceReaderIO

```go
func SequenceReaderIO[R, A any](
    ma ReaderIOResult[ReaderIO[R, A]]
) Kleisli[R, A]
```

**Type transformation:**
```
From: func(context.Context) func() Either[error, func(R) func() A]
To:   func(R) func(context.Context) func() Either[error, A]
```

### SequenceReaderResult

```go
func SequenceReaderResult[R, A any](
    ma ReaderIOResult[ReaderResult[R, A]]
) Kleisli[R, A]
```

**Type transformation:**
```
From: func(context.Context) func() Either[error, func(R) Either[error, A]]
To:   func(R) func(context.Context) func() Either[error, A]
```

## How Sequence Enables Point-Free Style

### 1. Partial Application

By moving the environment parameter first, we can partially apply it:

```go
type Config struct { Multiplier int }

computation := getComputation()  // ReaderIOResult[Reader[Config, int]]
sequenced := SequenceReader[Config, int](computation)

// Partially apply Config
cfg := Config{Multiplier: 5}
withConfig := sequenced(cfg)  // ✅ Now we have ReaderIOResult[int]

// Reuse with different contexts
result1 := withConfig(ctx1)()
result2 := withConfig(ctx2)()
```

### 2. Dependency Injection

Inject dependencies early in the pipeline:

```go
type Database struct { ConnectionString string }

makeQuery := func(ctx context.Context) func() Either[error, func(Database) string] {
    // ... implementation
}

// Sequence to enable DI
queryWithDB := SequenceReader[Database, string](makeQuery)

// Inject database
db := Database{ConnectionString: "localhost:5432"}
query := queryWithDB(db)  // ✅ Database injected

// Use query with any context
result := query(t.Context())()
```

### 3. Point-Free Composition

Build pipelines without mentioning intermediate values:

```go
var pipeline = F.Flow3(
    getComputation,                    // ReaderIOResult[Reader[Config, int]]
    SequenceReader[Config, int],       // func(Config) ReaderIOResult[int]
    applyConfig(cfg),                  // ReaderIOResult[int]
)

// Or with partial application:
var withConfig = F.Pipe1(
    getComputation(),
    SequenceReader[Config, int],
)

result := withConfig(cfg)(ctx)()
```

### 4. Reusable Computations

Create specialized versions of generic computations:

```go
// Generic computation
makeServiceInfo := func(ctx context.Context) func() Either[error, func(ServiceConfig) string] {
    // ... implementation
}

sequenced := SequenceReader[ServiceConfig, string](makeServiceInfo)

// Create specialized versions
authService := sequenced(ServiceConfig{Name: "Auth", Version: "1.0"})
userService := sequenced(ServiceConfig{Name: "User", Version: "2.0"})

// Reuse across contexts
authInfo := authService(ctx)()
userInfo := userService(ctx)()
```

## TraverseReader: Introducing Dependencies

While `SequenceReader` flips the parameter order of an existing nested structure, `TraverseReader` allows you to **introduce** a new Reader dependency into an existing computation.

### Function Signature

```go
func TraverseReader[R, A, B any](
    f reader.Kleisli[R, A, B],
) func(ReaderIOResult[A]) Kleisli[R, B]
```

**Type transformation:**
```
Input:  ReaderIOResult[A] = func(context.Context) func() Either[error, A]
With:   reader.Kleisli[R, A, B] = func(A) func(R) B
Output: Kleisli[R, B] = func(R) func(context.Context) func() Either[error, B]
```

### What It Does

`TraverseReader` takes:
1. A Reader-based transformation `f: func(A) func(R) B` that depends on environment `R`
2. Returns a function that transforms `ReaderIOResult[A]` into `Kleisli[R, B]`

This allows you to:
- Add environment dependencies to computations that don't have them yet
- Transform values within a ReaderIOResult using environment-dependent logic
- Build composable pipelines where transformations depend on configuration

### Key Difference from SequenceReader

- **SequenceReader**: Works with computations that **already contain** a Reader (`ReaderIOResult[Reader[R, A]]`)
  - Flips the order so `R` comes first
  - No transformation of the value itself

- **TraverseReader**: Works with computations that **don't have** a Reader yet (`ReaderIOResult[A]`)
  - Introduces a new Reader dependency via a transformation function
  - Transforms `A` to `B` using environment `R`

### Example: Adding Configuration to a Computation

```go
type Config struct {
    Multiplier int
    Prefix     string
}

// Original computation that just produces an int
getValue := func(ctx context.Context) func() Either[error, int] {
    return func() Either[error, int] {
        return Right[error](10)
    }
}

// A Reader-based transformation that depends on Config
formatWithConfig := func(n int) func(Config) string {
    return func(cfg Config) string {
        result := n * cfg.Multiplier
        return fmt.Sprintf("%s: %d", cfg.Prefix, result)
    }
}

// Use TraverseReader to introduce Config dependency
traversed := TraverseReader[Config, int, string](formatWithConfig)
withConfig := traversed(getValue)

// Now we can provide Config to get the final result
cfg := Config{Multiplier: 5, Prefix: "Result"}
ctx := t.Context()
result := withConfig(cfg)(ctx)() // Returns Right("Result: 50")
```

### Point-Free Composition with TraverseReader

```go
// Build a pipeline that introduces dependencies at each stage
var pipeline = F.Flow4(
    loadValue,                              // ReaderIOResult[int]
    TraverseReader(multiplyByConfig),       // Kleisli[Config, int]
    applyConfig(cfg),                       // ReaderIOResult[int]
    Chain(TraverseReader(formatWithStyle)), // Introduce another dependency
)
```

### When to Use TraverseReader vs SequenceReader

**Use SequenceReader when:**
- Your computation already returns a Reader: `ReaderIOResult[Reader[R, A]]`
- You just want to flip the parameter order
- No transformation of the value is needed

```go
// Already have Reader[Config, int]
computation := getComputation() // ReaderIOResult[Reader[Config, int]]
sequenced := SequenceReader[Config, int](computation)
result := sequenced(cfg)(ctx)()
```

**Use TraverseReader when:**
- Your computation doesn't have a Reader yet: `ReaderIOResult[A]`
- You want to transform the value using environment-dependent logic
- You're introducing a new dependency into the pipeline

```go
// Have ReaderIOResult[int], want to add Config dependency
computation := getValue() // ReaderIOResult[int]
traversed := TraverseReader[Config, int, string](formatWithConfig)
withDep := traversed(computation)
result := withDep(cfg)(ctx)()
```

### Practical Example: Multi-Stage Processing

```go
type DatabaseConfig struct {
    ConnectionString string
    Timeout          time.Duration
}

type FormattingConfig struct {
    DateFormat string
    Timezone   string
}

// Stage 1: Load raw data (no dependencies yet)
loadData := func(ctx context.Context) func() Either[error, RawData] {
    // ... implementation
}

// Stage 2: Process with database config
processWithDB := func(raw RawData) func(DatabaseConfig) ProcessedData {
    return func(cfg DatabaseConfig) ProcessedData {
        // Use cfg.ConnectionString, cfg.Timeout
        return ProcessedData{/* ... */}
    }
}

// Stage 3: Format with formatting config
formatData := func(processed ProcessedData) func(FormattingConfig) string {
    return func(cfg FormattingConfig) string {
        // Use cfg.DateFormat, cfg.Timezone
        return "formatted result"
    }
}

// Build pipeline introducing dependencies at each stage
var pipeline = F.Flow3(
    loadData,
    TraverseReader[DatabaseConfig, RawData, ProcessedData](processWithDB),
    // Now we have Kleisli[DatabaseConfig, ProcessedData]
    applyConfig(dbConfig),
    // Now we have ReaderIOResult[ProcessedData]
    TraverseReader[FormattingConfig, ProcessedData, string](formatData),
    // Now we have Kleisli[FormattingConfig, string]
)

// Execute with both configs
result := pipeline(fmtConfig)(ctx)()
```

### Combining TraverseReader and SequenceReader

You can combine both functions in complex pipelines:

```go
// Start with nested Reader
computation := getComputation() // ReaderIOResult[Reader[Config, User]]

var pipeline = F.Flow4(
    computation,
    SequenceReader[Config, User],           // Flip to get Kleisli[Config, User]
    applyConfig(cfg),                       // Apply config, get ReaderIOResult[User]
    TraverseReader(enrichWithDatabase),     // Add database dependency
    // Now have Kleisli[Database, EnrichedUser]
)

result := pipeline(db)(ctx)()
```

## Practical Benefits

### 1. **Performance: Eager Construction, Lazy Execution**

One of the most important but often overlooked benefits of point-free style is its performance characteristic: **the program structure is constructed eagerly (at definition time), but execution happens lazily (at runtime)**.

#### Construction Happens Once

When you define a pipeline using point-free style with `F.Flow`, `F.Pipe`, or function composition, the composition structure is built immediately at definition time:

```go
// Point-free style - composition built ONCE at definition time
var processUser = F.Flow3(
    getDatabase,
    SequenceReader[DatabaseConfig, Database],
    applyConfig(dbConfig),
)
// The pipeline structure is now fixed in memory
```

#### Execution Happens on Demand

The actual computation only runs when you provide the final parameters and invoke the result:

```go
// Execute multiple times - only execution cost, no re-composition
result1 := processUser(ctx1)()  // Fast - reuses pre-built pipeline
result2 := processUser(ctx2)()  // Fast - reuses pre-built pipeline
result3 := processUser(ctx3)()  // Fast - reuses pre-built pipeline
```

#### Performance Benefit for Repeated Execution

If a flow is executed multiple times, the point-free style is significantly more efficient because:

1. **Composition overhead is paid once** - The function composition happens at definition time
2. **No re-interpretation** - Each execution doesn't need to rebuild the pipeline
3. **Memory efficiency** - The composed function is created once and reused
4. **Better for hot paths** - Ideal for high-frequency operations

#### Comparison: Point-Free vs. Imperative

```go
// Imperative style - reconstruction on EVERY call
func processUserImperative(ctx context.Context) Either[error, Database] {
    // This function body is re-interpreted/executed every time
    dbComp := getDatabase()(ctx)()
    if dbReader, err := either.Unwrap(dbComp); err != nil {
        return Left[Database](err)
    }
    db := dbReader(dbConfig)
    // ... manual composition happens on every invocation
    return Right[error](db)
}

// Point-free style - composition built ONCE
var processUserPointFree = F.Flow3(
    getDatabase,
    SequenceReader[DatabaseConfig, Database],
    applyConfig(dbConfig),
)

// Benchmark scenario: 1000 executions
for i := 0; i < 1000; i++ {
    // Imperative: pays composition cost 1000 times
    result := processUserImperative(ctx)()
    
    // Point-free: pays composition cost once, execution cost 1000 times
    result := processUserPointFree(ctx)()
}
```

#### When This Matters Most

The performance benefit of eager construction is particularly important for:

- **High-frequency operations** - APIs, event handlers, request processors
- **Batch processing** - Same pipeline processes many items
- **Long-running services** - Pipelines defined once at startup, executed millions of times
- **Hot code paths** - Performance-critical sections that run repeatedly
- **Stream processing** - Processing continuous data streams

#### Example: API Handler

```go
// Define pipeline once at application startup
var handleUserRequest = F.Flow4(
    parseRequest,
    SequenceReader[Database, UserRequest],
    applyDatabase(db),
    Chain(validateAndProcess),
)

// Execute thousands of times per second
func apiHandler(w http.ResponseWriter, r *http.Request) {
    // No composition overhead - just execution
    result := handleUserRequest(r.Context())()
    // ... handle result
}
```

#### Memory and CPU Efficiency

```go
// Point-free: O(1) composition overhead
var pipeline = F.Flow5(step1, step2, step3, step4, step5)
// Composed once, stored in memory

// Execute N times: O(N) execution cost only
for i := 0; i < N; i++ {
    result := pipeline(input[i])
}

// Imperative: O(N) composition + execution cost
for i := 0; i < N; i++ {
    // Composition logic runs every iteration
    result := step5(step4(step3(step2(step1(input[i])))))
}
```

### 2. **Improved Testability**

Inject test dependencies easily:

```go
// Production
prodDB := Database{ConnectionString: "prod:5432"}
prodQuery := queryWithDB(prodDB)

// Testing
testDB := Database{ConnectionString: "test:5432"}
testQuery := queryWithDB(testDB)

// Same computation, different dependencies
```

### 3. **Better Separation of Concerns**

Separate configuration from execution:

```go
// Configuration phase (pure, no effects)
cfg := loadConfig()
computation := sequenced(cfg)

// Execution phase (with effects)
result := computation(ctx)()
```

### 4. **Enhanced Composability**

Build complex pipelines from simple pieces:

```go
var processUser = F.Flow4(
    loadUserConfig,           // ReaderIOResult[Reader[Database, User]]
    SequenceReader,           // func(Database) ReaderIOResult[User]
    applyDatabase(db),        // ReaderIOResult[User]
    Chain(validateUser),      // ReaderIOResult[ValidatedUser]
)
```

### 5. **Reduced Boilerplate**

No need to manually thread parameters:

```go
// Without Sequence - manual threading
func processWithConfig(cfg Config) ReaderIOResult[Result] {
    return func(ctx context.Context) func() Either[error, Result] {
        return func() Either[error, Result] {
            comp := getComputation()(ctx)()
            if reader, err := either.Unwrap(comp); err == nil {
                value := reader(cfg)
                // ... more processing
            }
            // ... error handling
        }
    }
}

// With Sequence - point-free
var processWithConfig = F.Flow2(
    getComputation,
    SequenceReader[Config, Result],
)
```

## Examples

### Example 1: Database Query with Configuration

```go
type QueryConfig struct {
    Timeout  time.Duration
    MaxRows  int
}

type Database struct {
    ConnectionString string
}

// Without Sequence
func executeQueryOld(cfg QueryConfig, db Database) ReaderIOResult[[]Row] {
    return func(ctx context.Context) func() Either[error, []Row] {
        return func() Either[error, []Row] {
            // Must manually handle all parameters
            // ...
        }
    }
}

// With Sequence
func makeQuery(ctx context.Context) func() Either[error, func(Database) []Row] {
    return func() Either[error, func(Database) []Row] {
        return Right[error](func(db Database) []Row {
            // Implementation
            return []Row{}
        })
    }
}

var executeQuery = F.Flow2(
    makeQuery,
    SequenceReader[Database, []Row],
)

// Usage
db := Database{ConnectionString: "localhost:5432"}
query := executeQuery(db)
result := query(ctx)()
```

### Example 2: Multi-Service Architecture

```go
type ServiceRegistry struct {
    AuthService  AuthService
    UserService  UserService
    EmailService EmailService
}

// Create computations that depend on services
makeAuthCheck := func(ctx context.Context) func() Either[error, func(ServiceRegistry) bool] {
    // ... implementation
}

makeSendEmail := func(ctx context.Context) func() Either[error, func(ServiceRegistry) error] {
    // ... implementation
}

// Sequence them
authCheck := SequenceReader[ServiceRegistry, bool](makeAuthCheck)
sendEmail := SequenceReader[ServiceRegistry, error](makeSendEmail)

// Inject services once
registry := ServiceRegistry{ /* ... */ }
checkAuth := authCheck(registry)
sendMail := sendEmail(registry)

// Use with different contexts
if isAuth, _ := either.Unwrap(checkAuth(ctx1)()); isAuth {
    sendMail(ctx2)()
}
```

### Example 3: Configuration-Driven Pipeline

```go
type PipelineConfig struct {
    Stage1Config Stage1Config
    Stage2Config Stage2Config
    Stage3Config Stage3Config
}

// Define stages
stage1 := SequenceReader[Stage1Config, IntermediateResult1](makeStage1)
stage2 := SequenceReader[Stage2Config, IntermediateResult2](makeStage2)
stage3 := SequenceReader[Stage3Config, FinalResult](makeStage3)

// Build pipeline with configuration
func buildPipeline(cfg PipelineConfig) ReaderIOResult[FinalResult] {
    return F.Pipe3(
        stage1(cfg.Stage1Config),
        Chain(func(r1 IntermediateResult1) ReaderIOResult[IntermediateResult2] {
            return stage2(cfg.Stage2Config)
        }),
        Chain(func(r2 IntermediateResult2) ReaderIOResult[FinalResult] {
            return stage3(cfg.Stage3Config)
        }),
    )
}

// Execute pipeline
cfg := loadPipelineConfig()
pipeline := buildPipeline(cfg)
result := pipeline(ctx)()
```

## Comparison: With and Without Sequence

### Without Sequence (Imperative Style)

```go
func processUser(userID string) ReaderIOResult[ProcessedUser] {
    return func(ctx context.Context) func() Either[error, ProcessedUser] {
        return func() Either[error, ProcessedUser] {
            // Get database
            dbComp := getDatabase()(ctx)()
            if dbReader, err := either.Unwrap(dbComp); err != nil {
                return Left[ProcessedUser](err)
            }
            db := dbReader(dbConfig)
            
            // Get user
            userComp := getUser(userID)(ctx)()
            if userReader, err := either.Unwrap(userComp); err != nil {
                return Left[ProcessedUser](err)
            }
            user := userReader(db)
            
            // Process user
            processComp := processUserData(user)(ctx)()
            if processReader, err := either.Unwrap(processComp); err != nil {
                return Left[ProcessedUser](err)
            }
            result := processReader(processingConfig)
            
            return Right[error](result)
        }
    }
}
```

### With Sequence (Point-Free Style)

```go
var processUser = func(userID string) ReaderIOResult[ProcessedUser] {
    return F.Pipe3(
        getDatabase,
        SequenceReader[DatabaseConfig, Database],
        applyConfig(dbConfig),
        Chain(func(db Database) ReaderIOResult[User] {
            return F.Pipe2(
                getUser(userID),
                SequenceReader[Database, User],
                applyDB(db),
            )
        }),
        Chain(func(user User) ReaderIOResult[ProcessedUser] {
            return F.Pipe2(
                processUserData(user),
                SequenceReader[ProcessingConfig, ProcessedUser],
                applyConfig(processingConfig),
            )
        }),
    )
}
```

## Key Takeaways

1. **Sequence functions flip parameter order** to enable partial application
2. **Dependencies come first**, making them easy to inject and test
3. **Point-free style** becomes natural and readable
4. **Composition** is enhanced through proper parameter ordering
5. **Reusability** increases as computations can be specialized early
6. **Testability** improves through easy dependency injection
7. **Separation of concerns** is clearer (configuration vs. execution)
8. **Performance benefit**: Eager construction (once) + lazy execution (many times) = efficiency for repeated operations

## When to Use Sequence

Use `Sequence*` functions when:

- ✅ You want to partially apply environment/configuration parameters
- ✅ You're building reusable computations with injected dependencies
- ✅ You need to test with different dependency implementations
- ✅ You're composing complex pipelines in point-free style
- ✅ You want to separate configuration from execution
- ✅ You're working with nested Reader-like structures

Don't use `Sequence*` when:

- ❌ The original parameter order is already optimal
- ❌ You're not doing any composition or partial application
- ❌ The added abstraction doesn't provide value
- ❌ The code is simpler without it

## Conclusion

The `Sequence*` functions are powerful tools for enabling point-free style programming in Go. By flipping the parameter order of nested monadic structures, they make it easy to:

- Partially apply dependencies
- Build composable pipelines
- Improve testability
- Write more declarative code

While they add a layer of abstraction, the benefits in terms of code reusability, testability, and composability make them invaluable for functional programming in Go.