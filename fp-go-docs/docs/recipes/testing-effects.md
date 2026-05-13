---
title: Testing Effects
hide_title: true
description: Test IO operations and side effects with fp-go using mocking, dependency injection, and the Reader pattern.
sidebar_position: 17
---

<PageHeader
  eyebrow="Recipes · 17 / 17"
  title="Testing"
  titleAccent="Effects"
  lede="Test IO operations and side effects with fp-go using mocking, dependency injection, and the Reader pattern for deterministic, reliable tests."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '6' },
    { label: 'Use Cases', value: 'IO testing, mocking, integration tests' }
  ]}
/>

<TLDR>
  <TLDRCard title="Lazy Evaluation" icon="clock">
    Effects don't execute until called—build test scenarios without triggering side effects.
  </TLDRCard>
  <TLDRCard title="Dependency Injection" icon="plug">
    Use Reader pattern to inject mock dependencies for isolated, fast tests.
  </TLDRCard>
  <TLDRCard title="Interface-Based Mocking" icon="layers">
    Define interfaces for all external dependencies to enable easy mocking and testing.
  </TLDRCard>
</TLDR>

<Section id="ioeither-testing" number="01" title="Testing" titleAccent="IOEither">

Test IO operations by executing them and asserting on the result.

<CodeCard file="ioeither_test.go">
{`package main

import (
    "testing"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func readConfig() IOE.IOEither[error, string] {
    return IOE.TryCatch(func() (string, error) {
        // Simulated file read
        return "config-value", nil
    })
}

func TestReadConfig(t *testing.T) {
    // Execute the IO operation
    result := readConfig()()
    
    if result.IsLeft() {
        t.Errorf("Expected success, got error: %v", result.Left())
    }
    
    if result.Right() != "config-value" {
        t.Errorf("Expected 'config-value', got %q", result.Right())
    }
}

func fetchData() IOE.IOEither[error, string] {
    return IOE.Right[error]("data")
}

func processData(data string) IOE.IOEither[error, string] {
    return IOE.Right[error](fmt.Sprintf("processed: %s", data))
}

func pipeline() IOE.IOEither[error, string] {
    return F.Pipe1(
        fetchData(),
        IOE.Chain(processData),
    )
}

func TestPipeline(t *testing.T) {
    result := pipeline()()
    
    if result.IsLeft() {
        t.Fatalf("Expected success, got error: %v", result.Left())
    }
    
    expected := "processed: data"
    if result.Right() != expected {
        t.Errorf("Expected %q, got %q", expected, result.Right())
    }
}
`}
</CodeCard>

<Callout type="info">
**Lazy Evaluation**: IOEither operations don't execute until you call them with `()`. This lets you build and compose operations without side effects.
</Callout>

</Section>

<Section id="mocking-dependencies" number="02" title="Mocking" titleAccent="Dependencies">

Use interfaces to enable mocking of external dependencies.

<CodeCard file="mocking_test.go">
{`package main

import (
    "context"
    "fmt"
    "testing"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

// Interface for database operations
type Database interface {
    Query(ctx context.Context, sql string) IOE.IOEither[error, []string]
}

// Real implementation
type PostgresDB struct {
    connString string
}

func (db *PostgresDB) Query(ctx context.Context, sql string) IOE.IOEither[error, []string] {
    return IOE.TryCatch(func() ([]string, error) {
        // Real database query
        return []string{"result1", "result2"}, nil
    })
}

// Mock implementation for testing
type MockDB struct {
    queryFunc func(ctx context.Context, sql string) IOE.IOEither[error, []string]
    calls     []string
}

func (m *MockDB) Query(ctx context.Context, sql string) IOE.IOEither[error, []string] {
    m.calls = append(m.calls, sql)
    if m.queryFunc != nil {
        return m.queryFunc(ctx, sql)
    }
    return IOE.Right[error]([]string{})
}

// Function under test
func getUsers(db Database) IOE.IOEither[error, []string] {
    return db.Query(context.Background(), "SELECT * FROM users")
}

func TestGetUsers(t *testing.T) {
    t.Run("successful query", func(t *testing.T) {
        mockDB := &MockDB{
            queryFunc: func(ctx context.Context, sql string) IOE.IOEither[error, []string] {
                return IOE.Right[error]([]string{"user1", "user2"})
            },
        }
        
        result := getUsers(mockDB)()
        
        if result.IsLeft() {
            t.Fatalf("Expected success, got error: %v", result.Left())
        }
        
        users := result.Right()
        if len(users) != 2 {
            t.Errorf("Expected 2 users, got %d", len(users))
        }
        
        // Verify the query was called
        if len(mockDB.calls) != 1 {
            t.Errorf("Expected 1 query call, got %d", len(mockDB.calls))
        }
    })
    
    t.Run("query error", func(t *testing.T) {
        mockDB := &MockDB{
            queryFunc: func(ctx context.Context, sql string) IOE.IOEither[error, []string] {
                return IOE.Left[[]string](fmt.Errorf("connection failed"))
            },
        }
        
        result := getUsers(mockDB)()
        
        if result.IsRight() {
            t.Error("Expected error, got success")
        }
    })
}
`}
</CodeCard>

</Section>

<Section id="reader-pattern" number="03" title="Reader Pattern for" titleAccent="Testing">

Use Reader pattern to inject dependencies for testable code.

<CodeCard file="reader_test.go">
{`package main

import (
    "context"
    "testing"
    RIE "github.com/IBM/fp-go/v2/readerioeither"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type Dependencies struct {
    DB     Database
    Logger Logger
}

type Logger interface {
    Info(msg string)
    Error(msg string)
}

type AppEffect[A any] = RIE.ReaderIOEither[Dependencies, error, A]

func getUsersWithLogging() AppEffect[[]string] {
    return RIE.Asks(func(deps Dependencies) IOE.IOEither[error, []string] {
        deps.Logger.Info("Fetching users")
        return deps.DB.Query(context.Background(), "SELECT * FROM users")
    })
}

// Mock logger for testing
type MockLogger struct {
    infos  []string
    errors []string
}

func (m *MockLogger) Info(msg string) {
    m.infos = append(m.infos, msg)
}

func (m *MockLogger) Error(msg string) {
    m.errors = append(m.errors, msg)
}

func TestGetUsersWithLogging(t *testing.T) {
    mockDB := &MockDB{
        queryFunc: func(ctx context.Context, sql string) IOE.IOEither[error, []string] {
            return IOE.Right[error]([]string{"user1", "user2"})
        },
    }
    
    mockLogger := &MockLogger{}
    
    deps := Dependencies{
        DB:     mockDB,
        Logger: mockLogger,
    }
    
    result := getUsersWithLogging()(deps)()
    
    if result.IsLeft() {
        t.Fatalf("Expected success, got error: %v", result.Left())
    }
    
    // Verify logging
    if len(mockLogger.infos) != 1 {
        t.Errorf("Expected 1 info log, got %d", len(mockLogger.infos))
    }
    
    if mockLogger.infos[0] != "Fetching users" {
        t.Errorf("Expected 'Fetching users', got %q", mockLogger.infos[0])
    }
}
`}
</CodeCard>

</Section>

<Section id="file-http-testing" number="04" title="Testing File & HTTP" titleAccent="Operations">

Mock file system and HTTP clients for isolated tests.

<CodeCard file="filesystem_test.go">
{`package main

import (
    "fmt"
    "testing"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type FileSystem interface {
    ReadFile(path string) IOE.IOEither[error, []byte]
    WriteFile(path string, data []byte) IOE.IOEither[error, int]
}

type MockFS struct {
    files map[string][]byte
}

func NewMockFS() *MockFS {
    return &MockFS{
        files: make(map[string][]byte),
    }
}

func (m *MockFS) ReadFile(path string) IOE.IOEither[error, []byte] {
    return IOE.TryCatch(func() ([]byte, error) {
        if data, ok := m.files[path]; ok {
            return data, nil
        }
        return nil, fmt.Errorf("file not found: %s", path)
    })
}

func (m *MockFS) WriteFile(path string, data []byte) IOE.IOEither[error, int] {
    return IOE.TryCatch(func() (int, error) {
        m.files[path] = data
        return len(data), nil
    })
}

func copyFile(fs FileSystem, src, dst string) IOE.IOEither[error, int] {
    return F.Pipe2(
        fs.ReadFile(src),
        IOE.Chain(func(data []byte) IOE.IOEither[error, int] {
            return fs.WriteFile(dst, data)
        }),
    )
}

func TestCopyFile(t *testing.T) {
    fs := NewMockFS()
    fs.files["source.txt"] = []byte("test content")
    
    result := copyFile(fs, "source.txt", "dest.txt")()
    
    if result.IsLeft() {
        t.Fatalf("Expected success, got error: %v", result.Left())
    }
    
    // Verify file was copied
    if data, ok := fs.files["dest.txt"]; !ok {
        t.Error("Destination file not created")
    } else if string(data) != "test content" {
        t.Errorf("Expected 'test content', got %q", string(data))
    }
}
`}
</CodeCard>

<CodeCard file="http_test.go">
{`package main

import (
    "encoding/json"
    "fmt"
    "testing"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type HTTPClient interface {
    Get(url string) IOE.IOEither[error, []byte]
    Post(url string, data []byte) IOE.IOEither[error, []byte]
}

type MockHTTPClient struct {
    getFunc  func(url string) IOE.IOEither[error, []byte]
    postFunc func(url string, data []byte) IOE.IOEither[error, []byte]
    calls    []string
}

func (m *MockHTTPClient) Get(url string) IOE.IOEither[error, []byte] {
    m.calls = append(m.calls, "GET "+url)
    if m.getFunc != nil {
        return m.getFunc(url)
    }
    return IOE.Right[error]([]byte("{}"))
}

func (m *MockHTTPClient) Post(url string, data []byte) IOE.IOEither[error, []byte] {
    m.calls = append(m.calls, "POST "+url)
    if m.postFunc != nil {
        return m.postFunc(url, data)
    }
    return IOE.Right[error]([]byte("{}"))
}

type User struct {
    ID   int    \`json:"id"\`
    Name string \`json:"name"\`
}

func fetchUser(client HTTPClient, id int) IOE.IOEither[error, User] {
    return F.Pipe2(
        client.Get(fmt.Sprintf("https://api.example.com/users/%d", id)),
        IOE.Chain(func(data []byte) IOE.IOEither[error, User] {
            return IOE.TryCatch(func() (User, error) {
                var user User
                err := json.Unmarshal(data, &user)
                return user, err
            })
        }),
    )
}

func TestFetchUser(t *testing.T) {
    mockClient := &MockHTTPClient{
        getFunc: func(url string) IOE.IOEither[error, []byte] {
            user := User{ID: 1, Name: "Alice"}
            data, _ := json.Marshal(user)
            return IOE.Right[error](data)
        },
    }
    
    result := fetchUser(mockClient, 1)()
    
    if result.IsLeft() {
        t.Fatalf("Expected success, got error: %v", result.Left())
    }
    
    user := result.Right()
    if user.Name != "Alice" {
        t.Errorf("Expected 'Alice', got %q", user.Name)
    }
    
    // Verify HTTP call was made
    if len(mockClient.calls) != 1 {
        t.Errorf("Expected 1 HTTP call, got %d", len(mockClient.calls))
    }
}
`}
</CodeCard>

</Section>

<Section id="async-testing" number="05" title="Testing Async" titleAccent="Operations">

Test parallel execution and verify timing characteristics.

<CodeCard file="async_test.go">
{`package main

import (
    "testing"
    "time"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func slowOperation(n int) IOE.IOEither[error, int] {
    return IOE.TryCatch(func() (int, error) {
        time.Sleep(10 * time.Millisecond)
        return n * 2, nil
    })
}

func TestParallelExecution(t *testing.T) {
    numbers := []int{1, 2, 3, 4, 5}
    
    start := time.Now()
    result := A.Traverse[int](IOE.ApplicativePar[error, int]())(
        slowOperation,
    )(numbers)()
    duration := time.Since(start)
    
    if result.IsLeft() {
        t.Fatalf("Expected success, got error: %v", result.Left())
    }
    
    // Parallel execution should be faster than sequential
    // 5 operations * 10ms = 50ms sequential
    // Should complete in ~10-20ms parallel
    if duration > 30*time.Millisecond {
        t.Errorf("Parallel execution too slow: %v", duration)
    }
    
    expected := []int{2, 4, 6, 8, 10}
    if !equalSlices(result.Right(), expected) {
        t.Errorf("Expected %v, got %v", expected, result.Right())
    }
}

func equalSlices[A comparable](a, b []A) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
`}
</CodeCard>

</Section>

<Section id="error-retry-testing" number="06" title="Testing Error & Retry" titleAccent="Logic">

Test error handling and retry mechanisms with controlled failures.

<CodeCard file="retry_test.go">
{`package main

import (
    "fmt"
    "testing"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func withRetry[A any](maxAttempts int, operation IOE.IOEither[error, A]) IOE.IOEither[error, A] {
    return func() IOE.Either[error, A] {
        var lastErr error
        for i := 0; i < maxAttempts; i++ {
            result := operation()
            if result.IsRight() {
                return result
            }
            lastErr = result.Left()
        }
        return IOE.Left[A](fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr))()
    }
}

func TestRetry(t *testing.T) {
    t.Run("succeeds on first attempt", func(t *testing.T) {
        attempts := 0
        operation := IOE.TryCatch(func() (string, error) {
            attempts++
            return "success", nil
        })
        
        result := withRetry(3, operation)()
        
        if result.IsLeft() {
            t.Errorf("Expected success, got error: %v", result.Left())
        }
        
        if attempts != 1 {
            t.Errorf("Expected 1 attempt, got %d", attempts)
        }
    })
    
    t.Run("succeeds on third attempt", func(t *testing.T) {
        attempts := 0
        operation := IOE.TryCatch(func() (string, error) {
            attempts++
            if attempts < 3 {
                return "", fmt.Errorf("attempt %d failed", attempts)
            }
            return "success", nil
        })
        
        result := withRetry(3, operation)()
        
        if result.IsLeft() {
            t.Errorf("Expected success, got error: %v", result.Left())
        }
        
        if attempts != 3 {
            t.Errorf("Expected 3 attempts, got %d", attempts)
        }
    })
    
    t.Run("fails after max attempts", func(t *testing.T) {
        attempts := 0
        operation := IOE.TryCatch(func() (string, error) {
            attempts++
            return "", fmt.Errorf("attempt %d failed", attempts)
        })
        
        result := withRetry(3, operation)()
        
        if result.IsRight() {
            t.Error("Expected error, got success")
        }
        
        if attempts != 3 {
            t.Errorf("Expected 3 attempts, got %d", attempts)
        }
    })
}
`}
</CodeCard>

</Section>
