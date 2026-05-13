---
title: StateReaderIOEither
hide_title: true
description: The ultimate monad transformer - combines State, Reader, IO, and Either for maximum expressiveness.
sidebar_position: 24
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="StateReaderIOEither"
  lede="The ultimate monad transformer. StateReaderIOEither[S, C, E, A] combines stateful computations (State), dependency injection (Reader), lazy evaluation (IO), and error handling (Either)."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/statereaderioeither' },
    { label: 'Type', value: 'func(C) func(S) IO[Either[E, Pair[A, S]]]' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

StateReaderIOEither combines all four powerful abstractions:
- **State[S]**: Stateful computations
- **Reader[C]**: Dependency injection
- **IO**: Lazy evaluation and side effects
- **Either[E]**: Error handling

<CodeCard file="type_definition.go">
{`package statereaderioeither

// StateReaderIOEither combines all four abstractions
type StateReaderIOEither[S, C, E, A any] = 
    func(C) func(S) IO[Either[E, pair.Pair[A, S]]]
`}
</CodeCard>

### When to Use

<ApiTable>
| Use Case | Example |
|----------|---------|
| Complex application state | Games, simulations |
| Stateful services | Session management, transactions |
| Advanced scenarios | When you need all four capabilities |
</ApiTable>

<Callout type="warn">

**Complexity Warning**: Most applications don't need this level of complexity. Consider simpler alternatives first:
- **Effect** (ReaderIOResult) for most applications
- **State** for pure stateful computations
- **ReaderIOEither** for effects with dependencies

</Callout>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[S, C, E, A any](value A) StateReaderIOEither[S, C, E, A]` | Wrap pure value |
| `Left` | `func Left[S, C, E, A any](err E) StateReaderIOEither[S, C, E, A]` | Create error |
| `Right` | `func Right[S, C, E, A any](value A) StateReaderIOEither[S, C, E, A]` | Create success (alias for Of) |
</ApiTable>

### State Operations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Get` | `func Get[S, C, E any]() StateReaderIOEither[S, C, E, S]` | Get current state |
| `Put` | `func Put[S, C, E any](s S) StateReaderIOEither[S, C, E, unit.Unit]` | Set new state |
| `Modify` | `func Modify[S, C, E any](f func(S) S) StateReaderIOEither[S, C, E, unit.Unit]` | Modify state |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[S, C, E, A, B any](f func(A) B) func(StateReaderIOEither[S, C, E, A]) StateReaderIOEither[S, C, E, B]` | Transform value |
| `Chain` | `func Chain[S, C, E, A, B any](f func(A) StateReaderIOEither[S, C, E, B]) func(StateReaderIOEither[S, C, E, A]) StateReaderIOEither[S, C, E, B]` | Sequence operations |
</ApiTable>

### Dependencies

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ask` | `func Ask[S, C, E any]() StateReaderIOEither[S, C, E, C]` | Access dependencies |
| `Asks` | `func Asks[S, C, E, A any](f func(C) A) StateReaderIOEither[S, C, E, A]` | Access and transform |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Game State Example

<CodeCard file="game.go">
{`package main

import (
    "fmt"
    SRIE "github.com/IBM/fp-go/v2/statereaderioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type GameState struct {
    Score int
    Lives int
    Level int
}

type Dependencies struct {
    Logger *log.Logger
}

type GameError struct {
    Message string
}

func IncreaseScore(points int) SRIE.StateReaderIOEither[GameState, Dependencies, GameError, unit.Unit] {
    return func(deps Dependencies) func(state GameState) ioeither.IOEither[GameError, pair.Pair[unit.Unit, GameState]] {
        return func(s GameState) ioeither.IOEither[GameError, pair.Pair[unit.Unit, GameState]] {
            return ioeither.TryCatch(func() (pair.Pair[unit.Unit, GameState], GameError) {
                newState := s
                newState.Score += points
                deps.Logger.Printf("Score increased by %d", points)
                return pair.MakePair(unit.VOID, newState), GameError{}
            })
        }
    }
}

func LoseLife() SRIE.StateReaderIOEither[GameState, Dependencies, GameError, bool] {
    return func(deps Dependencies) func(state GameState) ioeither.IOEither[GameError, pair.Pair[bool, GameState]] {
        return func(s GameState) ioeither.IOEither[GameError, pair.Pair[bool, GameState]] {
            return ioeither.TryCatch(func() (pair.Pair[bool, GameState], GameError) {
                if s.Lives <= 0 {
                    return pair.MakePair(false, s), GameError{Message: "Game Over"}
                }
                newState := s
                newState.Lives--
                deps.Logger.Printf("Life lost. Lives remaining: %d", newState.Lives)
                return pair.MakePair(true, newState), GameError{}
            })
        }
    }
}

func main() {
    deps := Dependencies{Logger: log.New(os.Stdout, "", 0)}
    initialState := GameState{Score: 0, Lives: 3, Level: 1}
    
    game := F.Pipe2(
        IncreaseScore(100),
        SRIE.Chain(func(_ unit.Unit) SRIE.StateReaderIOEither[GameState, Dependencies, GameError, bool] {
            return LoseLife()
        }),
    )
    
    result := game(deps)(initialState)()
    fmt.Printf("%+v\n", result)
}
`}
</CodeCard>

### Session Management

<CodeCard file="session.go">
{`package main

import (
    SRIE "github.com/IBM/fp-go/v2/statereaderioeither"
)

type SessionState struct {
    UserID    string
    Token     string
    ExpiresAt time.Time
}

type AppDeps struct {
    DB    *sql.DB
    Cache *Cache
}

type SessionError struct {
    Code    int
    Message string
}

func ValidateSession() SRIE.StateReaderIOEither[SessionState, AppDeps, SessionError, bool] {
    return func(deps AppDeps) func(state SessionState) ioeither.IOEither[SessionError, pair.Pair[bool, SessionState]] {
        return func(s SessionState) ioeither.IOEither[SessionError, pair.Pair[bool, SessionState]] {
            return ioeither.TryCatch(func() (pair.Pair[bool, SessionState], SessionError) {
                if time.Now().After(s.ExpiresAt) {
                    return pair.MakePair(false, s), SessionError{
                        Code:    401,
                        Message: "Session expired",
                    }
                }
                return pair.MakePair(true, s), SessionError{}
            })
        }
    }
}

func RefreshSession() SRIE.StateReaderIOEither[SessionState, AppDeps, SessionError, unit.Unit] {
    return func(deps AppDeps) func(state SessionState) ioeither.IOEither[SessionError, pair.Pair[unit.Unit, SessionState]] {
        return func(s SessionState) ioeither.IOEither[SessionError, pair.Pair[unit.Unit, SessionState]] {
            return ioeither.TryCatch(func() (pair.Pair[unit.Unit, SessionState], SessionError) {
                newState := s
                newState.ExpiresAt = time.Now().Add(30 * time.Minute)
                
                // Update in cache
                deps.Cache.Set(s.UserID, newState)
                
                return pair.MakePair(unit.VOID, newState), SessionError{}
            })
        }
    }
}
`}
</CodeCard>

### Transaction Processing

<CodeCard file="transaction.go">
{`package main

import (
    SRIE "github.com/IBM/fp-go/v2/statereaderioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type TransactionState struct {
    Balance float64
    History []Transaction
}

type BankDeps struct {
    DB     *sql.DB
    Logger *log.Logger
}

type BankError struct {
    Code    string
    Message string
}

func Debit(amount float64) SRIE.StateReaderIOEither[TransactionState, BankDeps, BankError, unit.Unit] {
    return func(deps BankDeps) func(state TransactionState) ioeither.IOEither[BankError, pair.Pair[unit.Unit, TransactionState]] {
        return func(s TransactionState) ioeither.IOEither[BankError, pair.Pair[unit.Unit, TransactionState]] {
            return ioeither.TryCatch(func() (pair.Pair[unit.Unit, TransactionState], BankError) {
                if s.Balance < amount {
                    return pair.MakePair(unit.VOID, s), BankError{
                        Code:    "INSUFFICIENT_FUNDS",
                        Message: "Insufficient balance",
                    }
                }
                
                newState := s
                newState.Balance -= amount
                newState.History = append(newState.History, Transaction{
                    Type:   "DEBIT",
                    Amount: amount,
                    Time:   time.Now(),
                })
                
                deps.Logger.Printf("Debited: $%.2f", amount)
                return pair.MakePair(unit.VOID, newState), BankError{}
            })
        }
    }
}

func Credit(amount float64) SRIE.StateReaderIOEither[TransactionState, BankDeps, BankError, unit.Unit] {
    return func(deps BankDeps) func(state TransactionState) ioeither.IOEither[BankError, pair.Pair[unit.Unit, TransactionState]] {
        return func(s TransactionState) ioeither.IOEither[BankError, pair.Pair[unit.Unit, TransactionState]] {
            return ioeither.TryCatch(func() (pair.Pair[unit.Unit, TransactionState], BankError) {
                newState := s
                newState.Balance += amount
                newState.History = append(newState.History, Transaction{
                    Type:   "CREDIT",
                    Amount: amount,
                    Time:   time.Now(),
                })
                
                deps.Logger.Printf("Credited: $%.2f", amount)
                return pair.MakePair(unit.VOID, newState), BankError{}
            })
        }
    }
}

func Transfer(amount float64) SRIE.StateReaderIOEither[TransactionState, BankDeps, BankError, unit.Unit] {
    return F.Pipe2(
        Debit(amount),
        SRIE.Chain(func(_ unit.Unit) SRIE.StateReaderIOEither[TransactionState, BankDeps, BankError, unit.Unit] {
            return Credit(amount)
        }),
    )
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Workflow Engine

<CodeCard file="workflow.go">
{`package main

import (
    SRIE "github.com/IBM/fp-go/v2/statereaderioeither"
)

type WorkflowState struct {
    CurrentStep int
    Data        map[string]any
    Completed   bool
}

type WorkflowDeps struct {
    DB     *sql.DB
    Logger *log.Logger
}

type WorkflowError struct {
    Step    int
    Message string
}

func ExecuteStep(step int) SRIE.StateReaderIOEither[WorkflowState, WorkflowDeps, WorkflowError, unit.Unit] {
    return func(deps WorkflowDeps) func(state WorkflowState) ioeither.IOEither[WorkflowError, pair.Pair[unit.Unit, WorkflowState]] {
        return func(s WorkflowState) ioeither.IOEither[WorkflowError, pair.Pair[unit.Unit, WorkflowState]] {
            return ioeither.TryCatch(func() (pair.Pair[unit.Unit, WorkflowState], WorkflowError) {
                deps.Logger.Printf("Executing step %d", step)
                
                newState := s
                newState.CurrentStep = step
                
                // Execute step logic
                // ...
                
                return pair.MakePair(unit.VOID, newState), WorkflowError{}
            })
        }
    }
}
`}
</CodeCard>

### Pattern 2: State Machine

<CodeCard file="state_machine.go">
{`package main

import (
    SRIE "github.com/IBM/fp-go/v2/statereaderioeither"
)

type MachineState struct {
    Current string
    History []string
}

type MachineDeps struct {
    Logger *log.Logger
}

type MachineError struct {
    From    string
    To      string
    Message string
}

func Transition(to string) SRIE.StateReaderIOEither[MachineState, MachineDeps, MachineError, unit.Unit] {
    return func(deps MachineDeps) func(state MachineState) ioeither.IOEither[MachineError, pair.Pair[unit.Unit, MachineState]] {
        return func(s MachineState) ioeither.IOEither[MachineError, pair.Pair[unit.Unit, MachineState]] {
            return ioeither.TryCatch(func() (pair.Pair[unit.Unit, MachineState], MachineError) {
                // Validate transition
                if !isValidTransition(s.Current, to) {
                    return pair.MakePair(unit.VOID, s), MachineError{
                        From:    s.Current,
                        To:      to,
                        Message: "Invalid transition",
                    }
                }
                
                newState := s
                newState.History = append(newState.History, s.Current)
                newState.Current = to
                
                deps.Logger.Printf("Transitioned: %s -> %s", s.Current, to)
                return pair.MakePair(unit.VOID, newState), MachineError{}
            })
        }
    }
}
`}
</CodeCard>

</Section>

<Callout type="info">

**When to Use**: StateReaderIOEither is powerful but complex. Use it only when you truly need:
1. **Stateful computations** that can't be avoided
2. **Dependency injection** for testability
3. **Lazy evaluation** for control flow
4. **Error handling** with specific error types

For most cases, **Effect** (ReaderIOResult) is sufficient and simpler.

</Callout>
