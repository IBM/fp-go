---
sidebar_position: 4
title: Comparison with Other Libraries
hide_title: true
description: How fp-go compares to samber/lo, go-functional, and mo — feature matrix, side-by-side code, and use-case recommendations.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Getting started · Section 05 / 05"
  title="Compare with"
  titleAccent="other libraries."
  lede="How fp-go stacks up against samber/lo, go-functional, and mo — features, code, and the right tool for the job."
  meta={[
    {label: '// Version', value: <>v2.2.82 <MetaPill>LATEST</MetaPill></>},
    {label: '// Updated', value: 'today'},
    {label: '// Reading time', value: '8 min · 5 sections'},
  ]}
/>

<TLDR>
  <TLDRCard label="// Libraries compared" value="4" description="fp-go · samber/lo · go-functional · mo." />
  <TLDRCard label="// Best for full FP" prose value={<><em>fp-go.</em> 20+ monadic types.</>} variant="up" />
  <TLDRCard label="// Best for simple ops" prose value={<><em>samber/lo.</em> Low learning curve.</>} />
</TLDR>

<Section
  id="matrix"
  number="01"
  title="At-a-glance"
  titleAccent="matrix."
>

<ApiTable
  columns={['Feature', 'fp-go', 'Others']}
  rows={[
    {symbol: 'Monadic types', signature: 'Full (Result, Either, Option, IO, …)', description: 'samber/lo: none · go-functional: limited · mo: some.'},
    {symbol: 'Type safety', signature: 'Excellent (generics)', description: 'samber/lo: uses any · go-functional: good · mo: good.'},
    {symbol: 'Error handling', signature: 'Built-in (Result, Either)', description: 'samber/lo: manual · go-functional: built-in · mo: built-in.'},
    {symbol: 'Lazy evaluation', signature: 'IO, Lazy types', description: 'samber/lo: no · go-functional: yes · mo: no.'},
    {symbol: 'Composition', signature: 'Pipe, Flow, Compose', description: 'samber/lo: limited · go-functional: yes · mo: limited.'},
    {symbol: 'Documentation', signature: 'Comprehensive', description: 'samber/lo: good · go-functional: limited · mo: good.'},
    {symbol: 'Learning curve', signature: 'Medium (FP concepts)', description: 'samber/lo: low · go-functional: medium · mo: low.'},
    {symbol: 'Performance', signature: 'Good (idiomatic: excellent)', description: 'samber/lo: excellent · go-functional: unknown · mo: good.'},
    {symbol: 'Production use', signature: 'IBM, others', description: 'samber/lo: wide · go-functional: unknown · mo: growing.'},
    {symbol: 'Active development', signature: 'Yes', description: 'samber/lo: yes · go-functional: sporadic · mo: yes.'},
    {symbol: 'Go version', signature: 'v2: 1.24+, v1: 1.18+', description: 'All others: 1.18+.'},
  ]}
/>

</Section>

<Section
  id="vs-lo"
  number="02"
  title="fp-go vs."
  titleAccent="samber/lo."
  tag="Difficulty · Beginner"
  lede="samber/lo is a Lodash-inspired utility library focused on collection operations."
>

### Collection operations

<Tabs>
  <TabItem value="lo" label="samber/lo" default>

<CodeCard file="lo.go">
{`import "github.com/samber/lo"

// Filter and map
result := lo.Map(
    lo.Filter(numbers, func(n int, _ int) bool {
        return n > 0
    }),
    func(n int, _ int) int {
        return n * 2
    },
)

// No built-in error handling
// Must handle errors manually`}
</CodeCard>

<Compare>
  <CompareCol kind="good" pill="pros">
    <p>Simple, intuitive API.</p>
    <p>Excellent performance.</p>
    <p>Wide adoption.</p>
    <p>Low learning curve.</p>
  </CompareCol>
  <CompareCol kind="bad" pill="cons">
    <p>No monadic types.</p>
    <p>No built-in error handling.</p>
    <p>Limited composition.</p>
    <p>Uses <code>any</code> for some operations.</p>
  </CompareCol>
</Compare>

  </TabItem>
  <TabItem value="fp-go" label="fp-go">

<CodeCard file="fp-go.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/array"
    "github.com/IBM/fp-go/v2/function"
    "github.com/IBM/fp-go/v2/result"
)

// Filter and map with error handling
result := function.Pipe2(
    array.Filter(func(n int) bool { return n > 0 }),
    array.TraverseResult(func(n int) result.Result[int] {
        if n > 100 {
            return result.Err[int](errors.New("too large"))
        }
        return result.Ok(n * 2)
    }),
)(numbers)

// Automatic error propagation`}
</CodeCard>

<Compare>
  <CompareCol kind="good" pill="pros">
    <p>Full monadic types.</p>
    <p>Built-in error handling.</p>
    <p>Excellent composition.</p>
    <p>Type-safe throughout.</p>
  </CompareCol>
  <CompareCol kind="bad" pill="cons">
    <p>Steeper learning curve.</p>
    <p>More verbose for simple operations.</p>
    <p>Requires FP knowledge.</p>
  </CompareCol>
</Compare>

  </TabItem>
</Tabs>

### Error handling

<Tabs>
  <TabItem value="lo" label="samber/lo">

<CodeCard file="lo-errs.go">
{`// Manual error handling required
results := make([]Result, 0)
for _, item := range items {
    result, err := process(item)
    if err != nil {
        return nil, err
    }
    results = append(results, result)
}`}
</CodeCard>

  </TabItem>
  <TabItem value="fp-go" label="fp-go">

<CodeCard file="fp-go-errs.go" status="tested">
{`// Automatic error handling
results := array.TraverseResult(process)(items)
// Returns Result[[]Result] - single error handling point`}
</CodeCard>

  </TabItem>
</Tabs>

<Compare>
  <CompareCol kind="bad" title="Choose samber/lo when" pill="best fit">
    <p>Simple collection operations.</p>
    <p>No complex error handling needed.</p>
    <p>Team unfamiliar with FP.</p>
    <p>Performance is critical.</p>
  </CompareCol>
  <CompareCol kind="good" title="Choose fp-go when" pill="best fit">
    <p>Complex error handling.</p>
    <p>Need monadic composition.</p>
    <p>Building data pipelines.</p>
    <p>Type safety is paramount.</p>
  </CompareCol>
</Compare>

</Section>

<Section
  id="vs-go-functional"
  number="03"
  title="fp-go vs."
  titleAccent="go-functional."
  lede="go-functional is a functional library built around iterators and monadic types."
>

<Tabs>
  <TabItem value="go-functional" label="go-functional" default>

<CodeCard file="go-functional.go">
{`import "github.com/BooleanCat/go-functional/v2/it"

// Iterator-based operations
result := it.Map(
    it.Filter(
        it.Lift(numbers),
        func(n int) bool { return n > 0 },
    ),
    func(n int) int { return n * 2 },
)`}
</CodeCard>

<Compare>
  <CompareCol kind="good" pill="pros">
    <p>Monadic types.</p>
    <p>Iterator-based (lazy).</p>
    <p>Good type safety.</p>
  </CompareCol>
  <CompareCol kind="bad" pill="cons">
    <p>Limited documentation.</p>
    <p>Smaller ecosystem.</p>
    <p>Less active development.</p>
    <p>Fewer types than fp-go.</p>
  </CompareCol>
</Compare>

  </TabItem>
  <TabItem value="fp-go" label="fp-go">

<CodeCard file="fp-go.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/array"
    "github.com/IBM/fp-go/v2/function"
)

// Array-based operations
result := function.Pipe2(
    array.Filter(func(n int) bool { return n > 0 }),
    array.Map(func(n int) int { return n * 2 }),
)(numbers)

// Also supports lazy iterators`}
</CodeCard>

<Compare>
  <CompareCol kind="good" pill="pros">
    <p>More monadic types (Result, Either, IO, Reader, …).</p>
    <p>Comprehensive documentation.</p>
    <p>Active development.</p>
    <p>Production-proven.</p>
  </CompareCol>
  <CompareCol kind="bad" pill="cons">
    <p>Larger API surface.</p>
    <p>More concepts to learn.</p>
  </CompareCol>
</Compare>

  </TabItem>
</Tabs>

</Section>

<Section
  id="vs-mo"
  number="04"
  title="fp-go vs."
  titleAccent="mo."
  lede="mo is a lightweight FP library focused on Option, Result, and Either."
>

<Tabs>
  <TabItem value="mo" label="mo" default>

<CodeCard file="mo.go">
{`import "github.com/samber/mo"

// Option type
opt := mo.Some(42)
value := opt.OrElse(0)

// Result type
result := mo.Ok[int](42)
value = result.OrElse(0)`}
</CodeCard>

<Compare>
  <CompareCol kind="good" pill="pros">
    <p>Simple API.</p>
    <p>Good documentation.</p>
    <p>Easy to learn.</p>
    <p>Good performance.</p>
  </CompareCol>
  <CompareCol kind="bad" pill="cons">
    <p>Limited types (Option, Result, Either).</p>
    <p>No IO or Reader types.</p>
    <p>Limited composition utilities.</p>
    <p>No lazy evaluation.</p>
  </CompareCol>
</Compare>

  </TabItem>
  <TabItem value="fp-go" label="fp-go">

<CodeCard file="fp-go.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/result"
)

// Option type
opt := option.Some(42)
value := option.GetOrElse(func() int { return 0 })(opt)

// Result type
res := result.Ok(42)
value = result.GetOrElse(func() int { return 0 })(res)

// Plus: IO, Reader, State, and more`}
</CodeCard>

<Compare>
  <CompareCol kind="good" pill="pros">
    <p>Full FP toolkit (20+ types).</p>
    <p>IO and Reader for effects.</p>
    <p>Lazy evaluation support.</p>
    <p>Comprehensive composition.</p>
  </CompareCol>
  <CompareCol kind="bad" pill="cons">
    <p>More complex API.</p>
    <p>Steeper learning curve.</p>
    <p>Overkill for simple use cases.</p>
  </CompareCol>
</Compare>

  </TabItem>
</Tabs>

</Section>

<Section
  id="features"
  number="05"
  title="Feature"
  titleAccent="deep dive."
>

### Error handling

<ApiTable
  columns={['Library', 'Approach', 'Composition']}
  rows={[
    {symbol: 'fp-go', signature: 'Result/Either types', description: 'Type safety: excellent. Composition: excellent.'},
    {symbol: 'samber/lo', signature: 'Manual (Go style)', description: 'Type safety: good. Composition: none.'},
    {symbol: 'go-functional', signature: 'Result type', description: 'Type safety: good. Composition: good.'},
    {symbol: 'mo', signature: 'Result/Either types', description: 'Type safety: good. Composition: limited.'},
  ]}
/>

### Collection operations

<ApiTable
  columns={['Library', 'Operations', 'Notes']}
  rows={[
    {symbol: 'fp-go', signature: '30+ operations', description: 'Good performance (idiomatic: excellent). Type safety: excellent.'},
    {symbol: 'samber/lo', signature: '100+ operations', description: 'Excellent performance. Uses any.'},
    {symbol: 'go-functional', signature: '20+ operations', description: 'Performance unknown. Type safety: good.'},
    {symbol: 'mo', signature: '10+ operations', description: 'Good performance. Type safety: good.'},
  ]}
/>

### Monadic types

<ApiTable
  columns={['Library', 'Types available', '']}
  rows={[
    {symbol: 'fp-go', signature: 'Result, Either, Option, IO, IOResult, IOEither, IOOption, Reader, ReaderIO, ReaderIOResult, State, Lazy, Identity, Effect, and more', description: ''},
    {symbol: 'samber/lo', signature: 'None', description: ''},
    {symbol: 'go-functional', signature: 'Option, Result, Iterator', description: ''},
    {symbol: 'mo', signature: 'Option, Result, Either, Future, Task', description: ''},
  ]}
/>

### Composition utilities

<ApiTable
  columns={['Library', 'Available', 'Notes']}
  rows={[
    {symbol: 'fp-go', signature: 'Pipe, Flow, Compose, custom', description: 'Many composition combinators.'},
    {symbol: 'samber/lo', signature: 'Limited', description: 'No first-class Pipe/Flow/Compose.'},
    {symbol: 'go-functional', signature: 'Pipe, Flow', description: 'Some.'},
    {symbol: 'mo', signature: 'Limited', description: 'No first-class Pipe/Flow/Compose.'},
  ]}
/>

</Section>

<Section
  id="recommendations"
  number="06"
  title="Use-case"
  titleAccent="recommendations."
>

<ApiTable
  columns={['Use case', 'Winner', 'Snippet']}
  rows={[
    {symbol: 'Simple collection ops', signature: 'samber/lo', description: 'lo.Map(lo.Filter(items, predicate), transform)'},
    {symbol: 'Complex error handling', signature: 'fp-go', description: 'function.Pipe4(step1(), result.Chain(step2), result.Chain(step3), result.Chain(step4))'},
    {symbol: 'Lightweight Option/Result', signature: 'mo', description: 'mo.Some(value) · mo.Ok[int](42)'},
    {symbol: 'Full FP application', signature: 'fp-go', description: 'readerioresult.Ask + Chain pipelines'},
    {symbol: 'Iterator-based processing', signature: 'go-functional or fp-go', description: 'it.Map(it.Filter(it.Lift(data), pred), transform)'},
  ]}
/>

</Section>

<Section
  id="migration"
  number="07"
  title="Migration"
  titleAccent="paths."
>

### From samber/lo to fp-go

<Tabs>
  <TabItem value="before" label="Before (lo)" default>

<CodeCard file="before.go">
{`import "github.com/samber/lo"

result := lo.Map(
    lo.Filter(numbers, func(n int, _ int) bool {
        return n > 0
    }),
    func(n int, _ int) int {
        return n * 2
    },
)`}
</CodeCard>

  </TabItem>
  <TabItem value="after" label="After (fp-go)">

<CodeCard file="after.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/array"
    "github.com/IBM/fp-go/v2/function"
)

result := function.Pipe2(
    array.Filter(func(n int) bool { return n > 0 }),
    array.Map(func(n int) int { return n * 2 }),
)(numbers)`}
</CodeCard>

  </TabItem>
</Tabs>

### From mo to fp-go

<Tabs>
  <TabItem value="before" label="Before (mo)" default>

<CodeCard file="before.go">
{`import "github.com/samber/mo"

opt := mo.Some(42)
value := opt.OrElse(0)`}
</CodeCard>

  </TabItem>
  <TabItem value="after" label="After (fp-go)">

<CodeCard file="after.go" status="tested">
{`import "github.com/IBM/fp-go/v2/option"

opt := option.Some(42)
value := option.GetOrElse(func() int { return 0 })(opt)`}
</CodeCard>

  </TabItem>
</Tabs>

</Section>

<Section
  id="perf"
  number="08"
  title="Performance"
  titleAccent="comparison."
>

<Bench
  title="Filter · 1M ops"
  command="go test -bench=. -benchmem"
  rows={[
    {label: 'stdlib', bar: 0.92, barKind: 'win', nsOp: '1,050', bOp: '—', delta: 'baseline'},
    {label: 'samber/lo', bar: 0.95, barKind: 'win', nsOp: '1,100', bOp: '—', delta: '+5%'},
    {label: 'fp-go (idiomatic)', bar: 0.94, barKind: 'win', nsOp: '1,080', bOp: '—', delta: '+3%', winner: true},
    {label: 'fp-go (standard)', bar: 1.0, barKind: 'lose', nsOp: '1,200', bOp: '—', delta: '+14%', deltaKind: 'bad'},
  ]}
/>

<Callout type="success" title="Bottom line.">
  <ul>
    <li><strong>samber/lo:</strong> Excellent performance, close to stdlib</li>
    <li><strong>fp-go standard:</strong> Small overhead for type safety</li>
    <li><strong>fp-go idiomatic:</strong> Near-native performance</li>
    <li><strong>All are fast enough</strong> for most use cases</li>
  </ul>
</Callout>

</Section>

<Section
  id="decision"
  number="09"
  title="Decision"
  titleAccent="matrix."
>

<Compare>
  <CompareCol kind="good" title="Choose fp-go" pill="if you need">
    <p>Full FP toolkit.</p>
    <p>Complex error handling.</p>
    <p>IO and effect management.</p>
    <p>Reader pattern (dependency injection).</p>
    <p>Comprehensive documentation.</p>
    <p>Production support.</p>
  </CompareCol>
  <CompareCol kind="bad" title="Choose samber/lo" pill="if you need">
    <p>Simple collection operations.</p>
    <p>Minimal learning curve.</p>
    <p>Maximum performance.</p>
    <p>Wide community adoption.</p>
    <p>Quick wins.</p>
  </CompareCol>
</Compare>

<Compare>
  <CompareCol kind="bad" title="Choose mo" pill="if you need">
    <p>Lightweight Option/Result.</p>
    <p>Simple API.</p>
    <p>Minimal dependencies.</p>
    <p>Quick adoption.</p>
  </CompareCol>
  <CompareCol kind="good" title="Choose go-functional" pill="if you need">
    <p>Iterator-based operations.</p>
    <p>Lazy evaluation.</p>
    <p>Simpler than fp-go.</p>
    <p>Basic monadic types.</p>
  </CompareCol>
</Compare>

</Section>

<Section
  id="coexist"
  number="10"
  title="Can you use"
  titleAccent="multiple libraries?"
>

<Callout type="success" title="Yes. They coexist.">
  <CodeCard file="coexist.go">
{`import (
    "github.com/IBM/fp-go/v2/result"
    "github.com/samber/lo"
)

// Use lo for simple operations
filtered := lo.Filter(items, predicate)

// Use fp-go for complex error handling
result := result.TraverseArray(func(item Item) result.Result[Processed] {
    return processWithErrors(item)
})(filtered)`}
  </CodeCard>
</Callout>

<Callout title="Best practice.">
  <ul>
    <li>Use samber/lo for simple collection operations.</li>
    <li>Use fp-go for error handling and composition.</li>
    <li>Use mo for lightweight Option/Result in simple modules.</li>
    <li>Choose one as your primary library.</li>
  </ul>
</Callout>

</Section>

<Section
  id="summary"
  number="11"
  title="Summary"
>

<ApiTable
  columns={['Aspect', 'Best choice', '']}
  rows={[
    {symbol: 'Simplicity', signature: 'samber/lo or mo', description: ''},
    {symbol: 'Type safety', signature: 'fp-go', description: ''},
    {symbol: 'Performance', signature: 'samber/lo or fp-go idiomatic', description: ''},
    {symbol: 'Error handling', signature: 'fp-go', description: ''},
    {symbol: 'Learning curve', signature: 'samber/lo', description: ''},
    {symbol: 'Composition', signature: 'fp-go', description: ''},
    {symbol: 'Documentation', signature: 'fp-go or samber/lo', description: ''},
    {symbol: 'Production use', signature: 'fp-go or samber/lo', description: ''},
    {symbol: 'Full FP', signature: 'fp-go', description: ''},
  ]}
/>

</Section>
