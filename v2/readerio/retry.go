// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package readerio

import (
	"github.com/IBM/fp-go/v2/retry"
	RG "github.com/IBM/fp-go/v2/retry/generic"
)

// Retrying retries a ReaderIO action according to a retry policy until it succeeds or the policy gives up.
//
// This function implements retry logic for ReaderIO computations that depend on an environment
// and perform side effects. It combines the Reader monad (for dependency injection) with the IO
// monad (for side effects) and adds retry capabilities. The retry behavior is controlled by three
// parameters: a policy, an action, and a check predicate.
//
// # Type Parameters
//
//   - R: The environment type (Reader context) - e.g., Config, Logger, Database connection
//   - A: The result type produced by the action
//
// # Parameters
//
//   - policy: A RetryPolicy that determines the delay between retries and when to stop.
//     Policies can be combined using retry.Monoid to create complex retry strategies.
//     Common policies include:
//
//   - retry.LimitRetries(n): Limit the number of retry attempts
//
//   - retry.ExponentialBackoff(delay): Exponentially increase delay between retries
//
//   - retry.ConstantDelay(delay): Use a constant delay between retries
//
//   - retry.CapDelay(max, policy): Cap the maximum delay
//
//   - action: A Kleisli arrow (function) that takes the current RetryStatus and returns a
//     ReaderIO[R, A]. The action is executed on each attempt, receiving updated status
//     information including the iteration number, cumulative delay, and previous delay.
//     The action has access to the environment R on each retry attempt.
//
//   - check: A predicate function that examines the result of the action and returns true if
//     the operation should be retried, or false if it succeeded. This allows you to define
//     custom success criteria based on the result value. Return true to retry, false to stop.
//
// # Returns
//
// A ReaderIO[R, A] that:
//   - Requires an environment of type R
//   - Performs the retry logic with side effects
//   - Produces a result of type A
//
// # Retry Flow
//
// The function will:
//  1. Execute the action with the current retry status and environment
//  2. Apply the check predicate to the result
//  3. If check returns false (success), return the result
//  4. If check returns true (should retry), apply the policy to get the next delay
//  5. If the policy returns None, stop retrying and return the last result
//  6. If the policy returns Some(delay), wait for that duration and retry from step 1
//
// The action receives RetryStatus information on each attempt, which includes:
//   - IterNumber: The current attempt number (0-indexed, so 0 is the first attempt)
//   - CumulativeDelay: The total time spent waiting between retries so far
//   - PreviousDelay: The delay from the last retry (None on the first attempt)
//
// # Comparison with Other Monads
//
// Unlike io.Retrying:
//   - Adds environment dependency (Reader aspect) available throughout retry attempts
//   - The environment R is passed to every retry attempt
//   - Useful when retry logic needs configuration, logging, or other context
//
// Unlike readerioeither.Retrying:
//   - No error channel (no Left/Right distinction)
//   - Simpler for cases where you don't need explicit error handling
//   - Check predicate operates on the success value directly
//
// # Use Cases
//
//  1. Retrying operations with environment-dependent configuration:
//     - HTTP requests with retry configuration from environment
//     - Database operations with connection pool from environment
//     - File operations with paths from configuration
//
//  2. Retrying with logging:
//     - Log each retry attempt using a logger from the environment
//     - Track retry metrics using monitoring from environment
//     - Debug retry behavior with environment-provided debug flags
//
//  3. Retrying with dynamic behavior:
//     - Adjust retry behavior based on environment state
//     - Use different strategies based on environment configuration
//     - Access shared resources during retry attempts
//
// # Example: HTTP Request with Retry and Logging
//
//	type Env struct {
//	    HTTPClient *http.Client
//	    Logger     func(string)
//	    MaxRetries int
//	}
//
//	policy := retry.Monoid.Concat(
//	    retry.LimitRetries(5),
//	    retry.ExponentialBackoff(100 * time.Millisecond),
//	)
//
//	fetchData := func(status retry.RetryStatus) readerio.ReaderIO[Env, *http.Response] {
//	    return func(env Env) io.IO[*http.Response] {
//	        return func() *http.Response {
//	            env.Logger(fmt.Sprintf("Attempt %d (cumulative delay: %v)",
//	                status.IterNumber, status.CumulativeDelay))
//	            resp, _ := env.HTTPClient.Get("https://api.example.com/data")
//	            return resp
//	        }
//	    }
//	}
//
//	checkResponse := func(resp *http.Response) bool {
//	    // Retry on server errors (5xx status codes)
//	    return resp.StatusCode >= 500
//	}
//
//	result := readerio.Retrying(policy, fetchData, checkResponse)
//	env := Env{HTTPClient: http.DefaultClient, Logger: log.Println, MaxRetries: 5}
//	response := result(env)() // Execute with environment
//
// # Example: Database Query with Retry
//
//	type Config struct {
//	    DB     *sql.DB
//	    Logger *log.Logger
//	}
//
//	policy := retry.Monoid.Concat(
//	    retry.LimitRetries(3),
//	    retry.ConstantDelay(500 * time.Millisecond),
//	)
//
//	queryUser := func(status retry.RetryStatus) readerio.ReaderIO[Config, *User] {
//	    return func(cfg Config) io.IO[*User] {
//	        return func() *User {
//	            cfg.Logger.Printf("Query attempt %d", status.IterNumber)
//	            var user User
//	            err := cfg.DB.QueryRow("SELECT * FROM users WHERE id = ?", 123).Scan(&user)
//	            if err != nil {
//	                return nil // Will retry on nil
//	            }
//	            return &user
//	        }
//	    }
//	}
//
//	checkUser := func(user *User) bool {
//	    return user == nil // Retry if user is nil
//	}
//
//	result := readerio.Retrying(policy, queryUser, checkUser)
//	config := Config{DB: db, Logger: logger}
//	user := result(config)() // Execute with config
//
// # Example: Retry Until Condition Met
//
//	type Env struct {
//	    StatusChecker func() string
//	    Logger        func(string)
//	}
//
//	policy := retry.Monoid.Concat(
//	    retry.LimitRetries(10),
//	    retry.ConstantDelay(1 * time.Second),
//	)
//
//	checkStatus := func(status retry.RetryStatus) readerio.ReaderIO[Env, string] {
//	    return func(env Env) io.IO[string] {
//	        return func() string {
//	            currentStatus := env.StatusChecker()
//	            env.Logger(fmt.Sprintf("Attempt %d: status is %s", status.IterNumber, currentStatus))
//	            return currentStatus
//	        }
//	    }
//	}
//
//	shouldRetry := func(status string) bool {
//	    return status != "ready" // Retry until status is "ready"
//	}
//
//	result := readerio.Retrying(policy, checkStatus, shouldRetry)
//	env := Env{StatusChecker: getStatus, Logger: log.Println}
//	finalStatus := result(env)() // Returns "ready" or last status
//
// # Stack Safety
//
// This function uses tail recursion via TailRec to ensure stack safety even with
// many retry attempts. Deep recursion (thousands of retries) will not cause stack overflow.
//
// # Performance Considerations
//
//   - Each retry creates a new IO action by calling the action with updated status
//   - The environment R is passed to every retry attempt
//   - Delays are implemented using time.Sleep, blocking the goroutine
//   - For high-frequency retries, consider the overhead of environment access
//
// # See Also
//
//   - [io.Retrying]: Retry without environment dependency
//   - [readerioeither.Retrying]: Retry with environment and error handling
//   - [retry.RetryPolicy]: For creating and combining retry policies
//   - [retry.RetryStatus]: For understanding retry status information
//   - [TailRec]: For understanding the tail recursion mechanism
//
//go:inline
func Retrying[R, A any](
	policy retry.RetryPolicy,
	action Kleisli[R, retry.RetryStatus, A],
	check Predicate[A],
) ReaderIO[R, A] {
	// Delegate to the generic retry implementation with trampoline-based tail recursion.
	// This provides stack-safe retry logic by using an iterative approach internally.
	return RG.Retrying(
		Chain[R, A, Trampoline[retry.RetryStatus, A]],
		Map[R, retry.RetryStatus, Trampoline[retry.RetryStatus, A]],
		Of[R, Trampoline[retry.RetryStatus, A]],
		Of[R, retry.RetryStatus],
		Delay[R, retry.RetryStatus],

		TailRec[R, retry.RetryStatus, A],

		policy,
		action,
		check,
	)
}
