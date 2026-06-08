---
sidebar_position: 2
title: Installation
hide_title: true
description: Install fp-go via go get. Requires Go 1.24+ for v2; v1 supports 1.18+.
---

<PageHeader
  eyebrow="Getting started · Section 02 / 02"
  title="Install"
  titleAccent="fp-go."
  lede="One go get and you're ready. v2 is the active line; v1 remains supported for Go 1.18–1.23."
  meta={[
    {label: '// Version', value: <>v2.2.82 <MetaPill>LATEST</MetaPill></>},
    {label: '// Go required', value: '1.24+ (v2)'},
    {label: '// Reading time', value: '2 min · 3 sections'},
  ]}
/>

<TLDR>
  <TLDRCard label="// Go required" value="1.24+" unit="v2" description="v1 still supports Go 1.18+." />
  <TLDRCard label="// Install command" prose value={<><code>go get github.com/IBM/fp-go/v2@latest</code></>} />
  <TLDRCard label="// Versioning" prose value={<>Follows <em>SemVer.</em></>} />
</TLDR>

<Section
  id="requirements"
  number="01"
  title="Requirements"
  tag="Difficulty · Beginner"
>

<Callout title="Go version.">
  fp-go v2 requires <strong>Go 1.24 or later</strong> for the latest generics features and improvements.
</Callout>

</Section>

<Section
  id="install"
  number="02"
  title="Install"
  titleAccent="fp-go."
>

### Latest version (v2.2.82)

<CodeCard file="shell" status="tested">
{`go get github.com/IBM/fp-go/v2@latest`}
</CodeCard>

### Pin a specific version

<CodeCard file="shell" status="tested">
{`go get github.com/IBM/fp-go/v2@v2.2.82`}
</CodeCard>

<p>The library follows <a href="https://semver.org/">semantic versioning</a>.</p>

### Check the installed version

<CodeCard file="shell">
{`go list -m github.com/IBM/fp-go/v2`}
</CodeCard>

</Section>

<Section
  id="verify"
  number="03"
  title="Verify"
  titleAccent="the install."
>

<CodeCard file="main.go" status="tested">
{`package main

import (
    "fmt"
    "github.com/IBM/fp-go/v2/option"
)

func main() {
    some := option.Some(42)
    fmt.Println(option.IsSome(some)) // true
}`}
</CodeCard>

<Callout type="success" title="You're set.">
  Next, head to <a href="./quickstart">the quickstart</a> to build your first fp-go program, or jump to <a href="./option">Option</a> to learn a core data type.
</Callout>

</Section>
