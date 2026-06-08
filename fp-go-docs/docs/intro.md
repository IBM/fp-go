---
sidebar_position: 1
title: What is fp-go?
hide_title: true
description: A comprehensive functional programming library for Go — type-safe error handling, composable effects, and monadic types.
---

<PageHeader
  eyebrow="Getting started · Section 01 / 01"
  title="What is"
  titleAccent="fp-go?"
  lede="A comprehensive functional programming library for Go — strongly influenced by fp-ts, built around small, pure, composable functions."
  meta={[
    {label: '// Version', value: <>v2.2.82 <MetaPill>LATEST</MetaPill></>},
    {label: '// Go required', value: '1.24+'},
    {label: '// Reading time', value: '2 min · 1 section'},
  ]}
/>

<TLDR>
  <TLDRCard
    label="// Approach"
    prose
    value={<>Many small, <em>pure functions</em> with no hidden side effects.</>}
  />
  <TLDRCard
    label="// Side effects"
    prose
    value={<>Isolated into <em>lazy</em> IO-style computations.</>}
  />
  <TLDRCard
    label="// Composition"
    prose
    value={<>A consistent set of combinators across <em>every</em> data type.</>}
  />
</TLDR>

<Section
  id="quick-example"
  number="01"
  title="Quick"
  titleAccent="example."
  tag="Difficulty · Beginner"
  lede="Handle errors functionally with Either. Map transforms the success branch; errors flow through untouched."
>

<CodeCard file="example.go" status="tested">
{`import (
    "errors"
    "github.com/IBM/fp-go/either"
    "github.com/IBM/fp-go/function"
)

// Pure function that can fail
func divide(a, b int) either.Either[error, int] {
    if b == 0 {
        return either.Left[int](errors.New("division by zero"))
    }
    return either.Right[error](a / b)
}

// Compose operations safely
result := function.Pipe2(
    divide(10, 2),
    either.Map(func(x int) int { return x * 2 }),
    either.GetOrElse(func() int { return 0 }),
)
// result = 10`}
</CodeCard>

<Callout title="What this demonstrates.">
  <ul>
    <li>Express operations that can fail using the <code>Either</code> type</li>
    <li>Chain operations together safely</li>
    <li>Handle errors explicitly without nested <code>if</code> statements</li>
    <li>Write pure, composable functions</li>
  </ul>
</Callout>

</Section>
