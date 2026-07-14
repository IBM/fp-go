---
sidebar_position: 9
title: v1 and v2 Interop
hide_title: true
description: Learn how to run fp-go v1 and v2 side-by-side during gradual migration, with patterns for bridging between versions.
---

<PageHeader
  eyebrow="Migration · 03 / 03"
  title="v1 and v2"
  titleAccent="Interop."
  lede="Run fp-go v1 and v2 side-by-side during gradual migration. Learn patterns for bridging between versions safely and efficiently."
  meta={[
    {label: '// Difficulty', value: 'Intermediate'},
    {label: '// Use case', value: 'Gradual migration'},
    {label: '// Reading time', value: '12 min · 6 sections'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Key pattern" prose value={<>Convert at <em>boundaries</em>, not throughout.</>} variant="up" />
  <TLDRCard label="// Bridge package" prose value={<>Isolate all <em>conversion logic</em>.</>} />
  <TLDRCard label="// Feature flags" prose value={<>Control which <em>version</em> to use.</>} />
</TLDR>

<Section id="why" number="01" title="Why run both" titleAccent="versions?">

### Use Cases

**Gradual Migration:**
- Migrate module by module
- Test incrementally
- Minimize risk

**Legacy Support:**
- Keep old code working
- Add new features with v2
- Eventual migration

**Team Coordination:**
- Different teams migrate at different speeds
- Shared libraries need both versions
- Smooth transition

</Section>

<Section id="setup" number="02" title="Setup" titleAccent="guide.">

### 1. Install Both Versions

<CodeCard file="install.sh" lang="bash">
{`# Install v1 and v2
go get github.com/IBM/fp-go
go get github.com/IBM/fp-go/v2

# Verify go.mod
cat go.mod`}
</CodeCard>

Your `go.mod` should have:

<CodeCard file="go.mod">
{`module myapp

go 1.24

require (
    github.com/IBM/fp-go v1.x.x
    github.com/IBM/fp-go/v2 v2.x.x
)`}
</CodeCard>

### 2. Use Import Aliases

<CodeCard file="imports.go">
{`import (
    // v1 imports with v1 prefix
    v1either "github.com/IBM/fp-go/either"
    v1option "github.com/IBM/fp-go/option"
    v1ioeither "github.com/IBM/fp-go/ioeither"
    
    // v2 imports with v2 prefix
    v2either "github.com/IBM/fp-go/v2/either"
    v2result "github.com/IBM/fp-go/v2/result"
    v2ioresult "github.com/IBM/fp-go/v2/ioresult"
)`}
</CodeCard>

### 3. Organize Code

<CodeCard file="structure.txt">
{`myapp/
├── legacy/          # v1 code
│   ├── user.go
│   └── auth.go
├── new/             # v2 code
│   ├── api.go
│   └── service.go
└── bridge/          # Interop code
    └── convert.go`}
</CodeCard>

</Section>

<Section id="conversion-patterns" number="03" title="Conversion" titleAccent="patterns.">

### Pattern 1: Either v1 → Result v2

Convert v1 Either to v2 Result:

<CodeCard file="either-to-result.go">
{`// Conversion function
func EitherToResult[A any](e v1either.Either[error, A]) v2result.Result[A] {
    return v1either.Fold(
        // Left (error) → Err
        func(err error) v2result.Result[A] {
            return v2result.Err[A](err)
        },
        // Right (value) → Ok
        func(val A) v2result.Result[A] {
            return v2result.Ok(val)
        },
    )(e)
}

// Usage
func legacyFunction() v1either.Either[error, User] {
    // v1 code
    return v1either.Right[error](User{ID: "123"})
}

func newFunction() v2result.Result[User] {
    v1Result := legacyFunction()
    return EitherToResult(v1Result)
}`}
</CodeCard>

### Pattern 2: Result v2 → Either v1

Convert v2 Result to v1 Either:

<CodeCard file="result-to-either.go">
{`// Conversion function using Fold
func ResultToEither[A any](r v2result.Result[A]) v1either.Either[error, A] {
    return r.Fold(
        // Err → Left
        func(err error) v1either.Either[error, A] {
            return v1either.Left[A](err)
        },
        // Ok → Right
        func(val A) v1either.Either[error, A] {
            return v1either.Right[error](val)
        },
    )
}

// Usage
func newFunction() v2result.Result[User] {
    // v2 code
    return v2result.Ok(User{ID: "123"})
}

func legacyFunction() v1either.Either[error, User] {
    v2Result := newFunction()
    return ResultToEither(v2Result)
}`}
</CodeCard>

### Pattern 3: Option v1 ↔ Option v2

Options are similar in both versions:

<CodeCard file="option-conversion.go">
{`// v1 Option → v2 Option
func OptionV1ToV2[A any](opt v1option.Option[A]) v2option.Option[A] {
    return v1option.Fold(
        func() v2option.Option[A] {
            return v2option.None[A]()
        },
        func(val A) v2option.Option[A] {
            return v2option.Some(val)
        },
    )(opt)
}

// v2 Option → v1 Option
func OptionV2ToV1[A any](opt v2option.Option[A]) v1option.Option[A] {
    return v2option.Fold(
        func() v1option.Option[A] {
            return v1option.None[A]()
        },
        func(val A) v1option.Option[A] {
            return v1option.Some(val)
        },
    )(opt)
}`}
</CodeCard>

### Pattern 4: IOEither v1 → IOResult v2

<CodeCard file="ioeither-to-ioresult.go">
{`// Conversion function
func IOEitherToIOResult[A any](
    ioe v1ioeither.IOEither[error, A],
) v2ioresult.IOResult[A] {
    return func() v2result.Result[A] {
        // Execute v1 IOEither
        e := ioe()
        
        // Convert Either to Result
        return EitherToResult(e)
    }
}

// Usage
func legacyReadFile(path string) v1ioeither.IOEither[error, []byte] {
    return func() v1either.Either[error, []byte] {
        data, err := os.ReadFile(path)
        if err != nil {
            return v1either.Left[[]byte](err)
        }
        return v1either.Right[error](data)
    }
}

func newReadFile(path string) v2ioresult.IOResult[[]byte] {
    v1IO := legacyReadFile(path)
    return IOEitherToIOResult(v1IO)
}`}
</CodeCard>

### Pattern 5: IOResult v2 → IOEither v1

<CodeCard file="ioresult-to-ioeither.go">
{`// Conversion function
func IOResultToIOEither[A any](
    ior v2ioresult.IOResult[A],
) v1ioeither.IOEither[error, A] {
    return func() v1either.Either[error, A] {
        // Execute v2 IOResult
        r := ior()
        
        // Convert Result to Either
        return ResultToEither(r)
    }
}

// Usage
func newFetchData(url string) v2ioresult.IOResult[Data] {
    return func() v2result.Result[Data] {
        // v2 implementation
        return v2result.Ok(Data{})
    }
}

func legacyFetchData(url string) v1ioeither.IOEither[error, Data] {
    v2IO := newFetchData(url)
    return IOResultToIOEither(v2IO)
}`}
</CodeCard>

</Section>

<Section id="real-world-examples" number="04" title="Real-world" titleAccent="examples.">

### Example 1: HTTP Handler Bridge

<CodeCard file="http-bridge.go">
{`// Legacy v1 service
type LegacyUserService struct{}

func (s *LegacyUserService) GetUser(id string) v1ioeither.IOEither[error, User] {
    return func() v1either.Either[error, User] {
        // v1 implementation
        user := User{ID: id, Name: "John"}
        return v1either.Right[error](user)
    }
}

// New v2 service
type UserService struct {
    legacy *LegacyUserService
}

func (s *UserService) GetUser(id string) v2ioresult.IOResult[User] {
    // Bridge to legacy service
    v1Result := s.legacy.GetUser(id)
    return IOEitherToIOResult(v1Result)
}

// HTTP handler using v2
func HandleGetUser(service *UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Query().Get("id")
        
        // Use v2 API
        result := service.GetUser(id)()
        
        result.Fold(
            func(err error) {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            },
            func(user User) {
                json.NewEncoder(w).Encode(user)
            },
        )
    }
}`}
</CodeCard>

### Example 2: Database Layer Migration

<CodeCard file="database-bridge.go">
{`// Legacy v1 repository
type LegacyUserRepo struct {
    db *sql.DB
}

func (r *LegacyUserRepo) FindByID(id string) v1ioeither.IOEither[error, User] {
    return func() v1either.Either[error, User] {
        var user User
        err := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
        if err != nil {
            return v1either.Left[User](err)
        }
        return v1either.Right[error](user)
    }
}

// New v2 repository
type UserRepo struct {
    db *sql.DB
}

func (r *UserRepo) FindByID(id string) v2ioresult.IOResult[User] {
    return func() v2result.Result[User] {
        var user User
        err := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
        return v2result.FromGoError(user, err)
    }
}

// Bridge service using both
type UserService struct {
    legacyRepo *LegacyUserRepo
    newRepo    *UserRepo
    useV2      bool // Feature flag
}

func (s *UserService) GetUser(id string) v2ioresult.IOResult[User] {
    if s.useV2 {
        return s.newRepo.FindByID(id)
    }
    
    // Bridge to v1
    v1Result := s.legacyRepo.FindByID(id)
    return IOEitherToIOResult(v1Result)
}`}
</CodeCard>

### Example 3: Shared Library

<CodeCard file="shared-library.go">
{`// shared/types.go - Common types
package shared

type User struct {
    ID    string
    Name  string
    Email string
}

// shared/v1/user.go - v1 API
package v1

import (
    v1either "github.com/IBM/fp-go/either"
    "myapp/shared"
)

func ValidateUser(user shared.User) v1either.Either[error, shared.User] {
    // v1 implementation
    return v1either.Right[error](user)
}

// shared/v2/user.go - v2 API
package v2

import (
    v2result "github.com/IBM/fp-go/v2/result"
    "myapp/shared"
)

func ValidateUser(user shared.User) v2result.Result[shared.User] {
    // v2 implementation
    return v2result.Ok(user)
}

// shared/bridge/user.go - Conversion
package bridge

import (
    v1either "github.com/IBM/fp-go/either"
    v2result "github.com/IBM/fp-go/v2/result"
    "myapp/shared"
    sharedv1 "myapp/shared/v1"
    sharedv2 "myapp/shared/v2"
)

// Use v1 from v2 code
func ValidateUserV1AsV2(user shared.User) v2result.Result[shared.User] {
    v1Result := sharedv1.ValidateUser(user)
    return EitherToResult(v1Result)
}

// Use v2 from v1 code
func ValidateUserV2AsV1(user shared.User) v1either.Either[error, shared.User] {
    v2Result := sharedv2.ValidateUser(user)
    return ResultToEither(v2Result)
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="05" title="Best" titleAccent="practices.">

### 1. Isolate Conversion Logic

Create a dedicated bridge package:

<CodeCard file="bridge-package.go">
{`// bridge/convert.go
package bridge

// All conversion functions in one place
func EitherToResult[A any](e v1either.Either[error, A]) v2result.Result[A] { ... }
func ResultToEither[A any](r v2result.Result[A]) v1either.Either[error, A] { ... }
func IOEitherToIOResult[A any](ioe v1ioeither.IOEither[error, A]) v2ioresult.IOResult[A] { ... }
// etc.`}
</CodeCard>

### 2. Use Feature Flags

Control which version to use:

<CodeCard file="feature-flags.go">
{`type Config struct {
    UseV2UserService bool
    UseV2AuthService bool
}

func NewUserService(cfg Config, v1Svc *V1UserService, v2Svc *V2UserService) UserService {
    if cfg.UseV2UserService {
        return v2Svc
    }
    return &BridgedUserService{v1: v1Svc}
}`}
</CodeCard>

### 3. Test Both Paths

<CodeCard file="test-both.go">
{`func TestUserService(t *testing.T) {
    t.Run("v1 implementation", func(t *testing.T) {
        svc := NewUserService(Config{UseV2UserService: false}, v1Svc, v2Svc)
        // Test v1 path
    })
    
    t.Run("v2 implementation", func(t *testing.T) {
        svc := NewUserService(Config{UseV2UserService: true}, v1Svc, v2Svc)
        // Test v2 path
    })
    
    t.Run("bridge conversion", func(t *testing.T) {
        // Test conversion functions
    })
}`}
</CodeCard>

### 4. Document Version Usage

<CodeCard file="documentation.go">
{`// Package user provides user management.
//
// This package is in transition from fp-go v1 to v2.
// Current status:
// - GetUser: v2 ✅
// - CreateUser: v1 (migrating)
// - UpdateUser: v1 (migrating)
// - DeleteUser: v1 (not started)
//
// See MIGRATION.md for details.
package user`}
</CodeCard>

### 5. Monitor Performance

<CodeCard file="benchmark.go">
{`func BenchmarkBridge(b *testing.B) {
    b.Run("v1 direct", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = v1Function()
        }
    })
    
    b.Run("v2 direct", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = v2Function()
        }
    })
    
    b.Run("v1 via bridge", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = IOResultToIOEither(v2Function())
        }
    })
}`}
</CodeCard>

</Section>

<Section id="common-pitfalls" number="06" title="Common" titleAccent="pitfalls.">

### Pitfall 1: Nested Conversions

<Compare>
<CompareCol kind="bad">

<CodeCard file="nested-conversions.go">
{`// Converting back and forth
v2Result := EitherToResult(v1Either)
v1Either2 := ResultToEither(v2Result)
v2Result2 := EitherToResult(v1Either2)
// Performance overhead!`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="boundary-conversion.go">
{`// Convert once at boundaries
func processData(input Data) v2result.Result[Data] {
    // Do all processing in v2
    return v2result.Ok(input)
}

// Convert only when interfacing with v1 code
func legacyAPI(input Data) v1either.Either[error, Data] {
    v2Result := processData(input)
    return ResultToEither(v2Result) // Convert once
}`}
</CodeCard>

</CompareCol>
</Compare>

### Pitfall 2: Forgetting Error Context

<Compare>
<CompareCol kind="bad">

<CodeCard file="lost-error.go">
{`func ResultToEither[A any](r v2result.Result[A]) v1either.Either[error, A] {
    if r.IsOk() {
        return v1either.Right[error](r.GetOrElse(func() A { var zero A; return zero }))
    }
    // Lost the actual error!
    return v1either.Left[A](errors.New("error"))
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="preserve-error.go">
{`func ResultToEither[A any](r v2result.Result[A]) v1either.Either[error, A] {
    return r.Fold(
        func(err error) v1either.Either[error, A] {
            return v1either.Left[A](err) // Preserve error
        },
        func(val A) v1either.Either[error, A] {
            return v1either.Right[error](val)
        },
    )
}`}
</CodeCard>

</CompareCol>
</Compare>

### Pitfall 3: Type Inference Issues

<Compare>
<CompareCol kind="bad">

<CodeCard file="inference-fail.go">
{`// Compiler can't infer types
result := EitherToResult(v1Either)
// Error: cannot infer A`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="explicit-types.go">
{`// Specify type explicitly
result := EitherToResult[User](v1Either)

// Or use type annotation
var result v2result.Result[User] = EitherToResult(v1Either)`}
</CodeCard>

</CompareCol>
</Compare>

</Section>
