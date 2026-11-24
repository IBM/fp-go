---
sidebar_position: 2
---

# Installation

## Requirements

fp-go v1 requires **Go 1.18 or later** for generics support.

## Install

To install fp-go v1, use the following command:

```bash
go get github.com/IBM/fp-go
```

This will download and install version 1 of fp-go into your Go module.

## Verify Installation

After installation, you can verify it's working by importing it in your Go code:

```go
package main

import (
    "fmt"
    "github.com/IBM/fp-go/option"
)

func main() {
    some := option.Some(42)
    fmt.Println(option.IsSome(some)) // true
}
```

## Upgrading to v2

If you want to use the latest features and improvements, consider upgrading to [v2](../intro.md) which requires Go 1.24+.

## Next Steps

Now that you have fp-go installed, check out the [Option](./option.md) documentation to learn about one of the core data types.