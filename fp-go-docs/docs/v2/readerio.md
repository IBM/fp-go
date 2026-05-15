---
title: ReaderIO
hide_title: true
description: Dependency injection with lazy side effects - Reader + IO without error handling.
sidebar_position: 19
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="ReaderIO"
  lede="Combine dependency injection (Reader) with lazy side effects (IO). ReaderIO[C, A] for effectful computations with dependencies."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/readerio' },
    { label: 'Type', value: 'Monad (func(C) IO[A])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

<CodeCard file="type_definition.go">
{`package readerio

// ReaderIO combines Reader and IO
type ReaderIO[C, A any] = Reader[C, IO[A]]
// Which expands to: func(C) func() A
`}
</CodeCard>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[C, A any](value A) ReaderIO[C, A]` | Wrap pure value |
| `Ask` | `func Ask[C any]() ReaderIO[C, C]` | Access context |
| `Asks` | `func Asks[C, A any](f func(C) A) ReaderIO[C, A]` | Access and transform context |
| `FromIO` | `func FromIO[C, A any](io IO[A]) ReaderIO[C, A]` | Lift IO to ReaderIO |
| `FromReader` | `func FromReader[C, A any](r Reader[C, A]) ReaderIO[C, A]` | Lift Reader to ReaderIO |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, A, B any](f func(A) B) func(ReaderIO[C, A]) ReaderIO[C, B]` | Transform result |
| `Chain` | `func Chain[C, A, B any](f func(A) ReaderIO[C, B]) func(ReaderIO[C, A]) ReaderIO[C, B]` | Sequence operations |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Usage

<CodeCard file="basic.go">
{`package main

import (
    RIO "github.com/IBM/fp-go/v2/readerio"
    IO "github.com/IBM/fp-go/v2/io"
)

type Dependencies struct {
    Logger *log.Logger
}

func logMessage(msg string) RIO.ReaderIO[Dependencies, unit.Unit] {
    return RIO.Ask[Dependencies, *log.Logger](func(deps Dependencies) *log.Logger {
        return deps.Logger
    }).Chain(func(logger *log.Logger) RIO.ReaderIO[Dependencies, unit.Unit] {
        return RIO.FromIO(IO.FromImpure(func() {
            logger.Println(msg)
        }))
    })
}

func main() {
    deps := Dependencies{Logger: log.New(os.Stdout, "", 0)}
    logMessage("Hello, World!")(deps)()
}
`}
</CodeCard>

</Section>
