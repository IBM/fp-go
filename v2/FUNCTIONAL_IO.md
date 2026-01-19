# Functional I/O in Go: Context, Errors, and the Reader Pattern

This document explores how functional programming principles apply to I/O operations in Go, comparing traditional imperative approaches with functional patterns using the `context/readerioresult` and `idiomatic/context/readerresult` packages.

## Table of Contents

- [Why Context in I/O Operations](#why-context-in-io-operations)
- [The Error-Value Tuple Pattern](#the-error-value-tuple-pattern)
- [Functional Approach: Reader Pattern](#functional-approach-reader-pattern)
- [Benefits of the Functional Approach](#benefits-of-the-functional-approach)
- [Side-by-Side Comparison](#side-by-side-comparison)
- [Advanced Patterns](#advanced-patterns)
- [When to Use Each Approach](#when-to-use-each-approach)

## Why Context in I/O Operations

In idiomatic Go, I/O operations conventionally take a `context.Context` as their first parameter:

```go
func QueryDatabase(ctx context.Context, query string) (Result, error)
func MakeHTTPRequest(ctx context.Context, url string) (*http.Response, error)
func ReadFile(ctx context.Context, path string) ([]byte, error)
```

### The Purpose of Context

The `context.Context` parameter serves several critical purposes:

1. **Cancellation Propagation**: Operations can be cancelled when the context is cancelled
2. **Deadline Management**: Operations respect timeouts and deadlines
3. **Request-Scoped Values**: Carry request metadata (trace IDs, user info, etc.)
4. **Resource Cleanup**: Signal to release resources when work is no longer needed

### Why Context Matters for I/O

I/O operations are inherently **effectful** - they interact with the outside world:
- Reading from disk, network, or database
- Writing to external systems
- Generating random numbers
- Reading the current time

These operations can:
- **Take time**: Network calls may be slow
- **Fail**: Connections drop, files don't exist
- **Block**: Waiting for external resources
- **Need cancellation**: User navigates away, request times out

Context provides a standard mechanism to control these operations across your entire application.

## The Error-Value Tuple Pattern

### Why Operations Must Return Errors

In Go, I/O operations return `(value, error)` tuples because:

1. **Context can be cancelled**: Even if the operation would succeed, cancellation must be represented
2. **External systems fail**: Networks fail, files are missing, permissions are denied
3. **Resources are exhausted**: Out of memory, disk full, connection pool exhausted
4. **Timeouts occur**: Operations exceed their deadline

**There cannot be I/O operations without error handling** because the context itself introduces a failure mode (cancellation) that must be represented in the return type.

### Traditional Go Pattern

```go
func ProcessUser(ctx context.Context, userID int) (User, error) {
    // Check context before starting
    if err := ctx.Err(); err != nil {
        return User{}, err
    }
    
    // Fetch user from database
    user, err := db.QueryUser(ctx, userID)
    if err != nil {
        return User{}, fmt.Errorf("query user: %w", err)
    }
    
    // Validate user
    if user.Age < 18 {
        return User{}, errors.New("user too young")
    }
    
    // Fetch user's posts
    posts, err := db.QueryPosts(ctx, user.ID)
    if err != nil {
        return User{}, fmt.Errorf("query posts: %w", err)
    }
    
    user.Posts = posts
    return user, nil
}
```

**Characteristics:**
- Explicit error checking at each step
- Manual error wrapping and propagation
- Context checked manually
- Imperative control flow
- Error handling mixed with business logic

## Functional Approach: Reader Pattern

### The Core Insight

In functional programming, we separate **what to compute** from **how to execute it**. Instead of functions that perform I/O directly, we create functions that **return descriptions of I/O operations**.

### Key Type: ReaderIOResult

```go
// A function that takes a context and returns a value or error
type ReaderIOResult[A any] = func(context.Context) (A, error)
```

This type represents:
- **Reader**: Depends on an environment (context.Context)
- **IO**: Performs side effects (I/O operations)
- **Result**: Can fail with an error

### Why This Is Better

The functional approach **carries the I/O aspect as the return value, not on the input**:

```go
// Traditional: I/O is implicit in the function execution
func fetchUser(ctx context.Context, id int) (User, error) {
    // Performs I/O immediately
}

// Functional: I/O is explicit in the return type
func fetchUser(id int) ReaderIOResult[User] {
    // Returns a description of I/O, doesn't execute yet
    return func(ctx context.Context) (User, error) {
        // I/O happens here when the function is called
    }
}
```

**Key difference**: The functional version is a **curried function** where:
1. Business parameters come first: `fetchUser(id)`
2. Context comes last: `fetchUser(id)(ctx)`
3. The intermediate result is composable: `ReaderIOResult[User]`

## Benefits of the Functional Approach

### 1. Separation of Pure and Impure Code

```go
// Pure computation - no I/O, no context needed
func validateAge(user User) (User, error) {
    if user.Age < 18 {
        return User{}, errors.New("user too young")
    }
    return user, nil
}

// Impure I/O operation - needs context
func fetchUser(id int) ReaderIOResult[User] {
    return func(ctx context.Context) (User, error) {
        return db.QueryUser(ctx, id)
    }
}

// Compose them - pure logic lifted into ReaderIOResult
pipeline := F.Pipe2(
    fetchUser(42),                           // ReaderIOResult[User]
    readerioresult.ChainEitherK(validateAge), // Lift pure function
)

// Execute when ready
user, err := pipeline(ctx)
```

**Benefits:**
- Pure functions are easier to test (no mocking needed)
- Pure functions are easier to reason about (no side effects)
- Clear boundary between logic and I/O
- Can test business logic independently

### 2. Composability

Functions compose naturally without manual error checking:

```go
// Traditional approach - manual error handling
func ProcessUserTraditional(ctx context.Context, userID int) (UserWithPosts, error) {
    user, err := fetchUser(ctx, userID)
    if err != nil {
        return UserWithPosts{}, err
    }
    
    validated, err := validateUser(user)
    if err != nil {
        return UserWithPosts{}, err
    }
    
    posts, err := fetchPosts(ctx, validated.ID)
    if err != nil {
        return UserWithPosts{}, err
    }
    
    return enrichUser(validated, posts), nil
}

// Functional approach - automatic error propagation
func ProcessUserFunctional(userID int) ReaderIOResult[UserWithPosts] {
    return F.Pipe3(
        fetchUser(userID),
        readerioresult.ChainEitherK(validateUser),
        readerioresult.Chain(func(user User) ReaderIOResult[UserWithPosts] {
            return F.Pipe2(
                fetchPosts(user.ID),
                readerioresult.Map(func(posts []Post) UserWithPosts {
                    return enrichUser(user, posts)
                }),
            )
        }),
    )
}
```

**Benefits:**
- No manual error checking
- Automatic short-circuiting on first error
- Clear data flow
- Easier to refactor and extend

### 3. Testability

```go
// Mock I/O operations by providing test implementations
func TestProcessUser(t *testing.T) {
    // Create a mock that returns test data
    mockFetchUser := func(id int) ReaderIOResult[User] {
        return func(ctx context.Context) (User, error) {
            return User{ID: id, Age: 25}, nil
        }
    }
    
    // Test with mock - no database needed
    result, err := mockFetchUser(42)(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, 25, result.Age)
}
```

### 4. Lazy Evaluation

Operations are not executed until you provide the context:

```go
// Build the pipeline - no I/O happens yet
pipeline := F.Pipe3(
    fetchUser(42),
    readerioresult.Map(enrichUser),
    readerioresult.Chain(saveUser),
)

// I/O only happens when we call it with a context
user, err := pipeline(ctx)
```

**Benefits:**
- Build complex operations as pure data structures
- Defer execution until needed
- Reuse pipelines with different contexts
- Test pipelines without executing I/O

### 5. Context Propagation

Context is automatically threaded through all operations:

```go
// Traditional - must pass context explicitly everywhere
func Process(ctx context.Context) error {
    user, err := fetchUser(ctx, 42)
    if err != nil {
        return err
    }
    posts, err := fetchPosts(ctx, user.ID)
    if err != nil {
        return err
    }
    return savePosts(ctx, posts)
}

// Functional - context provided once at execution
func Process() ReaderIOResult[any] {
    return F.Pipe2(
        fetchUser(42),
        readerioresult.Chain(func(user User) ReaderIOResult[any] {
            return F.Pipe2(
                fetchPosts(user.ID),
                readerioresult.Chain(savePosts),
            )
        }),
    )
}

// Context provided once
err := readerioresult.Fold(
    func(err error) error { return err },
    func(any) error { return nil },
)(Process())(ctx)
```

## Side-by-Side Comparison

### Example: User Service with Database Operations

#### Traditional Go Style

```go
package traditional

import (
    "context"
    "database/sql"
    "fmt"
)

type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

type UserService struct {
    db *sql.DB
}

// Fetch user from database
func (s *UserService) GetUser(ctx context.Context, id int) (User, error) {
    var user User
    
    // Check context
    if err := ctx.Err(); err != nil {
        return User{}, err
    }
    
    // Query database
    row := s.db.QueryRowContext(ctx, 
        "SELECT id, name, email, age FROM users WHERE id = ?", id)
    
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    if err != nil {
        return User{}, fmt.Errorf("scan user: %w", err)
    }
    
    return user, nil
}

// Validate user
func (s *UserService) ValidateUser(ctx context.Context, user User) (User, error) {
    if user.Age < 18 {
        return User{}, fmt.Errorf("user %d is too young", user.ID)
    }
    if user.Email == "" {
        return User{}, fmt.Errorf("user %d has no email", user.ID)
    }
    return user, nil
}

// Update user email
func (s *UserService) UpdateEmail(ctx context.Context, id int, email string) (User, error) {
    // Check context
    if err := ctx.Err(); err != nil {
        return User{}, err
    }
    
    // Update database
    _, err := s.db.ExecContext(ctx,
        "UPDATE users SET email = ? WHERE id = ?", email, id)
    if err != nil {
        return User{}, fmt.Errorf("update email: %w", err)
    }
    
    // Fetch updated user
    return s.GetUser(ctx, id)
}

// Process user: fetch, validate, update email
func (s *UserService) ProcessUser(ctx context.Context, id int, newEmail string) (User, error) {
    // Fetch user
    user, err := s.GetUser(ctx, id)
    if err != nil {
        return User{}, fmt.Errorf("get user: %w", err)
    }
    
    // Validate user
    validated, err := s.ValidateUser(ctx, user)
    if err != nil {
        return User{}, fmt.Errorf("validate user: %w", err)
    }
    
    // Update email
    updated, err := s.UpdateEmail(ctx, validated.ID, newEmail)
    if err != nil {
        return User{}, fmt.Errorf("update email: %w", err)
    }
    
    return updated, nil
}
```

**Characteristics:**
- ✗ Manual error checking at every step
- ✗ Context passed explicitly to every function
- ✗ Error wrapping is manual and verbose
- ✗ Business logic mixed with error handling
- ✗ Hard to test without database
- ✗ Difficult to compose operations
- ✓ Familiar to Go developers
- ✓ Explicit control flow

#### Functional Go Style (context/readerioresult)

```go
package functional

import (
    "context"
    "database/sql"
    "fmt"
    
    F "github.com/IBM/fp-go/v2/function"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
)

type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

type UserService struct {
    db *sql.DB
}

// Fetch user from database - returns a ReaderIOResult
func (s *UserService) GetUser(id int) RIO.ReaderIOResult[User] {
    return func(ctx context.Context) (User, error) {
        var user User
        row := s.db.QueryRowContext(ctx,
            "SELECT id, name, email, age FROM users WHERE id = ?", id)
        
        err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
        if err != nil {
            return User{}, fmt.Errorf("scan user: %w", err)
        }
        
        return user, nil
    }
}

// Validate user - pure function (no I/O, no context)
func ValidateUser(user User) (User, error) {
    if user.Age < 18 {
        return User{}, fmt.Errorf("user %d is too young", user.ID)
    }
    if user.Email == "" {
        return User{}, fmt.Errorf("user %d has no email", user.ID)
    }
    return user, nil
}

// Update user email - returns a ReaderIOResult
func (s *UserService) UpdateEmail(id int, email string) RIO.ReaderIOResult[User] {
    return func(ctx context.Context) (User, error) {
        _, err := s.db.ExecContext(ctx,
            "UPDATE users SET email = ? WHERE id = ?", email, id)
        if err != nil {
            return User{}, fmt.Errorf("update email: %w", err)
        }
        
        // Chain to GetUser
        return s.GetUser(id)(ctx)
    }
}

// Process user: fetch, validate, update email - composable pipeline
func (s *UserService) ProcessUser(id int, newEmail string) RIO.ReaderIOResult[User] {
    return F.Pipe3(
        s.GetUser(id),                          // Fetch user
        RIO.ChainEitherK(ValidateUser),         // Validate (pure function)
        RIO.Chain(func(user User) RIO.ReaderIOResult[User] {
            return s.UpdateEmail(user.ID, newEmail)  // Update email
        }),
    )
}

// Alternative: Using Do-notation for more complex flows
func (s *UserService) ProcessUserDo(id int, newEmail string) RIO.ReaderIOResult[User] {
    return RIO.Chain(func(user User) RIO.ReaderIOResult[User] {
        // Validate is pure, lift it into ReaderIOResult
        validated, err := ValidateUser(user)
        if err != nil {
            return RIO.Left[User](err)
        }
        // Update with validated user
        return s.UpdateEmail(validated.ID, newEmail)
    })(s.GetUser(id))
}
```

**Characteristics:**
- ✓ Automatic error propagation (no manual checking)
- ✓ Context threaded automatically
- ✓ Pure functions separated from I/O
- ✓ Business logic clear and composable
- ✓ Easy to test (mock ReaderIOResult)
- ✓ Operations compose naturally
- ✓ Lazy evaluation (build pipeline, execute later)
- ✗ Requires understanding of functional patterns
- ✗ Less familiar to traditional Go developers

#### Idiomatic Functional Style (idiomatic/context/readerresult)

For even better performance with the same functional benefits:

```go
package idiomatic

import (
    "context"
    "database/sql"
    "fmt"
    
    F "github.com/IBM/fp-go/v2/function"
    RR "github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
)

type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

type UserService struct {
    db *sql.DB
}

// ReaderResult is just: func(context.Context) (A, error)
// Same as ReaderIOResult but using native Go tuples

func (s *UserService) GetUser(id int) RR.ReaderResult[User] {
    return func(ctx context.Context) (User, error) {
        var user User
        row := s.db.QueryRowContext(ctx,
            "SELECT id, name, email, age FROM users WHERE id = ?", id)
        
        err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
        return user, err  // Native tuple return
    }
}

// Pure validation - returns native (User, error) tuple
func ValidateUser(user User) (User, error) {
    if user.Age < 18 {
        return User{}, fmt.Errorf("user %d is too young", user.ID)
    }
    if user.Email == "" {
        return User{}, fmt.Errorf("user %d has no email", user.ID)
    }
    return user, nil
}

func (s *UserService) UpdateEmail(id int, email string) RR.ReaderResult[User] {
    return func(ctx context.Context) (User, error) {
        _, err := s.db.ExecContext(ctx,
            "UPDATE users SET email = ? WHERE id = ?", email, id)
        if err != nil {
            return User{}, err
        }
        return s.GetUser(id)(ctx)
    }
}

// Composable pipeline with native tuples
func (s *UserService) ProcessUser(id int, newEmail string) RR.ReaderResult[User] {
    return F.Pipe3(
        s.GetUser(id),
        RR.ChainEitherK(ValidateUser),  // Lift pure function
        RR.Chain(func(user User) RR.ReaderResult[User] {
            return s.UpdateEmail(user.ID, newEmail)
        }),
    )
}
```

**Characteristics:**
- ✓ All benefits of functional approach
- ✓ **2-10x better performance** (native tuples)
- ✓ **Zero allocations** for many operations
- ✓ More familiar to Go developers (uses (value, error))
- ✓ Seamless integration with existing Go code
- ✓ Same composability as ReaderIOResult

### Usage Comparison

```go
// Traditional
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    service := &UserService{db: db}
    
    user, err := service.ProcessUser(ctx, 42, "new@email.com")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}

// Functional (both styles)
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    service := &UserService{db: db}
    
    // Build the pipeline (no execution yet)
    pipeline := service.ProcessUser(42, "new@email.com")
    
    // Execute with context
    user, err := pipeline(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}

// Or using Fold for cleaner error handling
func HandleRequestFold(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    service := &UserService{db: db}
    
    RR.Fold(
        func(err error) {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        },
        func(user User) {
            json.NewEncoder(w).Encode(user)
        },
    )(service.ProcessUser(42, "new@email.com"))(ctx)
}
```

## Advanced Patterns

### Resource Management with Bracket

```go
// Traditional
func ProcessFile(ctx context.Context, path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer file.Close()
    
    data, err := io.ReadAll(file)
    if err != nil {
        return "", err
    }
    
    return string(data), nil
}

// Functional - guaranteed cleanup even on panic
func ProcessFile(path string) RIO.ReaderIOResult[string] {
    return RIO.Bracket(
        // Acquire resource
        func(ctx context.Context) (*os.File, error) {
            return os.Open(path)
        },
        // Release resource (always called)
        func(file *os.File, err error) RIO.ReaderIOResult[any] {
            return func(ctx context.Context) (any, error) {
                return nil, file.Close()
            }
        },
        // Use resource
        func(file *os.File) RIO.ReaderIOResult[string] {
            return func(ctx context.Context) (string, error) {
                data, err := io.ReadAll(file)
                return string(data), err
            }
        },
    )
}
```

### Parallel Execution

```go
// Traditional - manual goroutines and sync
func FetchMultipleUsers(ctx context.Context, ids []int) ([]User, error) {
    var wg sync.WaitGroup
    users := make([]User, len(ids))
    errs := make([]error, len(ids))
    
    for i, id := range ids {
        wg.Add(1)
        go func(i, id int) {
            defer wg.Done()
            users[i], errs[i] = fetchUser(ctx, id)
        }(i, id)
    }
    
    wg.Wait()
    
    for _, err := range errs {
        if err != nil {
            return nil, err
        }
    }
    
    return users, nil
}

// Functional - automatic parallelization
func FetchMultipleUsers(ids []int) RIO.ReaderIOResult[[]User] {
    operations := A.Map(func(id int) RIO.ReaderIOResult[User] {
        return fetchUser(id)
    })(ids)
    
    return RIO.TraverseArrayPar(F.Identity[RIO.ReaderIOResult[User]])(operations)
}
```

### Retry Logic

```go
// Traditional
func FetchWithRetry(ctx context.Context, url string, maxRetries int) ([]byte, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        if ctx.Err() != nil {
            return nil, ctx.Err()
        }
        
        resp, err := http.Get(url)
        if err == nil {
            defer resp.Body.Close()
            return io.ReadAll(resp.Body)
        }
        
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return nil, lastErr
}

// Functional
func FetchWithRetry(url string, maxRetries int) RIO.ReaderIOResult[[]byte] {
    operation := func(ctx context.Context) ([]byte, error) {
        resp, err := http.Get(url)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        return io.ReadAll(resp.Body)
    }
    
    return RIO.Retry(
        maxRetries,
        func(attempt int) time.Duration {
            return time.Second * time.Duration(attempt)
        },
    )(operation)
}
```

## When to Use Each Approach

### Use Traditional Go Style When:

1. **Team familiarity**: Team is not familiar with functional programming
2. **Simple operations**: Single I/O operation with straightforward error handling
3. **Existing codebase**: Large codebase already using traditional patterns
4. **Learning curve**: Want to minimize onboarding time
5. **Explicit control**: Need very explicit control flow

### Use Functional Style (ReaderIOResult) When:

1. **Complex pipelines**: Multiple I/O operations that need composition
2. **Testability**: Need to test business logic separately from I/O
3. **Reusability**: Want to build reusable operation pipelines
4. **Error handling**: Want automatic error propagation
5. **Resource management**: Need guaranteed cleanup (Bracket)
6. **Parallel execution**: Need to parallelize operations easily
7. **Type safety**: Want the type system to track I/O effects

### Use Idiomatic Functional Style (idiomatic/context/readerresult) When:

1. **All functional benefits**: Want functional patterns with Go idioms
2. **Performance critical**: Need 2-10x better performance
3. **Zero allocations**: Memory efficiency is important
4. **Go integration**: Want seamless integration with existing Go code
5. **Production services**: Building high-throughput services
6. **Best of both worlds**: Want functional composition with Go's native patterns

## Summary

The functional approach to I/O in Go offers significant advantages:

1. **Separation of Concerns**: Pure logic separated from I/O effects
2. **Composability**: Operations compose naturally without manual error checking
3. **Testability**: Easy to test without mocking I/O
4. **Type Safety**: I/O effects visible in the type system
5. **Lazy Evaluation**: Build pipelines, execute when ready
6. **Context Propagation**: Automatic threading of context
7. **Performance**: Idiomatic version offers 2-10x speedup

The key insight is that **I/O operations return descriptions of effects** (ReaderIOResult) rather than performing effects immediately. This enables powerful composition patterns while maintaining Go's idiomatic error handling through the `(value, error)` tuple pattern.

For production Go services, the **idiomatic/context/readerresult** package provides the best balance: full functional programming capabilities with native Go performance and familiar error handling patterns.

## Further Reading

- [DESIGN.md](./DESIGN.md) - Design principles and patterns
- [IDIOMATIC_COMPARISON.md](./IDIOMATIC_COMPARISON.md) - Performance comparison
- [idiomatic/doc.go](./idiomatic/doc.go) - Idiomatic package overview
- [context/readerioresult](./context/readerioresult/) - ReaderIOResult package
- [idiomatic/context/readerresult](./idiomatic/context/readerresult/) - Idiomatic ReaderResult package