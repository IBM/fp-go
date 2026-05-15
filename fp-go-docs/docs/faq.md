---
sidebar_position: 5
title: Frequently Asked Questions
hide_title: true
description: Common questions about fp-go — versions, performance, error handling, patterns, migration, and contributing.
---

<PageHeader
  eyebrow="Reference · FAQ"
  title="Frequently asked"
  titleAccent="questions."
  lede="Common questions about fp-go, functional programming in Go, and when to use these patterns."
  meta={[
    {label: '// Version', value: <>v2.2.82 <MetaPill>LATEST</MetaPill></>},
    {label: '// Updated', value: 'today'},
    {label: '// Reading time', value: '10 min · 9 sections'},
  ]}
/>

<Section
  id="getting-started"
  number="01"
  title="Getting"
  titleAccent="started."
>

<Callout title="What is fp-go?">
  fp-go is a comprehensive functional programming library for Go providing monadic types (Result, Either, Option, IO, Reader, …) and utilities for composing operations, handling errors, and managing side effects in a type-safe, functional way.
  <ul>
    <li>Type-safe error handling with Result/Either</li>
    <li>Automatic error propagation through pipelines</li>
    <li>IO monad for managing side effects</li>
    <li>Reader monad for dependency injection</li>
    <li>Comprehensive collection operations</li>
    <li>Full composition utilities (Pipe, Flow, Compose)</li>
  </ul>
</Callout>

<Callout title="Do I need to know FP to use fp-go?">
  <p><strong>Short answer:</strong> not initially, but learning FP concepts will help you use it effectively.</p>
  <p>Start with Result/Either for error handling. Learn Pipe for composition. Gradually explore Map, Chain, and the rest.</p>
  <p>Resources: <a href="./quickstart">5-Minute Quickstart</a>, <a href="./concepts">Core Concepts</a>, <a href="./concepts/pure-functions">Pure Functions</a>, <a href="./concepts/monads">Monads Explained</a>.</p>
</Callout>

<Callout title="Which version: v1 or v2?">
  <Compare>
    <CompareCol kind="good" title="Use v2" pill="recommended">
      <p>You're on Go 1.24+.</p>
      <p>Starting a new project.</p>
      <p>Want latest features (Result, Effect, idiomatic packages).</p>
      <p>Want better type inference.</p>
    </CompareCol>
    <CompareCol kind="bad" title="Use v1" pill="legacy">
      <p>Stuck on Go 1.18–1.23.</p>
      <p>Have existing v1 codebase.</p>
      <p>Need Writer monad (v1 only).</p>
    </CompareCol>
  </Compare>
  <p><strong>Recommendation:</strong> use v2 for all new projects. See the <a href="./migration/v1-to-v2">migration guide</a> for upgrading.</p>
</Callout>

<Callout type="success" title="Is fp-go production-ready?">
  Yes. fp-go is used in production at IBM (creator and maintainer) and other companies. v2 is stable, actively maintained, and recommended. v1 is stable and in maintenance mode.
</Callout>

</Section>

<Section
  id="performance"
  number="02"
  title="Performance"
>

<Callout title="What's the performance impact of fp-go?">
  <p><strong>Standard packages:</strong> small overhead from function calls. Typically 5–15% slower than hand-written code. Negligible for most applications. Worth it for type safety and maintainability.</p>
  <p><strong>Idiomatic packages (v2 only):</strong> 2–32× faster than standard packages. Near-native performance. Use native Go tuples instead of generic types. Recommended for performance-critical code.</p>
</Callout>

<Bench
  title="Relative cost"
  command="go test -bench=. -benchmem"
  rows={[
    {label: 'Filter — stdlib', bar: 0.86, barKind: 'win', nsOp: '1.00x', bOp: '—', delta: 'baseline'},
    {label: 'Filter — fp-go (idiomatic)', bar: 0.89, barKind: 'win', nsOp: '1.03x', bOp: '—', delta: '+3%', winner: true},
    {label: 'Filter — fp-go (standard)', bar: 0.98, barKind: 'lose', nsOp: '1.14x', bOp: '—', delta: '+14%', deltaKind: 'bad'},
    {label: 'Map — fp-go (idiomatic)', bar: 0.88, barKind: 'win', nsOp: '1.02x', bOp: '—', delta: '+2%'},
    {label: 'Map — fp-go (standard)', bar: 0.96, barKind: 'lose', nsOp: '1.12x', bOp: '—', delta: '+12%', deltaKind: 'bad'},
    {label: 'Reduce — fp-go (idiomatic)', bar: 0.90, barKind: 'win', nsOp: '1.04x', bOp: '—', delta: '+4%'},
    {label: 'Reduce — fp-go (standard)', bar: 0.99, barKind: 'lose', nsOp: '1.15x', bOp: '—', delta: '+15%', deltaKind: 'bad'},
  ]}
/>

<Callout type="success" title="Bottom line.">
  Fast enough for 99% of use cases. Use idiomatic packages for hot paths.
</Callout>

<Callout title="When should I use idiomatic packages?">
  <Compare>
    <CompareCol kind="good" title="Idiomatic" pill="fast path">
      <p>Performance is critical.</p>
      <p>Processing large datasets.</p>
      <p>In tight loops.</p>
      <p>Hot code paths.</p>
    </CompareCol>
    <CompareCol kind="bad" title="Standard" pill="default">
      <p>Type safety matters most.</p>
      <p>Performance is adequate.</p>
      <p>Code clarity matters most.</p>
      <p>Not performance-critical.</p>
    </CompareCol>
  </Compare>
  <CodeCard file="example.go">
{`// Standard - type-safe, slightly slower
result := array.Map(func(x int) int { return x * 2 })(data)

// Idiomatic - near-native speed
result := idiomatic.Map(data, func(x int) int { return x * 2 })`}
  </CodeCard>
  <p>See the <a href="./advanced/performance">performance guide</a> for details.</p>
</Callout>

<Callout title="Does fp-go allocate a lot of memory?">
  No more than idiomatic Go. Result/Either/Option are single allocations per value. Pipelines have intermediate allocations (same as manual code). Idiomatic packages have minimal allocations. For memory-efficient patterns: use iterators for lazy evaluation, use idiomatic packages for large datasets, avoid unnecessary intermediate collections.
</Callout>

</Section>

<Section
  id="error-handling"
  number="03"
  title="Error"
  titleAccent="handling."
>

<Callout title="Should I use Result or Either?">
  <Compare>
    <CompareCol kind="good" title="Use Result (v2)" pill="recommended">
      <p>For error handling.</p>
      <p>Error type is always <code>error</code>.</p>
      <p>Simpler API, more idiomatic for Go.</p>
    </CompareCol>
    <CompareCol kind="bad" title="Use Either" pill="when needed">
      <p>You need non-error Left values.</p>
      <p>Sum types beyond error handling.</p>
      <p>Porting from other FP languages.</p>
    </CompareCol>
  </Compare>
  <CodeCard file="example.go">
{`// Result - recommended for errors
func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}

// Either - for generic sum types
func parseValue(s string) either.Either[ParseError, Value] {
    // Left can be any type, not just error
}`}
  </CodeCard>
</Callout>

<Callout title="How do I handle multiple errors?">
  <p><strong>Option 1: Stop at first error (default).</strong></p>
  <CodeCard file="option1.go">
{`result := function.Pipe3(
    step1(),
    result.Chain(step2),
    result.Chain(step3),
)
// Stops at first error`}
  </CodeCard>
  <p><strong>Option 2: Accumulate errors.</strong></p>
  <CodeCard file="option2.go">
{`results := array.TraverseResult(validate)(items)
// Returns Result[[]Item] - all or nothing

// For accumulating all errors, use Validation applicative
// (Advanced pattern - see recipes)`}
  </CodeCard>
  <p><strong>Option 3: Collect errors manually.</strong></p>
  <CodeCard file="option3.go">
{`var errors []error
for _, item := range items {
    if err := validate(item); err != nil {
        errors = append(errors, err)
    }
}`}
  </CodeCard>
</Callout>

<Callout title="Can I convert between Result and (value, error)?">
  <p>Yes — use the conversion helpers.</p>
  <CodeCard file="interop.go">
{`// From (value, error) to Result
func fetchUser(id string) result.Result[User] {
    user, err := db.Query(id)
    return result.FromGoError(user, err)
}

// From Result to (value, error)
func legacyAPI() (User, error) {
    result := fetchUser("123")
    return result.ToGoError()
}`}
  </CodeCard>
</Callout>

</Section>

<Section
  id="patterns"
  number="04"
  title="Usage"
  titleAccent="patterns."
>

<Callout title="When should I use fp-go vs. idiomatic Go?">
  <Compare>
    <CompareCol kind="good" title="Use fp-go for" pill="fits">
      <p>Complex error handling logic.</p>
      <p>Data transformation pipelines.</p>
      <p>Composable business logic.</p>
      <p>Side effect management.</p>
      <p>Dependency injection patterns.</p>
    </CompareCol>
    <CompareCol kind="bad" title="Use idiomatic Go for" pill="fits">
      <p>Simple CRUD operations.</p>
      <p>Straightforward logic.</p>
      <p>Performance-critical hot paths.</p>
      <p>When team is unfamiliar with FP.</p>
    </CompareCol>
  </Compare>
  <CodeCard file="fp-go-example.go">
{`// Complex pipeline with error handling
func processOrder(order Order) result.Result[Receipt] {
    return function.Pipe5(
        validateOrder(order),
        result.Chain(checkInventory),
        result.Chain(calculatePrice),
        result.Chain(applyDiscounts),
        result.Chain(generateReceipt),
    )
}`}
  </CodeCard>
  <CodeCard file="idiomatic-example.go">
{`// Simple database query
func getUser(id string) (*User, error) {
    return db.QueryUser(id)
}`}
  </CodeCard>
</Callout>

<Callout type="success" title="Can I mix fp-go with idiomatic Go?">
  <p>Yes — they work together seamlessly.</p>
  <CodeCard file="mix.go">
{`// Idiomatic Go function
func fetchFromDB(id string) (*User, error) {
    return db.Query(id)
}

// Wrap in Result for fp-go pipeline
func getUser(id string) result.Result[*User] {
    user, err := fetchFromDB(id)
    return result.FromGoError(user, err)
}

// Use in pipeline
result := function.Pipe2(
    getUser("123"),
    result.Map(enrichUser),
    result.Chain(validateUser),
)`}
  </CodeCard>
</Callout>

<Callout title="How do I handle side effects?">
  <p>Use IO types.</p>
  <CodeCard file="io.go">
{`// Pure function returning IO
func readFile(path string) io.IO[[]byte] {
    return func() []byte {
        data, _ := os.ReadFile(path)
        return data
    }
}

// Compose IO operations
program := function.Pipe2(
    readFile("config.json"),
    io.Map(parseConfig),
    io.Map(validateConfig),
)

// Execute when ready
config := program() // Side effect happens here`}
  </CodeCard>
  <p>Use IOResult when the effect can fail:</p>
  <CodeCard file="ioresult.go">
{`func readFile(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}`}
  </CodeCard>
  <p>See <a href="./concepts/effects-and-io">Effects and IO</a>.</p>
</Callout>

</Section>

<Section
  id="learning"
  number="05"
  title="Learning &"
  titleAccent="adoption."
>

<ApiTable
  columns={['Phase', 'You learn', 'Difficulty']}
  rows={[
    {symbol: 'Week 1', signature: 'Result/Either, Pipe, Map', description: 'Low–Medium.'},
    {symbol: 'Week 2–4', signature: 'Chain (FlatMap), IO, Reader', description: 'Medium.'},
    {symbol: 'Month 2+', signature: 'All monadic types, monad laws, optics', description: 'Medium–High.'},
  ]}
/>

<Callout title="How do I convince my team to use fp-go?">
  <ol>
    <li>Use for new features only.</li>
    <li>Show concrete benefits (fewer bugs, easier testing).</li>
    <li>Provide training and examples.</li>
    <li>Let the team see the value.</li>
  </ol>
  <p>Address concerns: show benchmarks for performance; provide training for the learning curve; start with simple patterns for complexity; do gradual migration for adoption.</p>
  <p>Success metrics: reduced bug count, faster development, easier code reviews, better test coverage.</p>
</Callout>

<Callout title="Are there good tutorials or courses?">
  <p><strong>Official:</strong> <a href="./quickstart">Quickstart</a>, <a href="./concepts">Core Concepts</a>, <a href="./recipes/error-handling">Recipes</a>, <a href="https://pkg.go.dev/github.com/IBM/fp-go/v2">API docs</a>.</p>
  <p><strong>Learning path:</strong> Read <a href="./why-fp-go">Why fp-go?</a> → Complete <a href="./quickstart">Quickstart</a> → Study <a href="./concepts/pure-functions">Pure Functions</a> → Learn <a href="./concepts/monads">Monads</a> → Practice with <a href="./recipes/error-handling">Recipes</a>.</p>
</Callout>

</Section>

<Section
  id="comparison"
  number="06"
  title="Comparison"
>

<Callout title="How does fp-go compare to samber/lo?">
  <Compare>
    <CompareCol kind="bad" title="samber/lo" pill="utility">
      <p>Simple collection operations.</p>
      <p>Low learning curve.</p>
      <p>Excellent performance.</p>
      <p>No monadic types, no error handling.</p>
    </CompareCol>
    <CompareCol kind="good" title="fp-go" pill="FP toolkit">
      <p>Full FP toolkit.</p>
      <p>Built-in error handling.</p>
      <p>Monadic composition.</p>
      <p>Steeper learning curve; idiomatic packages for speed.</p>
    </CompareCol>
  </Compare>
  <CodeCard file="both.go">
{`// Use lo for simple operations
filtered := lo.Filter(items, predicate)

// Use fp-go for error handling
result := result.TraverseArray(processWithErrors)(filtered)`}
  </CodeCard>
  <p>See the <a href="./comparison">comparison guide</a>.</p>
</Callout>

<Callout title="Is fp-go similar to fp-ts (TypeScript)?">
  <p>Yes — fp-go is heavily inspired by fp-ts. Same monadic types, similar API design, same composition patterns, same concepts. Main differences: Go's type system limits (no HKT), different syntax, performance characteristics, ecosystem.</p>
</Callout>

</Section>

<Section
  id="troubleshooting"
  number="07"
  title="Troubleshooting"
>

<Callout type="warn" title="Why am I getting type inference errors?">
  <p><strong>Common causes:</strong></p>
  <p><strong>1. Type parameters in wrong order.</strong></p>
  <CodeCard file="case1.go">
{`// Wrong - B cannot be inferred
result.Map(func(x int) string { return fmt.Sprintf("%d", x) })

// Right - specify B explicitly
result.Map[string](func(x int) string { return fmt.Sprintf("%d", x) })`}
  </CodeCard>
  <p><strong>2. Missing type parameters.</strong></p>
  <CodeCard file="case2.go">
{`// Wrong
return result.Err(errors.New("error"))

// Right
return result.Err[int](errors.New("error"))`}
  </CodeCard>
  <p><strong>3. Ambiguous types.</strong></p>
  <CodeCard file="case3.go">
{`// Wrong - compiler can't infer
result := result.Ok(nil)

// Right - specify type
result := result.Ok[*User](nil)`}
  </CodeCard>
</Callout>

<Callout title="How do I debug pipelines?">
  <p>Use intermediate logging:</p>
  <CodeCard file="debug.go">
{`result := function.Pipe3(
    step1(),
    result.Map(func(x T) T {
        fmt.Printf("After step1: %+v\\n", x)
        return x
    }),
    result.Chain(step2),
    result.Map(func(x T) T {
        fmt.Printf("After step2: %+v\\n", x)
        return x
    }),
)`}
  </CodeCard>
  <p>Use Fold to inspect:</p>
  <CodeCard file="fold.go">
{`result.Fold(
    func(err error) {
        fmt.Printf("Error: %v\\n", err)
    },
    func(val T) {
        fmt.Printf("Success: %+v\\n", val)
    },
)`}
  </CodeCard>
  <p>Use the logging package:</p>
  <CodeCard file="logging.go">
{`import "github.com/IBM/fp-go/v2/logging"

result := logging.WithLogging(
    "operation",
    func() result.Result[T] {
        return operation()
    },
)`}
  </CodeCard>
</Callout>

</Section>

<Section
  id="migration"
  number="08"
  title="Migration"
>

<Callout title="How do I migrate from v1 to v2?">
  <p><strong>5 breaking changes:</strong></p>
  <ol>
    <li><strong>Generic type aliases</strong> — <code>type X = Y</code> instead of <code>type X Y</code></li>
    <li><strong>Type parameter reordering</strong> — non-inferrable params first</li>
    <li><strong>Pair operates on second</strong> — v1 was first, v2 is second</li>
    <li><strong>Compose is right-to-left</strong> — mathematical composition</li>
    <li><strong>No <code>generic/</code> subpackages</strong> — removed</li>
  </ol>
  <p><strong>Steps:</strong> update imports → fix type-parameter order → update Pair usage → update Compose usage → remove <code>generic/</code> imports.</p>
  <p>See the <a href="./migration/v1-to-v2">complete migration guide</a>.</p>
</Callout>

<Callout type="success" title="Can I run v1 and v2 side-by-side?">
  <p>Yes — they can coexist.</p>
  <CodeCard file="side-by-side.go">
{`import (
    v1 "github.com/IBM/fp-go/either"
    v2 "github.com/IBM/fp-go/v2/result"
)

// Use both in same codebase
v1Result := v1.Right[error](42)
v2Result := v2.Ok(42)`}
  </CodeCard>
  <p>See <a href="./migration/interop">Interop during migration</a>.</p>
</Callout>

</Section>

<Section
  id="contributing"
  number="09"
  title="Contributing"
>

<Callout type="success" title="Ways to contribute.">
  <ul>
    <li>Report bugs</li>
    <li>Suggest features</li>
    <li>Improve documentation</li>
    <li>Submit pull requests</li>
    <li>Help in discussions</li>
    <li>Star the repository</li>
  </ul>
  <p>Read the <a href="https://github.com/IBM/fp-go/blob/main/CONTRIBUTING.md">contributing guide</a>; check the <a href="https://github.com/IBM/fp-go/labels/good%20first%20issue">good first issues</a>; join <a href="https://github.com/IBM/fp-go/discussions">GitHub Discussions</a>.</p>
</Callout>

<Callout title="Where can I get help?">
  <ul>
    <li><a href="https://github.com/IBM/fp-go/discussions">GitHub Discussions</a> — ask questions</li>
    <li><a href="https://github.com/IBM/fp-go/issues">GitHub Issues</a> — report bugs</li>
    <li><a href="./intro">Documentation</a> — comprehensive guides</li>
    <li><a href="https://pkg.go.dev/github.com/IBM/fp-go/v2">API reference</a></li>
    <li><a href="./recipes/error-handling">Recipes</a> — practical examples</li>
  </ul>
</Callout>

</Section>
