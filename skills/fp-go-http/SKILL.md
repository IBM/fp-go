---
name: fp-go-http
description: Use this skill when making HTTP requests in fp-go using the ReaderIOResult-based HTTP client (github.com/IBM/fp-go/v2/context/readerioresult/http). Trigger on mentions of fp-go HTTP, MakeClient, MakeGetRequest, MakeRequest, ReadJSON, ReadText, ReadAll, ReadFullResponse, the HTTP request builder (WithURL, WithJSON, WithBearer, WithHeader, WithQueryArg), HTTP header name constants (http/headers: ContentType, Accept, Authorization, XRequestID, …) or content type / media type constants (http/content: JSON, ProblemJSON, FormEncoded, OctetStream, …), parallel requests with TraverseArray or TraverseTuple2, or building context-aware, composable HTTP pipelines that propagate errors through the Result monad.
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

### 4. Timeouts and Request-Scoped Values

`HTTP.Client.Timeout` is a global cap. For a per-request or per-pipeline bound, scope the context instead: requests honour context cancellation, and the cancel func is released automatically.

```go
bounded := F.Pipe2(
    H.ReadJSON[User](client)(H.MakeGetRequest("https://api.example.com/users/1")),
    RIO.WithTimeout[User](2*time.Second),
    RIO.WithValue[User](traceIDKey, traceID), // read downstream with RIO.AskValue[string](traceIDKey)
)
```

Never write `ctx, cancel := context.WithTimeout(...)` around the call or `ctx.Value(key).(V)` inside it. See the `fp-go-context` skill for the full context-handling guide.

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
    N   "github.com/IBM/fp-go/v2/number"
    S   "github.com/IBM/fp-go/v2/string"
)

type PostItem struct {
    UserID uint   `json:"userId"`
    ID     uint   `json:"id"`
    Title  string `json:"title"`
}

client     := H.MakeClient(HTTP.DefaultClient)
readPost   := H.ReadJSON[PostItem](client)

// index -> URL, point-free
postURL := F.Flow2(N.Add(1), S.Format[int]("https://jsonplaceholder.typicode.com/posts/%d"))

// Fetch 10 posts in parallel
data := F.Pipe3(
    A.MakeBy(10, postURL),
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
    HD "github.com/IBM/fp-go/v2/http/headers"
    C  "github.com/IBM/fp-go/v2/http/content"
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
    B.WithPost,
    B.WithJSON(map[string]string{"name": "Alice"}),
    // sets Content-Type: application/json (C.JSON) automatically
)
requester := RB.Requester(req)

// With authentication and custom headers — use the constants, not string literals
req := F.Pipe4(
    B.Default,
    B.WithURL("https://api.example.com/protected"),
    B.WithBearer("my-token"),               // sets Authorization: Bearer my-token
    B.WithHeader(HD.Accept)(C.JSON),
    B.WithHeader(HD.XRequestID)("123"),
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

## Header Names and Content Types

Never spell header names or media types as string literals — fp-go ships constants for both.

### Header names — `http/headers` (alias `HD`)

`github.com/IBM/fp-go/v2/http/headers` defines the commonly used header names as **lower-case** constants (the form mandated by HTTP/2 and HTTP/3):

| Group | Constants |
|-------|-----------|
| Representation | `ContentType`, `ContentLength`, `ContentEncoding`, `ContentLanguage`, `ContentDisposition`, `ContentRange`, `ContentLocation`, `TransferEncoding`, `Date`, `Link` |
| Request | `Accept`, `AcceptCharset`, `AcceptEncoding`, `AcceptLanguage`, `Authorization`, `ProxyAuthorization`, `Cookie`, `Host`, `UserAgent`, `Referer`, `Origin`, `Range`, `Expect`, `Forwarded` |
| Response | `Location`, `Server`, `SetCookie`, `WWWAuthenticate`, `ProxyAuthenticate`, `RetryAfter`, `Allow`, `AcceptRanges` |
| Caching / conditional | `CacheControl`, `ETag`, `LastModified`, `Expires`, `Age`, `Vary`, `Pragma`, `IfMatch`, `IfNoneMatch`, `IfModifiedSince`, `IfUnmodifiedSince`, `IfRange` |
| CORS | `AccessControlAllowOrigin`, `AccessControlAllowMethods`, `AccessControlAllowHeaders`, `AccessControlAllowCredentials`, `AccessControlExposeHeaders`, `AccessControlMaxAge`, `AccessControlRequestMethod`, `AccessControlRequestHeaders` |
| Security | `StrictTransportSecurity`, `ContentSecurityPolicy`, `XContentTypeOptions`, `XFrameOptions`, `ReferrerPolicy` |
| Proxy / tracing | `XForwardedFor`, `XForwardedHost`, `XForwardedProto`, `XRequestID`, `XCorrelationID`, `Traceparent`, `Tracestate` |

Lower case is safe with HTTP/1.1 too: `http.Header.Get/Set/Add/Values/Del` and the builder canonicalize keys. **Only raw map indexing does not** — `h[HD.ContentType]` misses `"Content-Type"`. Use `h.Get(HD.ContentType)` or the lenses below.

The package also provides functional access to `http.Header`:

```go
HD.AtValue(HD.Authorization).Get(h)                      // Option[string] — first value
HD.AtValues(HD.Accept).Get(h)                            // Option[[]string] — all values, None if absent
h2 := HD.AtValue(HD.ContentType).Set(O.Some(C.JSON))(h)  // new header map; O.None removes the header
merged := HD.Monoid.Concat(defaults, overrides)          // union; values of shared keys are concatenated
```

Both lenses canonicalize the header name, and `Set` returns a new `http.Header`, leaving `h` untouched.

### Content types — `http/content` (alias `C`)

`github.com/IBM/fp-go/v2/http/content` defines media type constants (bare type, no parameters):

| Group | Constants |
|-------|-----------|
| JSON family | `JSON`, `ProblemJSON`, `JSONPatch`, `MergePatch`, `NDJSON`, `JSONLD`, `HALJSON`, `JSONAPI` |
| XML family | `XML`, `TextXML`, `ProblemXML`, `SOAP` |
| Other structured | `YAML`, `CBOR`, `Protobuf`, `GRPC` |
| Text | `TextPlain`, `TextHTML`, `TextCSS`, `TextCSV`, `TextJavaScript`, `TextMarkdown`, `TextEventStream` |
| Forms / multipart | `FormEncoded`, `MultipartFormData`, `MultipartMixed`, `MultipartByteRanges` |
| Binary | `OctetStream`, `PDF`, `ZIP`, `Gzip`, `WASM` |
| Media | `ImagePNG`, `ImageJPEG`, `ImageGIF`, `ImageWebP`, `ImageAVIF`, `ImageSVG`, `AudioMPEG`, `AudioOGG`, `VideoMP4`, `VideoWebM`, `FontWOFF2` |

`C.Json` is deprecated — use `C.JSON`.

The builder already uses these: `B.WithJSON` sets `C.JSON`, `B.WithFormData` sets `C.FormEncoded`; `B.WithContentType(ct)` and `B.WithAuthorization(v)` are `B.WithHeader(HD.ContentType)` / `B.WithHeader(HD.Authorization)`.

A received `Content-Type` usually carries parameters (`application/json; charset=utf-8`), so parse before comparing against a constant:

```go
import (
    FH "github.com/IBM/fp-go/v2/http"
    PA "github.com/IBM/fp-go/v2/pair"
    R  "github.com/IBM/fp-go/v2/result"
    S  "github.com/IBM/fp-go/v2/string"
)

isJSON := F.Flow3(
    FH.ParseMediaType,                       // string -> Result[Pair[mediaType, params]]
    R.Map(PA.Head[string, map[string]string]),
    R.Fold(F.Constant1[error](false), S.Equals(C.JSON)),
)
isJSON(resp.Header.Get(HD.ContentType))
```

`H.ReadJSON` already performs this validation (JSON or any `+json` type) — no manual check is needed on that path.

## Error Handling

Errors from request creation, HTTP status codes, Content-Type validation, and JSON parsing all propagate automatically through the `Result` monad. You only handle errors at the call site:

```go
// Pattern 1: run it, then unwrap. pipeline(ctx)() yields a single Result[A];
// result.Unwrap turns that into idiomatic (A, error).
value, err := R.Unwrap(pipeline(ctx)())
if err != nil { /* handle */ }

// Pattern 2: run the pipeline, then eliminate the Result with result.Fold,
// passing named handlers (defined once, see the full example below)
F.Pipe1(
    pipeline(ctx)(),          // Result[MyType]
    R.Fold(writeError(w, http.StatusInternalServerError), writeJSON[MyType](w)),
)
```

**Do not use `RIO.Fold` for this.** In `context/readerioresult`, `Fold[A, B](onLeft Kleisli[error, B], onRight Kleisli[A, B]) Operator[A, B]` stays *inside* the monad — both branches must return a `ReaderIOResult[B]`, so it cannot take plain `func(error)` / `func(A)` side-effecting handlers. To leave the monad, execute the pipeline and fold the resulting `Result` (as above), or use `readerioresult.Fold` from the non-context package, which lands in `ReaderIO`.

## Full HTTP Handler Example

```go
package main

import (
    "encoding/json"
    "net/http"

    F   "github.com/IBM/fp-go/v2/function"
    H   "github.com/IBM/fp-go/v2/context/readerioresult/http"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    R   "github.com/IBM/fp-go/v2/result"
    IO  "github.com/IBM/fp-go/v2/io"
    S   "github.com/IBM/fp-go/v2/string"
    HD  "github.com/IBM/fp-go/v2/http/headers"
    C   "github.com/IBM/fp-go/v2/http/content"
)

type Post struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
}

var client = H.MakeClient(http.DefaultClient)

// fetchPost: int -> ReaderIOResult[Post], point-free
func fetchPost() RIO.Kleisli[int, Post] {
    return F.Flow4(
        S.Format[int]("https://jsonplaceholder.typicode.com/posts/%d"),
        H.MakeGetRequest,
        H.ReadJSON[Post](client),
        RIO.ChainFirstIOK(IO.Logf[Post]("fetched: %v")),
    )
}

// Leaf handlers — the only place that touches the ResponseWriter
func writeError(w http.ResponseWriter, status int) func(error) F.Void {
    return func(err error) F.Void {
        http.Error(w, err.Error(), status)
        return F.VOID
    }
}

func writeJSON[A any](w http.ResponseWriter) func(A) F.Void {
    return func(a A) F.Void {
        w.Header().Set(HD.ContentType, C.JSON)
        json.NewEncoder(w).Encode(a)
        return F.VOID
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    F.Pipe1(
        fetchPost()(1)(r.Context())(), // run it: Result[Post]
        R.Fold(writeError(w, http.StatusBadGateway), writeJSON[Post](w)),
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
    HD  "github.com/IBM/fp-go/v2/http/headers"   // header name constants, AtValue/AtValues lenses, Monoid
    C   "github.com/IBM/fp-go/v2/http/content"   // content type constants
    FH  "github.com/IBM/fp-go/v2/http"           // ParseMediaType, StatusCodeError, …
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
    R   "github.com/IBM/fp-go/v2/result"   // Unwrap, Fold — leaving the monad
    A   "github.com/IBM/fp-go/v2/array"
    T   "github.com/IBM/fp-go/v2/tuple"
    IO  "github.com/IBM/fp-go/v2/io"
    N   "github.com/IBM/fp-go/v2/number"
    S   "github.com/IBM/fp-go/v2/string"
    PA  "github.com/IBM/fp-go/v2/pair"
)
```

Requires Go 1.24+.
