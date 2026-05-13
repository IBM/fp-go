---
title: Middleware Patterns
hide_title: true
description: Build composable middleware using functional patterns with fp-go for cross-cutting concerns like logging, auth, and caching.
sidebar_position: 15
---

<PageHeader
  eyebrow="Recipes · 15 / 17"
  title="Middleware"
  titleAccent="Patterns"
  lede="Build composable middleware using functional patterns with fp-go for cross-cutting concerns like logging, authentication, caching, and error handling."
  meta={[
    { label: 'Difficulty', value: 'Advanced' },
    { label: 'Patterns', value: '6' },
    { label: 'Use Cases', value: 'HTTP handlers, request processing, cross-cutting concerns' }
  ]}
/>

<TLDR>
  <TLDRCard title="Function Composition" icon="layers">
    Chain middleware using function composition for clean, declarative request processing pipelines.
  </TLDRCard>
  <TLDRCard title="Type-Safe Wrapping" icon="shield">
    Wrap operations with compile-time guarantees that middleware preserves input/output types.
  </TLDRCard>
  <TLDRCard title="Reusable & Testable" icon="check">
    Build small, focused middleware pieces that are easy to test and compose into complex behaviors.
  </TLDRCard>
</TLDR>

<Section id="basic-middleware" number="01" title="Basic" titleAccent="Middleware">

Middleware wraps operations to add cross-cutting concerns without modifying the core logic.

<CodeCard file="logging_middleware.go">
{`package main

import (
    "fmt"
    "time"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type Middleware[A any] func(IOE.IOEither[error, A]) IOE.IOEither[error, A]

func withLogging[A any](name string) Middleware[A] {
    return func(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        return func() IOE.Either[error, A] {
            fmt.Printf("[%s] Starting...\\n", name)
            start := time.Now()
            
            result := operation()
            
            duration := time.Since(start)
            if result.IsLeft() {
                fmt.Printf("[%s] Failed after %v: %v\\n", name, duration, result.Left())
            } else {
                fmt.Printf("[%s] Completed in %v\\n", name, duration)
            }
            
            return result
        }
    }
}

func fetchData() IOE.IOEither[error, string] {
    return IOE.TryCatch(func() (string, error) {
        time.Sleep(100 * time.Millisecond)
        return "data", nil
    })
}

func main() {
    // Apply logging middleware
    operation := withLogging[string]("fetchData")(fetchData())
    
    result := operation()
    fmt.Println("Result:", result.Right())
    // [fetchData] Starting...
    // [fetchData] Completed in 100ms
    // Result: data
}
`}
</CodeCard>

<Callout type="info">
**Middleware Pattern**: A middleware is a higher-order function that takes an operation and returns a wrapped version with additional behavior.
</Callout>

</Section>

<Section id="http-middleware" number="02" title="HTTP" titleAccent="Middleware">

Build HTTP middleware for authentication, rate limiting, and request processing.

<CodeCard file="auth_middleware.go">
{`package main

import (
    "context"
    "fmt"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type Request struct {
    Context context.Context
    Headers map[string]string
    Body    []byte
}

type Response struct {
    Status int
    Body   []byte
}

type Handler func(Request) IOE.IOEither[error, Response]

func withAuth(handler Handler) Handler {
    return func(req Request) IOE.IOEither[error, Response] {
        token := req.Headers["Authorization"]
        
        if token == "" {
            return IOE.Left[Response](fmt.Errorf("missing authorization"))
        }
        
        if !isValidToken(token) {
            return IOE.Left[Response](fmt.Errorf("invalid token"))
        }
        
        // Add user to context
        ctx := context.WithValue(req.Context, "user", getUserFromToken(token))
        req.Context = ctx
        
        return handler(req)
    }
}

func isValidToken(token string) bool {
    return token == "valid-token"
}

func getUserFromToken(token string) string {
    return "user-123"
}

func handleRequest(req Request) IOE.IOEither[error, Response] {
    return IOE.Right[error](Response{
        Status: 200,
        Body:   []byte("Success"),
    })
}

func main() {
    handler := withAuth(handleRequest)
    
    // Valid request
    req1 := Request{
        Context: context.Background(),
        Headers: map[string]string{"Authorization": "valid-token"},
    }
    result1 := handler(req1)()
    fmt.Println("Valid:", result1.IsRight())
    
    // Invalid request
    req2 := Request{
        Context: context.Background(),
        Headers: map[string]string{},
    }
    result2 := handler(req2)()
    fmt.Println("Invalid:", result2.Left())
}
`}
</CodeCard>

</Section>

<Section id="composing-middleware" number="03" title="Composing" titleAccent="Middleware">

Chain multiple middleware together for complex request processing pipelines.

<CodeCard file="middleware_chain.go">
{`package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func compose[A any](middlewares ...Middleware[A]) Middleware[A] {
    return func(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        result := operation
        // Apply middlewares in reverse order (right to left)
        for i := len(middlewares) - 1; i >= 0; i-- {
            result = middlewares[i](result)
        }
        return result
    }
}

func withRetry[A any](maxAttempts int) Middleware[A] {
    return func(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        return func() IOE.Either[error, A] {
            var lastErr error
            for i := 0; i < maxAttempts; i++ {
                result := operation()
                if result.IsRight() {
                    return result
                }
                lastErr = result.Left()
                fmt.Printf("Attempt %d failed, retrying...\\n", i+1)
            }
            return IOE.Left[A](fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr))()
        }
    }
}

func withCache[A any](cache map[string]A, key string) Middleware[A] {
    return func(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        return func() IOE.Either[error, A] {
            if cached, ok := cache[key]; ok {
                fmt.Println("Cache hit!")
                return IOE.Right[error](cached)()
            }
            
            result := operation()
            if result.IsRight() {
                cache[key] = result.Right()
            }
            return result
        }
    }
}

func main() {
    cache := make(map[string]string)
    
    // Compose multiple middleware
    middleware := compose(
        withLogging[string]("operation"),
        withRetry[string](3),
        withCache(cache, "data"),
    )
    
    operation := middleware(fetchData())
    
    // First call: cache miss, fetches data
    result1 := operation()
    fmt.Println("First:", result1.Right())
    
    // Second call: cache hit
    result2 := operation()
    fmt.Println("Second:", result2.Right())
}
`}
</CodeCard>

<CodeCard file="pipeline_middleware.go">
{`package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type Pipeline[A any] struct {
    middlewares []Middleware[A]
}

func NewPipeline[A any]() *Pipeline[A] {
    return &Pipeline[A]{
        middlewares: []Middleware[A]{},
    }
}

func (p *Pipeline[A]) Use(middleware Middleware[A]) *Pipeline[A] {
    p.middlewares = append(p.middlewares, middleware)
    return p
}

func (p *Pipeline[A]) Execute(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
    return compose(p.middlewares...)(operation)
}

func main() {
    pipeline := NewPipeline[string]().
        Use(withLogging[string]("step1")).
        Use(withRetry[string](2)).
        Use(withTiming[string])
    
    result := pipeline.Execute(fetchData())()
    fmt.Println("Result:", result.Right())
}
`}
</CodeCard>

</Section>

<Section id="reader-middleware" number="04" title="Reader-Based" titleAccent="Middleware">

Use Reader pattern for middleware that needs dependencies.

<CodeCard file="reader_middleware.go">
{`package main

import (
    "context"
    "fmt"
    RIE "github.com/IBM/fp-go/v2/readerioeither"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type Dependencies struct {
    Logger Logger
    Cache  Cache
    DB     Database
}

type Logger interface {
    Info(msg string)
    Error(msg string)
}

type Cache interface {
    Get(key string) (string, bool)
    Set(key string, value string)
}

type Database interface {
    Query(sql string) ([]string, error)
}

type AppHandler[A any] = RIE.ReaderIOEither[Dependencies, error, A]

func withLogging[A any](name string, handler AppHandler[A]) AppHandler[A] {
    return F.Pipe2(
        RIE.Asks(func(deps Dependencies) IOE.IOEither[error, struct{}] {
            return IOE.TryCatch(func() (struct{}, error) {
                deps.Logger.Info(fmt.Sprintf("[%s] Starting", name))
                return struct{}{}, nil
            })
        }),
        RIE.Chain(func(_ struct{}) AppHandler[A] {
            return handler
        }),
    )
}

func withCaching[A any](key string, handler AppHandler[A]) AppHandler[A] {
    return RIE.Asks(func(deps Dependencies) IOE.IOEither[error, A] {
        // Check cache
        if cached, ok := deps.Cache.Get(key); ok {
            deps.Logger.Info("Cache hit")
            // Type assertion needed here
            return IOE.Right[error](cached.(A))
        }
        
        // Execute handler
        result := handler(deps)()
        
        // Store in cache if successful
        if result.IsRight() {
            deps.Cache.Set(key, fmt.Sprint(result.Right()))
        }
        
        return result
    })
}

func getUsers() AppHandler[[]string] {
    return RIE.Asks(func(deps Dependencies) IOE.IOEither[error, []string] {
        return IOE.TryCatch(func() ([]string, error) {
            return deps.DB.Query("SELECT * FROM users")
        })
    })
}

func main() {
    deps := Dependencies{
        Logger: &ConsoleLogger{},
        Cache:  &MemoryCache{},
        DB:     &MockDB{},
    }
    
    handler := withLogging("getUsers", withCaching("users", getUsers()))
    
    result := handler(deps)()
    fmt.Println("Users:", result.Right())
}
`}
</CodeCard>

</Section>

<Section id="error-middleware" number="05" title="Error Handling" titleAccent="Middleware">

Build middleware for error recovery and transformation.

<CodeCard file="error_middleware.go">
{`package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func withErrorRecovery[A any](fallback A) Middleware[A] {
    return func(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        return func() IOE.Either[error, A] {
            result := operation()
            if result.IsLeft() {
                fmt.Printf("Error occurred: %v, using fallback\\n", result.Left())
                return IOE.Right[error](fallback)()
            }
            return result
        }
    }
}

func withErrorMapping[A any](mapError func(error) error) Middleware[A] {
    return func(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        return func() IOE.Either[error, A] {
            result := operation()
            if result.IsLeft() {
                return IOE.Left[A](mapError(result.Left()))()
            }
            return result
        }
    }
}

func main() {
    // With recovery
    operation1 := withErrorRecovery("default")(
        IOE.Left[string](fmt.Errorf("failed")),
    )
    result1 := operation1()
    fmt.Println("Recovered:", result1.Right())
    
    // With error mapping
    operation2 := withErrorMapping[string](func(err error) error {
        return fmt.Errorf("wrapped: %w", err)
    })(IOE.Left[string](fmt.Errorf("original")))
    result2 := operation2()
    fmt.Println("Mapped:", result2.Left())
}
`}
</CodeCard>

</Section>

<Section id="best-practices" number="06" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Keep middleware focused** — Each middleware should have a single, clear responsibility
  </ChecklistItem>
  <ChecklistItem status="required">
    **Order matters** — Apply middleware in logical order (auth → cache → logging → retry)
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Make configurable** — Accept configuration parameters instead of hardcoding values
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Test independently** — Test each middleware separately before testing composition
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Document behavior** — Clearly document what each middleware does and its side effects
  </ChecklistItem>
</Checklist>

<Compare>
<CompareCol kind="good">
<CodeCard file="good_middleware.go">
{`// ✅ Good: Single responsibility
func withLogging[A any](name string) Middleware[A] { /* ... */ }
func withAuth[A any](token string) Middleware[A] { /* ... */ }
func withCache[A any](key string) Middleware[A] { /* ... */ }

// ✅ Good: Logical order
pipeline := NewPipeline[string]().
    Use(withAuth).        // Check auth first
    Use(withCache).       // Then check cache
    Use(withLogging).     // Log the actual operation
    Use(withRetry)        // Retry if needed

// ✅ Good: Configurable
func withRetry[A any](config RetryConfig) Middleware[A] {
    return func(operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        // Use config.MaxAttempts, config.Delay, etc.
    }
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="bad">
<CodeCard file="bad_middleware.go">
{`// ❌ Avoid: Doing too much
func withEverything[A any]() Middleware[A] {
    // Logging, auth, caching, retry, metrics...
}

// ❌ Avoid: Illogical order
pipeline := NewPipeline[string]().
    Use(withRetry).       // Retry before auth?
    Use(withCache).       // Cache before auth?
    Use(withAuth)

// ❌ Avoid: Hardcoded values
func withRetry[A any]() Middleware[A] {
    maxAttempts := 3 // Hardcoded
    // ...
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
