---
title: ReaderEither
hide_title: true
description: Dependency injection with custom error handling - Reader + Either without IO.
sidebar_position: 18
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="ReaderEither"
  lede="Combine dependency injection (Reader) with custom error handling (Either). ReaderEither[C, E, A] for pure computations with dependencies and custom errors."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/readereither' },
    { label: 'Type', value: 'Monad (func(C) Either[E, A])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

<CodeCard file="type_definition.go">
{`package readereither

// ReaderEither combines Reader and Either
type ReaderEither[C, E, A any] = Reader[C, Either[E, A]]
// Which expands to: func(C) Either[E, A]
`}
</CodeCard>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Right` | `func Right[C, E, A any](value A) ReaderEither[C, E, A]` | Create successful value |
| `Left` | `func Left[C, E, A any](err E) ReaderEither[C, E, A]` | Create error value |
| `Of` | `func Of[C, E, A any](value A) ReaderEither[C, E, A]` | Alias for Right |
| `Ask` | `func Ask[C, E any]() ReaderEither[C, E, C]` | Access context |
| `Asks` | `func Asks[C, E, A any](f func(C) A) ReaderEither[C, E, A]` | Access and transform context |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, E, A, B any](f func(A) B) func(ReaderEither[C, E, A]) ReaderEither[C, E, B]` | Transform success value |
| `MapLeft` | `func MapLeft[C, A, E1, E2 any](f func(E1) E2) func(ReaderEither[C, E1, A]) ReaderEither[C, E2, A]` | Transform error |
| `Chain` | `func Chain[C, E, A, B any](f func(A) ReaderEither[C, E, B]) func(ReaderEither[C, E, A]) ReaderEither[C, E, B]` | Sequence operations |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Usage

<CodeCard file="basic.go">
{`package main

import (
    RE "github.com/IBM/fp-go/v2/readereither"
)

type Config struct {
    MaxRetries int
    Timeout    time.Duration
}

type ValidationError struct {
    Field   string
    Message string
}

func validateRetries() RE.ReaderEither[Config, ValidationError, int] {
    return RE.Asks(func(cfg Config) either.Either[ValidationError, int] {
        if cfg.MaxRetries < 1 {
            return either.Left[int](ValidationError{
                Field:   "maxRetries",
                Message: "must be at least 1",
            })
        }
        return either.Right[ValidationError](cfg.MaxRetries)
    })
}

func main() {
    cfg := Config{MaxRetries: 3, Timeout: time.Second}
    result := validateRetries()(cfg)
    // Either[ValidationError, int]
}
`}
</CodeCard>

</Section>
