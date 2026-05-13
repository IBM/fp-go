---
sidebar_position: 1
title: Recipes Overview
hide_title: true
description: Practical examples and patterns for using fp-go
---

<PageHeader
  eyebrow="Recipes · Introduction"
  title="Recipes"
  titleAccent="Overview."
  lede="Practical recipes and patterns for common tasks using fp-go. Each recipe demonstrates real-world usage with complete, runnable examples."
  meta={[
    {label: '// Format', value: 'Problem → Solution → Explanation'},
    {label: '// Code', value: 'Complete & runnable'},
    {label: '// Difficulty', value: 'All levels'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Each recipe includes" prose value={<>Problem, solution, <em>explanation</em>, variations.</>} />
  <TLDRCard label="// Prerequisites" prose value={<>Go 1.24+, fp-go v2 <em>installed</em>.</>} />
  <TLDRCard label="// Categories" value="5" description="Error handling, data, I/O, composition, testing." variant="up" />
</TLDR>

<Section id="error-handling" number="01" title="Error" titleAccent="Handling.">

Learn how to handle errors functionally:

- **[Error Handling Patterns](./error-handling)** - Common error handling strategies
- **[Validation](./validation)** - Input validation with accumulating errors
- **[Retry Logic](./retry)** - Implementing retry with exponential backoff
- **[Error Recovery](./error-recovery)** - Graceful error recovery strategies

</Section>

<Section id="data-processing" number="02" title="Data" titleAccent="Processing.">

Transform and process data functionally:

- **[Data Transformation](./data-transformation)** - Pipeline-based data processing
- **[Filtering and Mapping](./filtering-mapping)** - Working with collections
- **[Aggregation](./aggregation)** - Reducing and aggregating data
- **[Parsing](./parsing)** - Safe parsing with error handling

</Section>

<Section id="io-operations" number="03" title="I/O" titleAccent="Operations.">

Handle asynchronous operations:

- **[HTTP Requests](./http-requests)** - Making HTTP calls with IOResult
- **[File Operations](./file-operations)** - Reading and writing files safely
- **[Parallel Tasks](./parallel-tasks)** - Running operations in parallel

</Section>

<Section id="composition" number="04" title="Composition" titleAccent="Patterns.">

Compose functions and effects:

- **[Pipeline Building](./pipelines)** - Building data processing pipelines
- **[Dependency Injection](./dependency-injection)** - Using Reader for DI
- **[Middleware Patterns](./middleware)** - Building middleware chains

</Section>

<Section id="testing" number="05" title="Testing" titleAccent="Patterns.">

Test functional code effectively:

- **[Testing Pure Functions](./testing-pure)** - Testing without side effects
- **[Testing Effects](./testing-effects)** - Testing IO and IOResult

</Section>

<Section id="getting-started" number="06" title="Getting" titleAccent="Started.">

### Recipe Structure

Each recipe includes:

- **Problem Statement** - What problem does this solve?
- **Solution** - Complete, runnable code example
- **Explanation** - Step-by-step breakdown
- **Variations** - Alternative approaches
- **Best Practices** - Tips and recommendations

### Prerequisites

All recipes assume you have:

<CodeCard file="install.sh" lang="bash">
{`go get github.com/IBM/fp-go/v2`}
</CodeCard>

And basic familiarity with fp-go core types like `Option`, `Result`, `Either`, and `IO`.

### Contributing

Have a recipe to share? Contributions are welcome! See our [contribution guidelines](https://github.com/IBM/fp-go/blob/main/CONTRIBUTING.md).

</Section>
