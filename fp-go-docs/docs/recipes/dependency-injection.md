---
title: Dependency Injection
hide_title: true
description: Implement dependency injection using the Reader pattern with fp-go for testable, modular code.
sidebar_position: 13
---

<PageHeader
  eyebrow="Recipes · 13 / 17"
  title="Dependency"
  titleAccent="Injection"
  lede="Implement dependency injection using the Reader pattern with fp-go for testable, modular code without global state."
  meta={[
    { label: 'Difficulty', value: 'Advanced' },
    { label: 'Patterns', value: '6' },
    { label: 'Use Cases', value: 'Service layers, HTTP APIs, testing' }
  ]}
/>

<TLDR>
  <TLDRCard title="Type-Safe Dependencies" icon="shield">
    Reader pattern provides compile-time guarantees that all dependencies are provided before execution.
  </TLDRCard>
  <TLDRCard title="Pure & Testable" icon="check">
    No global state means easy mocking and testing—just pass different dependencies to the same code.
  </TLDRCard>
  <TLDRCard title="Composable Services" icon="layers">
    Combine multiple readers to build complex operations from simple, focused functions.
  </TLDRCard>
</TLDR>

<Section id="basic-reader" number="01" title="Basic Reader" titleAccent="Pattern">

The Reader pattern allows functions to access dependencies without passing them explicitly through every function call.

<CodeCard file="reader_basic.go">
{`package main

import (
    "fmt"
    R "github.com/IBM/fp-go/v2/reader"
)

type Config struct {
    APIKey string
    BaseURL string
}

// Reader that needs Config
type AppReader[A any] = R.Reader[Config, A]

func getAPIKey() AppReader[string] {
    return R.Asks(func(cfg Config) string {
        return cfg.APIKey
    })
}

func getBaseURL() AppReader[string] {
    return R.Asks(func(cfg Config) string {
        return cfg.BaseURL
    })
}

func buildURL(path string) AppReader[string] {
    return R.Map(func(base string) string {
        return base + path
    })(getBaseURL())
}

func main() {
    config := Config{
        APIKey: "secret-key-123",
        BaseURL: "https://api.example.com",
    }
    
    // Execute readers with config
    apiKey := getAPIKey()(config)
    url := buildURL("/users")(config)
    
    fmt.Println("API Key:", apiKey)
    fmt.Println("URL:", url)
    // API Key: secret-key-123
    // URL: https://api.example.com/users
}
`}
</CodeCard>

<Callout type="info">
**Reader vs Global State**: Reader makes dependencies explicit and testable. Global variables hide dependencies and make testing difficult.
</Callout>

</Section>

<Section id="service-layer" number="02" title="Service Layer" titleAccent="Pattern">

Build service layers with multiple dependencies using ReaderEither for error handling.

<CodeCard file="service_layer.go">
{`package main

import (
    "context"
    "fmt"
    RE "github.com/IBM/fp-go/v2/readereither"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

// Service interfaces
type UserService interface {
    GetUser(ctx context.Context, id int) IOE.IOEither[error, User]
    CreateUser(ctx context.Context, user User) IOE.IOEither[error, User]
}

type EmailService interface {
    SendEmail(ctx context.Context, to, subject, body string) IOE.IOEither[error, struct{}]
}

// Dependencies container
type Services struct {
    Users  UserService
    Emails EmailService
}

// Reader type for services
type ServiceReader[A any] = RE.ReaderEither[Services, error, A]

type User struct {
    ID    int
    Name  string
    Email string
}

// Service operations
func getUser(ctx context.Context, id int) ServiceReader[User] {
    return RE.Asks(func(services Services) IOE.Either[error, User] {
        return services.Users.GetUser(ctx, id)()
    })
}

func sendWelcomeEmail(ctx context.Context, user User) ServiceReader[struct{}] {
    return RE.Asks(func(services Services) IOE.Either[error, struct{}] {
        return services.Emails.SendEmail(
            ctx,
            user.Email,
            "Welcome!",
            fmt.Sprintf("Hello %s, welcome to our platform!", user.Name),
        )()
    })
}

// Composed operation
func createUserAndSendWelcome(ctx context.Context, user User) ServiceReader[User] {
    return F.Pipe3(
        RE.Do[Services, error](RE.Monad[Services, error, User]()),
        RE.Bind("user", func() ServiceReader[User] {
            return RE.Asks(func(services Services) IOE.Either[error, User] {
                return services.Users.CreateUser(ctx, user)()
            })
        }),
        RE.ChainFirst(func(u User) ServiceReader[struct{}] {
            return sendWelcomeEmail(ctx, u)
        }),
        RE.Map(func(data struct{ user User }) User {
            return data.user
        }),
    )
}

func main() {
    services := Services{
        Users:  &MockUserService{},
        Emails: &MockEmailService{},
    }
    
    ctx := context.Background()
    newUser := User{Name: "Alice", Email: "alice@example.com"}
    
    result := createUserAndSendWelcome(ctx, newUser)(services)
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Created user: %+v\\n", result.Right())
    }
}
`}
</CodeCard>

</Section>

<Section id="readerioeither" number="03" title="ReaderIOEither" titleAccent="Pattern">

Combine Reader with IO and error handling for real-world applications.

<CodeCard file="readerioeither.go">
{`package main

import (
    "context"
    "fmt"
    RIE "github.com/IBM/fp-go/v2/readerioeither"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type AppDeps struct {
    Config   Config
    Database Database
    Logger   Logger
}

type AppEffect[A any] = RIE.ReaderIOEither[AppDeps, error, A]

func logInfo(msg string) AppEffect[struct{}] {
    return RIE.Asks(func(deps AppDeps) IOE.IOEither[error, struct{}] {
        return IOE.TryCatch(func() (struct{}, error) {
            deps.Logger.Info(msg)
            return struct{}{}, nil
        })
    })
}

func queryDB(sql string) AppEffect[[]string] {
    return RIE.Asks(func(deps AppDeps) IOE.IOEither[error, []string] {
        return IOE.TryCatch(func() ([]string, error) {
            return deps.Database.Query(sql)
        })
    })
}

func getUsersWithLogging() AppEffect[[]string] {
    return F.Pipe3(
        logInfo("Starting user query"),
        RIE.Chain(func(_ struct{}) AppEffect[[]string] {
            return queryDB("SELECT * FROM users")
        }),
        RIE.ChainFirst(func(users []string) AppEffect[struct{}] {
            return logInfo(fmt.Sprintf("Found %d users", len(users)))
        }),
    )
}

func main() {
    deps := AppDeps{
        Config:   Config{},
        Database: &MockDatabase{},
        Logger:   &ConsoleLogger{},
    }
    
    result := getUsersWithLogging()(deps)()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Println("Users:", result.Right())
    }
}
`}
</CodeCard>

</Section>

<Section id="testing" number="04" title="Testing with" titleAccent="DI">

Dependency injection makes testing trivial—just provide mock implementations.

<CodeCard file="testing_di.go">
{`package main

import (
    "context"
    "testing"
)

// Mock implementations
type MockUserService struct {
    users map[int]User
}

func (m *MockUserService) GetUser(ctx context.Context, id int) IOE.IOEither[error, User] {
    return IOE.TryCatch(func() (User, error) {
        if user, ok := m.users[id]; ok {
            return user, nil
        }
        return User{}, fmt.Errorf("user not found: %d", id)
    })
}

func (m *MockUserService) CreateUser(ctx context.Context, user User) IOE.IOEither[error, User] {
    return IOE.TryCatch(func() (User, error) {
        user.ID = len(m.users) + 1
        m.users[user.ID] = user
        return user, nil
    })
}

type MockEmailService struct {
    sentEmails []string
}

func (m *MockEmailService) SendEmail(ctx context.Context, to, subject, body string) IOE.IOEither[error, struct{}] {
    return IOE.TryCatch(func() (struct{}, error) {
        m.sentEmails = append(m.sentEmails, to)
        return struct{}{}, nil
    })
}

// Test
func TestCreateUserAndSendWelcome(t *testing.T) {
    mockUsers := &MockUserService{
        users: make(map[int]User),
    }
    mockEmails := &MockEmailService{
        sentEmails: []string{},
    }
    
    services := Services{
        Users:  mockUsers,
        Emails: mockEmails,
    }
    
    ctx := context.Background()
    newUser := User{Name: "Bob", Email: "bob@example.com"}
    
    result := createUserAndSendWelcome(ctx, newUser)(services)
    
    if result.IsLeft() {
        t.Fatalf("Expected success, got error: %v", result.Left())
    }
    
    user := result.Right()
    if user.ID == 0 {
        t.Error("Expected user ID to be set")
    }
    
    if len(mockEmails.sentEmails) != 1 {
        t.Errorf("Expected 1 email sent, got %d", len(mockEmails.sentEmails))
    }
    
    if mockEmails.sentEmails[0] != "bob@example.com" {
        t.Errorf("Expected email to bob@example.com, got %s", mockEmails.sentEmails[0])
    }
}
`}
</CodeCard>

</Section>

<Section id="http-api" number="05" title="HTTP API" titleAccent="Example">

Complete HTTP API with dependency injection for clean, testable handlers.

<CodeCard file="http_api.go">
{`package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    RIE "github.com/IBM/fp-go/v2/readerioeither"
    F "github.com/IBM/fp-go/v2/function"
)

// Application dependencies
type AppContext struct {
    UserRepo  UserRepository
    Logger    Logger
    Config    AppConfig
}

type AppConfig struct {
    Port     int
    LogLevel string
}

type UserRepository interface {
    FindByID(ctx context.Context, id int) IOE.IOEither[error, User]
    Save(ctx context.Context, user User) IOE.IOEither[error, User]
    List(ctx context.Context) IOE.IOEither[error, []User]
}

type AppHandler[A any] = RIE.ReaderIOEither[AppContext, error, A]

// Handler operations
func logRequest(method, path string) AppHandler[struct{}] {
    return RIE.Asks(func(app AppContext) IOE.IOEither[error, struct{}] {
        return IOE.TryCatch(func() (struct{}, error) {
            app.Logger.Info(fmt.Sprintf("%s %s", method, path))
            return struct{}{}, nil
        })
    })
}

func getUserByID(ctx context.Context, id int) AppHandler[User] {
    return RIE.Asks(func(app AppContext) IOE.IOEither[error, User] {
        return app.UserRepo.FindByID(ctx, id)
    })
}

func createUser(ctx context.Context, user User) AppHandler[User] {
    return RIE.Asks(func(app AppContext) IOE.IOEither[error, User] {
        return app.UserRepo.Save(ctx, user)
    })
}

// HTTP handlers
func makeGetUserHandler(app AppContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := 1 // Parse from URL
        
        handler := F.Pipe2(
            logRequest("GET", r.URL.Path),
            RIE.Chain(func(_ struct{}) AppHandler[User] {
                return getUserByID(r.Context(), id)
            }),
        )
        
        result := handler(app)()
        
        if result.IsLeft() {
            http.Error(w, result.Left().Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(result.Right())
    }
}

func main() {
    app := AppContext{
        UserRepo: &InMemoryUserRepo{users: make(map[int]User)},
        Logger:   &ConsoleLogger{},
        Config:   AppConfig{Port: 8080, LogLevel: "info"},
    }
    
    http.HandleFunc("/users", makeCreateUserHandler(app))
    http.HandleFunc("/users/", makeGetUserHandler(app))
    
    fmt.Printf("Server starting on port %d\\n", app.Config.Port)
    http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
}
`}
</CodeCard>

</Section>

<Section id="best-practices" number="06" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Keep dependencies minimal** — Only include what each function actually needs
  </ChecklistItem>
  <ChecklistItem status="required">
    **Use interfaces** — Define interfaces for all dependencies to enable mocking
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Compose small readers** — Build complex operations from simple, focused functions
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Avoid kitchen sink** — Don't pass every possible dependency to every function
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Scope dependencies** — Use nested readers for request-scoped vs app-scoped dependencies
  </ChecklistItem>
</Checklist>

<Compare>
<CompareCol kind="good">
<CodeCard file="good_di.go">
{`// ✅ Good: Minimal, focused dependencies
type UserHandlerDeps struct {
    UserRepo UserRepository
    Logger   Logger
}

// ✅ Good: Interface for testability
type Logger interface {
    Info(msg string)
    Error(msg string)
}

// ✅ Good: Small, composable readers
func getUser(id int) AppReader[User] { /* ... */ }
func validateUser(user User) AppReader[User] { /* ... */ }
func saveUser(user User) AppReader[User] { /* ... */ }

func createUser(user User) AppReader[User] {
    return F.Pipe2(
        validateUser(user),
        R.Chain(saveUser),
    )
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="bad">
<CodeCard file="bad_di.go">
{`// ❌ Avoid: Kitchen sink dependencies
type UserHandlerDeps struct {
    UserRepo    UserRepository
    EmailRepo   EmailRepository
    PaymentRepo PaymentRepository
    Logger      Logger
    Cache       Cache
    Queue       Queue
    // ... everything
}

// ❌ Avoid: Concrete types
type Logger struct {
    file *os.File
}

// ❌ Avoid: Monolithic readers
func createUser(user User) AppReader[User] {
    // 100 lines of logic
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
