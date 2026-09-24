---
sidebar_position: 8
title: v1 to v2 Migration
hide_title: true
description: Complete step-by-step guide for migrating from fp-go v1 to v2, covering all breaking changes with detailed examples and solutions.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


<PageHeader
  eyebrow="Migration · 02 / 03"
  title="v1 to v2"
  titleAccent="Migration."
  lede="Complete step-by-step guide for migrating from fp-go v1 to v2. Detailed examples, solutions, and testing strategies for each breaking change."
  meta={[
    {label: '// Steps', value: '5 main steps'},
    {label: '// Time estimate', value: '1-4 weeks'},
    {label: '// Reading time', value: '20 min · 8 sections'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Prerequisites" prose value={<>Go 1.24+, backup code, <em>review usage</em>.</>} />
  <TLDRCard label="// Main work" prose value={<>Update imports, fix <em>5 breaking changes</em>.</>} />
  <TLDRCard label="// Testing" prose value={<>Unit, integration, <em>performance</em> tests.</>} variant="up" />
</TLDR>

<Section id="prerequisites" number="01" title="Prerequisites" titleAccent="checklist.">

Before starting migration:

### 1. Upgrade Go Version

<CodeCard file="check-go.sh" lang="bash">
{`# Check current version
go version

# Must be 1.24 or higher
# If not, upgrade Go first`}
</CodeCard>

<Callout type="warn">
**Why:** v2 requires Go 1.24+ for generic type alias support.
</Callout>

### 2. Backup Your Code

<CodeCard file="backup.sh" lang="bash">
{`# Create a migration branch
git checkout -b migrate-to-fp-go-v2

# Or tag current state
git tag pre-fp-go-v2-migration`}
</CodeCard>

### 3. Review Current Usage

<CodeCard file="review.sh" lang="bash">
{`# Find all fp-go imports
grep -r "github.com/IBM/fp-go" . --include="*.go"

# Count usages
grep -r "github.com/IBM/fp-go" . --include="*.go" | wc -l`}
</CodeCard>

</Section>

<Section id="step-1" number="02" title="Step 1: Update" titleAccent="dependencies.">

### Add v2 Dependency

<CodeCard file="add-v2.sh" lang="bash">
{`# Add v2 (keeps v1 if already present)
go get github.com/IBM/fp-go/v2

# Update go.mod
go mod tidy`}
</CodeCard>

Your `go.mod` should now have:

<CodeCard file="go.mod">
{`require (
    github.com/IBM/fp-go v1.x.x        // Optional: keep for gradual migration
    github.com/IBM/fp-go/v2 v2.x.x     // New v2 dependency
)`}
</CodeCard>

### Remove v1 (Optional)

<Callout type="info">
Only after full migration is complete.
</Callout>

<CodeCard file="remove-v1.sh" lang="bash">
{`# Only after full migration
go get github.com/IBM/fp-go@none
go mod tidy`}
</CodeCard>

</Section>

<Section id="step-2" number="03" title="Step 2: Update" titleAccent="imports.">

### Automated Approach (Recommended)

Create a script to update imports:

<CodeCard file="migrate-imports.sh" lang="bash">
{`#!/bin/bash
# migrate-imports.sh

# Update all .go files
find . -name "*.go" -type f -exec sed -i '' \\
  's|github.com/IBM/fp-go/|github.com/IBM/fp-go/v2/|g' {} +

echo "Import paths updated. Run 'go build' to check for issues."`}
</CodeCard>

Run it:

<CodeCard file="run-script.sh" lang="bash">
{`chmod +x migrate-imports.sh
./migrate-imports.sh`}
</CodeCard>

### Manual Approach

Update each import:

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-imports.go">
{`// Before (v1)
import (
    "github.com/IBM/fp-go/either"
    "github.com/IBM/fp-go/option"
    "github.com/IBM/fp-go/ioeither"
)`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-imports.go">
{`// After (v2)
import (
    "github.com/IBM/fp-go/v2/either"
    "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/ioeither"
)`}
</CodeCard>

</CompareCol>
</Compare>

### Gradual Migration (v1 + v2)

Use import aliases:

<CodeCard file="aliases.go">
{`import (
    v1either "github.com/IBM/fp-go/either"
    v2either "github.com/IBM/fp-go/v2/either"
    
    v1option "github.com/IBM/fp-go/option"
    v2option "github.com/IBM/fp-go/v2/option"
)`}
</CodeCard>

</Section>

<Section id="step-3" number="04" title="Step 3: Fix breaking" titleAccent="changes.">

### Breaking Change 1: Generic Type Aliases

**What Changed:**

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-type-def.go">
{`// v1 - type definition
type IOEither[E, A any] func() E.Either[E, A]`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-type-alias.go">
{`// v2 - type alias
type IOEither[E, A any] = func() E.Either[E, A]`}
</CodeCard>

</CompareCol>
</Compare>

**Impact:** Mostly internal. Your code likely works without changes.

**Action Required:**

If you defined custom types based on fp-go types:

<Compare>
<CompareCol kind="bad">

<CodeCard file="custom-v1.go">
{`// v1 - Update this
type MyEither[E, A any] either.Either[E, A]`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="custom-v2.go">
{`// v2 - To this
type MyEither[E, A any] = either.Either[E, A]`}
</CodeCard>

</CompareCol>
</Compare>

### Breaking Change 2: Type Parameter Reordering

**What Changed:**

Type parameters that cannot be inferred are now first.

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-params.go">
{`// v1 signature
func Map[A, B any](f func(A) B) func(Either[error, A]) Either[error, B]`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-params.go">
{`// v2 signature
func Map[B, A any](f func(A) B) func(Either[error, A]) Either[error, B]
//       ^  ^
//       |  Inferred from function argument
//       Cannot be inferred, so comes first`}
</CodeCard>

</CompareCol>
</Compare>

**Impact:** Most code works without changes due to type inference.

**Migration Pattern:**

<Tabs groupId="migration">
<TabItem value="v1" label="v1 Code">

<CodeCard file="v1-explicit.go">
{`// v1 - explicit types
result := either.Chain[User, UserProfile](func(u User) either.Either[error, UserProfile] {
    return fetchProfile(u.ID)
})(userEither)`}
</CodeCard>

</TabItem>
<TabItem value="v2-explicit" label="v2 (Explicit)">

<CodeCard file="v2-explicit.go">
{`// v2 - reordered types
result := either.Chain[UserProfile, User](func(u User) either.Either[error, UserProfile] {
    return fetchProfile(u.ID)
})(userEither)`}
</CodeCard>

</TabItem>
<TabItem value="v2-inferred" label="v2 (Inferred)">

<CodeCard file="v2-inferred.go">
{`// v2 - inferred (recommended)
result := either.Chain(func(u User) either.Either[error, UserProfile] {
    return fetchProfile(u.ID)
})(userEither)`}
</CodeCard>

</TabItem>
</Tabs>

<Callout type="success">
**Action Required:**
- ✅ Remove explicit type parameters (let Go infer)
- ⚠️ If you must specify types, reverse the order
</Callout>

### Breaking Change 3: Pair Operates on Second Element

**What Changed:**

Pair operations now target the second element instead of the first.

**Why:** Aligns with Haskell and other FP languages.

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-pair.go">
{`// v1 - operates on FIRST element
pair := pair.MakePair(1, "hello")
mapped := pair.Map(func(x int) int { 
    return x * 2 
})(pair)
// Result: Pair(2, "hello")`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-pair.go">
{`// v2 - operates on SECOND element
pair := pair.MakePair(1, "hello")
mapped := pair.Map(func(s string) string { 
    return strings.ToUpper(s) 
})(pair)
// Result: Pair(1, "HELLO")`}
</CodeCard>

</CompareCol>
</Compare>

**Action Required:**
- ⚠️ Review all Pair usage
- ⚠️ Update Map, Chain, etc. to target second element
- ⚠️ Or swap pair elements if needed

### Breaking Change 4: Compose is Right-to-Left

**What Changed:**

Compose now applies functions right-to-left (mathematical composition).

**Why:** Aligns with mathematical notation: (f ∘ g)(x) = f(g(x))

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

**Action Required:**
- ⚠️ Reverse function order in Compose calls
- ✅ Or switch to Flow (left-to-right, unchanged)

### Breaking Change 5: No generic/ Subpackages

**What Changed:**

Removed `generic/` subpackages from all modules.

**Why:** Generic type aliases make them unnecessary.

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

**Action Required:**
- ⚠️ Remove `/generic` from import paths
- ⚠️ Update function calls if needed

</Section>

<Section id="step-4" number="05" title="Step 4: Adopt v2" titleAccent="best practices.">

### Use Result Instead of Either

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

### Use IOResult Instead of IOEither

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

### Use Idiomatic Packages

<CodeCard file="idiomatic.go">
{`// fp-go provides idiomatic (faster) versions
import "github.com/IBM/fp-go/v2/array/idiomatic"

// 2-32x faster for array operations
filtered := idiomatic.Filter(predicate)(array)
mapped := idiomatic.Map(transform)(filtered)`}
</CodeCard>

</Section>

<Section id="step-5" number="06" title="Step 5: Test" titleAccent="thoroughly.">

### Unit Tests

<CodeCard file="unit-test.go">
{`func TestMigration(t *testing.T) {
    // Test v2 behavior
    result := divide(10, 2)
    
    assert.True(t, result.IsOk())
    assert.Equal(t, 5, result.GetOrElse(func() int { return 0 }))
    
    // Test error case
    errorResult := divide(10, 0)
    assert.True(t, errorResult.IsErr())
}`}
</CodeCard>

### Integration Tests

<CodeCard file="integration-test.go">
{`func TestEndToEnd(t *testing.T) {
    // Test full pipeline with v2
    result := function.Pipe3(
        fetchData(),
        result.Chain(processData),
        result.Chain(saveData),
    )
    
    assert.True(t, result.IsOk())
}`}
</CodeCard>

### Performance Tests

<CodeCard file="benchmark.go">
{`func BenchmarkV1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = v1Function()
    }
}

func BenchmarkV2(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = v2Function()
    }
}

func BenchmarkV2Idiomatic(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = v2IdiomaticFunction()
    }
}`}
</CodeCard>

</Section>

<Section id="common-issues" number="07" title="Common migration" titleAccent="issues.">

### Issue 1: Type Inference Fails

**Problem:**
<CodeCard file="inference-fail.go">
{`// Compiler can't infer types
result := either.Map(transform)(myEither)
// Error: cannot infer type parameters`}
</CodeCard>

**Solution:**
<CodeCard file="inference-fix.go">
{`// Specify types explicitly
result := either.Map[OutputType, InputType](transform)(myEither)

// Or use type annotation
var result either.Either[error, OutputType] = either.Map(transform)(myEither)`}
</CodeCard>

### Issue 2: Pair Behavior Changed

**Problem:**
<CodeCard file="pair-problem.go">
{`// v1 code that operated on first element
mapped := pair.Map(func(x int) int { return x * 2 })(myPair)`}
</CodeCard>

**Solution:**
<CodeCard file="pair-fix.go">
{`// v2: Update to operate on second element
mapped := pair.Map(func(s string) string { return strings.ToUpper(s) })(myPair)

// Or swap the pair elements
swapped := pair.Swap(myPair)
mapped := pair.Map(func(x int) int { return x * 2 })(swapped)`}
</CodeCard>

### Issue 3: Compose Order Reversed

**Problem:**
<CodeCard file="compose-problem.go">
{`// v1 code with left-to-right composition
composed := function.Compose2(step1, step2)`}
</CodeCard>

**Solution:**
<CodeCard file="compose-fix.go">
{`// v2: Reverse order
composed := function.Compose2(step2, step1)

// Or use Flow (unchanged)
pipeline := function.Flow2(step1, step2)`}
</CodeCard>

### Issue 4: Generic Import Not Found

**Problem:**
<CodeCard file="generic-problem.go">
{`// v1 code with generic subpackage
import "github.com/IBM/fp-go/ioeither/generic"`}
</CodeCard>

**Solution:**
<CodeCard file="generic-fix.go">
{`// v2: Remove /generic
import "github.com/IBM/fp-go/v2/ioeither"`}
</CodeCard>

### Issue 5: Performance Regression

**Problem:**
Performance slower after migration.

**Solution:**
<CodeCard file="performance-fix.go">
{`// Use idiomatic packages for better performance
import "github.com/IBM/fp-go/v2/array/idiomatic"

// 2-32x faster
result := idiomatic.Map(transform)(array)`}
</CodeCard>

</Section>

<Section id="verification" number="08" title="Verification" titleAccent="checklist.">

<Checklist>
  <ChecklistItem status="required">All imports updated to v2</ChecklistItem>
  <ChecklistItem status="required">Breaking changes fixed</ChecklistItem>
  <ChecklistItem status="required">All tests passing</ChecklistItem>
  <ChecklistItem status="required">Code compiles without errors</ChecklistItem>
  <ChecklistItem status="optional">Performance benchmarks run</ChecklistItem>
  <ChecklistItem status="optional">Documentation updated</ChecklistItem>
  <ChecklistItem status="optional">Team trained on v2 changes</ChecklistItem>
  <ChecklistItem status="optional">Rollback plan documented</ChecklistItem>
</Checklist>

</Section>
