# 🎯 ReaderResult: Context-Aware Functional Composition

## 📖 Overview

The [`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult) monad is a specialized implementation of the Reader monad pattern for Go, designed specifically for functions that:
- Depend on [`context.Context`](https://pkg.go.dev/context#Context) (for cancellation, deadlines, or context values)
- May fail with an error
- Need to be composed in a functional, declarative style

```go
type ReaderResult[A any] func(context.Context) (A, error)
```

This is equivalent to the common Go pattern `func(ctx context.Context) (A, error)`, but wrapped in a way that enables powerful functional composition.

## 🔄 ReaderResult as an Effectful Operation

**Important:** [`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult) represents an **effectful operation** because it depends on [`context.Context`](https://pkg.go.dev/context#Context), which is inherently mutable and can change during execution.

### Why Context Makes ReaderResult Effectful

The [`context.Context`](https://pkg.go.dev/context#Context) type in Go is designed to be mutable in the following ways:

1. **Cancellation State**: A context can transition from active to cancelled at any time
2. **Deadline Changes**: Timeouts and deadlines can expire during execution
3. **Value Storage**: Context values can be added or modified through [`context.WithValue`](https://pkg.go.dev/context#WithValue)
4. **Parent-Child Relationships**: Derived contexts inherit and can override parent behavior

This mutability means that:
- **Same ReaderResult, Different Results**: Executing the same [`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult) with different contexts (or the same context at different times) can produce different outcomes
- **Non-Deterministic Behavior**: Context cancellation or timeout can interrupt execution at any point
- **Side Effects**: The context carries runtime state that affects computation behavior

### Use Case Examples

#### 1. **HTTP Request Handling with Timeout**
```go
// The context carries a deadline that can expire during execution
func FetchUserProfile(userID int) ReaderResult[UserProfile] {
    return func(ctx context.Context) (UserProfile, error) {
        // Context deadline affects when this operation fails
        req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("/users/%d", userID), nil)
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return UserProfile{}, err // May fail due to context timeout
        }
        defer resp.Body.Close()
        
        var profile UserProfile
        json.NewDecoder(resp.Body).Decode(&profile)
        return profile, nil
    }
}

// Same function, different contexts = different behavior
ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel1()
profile1, _ := FetchUserProfile(123)(ctx1) // Has 5 seconds to complete

ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel2()
profile2, _ := FetchUserProfile(123)(ctx2) // Has only 100ms - likely to timeout
```

#### 2. **Database Transactions with Cancellation**
```go
// Context cancellation can abort the transaction at any point
func TransferFunds(from, to int, amount float64) ReaderResult[Transaction] {
    return func(ctx context.Context) (Transaction, error) {
        tx, err := db.BeginTx(ctx, nil) // Context controls transaction lifetime
        if err != nil {
            return Transaction{}, err
        }
        defer tx.Rollback()
        
        // If context is cancelled here, the debit fails
        if err := debitAccount(ctx, tx, from, amount); err != nil {
            return Transaction{}, err
        }
        
        // Or cancellation could happen here, before credit
        if err := creditAccount(ctx, tx, to, amount); err != nil {
            return Transaction{}, err
        }
        
        if err := tx.Commit(); err != nil {
            return Transaction{}, err
        }
        
        return Transaction{From: from, To: to, Amount: amount}, nil
    }
}

// User cancels the request mid-transaction
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(50 * time.Millisecond)
    cancel() // Cancellation affects the running operation
}()
result, err := TransferFunds(100, 200, 50.0)(ctx) // May be interrupted
```

#### 3. **Context Values for Request Tracing**
```go
// Context values affect logging and tracing behavior
func ProcessOrder(orderID string) ReaderResult[Order] {
    return func(ctx context.Context) (Order, error) {
        // Extract trace ID from context (mutable state)
        traceID := ctx.Value("trace-id")
        log.Printf("[%v] Processing order %s", traceID, orderID)
        
        // The same function behaves differently based on context values
        if ctx.Value("debug") == true {
            log.Printf("[%v] Debug mode: detailed order processing", traceID)
        }
        
        return fetchOrder(ctx, orderID)
    }
}

// Different contexts = different tracing behavior
ctx1 := context.WithValue(context.Background(), "trace-id", "req-001")
order1, _ := ProcessOrder("ORD-123")(ctx1) // Logs with trace-id: req-001

ctx2 := context.WithValue(context.Background(), "trace-id", "req-002")
ctx2 = context.WithValue(ctx2, "debug", true)
order2, _ := ProcessOrder("ORD-123")(ctx2) // Logs with trace-id: req-002 + debug info
```

#### 4. **Parallel Operations with Shared Cancellation**
```go
// Multiple operations share the same cancellable context
func FetchDashboardData(userID int) ReaderResult[Dashboard] {
    return func(ctx context.Context) (Dashboard, error) {
        // All these operations can be cancelled together
        userCh := make(chan User)
        postsCh := make(chan []Post)
        statsCh := make(chan Stats)
        errCh := make(chan error, 3)
        
        go func() {
            user, err := FetchUser(userID)(ctx) // Shares cancellation
            if err != nil {
                errCh <- err
                return
            }
            userCh <- user
        }()
        
        go func() {
            posts, err := FetchPosts(userID)(ctx) // Shares cancellation
            if err != nil {
                errCh <- err
                return
            }
            postsCh <- posts
        }()
        
        go func() {
            stats, err := FetchStats(userID)(ctx) // Shares cancellation
            if err != nil {
                errCh <- err
                return
            }
            statsCh <- stats
        }()
        
        // If context is cancelled, all goroutines stop
        select {
        case err := <-errCh:
            return Dashboard{}, err
        case <-ctx.Done():
            return Dashboard{}, ctx.Err() // Context cancellation is an effect
        case user := <-userCh:
            // ... collect results
        }
    }
}
```

#### 5. **Retry Logic with Context Awareness**
```go
import (
    "time"
    R "github.com/IBM/fp-go/v2/retry"
)

// Context state affects retry behavior using the built-in Retrying method
func FetchWithRetry[A any](operation ReaderResult[A]) ReaderResult[A] {
    // Create a retry policy: exponential backoff with a cap, limited to 5 retries
    policy := R.Monoid.Concat(
        R.LimitRetries(5),
        R.CapDelay(10*time.Second, R.ExponentialBackoff(100*time.Millisecond)),
    )
    
    // Check function: retry on any error
    // Note: context cancellation is automatically handled by Retrying
    shouldRetry := func(val A, err error) bool {
        return err != nil
    }
    
    // Use the built-in Retrying method with automatic context awareness
    // R.Always creates a constant function that ignores RetryStatus and always returns the operation
    return Retrying(policy, R.Always(operation), shouldRetry)
}

// Example usage:
// fetchUser := FetchWithRetry(GetUser(123))
// user, err := fetchUser(ctx) // Automatically retries with exponential backoff
//                              // and respects context cancellation
```

### Key Takeaway

Because `ReaderResult` depends on the mutable `context.Context`, it represents **effectful computations** where:
- Execution behavior can change based on context state
- The same operation can produce different results with different contexts
- External factors (timeouts, cancellations, context values) influence outcomes
- Side effects are inherent due to context's runtime nature

This makes [`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult) ideal for modeling real-world operations that interact with external systems, respect cancellation, and need to be composed in a functional style while acknowledging their effectful nature.

## 🤔 Why Use ReaderResult Instead of Traditional Go Methods?

### 1. ✨ **Simplified API Design**

**Traditional Go approach:**
```go
func GetUser(ctx context.Context, id int) (User, error)
func GetPosts(ctx context.Context, userID int) ([]Post, error)
func FormatUser(ctx context.Context, user User) (string, error)
```

Every function must explicitly accept and thread [`context.Context`](https://pkg.go.dev/context#Context) through the call chain, leading to repetitive boilerplate.

**ReaderResult approach:**
```go
func GetUser(id int) ReaderResult[User]
func GetPosts(userID int) ReaderResult[[]Post]
func FormatUser(user User) ReaderResult[string]
```

The context dependency is implicit in the return type. Functions are cleaner and focus on their core logic.

### 2. 🔗 **Composability Through Monadic Operations**

**Traditional Go approach:**
```go
func GetUserWithPosts(ctx context.Context, userID int) (UserWithPosts, error) {
    user, err := GetUser(ctx, userID)
    if err != nil {
        return UserWithPosts{}, err
    }
    
    posts, err := GetPosts(ctx, user.ID)
    if err != nil {
        return UserWithPosts{}, err
    }
    
    formatted, err := FormatUser(ctx, user)
    if err != nil {
        return UserWithPosts{}, err
    }
    
    return UserWithPosts{
        User:      user,
        Posts:     posts,
        Formatted: formatted,
    }, nil
}
```

Manual error handling at every step, repetitive context threading, and imperative style.

**ReaderResult approach:**
```go
func GetUserWithPosts(userID int) ReaderResult[UserWithPosts] {
    return F.Pipe3(
        Do(UserWithPosts{}),
        Bind(setUser, func(s UserWithPosts) ReaderResult[User] {
            return GetUser(userID)
        }),
        Bind(setPosts, func(s UserWithPosts) ReaderResult[[]Post] {
            return GetPosts(s.User.ID)
        }),
        Bind(setFormatted, func(s UserWithPosts) ReaderResult[string] {
            return FormatUser(s.User)
        }),
    )
}
```

Declarative pipeline, automatic error propagation, and clear data flow.

### 3. 🎨 **Pure Composition - Side Effects Deferred**

**💡 Key Insight:** [`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult) separates *building* computations from *executing* them.

```go
// Building the computation (pure, no side effects)
getUserPipeline := F.Pipe2(
    GetUser(123),
    Chain(func(user User) ReaderResult[[]Post] {
        return GetPosts(user.ID)
    }),
    Map(len[[]Post]),
)

// Execution happens later, at the edge of your system
postCount, err := getUserPipeline(ctx)
```

**Benefits:**
- Computations can be built, tested, and reasoned about without executing side effects
- Easy to mock and test individual components
- Clear separation between business logic and execution
- Computations are reusable with different contexts

### 4. 🧪 **Improved Testability**

**Traditional approach:**
```go
func TestGetUserWithPosts(t *testing.T) {
    // Need to mock database, HTTP clients, etc.
    // Tests are tightly coupled to implementation
    ctx := context.Background()
    result, err := GetUserWithPosts(ctx, 123)
    // ...
}
```

**ReaderResult approach:**
```go
func TestGetUserWithPosts(t *testing.T) {
    // Test the composition logic without executing side effects
    pipeline := GetUserWithPosts(123)
    
    // Can test with a mock context that provides test data
    testCtx := context.WithValue(context.Background(), "test", true)
    result, err := pipeline(testCtx)
    
    // Or test individual components in isolation
    mockGetUser := func(id int) ReaderResult[User] {
        return Of(User{ID: id, Name: "Test User"})
    }
}
```

You can test the composition logic separately from the actual I/O operations.

### 5. 📝 **Better Error Context Accumulation**

ReaderResult makes it easy to add context to errors as they propagate:

```go
getUserWithContext := F.Pipe2(
    GetUser(userID),
    MapError(func(err error) error {
        return fmt.Errorf("failed to get user %d: %w", userID, err)
    }),
    Chain(func(user User) ReaderResult[UserWithPosts] {
        return F.Pipe1(
            GetPosts(user.ID),
            MapError(func(err error) error {
                return fmt.Errorf("failed to get posts for user %s: %w", user.Name, err)
            }),
            Map(func(posts []Post) UserWithPosts {
                return UserWithPosts{User: user, Posts: posts}
            }),
        )
    }),
)
```

Errors automatically accumulate context as they bubble up through the composition.

### 6. ⚡ **Natural Parallel Execution**

With applicative functors, independent operations can be expressed naturally:

```go
// These operations don't depend on each other
getUserData := F.Pipe2(
    Do(UserData{}),
    ApS(setUser, GetUser(userID)),      // Can run in parallel
    ApS(setSettings, GetSettings(userID)), // Can run in parallel
    ApS(setPreferences, GetPreferences(userID)), // Can run in parallel
)
```

The structure makes it clear which operations are independent, enabling potential optimization.

### 7. 🔄 **Retry and Recovery Patterns**

ReaderResult makes retry logic composable. For production use, leverage the built-in [`Retrying`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Retrying) function with configurable retry policies:

```go
import (
    "time"
    R "github.com/IBM/fp-go/v2/retry"
)

// Production-ready retry with exponential backoff and context awareness
func WithRetry[A any](operation ReaderResult[A]) ReaderResult[A] {
    // Create a retry policy: exponential backoff with a cap, limited to 5 retries
    policy := R.Monoid.Concat(
        R.LimitRetries(5),
        R.CapDelay(10*time.Second, R.ExponentialBackoff(100*time.Millisecond)),
    )
    
    // Retry on any error (context cancellation is automatically handled)
    shouldRetry := func(val A, err error) bool {
        return err != nil
    }
    
    // Use built-in Retrying with automatic context cancellation support
    return Retrying(policy, R.Always(operation), shouldRetry)
}

// Use it:
reliableGetUser := WithRetry(GetUser(userID))
user, err := reliableGetUser(ctx) // Automatically retries with exponential backoff

// Or for simple cases, implement custom retry logic:
func SimpleRetry[A any](maxAttempts int, operation ReaderResult[A]) ReaderResult[A] {
    return func(ctx context.Context) (A, error) {
        var lastErr error
        for i := 0; i < maxAttempts; i++ {
            result, err := operation(ctx)
            if err == nil {
                return result, nil
            }
            lastErr = err
            time.Sleep(time.Second * time.Duration(i+1))
        }
        return *new(A), fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
    }
}
```

### 8. 🎭 **Middleware/Aspect-Oriented Programming**

Cross-cutting concerns can be added as higher-order functions:

```go
import (
    "time"
    F "github.com/IBM/fp-go/v2/function"
    RR "github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
)

// Logging middleware
func WithLogging[A any](name string) RR.Operator[A, A] {
    return func(operation RR.ReaderResult[A]) RR.ReaderResult[A] {
        return func(ctx context.Context) (A, error) {
            log.Printf("Starting %s", name)
            start := time.Now()
            result, err := operation(ctx)
            log.Printf("Finished %s in %v (error: %v)", name, time.Since(start), err)
            return result, err
        }
    }
}

// Compose middleware using built-in functions:
robustGetUser := F.Pipe3(
    GetUser(userID),
    WithLogging[User]("GetUser"),
    RR.WithTimeout[User](5 * time.Second),  // Built-in timeout support
    WithRetry[User],
)
```

**Built-in Middleware Functions:**
- [`WithTimeout`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#WithTimeout) - Add timeout to operations
- [`WithDeadline`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#WithDeadline) - Add absolute deadline to operations
- [`Local`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Local) - Transform context for specific operations

### 9. 🛡️ **Resource Management with Bracket**

Safe resource handling with guaranteed cleanup:

```go
readFile := Bracket(
    // Acquire
    func() ReaderResult[*os.File] {
        return func(ctx context.Context) (*os.File, error) {
            return os.Open("data.txt")
        }
    },
    // Use
    func(file *os.File) ReaderResult[string] {
        return func(ctx context.Context) (string, error) {
            data, err := io.ReadAll(file)
            return string(data), err
        }
    },
    // Release (always called)
    func(file *os.File, content string, err error) ReaderResult[any] {
        return func(ctx context.Context) (any, error) {
            return nil, file.Close()
        }
    },
)
```

The bracket pattern ensures resources are always cleaned up, even on errors.

### 10. 🔒 **Type-Safe State Threading**

Do-notation provides type-safe accumulation of state:

```go
type UserProfile struct {
    User     User
    Posts    []Post
    Comments []Comment
    Stats    Statistics
}

buildProfile := F.Pipe4(
    Do(UserProfile{}),
    Bind(setUser, func(s UserProfile) ReaderResult[User] {
        return GetUser(userID)
    }),
    Bind(setPosts, func(s UserProfile) ReaderResult[[]Post] {
        return GetPosts(s.User.ID) // Can access previous results
    }),
    Bind(setComments, func(s UserProfile) ReaderResult[[]Comment] {
        return GetComments(s.User.ID)
    }),
    Bind(setStats, func(s UserProfile) ReaderResult[Statistics] {
        return CalculateStats(s.Posts, s.Comments) // Can use multiple previous results
    }),
)
```

The compiler ensures you can't access fields that haven't been set yet.

### 🔍 **Using Optics with Bind and BindTo**

[Optics](../../../optics/README.md) provide powerful, composable abstractions for working with data structures. They integrate seamlessly with [`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult)'s [`Bind`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Bind) and [`BindTo`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Bind) methods, enabling elegant state accumulation patterns.

#### Lenses for Product Types (Structs)

[Lenses](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens) focus on struct fields and can be used as setters in [`Bind`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Bind) operations:

```go
import (
    "github.com/IBM/fp-go/v2/optics/lens"
    F "github.com/IBM/fp-go/v2/function"
)

type User struct {
    ID   int
    Name string
}

type UserProfile struct {
    User     User
    Posts    []Post
    Comments []Comment
}

// Auto-generated or manually created lenses
var (
    userLens = lens.MakeLens(
        func(p UserProfile) User { return p.User },
        func(p UserProfile, u User) UserProfile {
            p.User = u
            return p
        },
    )
    postsLens = lens.MakeLens(
        func(p UserProfile) []Post { return p.Posts },
        func(p UserProfile, posts []Post) UserProfile {
            p.Posts = posts
            return p
        },
    )
    commentsLens = lens.MakeLens(
        func(p UserProfile) []Comment { return p.Comments },
        func(p UserProfile, comments []Comment) UserProfile {
            p.Comments = comments
            return p
        },
    )
    
    // Lens for User.ID field
    userIDLens = lens.MakeLens(
        func(u User) int { return u.ID },
        func(u User, id int) User {
            u.ID = id
            return u
        },
    )
    
    // Composed lens: UserProfile -> User -> ID
    // This demonstrates lens composition - a key benefit of optics!
    profileUserIDLens = F.Pipe1(
        userLens,
        lens.Compose[UserProfile](userIDLens),
    )
)

// Use lenses as setters in Bind
buildProfile := F.Pipe3(
    Do(UserProfile{}),
    Bind(userLens.Set, func(s UserProfile) ReaderResult[User] {
        return GetUser(userID)
    }),
    Bind(postsLens.Set, func(s UserProfile) ReaderResult[[]Post] {
        // Use composed lens to access nested User.ID directly
        return GetPosts(profileUserIDLens.Get(s))
    }),
    Bind(commentsLens.Set, func(s UserProfile) ReaderResult[[]Comment] {
        // Composed lens makes nested access clean and type-safe
        return GetComments(profileUserIDLens.Get(s))
    }),
)
```

**Benefits:**
- Type-safe field updates
- Reusable lens definitions
- Clear separation of data access from business logic
- Can be auto-generated with `go generate` (see [optics documentation](../../../optics/README.md#-auto-generation-with-go-generate))

#### Prisms for Sum Types (Variants)

[Prisms](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism) are particularly powerful in [`Bind`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Bind) operations as they act as **generalized constructors**. The prism's `ReverseGet` function constructs values of sum types, making them ideal for building up complex results:

```go
import (
    "github.com/IBM/fp-go/v2/optics/prism"
    O "github.com/IBM/fp-go/v2/option"
)

type APIResponse struct {
    UserData   O.Option[User]
    PostsData  O.Option[[]Post]
    StatsData  O.Option[Stats]
}

// Prisms for Option fields - ReverseGet acts as a constructor
var (
    userDataPrism = prism.MakePrism(
        func(r APIResponse) O.Option[User] { return r.UserData },
        func(u User) APIResponse {
            return APIResponse{UserData: O.Some(u)}
        },
    )
    postsDataPrism = prism.MakePrism(
        func(r APIResponse) O.Option[[]Post] { return r.PostsData },
        func(posts []Post) APIResponse {
            return APIResponse{PostsData: O.Some(posts)}
        },
    )
)

// Use prisms to construct and accumulate optional data
fetchAPIData := F.Pipe3(
    Do(APIResponse{}),
    // ReverseGet constructs APIResponse with UserData set
    Bind(userDataPrism.ReverseGet, func(s APIResponse) ReaderResult[User] {
        return GetUser(userID)
    }),
    // ReverseGet constructs APIResponse with PostsData set
    Bind(postsDataPrism.ReverseGet, func(s APIResponse) ReaderResult[[]Post] {
        return O.Fold(
            func() ReaderResult[[]Post] { return Of([]Post{}) },
            func(user User) ReaderResult[[]Post] { return GetPosts(user.ID) },
        )(s.UserData)
    }),
    // Handle optional stats
    BindTo(func(s APIResponse) ReaderResult[APIResponse] {
        return F.Pipe1(
            GetStats(userID),
            Map(func(stats Stats) APIResponse {
                s.StatsData = O.Some(stats)
                return s
            }),
        )
    }),
)
```

**Why Prisms Excel in Bind:**
- **Generalized Constructors**: `ReverseGet` creates values from variants, perfect for building sum types
- **Partial Construction**: Build complex structures incrementally
- **Type Safety**: Compiler ensures correct variant handling
- **Composability**: Prisms compose naturally with monadic operations

#### Combining Lenses and Prisms

For maximum flexibility, combine both optics in a single pipeline:

```go
type ComplexState struct {
    Config   Config
    Result   O.Option[ProcessingResult]
    Metadata Metadata
}

var (
    configLens = lens.MakeLens(
        func(s ComplexState) Config { return s.Config },
        func(s ComplexState, c Config) ComplexState {
            s.Config = c
            return s
        },
    )
    resultPrism = prism.MakePrism(
        func(s ComplexState) O.Option[ProcessingResult] { return s.Result },
        func(r ProcessingResult) ComplexState {
            return ComplexState{Result: O.Some(r)}
        },
    )
)

pipeline := F.Pipe3(
    Do(ComplexState{}),
    Bind(configLens.Set, LoadConfig),           // Lens for required field
    Bind(resultPrism.ReverseGet, ProcessData),  // Prism for optional result
    Bind(metadataLens.Set, ExtractMetadata),    // Lens for metadata
)
```

**Learn More:**
- [Optics Overview](../../optics/README.md) - Complete guide to lenses, prisms, and other optics
- [Lens Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens) - Detailed lens API
- [Prism Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism) - Prism patterns and usage
- [Auto-generation](../../optics/README.md#-auto-generation-with-go-generate) - Generate optics automatically

### 11. ⏱️ **Automatic Context Cancellation Checks**

ReaderResult automatically checks for context cancellation at composition boundaries through `WithContextK`, ensuring fail-fast behavior without manual checks:

**Traditional approach:**
```go
func ProcessData(ctx context.Context, data Data) (Result, error) {
    // Manual cancellation check
    if ctx.Err() != nil {
        return Result{}, ctx.Err()
    }
    
    step1, err := Step1(ctx, data)
    if err != nil {
        return Result{}, err
    }
    
    // Manual cancellation check again
    if ctx.Err() != nil {
        return Result{}, ctx.Err()
    }
    
    step2, err := Step2(ctx, step1)
    if err != nil {
        return Result{}, err
    }
    
    // And again...
    if ctx.Err() != nil {
        return Result{}, ctx.Err()
    }
    
    return Step3(ctx, step2)
}
```

**ReaderResult approach:**
```go
func ProcessData(data Data) ReaderResult[Result] {
    return F.Pipe3(
        Step1(data),
        Chain(Step2),  // Automatic cancellation check before Step2
        Chain(Step3),  // Automatic cancellation check before Step3
    )
}
```

**How it works:**
- All `Bind` operations use `WithContextK` internally
- `WithContextK` wraps each Kleisli arrow with a cancellation check
- Before executing each step, it checks `ctx.Err()` and fails fast if cancelled
- No manual cancellation checks needed in your business logic
- Ensures long-running pipelines respect context cancellation at every step

**Example with timeout:**
```go
pipeline := F.Pipe3(
    FetchUser(userID),        // Step 1
    Chain(FetchPosts),        // Cancellation checked before Step 2
    Chain(EnrichWithMetadata), // Cancellation checked before Step 3
)

// Set a timeout
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

// If Step 1 takes too long, Steps 2 and 3 won't execute
result, err := pipeline(ctx)
```

This makes ReaderResult ideal for:
- Long-running pipelines that should respect timeouts
- Operations that need to be cancellable at any point
- Composing third-party functions that don't check context themselves
- Building responsive services that handle request cancellation properly

## 🎯 When to Use ReaderResult

**✅ Use ReaderResult when:**
- You have complex composition of context-dependent operations
- You want to separate business logic from execution
- You need better testability and mockability
- You want declarative, pipeline-style code
- You need to add cross-cutting concerns (logging, retry, timeout)
- You want type-safe state accumulation
- You need automatic context cancellation checks at composition boundaries
- You're building long-running pipelines that should respect timeouts

**❌ Stick with traditional Go when:**
- You have simple, one-off operations
- The team is unfamiliar with functional patterns
- You're writing library code that needs to be idiomatic Go
- Performance is absolutely critical (though the overhead is minimal)

## 🚀 Quick Start

```go
import (
    "context"
    F "github.com/IBM/fp-go/v2/function"
    RR "github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
)

// Define your operations
func GetUser(id int) RR.ReaderResult[User] {
    return func(ctx context.Context) (User, error) {
        // Your implementation
    }
}

// Compose them
pipeline := F.Pipe2(
    GetUser(123),
    RR.Chain(func(user User) RR.ReaderResult[[]Post] {
        return GetPosts(user.ID)
    }),
    RR.Map(CreateSummary),
)

// Execute at the edge
summary, err := pipeline(context.Background())
```

**Key functions used:**
- [`Chain`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Chain) - Sequence dependent operations
- [`Map`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Map) - Transform success values

## 🔄 Converting Traditional Go Functions

[`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult) provides convenient functions to convert traditional Go functions (that take [`context.Context`](https://pkg.go.dev/context#Context) as their first parameter) into functional ReaderResult operations. This makes it easy to integrate existing code into functional pipelines.

### Using `FromXXX` Functions (Uncurried)

The [`From0`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From0), [`From1`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From1), [`From2`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From2), [`From3`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From3) functions convert traditional Go functions into ReaderResult-returning functions that take all parameters at once. This is the most straightforward conversion for direct use.

```go
// Traditional Go function
func getUser(ctx context.Context, id int) (User, error) {
    // ... database query
    return User{ID: id, Name: "Alice"}, nil
}

func updateUser(ctx context.Context, id int, name string) (User, error) {
    // ... database update
    return User{ID: id, Name: name}, nil
}

// Convert using From1 (1 parameter besides context)
getUserRR := RR.From1(getUser)

// Convert using From2 (2 parameters besides context)
updateUserRR := RR.From2(updateUser)

// Use in a pipeline
pipeline := F.Pipe2(
    getUserRR(123),                    // Returns ReaderResult[User]
    RR.Chain(func(user User) RR.ReaderResult[User] {
        return updateUserRR(user.ID, "Bob")  // All params at once
    }),
)

result, err := pipeline(ctx)
```

**Available From functions:**
- `From0`: Converts `func(context.Context) (A, error)` → `func() ReaderResult[A]`
- `From1`: Converts `func(context.Context, T1) (A, error)` → `func(T1) ReaderResult[A]`
- `From2`: Converts `func(context.Context, T1, T2) (A, error)` → `func(T1, T2) ReaderResult[A]`
- `From3`: Converts `func(context.Context, T1, T2, T3) (A, error)` → `func(T1, T2, T3) ReaderResult[A]`

### Using `CurryXXX` Functions (Curried)

The [`Curry0`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry0), [`Curry1`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry1), [`Curry2`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry2), [`Curry3`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry3) functions convert traditional Go functions into curried ReaderResult-returning functions. This enables partial application, which is useful for building reusable function pipelines.

```go
// Traditional Go function
func createPost(ctx context.Context, userID int, title string, body string) (Post, error) {
    return Post{UserID: userID, Title: title, Body: body}, nil
}

// Convert using Curry3 (3 parameters besides context)
createPostRR := RR.Curry3(createPost)

// Partial application - build specialized functions
createPostForUser42 := createPostRR(42)
createPostWithTitle := createPostForUser42("My Title")

// Complete the application
rr := createPostWithTitle("Post body content")
post, err := rr(ctx)

// Or apply all at once
post2, err := createPostRR(42)("Another Title")("Another body")(ctx)
```

**Available Curry functions:**
- `Curry0`: Converts `func(context.Context) (A, error)` → `ReaderResult[A]`
- `Curry1`: Converts `func(context.Context, T1) (A, error)` → `func(T1) ReaderResult[A]`
- `Curry2`: Converts `func(context.Context, T1, T2) (A, error)` → `func(T1) func(T2) ReaderResult[A]`
- `Curry3`: Converts `func(context.Context, T1, T2, T3) (A, error)` → `func(T1) func(T2) func(T3) ReaderResult[A]`

### Practical Example: Integrating Existing Code

```go
// Existing traditional Go functions (e.g., from a database package)
func fetchUser(ctx context.Context, id int) (User, error) { /* ... */ }
func fetchPosts(ctx context.Context, userID int) ([]Post, error) { /* ... */ }
func fetchComments(ctx context.Context, postID int) ([]Comment, error) { /* ... */ }

// Convert them all to ReaderResult
var (
    GetUser     = RR.From1(fetchUser)
    GetPosts    = RR.From1(fetchPosts)
    GetComments = RR.From1(fetchComments)
)

// Now compose them functionally
func GetUserWithData(userID int) RR.ReaderResult[UserData] {
    return F.Pipe3(
        GetUser(userID),
        RR.Chain(func(user User) RR.ReaderResult[[]Post] {
            return GetPosts(user.ID)
        }),
        RR.Map(func(posts []Post) UserData {
            return UserData{
                User:  user,
                Posts: posts,
            }
        }),
    )
}

// Execute
userData, err := GetUserWithData(123)(ctx)
```

### When to Use From vs Curry

**Use `FromXXX` when:**
- You want straightforward conversion for immediate use
- You're calling functions with all parameters at once
- You prefer a more familiar, uncurried style
- You're converting functions for one-time use in a pipeline

**Use `CurryXXX` when:**
- You want to partially apply parameters
- You're building reusable, specialized functions
- You need maximum composability
- You're working with higher-order functions that expect curried inputs

### Converting Back: Uncurry Functions

If you need to convert a [`ReaderResult`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ReaderResult) function back to traditional Go style (e.g., for interfacing with non-functional code), use the [`Uncurry1`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Uncurry1), [`Uncurry2`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Uncurry2), [`Uncurry3`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Uncurry3) functions:

```go
// Functional style
getUserRR := func(id int) RR.ReaderResult[User] {
    return func(ctx context.Context) (User, error) {
        return User{ID: id}, nil
    }
}

// Convert back to traditional Go
getUser := RR.Uncurry1(getUserRR)

// Now callable in traditional style
user, err := getUser(ctx, 123)
```

## 📚 API Reference

### Core Functions
- [`Map`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Map) - Transform the success value
- [`Chain`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Chain) - Sequence operations (also known as `FlatMap` or `Bind`)
- [`Bind`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Bind) - Do-notation binding for state accumulation
- [`Do`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Do) - Start a Do-notation pipeline
- [`ApS`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#ApS) - Applicative sequencing for parallel operations
- [`Of`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Of) - Lift a pure value into ReaderResult

### Resource Management
- [`Bracket`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Bracket) - Safe resource acquisition and cleanup

### Conversion Functions
- [`From0`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From0), [`From1`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From1), [`From2`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From2), [`From3`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#From3) - Convert traditional Go functions to ReaderResult
- [`Curry0`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry0), [`Curry1`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry1), [`Curry2`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry2), [`Curry3`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Curry3) - Convert to curried ReaderResult functions
- [`Uncurry1`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Uncurry1), [`Uncurry2`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Uncurry2), [`Uncurry3`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#Uncurry3) - Convert back to traditional Go functions

### Error Handling
- [`MapError`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#MapError) - Transform error values
- [`OrElse`](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult#OrElse) - Provide fallback on error

### Full Documentation
- [Package Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/idiomatic/context/readerresult) - Complete API reference

## 📚 See Also

- [bind.go](bind.go) - Do-notation and composition operators
- [bracket.go](bracket.go) - Resource management patterns
- [examples_bind_test.go](examples_bind_test.go) - Comprehensive examples