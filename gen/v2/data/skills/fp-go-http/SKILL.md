---
name: fp-go-http
description: Use this skill when making HTTP requests in fp-go using the ReaderIOResult-based HTTP client (github.com/IBM/fp-go/v2/context/readerioresult/http). Trigger on mentions of fp-go HTTP, MakeClient, MakeGetRequest, MakeRequest, ReadJSON, ReadText, ReadAll, ReadFullResponse, the HTTP request builder (WithURL, WithJSON, WithBearer, WithHeader, WithQueryArg), parallel requests with TraverseArray or TraverseTuple2, or building context-aware, composable HTTP pipelines that propagate errors through the Result monad.
---

# fp-go HTTP Requests

## Overview

fp-go wraps `net/http` in the `ReaderIOResult` monad, giving you composable, context-aware HTTP operations with automatic error propagation. The core package is:

```
github.com/IBM/fp-go/v2/context/readerioresult/http
```

All HTTP operations are lazy — they describe what to do but do not execute until you call the resulting function with a `context.Context`.

**Executing a pipeline.** `ReaderIOResult[A]` is `func(context.Context) func() Result[A]`, so running it takes two calls and yields **one** value:

```go
res := pipeline(ctx)()            // Result[A] — a single value, NOT (A, error)
value, err := R.Unwrap(res)       // R = github.com/IBM/fp-go/v2/result
```

Examples below use the shorthand `value, err := R.Unwrap(pipeline(ctx)())`. Writing
`value, err := pipeline(ctx)()` is a compile error ("2 variables but … returns 1 value").

## Core Types

```go
// Requester builds an *http.Request given a context.
type Requester = ReaderIOResult[*http.Request]  // func(context.Context) func() result.Result[*http.Request]

// Client executes a Requester and returns the response wrapped in ReaderIOResult.
type Client interface {
    Do(Requester) ReaderIOResult[*http.Response]
}
```

## Basic Usage

### 1. Create a Client

```go
import (
    HTTP "net/http"
    H    "github.com/IBM/fp-go/v2/context/readerioresult/http"
)

client := H.MakeClient(HTTP.DefaultClient)

// Or with a custom client:
custom := &HTTP.Client{Timeout: 10 * time.Second}
client := H.MakeClient(custom)
```

### 2. Build a Request

```go
// GET request (most common)
req := H.MakeGetRequest("https://api.example.com/users/1")

// Arbitrary method + body
req := H.MakeRequest("POST", "https://api.example.com/users", bodyReader)
```

### 3. Execute and Parse

```go
import (
    "context"
    H "github.com/IBM/fp-go/v2/context/readerioresult/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

client := H.MakeClient(HTTP.DefaultClient)

// ReadJSON validates status, Content-Type, then unmarshals JSON
result := H.ReadJSON[User](client)(H.MakeGetRequest("https://api.example.com/users/1"))

// Execute — provide context once. The inner () yields a Result[User] (one value).
user, err := R.Unwrap(result(context.Background())())
```

## Response Readers

All accept a `Client` and return a function `Requester → ReaderIOResult[A]`:

| Function | Returns | Notes |
|----------|---------|-------|
| `ReadJSON[A](client)` | `ReaderIOResult[A]` | Validates status + Content-Type, unmarshals JSON |
| `ReadText(client)` | `ReaderIOResult[string]` | Validates status, reads body as UTF-8 string |
| `ReadAll(client)` | `ReaderIOResult[[]byte]` | Validates status, returns raw body bytes |
| `ReadFullResponse(client)` | `ReaderIOResult[FullResponse]` | Returns `Pair[*http.Response, []byte]` |

`FullResponse = Pair[*http.Response, []byte]` — use `pair.First` / `pair.Second` to access components.

## Composing Requests in Pipelines

```go
import (
    F   "github.com/IBM/fp-go/v2/function"
    H   "github.com/IBM/fp-go/v2/context/readerioresult/http"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    IO  "github.com/IBM/fp-go/v2/io"
)

client     := H.MakeClient(HTTP.DefaultClient)
readPost   := H.ReadJSON[Post](client)

pipeline := F.Pipe2(
    H.MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1"),
    readPost,
    RIO.ChainFirstIOK(IO.Logf[Post]("Got post: %v")),
)

post, err := R.Unwrap(pipeline(context.Background())())
```

## Parallel Requests — Homogeneous Types

Use `RIO.TraverseArray` when all requests return the same type:

```go
import (
    A   "github.com/IBM/fp-go/v2/array"
    F   "github.com/IBM/fp-go/v2/function"
    H   "github.com/IBM/fp-go/v2/context/readerioresult/http"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    IO  "github.com/IBM/fp-go/v2/io"
)

type PostItem struct {
    UserID uint   `json:"userId"`
    ID     uint   `json:"id"`
    Title  string `json:"title"`
}

client     := H.MakeClient(HTTP.DefaultClient)
readPost   := H.ReadJSON[PostItem](client)

// Fetch 10 posts in parallel
data := F.Pipe3(
    A.MakeBy(10, func(i int) string {
        return fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", i+1)
    }),
    RIO.TraverseArray(F.Flow3(
        H.MakeGetRequest,
        readPost,
        RIO.ChainFirstIOK(IO.Logf[PostItem]("Post: %v")),
    )),
    RIO.ChainFirstIOK(IO.Logf[[]PostItem]("All posts: %v")),
    RIO.Map(A.Size[PostItem]),
)

count, err := R.Unwrap(data(context.Background())())
```

## Parallel Requests — Heterogeneous Types

Use `RIO.TraverseTuple2` (or `Tuple3`, etc.) when requests return different types:

```go
import (
    T   "github.com/IBM/fp-go/v2/tuple"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    H   "github.com/IBM/fp-go/v2/context/readerioresult/http"
    F   "github.com/IBM/fp-go/v2/function"
)

type CatFact struct {
    Fact string `json:"fact"`
}

client         := H.MakeClient(HTTP.DefaultClient)
readPost       := H.ReadJSON[PostItem](client)
readCatFact    := H.ReadJSON[CatFact](client)

// Execute both requests in parallel with different response types
data := F.Pipe3(
    T.MakeTuple2(
        "https://jsonplaceholder.typicode.com/posts/1",
        "https://catfact.ninja/fact",
    ),
    T.Map2(H.MakeGetRequest, H.MakeGetRequest), // build both requesters
    RIO.TraverseTuple2(readPost, readCatFact),  // run in parallel, typed
    RIO.ChainFirstIOK(IO.Logf[T.Tuple2[PostItem, CatFact]]("Result: %v")),
)

both, err := R.Unwrap(data(context.Background())())
// both.F1 is PostItem, both.F2 is CatFact
```

## Building Requests with the Builder API

For complex requests (custom headers, query params, JSON body), use the builder:

```go
import (
    B  "github.com/IBM/fp-go/v2/http/builder"
    RB "github.com/IBM/fp-go/v2/context/readerioresult/http/builder"
    F  "github.com/IBM/fp-go/v2/function"
)

// GET with query parameters
req := F.Pipe2(
    B.Default,
    B.WithURL("https://api.example.com/items?page=1"),
    B.WithQueryArg("limit")("50"),
)
requester := RB.Requester(req)

// POST with JSON body
req := F.Pipe3(
    B.Default,
    B.WithURL("https://api.example.com/users"),
    B.WithMethod("POST"),
    B.WithJSON(map[string]string{"name": "Alice"}),
    // sets Content-Type: application/json automatically
)
requester := RB.Requester(req)

// With authentication and custom headers
req := F.Pipe3(
    B.Default,
    B.WithURL("https://api.example.com/protected"),
    B.WithBearer("my-token"),           // sets Authorization: Bearer my-token
    B.WithHeader("X-Request-ID")("123"),
)
requester := RB.Requester(req)

// Execute
result := H.ReadJSON[Response](client)(requester)
data, err := R.Unwrap(result(ctx)())
```

### Builder Functions

| Function | Effect |
|----------|--------|
| `B.WithURL(url)` | Set the target URL |
| `B.WithMethod(method)` | Set HTTP method (GET, POST, PUT, DELETE, …) |
| `B.WithJSON(v)` | Marshal `v` as JSON body, set `Content-Type: application/json` |
| `B.WithBytes(data)` | Set raw bytes body, set `Content-Length` automatically |
| `B.WithHeader(key)(value)` | Add a request header |
| `B.WithBearer(token)` | Set `Authorization: Bearer <token>` |
| `B.WithQueryArg(key)(value)` | Append a query parameter |
| `B.WithGet` / `B.WithPost` / `B.WithPut` / `B.WithDelete` | Pre-bound `WithMethod` for the common verbs |
| `B.WithContentType(ct)` | Set the `Content-Type` header |
| `B.WithFormData(values)` | `url.Values` body + `Content-Type: application/x-www-form-urlencoded` |
| `B.WithoutBody` | Remove any body previously set |
| `B.WithoutHeader(key)` / `B.WithoutQueryArg(key)` | Remove a header / query parameter |

Every `With*` is an `Endomorphism[*Builder]`, so they chain freely inside `F.PipeN(B.Default, …)`.

## Error Handling

Errors from request creation, HTTP status codes, Content-Type validation, and JSON parsing all propagate automatically through the `Result` monad. You only handle errors at the call site:

```go
// Pattern 1: run it, then unwrap. pipeline(ctx)() yields a single Result[A];
// result.Unwrap turns that into idiomatic (A, error).
value, err := R.Unwrap(pipeline(ctx)())
if err != nil { /* handle */ }

// Pattern 2: run the pipeline, then eliminate the Result with result.Fold
F.Pipe1(
    pipeline(ctx)(),          // Result[MyType]
    R.Fold(
        func(err error) F.Void {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return F.VOID
        },
        func(data MyType) F.Void {
            json.NewEncoder(w).Encode(data)
            return F.VOID
        },
    ),
)
```

**Do not use `RIO.Fold` for this.** In `context/readerioresult`, `Fold[A, B](onLeft Kleisli[error, B], onRight Kleisli[A, B]) Operator[A, B]` stays *inside* the monad — both branches must return a `ReaderIOResult[B]`, so it cannot take plain `func(error)` / `func(A)` side-effecting handlers. To leave the monad, execute the pipeline and fold the resulting `Result` (as above), or use `readerioresult.Fold` from the non-context package, which lands in `ReaderIO`.

## Full HTTP Handler Example

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    F   "github.com/IBM/fp-go/v2/function"
    H   "github.com/IBM/fp-go/v2/context/readerioresult/http"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    R   "github.com/IBM/fp-go/v2/result"
    IO  "github.com/IBM/fp-go/v2/io"
)

type Post struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
}

var client = H.MakeClient(http.DefaultClient)

func fetchPost(id int) RIO.ReaderIOResult[Post] {
    url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id)
    return F.Pipe2(
        H.MakeGetRequest(url),
        H.ReadJSON[Post](client),
        RIO.ChainFirstIOK(IO.Logf[Post]("fetched: %v")),
    )
}

func handler(w http.ResponseWriter, r *http.Request) {
    F.Pipe1(
        fetchPost(1)(r.Context())(), // run it: Result[Post]
        R.Fold(
            func(err error) F.Void {
                http.Error(w, err.Error(), http.StatusBadGateway)
                return F.VOID
            },
            func(post Post) F.Void {
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(post)
                return F.VOID
            },
        ),
    )
}
```

## Import Reference

```go
import (
    HTTP "net/http"

    H   "github.com/IBM/fp-go/v2/context/readerioresult/http"
    RB  "github.com/IBM/fp-go/v2/context/readerioresult/http/builder"
    B   "github.com/IBM/fp-go/v2/http/builder"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
    R   "github.com/IBM/fp-go/v2/result"   // Unwrap, Fold — leaving the monad
    A   "github.com/IBM/fp-go/v2/array"
    T   "github.com/IBM/fp-go/v2/tuple"
    IO  "github.com/IBM/fp-go/v2/io"
)
```

Requires Go 1.24+.
