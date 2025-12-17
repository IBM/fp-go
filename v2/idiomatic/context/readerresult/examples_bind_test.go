// Copyright (c) 2023 - 2025 IBM Corp.
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

package readerresult

import (
	"context"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/lens"
	RES "github.com/IBM/fp-go/v2/result"
)

// Post represents a blog post
// fp-go:Lens
type Post struct {
	ID     int
	UserID int
	Title  string
}

// State represents accumulated state in do-notation
// fp-go:Lens
type State struct {
	User     User
	Posts    []Post
	FullName string
	Status   string
}

// getUser simulates fetching a user by ID
func getUser(id int) ReaderResult[User] {
	return func(ctx context.Context) (User, error) {
		return User{ID: id, Name: "Alice"}, nil
	}
}

// getPosts simulates fetching posts for a user
func getPosts(userID int) ReaderResult[[]Post] {
	return func(ctx context.Context) ([]Post, error) {
		return []Post{
			{ID: 1, UserID: userID, Title: "First Post"},
			{ID: 2, UserID: userID, Title: "Second Post"},
		}, nil
	}
}

// fp-go:Lens
type SimpleState struct {
	Value int
}

// ExampleDo demonstrates initializing a do-notation context with an empty state.
// This is the starting point for do-notation style composition, which allows
// imperative-style sequencing of ReaderResult computations while maintaining
// functional purity.
//
// Step-by-step breakdown:
//
//  1. Do(SimpleState{}) - Initialize the do-notation chain with an empty SimpleState.
//     This creates a ReaderResult that, when executed, will return the initial state.
//     The state acts as an accumulator that will be threaded through subsequent operations.
//
//  2. LetToL(simpleStateLenses.Value, 42) - Set the Value field to the constant 42.
//     LetToL uses a lens to focus on a specific field in the state and assign a constant value.
//     The "L" suffix indicates this is the lens-based version of LetTo.
//     After this step, state.Value = 42.
//
//  3. LetL(simpleStateLenses.Value, N.Add(8)) - Transform the Value field by adding 8.
//     LetL uses a lens to focus on a field and apply a transformation function to it.
//     N.Add(8) creates a function that adds 8 to its input.
//     After this step, state.Value = 42 + 8 = 50.
//
//  4. result(context.Background()) - Execute the composed ReaderResult computation.
//     This runs the entire chain with the provided context and returns the final state
//     and any error that occurred during execution.
//
// The key insight: Do-notation allows you to build complex stateful computations
// in a declarative, pipeline style while maintaining immutability and composability.
func ExampleDo() {

	simpleStateLenses := MakeSimpleStateLenses()

	result := F.Pipe2(
		Do(SimpleState{}),
		LetToL(
			simpleStateLenses.Value,
			42,
		),
		LetL(
			simpleStateLenses.Value,
			N.Add(8),
		),
	)

	state, err := result(context.Background())
	fmt.Printf("Value: %d, Error: %v\n", state.Value, err)
	// Output: Value: 50, Error: <nil>
}

// ExampleBind demonstrates sequencing a ReaderResult computation and updating
// the state with its result. This is the core operation for do-notation,
// allowing you to chain computations where each step can depend on the
// accumulated state and update it with new values.
//
// Step-by-step breakdown:
//
// 1. Setup lenses for accessing nested state fields:
//
//   - userLenses: Provides lenses for User fields (ID, Name)
//
//   - stateLenses: Provides lenses for State fields (User, Posts, FullName, Status)
//
//   - userIdLens: A composed lens that focuses on state.User.ID
//     Created by composing stateLenses.User with userLenses.ID
//
//     2. Do(State{}) - Initialize the do-notation chain with an empty State.
//     This creates the initial ReaderResult that will accumulate data through
//     subsequent operations.
//
//     3. ApSL(stateLenses.User, getUser(42)) - Fetch user and store in state.User field.
//     ApSL (Applicative Set Lens) executes the getUser(42) ReaderResult computation
//     and uses the lens to set the result into state.User.
//     After this step: state.User = User{ID: 42, Name: "Alice"}
//
//     4. Bind(stateLenses.Posts.Set, F.Flow2(userIdLens.Get, getPosts)) - Fetch posts
//     based on the user ID from state and store them in state.Posts.
//
//     Breaking down the Bind operation:
//     a) First parameter: stateLenses.Posts.Set - A setter function that will update
//     the Posts field in the state with the result of the computation.
//
//     b) Second parameter: F.Flow2(userIdLens.Get, getPosts) - A composed function that:
//
//   - Takes the current state as input
//
//   - Extracts the user ID using userIdLens.Get (gets state.User.ID)
//
//   - Passes the user ID to getPosts, which returns a ReaderResult[[]Post]
//
//   - The result is then set into state.Posts using the setter
//
//     After this step: state.Posts = [{ID: 1, UserID: 42, ...}, {ID: 2, UserID: 42, ...}]
//
//     5. result(context.Background()) - Execute the entire computation chain.
//     This runs all the ReaderResult operations in sequence, threading the context
//     through each step and accumulating the state.
//
// Key concepts demonstrated:
// - Lens composition: Building complex accessors from simple ones
// - Sequential effects: Each step can depend on previous results
// - State accumulation: Building up a complex state object step by step
// - Context threading: The context.Context flows through all operations
// - Error handling: Any error in the chain short-circuits execution
func ExampleBind() {

	userLenses := MakeUserLenses()
	stateLenses := MakeStateLenses()

	userIdLens := F.Pipe1(
		stateLenses.User,
		lens.Compose[State](userLenses.ID),
	)

	result := F.Pipe2(
		Do(State{}),
		ApSL(
			stateLenses.User,
			getUser(42),
		),
		Bind(
			stateLenses.Posts.Set,
			F.Flow2(
				userIdLens.Get,
				getPosts,
			),
		),
	)

	state, err := result(context.Background())
	fmt.Printf("User: %s, Posts: %d, Error: %v\n", state.User.Name, len(state.Posts), err)
	// Output: User: Alice, Posts: 2, Error: <nil>
}

// fp-go:Lens
type NameState struct {
	FirstName string
	LastName  string
	FullName  string
}

// ExampleLet demonstrates attaching the result of a pure computation to a state.
// Unlike Bind, Let works with pure functions (not ReaderResult computations).
// This is useful for deriving values from the current state without performing
// any effects.
//
// Step-by-step breakdown:
//
//  1. nameStateLenses := MakeNameStateLenses() - Create lenses for accessing NameState fields.
//     Lenses provide a functional way to get and set nested fields in immutable data structures.
//     This gives us lenses for FirstName, LastName, and FullName fields.
//
//  2. Do(NameState{FirstName: "John", LastName: "Doe"}) - Initialize the do-notation
//     chain with a NameState containing first and last names.
//     Initial state: {FirstName: "John", LastName: "Doe", FullName: ""}
//
//  3. Let(nameStateLenses.FullName.Set, func(s NameState) string {...}) - Compute a
//     derived value from the current state and update the state with it.
//
//     Let takes two parameters:
//
//     a) First parameter: nameStateLenses.FullName.Set
//     This is a setter function (from the lens) that takes a value and returns a
//     function to update the FullName field in the state. The lens-based setter
//     ensures immutable updates.
//
//     b) Second parameter: func(s NameState) string
//     This is a pure "getter" or "computation" function that derives a value from
//     the current state. Here it concatenates FirstName and LastName with a space.
//     This function has no side effects - it just computes a value.
//
//     The Let operation flow:
//     - Takes the current state: {FirstName: "John", LastName: "Doe", FullName: ""}
//     - Calls the computation function: "John" + " " + "Doe" = "John Doe"
//     - Passes "John Doe" to the setter (nameStateLenses.FullName.Set)
//     - The setter creates a new state with FullName updated
//     After this step: {FirstName: "John", LastName: "Doe", FullName: "John Doe"}
//
//  4. Map(nameStateLenses.FullName.Get) - Transform the final state to extract just
//     the FullName field using the lens getter. This changes the result type from
//     ReaderResult[NameState] to ReaderResult[string].
//
//  5. result(context.Background()) - Execute the computation chain and return the
//     final extracted value ("John Doe") and any error.
//
// Key differences between Let and Bind:
// - Let: Works with pure functions (State -> Value), no effects or errors
// - Bind: Works with effectful computations (State -> ReaderResult[Value])
// - Let: Used for deriving/computing values from existing state
// - Bind: Used for operations that may fail, need context, or have side effects
//
// Use Let when you need to:
// - Compute derived values from existing state fields
// - Transform or combine state values without side effects
// - Add computed fields to your state for later use in the pipeline
// - Perform pure calculations that don't require context or error handling
func ExampleLet() {

	nameStateLenses := MakeNameStateLenses()

	result := F.Pipe2(
		Do(NameState{FirstName: "John", LastName: "Doe"}),
		Let(nameStateLenses.FullName.Set,
			func(s NameState) string {
				return s.FirstName + " " + s.LastName
			},
		),
		Map(nameStateLenses.FullName.Get),
	)

	fullName, err := result(context.Background())
	fmt.Printf("Full Name: %s, Error: %v\n", fullName, err)
	// Output: Full Name: John Doe, Error: <nil>
}

// fp-go:Lens
type StatusState struct {
	Status string
}

// ExampleLetTo demonstrates attaching a constant value to a state.
// This is a simplified version of Let for when you want to add a constant
// value to the state without computing it.
//
// Step-by-step breakdown:
//
//  1. statusStateLenses := MakeStatusStateLenses() - Create lenses for accessing
//     StatusState fields. This provides functional accessors (getters and setters)
//     for the Status field.
//
//  2. Do(StatusState{}) - Initialize the do-notation chain with an empty StatusState.
//     Initial state: {Status: ""}
//
//  3. LetToL(statusStateLenses.Status, "active") - Set the Status field to the
//     constant value "active".
//
//     LetToL is the lens-based version of LetTo and takes two parameters:
//
//     a) First parameter: statusStateLenses.Status
//     This is a lens that focuses on the Status field. The lens provides both
//     a getter and setter for the field, enabling immutable updates.
//
//     b) Second parameter: "active"
//     This is the constant value to assign to the Status field. Unlike Let,
//     which takes a function to compute the value, LetToL directly takes the
//     value itself.
//
//     The LetToL operation:
//     - Takes the constant value "active"
//     - Uses the lens setter to create a new state with Status = "active"
//     - Returns the updated state
//     After this step: {Status: "active"}
//
//  4. Map(statusStateLenses.Status.Get) - Transform the final state to extract
//     just the Status field using the lens getter. This changes the result type
//     from ReaderResult[StatusState] to ReaderResult[string].
//
//  5. result(context.Background()) - Execute the computation chain and return
//     the final extracted value ("active") and any error.
//
// Comparison of state-setting operations:
// - LetToL: Set a field to a constant value using a lens (simplest)
// - LetL: Transform a field using a function and a lens
// - Let: Compute a value from state and update using a custom setter
// - Bind: Execute an effectful computation and update state with the result
//
// Use LetToL when you need to:
// - Set a field to a known constant value
// - Initialize state fields with default values
// - Update configuration or status flags
// - Assign literal values without any computation
//
// LetToL is the most straightforward way to set a constant value in do-notation,
// combining the simplicity of LetTo with the power of lenses for type-safe,
// immutable field updates.
func ExampleLetTo() {

	statusStateLenses := MakeStatusStateLenses()

	result := F.Pipe2(
		Do(StatusState{}),
		LetToL(
			statusStateLenses.Status,
			"active",
		),
		Map(statusStateLenses.Status.Get),
	)

	status, err := result(context.Background())
	fmt.Printf("Status: %s, Error: %v\n", status, err)
	// Output: Status: active, Error: <nil>
}

// fp-go:Lens
type UserState struct {
	User User
}

// ExampleBindTo demonstrates initializing do-notation by binding a value to a state.
// This is typically used as the first operation after a computation to
// start building up a state structure.
func ExampleBindTo() {

	userStatePrisms := MakeUserStatePrisms()

	result := F.Pipe1(
		getUser(42),
		BindToP(userStatePrisms.User),
	)

	state, err := result(context.Background())
	fmt.Printf("User: %s, Error: %v\n", state.User.Name, err)
	// Output: User: Alice, Error: <nil>
}

// fp-go:Lens
type ConfigState struct {
	Config string
}

// ExampleBindReaderK demonstrates binding a Reader computation (context-dependent
// but error-free) into a ReaderResult do-notation chain.
func ExampleBindReaderK() {

	configStateLenses := MakeConfigStateLenses()

	// A Reader that extracts a value from context
	getConfig := func(ctx context.Context) string {
		if val := ctx.Value("config"); val != nil {
			return val.(string)
		}
		return "default"
	}

	result := F.Pipe1(
		Do(ConfigState{}),
		BindReaderK(configStateLenses.Config.Set,
			func(s ConfigState) Reader[context.Context, string] {
				return getConfig
			},
		),
	)

	ctx := context.WithValue(context.Background(), "config", "production")
	state, err := result(ctx)
	fmt.Printf("Config: %s, Error: %v\n", state.Config, err)
	// Output: Config: production, Error: <nil>
}

// fp-go:Lens
type NumberState struct {
	Number int
}

// ExampleBindEitherK demonstrates binding a Result (Either) computation into
// a ReaderResult do-notation chain. This is useful for integrating pure
// error-handling logic that doesn't need context.
func ExampleBindEitherK() {

	numberStateLenses := MakeNumberStateLenses()

	// A pure function that returns a Result
	parseNumber := func(s NumberState) RES.Result[int] {
		return RES.Of(42)
	}

	result := F.Pipe1(
		Do(NumberState{}),
		BindEitherK(
			numberStateLenses.Number.Set,
			parseNumber,
		),
	)

	state, err := result(context.Background())
	fmt.Printf("Number: %d, Error: %v\n", state.Number, err)
	// Output: Number: 42, Error: <nil>
}

// fp-go:Lens
type DataState struct {
	Data string
}

// ExampleBindResultK demonstrates binding an idiomatic Go function (returning
// value and error) into a ReaderResult do-notation chain. This is particularly
// useful for integrating existing Go code that follows the standard (value, error)
// return pattern into functional pipelines.
//
// Step-by-step breakdown:
//
//  1. dataStateLenses := MakeDataStateLenses() - Create lenses for accessing
//     DataState fields. This provides functional accessors (getters and setters)
//     for the Data field, enabling type-safe, immutable field updates.
//
//  2. fetchData := func(s DataState) (string, error) - Define an idiomatic Go
//     function that takes the current state and returns a tuple of (value, error).
//
//     IMPORTANT: This function represents a PURE READER COMPOSITION - it reads from
//     the state and performs computations that don't require a context.Context.
//     This is suitable for:
//     - Pure computations that may fail (parsing, validation, calculations)
//     - Operations that only depend on the state, not external context
//     - Stateless transformations with error handling
//     - Synchronous operations that don't need cancellation or timeouts
//
//     For EFFECTFUL COMPOSITION (operations that need context), use the full
//     ReaderResult type instead: func(context.Context) (Value, error)
//     Use ReaderResult when you need:
//     - Context cancellation or timeouts
//     - Context values (request IDs, trace IDs, etc.)
//     - Operations that depend on external context state
//     - Async operations that should respect context lifecycle
//
//     In this example, fetchData always succeeds with "fetched data", but in real
//     code it might perform pure operations like:
//     - Parsing or validating data from the state
//     - Performing calculations that could fail
//     - Calling pure functions from external libraries
//     - Data transformations that don't require context
//
//  3. Do(DataState{}) - Initialize the do-notation chain with an empty DataState.
//     This creates the initial ReaderResult that will accumulate data through
//     subsequent operations.
//     Initial state: {Data: ""}
//
//  4. BindResultK(dataStateLenses.Data.Set, fetchData) - Bind the idiomatic Go
//     function into the ReaderResult chain.
//
//     BindResultK takes two parameters:
//
//     a) First parameter: dataStateLenses.Data.Set
//     This is a setter function from the lens that will update the Data field
//     with the result of the computation. The lens ensures immutable updates.
//
//     b) Second parameter: fetchData
//     This is the idiomatic Go function (State -> (Value, error)) that will be
//     lifted into the ReaderResult context.
//
//     The BindResultK operation flow:
//     - Takes the current state: {Data: ""}
//     - Calls fetchData with the state: fetchData(DataState{})
//     - Gets the result tuple: ("fetched data", nil)
//     - If error is not nil, short-circuits the chain and returns the error
//     - If error is nil, uses the setter to update state.Data with "fetched data"
//     - Returns the updated state: {Data: "fetched data"}
//     After this step: {Data: "fetched data"}
//
//  5. result(context.Background()) - Execute the computation chain with a context.
//     Even though fetchData doesn't use the context, the ReaderResult still needs
//     one to maintain the uniform interface. This runs all operations in sequence
//     and returns the final state and any error.
//
// Key concepts demonstrated:
// - Integration of idiomatic Go code: BindResultK bridges functional and imperative styles
// - Error propagation: Errors from the Go function automatically propagate through the chain
// - State transformation: The result updates the state using lens-based setters
// - Context independence: The function doesn't need context but still works in ReaderResult
//
// Comparison with other bind operations:
// - BindResultK: For idiomatic Go functions (State -> (Value, error))
// - Bind: For full ReaderResult computations (State -> ReaderResult[Value])
// - BindEitherK: For pure Result/Either values (State -> Result[Value])
// - BindReaderK: For context-dependent functions (State -> Reader[Context, Value])
//
// Use BindResultK when you need to:
// - Integrate existing Go code that returns (value, error)
// - Call functions that may fail but don't need context
// - Perform stateful computations with standard Go error handling
// - Bridge between functional pipelines and imperative Go code
// - Work with libraries that follow Go conventions
//
// Real-world example scenarios:
// - Parsing JSON from a state field: func(s State) (ParsedData, error)
// - Validating user input: func(s State) (ValidatedInput, error)
// - Performing calculations: func(s State) (Result, error)
// - Calling third-party libraries: func(s State) (APIResponse, error)
func ExampleBindResultK() {

	dataStateLenses := MakeDataStateLenses()

	// An idiomatic Go function returning (value, error)
	fetchData := func(s DataState) (string, error) {
		return "fetched data", nil
	}

	result := F.Pipe1(
		Do(DataState{}),
		BindResultK(
			dataStateLenses.Data.Set,
			fetchData,
		),
	)

	state, err := result(context.Background())
	fmt.Printf("Data: %s, Error: %v\n", state.Data, err)
	// Output: Data: fetched data, Error: <nil>
}

// fp-go:Lens
type RequestState struct {
	RequestID string
}

// ExampleBindToReader demonstrates converting a Reader computation into a
// ReaderResult and binding it to create an initial state.
func ExampleBindToReader() {
	// A Reader that extracts request ID from context
	getRequestID := func(ctx context.Context) string {
		if val := ctx.Value("requestID"); val != nil {
			return val.(string)
		}
		return "unknown"
	}

	result := F.Pipe1(
		getRequestID,
		BindToReader(func(id string) RequestState {
			return RequestState{RequestID: id}
		}),
	)

	ctx := context.WithValue(context.Background(), "requestID", "req-123")
	state, err := result(ctx)
	fmt.Printf("Request ID: %s, Error: %v\n", state.RequestID, err)
	// Output: Request ID: req-123, Error: <nil>
}

// fp-go:Lens
type ValueState struct {
	Value int
}

// ExampleBindToEither demonstrates converting a Result (Either) into a
// ReaderResult and binding it to create an initial state.
func ExampleBindToEither() {
	// A Result value
	resultValue := RES.Of(100)

	result := F.Pipe1(
		resultValue,
		BindToEither(func(v int) ValueState {
			return ValueState{Value: v}
		}),
	)

	state, err := result(context.Background())
	fmt.Printf("Value: %d, Error: %v\n", state.Value, err)
	// Output: Value: 100, Error: <nil>
}

// fp-go:Lens
type ResultState struct {
	Result string
}

// ExampleBindToResult demonstrates converting an idiomatic Go tuple (value, error)
// into a ReaderResult and binding it to create an initial state.
func ExampleBindToResult() {

	// Simulate an idiomatic Go function result
	value, err := "success", error(nil)

	result := F.Pipe1(
		BindToResult(func(v string) ResultState {
			return ResultState{Result: v}
		}),
		func(f func(string, error) ReaderResult[ResultState]) ReaderResult[ResultState] {
			return f(value, err)
		},
	)

	state, resultErr := result(context.Background())
	fmt.Printf("Result: %s, Error: %v\n", state.Result, resultErr)
	// Output: Result: success, Error: <nil>
}

// fp-go:Lens
type EnvState struct {
	Environment string
}

// ExampleApReaderS demonstrates applying a Reader computation in applicative style,
// combining it with the current state in a do-notation chain.
func ExampleApReaderS() {

	// A Reader that gets environment from context
	getEnv := func(ctx context.Context) string {
		if val := ctx.Value("env"); val != nil {
			return val.(string)
		}
		return "development"
	}

	result := F.Pipe1(
		Do(EnvState{}),
		ApReaderS(
			func(env string) Endomorphism[EnvState] {
				return func(s EnvState) EnvState {
					s.Environment = env
					return s
				}
			},
			getEnv,
		),
	)

	ctx := context.WithValue(context.Background(), "env", "staging")
	state, err := result(ctx)
	fmt.Printf("Environment: %s, Error: %v\n", state.Environment, err)
	// Output: Environment: staging, Error: <nil>
}

// fp-go:Lens
type ScoreState struct {
	Score int
}

// ExampleApEitherS demonstrates applying a Result (Either) in applicative style,
// combining it with the current state in a do-notation chain.
func ExampleApEitherS() {
	// A Result value
	scoreResult := RES.Of(95)

	result := F.Pipe1(
		Do(ScoreState{}),
		ApEitherS(
			func(score int) Endomorphism[ScoreState] {
				return func(s ScoreState) ScoreState {
					s.Score = score
					return s
				}
			},
			scoreResult,
		),
	)

	state, err := result(context.Background())
	fmt.Printf("Score: %d, Error: %v\n", state.Score, err)
	// Output: Score: 95, Error: <nil>
}

// fp-go:Lens
type MessageState struct {
	Message string
}

// ExampleApResultS demonstrates applying an idiomatic Go tuple (value, error)
// in applicative style, combining it with the current state in a do-notation chain.
func ExampleApResultS() {
	// Simulate an idiomatic Go function result
	value, err := "Hello, World!", error(nil)

	result := F.Pipe1(
		Do(MessageState{}),
		func(rr ReaderResult[MessageState]) ReaderResult[MessageState] {
			return F.Pipe1(
				rr,
				ApResultS(
					func(msg string) Endomorphism[MessageState] {
						return func(s MessageState) MessageState {
							s.Message = msg
							return s
						}
					},
				)(value, err),
			)
		},
	)

	state, resultErr := result(context.Background())
	fmt.Printf("Message: %s, Error: %v\n", state.Message, resultErr)
	// Output: Message: Hello, World!, Error: <nil>
}
