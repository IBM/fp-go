---
sidebar_position: 16
title: The Zen of Go
hide_title: true
description: Learn how fp-go aligns with Go's philosophy and when to use functional patterns vs standard Go idioms.
---

<PageHeader
  eyebrow="Concepts · 06 / 06"
  title="The Zen of"
  titleAccent="Go."
  lede="Balance functional programming with Go's philosophy. Learn when to use fp-go patterns and when standard Go is better."
  meta={[
    {label: '// Difficulty', value: 'Intermediate'},
    {label: '// Reading time', value: '14 min · 8 sections'},
    {label: '// Prereqs', value: 'All previous concepts'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Philosophy" prose value={<>fp-go <em>complements</em> Go, doesn't replace it.</>} variant="up" />
  <TLDRCard label="// When to use" prose value={<>Complex errors, optional values, <em>data transforms</em>.</>} />
  <TLDRCard label="// When not to" prose value={<>Simple ops, hot paths, <em>team unfamiliar</em>.</>} />
</TLDR>

<Section id="philosophy" number="01" title="Go's" titleAccent="philosophy.">

Go values:

1. **Simplicity** - Easy to read and understand
2. **Clarity** - Explicit over implicit
3. **Pragmatism** - Practical over theoretical
4. **Composition** - Small pieces that work together
5. **Concurrency** - Built-in support for concurrent programming

<Callout type="success">
**fp-go embraces these values.**
</Callout>

</Section>

<Section id="go-first" number="02" title="fp-go is" titleAccent="Go-first.">

### Not a Replacement

<Compare>
<CompareCol kind="bad">

<CodeCard file="replace-all.go">
{`// ❌ Don't replace all Go code with fp-go
func processUser(id string) result.Result[User] {
    return function.Pipe10(
        step1, step2, step3, step4, step5,
        step6, step7, step8, step9, step10,
    )
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="use-where-valuable.go">
{`// ✅ Use fp-go where it adds value
func processUser(id string) (User, error) {
    user, err := fetchUser(id)
    if err != nil {
        return User{}, err
    }
    
    // Use fp-go for complex transformations
    validated := result.Chain(validateUser)(result.Ok(user))
    
    return validated.Fold(
        func(err error) (User, error) { return User{}, err },
        func(u User) (User, error) { return u, nil },
    )
}`}
</CodeCard>

</CompareCol>
</Compare>

### Complement, Don't Replace

<CodeCard file="complement.go">
{`// Standard Go for simple cases
func add(a, b int) int {
    return a + b
}

// fp-go for complex error handling
func processOrder(id string) result.Result[Order] {
    return function.Pipe3(
        fetchOrder(id),
        result.Chain(validateOrder),
        result.Chain(enrichOrder),
    )
}`}
</CodeCard>

</Section>

<Section id="when-to-use" number="03" title="When to use" titleAccent="fp-go.">

### ✅ Use fp-go When:

#### 1. Complex Error Handling

<Compare>
<CompareCol kind="bad">

<CodeCard file="repetitive-errors.go">
{`func processData(input string) (Output, error) {
    parsed, err := parse(input)
    if err != nil {
        return Output{}, err
    }
    
    validated, err := validate(parsed)
    if err != nil {
        return Output{}, err
    }
    
    transformed, err := transform(validated)
    if err != nil {
        return Output{}, err
    }
    
    return transformed, nil
}

// Repetitive error checking`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="clean-pipeline.go">
{`func processData(input string) result.Result[Output] {
    return function.Pipe3(
        parse(input),
        result.Chain(validate),
        result.Chain(transform),
    )
}

// Clean and composable`}
</CodeCard>

</CompareCol>
</Compare>

#### 2. Optional Values

<Compare>
<CompareCol kind="bad">

<CodeCard file="nil-checks.go">
{`func findUser(id string) *User {
    // nil means not found
    return db.FindByID(id)
}

// Lots of nil checks
user := findUser("123")
if user != nil {
    email := user.Email
    if email != "" {
        // Use email
    }
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="option-chain.go">
{`func findUser(id string) option.Option[User] {
    user := db.FindByID(id)
    if user == nil {
        return option.None[User]()
    }
    return option.Some(*user)
}

// Chain operations
email := function.Pipe2(
    findUser("123"),
    option.Chain(getEmail),
)`}
</CodeCard>

</CompareCol>
</Compare>

#### 3. Data Transformations

<Compare>
<CompareCol kind="bad">

<CodeCard file="imperative-transform.go">
{`func processUsers(users []User) []UserDTO {
    result := make([]UserDTO, 0, len(users))
    for _, user := range users {
        if user.Active {
            dto := UserDTO{
                ID:   user.ID,
                Name: strings.ToUpper(user.Name),
            }
            result = append(result, dto)
        }
    }
    return result
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="functional-transform.go">
{`import "github.com/IBM/fp-go/v2/array"

func processUsers(users []User) []UserDTO {
    return function.Pipe2(
        array.Filter(func(u User) bool { return u.Active }),
        array.Map(toDTO),
    )(users)
}

func toDTO(u User) UserDTO {
    return UserDTO{
        ID:   u.ID,
        Name: strings.ToUpper(u.Name),
    }
}`}
</CodeCard>

</CompareCol>
</Compare>

#### 4. Composable Pipelines

<CodeCard file="pipelines.go">
{`// Build reusable pipelines
var processUser = function.Flow3(
    normalize,
    validate,
    enrich,
)

// Use in different contexts
user1 := processUser(rawUser1)
user2 := processUser(rawUser2)`}
</CodeCard>

### ❌ Don't Use fp-go When:

#### 1. Simple Operations

<Compare>
<CompareCol kind="bad">

<CodeCard file="overkill.go">
{`// ❌ Overkill
result := result.Map(func(x int) int {
    return x + 1
})(result.Ok(5))`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="simple-go.go">
{`// ✅ Just use standard Go
x := 5 + 1`}
</CodeCard>

</CompareCol>
</Compare>

#### 2. Performance-Critical Code

<Compare>
<CompareCol kind="bad">

<CodeCard file="hot-path-fp.go">
{`// ❌ Hot path with millions of calls
for i := 0; i < 1000000; i++ {
    result := option.Map(transform)(opt)
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="hot-path-go.go">
{`// ✅ Use standard Go for hot paths
for i := 0; i < 1000000; i++ {
    if opt.IsSome() {
        value := transform(opt.Value())
    }
}`}
</CodeCard>

</CompareCol>
</Compare>

#### 3. Team Unfamiliar with FP

<Compare>
<CompareCol kind="bad">

<CodeCard file="confusing.go">
{`// ❌ Confusing for team
var pipeline = function.Flow5(
    step1, step2, step3, step4, step5,
)`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="gradual.go">
{`// ✅ Start simple, introduce gradually
func process(input Input) Output {
    step1Result := step1(input)
    step2Result := step2(step1Result)
    // ...
    return step5Result
}`}
</CodeCard>

</CompareCol>
</Compare>

#### 4. Standard Go is Clearer

<Compare>
<CompareCol kind="bad">

<CodeCard file="forced-fp.go">
{`// ❌ Forced FP
result := option.Fold(
    func() string { return "" },
    func(s string) string { return s },
)(opt)`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="clear-go.go">
{`// ✅ Clear Go
var result string
if opt.IsSome() {
    result = opt.Value()
}`}
</CodeCard>

</CompareCol>
</Compare>

</Section>

<Section id="balancing" number="04" title="Balancing FP and" titleAccent="Go idioms.">

### Pattern 1: FP Core, Go Boundaries

<CodeCard file="fp-core-go-boundaries.go">
{`// Pure FP core
func processOrder(order Order) result.Result[Order] {
    return function.Pipe3(
        result.Ok(order),
        result.Chain(validateOrder),
        result.Chain(enrichOrder),
    )
}

// Go-style API
func ProcessOrder(order Order) (Order, error) {
    result := processOrder(order)
    return result.Fold(
        func(err error) (Order, error) { return Order{}, err },
        func(o Order) (Order, error) { return o, nil },
    )
}`}
</CodeCard>

### Pattern 2: Gradual Adoption

<CodeCard file="gradual-adoption.go">
{`// Phase 1: Start with error handling
func fetchUser(id string) result.Result[User] {
    user, err := db.FindByID(id)
    return result.FromGoError(user, err)
}

// Phase 2: Add composition
func processUser(id string) result.Result[User] {
    return function.Pipe2(
        fetchUser(id),
        result.Chain(validateUser),
    )
}

// Phase 3: Full pipeline
func processUser(id string) result.Result[UserDTO] {
    return function.Pipe4(
        fetchUser(id),
        result.Chain(validateUser),
        result.Chain(enrichUser),
        result.Map(toDTO),
    )
}`}
</CodeCard>

### Pattern 3: Hybrid Approach

<CodeCard file="hybrid.go">
{`// Mix FP and standard Go
func processUsers(users []User) ([]UserDTO, error) {
    // Use fp-go for transformation
    active := array.Filter(func(u User) bool {
        return u.Active
    })(users)
    
    // Standard Go for complex logic
    var dtos []UserDTO
    for _, user := range active {
        // Complex validation
        if err := complexValidation(user); err != nil {
            return nil, err
        }
        
        // Use fp-go for transformation
        dto := toDTO(user)
        dtos = append(dtos, dto)
    }
    
    return dtos, nil
}`}
</CodeCard>

</Section>

<Section id="go-idioms" number="05" title="Go idioms with" titleAccent="fp-go.">

### Idiom 1: Errors are Values

<CodeCard file="errors-as-values.go">
{`// Go idiom: return errors
func fetchUser(id string) (User, error) {
    // ...
}

// fp-go: errors in Result
func fetchUser(id string) result.Result[User] {
    // ...
}

// Both are valid - choose based on context`}
</CodeCard>

### Idiom 2: Accept Interfaces, Return Structs

<CodeCard file="interfaces.go">
{`// Go idiom
type UserService interface {
    GetUser(id string) (User, error)
}

type userService struct {
    db Database
}

func (s *userService) GetUser(id string) (User, error) {
    // Can use fp-go internally
    result := function.Pipe2(
        s.fetchUser(id),
        result.Chain(s.validateUser),
    )
    
    // Return Go-style
    return result.Fold(
        func(err error) (User, error) { return User{}, err },
        func(u User) (User, error) { return u, nil },
    )
}`}
</CodeCard>

### Idiom 3: Make Zero Values Useful

<CodeCard file="zero-values.go">
{`// Go idiom: zero values work
type Config struct {
    Port    int    // 0 is valid
    Host    string // "" is valid
    Timeout time.Duration // 0 is valid
}

// fp-go: use Option for truly optional values
type Config struct {
    Port    int
    Host    string
    Timeout option.Option[time.Duration] // None means use default
}`}
</CodeCard>

### Idiom 4: Keep Packages Focused

<CodeCard file="focused-packages.go">
{`// Go idiom: small, focused packages
package user

// Don't mix everything
// ❌ user/fp.go, user/imperative.go

// ✅ Use fp-go where it helps
func (s *Service) GetUser(id string) (User, error) {
    // Use fp-go internally
    result := s.fetchAndValidate(id)
    
    // Return Go-style
    return result.Fold(
        func(err error) (User, error) { return User{}, err },
        func(u User) (User, error) { return u, nil },
    )
}`}
</CodeCard>

</Section>

<Section id="testing" number="06" title="Testing with" titleAccent="fp-go.">

### Test Pure Functions

<CodeCard file="test-pure.go">
{`// Pure functions are easy to test
func TestNormalizeUser(t *testing.T) {
    user := User{Name: "  JOHN  "}
    normalized := normalizeUser(user)
    
    assert.Equal(t, "john", normalized.Name)
}`}
</CodeCard>

### Test Pipelines

<CodeCard file="test-pipelines.go">
{`// Test the pipeline structure
func TestProcessUser(t *testing.T) {
    // Mock dependencies
    fetchUser = func(id string) result.Result[User] {
        return result.Ok(User{ID: id})
    }
    
    // Test pipeline
    result := processUser("123")
    
    assert.True(t, result.IsOk())
}`}
</CodeCard>

### Test with Table Tests (Go Idiom)

<CodeCard file="table-tests.go">
{`func TestValidateUser(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {"valid user", User{Email: "test@example.com"}, false},
        {"invalid email", User{Email: "invalid"}, true},
        {"empty email", User{Email: ""}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validateUser(tt.user)
            
            if tt.wantErr {
                assert.True(t, result.IsErr())
            } else {
                assert.True(t, result.IsOk())
            }
        })
    }
}`}
</CodeCard>

</Section>

<Section id="performance" number="07" title="Performance" titleAccent="considerations.">

### Measure First

<CodeCard file="benchmark.go">
{`// Don't assume - benchmark
func BenchmarkStandard(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = standardApproach()
    }
}

func BenchmarkFPGo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = fpgoApproach()
    }
}

// Then decide based on results`}
</CodeCard>

### Use Idiomatic Packages

<CodeCard file="idiomatic.go">
{`// fp-go provides idiomatic (faster) versions
import "github.com/IBM/fp-go/v2/array/idiomatic"

// 2-32x faster for array operations
filtered := idiomatic.Filter(predicate)(array)`}
</CodeCard>

### Optimize Hot Paths

<CodeCard file="optimize-hot-paths.go">
{`// Cold path: use fp-go for clarity
func processConfig(config Config) result.Result[Config] {
    return function.Pipe3(
        result.Ok(config),
        result.Chain(validate),
        result.Chain(enrich),
    )
}

// Hot path: optimize
func processRequest(req Request) Response {
    // Direct implementation for performance
    if !isValid(req) {
        return errorResponse
    }
    return successResponse
}`}
</CodeCard>

</Section>

<Section id="team-adoption" number="08" title="Team" titleAccent="adoption.">

### Start Small

<CodeCard file="start-small.go">
{`// Phase 1: Use Result for error handling
func fetchUser(id string) result.Result[User] {
    user, err := db.FindByID(id)
    return result.FromGoError(user, err)
}

// Phase 2: Add Option for optional values
func findConfig(key string) option.Option[string] {
    // ...
}

// Phase 3: Introduce composition
func processUser(id string) result.Result[User] {
    return function.Pipe2(
        fetchUser(id),
        result.Chain(validateUser),
    )
}`}
</CodeCard>

### Provide Training

<CodeCard file="training.go">
{`// Create examples for your team
// examples/user_processing.go

// Example 1: Simple Result usage
func example1() {
    result := fetchUser("123")
    result.Fold(
        func(err error) { fmt.Println("Error:", err) },
        func(user User) { fmt.Println("User:", user) },
    )
}

// Example 2: Chaining operations
func example2() {
    result := function.Pipe2(
        fetchUser("123"),
        result.Chain(validateUser),
    )
}`}
</CodeCard>

### Set Guidelines

<Callout type="info">
**Document when to use fp-go:**

Use fp-go for:
- Complex error handling (3+ sequential operations)
- Optional values (instead of nil pointers)
- Data transformations (map, filter, reduce)

Use standard Go for:
- Simple operations
- Performance-critical code
- Public APIs (convert at boundary)
</Callout>

</Section>
