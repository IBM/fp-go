---
name: fp-go-pr-review
description: Use this skill when reviewing pull requests for fp-go code (github.com/IBM/fp-go/v2). Trigger on mentions of PR review, code review, pull request validation, fp-go best practices validation, functional programming review, or when the user asks to review changes on a PR branch. This skill validates that changes follow fp-go conventions including data-last composition, point-free style, proper monad usage, lens patterns, and idiomatic functional patterns.
---

# fp-go PR Review

## Overview

This skill assists with reviewing pull requests that use the fp-go library (github.com/IBM/fp-go/v2). It validates that code changes follow fp-go best practices and functional programming conventions. Requires Go 1.24+ for generic type alias support.

## When to Use This Skill

- Reviewing pull requests with fp-go code
- Validating that changes follow fp-go best practices
- Checking for common fp-go anti-patterns
- Ensuring proper functional composition patterns
- Verifying correct monad usage and error handling

## Review Checklist

### 1. Import Path Validation

**Rule**: All imports MUST use `github.com/IBM/fp-go/v2/...`, never `github.com/IBM/fp-go/...` (v1).

**Check for**:
```go
// ❌ WRONG - v1 import
import "github.com/IBM/fp-go/option"

// ✅ CORRECT - v2 import
import O "github.com/IBM/fp-go/v2/option"
```

**Severity**: Critical — v1 and v2 are incompatible

### 2. Data-Last Principle

**Rule**: All fp-go operations use data-last. The data being transformed is always the last argument.

**Check for**:
```go
// ❌ WRONG - data-first
option.Map(myOption, transformFunc)

// ✅ CORRECT - data-last
option.Map(transformFunc)(myOption)

// ✅ CORRECT - in pipeline
F.Pipe2(
    myOption,
    O.Map(transformFunc),
    O.GetOrElse(F.Constant("default")),
)
```

**Severity**: High — breaks composition

### 3. Point-Free Style

**Rule**: Prefer composing named functions with `Flow` and `Pipe` over inline anonymous functions.

**Check for**:
```go
// ❌ AVOID - unnecessary lambda wrapping
pipeline := F.Flow2(
    func(s string) O.Option[int] { return O.FromPredicate(S.IsNonEmpty)(s) },
    func(o O.Option[int]) int { return O.GetOrElse(F.Constant(0))(o) },
)

// ✅ CORRECT - point-free composition
pipeline := F.Flow2(
    O.FromPredicate(S.IsNonEmpty),
    O.GetOrElse(F.Constant(0)),
)

// ❌ AVOID - inline comparison
A.Filter(func(x int) bool { return x > 18 })

// ✅ CORRECT - use numeric combinator
A.Filter(N.MoreThan(18))
```

**Severity**: Medium — impacts readability and maintainability

### 4. Prefer Result Over Either

**Rule**: Use `Result[A]` (which is `Either[error, A]`) when the error type is Go's `error`. Reserve `Either` for custom error types.

**Check for**:
```go
// ❌ AVOID - Either with error
func fetchData() E.Either[error, Data] { ... }

// ✅ CORRECT - use Result
func fetchData() R.Result[Data] { ... }

// ✅ CORRECT - Either with custom error type
func validate() E.Either[ValidationError, Data] { ... }
```

**Severity**: Medium — Result is more idiomatic for Go errors

### 5. IO Laziness

**Rule**: IO values are lazy (`IO[A]` is `func() A`). They must be called with `()` to execute.

**Check for**:
```go
// ❌ WRONG - forgot to execute
result := readConfig("config.json")  // returns IO[Config], not Config

// ✅ CORRECT - execute with ()
result := readConfig("config.json")()

// ❌ WRONG - in ReaderIOResult, forgot inner ()
value, err := pipeline(ctx)  // returns func() Result[A]

// ✅ CORRECT - execute both context and IO
value, err := pipeline(ctx)()
```

**Severity**: Critical — code won't execute

### 6. Monad Selection

**Rule**: Use the simplest monad that covers your needs. Escalate only when necessary.

**Check for**:
```go
// ❌ AVOID - using ReaderIOResult for pure computation
func processUsers(users []User) RIO.ReaderIOResult[string] {
    return F.Pipe2(
        RIO.Of[context.Context](users),
        RIO.Map[context.Context, []User, string](pureTransform),
    )
}

// ✅ CORRECT - pure computation, no monad needed
func processUsers() func([]User) string {
    return F.Flow2(
        A.FilterMap(toAdultName()),
        A.Intercalate(S.Monoid)(","),
    )
}
```

**Severity**: Medium — unnecessary complexity

**Escalation path**: `Option` → `Result` → `IOResult` → `ReaderIOResult` → `Effect`

### 7. Effect vs ReaderIOResult

**Rule**: Use `Effect[C, A]` for services with typed dependencies. Use `ReaderIOResult` only when you truly only need `context.Context`.

**Check for**:
```go
// ❌ AVOID - stuffing deps into context.Context
func fetchUser(id int) RIO.ReaderIOResult[User] {
    return func(ctx context.Context) func() R.Result[User] {
        db := ctx.Value("db").(DBClient)  // runtime type assertion
        // ...
    }
}

// ✅ CORRECT - typed dependencies with Effect
type Deps struct {
    DB     DBClient
    Logger Logger
}

func fetchUser(id int) EF.Effect[Deps, User] {
    return EF.Asks(func(deps Deps) EF.ReaderIOResult[User] {
        // deps.DB is compile-time checked
        return queryUser(deps.DB, id)
    })
}
```

**Severity**: High — type safety and testability

### 8. Lifting Go Functions

**Rule**: Use `Eitherize1`..`EitherizeN` to lift Go functions returning `(T, error)` into Result.

**Check for**:
```go
// ❌ AVOID - manual error handling
func parseNumber(s string) R.Result[int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return R.Left[int](err)
    }
    return R.Right[error](n)
}

// ✅ CORRECT - use Eitherize
var parseNumber = R.Eitherize1(strconv.Atoi)

// ✅ CORRECT - in pipeline
pipeline := F.Flow2(
    R.Eitherize1(strconv.Atoi),
    R.Map(N.Mul(2)),
)
```

**Severity**: Medium — reduces boilerplate

### 9. Do-Notation with Lenses

**Rule**: Use lenses with `Bind`/`ApS` instead of manual setter functions.

**Check for**:
```go
// ❌ AVOID - manual setter functions
func setUser(u User) func(State) State {
    return func(s State) State { s.User = u; return s }
}

pipeline := F.Pipe2(
    RIO.Do(State{}),
    RIO.Bind(setUser, fetchUser),
)

// ✅ CORRECT - use lens
var userLens = L.MakeLens(
    func(s State) User { return s.User },
    func(s State, u User) State { s.User = u; return s },
)

pipeline := F.Pipe2(
    RIO.Do(State{}),
    RIO.Bind(userLens.Set, fetchUser),
)

// ✅ EVEN BETTER - use code generation
//go:generate go run github.com/IBM/fp-go/v2/main lens --dir . --filename gen_lens.go

// fp-go:Lens
type State struct {
    User User
}

// Then use generated lens
lenses := MakeStateLenses()
pipeline := F.Pipe2(
    RIO.Do(State{}),
    RIO.Bind(lenses.User.Set, fetchUser),
)
```

**Severity**: Medium — maintainability and consistency

### 10. Bind vs ApS

**Rule**: Use `Bind` when the step depends on accumulated state; use `ApS` when steps are independent.

**Check for**:
```go
// ❌ WRONG - using Bind when steps are independent
pipeline := F.Pipe2(
    RIO.Do(Summary{}),
    RIO.Bind(userLens.Set, func(_ Summary) RIO.ReaderIOResult[User] {
        return fetchUser(42)  // doesn't use state
    }),
    RIO.Bind(weatherLens.Set, func(_ Summary) RIO.ReaderIOResult[Weather] {
        return fetchWeather("NYC")  // doesn't use state
    }),
)

// ✅ CORRECT - use ApS for independent steps
pipeline := F.Pipe2(
    RIO.Do(Summary{}),
    RIO.ApS(userLens.Set, fetchUser(42)),
    RIO.ApS(weatherLens.Set, fetchWeather("NYC")),
)

// ✅ CORRECT - use Bind when dependent
pipeline := F.Pipe2(
    RIO.Do(Pipeline{}),
    RIO.Bind(userLens.Set, func(_ Pipeline) RIO.ReaderIOResult[User] {
        return fetchUser(42)
    }),
    RIO.Bind(configLens.Set, F.Flow2(userLens.Get, fetchConfigForUser)),
)
```

**Severity**: Medium — semantic clarity

### 11. TraverseArray Usage

**Rule**: Use `TraverseArray` to process slices monadically, not manual loops with error accumulation.

**Check for**:
```go
// ❌ AVOID - manual loop with error handling
func fetchAll(ids []int) RIO.ReaderIOResult[[]User] {
    return func(ctx context.Context) func() R.Result[[]User] {
        return func() R.Result[[]User] {
            users := make([]User, 0, len(ids))
            for _, id := range ids {
                user, err := fetchUser(id)(ctx)()
                if err != nil {
                    return R.Left[[]User](err)
                }
                users = append(users, user)
            }
            return R.Right[error](users)
        }
    }
}

// ✅ CORRECT - use TraverseArray
func fetchAll(ids []int) RIO.ReaderIOResult[[]User] {
    return RIO.TraverseArray(fetchUser)(ids)
}
```

**Severity**: High — idiomatic functional pattern

### 12. Logging Side Effects

**Rule**: Use `ChainFirstIOK` with `IO.Logf` for logging without breaking the pipeline.

**Check for**:
```go
// ❌ AVOID - breaking the pipeline for logging
pipeline := F.Pipe2(
    fetchUser(42),
    RIO.Chain(func(user User) RIO.ReaderIOResult[User] {
        log.Printf("Fetched user: %v", user)
        return RIO.Of[context.Context](user)
    }),
)

// ✅ CORRECT - use ChainFirstIOK
pipeline := F.Pipe2(
    fetchUser(42),
    RIO.ChainFirstIOK(IO.Logf[User]("Fetched user: %v")),
)

// ✅ CORRECT - structured logging with TapSLog
pipeline := F.Pipe2(
    fetchUser(42),
    RIO.TapSLog[User]("User fetched"),
)
```

**Severity**: Low — code quality

### 13. Prefer Functions Over Variables

**Rule**: Wrap pipeline results in functions, not package-level vars.

**Check for**:
```go
// ❌ WRONG - var is allocated even if never called
var processUser = F.Flow2(getName, strings.ToUpper)

// ✅ CORRECT - zero cost until called
func processUser() func(User) string {
    return F.Flow2(getName, strings.ToUpper)
}
```

**Severity**: Low — performance and dead code elimination

### 14. Type Parameter Order

**Rule**: Non-inferrable type parameters come first. The compiler infers trailing params from arguments.

**Check for**:
```go
// ❌ WRONG - compiler can't infer B
O.Map[string, int](toLength)

// ✅ CORRECT - B comes first, A is inferred
O.Map[int](toLength)

// ✅ CORRECT - let compiler infer both when possible
O.Map(toLength)
```

**Severity**: Low — compilation errors or verbosity

### 15. Lens Composition

**Rule**: Use `Compose`/`ComposeRef` for nested struct access, not manual chaining.

**Check for**:
```go
// ❌ AVOID - manual nested access
func getStreetName(p Person) string {
    if p.Address != nil && p.Address.Street != nil {
        return p.Address.Street.Name
    }
    return ""
}

// ✅ CORRECT - compose lenses
streetNameInPerson := F.Pipe2(
    personAddressLens,
    LO.Compose[Person, *Street](defaultAddress)(addressStreetLens),
    LO.ComposeOption[Person, string](defaultStreet)(streetNameLens),
)

name := streetNameInPerson.Get(person)  // Option[string]
```

**Severity**: Medium — immutability and composability

## Review Process

### Step 1: Obtain Git Diff

Get the changes on the PR branch relative to main:

```bash
git diff main...HEAD
```

To list only changed file paths:

```bash
git diff --name-only main...HEAD
```

For a GitHub PR, fetch it first:

```bash
gh pr checkout <PR-number>
git diff main...HEAD
```

### Step 2: Analyze Changes

For each modified file:
1. Check import paths (v2 requirement)
2. Validate data-last usage
3. Check for point-free style opportunities
4. Verify monad selection appropriateness
5. Check IO execution (trailing `()`)
6. Validate error handling patterns
7. Check lens usage in do-notation
8. Verify Bind vs ApS usage
9. Look for TraverseArray opportunities
10. Check logging patterns

### Step 3: Submit Findings

Post a review comment on the GitHub PR:

```bash
gh pr review <PR-number> --comment -b "$(cat <<'EOF'
## fp-go Review

**Overall**: Needs Changes

### Critical
- ❌ ...

### High
- ⚠️ ...

### Recommendations
1. ...
EOF
)"
```

For inline annotations on specific lines, use:

```bash
gh api repos/{owner}/{repo}/pulls/<PR-number>/comments \
  -f body="Replace inline lambda with point-free: \`F.Flow2(O.FromPredicate(S.IsNonEmpty), O.GetOrElse(F.Constant(0)))\`" \
  -f commit_id="$(git rev-parse HEAD)" \
  -f path="src/user/handler.go" \
  -F line=42 \
  -f side=RIGHT
```

Alternatively, use the `/code-review --comment` skill to post inline PR annotations automatically.

## Common Issue Categories

| Category | Type | Example |
|----------|------|---------|
| maintainability | dry-principle-violation | Inline lambdas instead of point-free |
| maintainability | naming-intent-review | Non-descriptive variable names |
| functionality | error-handling-review | Missing error propagation |
| performance | inefficient-algorithm | Manual loops instead of TraverseArray |
| style | style-consistency-check | Inconsistent import aliases |
| security | sensitive-data-logging | Logging sensitive information |

## Severity Guidelines

- **Critical**: Code won't compile or execute (wrong import path, missing `()`)
- **High**: Type safety issues, incorrect monad usage, breaks composition
- **Medium**: Readability, maintainability, non-idiomatic patterns
- **Low**: Style preferences, minor optimizations

## Example Review Comments

### Import Path Issue

> **Severity**: Critical
> **Issue**: Using v1 import path
>
> The import `github.com/IBM/fp-go/option` is the v1 path. All imports must use v2:
> `github.com/IBM/fp-go/v2/option`
>
> v1 and v2 are incompatible. This will cause compilation errors or runtime issues.

### Point-Free Style

> **Severity**: Medium
> **Issue**: Unnecessary lambda wrapping
>
> This inline lambda can be replaced with point-free composition:
>
> Current:
> ```go
> option.Filter(func(s string) bool { return s != "" })
> ```
>
> Suggested:
> ```go
> option.Filter(S.IsNonEmpty)
> ```
>
> Point-free style is more readable and idiomatic in fp-go.

### Monad Selection

> **Severity**: Medium
> **Issue**: Unnecessary monad for pure computation
>
> This function uses `ReaderIOResult` but performs only pure transformations without IO or context:
>
> ```go
> func processUsers(users []User) RIO.ReaderIOResult[string] {
>     return F.Pipe2(
>         RIO.Of[context.Context](users),
>         RIO.Map[context.Context, []User, string](pureTransform),
>     )
> }
> ```
>
> Suggested:
> ```go
> func processUsers() func([]User) string {
>     return F.Flow2(
>         A.FilterMap(toAdultName()),
>         A.Intercalate(S.Monoid)(","),
>     )
> }
> ```
>
> Use the simplest abstraction that covers your needs.

## Integration with Other Skills

This skill can reference and include:
- `fp-go` — Core fp-go patterns and best practices
- `fp-go-pipe-flow` — Pipe/Flow composition patterns
- `fp-go-http` — HTTP request patterns
- `fp-go-logging` — Logging patterns
- `fp-go-lens` — Lens and optics patterns

## Automated Checks

When reviewing, automatically check for:

1. ✅ All imports use `v2` path
2. ✅ No data-first function calls
3. ✅ IO values are executed with `()`
4. ✅ `Result` used instead of `Either[error, A]`
5. ✅ Point-free style where applicable
6. ✅ Appropriate monad selection
7. ✅ Lenses used in do-notation
8. ✅ `Bind` vs `ApS` used correctly
9. ✅ `TraverseArray` for slice processing
10. ✅ `ChainFirstIOK` for logging

## Output Format

Provide a summary with:
1. **Overall Assessment**: Pass/Needs Changes/Blocked
2. **Critical Issues**: Count and list
3. **High Priority Issues**: Count and list
4. **Medium Priority Issues**: Count and list
5. **Low Priority Issues**: Count and list
6. **Positive Observations**: What was done well
7. **Recommendations**: Suggested improvements

## Example Summary

```markdown
## PR Review Summary

**Overall Assessment**: Needs Changes

### Critical Issues (1)
- ❌ Using v1 import path in `user/handler.go:5`

### High Priority Issues (2)
- ⚠️ Missing IO execution in `config/loader.go:42`
- ⚠️ Manual error handling instead of Eitherize in `api/client.go:78`

### Medium Priority Issues (3)
- 💡 Inline lambda instead of point-free in `user/service.go:23`
- 💡 Using ReaderIOResult for pure computation in `utils/format.go:15`
- 💡 Manual setter instead of lens in `state/pipeline.go:56`

### Low Priority Issues (1)
- 📝 Inconsistent import alias in `handler/http.go:8`

### Positive Observations
- ✅ Excellent use of TraverseArray for parallel requests
- ✅ Proper Effect usage with typed dependencies
- ✅ Good lens composition for nested struct access

### Recommendations
1. Update all imports to v2 path
2. Add trailing `()` to execute IO values
3. Consider using `R.Eitherize1` for Go function lifting
4. Refactor pure computations to use Flow instead of ReaderIOResult
```

## References

- [fp-go v2 Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2)
- [fp-go GitHub Repository](https://github.com/IBM/fp-go)
- [Functional Programming in Go](https://github.com/IBM/fp-go/blob/main/README.md)