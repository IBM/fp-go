---
title: Record - Chain
hide_title: true
description: FlatMap operations for maps - transform values to maps and flatten the result.
sidebar_position: 16
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Record"
  titleAccent="Chain"
  lede="FlatMap operations for maps. Transform map values into new maps and flatten the result into a single map."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Operations', value: 'Chain, ChainWithIndex' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Chain` | `func Chain[K comparable, A, B any](Monoid[K, B]) func(func(A) map[K]B) func(map[K]A) map[K]B` | Transform and flatten |
| `ChainWithIndex` | `func ChainWithIndex[K comparable, A, B any](Monoid[K, B]) func(func(K, A) map[K]B) func(map[K]A) map[K]B` | Chain with keys |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Chain - Basic

<CodeCard file="chain.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    M "github.com/IBM/fp-go/v2/monoid"
    F "github.com/IBM/fp-go/v2/function"
)

m := map[string]int{"a": 1, "b": 2}

// Expand each entry into multiple entries
expanded := F.Pipe2(
    m,
    R.Chain(M.MergeMonoid[string, int]())(func(v int) map[string]int {
        return map[string]int{
            fmt.Sprintf("val_%d", v):    v,
            fmt.Sprintf("double_%d", v): v * 2,
        }
    }),
)
// map[string]int{
//   "val_1": 1, "double_1": 2,
//   "val_2": 2, "double_2": 4,
// }
`}
</CodeCard>

### ChainWithIndex

<CodeCard file="chain_index.go">
{`m := map[string]int{"a": 1, "b": 2}

// Use keys in expansion
result := F.Pipe2(
    m,
    R.ChainWithIndex(M.MergeMonoid[string, int]())(
        func(k string, v int) map[string]int {
            return map[string]int{
                k + "_original": v,
                k + "_doubled":  v * 2,
            }
        },
    ),
)
// map[string]int{
//   "a_original": 1, "a_doubled": 2,
//   "b_original": 2, "b_doubled": 4,
// }
`}
</CodeCard>

### Expanding Configuration

<CodeCard file="config.go">
{`type Config struct {
    Replicas int
    Prefix   string
}

configs := map[string]Config{
    "web": {Replicas: 3, Prefix: "web"},
    "api": {Replicas: 2, Prefix: "api"},
}

// Generate instance names
instances := F.Pipe2(
    configs,
    R.ChainWithIndex(M.MergeMonoid[string, string]())(
        func(service string, cfg Config) map[string]string {
            result := make(map[string]string)
            for i := 0; i < cfg.Replicas; i++ {
                key := fmt.Sprintf("%s-%d", service, i)
                result[key] = fmt.Sprintf("%s-%d.example.com", cfg.Prefix, i)
            }
            return result
        },
    ),
)
// map[string]string{
//   "web-0": "web-0.example.com",
//   "web-1": "web-1.example.com",
//   "web-2": "web-2.example.com",
//   "api-0": "api-0.example.com",
//   "api-1": "api-1.example.com",
// }
`}
</CodeCard>

### Dependency Resolution

<CodeCard file="deps.go">
{`type Package struct {
    Name         string
    Dependencies []string
}

packages := map[string]Package{
    "app": {
        Name:         "app",
        Dependencies: []string{"lib1", "lib2"},
    },
}

// Flatten dependencies
allDeps := F.Pipe2(
    packages,
    R.Chain(M.MergeMonoid[string, bool]())(
        func(pkg Package) map[string]bool {
            deps := make(map[string]bool)
            for _, dep := range pkg.Dependencies {
                deps[dep] = true
            }
            return deps
        },
    ),
)
// map[string]bool{"lib1": true, "lib2": true}
`}
</CodeCard>

### Tag Expansion

<CodeCard file="tags.go">
{`type Resource struct {
    Name string
    Tags []string
}

resources := map[string]Resource{
    "server1": {
        Name: "server1",
        Tags: []string{"prod", "web"},
    },
    "server2": {
        Name: "server2",
        Tags: []string{"prod", "api"},
    },
}

// Create tag index
tagIndex := F.Pipe2(
    resources,
    R.ChainWithIndex(M.MergeMonoid[string, []string]())(
        func(id string, res Resource) map[string][]string {
            result := make(map[string][]string)
            for _, tag := range res.Tags {
                result[tag] = []string{id}
            }
            return result
        },
    ),
)
// map[string][]string{
//   "prod": ["server1", "server2"],
//   "web": ["server1"],
//   "api": ["server2"],
// }
`}
</CodeCard>

</Section>

<Section id="patterns" number="03" title="Common" titleAccent="Patterns">

### Nested Data Flattening

<CodeCard file="flatten.go">
{`type Category struct {
    Name  string
    Items []string
}

categories := map[string]Category{
    "fruits": {
        Name:  "Fruits",
        Items: []string{"apple", "banana"},
    },
    "veggies": {
        Name:  "Vegetables",
        Items: []string{"carrot", "lettuce"},
    },
}

// Flatten to item -> category mapping
itemMap := F.Pipe2(
    categories,
    R.ChainWithIndex(M.MergeMonoid[string, string]())(
        func(catKey string, cat Category) map[string]string {
            result := make(map[string]string)
            for _, item := range cat.Items {
                result[item] = catKey
            }
            return result
        },
    ),
)
// map[string]string{
//   "apple": "fruits",
//   "banana": "fruits",
//   "carrot": "veggies",
//   "lettuce": "veggies",
// }
`}
</CodeCard>

### Permission Expansion

<CodeCard file="permissions.go">
{`type Role struct {
    Name        string
    Permissions []string
}

roles := map[string]Role{
    "admin": {
        Name:        "admin",
        Permissions: []string{"read", "write", "delete"},
    },
    "user": {
        Name:        "user",
        Permissions: []string{"read"},
    },
}

// Create permission -> roles mapping
permMap := F.Pipe2(
    roles,
    R.ChainWithIndex(M.MergeMonoid[string, []string]())(
        func(roleKey string, role Role) map[string][]string {
            result := make(map[string][]string)
            for _, perm := range role.Permissions {
                result[perm] = []string{roleKey}
            }
            return result
        },
    ),
)
// map[string][]string{
//   "read": ["admin", "user"],
//   "write": ["admin"],
//   "delete": ["admin"],
// }
`}
</CodeCard>

</Section>

<Callout type="info">

**Monoid Strategy**: The monoid parameter determines how duplicate keys are handled when flattening. Use `MergeMonoid` for last-wins, or custom monoids for other strategies like concatenation or summation.

</Callout>
