---
sidebar_position: 6
title: Data Transformation
description: Pipeline-based data processing
hide_title: true
---

<PageHeader
  eyebrow="Recipes · 06 / 17"
  title="Data"
  titleAccent="Transformation"
  lede="Transform data using functional pipelines, composing operations for clean and maintainable data processing."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '6' },
    { label: 'Use Cases', value: 'ETL, Normalization, Enrichment' }
  ]}
/>

<TLDR>
  <TLDRCard title="Use Pipelines" icon="git-branch">
    Chain operations for clarity and maintainability—compose small functions into powerful transformations.
  </TLDRCard>
  <TLDRCard title="Keep Transformations Pure" icon="shield">
    No side effects in transformation functions—makes code predictable, testable, and composable.
  </TLDRCard>
  <TLDRCard title="Normalize Early" icon="filter">
    Convert to common format at system boundaries—simplifies downstream processing and reduces complexity.
  </TLDRCard>
</TLDR>

<Section id="basic-pipeline" number="01" title="Basic" titleAccent="Pipeline">

Transform data through a series of operations using function composition.

<CodeCard file="basic-pipeline.go">
{`package main

import (
    "fmt"
    "strings"
    
    F "github.com/IBM/fp-go/v2/function"
)

// Transform functions
func trim(s string) string {
    return strings.TrimSpace(s)
}

func toLowerCase(s string) string {
    return strings.ToLower(s)
}

func addPrefix(prefix string) func(string) string {
    return func(s string) string {
        return prefix + s
    }
}

func main() {
    input := "  HELLO WORLD  "
    
    // Pipeline: trim -> lowercase -> add prefix
    result := F.Pipe3(
        input,
        trim,
        toLowerCase,
        addPrefix("processed: "),
    )
    
    fmt.Println(result)
    // Output: processed: hello world
}`}
</CodeCard>

</Section>

<Section id="array-pipeline" number="02" title="Array Transformation" titleAccent="Pipeline">

Process collections through transformation pipelines with filtering and mapping.

<CodeCard file="array-pipeline.go">
{`package main

import (
    "fmt"
    "strings"
    
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

type Product struct {
    Name  string
    Price float64
}

func main() {
    products := []Product{
        {Name: "  Apple  ", Price: 1.50},
        {Name: "BANANA", Price: 0.75},
        {Name: "cherry", Price: 2.00},
    }
    
    // Pipeline: normalize names -> filter by price -> extract names
    result := F.Pipe3(
        products,
        // Normalize product names
        A.Map(func(p Product) Product {
            return Product{
                Name:  strings.TrimSpace(strings.ToLower(p.Name)),
                Price: p.Price,
            }
        }),
        // Filter products under $2
        A.Filter(func(p Product) bool {
            return p.Price < 2.00
        }),
        // Extract just the names
        A.Map(func(p Product) string {
            return p.Name
        }),
    )
    
    fmt.Printf("Affordable products: %v\\n", result)
    // Output: Affordable products: [apple banana]
}`}
</CodeCard>

</Section>

<Section id="nested-data" number="03" title="Nested Data" titleAccent="Transformation">

Transform nested structures by mapping over multiple levels.

<CodeCard file="nested-transformation.go">
{`package main

import (
    "fmt"
    
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

type Order struct {
    ID    string
    Items []OrderItem
}

type OrderItem struct {
    Product  string
    Quantity int
    Price    float64
}

type OrderSummary struct {
    ID         string
    TotalItems int
    TotalPrice float64
}

func calculateTotal(items []OrderItem) float64 {
    return A.Reduce(
        func(acc float64, item OrderItem) float64 {
            return acc + (item.Price * float64(item.Quantity))
        },
        0.0,
    )(items)
}

func toSummary(order Order) OrderSummary {
    return OrderSummary{
        ID:         order.ID,
        TotalItems: len(order.Items),
        TotalPrice: calculateTotal(order.Items),
    }
}

func main() {
    orders := []Order{
        {
            ID: "ORD-001",
            Items: []OrderItem{
                {Product: "Apple", Quantity: 3, Price: 1.50},
                {Product: "Banana", Quantity: 2, Price: 0.75},
            },
        },
        {
            ID: "ORD-002",
            Items: []OrderItem{
                {Product: "Cherry", Quantity: 1, Price: 2.00},
            },
        },
    }
    
    // Transform orders to summaries
    summaries := A.Map(toSummary)(orders)
    
    for _, s := range summaries {
        fmt.Printf("%s: %d items, $%.2f\\n", s.ID, s.TotalItems, s.TotalPrice)
    }
}`}
</CodeCard>

</Section>

<Section id="normalization" number="04" title="Data" titleAccent="Normalization">

Normalize data from different sources into a common format for unified processing.

<CodeCard file="normalization.go">
{`package main

import (
    "fmt"
    "time"
    
    A "github.com/IBM/fp-go/v2/array"
)

// Different source formats
type APIUser struct {
    UserID    int    \`json:"user_id"\`
    FullName  string \`json:"full_name"\`
    EmailAddr string \`json:"email_addr"\`
}

type DBUser struct {
    ID    int
    Name  string
    Email string
}

// Normalized format
type User struct {
    ID    int
    Name  string
    Email string
}

// Converters
func fromAPI(apiUser APIUser) User {
    return User{
        ID:    apiUser.UserID,
        Name:  apiUser.FullName,
        Email: apiUser.EmailAddr,
    }
}

func fromDB(dbUser DBUser) User {
    return User{
        ID:    dbUser.ID,
        Name:  dbUser.Name,
        Email: dbUser.Email,
    }
}

func main() {
    apiUsers := []APIUser{
        {UserID: 1, FullName: "Alice Smith", EmailAddr: "alice@example.com"},
        {UserID: 2, FullName: "Bob Jones", EmailAddr: "bob@example.com"},
    }
    
    dbUsers := []DBUser{
        {ID: 3, Name: "Charlie Brown", Email: "charlie@example.com"},
    }
    
    // Normalize both sources
    normalizedAPI := A.Map(fromAPI)(apiUsers)
    normalizedDB := A.Map(fromDB)(dbUsers)
    
    // Combine
    allUsers := append(normalizedAPI, normalizedDB...)
    
    for _, u := range allUsers {
        fmt.Printf("User %d: %s (%s)\\n", u.ID, u.Name, u.Email)
    }
}`}
</CodeCard>

</Section>

<Section id="enrichment" number="05" title="Data" titleAccent="Enrichment">

Enrich data with additional information from related sources.

<CodeCard file="enrichment.go">
{`package main

import (
    "fmt"
    
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
)

type User struct {
    ID   int
    Name string
}

type Post struct {
    ID       int
    UserID   int
    Content  string
    UserName O.Option[string] // Enriched field
}

// Lookup user by ID
func findUser(users []User, userID int) O.Option[User] {
    for _, u := range users {
        if u.ID == userID {
            return O.Some(u)
        }
    }
    return O.None[User]()
}

// Enrich post with user name
func enrichPost(users []User) func(Post) Post {
    return func(post Post) Post {
        user := findUser(users, post.UserID)
        return Post{
            ID:       post.ID,
            UserID:   post.UserID,
            Content:  post.Content,
            UserName: O.Map(func(u User) string { return u.Name })(user),
        }
    }
}

func main() {
    users := []User{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }
    
    posts := []Post{
        {ID: 101, UserID: 1, Content: "Hello world"},
        {ID: 102, UserID: 2, Content: "Functional programming rocks"},
        {ID: 103, UserID: 999, Content: "From deleted user"},
    }
    
    // Enrich posts with user names
    enrichedPosts := A.Map(enrichPost(users))(posts)
    
    for _, p := range enrichedPosts {
        userName := O.GetOrElse(func() string { return "Unknown" })(p.UserName)
        fmt.Printf("Post %d by %s: %s\\n", p.ID, userName, p.Content)
    }
}`}
</CodeCard>

</Section>

<Section id="flattening" number="06" title="Flattening Nested" titleAccent="Structures">

Flatten nested arrays and structures into a single level.

<CodeCard file="flattening.go">
{`package main

import (
    "fmt"
    
    A "github.com/IBM/fp-go/v2/array"
)

type Department struct {
    Name      string
    Employees []string
}

func main() {
    departments := []Department{
        {Name: "Engineering", Employees: []string{"Alice", "Bob"}},
        {Name: "Sales", Employees: []string{"Charlie", "Diana"}},
        {Name: "HR", Employees: []string{"Eve"}},
    }
    
    // Flatten: extract all employees
    allEmployees := A.Chain(func(dept Department) []string {
        return dept.Employees
    })(departments)
    
    fmt.Printf("All employees: %v\\n", allEmployees)
    // Output: All employees: [Alice Bob Charlie Diana Eve]
    
    // With department prefix
    employeesWithDept := A.Chain(func(dept Department) []string {
        return A.Map(func(emp string) string {
            return fmt.Sprintf("%s (%s)", emp, dept.Name)
        })(dept.Employees)
    })(departments)
    
    for _, emp := range employeesWithDept {
        fmt.Println(emp)
    }
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="07" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Use pipelines** — Chain operations for clarity and maintainability
  </ChecklistItem>
  <ChecklistItem status="required">
    **Keep transformations pure** — No side effects in transformation functions
  </ChecklistItem>
  <ChecklistItem status="required">
    **Compose small functions** — Build complex transformations from simple ones
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Handle missing data** — Use Option for optional fields
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Normalize early** — Convert to common format at system boundaries
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Separate concerns** — Keep transformation logic separate from I/O
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Make transformations reusable** — Extract common patterns into functions
  </ChecklistItem>
</Checklist>

</Section>
