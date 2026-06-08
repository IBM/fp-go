# fp-go Documentation Checklist

Tracking sheet for the Docusaurus docs site (`fp-go-docs`). Use this to add
documentation incrementally — one (or a few) packages per PR.

- ✅ = doc page exists and is wired into the sidebar (`sidebars.ts`)
- ⬜ = not documented yet

> Source of truth for API packages is `v2/<package>`. Internal/support packages
> (`assert`, `cli`, `constraints`, `erasure`, `idiomatic`, `reflect`, `samples`)
> and non-code dirs (`lambda`, `resources`) are intentionally excluded.

---

# Guides & Narrative Docs

## Getting Started

| Status | Page | Description |
| --- | --- | --- |
| ⬜ | `intro` | "What is fp-go?" — landing/overview page. |
| ⬜ | `installation` | How to install fp-go and set up your module. |
| ⬜ | `quickstart` | 5-minute hands-on introduction. |
| ⬜ | `why-fp-go` | Motivation and benefits of using fp-go. |
| ⬜ | `comparison` | Comparison with other FP libraries and idiomatic Go. |
| ⬜ | `faq` | Frequently asked questions. |
| ⬜ | `glossary` | Definitions of functional-programming terms. |
| ⬜ | `design-kit` | Content/branding design kit. |
| ⬜ | `api-reference` | API reference landing page. |

## Concepts

| Status | Page | Description |
| --- | --- | --- |
| ⬜ | `concepts/index` | Core concepts overview. |
| ⬜ | `concepts/pure-functions` | Pure functions and referential transparency. |
| ⬜ | `concepts/monads` | What monads are and how they're used in fp-go. |
| ⬜ | `concepts/composition` | Function composition with `Pipe` / `Flow`. |
| ⬜ | `concepts/effects-and-io` | Modeling side effects with the IO family. |
| ⬜ | `concepts/higher-kinded-types` | Higher-kinded types and how fp-go emulates them in Go. |
| ⬜ | `concepts/zen-of-go` | Reconciling FP with idiomatic ("Zen of") Go. |

## Migration

| Status | Page | Description |
| --- | --- | --- |
| ⬜ | `migration/index` | Migration guide overview. |
| ⬜ | `migration/v1-to-v2` | Upgrading from fp-go v1 to v2. |
| ⬜ | `migration/interop` | Interop between fp-go and existing/idiomatic Go code. |

---

# API Reference

## Core Effect & Data Types

| Status | Package | Description |
| --- | --- | --- |
| ⬜ | `either` | `Either[E, A]` — a value that is one of two types; the foundation for typed error handling (`Left` = failure `E`, `Right` = success `A`). |
| ⬜ | `result` | `Result[A]` = `Either[error, A]` — Go-idiomatic error handling: a success value or a standard `error`. |
| ⬜ | `option` | `Option[A]` — an optional value (`Some`/`None`); a type-safe replacement for `nil` pointers. |
| ⬜ | `io` | `IO[A]` — a lazy, synchronous computation that performs side effects and always succeeds. |
| ⬜ | `ioresult` | `IOResult[A]` = `IO[Result[A]]` — a lazy effect that may fail with an `error`. The idiomatic effect type. |
| ⬜ | `ioeither` | `IOEither[E, A]` = `IO[Either[E, A]]` — a lazy effect that may fail with a typed error `E`. |
| ⬜ | `iooption` | `IOOption[A]` = `IO[Option[A]]` — a lazy effect that produces an optional value. |
| ⬜ | `ioref` | `IORef[A]` — a mutable reference cell read and written within the `IO` monad. |
| ⬜ | `effect` | A higher-level functional effect system for managing and composing side effects. |

## Reader Types

| Status | Package | Description |
| --- | --- | --- |
| ⬜ | `reader` | `Reader[R, A]` — a computation that depends on a shared environment/config `R`. |
| ⬜ | `readereither` | `ReaderEither[R, E, A]` — a `Reader` that returns an `Either` (env + typed failure). |
| ⬜ | `readerio` | `ReaderIO[R, A]` — a `Reader` that returns an `IO` effect. |
| ⬜ | `readerioeither` | `ReaderIOEither[R, E, A]` — env + IO + typed error; the full DI effect with typed failures. |
| ⬜ | `readerioresult` | `ReaderIOResult[R, A]` — env + IO + `Result`; the idiomatic dependency-injected effect. |
| ⬜ | `readeroption` | `ReaderOption[R, A]` — a `Reader` that returns an `Option`. |
| ⬜ | `readeriooption` | `ReaderIOOption[R, A]` — env + IO producing an optional value. |
| ⬜ | `readerresult` | `ReaderResult[R, A]` — a `Reader` that returns a `Result`. |
| ⬜ | `readerreaderioeither` | A `Reader` nested over `ReaderIOEither` — composes two distinct environments. |

## State & Advanced Types

| Status | Package | Description |
| --- | --- | --- |
| ⬜ | `state` | `State[S, A]` — threads mutable state `S` through a computation in a pure way. |
| ⬜ | `stateio` | `StateIO[S, A]` — combines stateful computation with `IO` side effects. |
| ⬜ | `statereaderioeither` | Combines `State` + `Reader` + `IO` + `Either` in a single monad. |
| ⬜ | `lazy` | `Lazy[A]` — a deferred, memoizable computation. |
| ⬜ | `constant` | `Const[E, A]` — the constant functor; carries an `E` and ignores the `A`. |
| ⬜ | `identity` | `Identity[A]` — the trivial wrapper / identity monad. |
| ⬜ | `endomorphism` | `Endomorphism[A]` — functions `A -> A`, forming a monoid under composition. |
| ⬜ | `tailrec` | A trampoline for stack-safe tail-recursion (Go has no TCO). |

## Collections

| Status | Package | Description |
| --- | --- | --- |
| ⬜ | `array` | Immutable, functional operations over slices (`Map`, `Filter`, `Reduce`, `Sort`, `Uniq`, `Zip`, …). |
| ⬜ | `record` | Functional operations over `map[K]V` (`Map`, `Chain`, `Traverse`, conversions, instances). |
| ⬜ | `iterator` | Lazy sequence/iterator abstractions (`iter`, `itereither`, `iterresult`, `stateless`). |

## Utilities

| Status | Package | Description |
| --- | --- | --- |
| ⬜ | `function` | Core combinators — `Pipe`, `Flow`, `Identity`, `Constant`, currying, etc. |
| ⬜ | `predicate` | `Predicate[A]` = `func(A) bool` with `And` / `Or` / `Not` combinators. |
| ⬜ | `boolean` | Boolean helpers and folds. |
| ⬜ | `number` | Numeric helpers and type-class instances. |
| ⬜ | `string` | Functional string helpers. |
| ⬜ | `tuple` | `Tuple` types and helpers for fixed-size products. |
| ⬜ | `pair` | `Pair[A, B]` — a two-element product type. |
| ⬜ | `eq` | `Eq[A]` — the equality type class. |
| ⬜ | `ord` | `Ord[A]` — the ordering/comparison type class. |
| ⬜ | `bounded` | `Bounded[A]` — a type class providing `Min`/`Max` bounds. |
| ⬜ | `semigroup` | `Semigroup[A]` — an associative `Concat` operation. |
| ⬜ | `monoid` | `Monoid[A]` — a `Semigroup` with an identity element. |
| ⬜ | `magma` | `Magma[A]` — a binary combine with no laws (weakest algebraic structure). |
| ⬜ | `consumer` | `Consumer[A]` = `func(A)` — functions that consume a value without returning a result. |

## Type-Class / Abstractions

| Status | Package | Description |
| --- | --- | --- |
| ⬜ | `optics` | Composable `Lens` / `Prism` / `Iso` / `Optional` for accessing & updating nested immutable data. |

## Application & I/O

| Status | Package | Description |
| --- | --- | --- |
| ⬜ | `context` | Integration with `context.Context` for cancellation/deadlines inside effects. |
| ⬜ | `http` | Functional HTTP client helpers built on the effect types. |
| ⬜ | `file` | File-system operations modeled as effects. |
| ⬜ | `exec` | Running external commands as effects. |
| ⬜ | `json` | Functional JSON encode/decode helpers. |
| ⬜ | `errors` | Error construction, wrapping, and combinators. |
| ⬜ | `logging` | Logging helpers/integration for effects. |
| ⬜ | `retry` | Retry policies and combinators for fallible effects. |
| ⬜ | `circuitbreaker` | Circuit-breaker error types and utilities for effects. |
| ⬜ | `di` | Functions and types supporting dependency-injection patterns. |
| ⬜ | `builder` | A generic Builder-pattern interface for constructing values. |
| ⬜ | `bytes` | Functional operations over byte slices. |

---

# Recipes

| Status | Page | Description |
| --- | --- | --- |
| ⬜ | `recipes-index` / `recipes/index` | Recipes landing page. |
| ⬜ | `recipes/validation` | Validating input and accumulating errors. |
| ⬜ | `recipes/error-recovery` | Recovering from failures and providing fallbacks. |
| ⬜ | `recipes/error-handling` | End-to-end error-handling patterns. |
| ⬜ | `recipes/retry` | Retrying fallible operations with policies. |
| ⬜ | `recipes/data-transformation` | Transforming data through pipelines. |
| ⬜ | `recipes/filtering-mapping` | Filtering and mapping collections. |
| ⬜ | `recipes/aggregation` | Aggregating/reducing collections. |
| ⬜ | `recipes/parsing` | Parsing input functionally. |
| ⬜ | `recipes/file-operations` | Working with files as effects. |
| ⬜ | `recipes/http-requests` | Making HTTP requests functionally. |
| ⬜ | `recipes/parallel-tasks` | Running tasks concurrently/in parallel. |
| ⬜ | `recipes/dependency-injection` | Dependency injection with Reader-style effects. |
| ⬜ | `recipes/pipelines` | Building composable processing pipelines. |
| ⬜ | `recipes/middleware` | Middleware-style composition. |
| ⬜ | `recipes/testing-pure` | Testing pure functions. |
| ⬜ | `recipes/testing-effects` | Testing effectful code. |

---

# Advanced

| Status | Page | Description |
| --- | --- | --- |
| ⬜ | `advanced/patterns` | Advanced composition patterns. |
| ⬜ | `advanced/type-theory` | The type theory behind fp-go. |
| ⬜ | `advanced/performance` | Performance optimization and trade-offs. |
| ⬜ | `advanced/architecture` | Architecture patterns for FP-style Go apps. |

---

## How to add a doc page (per-PR workflow)

1. Pick a ⬜ row from this checklist.
2. Create the page in the matching folder (use `.mdx` so you can use the design
   kit — see below). API pages go in `docs/v2/<package>.mdx`.
3. Embed runnable code with **embedmd** instead of pasting snippets (see below).
4. Register the page in `sidebars.ts` under the right category.
5. Run `npm run preprocess` then `npm start` and verify the page renders.
6. Flip the row in this file to ✅.
7. Open the PR.

---

# Documentation Contribution Guidelines

Every new page should follow the two house rules below: **build it from the
design kit**, and **embed real code instead of pasting it**.

## 1. Use the design kit (`design-kit.mdx`)

The site ships a set of reusable MDX components (a "Carbon look" kit). They are
registered globally via `src/theme/MDXComponents`, so **no imports are needed** —
but the file must be `.mdx`, not `.md`, to use them.

See `docs/design-kit.mdx` for the live gallery and copy-paste examples. The
available building blocks:

| Component | Use it for |
| --- | --- |
| `<PageHeader>` | The hero header at the top of every page (eyebrow, title, lede, meta). |
| `<TLDR>` / `<TLDRCard>` | A "too long; didn't read" summary strip near the top. |
| `<Section>` | A numbered top-level section wrapper. |
| `<CodeCard>` | A styled code block with a filename + status badge. |
| `<Callout>` | Notes / `type="warn"` / `type="success"` admonitions. |
| `<Compare>` / `<CompareCol>` | Before-and-after "bad path vs recommended" strips. |
| `<Bench>` | Benchmark comparison tables with bars. |
| `<ApiTable>` | API reference tables (symbol, signature, description, since). |
| `<Checklist>` | Action-item checklists with impact tags. |
| Pager | Next/previous navigation. |

**Page skeleton** — start every page like this:

```mdx
---
id: my-package
title: My Package
sidebar_label: My Package
description: One-line summary used for SEO and previews.
---

<PageHeader
  eyebrow="API · Core types"
  title="My"
  titleAccent="package."
  lede="One-sentence description of what this package does."
  meta={[
    {label: '// Version', value: <>v2 <MetaPill>LATEST</MetaPill></>},
    {label: '// Reading time', value: '5 min'},
  ]}
/>

<TLDR>
  <TLDRCard label="// Use when" prose value={<>You need <em>…</em>.</>} />
</TLDR>

<Section id="overview" number="01" title="Overview" titleAccent=".">
  ...
</Section>
```

## 2. Embed code with embedmd (don't paste snippets)

Code examples must come from **real, compiling Go files** (the `v2/samples/`
tree) so they never drift from the API. We use
[`embedmd`](https://github.com/campoy/embedmd), which is wired into the build:

- `package.json` → `"preprocess": "embedmd -w docs/**/*.md"`
- `"prebuild": "npm run preprocess"` runs it automatically before every build,
  and CI re-runs it to verify the embeds are in sync.

**How to embed:** put an embedmd directive on its own line, immediately above an
(empty) fenced code block. Running the preprocessor fills the block in.

Whole file:

```text
[embedmd]:# (../../../v2/samples/option/basic.go go)
```

A specific range, matched by regex (start marker / end marker):

```text
[embedmd]:# (../../../v2/samples/option/basic.go go /func Example/ /^}/)
```

Rules of thumb:

- The path is **relative to the markdown file**, and must point to a file that
  actually compiles (add a sample under `v2/samples/<topic>/` if one doesn't
  exist, ideally with a `_test.go` so `go test ./...` covers it).
- Use the second arg (`go`) for syntax highlighting.
- Prefer `/start/ /end/` ranges over whole files so examples stay focused.
- After editing, run `npm run preprocess` so the generated code block is updated;
  commit the result. CI fails if the committed output is stale.
- For decorative/illustrative snippets that aren't meant to compile, a plain
  fenced block (or `<CodeCard>`) is fine — reserve embedmd for real examples.

> ⚠️ **`.md` vs `.mdx` gotcha:** the preprocess glob is `docs/**/*.md`, which
> does **not** match `.mdx` files. So today a page can use the design kit
> (`.mdx`) **or** be auto-embedded by embedmd (`.md`), not both. If you want
> both on one page, update the script to also cover `.mdx`
> (e.g. `embedmd -w docs/**/*.md docs/**/*.mdx`) in the same PR.

---

_Last updated: 2026-06-08_
