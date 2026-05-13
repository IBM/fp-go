---
sidebar_position: 10
title: HTTP Requests
description: Functional HTTP client patterns
hide_title: true
---

<PageHeader
  eyebrow="Recipes · 10 / 17"
  title="HTTP"
  titleAccent="Requests"
  lede="Make HTTP requests using functional patterns with IOEither for lazy evaluation, composability, and type-safe error handling."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '6' },
    { label: 'Use Cases', value: 'APIs, Microservices, Integration' }
  ]}
/>

<TLDR>
  <TLDRCard title="Lazy Evaluation" icon="clock">
    Requests don't execute until called—compose operations without triggering side effects.
  </TLDRCard>
  <TLDRCard title="Composable Requests" icon="layers">
    Chain multiple requests with proper error handling—build complex workflows from simple operations.
  </TLDRCard>
  <TLDRCard title="Built-in Retry" icon="refresh-cw">
    Add retry logic with exponential backoff and circuit breakers—handle transient failures gracefully.
  </TLDRCard>
</TLDR>

<Section id="basic-requests" number="01" title="Basic HTTP" titleAccent="Requests">

Simple GET and POST requests with proper error handling.

<CodeCard file="simple-get.go">
{`package main

import (
    "fmt"
    "io"
    "net/http"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func fetchURL(url string) IOE.IOEither[error, string] {
    return IOE.TryCatch(func() (string, error) {
        resp, err := http.Get(url)
        if err != nil {
            return "", err
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
        }
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return "", err
        }
        
        return string(body), nil
    })
}

func main() {
    result := fetchURL("https://api.github.com/users/octocat")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Println("Response:", result.Right()[:100], "...")
    }
}`}
</CodeCard>

<CodeCard file="post-json.go">
{`package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type CreateUserRequest struct {
    Name  string \`json:"name"\`
    Email string \`json:"email"\`
}

type CreateUserResponse struct {
    ID    int    \`json:"id"\`
    Name  string \`json:"name"\`
    Email string \`json:"email"\`
}

func postJSON[Req, Resp any](url string, data Req) IOE.IOEither[error, Resp] {
    return IOE.TryCatch(func() (Resp, error) {
        var result Resp
        
        jsonData, err := json.Marshal(data)
        if err != nil {
            return result, fmt.Errorf("marshal error: %w", err)
        }
        
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            return result, fmt.Errorf("request error: %w", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
            body, _ := io.ReadAll(resp.Body)
            return result, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
        }
        
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            return result, fmt.Errorf("decode error: %w", err)
        }
        
        return result, nil
    })
}

func main() {
    request := CreateUserRequest{
        Name:  "Alice",
        Email: "alice@example.com",
    }
    
    result := postJSON[CreateUserRequest, CreateUserResponse](
        "https://api.example.com/users",
        request,
    )()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        user := result.Right()
        fmt.Printf("Created user: ID=%d, Name=%s\\n", user.ID, user.Name)
    }
}`}
</CodeCard>

</Section>

<Section id="request-builder" number="02" title="Fluent Request" titleAccent="Builder">

Build complex requests with headers, query parameters, and custom configuration.

<CodeCard file="request-builder.go">
{`package main

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type RequestBuilder struct {
    method  string
    url     string
    headers map[string]string
    query   url.Values
    body    io.Reader
}

func NewRequest(method, urlStr string) *RequestBuilder {
    return &RequestBuilder{
        method:  method,
        url:     urlStr,
        headers: make(map[string]string),
        query:   url.Values{},
    }
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
    rb.headers[key] = value
    return rb
}

func (rb *RequestBuilder) WithQuery(key, value string) *RequestBuilder {
    rb.query.Add(key, value)
    return rb
}

func (rb *RequestBuilder) Execute() IOE.IOEither[error, *http.Response] {
    return IOE.TryCatch(func() (*http.Response, error) {
        u, err := url.Parse(rb.url)
        if err != nil {
            return nil, err
        }
        u.RawQuery = rb.query.Encode()
        
        req, err := http.NewRequest(rb.method, u.String(), rb.body)
        if err != nil {
            return nil, err
        }
        
        for key, value := range rb.headers {
            req.Header.Set(key, value)
        }
        
        return http.DefaultClient.Do(req)
    })
}

func main() {
    result := F.Pipe2(
        NewRequest("GET", "https://api.github.com/search/repositories").
            WithQuery("q", "language:go").
            WithQuery("sort", "stars").
            WithHeader("Accept", "application/vnd.github.v3+json").
            Execute(),
        IOE.ChainFirst(func(resp *http.Response) IOE.IOEither[error, *http.Response] {
            if resp.StatusCode != http.StatusOK {
                return IOE.Left[*http.Response](
                    fmt.Errorf("HTTP %d", resp.StatusCode),
                )
            }
            return IOE.Right[error](resp)
        }),
        IOE.Map(func(resp *http.Response) string {
            defer resp.Body.Close()
            body, _ := io.ReadAll(resp.Body)
            return string(body)
        }),
    )()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Println("Response:", result.Right()[:200], "...")
    }
}`}
</CodeCard>

</Section>

<Section id="retry-backoff" number="03" title="Retry with Exponential" titleAccent="Backoff">

Handle transient failures with retry logic and exponential backoff.

<CodeCard file="retry-backoff.go">
{`package main

import (
    "fmt"
    "time"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func retryWithBackoff[A any](
    maxRetries int,
    initialDelay time.Duration,
) func(IOE.IOEither[error, A]) IOE.IOEither[error, A] {
    return func(io IOE.IOEither[error, A]) IOE.IOEither[error, A] {
        return func() IOE.Either[error, A] {
            var lastErr error
            delay := initialDelay
            
            for i := 0; i <= maxRetries; i++ {
                result := io()
                
                if result.IsRight() {
                    return result
                }
                
                lastErr = result.Left()
                
                if i < maxRetries {
                    fmt.Printf("Attempt %d failed, retrying in %v...\\n", i+1, delay)
                    time.Sleep(delay)
                    delay *= 2
                }
            }
            
            return IOE.Left[A](fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr))()
        }
    }
}

func fetchWithRetry(url string) IOE.IOEither[error, string] {
    return F.Pipe1(
        fetchURL(url),
        retryWithBackoff[string](3, 1*time.Second),
    )
}

func main() {
    result := fetchWithRetry("https://api.example.com/unstable-endpoint")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Println("Success:", result.Right())
    }
}`}
</CodeCard>

</Section>

<Section id="circuit-breaker" number="04" title="Circuit Breaker" titleAccent="Pattern">

Prevent cascading failures with circuit breaker pattern.

<CodeCard file="circuit-breaker.go">
{`package main

import (
    "fmt"
    "sync"
    "time"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    failures     int
    lastFailTime time.Time
    state        string // "closed", "open", "half-open"
    mu           sync.Mutex
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
        state:        "closed",
    }
}

func (cb *CircuitBreaker) Execute[A any](
    io IOE.IOEither[error, A],
) IOE.IOEither[error, A] {
    return func() IOE.Either[error, A] {
        cb.mu.Lock()
        
        if cb.state == "open" && time.Since(cb.lastFailTime) > cb.resetTimeout {
            cb.state = "half-open"
            cb.failures = 0
        }
        
        if cb.state == "open" {
            cb.mu.Unlock()
            return IOE.Left[A](fmt.Errorf("circuit breaker is open"))()
        }
        
        cb.mu.Unlock()
        
        result := io()
        
        cb.mu.Lock()
        defer cb.mu.Unlock()
        
        if result.IsLeft() {
            cb.failures++
            cb.lastFailTime = time.Now()
            
            if cb.failures >= cb.maxFailures {
                cb.state = "open"
                fmt.Println("Circuit breaker opened")
            }
        } else {
            if cb.state == "half-open" {
                cb.state = "closed"
                fmt.Println("Circuit breaker closed")
            }
            cb.failures = 0
        }
        
        return result
    }
}

func main() {
    cb := NewCircuitBreaker(3, 5*time.Second)
    
    for i := 0; i < 10; i++ {
        result := cb.Execute(fetchURL("https://api.example.com/failing"))()
        
        if result.IsLeft() {
            fmt.Printf("Request %d failed: %v\\n", i+1, result.Left())
        } else {
            fmt.Printf("Request %d succeeded\\n", i+1)
        }
        
        time.Sleep(500 * time.Millisecond)
    }
}`}
</CodeCard>

</Section>

<Section id="parallel-requests" number="05" title="Parallel" titleAccent="Requests">

Fetch multiple URLs concurrently for improved performance.

<CodeCard file="parallel-fetch.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func fetchAll(urls []string) IOE.IOEither[error, []string] {
    return A.Traverse[string](IOE.ApplicativePar[error, string]())(
        fetchURL,
    )(urls)
}

func main() {
    urls := []string{
        "https://api.github.com/users/octocat",
        "https://api.github.com/users/torvalds",
        "https://api.github.com/users/gvanrossum",
    }
    
    result := fetchAll(urls)()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        responses := result.Right()
        fmt.Printf("Fetched %d responses\\n", len(responses))
        for i, resp := range responses {
            fmt.Printf("Response %d: %d bytes\\n", i+1, len(resp))
        }
    }
}`}
</CodeCard>

</Section>

<Section id="sequential-composition" number="06" title="Sequential Request" titleAccent="Composition">

Chain dependent requests where later requests depend on earlier results.

<CodeCard file="sequential-requests.go">
{`package main

import (
    "encoding/json"
    "fmt"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type User struct {
    ID       int    \`json:"id"\`
    Username string \`json:"login"\`
}

type Repository struct {
    Name        string \`json:"name"\`
    Description string \`json:"description"\`
}

func getUser(username string) IOE.IOEither[error, User] {
    return F.Pipe2(
        fetchURL(fmt.Sprintf("https://api.github.com/users/%s", username)),
        IOE.Chain(func(body string) IOE.IOEither[error, User] {
            var user User
            if err := json.Unmarshal([]byte(body), &user); err != nil {
                return IOE.Left[User](err)
            }
            return IOE.Right[error](user)
        }),
    )
}

func getUserRepos(username string) IOE.IOEither[error, []Repository] {
    return F.Pipe2(
        fetchURL(fmt.Sprintf("https://api.github.com/users/%s/repos", username)),
        IOE.Chain(func(body string) IOE.IOEither[error, []Repository] {
            var repos []Repository
            if err := json.Unmarshal([]byte(body), &repos); err != nil {
                return IOE.Left[[]Repository](err)
            }
            return IOE.Right[error](repos)
        }),
    )
}

func getUserWithRepos(username string) IOE.IOEither[error, struct {
    User  User
    Repos []Repository
}] {
    return F.Pipe3(
        IOE.Do[error](IOE.Monad[error, struct {
            User  User
            Repos []Repository
        }]()),
        IOE.Bind("user", func() IOE.IOEither[error, User] {
            return getUser(username)
        }),
        IOE.Bind("repos", func() IOE.IOEither[error, []Repository] {
            return getUserRepos(username)
        }),
        IOE.Map(func(data struct {
            user  User
            repos []Repository
        }) struct {
            User  User
            Repos []Repository
        } {
            return struct {
                User  User
                Repos []Repository
            }{
                User:  data.user,
                Repos: data.repos,
            }
        }),
    )
}

func main() {
    result := getUserWithRepos("octocat")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        data := result.Right()
        fmt.Printf("User: %s (ID: %d)\\n", data.User.Username, data.User.ID)
        fmt.Printf("Repositories: %d\\n", len(data.Repos))
    }
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="07" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Use context for cancellation** — Support request cancellation with context.Context
  </ChecklistItem>
  <ChecklistItem status="required">
    **Set timeouts** — Always configure client timeouts to prevent hanging
  </ChecklistItem>
  <ChecklistItem status="required">
    **Handle rate limiting** — Implement rate limiting to respect API quotas
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Add retry logic** — Use exponential backoff for transient failures
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Implement circuit breakers** — Prevent cascading failures in distributed systems
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Log requests** — Track request/response for debugging and monitoring
  </ChecklistItem>
</Checklist>

</Section>
