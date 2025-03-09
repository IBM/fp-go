# Functional programming library for golang V2

Go 1.24 introduces [generic type aliases](https://github.com/golang/go/issues/46477) which are leveraged by V2.

## ⚠️ Breaking Changes

- use of [generic type aliases](https://github.com/golang/go/issues/46477) which requires [go1.24](https://tip.golang.org/doc/go1.24)
- order of generic type arguments adjusted such that types that _cannot_ be inferred by the method argument come first, e.g. in the `Ap` methods
- monadic operations for `Pair` operate on the second argument, to be compatible with the [Haskell](https://hackage.haskell.org/package/TypeCompose-0.9.14/docs/Data-Pair.html) definition

## Simplifications

- use type aliases to get rid of namespace imports for type declarations, e.g. instead of

```go
import (
    ET "github.com/IBM/fp-go/v2/either"
)

func doSth() ET.Either[error, string] {
    ...
}
```

you can declare your type once 

```go
import (
    "github.com/IBM/fp-go/v2/either"
)

type Either[A any] = either.Either[error, A]
```

and then use it across your codebase

```go
func doSth() Either[string] {
    ...
}
```

- library implementation does no longer need to use the `generic` subpackage, this simplifies reading and understanding of the code