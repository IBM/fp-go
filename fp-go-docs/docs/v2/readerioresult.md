---
title: ReaderIOResult
hide_title: true
description: Dependency injection + lazy effects + Go error handling - Reader + IO + Result.
sidebar_position: 21
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="ReaderIOResult"
  lede="Combine dependency injection (Reader) + lazy effects (IO) + Go error handling (Result). ReaderIOResult[C, A] is ReaderIOEither specialized for error."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/readerioresult' },
    { label: 'Type', value: 'Monad (func(C) IO[Result[A]])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

<CodeCard file="type_definition.go">
{`package readerioresult

// ReaderIOResult is ReaderIOEither specialized for error
type ReaderIOResult[C, A any] = ReaderIOEither[C, error, A]
// Which expands to: func(C) func() Either[error, A]
`}
</CodeCard>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ok` | `func Ok[C, A any](value A) ReaderIOResult[C, A]` | Create successful value |
| `Error` | `func Error[C, A any](err error) ReaderIOResult[C, A]` | Create error value |
| `Of` | `func Of[C, A any](value A) ReaderIOResult[C, A]` | Alias for Ok |
| `Ask` | `func Ask[C any]() ReaderIOResult[C, C]` | Access context |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, A, B any](f func(A) B) func(ReaderIOResult[C, A]) ReaderIOResult[C, B]` | Transform success value |
| `Chain` | `func Chain[C, A, B any](f func(A) ReaderIOResult[C, B]) func(ReaderIOResult[C, A]) ReaderIOResult[C, B]` | Sequence operations |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Usage

<CodeCard file="basic.go">
{`package main

import (
    RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

type Dependencies struct {
    DB *sql.DB
}

func fetchUser(id string) RIOR.ReaderIOResult[Dependencies, User] {
    return RIOR.Ask[Dependencies, *sql.DB](func(deps Dependencies) *sql.DB {
        return deps.DB
    }).Chain(func(db *sql.DB) RIOR.ReaderIOResult[Dependencies, User] {
        return RIOR.TryCatchError(func() (User, error) {
            return db.QueryUser(id)
        })
    })
}
`}
</CodeCard>

</Section>
