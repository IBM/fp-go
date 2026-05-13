---
sidebar_position: 5
title: Error Recovery
description: Graceful error recovery strategies
hide_title: true
---

<PageHeader
  eyebrow="Recipes · 05 / 17"
  title="Error"
  titleAccent="Recovery"
  lede="Graceful error recovery strategies with fallbacks, cascading sources, and resilient degradation patterns."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '7' },
    { label: 'Use Cases', value: 'Resilience, Fallbacks, Degradation' }
  ]}
/>

<TLDR>
  <TLDRCard title="Always Have a Fallback" icon="shield">
    Never leave users with no response—provide default values, cached data, or reduced functionality when operations fail.
  </TLDRCard>
  <TLDRCard title="Degrade Gracefully" icon="layers">
    Provide reduced functionality rather than complete failure—try multiple sources before giving up.
  </TLDRCard>
  <TLDRCard title="Log and Monitor" icon="activity">
    Track errors for debugging while maintaining user experience—high fallback usage indicates underlying issues.
  </TLDRCard>
</TLDR>

<Section id="fallback-values" number="01" title="Fallback" titleAccent="Values">

Provide default values when operations fail, ensuring your application always has a valid response.

<CodeCard file="fallback-values.go">
{`package main

import (
    "fmt"
    
    O "github.com/IBM/fp-go/v2/option"
    R "github.com/IBM/fp-go/v2/result"
)

// Get user preference with fallback
func getUserPreference(userID string) R.Result[string] {
    // Simulate database lookup failure
    return R.Left[string](fmt.Errorf("user not found"))
}

func getDefaultPreference() string {
    return "default-theme"
}

func main() {
    // Try to get user preference, fall back to default
    result := getUserPreference("user123")
    preference := R.GetOrElse(getDefaultPreference)(result)
    
    fmt.Printf("Using preference: %s\\n", preference)
    // Output: Using preference: default-theme
}`}
</CodeCard>

</Section>

<Section id="cascading-fallbacks" number="02" title="Cascading" titleAccent="Fallbacks">

Try multiple sources before giving up, creating a resilient chain of fallback options.

<CodeCard file="cascading-fallbacks.go">
{`package main

import (
    "fmt"
    
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

// Try to get config from multiple sources
func getConfigFromFile() R.Result[string] {
    return R.Left[string](fmt.Errorf("file not found"))
}

func getConfigFromEnv() R.Result[string] {
    return R.Left[string](fmt.Errorf("env var not set"))
}

func getConfigFromRemote() R.Result[string] {
    return R.Right[error]("remote-config")
}

func getConfig() R.Result[string] {
    return F.Pipe3(
        getConfigFromFile(),
        R.OrElse(getConfigFromEnv),
        R.OrElse(getConfigFromRemote),
    )
}

func main() {
    config := getConfig()
    value := R.GetOrElse(func() string { return "fallback" })(config)
    
    fmt.Printf("Config: %s\\n", value)
    // Output: Config: remote-config
}`}
</CodeCard>

</Section>

<Section id="partial-success" number="03" title="Partial" titleAccent="Success">

Handle partial failures in batch operations by collecting both successes and failures.

<CodeCard file="partial-success.go">
{`package main

import (
    "fmt"
    
    A "github.com/IBM/fp-go/v2/array"
    E "github.com/IBM/fp-go/v2/either"
)

type ProcessResult struct {
    Successes []string
    Failures  []error
}

// Process items, collecting both successes and failures
func processItems(items []string) ProcessResult {
    var successes []string
    var failures []error
    
    for _, item := range items {
        result := processItem(item)
        if E.IsRight(result) {
            successes = append(successes, E.GetRight(result))
        } else {
            failures = append(failures, E.GetLeft(result))
        }
    }
    
    return ProcessResult{
        Successes: successes,
        Failures:  failures,
    }
}

func processItem(item string) E.Either[error, string] {
    if len(item) < 3 {
        return E.Left[string](fmt.Errorf("item too short: %s", item))
    }
    return E.Right[error](fmt.Sprintf("processed-%s", item))
}

func main() {
    items := []string{"apple", "ab", "banana", "x", "cherry"}
    result := processItems(items)
    
    fmt.Printf("Successes: %d\\n", len(result.Successes))
    fmt.Printf("Failures: %d\\n", len(result.Failures))
    
    for _, s := range result.Successes {
        fmt.Printf("  ✓ %s\\n", s)
    }
    
    for _, f := range result.Failures {
        fmt.Printf("  ✗ %s\\n", f.Error())
    }
}`}
</CodeCard>

</Section>

<Section id="graceful-degradation" number="04" title="Graceful" titleAccent="Degradation">

Provide reduced functionality when full functionality fails, maintaining user experience.

<CodeCard file="graceful-degradation.go">
{`package main

import (
    "fmt"
    
    R "github.com/IBM/fp-go/v2/result"
)

type UserProfile struct {
    Name   string
    Avatar string
    Bio    string
}

type BasicProfile struct {
    Name string
}

// Try to get full profile
func getFullProfile(userID string) R.Result[UserProfile] {
    return R.Left[UserProfile](fmt.Errorf("profile service unavailable"))
}

// Fallback to basic profile
func getBasicProfile(userID string) R.Result[BasicProfile] {
    return R.Right[error](BasicProfile{Name: "User " + userID})
}

// Convert basic to full profile with defaults
func basicToFull(basic BasicProfile) UserProfile {
    return UserProfile{
        Name:   basic.Name,
        Avatar: "default-avatar.png",
        Bio:    "No bio available",
    }
}

func getUserProfile(userID string) UserProfile {
    fullProfile := getFullProfile(userID)
    
    if R.IsRight(fullProfile) {
        return R.GetRight(fullProfile)
    }
    
    // Degrade to basic profile
    basicProfile := getBasicProfile(userID)
    if R.IsRight(basicProfile) {
        return basicToFull(R.GetRight(basicProfile))
    }
    
    // Ultimate fallback
    return UserProfile{
        Name:   "Anonymous",
        Avatar: "default-avatar.png",
        Bio:    "Profile unavailable",
    }
}

func main() {
    profile := getUserProfile("123")
    fmt.Printf("Name: %s\\n", profile.Name)
    fmt.Printf("Avatar: %s\\n", profile.Avatar)
    fmt.Printf("Bio: %s\\n", profile.Bio)
}`}
</CodeCard>

</Section>

<Section id="error-logging" number="05" title="Error Logging with" titleAccent="Recovery">

Log errors while providing fallback values, maintaining observability without sacrificing user experience.

<CodeCard file="error-logging.go">
{`package main

import (
    "fmt"
    "log"
    
    R "github.com/IBM/fp-go/v2/result"
)

// Tap into error for logging without changing the flow
func tapError[A any](onError func(error)) func(R.Result[A]) R.Result[A] {
    return func(result R.Result[A]) R.Result[A] {
        if R.IsLeft(result) {
            onError(R.GetLeft(result))
        }
        return result
    }
}

func fetchData(id string) R.Result[string] {
    return R.Left[string](fmt.Errorf("network error: timeout"))
}

func main() {
    result := fetchData("123")
    
    // Log error and provide fallback
    logged := tapError[string](func(err error) {
        log.Printf("Error fetching data: %v", err)
    })(result)
    
    value := R.GetOrElse(func() string { return "cached-data" })(logged)
    
    fmt.Printf("Using: %s\\n", value)
}`}
</CodeCard>

</Section>

<Section id="timeout-fallback" number="06" title="Timeout with" titleAccent="Fallback">

Implement timeout with graceful fallback to prevent indefinite waiting.

<CodeCard file="timeout-fallback.go">
{`package main

import (
    "context"
    "fmt"
    "time"
    
    IO "github.com/IBM/fp-go/v2/io"
    IOR "github.com/IBM/fp-go/v2/ioresult"
)

// Execute with timeout
func withTimeout[A any](
    timeout time.Duration,
    operation func() IOR.IOResult[A],
    fallback func() A,
) IO.IO[A] {
    return IO.MakeIO(func() A {
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()
        
        resultChan := make(chan A, 1)
        
        go func() {
            result := operation()
            outcome := result()
            if outcome.IsRight() {
                resultChan <- outcome.GetRight()
            }
        }()
        
        select {
        case result := <-resultChan:
            return result
        case <-ctx.Done():
            fmt.Println("Operation timed out, using fallback")
            return fallback()
        }
    })
}

func slowOperation() IOR.IOResult[string] {
    return IOR.FromIO[error](IO.MakeIO(func() string {
        time.Sleep(2 * time.Second)
        return "slow-result"
    }))
}

func main() {
    operation := withTimeout(
        500*time.Millisecond,
        slowOperation,
        func() string { return "fallback-result" },
    )
    
    result := operation()
    fmt.Printf("Result: %s\\n", result)
}`}
</CodeCard>

</Section>

<Section id="retry-fallback" number="07" title="Retry with" titleAccent="Fallback">

Combine retry logic with fallback values for maximum resilience.

<CodeCard file="retry-fallback.go">
{`package main

import (
    "fmt"
    "time"
    
    IO "github.com/IBM/fp-go/v2/io"
    IOR "github.com/IBM/fp-go/v2/ioresult"
)

func retryWithFallback[A any](
    maxAttempts int,
    operation func() IOR.IOResult[A],
    fallback func() A,
) IO.IO[A] {
    return IO.MakeIO(func() A {
        for i := 0; i < maxAttempts; i++ {
            result := operation()
            outcome := result()
            
            if outcome.IsRight() {
                return outcome.GetRight()
            }
            
            if i < maxAttempts-1 {
                time.Sleep(100 * time.Millisecond)
            }
        }
        
        fmt.Printf("All %d attempts failed, using fallback\\n", maxAttempts)
        return fallback()
    })
}

var attemptNum = 0

func unreliableOp() IOR.IOResult[string] {
    return IOR.FromIO[error](IO.MakeIO(func() string {
        attemptNum++
        fmt.Printf("Attempt %d failed\\n", attemptNum)
        return ""
    }))
}

func main() {
    operation := retryWithFallback(
        3,
        unreliableOp,
        func() string { return "fallback-value" },
    )
    
    result := operation()
    fmt.Printf("Final result: %s\\n", result)
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="08" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Always have a fallback** — Never leave users with no response
  </ChecklistItem>
  <ChecklistItem status="required">
    **Log failures** — Track errors for debugging and monitoring
  </ChecklistItem>
  <ChecklistItem status="required">
    **Degrade gracefully** — Provide reduced functionality rather than complete failure
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Use timeouts** — Don't wait forever for failing operations
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Cache when possible** — Use cached data as fallback for fresh data
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Communicate degradation** — Let users know when using fallback/cached data
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Monitor fallback usage** — High fallback usage indicates underlying issues
  </ChecklistItem>
</Checklist>

</Section>
