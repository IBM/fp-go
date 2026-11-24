---
sidebar_position: 10
title: Core Concepts
hide_title: true
description: Essential functional programming concepts and how they apply to Go with fp-go.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Concepts · Overview"
  title="Core"
  titleAccent="concepts."
  lede="The functional programming ideas behind fp-go — six topics, taught from first principles with Go-specific examples."
  meta={[
    {label: '// Topics', value: '6'},
    {label: '// FP experience', value: 'Not required'},
    {label: '// Reading time', value: '~45 min total'},
  ]}
/>

<TLDR>
  <TLDRCard label="// Foundations" prose value={<><em>Pure functions</em> and composition.</>} />
  <TLDRCard label="// Patterns" prose value={<>Monads, effects, <em>HKTs</em>.</>} />
  <TLDRCard label="// Integration" prose value={<>How FP fits the <em>Zen of Go.</em></>} />
</TLDR>

<Section id="topics" number="01" title="The six" titleAccent="topics.">

<ApiTable
  columns={['Topic', 'You learn', 'Why it matters']}
  rows={[
    {symbol: <a href="./pure-functions">Pure functions</a>, signature: 'Definition, properties, benefits, pure vs impure in Go, practical examples', description: 'The foundation. Easier to test, reason about, compose.'},
    {symbol: <a href="./monads">Monads</a>, signature: 'What monads are (and aren\'t), the monad laws, Option/Either/Result/IO, usage patterns', description: 'A consistent way to handle effects, errors, optional values — without exceptions or nil checks.'},
    {symbol: <a href="./composition">Composition</a>, signature: 'Function composition, Pipe and Flow, point-free style, pipelines', description: 'Build maintainable systems from small, focused functions.'},
    {symbol: <a href="./effects-and-io">Effects & IO</a>, signature: 'What effects are, lazy evaluation, IO monad, Effect type, separating description from execution', description: 'Explicit effect management makes code predictable and testable.'},
    {symbol: <a href="./higher-kinded-types">Higher-kinded types</a>, signature: 'What HKTs are, why Go lacks them, how fp-go works around it, generic parameters', description: 'Helps you use fp-go\'s generic APIs effectively.'},
    {symbol: <a href="./zen-of-go">The Zen of Go</a>, signature: 'Simplicity, explicit over implicit, composition, error handling, when to use fp-go', description: 'fp-go enhances Go — it doesn\'t replace it.'},
  ]}
/>

</Section>

<Section id="paths" number="02" title="Learning" titleAccent="paths.">

<Compare>
  <CompareCol kind="good" title="For beginners" pill="recommended order">
    <p>1. <a href="./pure-functions">Pure functions</a> — the basics.</p>
    <p>2. <a href="./effects-and-io">Effects & IO</a> — why special types are needed.</p>
    <p>3. <a href="./composition">Composition</a> — building pipelines.</p>
    <p>4. <a href="./monads">Monads</a> — the underlying pattern.</p>
    <p>5. <a href="./zen-of-go">The Zen of Go</a> — integrating with Go idioms.</p>
  </CompareCol>
  <CompareCol kind="good" title="For FP practitioners" pill="from other languages">
    <p>1. <a href="./zen-of-go">Zen of Go</a> — understand the constraints.</p>
    <p>2. <a href="./higher-kinded-types">HKTs</a> — see fp-go's approach.</p>
    <p>3. <a href="./monads">Monads</a> — Go-specific implementations.</p>
    <p>4. <a href="./effects-and-io">Effects & IO</a> — lazy evaluation in Go.</p>
  </CompareCol>
</Compare>

<Callout title="For Go developers new to FP.">
  Start with <a href="./zen-of-go">The Zen of Go</a> → <a href="./pure-functions">Pure functions</a> → <a href="./composition">Composition</a> → <a href="./effects-and-io">Effects & IO</a> → <a href="./monads">Monads</a>.
</Callout>

</Section>

<Section id="examples" number="03" title="Practical" titleAccent="examples." lede="Each concept page includes a three-tab comparison: without fp-go, fp-go v2, fp-go v1.">

<Tabs groupId="approach">
<TabItem value="standard" label="Without fp-go">

<CodeCard file="processUser.go">
{`// Standard Go approach
func processUser(id string) (*User, error) {
    user, err := fetchUser(id)
    if err != nil {
        return nil, err
    }

    validated, err := validateUser(user)
    if err != nil {
        return nil, err
    }

    return saveUser(validated)
}`}
</CodeCard>

</TabItem>
<TabItem value="v2" label="With fp-go v2">

<CodeCard file="processUser.go" status="tested">
{`// fp-go v2 approach
func processUser(id string) ioresult.IOResult[User] {
    return function.Pipe3(
        fetchUser(id),
        ioresult.Chain(validateUser),
        ioresult.Chain(saveUser),
    )
}`}
</CodeCard>

</TabItem>
<TabItem value="v1" label="With fp-go v1">

<CodeCard file="processUser.go">
{`// fp-go v1 approach
func processUser(id string) ioeither.IOEither[error, User] {
    return function.Pipe3(
        fetchUser(id),
        ioeither.Chain(validateUser),
        ioeither.Chain(saveUser),
    )
}`}
</CodeCard>

</TabItem>
</Tabs>

<Callout title="Real-world scenarios covered.">
  HTTP API handlers · database operations · file processing · configuration management · error-handling patterns.
</Callout>

</Section>

<Section id="faq" number="04" title="Common" titleAccent="questions.">

<Callout title="Isn't this overengineering?">
  Not if used appropriately. fp-go shines for complex error handling, multiple sequential operations, composability needs, and testing-heavy codebases. For simple cases, standard Go is often better — see <a href="./zen-of-go">The Zen of Go</a>.
</Callout>

<Callout title="Will this make my code slower?">
  Usually no, sometimes yes. Pure functions are often faster (easier to optimize); composition has minimal overhead; IO types add one function call. Use idiomatic packages for performance-critical code.
</Callout>

<Callout title="How do I convince my team?">
  Start small, show value: use fp-go for new features, demonstrate improved testability and clearer error handling, share success stories, provide training. See <a href="../why-fp-go">Why fp-go?</a>.
</Callout>

<Callout type="success" title="What if I don't understand monads?">
  You don't need to. Think: <strong>Option</strong> = "maybe has a value"; <strong>Result</strong> = "success or error"; <strong>IO</strong> = "lazy computation". The theory helps but isn't required.
</Callout>

</Section>

<Section id="relationships" number="05" title="Concept" titleAccent="relationships.">

<CodeCard file="map" copy={false}>
{`Pure Functions
    ↓
Composition ←→ Effects & IO
    ↓              ↓
Monads ←──────────┘
    ↓
Higher-Kinded Types
    ↓
The Zen of Go`}
</CodeCard>

<Callout title="Flow.">
  Pure functions are composable → composition builds pipelines → effects need special handling → monads provide the pattern → HKTs enable generic implementations → the Zen of Go guides usage.
</Callout>

</Section>

<Section id="prereqs" number="06" title="Prerequisites">

<Compare>
  <CompareCol kind="good" title="Required" pill="Go basics">
    <p>Basic Go syntax.</p>
    <p>Functions and closures.</p>
    <p>Interfaces.</p>
    <p>Generics (Go 1.18+).</p>
  </CompareCol>
  <CompareCol kind="good" title="Recommended" pill="reading">
    <p><a href="https://go.dev/tour/">Go Tour</a> — if new to Go.</p>
    <p><a href="https://go.dev/doc/effective_go">Effective Go</a> — best practices.</p>
    <p><a href="../quickstart">Quickstart</a> — fp-go basics.</p>
  </CompareCol>
</Compare>

<Callout type="success" title="No FP experience needed.">
  These guides assume no prior functional programming knowledge. Concepts are taught from first principles with Go examples.
</Callout>

</Section>

<Section id="summary" number="07" title="Summary">

<Checklist
  title="What you'll learn"
  items={[
    {label: 'Pure Functions — predictable, testable code', done: true},
    {label: 'Monads — pattern for sequencing with context', done: true},
    {label: 'Composition — building from simple pieces', done: true},
    {label: 'Effects & IO — explicit side effect management', done: true},
    {label: 'Higher-Kinded Types — generic abstractions', done: true},
    {label: 'The Zen of Go — idiomatic usage', done: true},
  ]}
/>

<Callout type="success" title="Key takeaway.">
  You don't need to understand all the theory to use fp-go effectively. Start with practical patterns and deepen your understanding over time.
</Callout>

</Section>
