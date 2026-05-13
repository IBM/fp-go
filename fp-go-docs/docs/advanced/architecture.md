---
title: Architecture Patterns
hide_title: true
description: Design scalable functional applications with hexagonal architecture, clean architecture, domain-driven design, and CQRS patterns.
sidebar_position: 4
---

<PageHeader
  eyebrow="Advanced · 04 / 04"
  title="Architecture"
  titleAccent="Patterns"
  lede="Design scalable functional applications with hexagonal architecture, clean architecture, domain-driven design, CQRS, event sourcing, and effect systems."
  meta={[
    { label: 'Difficulty', value: 'Expert' },
    { label: 'Topics', value: '6' },
    { label: 'Prerequisites', value: 'DDD, CQRS, Event Sourcing' }
  ]}
/>

<TLDR>
  <TLDRCard title="Hexagonal Architecture" icon="box">
    Isolate business logic from infrastructure with ports and adapters for testability.
  </TLDRCard>
  <TLDRCard title="Clean Architecture" icon="layers">
    Organize code in concentric layers with dependency inversion for maintainability.
  </TLDRCard>
  <TLDRCard title="Effect Systems" icon="zap">
    Track side effects in types with IO, Reader, and monad transformers for composability.
  </TLDRCard>
</TLDR>

<Section id="hexagonal" number="01" title="Hexagonal" titleAccent="Architecture">

Isolate business logic from infrastructure concerns using ports and adapters.

<CodeCard file="hexagonal_domain.go">
{`package domain

import (
    "context"
    "github.com/your-org/fp-go/either"
)

// Domain model (pure)
type User struct {
    ID    string
    Email string
    Name  string
}

type UserError struct {
    Code    string
    Message string
}

// Port (interface)
type UserRepository interface {
    FindByID(ctx context.Context, id string) either.Either[UserError, User]
    Save(ctx context.Context, user User) either.Either[UserError, User]
}

// Domain service (pure business logic)
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) either.Either[UserError, User] {
    return s.repo.FindByID(ctx, id)
}

func (s *UserService) UpdateEmail(ctx context.Context, id, newEmail string) either.Either[UserError, User] {
    return either.Chain(
        s.repo.FindByID(ctx, id),
        func(user User) either.Either[UserError, User] {
            user.Email = newEmail
            return s.repo.Save(ctx, user)
        },
    )
}
`}
</CodeCard>

<CodeCard file="hexagonal_adapter.go">
{`package postgres

import (
    "context"
    "database/sql"
    "github.com/your-org/fp-go/either"
    "myapp/domain"
)

// Adapter (infrastructure)
type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) either.Either[domain.UserError, domain.User] {
    var user domain.User
    err := r.db.QueryRowContext(ctx, "SELECT id, email, name FROM users WHERE id = $1", id).
        Scan(&user.ID, &user.Email, &user.Name)
    
    if err == sql.ErrNoRows {
        return either.Left[domain.UserError, domain.User](domain.UserError{
            Code:    "NOT_FOUND",
            Message: "User not found",
        })
    }
    
    if err != nil {
        return either.Left[domain.UserError, domain.User](domain.UserError{
            Code:    "DATABASE_ERROR",
            Message: err.Error(),
        })
    }
    
    return either.Right[domain.UserError](user)
}

func (r *PostgresUserRepository) Save(ctx context.Context, user domain.User) either.Either[domain.UserError, domain.User] {
    _, err := r.db.ExecContext(ctx,
        "UPDATE users SET email = $1, name = $2 WHERE id = $3",
        user.Email, user.Name, user.ID,
    )
    
    if err != nil {
        return either.Left[domain.UserError, domain.User](domain.UserError{
            Code:    "DATABASE_ERROR",
            Message: err.Error(),
        })
    }
    
    return either.Right[domain.UserError](user)
}
`}
</CodeCard>

</Section>

<Section id="clean-architecture" number="02" title="Clean" titleAccent="Architecture">

Organize code in concentric layers with dependency inversion.

<CodeCard file="clean_layers.go">
{`// Layer 1: Entities (innermost - pure domain)
package entities

type Order struct {
    ID         string
    CustomerID string
    Items      []OrderItem
    Total      float64
}

type OrderItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

// Layer 2: Use Cases (application logic)
package usecases

import (
    "context"
    "myapp/entities"
    "github.com/your-org/fp-go/either"
)

type OrderRepository interface {
    Save(ctx context.Context, order entities.Order) either.Either[error, entities.Order]
}

type PaymentGateway interface {
    Charge(ctx context.Context, amount float64, customerID string) either.Either[error, string]
}

type PlaceOrderUseCase struct {
    orderRepo OrderRepository
    payment   PaymentGateway
}

func NewPlaceOrderUseCase(repo OrderRepository, payment PaymentGateway) *PlaceOrderUseCase {
    return &PlaceOrderUseCase{
        orderRepo: repo,
        payment:   payment,
    }
}

func (uc *PlaceOrderUseCase) Execute(ctx context.Context, order entities.Order) either.Either[error, entities.Order] {
    // Charge payment
    chargeResult := uc.payment.Charge(ctx, order.Total, order.CustomerID)
    
    return either.Chain(
        chargeResult,
        func(transactionID string) either.Either[error, entities.Order] {
            // Save order
            return uc.orderRepo.Save(ctx, order)
        },
    )
}

// Layer 3: Interface Adapters (controllers, presenters)
package controllers

import (
    "encoding/json"
    "net/http"
    "myapp/usecases"
    "github.com/your-org/fp-go/either"
)

type OrderController struct {
    placeOrder *usecases.PlaceOrderUseCase
}

func NewOrderController(placeOrder *usecases.PlaceOrderUseCase) *OrderController {
    return &OrderController{placeOrder: placeOrder}
}

func (c *OrderController) PlaceOrder(w http.ResponseWriter, r *http.Request) {
    var req struct {
        CustomerID string  \`json:"customer_id"\`
        Items      []struct {
            ProductID string  \`json:"product_id"\`
            Quantity  int     \`json:"quantity"\`
            Price     float64 \`json:"price"\`
        } \`json:"items"\`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Convert to domain entity
    order := entities.Order{
        CustomerID: req.CustomerID,
        Items:      make([]entities.OrderItem, len(req.Items)),
    }
    
    for i, item := range req.Items {
        order.Items[i] = entities.OrderItem{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     item.Price,
        }
        order.Total += item.Price * float64(item.Quantity)
    }
    
    // Execute use case
    result := c.placeOrder.Execute(r.Context(), order)
    
    either.Match(
        result,
        func(err error) {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        },
        func(order entities.Order) {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(order)
        },
    )
}

// Layer 4: Frameworks & Drivers (outermost - infrastructure)
package main

import (
    "database/sql"
    "net/http"
    "myapp/controllers"
    "myapp/infrastructure"
    "myapp/usecases"
)

func main() {
    db, _ := sql.Open("postgres", "connection-string")
    
    // Wire dependencies (outer to inner)
    orderRepo := infrastructure.NewPostgresOrderRepository(db)
    paymentGateway := infrastructure.NewStripePaymentGateway("api-key")
    
    placeOrderUC := usecases.NewPlaceOrderUseCase(orderRepo, paymentGateway)
    orderController := controllers.NewOrderController(placeOrderUC)
    
    http.HandleFunc("/orders", orderController.PlaceOrder)
    http.ListenAndServe(":8080", nil)
}
`}
</CodeCard>

</Section>

<Section id="ddd" number="03" title="Domain-Driven" titleAccent="Design">

Model complex domains with aggregates, value objects, and domain events.

<CodeCard file="ddd_aggregate.go">
{`package domain

import (
    "time"
    "github.com/your-org/fp-go/either"
)

// Value Object (immutable)
type Money struct {
    Amount   float64
    Currency string
}

func NewMoney(amount float64, currency string) either.Either[error, Money] {
    if amount < 0 {
        return either.Left[error, Money](errors.New("amount cannot be negative"))
    }
    return either.Right[error](Money{Amount: amount, Currency: currency})
}

// Entity
type OrderLine struct {
    ProductID string
    Quantity  int
    Price     Money
}

// Aggregate Root
type Order struct {
    id         string
    customerID string
    lines      []OrderLine
    status     OrderStatus
    events     []DomainEvent
}

type OrderStatus string

const (
    OrderPending   OrderStatus = "PENDING"
    OrderConfirmed OrderStatus = "CONFIRMED"
    OrderShipped   OrderStatus = "SHIPPED"
)

// Domain Event
type DomainEvent interface {
    OccurredAt() time.Time
}

type OrderConfirmed struct {
    OrderID    string
    occurredAt time.Time
}

func (e OrderConfirmed) OccurredAt() time.Time {
    return e.occurredAt
}

// Factory
func NewOrder(id, customerID string) *Order {
    return &Order{
        id:         id,
        customerID: customerID,
        lines:      []OrderLine{},
        status:     OrderPending,
        events:     []DomainEvent{},
    }
}

// Business logic (domain methods)
func (o *Order) AddLine(productID string, quantity int, price Money) either.Either[error, *Order] {
    if o.status != OrderPending {
        return either.Left[error, *Order](errors.New("cannot modify confirmed order"))
    }
    
    o.lines = append(o.lines, OrderLine{
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
    })
    
    return either.Right[error](o)
}

func (o *Order) Confirm() either.Either[error, *Order] {
    if len(o.lines) == 0 {
        return either.Left[error, *Order](errors.New("cannot confirm empty order"))
    }
    
    if o.status != OrderPending {
        return either.Left[error, *Order](errors.New("order already confirmed"))
    }
    
    o.status = OrderConfirmed
    o.events = append(o.events, OrderConfirmed{
        OrderID:    o.id,
        occurredAt: time.Now(),
    })
    
    return either.Right[error](o)
}

func (o *Order) Total() Money {
    total := 0.0
    currency := "USD"
    
    for _, line := range o.lines {
        total += line.Price.Amount * float64(line.Quantity)
        currency = line.Price.Currency
    }
    
    return Money{Amount: total, Currency: currency}
}

func (o *Order) DomainEvents() []DomainEvent {
    return o.events
}

func (o *Order) ClearEvents() {
    o.events = []DomainEvent{}
}
`}
</CodeCard>

</Section>

<Section id="cqrs" number="04" title="CQRS" titleAccent="Pattern">

Separate read and write models for scalability and optimization.

<CodeCard file="cqrs_commands.go">
{`package commands

import (
    "context"
    "github.com/your-org/fp-go/either"
)

// Command (write model)
type CreateUserCommand struct {
    Email string
    Name  string
}

type CommandHandler[C any, R any] interface {
    Handle(ctx context.Context, cmd C) either.Either[error, R]
}

type CreateUserHandler struct {
    repo UserRepository
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) either.Either[error, string] {
    // Validate
    if cmd.Email == "" {
        return either.Left[error, string](errors.New("email required"))
    }
    
    // Create user
    user := User{
        ID:    generateID(),
        Email: cmd.Email,
        Name:  cmd.Name,
    }
    
    // Save
    return either.Map(
        h.repo.Save(ctx, user),
        func(u User) string { return u.ID },
    )
}
`}
</CodeCard>

<CodeCard file="cqrs_queries.go">
{`package queries

import (
    "context"
    "github.com/your-org/fp-go/either"
)

// Query (read model - optimized for reads)
type UserListQuery struct {
    Page     int
    PageSize int
}

type UserListItem struct {
    ID    string
    Email string
    Name  string
}

type QueryHandler[Q any, R any] interface {
    Handle(ctx context.Context, query Q) either.Either[error, R]
}

type UserListHandler struct {
    readDB ReadDatabase
}

func (h *UserListHandler) Handle(ctx context.Context, query UserListQuery) either.Either[error, []UserListItem] {
    offset := (query.Page - 1) * query.PageSize
    
    // Query optimized read model (could be denormalized, cached, etc.)
    rows, err := h.readDB.Query(ctx,
        "SELECT id, email, name FROM user_list_view LIMIT $1 OFFSET $2",
        query.PageSize, offset,
    )
    
    if err != nil {
        return either.Left[error, []UserListItem](err)
    }
    defer rows.Close()
    
    var users []UserListItem
    for rows.Next() {
        var user UserListItem
        if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
            return either.Left[error, []UserListItem](err)
        }
        users = append(users, user)
    }
    
    return either.Right[error](users)
}
`}
</CodeCard>

</Section>

<Section id="event-sourcing" number="05" title="Event" titleAccent="Sourcing">

Store state changes as a sequence of events for auditability and time travel.

<CodeCard file="event_sourcing.go">
{`package eventsourcing

import (
    "time"
)

// Event
type Event interface {
    AggregateID() string
    OccurredAt() time.Time
    EventType() string
}

type AccountCreated struct {
    AccountID  string
    Owner      string
    occurredAt time.Time
}

func (e AccountCreated) AggregateID() string { return e.AccountID }
func (e AccountCreated) OccurredAt() time.Time { return e.occurredAt }
func (e AccountCreated) EventType() string { return "AccountCreated" }

type MoneyDeposited struct {
    AccountID  string
    Amount     float64
    occurredAt time.Time
}

func (e MoneyDeposited) AggregateID() string { return e.AccountID }
func (e MoneyDeposited) OccurredAt() time.Time { return e.occurredAt }
func (e MoneyDeposited) EventType() string { return "MoneyDeposited" }

type MoneyWithdrawn struct {
    AccountID  string
    Amount     float64
    occurredAt time.Time
}

func (e MoneyWithdrawn) AggregateID() string { return e.AccountID }
func (e MoneyWithdrawn) OccurredAt() time.Time { return e.occurredAt }
func (e MoneyWithdrawn) EventType() string { return "MoneyWithdrawn" }

// Aggregate (rebuilt from events)
type Account struct {
    ID      string
    Owner   string
    Balance float64
    version int
}

func NewAccount() *Account {
    return &Account{}
}

// Apply events to rebuild state
func (a *Account) Apply(event Event) {
    switch e := event.(type) {
    case AccountCreated:
        a.ID = e.AccountID
        a.Owner = e.Owner
        a.Balance = 0
    case MoneyDeposited:
        a.Balance += e.Amount
    case MoneyWithdrawn:
        a.Balance -= e.Amount
    }
    a.version++
}

// Event Store
type EventStore interface {
    Save(events []Event) error
    Load(aggregateID string) ([]Event, error)
}

// Repository
type AccountRepository struct {
    store EventStore
}

func (r *AccountRepository) Load(id string) (*Account, error) {
    events, err := r.store.Load(id)
    if err != nil {
        return nil, err
    }
    
    account := NewAccount()
    for _, event := range events {
        account.Apply(event)
    }
    
    return account, nil
}

func (r *AccountRepository) Save(account *Account, events []Event) error {
    return r.store.Save(events)
}

// Usage
func main() {
    store := NewInMemoryEventStore()
    repo := &AccountRepository{store: store}
    
    // Create account
    events := []Event{
        AccountCreated{
            AccountID:  "acc-1",
            Owner:      "Alice",
            occurredAt: time.Now(),
        },
        MoneyDeposited{
            AccountID:  "acc-1",
            Amount:     100.0,
            occurredAt: time.Now(),
        },
    }
    
    repo.Save(nil, events)
    
    // Load account (rebuild from events)
    account, _ := repo.Load("acc-1")
    fmt.Println("Balance:", account.Balance) // 100.0
}
`}
</CodeCard>

</Section>

<Section id="effect-systems" number="06" title="Effect" titleAccent="Systems">

Track side effects in types for composability and testability.

<CodeCard file="effect_system.go">
{`package effects

import (
    "context"
    "github.com/your-org/fp-go/io"
    "github.com/your-org/fp-go/reader"
    "github.com/your-org/fp-go/readerioeither"
)

// Environment (dependencies)
type Env struct {
    DB     Database
    Logger Logger
    Config Config
}

// Effect type: ReaderIOEither[Env, Error, A]
// - Reader: needs Env
// - IO: performs side effects
// - Either: can fail with Error
type Effect[A any] = readerioeither.ReaderIOEither[Env, error, A]

// Pure computation (no effects)
func pure[A any](value A) Effect[A] {
    return readerioeither.Right[Env, error, A](value)
}

// Database effect
func findUser(id string) Effect[User] {
    return readerioeither.Ask[Env, error, User](func(env Env) io.IO[either.Either[error, User]] {
        return io.Of(func() either.Either[error, User] {
            return env.DB.FindUser(id)
        })
    })
}

// Logging effect
func logInfo(message string) Effect[unit.Unit] {
    return readerioeither.Ask[Env, error, unit.Unit](func(env Env) io.IO[either.Either[error, unit.Unit]] {
        return io.Of(func() either.Either[error, unit.Unit] {
            env.Logger.Info(message)
            return either.Right[error](unit.Unit{})
        })
    })
}

// Compose effects
func getUserWithLogging(id string) Effect[User] {
    return readerioeither.Chain(
        logInfo("Finding user: " + id),
        func(_ unit.Unit) Effect[User] {
            return readerioeither.Chain(
                findUser(id),
                func(user User) Effect[User] {
                    return readerioeither.Chain(
                        logInfo("Found user: " + user.Name),
                        func(_ unit.Unit) Effect[User] {
                            return pure(user)
                        },
                    )
                },
            )
        },
    )
}

// Run effect with environment
func main() {
    env := Env{
        DB:     NewPostgresDB(),
        Logger: NewLogger(),
        Config: LoadConfig(),
    }
    
    effect := getUserWithLogging("user-123")
    
    // Execute effect
    result := effect(env)() // Reader -> IO -> Either
    
    either.Match(
        result,
        func(err error) {
            fmt.Println("Error:", err)
        },
        func(user User) {
            fmt.Println("User:", user.Name)
        },
    )
}
`}
</CodeCard>

</Section>
