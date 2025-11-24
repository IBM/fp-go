# Idiomatic ReadIOResult Functions - Implementation Plan

## Overview

This document outlines the idiomatic functions that should be added to the `readerioresult` package to support Go's native `(value, error)` pattern, similar to what was implemented for `readerresult`.

## Key Concepts

The idiomatic package `github.com/IBM/fp-go/v2/idiomatic/readerioresult` defines:
- `ReaderIOResult[R, A]` as `func(R) func() (A, error)` (idiomatic style)
- This contrasts with `readerioresult.ReaderIOResult[R, A]` which is `Reader[R, IOResult[A]]` (functional style)

## Functions to Add

### In `readerioresult/reader.go`

Add helper functions at the top:
```go
func fromReaderIOResultKleisliI[R, A, B any](f RIORI.Kleisli[R, A, B]) Kleisli[R, A, B] {
	return function.Flow2(f, FromReaderIOResultI[R, B])
}

func fromIOResultKleisliI[A, B any](f IORI.Kleisli[A, B]) ioresult.Kleisli[A, B] {
	return ioresult.Eitherize1(f)
}
```

### Core Conversion Functions

1. **FromResultI** - Lift `(value, error)` to ReaderIOResult
   ```go
   func FromResultI[R, A any](a A, err error) ReaderIOResult[R, A]
   ```

2. **FromIOResultI** - Lift idiomatic IOResult to functional
   ```go
   func FromIOResultI[R, A any](ioe func() (A, error)) ReaderIOResult[R, A]
   ```

3. **FromReaderIOResultI** - Convert idiomatic ReaderIOResult to functional
   ```go
   func FromReaderIOResultI[R, A any](rr RIORI.ReaderIOResult[R, A]) ReaderIOResult[R, A]
   ```

### Chain Functions

4. **MonadChainI** / **ChainI** - Chain with idiomatic Kleisli
   ```go
   func MonadChainI[R, A, B any](ma ReaderIOResult[R, A], f RIORI.Kleisli[R, A, B]) ReaderIOResult[R, B]
   func ChainI[R, A, B any](f RIORI.Kleisli[R, A, B]) Operator[R, A, B]
   ```

5. **MonadChainEitherIK** / **ChainEitherIK** - Chain with idiomatic Result functions
   ```go
   func MonadChainEitherIK[R, A, B any](ma ReaderIOResult[R, A], f func(A) (B, error)) ReaderIOResult[R, B]
   func ChainEitherIK[R, A, B any](f func(A) (B, error)) Operator[R, A, B]
   ```

6. **MonadChainIOResultIK** / **ChainIOResultIK** - Chain with idiomatic IOResult
   ```go
   func MonadChainIOResultIK[R, A, B any](ma ReaderIOResult[R, A], f func(A) func() (B, error)) ReaderIOResult[R, B]
   func ChainIOResultIK[R, A, B any](f func(A) func() (B, error)) Operator[R, A, B]
   ```

### Applicative Functions

7. **MonadApI** / **ApI** - Apply with idiomatic value
   ```go
   func MonadApI[B, R, A any](fab ReaderIOResult[R, func(A) B], fa RIORI.ReaderIOResult[R, A]) ReaderIOResult[R, B]
   func ApI[B, R, A any](fa RIORI.ReaderIOResult[R, A]) Operator[R, func(A) B, B]
   ```

### Error Handling Functions

8. **OrElseI** - Fallback with idiomatic computation
   ```go
   func OrElseI[R, A any](onLeft RIORI.Kleisli[R, error, A]) Operator[R, A, A]
   ```

9. **MonadAltI** / **AltI** - Alternative with idiomatic computation
   ```go
   func MonadAltI[R, A any](first ReaderIOResult[R, A], second Lazy[RIORI.ReaderIOResult[R, A]]) ReaderIOResult[R, A]
   func AltI[R, A any](second Lazy[RIORI.ReaderIOResult[R, A]]) Operator[R, A, A]
   ```

### Flatten Functions

10. **FlattenI** - Flatten nested idiomatic ReaderIOResult
    ```go
    func FlattenI[R, A any](mma ReaderIOResult[R, RIORI.ReaderIOResult[R, A]]) ReaderIOResult[R, A]
    ```

### In `readerioresult/bind.go`

11. **BindI** - Bind with idiomatic Kleisli
    ```go
    func BindI[R, S1, S2, T any](setter func(T) func(S1) S2, f RIORI.Kleisli[R, S1, T]) Operator[R, S1, S2]
    ```

12. **ApIS** - Apply idiomatic value to state
    ```go
    func ApIS[R, S1, S2, T any](setter func(T) func(S1) S2, fa RIORI.ReaderIOResult[R, T]) Operator[R, S1, S2]
    ```

13. **ApISL** - Apply idiomatic value using lens
    ```go
    func ApISL[R, S, T any](lens L.Lens[S, T], fa RIORI.ReaderIOResult[R, T]) Operator[R, S, S]
    ```

14. **BindIL** - Bind idiomatic with lens
    ```go
    func BindIL[R, S, T any](lens L.Lens[S, T], f RIORI.Kleisli[R, T, T]) Operator[R, S, S]
    ```

15. **BindEitherIK** / **BindResultIK** - Bind idiomatic Result
    ```go
    func BindEitherIK[R, S1, S2, T any](setter func(T) func(S1) S2, f func(S1) (T, error)) Operator[R, S1, S2]
    func BindResultIK[R, S1, S2, T any](setter func(T) func(S1) S2, f func(S1) (T, error)) Operator[R, S1, S2]
    ```

16. **BindIOResultIK** - Bind idiomatic IOResult
    ```go
    func BindIOResultIK[R, S1, S2, T any](setter func(T) func(S1) S2, f func(S1) func() (T, error)) Operator[R, S1, S2]
    ```

17. **BindToEitherI** / **BindToResultI** - Initialize from idiomatic pair
    ```go
    func BindToEitherI[R, S1, T any](setter func(T) S1) func(T, error) ReaderIOResult[R, S1]
    func BindToResultI[R, S1, T any](setter func(T) S1) func(T, error) ReaderIOResult[R, S1]
    ```

18. **BindToIOResultI** - Initialize from idiomatic IOResult
    ```go
    func BindToIOResultI[R, S1, T any](setter func(T) S1) func(func() (T, error)) ReaderIOResult[R, S1]
    ```

19. **ApEitherIS** / **ApResultIS** - Apply idiomatic pair to state
    ```go
    func ApEitherIS[R, S1, S2, T any](setter func(T) func(S1) S2) func(T, error) Operator[R, S1, S2]
    func ApResultIS[R, S1, S2, T any](setter func(T) func(S1) S2) func(T, error) Operator[R, S1, S2]
    ```

20. **ApIOResultIS** - Apply idiomatic IOResult to state
    ```go
    func ApIOResultIS[R, S1, S2, T any](setter func(T) func(S1) S2, fa func() (T, error)) Operator[R, S1, S2]
    ```

## Testing Strategy

Create `readerioresult/idiomatic_test.go` with:
- Tests for each idiomatic function
- Success and error cases
- Integration tests showing real-world usage patterns
- Parallel execution tests where applicable
- Complex scenarios combining multiple idiomatic functions

## Implementation Priority

1. **High Priority** - Core conversion and chain functions (1-6)
2. **Medium Priority** - Bind functions for do-notation (11-16)
3. **Low Priority** - Advanced applicative and error handling (7-10, 17-20)

## Benefits

1. **Seamless Integration** - Mix Go idiomatic code with functional pipelines
2. **Gradual Adoption** - Convert code incrementally from idiomatic to functional
3. **Interoperability** - Work with existing Go libraries that return `(value, error)`
4. **Consistency** - Mirrors the successful pattern from `readerresult`

## References

- See `readerresult` package for similar implementations
- See `idiomatic/readerresult` for the idiomatic types
- See `idiomatic/ioresult` for IO-level idiomatic patterns
