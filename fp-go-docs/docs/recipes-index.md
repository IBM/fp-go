---
sidebar_position: 1
title: Recipes
hide_title: true
description: Practical guides for common functional programming patterns in fp-go.
---

<PageHeader
  eyebrow="Recipes · Overview"
  title="Practical"
  titleAccent="Recipes."
  lede="Real-world patterns and solutions for common functional programming tasks. Complete, runnable examples for error handling, data processing, I/O, and more."
  meta={[
    {label: '// Categories', value: '5 main areas'},
    {label: '// Recipes', value: '16 patterns'},
    {label: '// Difficulty', value: 'Beginner → Advanced'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Error handling" value="5" description="Validation, retry, recovery, fallback, accumulation." />
  <TLDRCard label="// Data processing" value="4" description="Transform, filter, aggregate, parse data." />
  <TLDRCard label="// I/O operations" value="3" description="Files, HTTP, parallel tasks." />
</TLDR>

<Section id="error-handling" number="01" title="Error" titleAccent="Handling.">

Learn how to handle errors functionally:

### [Validation](./recipes/validation)
Validate data with composable validators. Build complex validation logic from simple rules.

**Topics:** Input validation, form validation, business rules, validation combinators

### [Recovery](./recipes/error-recovery)
Recover from errors gracefully. Provide fallbacks and alternative paths.

**Topics:** Error recovery, fallback values, alternative computations, graceful degradation

### [Retry](./recipes/retry)
Implement retry logic for transient failures. Configure retry strategies.

**Topics:** Exponential backoff, retry policies, circuit breakers, timeout handling

</Section>

<Section id="data-processing" number="02" title="Data" titleAccent="Processing.">

Transform and process data functionally:

### [Transformation](./recipes/data-transformation)
Transform data structures with map, filter, and reduce. Build data pipelines.

**Topics:** Data mapping, filtering, reducing, pipeline composition

### [Filtering & Mapping](./recipes/filtering-mapping)
Filter collections with predicates. Combine multiple filter conditions.

**Topics:** Predicate composition, complex filters, partition, takeWhile, dropWhile

### [Aggregation](./recipes/aggregation)
Aggregate data with monoids and semigroups. Compute statistics and summaries.

**Topics:** Sum, product, min, max, average, grouping, statistical aggregation

### [Parsing](./recipes/parsing)
Parse structured data safely. Handle parsing errors functionally.

**Topics:** JSON parsing, CSV parsing, parser combinators, validation during parsing

</Section>

<Section id="io-operations" number="03" title="I/O" titleAccent="Operations.">

Handle side effects and I/O functionally:

### [File Operations](./recipes/file-operations)
Read and write files with proper error handling. Manage file resources safely.

**Topics:** File reading, file writing, resource management, bracket pattern, streaming

### [HTTP Requests](./recipes/http-requests)
Make HTTP requests with retry logic and error handling. Build resilient API clients.

**Topics:** REST APIs, retry logic, circuit breaker, rate limiting, request composition

### [Parallel Tasks](./recipes/parallel-tasks)
Execute tasks in parallel safely. Coordinate concurrent operations.

**Topics:** Concurrency, worker pools, parallel map, race conditions, synchronization

</Section>

<Section id="composition" number="04" title="Composition" titleAccent="Patterns.">

Compose complex operations from simple ones:

### [Dependency Injection](./recipes/dependency-injection)
Inject dependencies using the Reader pattern. Build testable applications.

**Topics:** Reader monad, dependency injection, testing, mocking, service architecture

### [Pipelines](./recipes/pipelines)
Build data processing pipelines. Compose operations with pipe and flow.

**Topics:** Pipe, flow, composition, data pipelines, transformation chains

### [Middleware](./recipes/middleware)
Create composable middleware for HTTP handlers. Build middleware stacks.

**Topics:** HTTP middleware, composition, logging, authentication, error handling

</Section>

<Section id="testing" number="05" title="Testing" titleAccent="Patterns.">

Test functional code effectively:

### [Testing Pure Functions](./recipes/testing-pure)
Test pure functions with property-based testing. Verify function laws.

**Topics:** Unit testing, property-based testing, test generation, function laws

### [Testing Effects](./recipes/testing-effects)
Test effectful code with IO and IOEither. Mock dependencies and side effects.

**Topics:** IO testing, mocking, integration testing, test doubles, dependency injection

</Section>
