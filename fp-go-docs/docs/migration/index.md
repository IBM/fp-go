---
sidebar_position: 7
title: Migration Guide
hide_title: true
description: Comprehensive guide for migrating from fp-go v1 to v2, understanding breaking changes, and strategies for smooth transitions.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Migration · 01 / 03"
  title="Migration"
  titleAccent="Guide."
  lede="Migrate from fp-go v1 to v2 with confidence. Understand breaking changes, choose your strategy, and execute a smooth transition."
  meta={[
    {label: '// Difficulty', value: 'Low to Medium'},
    {label: '// Time estimate', value: '1-4 weeks'},
    {label: '// Reading time', value: '8 min · 5 sections'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Breaking changes" value="5" description="Generic aliases, type params, Pair, Compose, no generic/ packages." />
  <TLDRCard label="// Recommended strategy" prose value={<><em>Gradual</em> migration for production.</>} variant="up" />
  <TLDRCard label="// Risk level" prose value={<>Low with <em>incremental</em> testing.</>} />
</TLDR>

<Section id="should-you-migrate" number="01" title="Should you" titleAccent="migrate?">

### ✅ Reasons to Migrate to v2

**New Features:**
- **Result type** - Recommended over Either for error handling
- **Effect type** - Combines Reader + IO + Result
- **Idiomatic packages** - 2-32x faster performance
- **Better type inference** - Improved type parameter ordering
- **Generic type aliases** - Cleaner type definitions

**Improvements:**
- More intuitive API
- Better documentation
- Active development
- Future-proof

**Requirements:**
- Go 1.24+ (uses new generic features)

### ⚠️ Reasons to Stay on v1

**Valid reasons:**
- Stuck on Go 1.18-1.23
- Large existing v1 codebase
- Need Writer monad (v1 only, not in v2)
- Team bandwidth constraints

<Callout type="info">
**Note:** v1 is in maintenance mode but still supported.
</Callout>

</Section>

<Section id="overview" number="02" title="Migration" titleAccent="overview.">

### The 5 Breaking Changes

v2 introduces 5 breaking changes that require code updates:

1. **Generic Type Aliases** - `type X = Y` instead of `type X Y`
2. **Type Parameter Reordering** - Non-inferrable parameters first
3. **Pair Operates on Second Element** - v1 was first, v2 is second (Haskell-aligned)
4. **Compose is Right-to-Left** - Mathematical composition order
5. **No generic/ Subpackages** - Removed internal generic packages

<Callout type="success">
**Impact:** Most code requires only import path changes. Some code needs minor adjustments.
</Callout>

### Quick Migration Checklist

<Checklist>
  <ChecklistItem status="required">Upgrade to Go 1.24+</ChecklistItem>
  <ChecklistItem status="required">Review breaking changes</ChecklistItem>
  <ChecklistItem status="required">Plan migration strategy</ChecklistItem>
  <ChecklistItem status="required">Set up testing environment</ChecklistItem>
  <ChecklistItem status="optional">Update dependencies</ChecklistItem>
  <ChecklistItem status="optional">Update imports</ChecklistItem>
  <ChecklistItem status="optional">Fix breaking changes</ChecklistItem>
  <ChecklistItem status="optional">Test thoroughly</ChecklistItem>
</Checklist>

</Section>

<Section id="strategies" number="03" title="Migration" titleAccent="strategies.">

### Strategy 1: Big Bang (Small Codebases)

**Best for:**
- Small codebases (`<10k` lines)
- Few fp-go usages
- Can afford downtime

**Steps:**
1. Update all imports at once
2. Fix breaking changes
3. Test everything
4. Deploy

**Pros:**
- ✅ Clean, no mixed versions
- ✅ Fast migration
- ✅ Simple

**Cons:**
- ❌ Risky for large codebases
- ❌ All-or-nothing

### Strategy 2: Gradual Migration (Recommended)

**Best for:**
- Large codebases
- Production systems
- Risk-averse teams

**Steps:**
1. Run v1 and v2 side-by-side
2. Migrate module by module
3. Test each module
4. Remove v1 when done

**Pros:**
- ✅ Low risk
- ✅ Incremental testing
- ✅ Can pause/resume

**Cons:**
- ⚠️ Longer timeline
- ⚠️ Mixed versions temporarily

<CodeCard file="gradual-migration.go">
{`import (
    v1either "github.com/IBM/fp-go/either"
    v2result "github.com/IBM/fp-go/v2/result"
)

// Old code uses v1
func oldFunction() v1either.Either[error, int] {
    // ...
}

// New code uses v2
func newFunction() v2result.Result[int] {
    // ...
}

// Bridge between versions
func bridge() v2result.Result[int] {
    v1Result := oldFunction()
    return v1either.Fold(
        func(err error) v2result.Result[int] {
            return v2result.Err[int](err)
        },
        func(val int) v2result.Result[int] {
            return v2result.Ok(val)
        },
    )(v1Result)
}`}
</CodeCard>

### Strategy 3: New Code Only

**Best for:**
- Maintaining legacy code
- Limited resources
- Long-term migration

**Steps:**
1. Keep v1 for existing code
2. Use v2 for all new code
3. Gradually refactor when touching old code
4. Eventually remove v1

**Pros:**
- ✅ Minimal disruption
- ✅ Natural migration
- ✅ Low risk

**Cons:**
- ⚠️ Very long timeline
- ⚠️ Mixed versions indefinitely

</Section>

<Section id="breaking-changes" number="04" title="Breaking change" titleAccent="details.">

### 1. Generic Type Aliases

**What Changed:**
v2 uses generic type aliases (`type X = Y`) instead of type definitions (`type X Y`).

**Why:**
Go 1.24 added support for generic type aliases, allowing cleaner type definitions.

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-type-def.go">
{`// v1 - type definition
type ReaderIOEither[R, E, A any] RD.Reader[R, IOE.IOEither[E, A]]`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-type-alias.go">
{`// v2 - type alias
type ReaderIOEither[R, E, A any] = RD.Reader[R, IOE.IOEither[E, A]]`}
</CodeCard>

</CompareCol>
</Compare>

**Action Required:**
- ✅ None for most users
- ⚠️ Update custom type definitions if you created them

### 2. Type Parameter Reordering

**What Changed:**
Type parameters that cannot be inferred are now first.

**Why:**
Better type inference. Go can infer trailing type parameters but not leading ones.

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-params.go">
{`// v1
func Map[A, B any](f func(A) B) func(Either[error, A]) Either[error, B]`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-params.go">
{`// v2
func Map[B, A any](f func(A) B) func(Either[error, A]) Either[error, B]
//       ^  ^
//       |  Can be inferred from function argument
//       Cannot be inferred, so comes first`}
</CodeCard>

</CompareCol>
</Compare>

### 3. Pair Operates on Second Element

**What Changed:**
Pair operations now target the second element instead of the first.

**Why:**
Aligns with Haskell and other FP languages. More intuitive for most use cases.

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-pair.go">
{`// v1 - operates on FIRST element
pair := pair.MakePair(1, "hello")
mapped := pair.Map(func(x int) int { return x * 2 })(pair)
// Result: Pair(2, "hello")`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-pair.go">
{`// v2 - operates on SECOND element
pair := pair.MakePair(1, "hello")
mapped := pair.Map(func(s string) string { return strings.ToUpper(s) })(pair)
// Result: Pair(1, "HELLO")`}
</CodeCard>

</CompareCol>
</Compare>

### 4. Compose is Right-to-Left

**What Changed:**
Compose now applies functions right-to-left (mathematical composition).

**Why:**
Aligns with mathematical notation: (f ∘ g)(x) = f(g(x))

<Tabs groupId="compose">
<TabItem value="v1" label="v1 - Left-to-Right">

<CodeCard file="v1-compose.go">
{`composed := function.Compose2(
    func(x int) int { return x + 1 },  // Applied first
    func(x int) int { return x * 2 },  // Applied second
)
result := composed(5) // (5 + 1) * 2 = 12`}
</CodeCard>

</TabItem>
<TabItem value="v2" label="v2 - Right-to-Left">

<CodeCard file="v2-compose.go">
{`composed := function.Compose2(
    func(x int) int { return x * 2 },  // Applied second
    func(x int) int { return x + 1 },  // Applied first
)
result := composed(5) // (5 + 1) * 2 = 12`}
</CodeCard>

</TabItem>
<TabItem value="flow" label="Or Use Flow">

<CodeCard file="flow.go">
{`// Flow is left-to-right, unchanged
pipeline := function.Flow2(
    func(x int) int { return x + 1 },  // Applied first
    func(x int) int { return x * 2 },  // Applied second
)
result := pipeline(5) // (5 + 1) * 2 = 12`}
</CodeCard>

</TabItem>
</Tabs>

### 5. No generic/ Subpackages

**What Changed:**
Removed `generic/` subpackages from all modules.

**Why:**
Generic type aliases make them unnecessary. Cleaner API.

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-generic.go">
{`// v1 - generic subpackage
import "github.com/IBM/fp-go/ioeither/generic"`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-no-generic.go">
{`// v2 - no generic subpackage
import "github.com/IBM/fp-go/v2/ioeither"`}
</CodeCard>

</CompareCol>
</Compare>

</Section>

<Section id="common-patterns" number="05" title="Common migration" titleAccent="patterns.">

### Pattern 1: Either → Result

<Callout type="success">
**Recommended:** Use Result for error handling in v2.
</Callout>

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-either.go">
{`// v1 - Either
func divide(a, b int) either.Either[error, int] {
    if b == 0 {
        return either.Left[int](errors.New("division by zero"))
    }
    return either.Right[error](a / b)
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-result.go">
{`// v2 - Result (recommended)
func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}`}
</CodeCard>

</CompareCol>
</Compare>

### Pattern 2: IOEither → IOResult

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-ioeither.go">
{`// v1 - IOEither
func readFile(path string) ioeither.IOEither[error, []byte] {
    return func() either.Either[error, []byte] {
        data, err := os.ReadFile(path)
        if err != nil {
            return either.Left[[]byte](err)
        }
        return either.Right[error](data)
    }
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-ioresult.go">
{`// v2 - IOResult (recommended)
func readFile(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}`}
</CodeCard>

</CompareCol>
</Compare>

</Section>
