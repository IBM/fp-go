---
title: String
hide_title: true
description: String utilities and type class instances for text operations in fp-go
sidebar_position: 31
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="String"
  titleAccent="Type Classes"
  lede="The string package provides type class instances and utilities for string values, enabling functional operations on text."
  meta={[
    { label: 'Package', value: 'string' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Type', value: 'string' }
  ]}
/>

---

<Section num="1" title="Overview">

The **String** package provides type class instances for the `string` type, enabling:
- **Eq**: Equality comparison
- **Ord**: Lexicographic ordering
- **Semigroup**: Concatenation
- **Monoid**: Concatenation with empty string identity

These instances allow strings to be used with generic functional operations.

</Section>

---

<Section num="2" title="Equality">

<CodeCard file="string_eq.go" tag="example">
{`import S "github.com/IBM/fp-go/string"

// Compare strings for equality
S.Eq.Equals("hello", "hello")  // true
S.Eq.Equals("hello", "world")  // false
S.Eq.Equals("", "")            // true

// Case-sensitive comparison
S.Eq.Equals("Hello", "hello")  // false`}
</CodeCard>

</Section>

---

<Section num="3" title="Ordering">

Lexicographic (dictionary) ordering:

<CodeCard file="string_ord.go" tag="example">
{`import S "github.com/IBM/fp-go/string"

// Compare strings lexicographically
S.Ord.Compare("a", "b")     // -1 (less than)
S.Ord.Compare("b", "a")     // 1  (greater than)
S.Ord.Compare("a", "a")     // 0  (equal)
S.Ord.Compare("apple", "banana")  // -1

// Derived operations
S.Ord.LessThan("a", "b")           // true
S.Ord.GreaterThan("z", "a")        // true
S.Ord.LessThanOrEqual("a", "a")    // true

// Min and Max
S.Ord.Min("zebra", "apple")  // "apple"
S.Ord.Max("zebra", "apple")  // "zebra"`}
</CodeCard>

</Section>

---

<Section num="4" title="Semigroup Concatenation">

Combine strings with concatenation:

<CodeCard file="string_semigroup.go" tag="example">
{`import S "github.com/IBM/fp-go/string"

// Concatenate strings
S.Semigroup.Concat("Hello", " World")  // "Hello World"
S.Semigroup.Concat("foo", "bar")       // "foobar"

// Associativity holds
s1 := S.Semigroup.Concat(
    S.Semigroup.Concat("a", "b"),
    "c",
)
s2 := S.Semigroup.Concat(
    "a",
    S.Semigroup.Concat("b", "c"),
)
// s1 == s2 == "abc"`}
</CodeCard>

</Section>

---

<Section num="5" title="Monoid Operations">

Concatenation with empty string identity:

<CodeCard file="string_monoid.go" tag="example">
{`import S "github.com/IBM/fp-go/string"

// Concatenate with monoid
S.Monoid.Concat("Hello", " World")  // "Hello World"
S.Monoid.Empty()                    // ""

// Identity laws
text := "test"
S.Monoid.Concat(S.Monoid.Empty(), text)  // "test"
S.Monoid.Concat(text, S.Monoid.Empty())  // "test"`}
</CodeCard>

</Section>

---

<Section num="6" title="Sorting Strings">

Sort string arrays alphabetically:

<CodeCard file="string_sort.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    S "github.com/IBM/fp-go/string"
)

words := []string{"zebra", "apple", "mango", "banana"}

// Sort alphabetically (ascending)
sorted := F.Pipe2(
    words,
    A.Sort(S.Ord),
)
// []string{"apple", "banana", "mango", "zebra"}

// Sort descending
import O "github.com/IBM/fp-go/ord"

sortedDesc := F.Pipe2(
    words,
    A.Sort(O.Reverse(S.Ord)),
)
// []string{"zebra", "mango", "banana", "apple"}`}
</CodeCard>

</Section>

---

<Section num="7" title="Concatenating Arrays">

Join string arrays:

<CodeCard file="string_concat.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    S "github.com/IBM/fp-go/string"
)

words := []string{"Hello", "functional", "world"}

// Concatenate all strings
combined := F.Pipe2(
    words,
    A.Fold(S.Monoid),
)
// "Hellofunctionalworld"

// Join with separator
joined := F.Pipe3(
    words,
    A.Intersperse(" "),
    A.Fold(S.Monoid),
)
// "Hello functional world"

// Join with custom separator
withCommas := F.Pipe3(
    words,
    A.Intersperse(", "),
    A.Fold(S.Monoid),
)
// "Hello, functional, world"`}
</CodeCard>

</Section>

---

<Section num="8" title="Building Strings">

Use monoid to build strings from parts:

<CodeCard file="string_build.go" tag="example">
{`type User struct {
    FirstName string
    LastName  string
    Age       int
}

user := User{FirstName: "Alice", LastName: "Smith", Age: 30}

// Build formatted string
parts := []string{
    "User: ",
    user.FirstName,
    " ",
    user.LastName,
    " (age ",
    fmt.Sprintf("%d", user.Age),
    ")",
}

result := F.Pipe2(
    parts,
    A.Fold(S.Monoid),
)
// "User: Alice Smith (age 30)"

// Build CSV row
csvRow := F.Pipe3(
    []string{user.FirstName, user.LastName, fmt.Sprintf("%d", user.Age)},
    A.Intersperse(","),
    A.Fold(S.Monoid),
)
// "Alice,Smith,30"`}
</CodeCard>

</Section>

---

<Section num="9" title="Filtering and Mapping">

Combine with array operations:

<CodeCard file="string_filter.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    S "github.com/IBM/fp-go/string"
)

words := []string{"apple", "banana", "apricot", "cherry", "avocado"}

// Filter words starting with 'a'
aWords := F.Pipe2(
    words,
    A.Filter(func(s string) bool {
        return len(s) > 0 && s[0] == 'a'
    }),
)
// []string{"apple", "apricot", "avocado"}

// Sort and join
result := F.Pipe3(
    aWords,
    A.Sort(S.Ord),
    A.Intersperse(", "),
    A.Fold(S.Monoid),
)
// "apple, apricot, avocado"`}
</CodeCard>

</Section>

---

<Section num="10" title="API Reference">

<ApiTable>
| Instance | Type | Description |
|----------|------|-------------|
| `Eq` | `Eq[string]` | Equality comparison |
| `Ord` | `Ord[string]` | Lexicographic ordering |
| `Semigroup` | `Semigroup[string]` | Concatenation |
| `Monoid` | `Monoid[string]` | Concatenation with empty identity |
</ApiTable>

**Monoid Identity:**
- `Monoid.Empty()` returns `""` (empty string)

</Section>

---

<Section num="11" title="Related Concepts">

**Common Use Cases:**
- Sorting string arrays
- Joining strings with separators
- Building formatted strings
- String comparison and filtering

**See Also:**
- [Number](./number.md) - Numeric type class instances
- [Eq](./eq.md) - Equality type class
- [Ord](./ord.md) - Ordering type class
- [Monoid](./monoid.md) - Understanding monoid operations

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/number', title: 'Number' }}
  next={{ to: '/docs/v2/utilities/tuple', title: 'Tuple' }}
/>