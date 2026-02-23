# CLI Package - Functional Wrappers for urfave/cli/v3

This package provides functional programming wrappers for the `github.com/urfave/cli/v3` library, enabling Effect-based command actions and type-safe flag handling through Prisms.

## Features

### 1. Effect-Based Command Actions

Transform CLI command actions into composable Effects that follow functional programming principles.

#### Key Functions

- **`ToAction(effect CommandEffect) func(context.Context, *C.Command) error`**
  - Converts a CommandEffect into a standard urfave/cli Action function
  - Enables Effect-based command handlers to work with cli/v3 framework

- **`FromAction(action func(context.Context, *C.Command) error) CommandEffect`**
  - Lifts existing cli/v3 action handlers into the Effect type
  - Allows gradual migration to functional style

- **`MakeCommand(name, usage string, flags []C.Flag, effect CommandEffect) *C.Command`**
  - Creates a new Command with an Effect-based action
  - Convenience function combining command creation with Effect conversion

- **`MakeCommandWithSubcommands(...) *C.Command`**
  - Creates a Command with subcommands and an Effect-based action

#### Example Usage

```go
import (
    "context"
    
    E "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/cli"
    C "github.com/urfave/cli/v3"
)

// Define an Effect-based command action
processEffect := func(cmd *C.Command) E.Thunk[F.Void] {
    return func(ctx context.Context) E.IOResult[F.Void] {
        return func() R.Result[F.Void] {
            input := cmd.String("input")
            // Process input...
            return R.Of(F.Void{})
        }
    }
}

// Create command with Effect
command := cli.MakeCommand(
    "process",
    "Process input files",
    []C.Flag{
        &C.StringFlag{Name: "input", Usage: "Input file path"},
    },
    processEffect,
)

// Or convert existing action to Effect
existingAction := func(ctx context.Context, cmd *C.Command) error {
    // Existing logic...
    return nil
}
effect := cli.FromAction(existingAction)
```

### 2. Flag Type Prisms

Type-safe extraction and manipulation of CLI flags using Prisms from the optics package.

#### Available Prisms

- `StringFlagPrism()` - Extract `*C.StringFlag` from `C.Flag`
- `IntFlagPrism()` - Extract `*C.IntFlag` from `C.Flag`
- `BoolFlagPrism()` - Extract `*C.BoolFlag` from `C.Flag`
- `Float64FlagPrism()` - Extract `*C.Float64Flag` from `C.Flag`
- `DurationFlagPrism()` - Extract `*C.DurationFlag` from `C.Flag`
- `TimestampFlagPrism()` - Extract `*C.TimestampFlag` from `C.Flag`
- `StringSliceFlagPrism()` - Extract `*C.StringSliceFlag` from `C.Flag`
- `IntSliceFlagPrism()` - Extract `*C.IntSliceFlag` from `C.Flag`
- `Float64SliceFlagPrism()` - Extract `*C.Float64SliceFlag` from `C.Flag`
- `UintFlagPrism()` - Extract `*C.UintFlag` from `C.Flag`
- `Uint64FlagPrism()` - Extract `*C.Uint64Flag` from `C.Flag`
- `Int64FlagPrism()` - Extract `*C.Int64Flag` from `C.Flag`

#### Example Usage

```go
import (
    O "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/cli"
    C "github.com/urfave/cli/v3"
)

// Extract a StringFlag from a Flag interface
var flag C.Flag = &C.StringFlag{Name: "input", Value: "default"}
prism := cli.StringFlagPrism()

// Safe extraction returns Option
result := prism.GetOption(flag)
if O.IsSome(result) {
    strFlag := O.MonadFold(result, 
        func() *C.StringFlag { return nil },
        func(f *C.StringFlag) *C.StringFlag { return f },
    )
    // Use strFlag...
}

// Type mismatch returns None
var intFlag C.Flag = &C.IntFlag{Name: "count"}
result = prism.GetOption(intFlag)  // Returns None

// Convert back to Flag
strFlag := &C.StringFlag{Name: "output"}
flag = prism.ReverseGet(strFlag)
```

## Type Definitions

### CommandEffect

```go
type CommandEffect = E.Effect[*C.Command, F.Void]
```

A CommandEffect represents a CLI command action as an Effect. It takes a `*C.Command` as context and produces a result wrapped in the Effect monad.

The Effect structure is:
```
func(*C.Command) -> func(context.Context) -> func() -> Result[Void]
```

This allows for:
- **Composability**: Effects can be composed using standard functional combinators
- **Testability**: Pure functions are easier to test
- **Error Handling**: Errors are explicitly represented in the Result type
- **Context Management**: Context flows naturally through the Effect

## Benefits

### 1. Functional Composition

Effects can be composed using standard functional programming patterns:

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    RRIOE "github.com/IBM/fp-go/v2/context/readerreaderioresult"
)

// Compose multiple effects
validateInput := func(cmd *C.Command) E.Thunk[F.Void] { /* ... */ }
processData := func(cmd *C.Command) E.Thunk[F.Void] { /* ... */ }
saveResults := func(cmd *C.Command) E.Thunk[F.Void] { /* ... */ }

// Chain effects together
pipeline := F.Pipe3(
    validateInput,
    RRIOE.Chain(func(F.Void) E.Effect[*C.Command, F.Void] { return processData }),
    RRIOE.Chain(func(F.Void) E.Effect[*C.Command, F.Void] { return saveResults }),
)
```

### 2. Type Safety

Prisms provide compile-time type safety when working with flags:

```go
// Type-safe flag extraction
flags := []C.Flag{
    &C.StringFlag{Name: "input"},
    &C.IntFlag{Name: "count"},
}

for _, flag := range flags {
    // Safe extraction with pattern matching
    O.MonadFold(
        cli.StringFlagPrism().GetOption(flag),
        func() { /* Not a string flag */ },
        func(sf *C.StringFlag) { /* Handle string flag */ },
    )
}
```

### 3. Error Handling

Errors are explicitly represented in the Result type:

```go
effect := func(cmd *C.Command) E.Thunk[F.Void] {
    return func(ctx context.Context) E.IOResult[F.Void] {
        return func() R.Result[F.Void] {
            if err := validateInput(cmd); err != nil {
                return R.Left[F.Void](err)  // Explicit error
            }
            return R.Of(F.Void{})  // Success
        }
    }
}
```

### 4. Testability

Pure functions are easier to test:

```go
func TestCommandEffect(t *testing.T) {
    cmd := &C.Command{Name: "test"}
    effect := myCommandEffect(cmd)
    
    // Execute effect
    result := effect(context.Background())()
    
    // Assert on result
    assert.True(t, R.IsRight(result))
}
```

## Migration Guide

### From Standard Actions to Effects

**Before:**
```go
command := &C.Command{
    Name: "process",
    Action: func(ctx context.Context, cmd *C.Command) error {
        input := cmd.String("input")
        // Process...
        return nil
    },
}
```

**After:**
```go
effect := func(cmd *C.Command) E.Thunk[F.Void] {
    return func(ctx context.Context) E.IOResult[F.Void] {
        return func() R.Result[F.Void] {
            input := cmd.String("input")
            // Process...
            return R.Of(F.Void{})
        }
    }
}

command := cli.MakeCommand("process", "Process files", flags, effect)
```

### Gradual Migration

You can mix both styles during migration:

```go
// Wrap existing action
existingAction := func(ctx context.Context, cmd *C.Command) error {
    // Legacy code...
    return nil
}

// Use as Effect
effect := cli.FromAction(existingAction)
command := cli.MakeCommand("legacy", "Legacy command", flags, effect)
```

## See Also

- [Effect Package](../effect/) - Core Effect type definitions
- [Optics Package](../optics/) - Prism and other optics
- [urfave/cli/v3](https://github.com/urfave/cli) - Underlying CLI framework