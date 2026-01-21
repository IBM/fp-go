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

package consumer

// Local transforms a Consumer by preprocessing its input through a function.
// This is the contravariant map operation for Consumers, analogous to reader.Local
// but operating on the input side rather than the output side.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Given a Consumer[R1] that consumes values of type R1, and a function f that
// converts R2 to R1, Local creates a new Consumer[R2] that:
//  1. Takes a value of type R2
//  2. Applies f to convert it to R1
//  3. Passes the result to the original Consumer[R1]
//
// This is particularly useful for adapting consumers to work with different input types,
// similar to how reader.Local adapts readers to work with different environment types.
//
// Comparison with reader.Local:
//   - reader.Local: Transforms the environment BEFORE passing it to a Reader (preprocessing input)
//   - consumer.Local: Transforms the value BEFORE passing it to a Consumer (preprocessing input)
//   - Both are contravariant operations on the input type
//   - Reader produces output, Consumer performs side effects
//
// Type Parameters:
//   - R2: The input type of the new Consumer (what you have)
//   - R1: The input type of the original Consumer (what it expects)
//
// Parameters:
//   - f: A function that converts R2 to R1 (preprocessing function)
//
// Returns:
//   - An Operator that transforms Consumer[R1] into Consumer[R2]
//
// Example - Basic type adaptation:
//
//	// Consumer that logs integers
//	logInt := func(x int) {
//	    fmt.Printf("Value: %d\n", x)
//	}
//
//	// Adapt it to consume strings by parsing them first
//	parseToInt := func(s string) int {
//	    n, _ := strconv.Atoi(s)
//	    return n
//	}
//
//	logString := consumer.Local(parseToInt)(logInt)
//	logString("42") // Logs: "Value: 42"
//
// Example - Extracting fields from structs:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	// Consumer that logs names
//	logName := func(name string) {
//	    fmt.Printf("Name: %s\n", name)
//	}
//
//	// Adapt it to consume User structs
//	extractName := func(u User) string {
//	    return u.Name
//	}
//
//	logUser := consumer.Local(extractName)(logName)
//	logUser(User{Name: "Alice", Age: 30}) // Logs: "Name: Alice"
//
// Example - Simplifying complex types:
//
//	type DetailedConfig struct {
//	    Host     string
//	    Port     int
//	    Timeout  time.Duration
//	    MaxRetry int
//	}
//
//	type SimpleConfig struct {
//	    Host string
//	    Port int
//	}
//
//	// Consumer that logs simple configs
//	logSimple := func(c SimpleConfig) {
//	    fmt.Printf("Server: %s:%d\n", c.Host, c.Port)
//	}
//
//	// Adapt it to consume detailed configs
//	simplify := func(d DetailedConfig) SimpleConfig {
//	    return SimpleConfig{Host: d.Host, Port: d.Port}
//	}
//
//	logDetailed := consumer.Local(simplify)(logSimple)
//	logDetailed(DetailedConfig{
//	    Host:     "localhost",
//	    Port:     8080,
//	    Timeout:  time.Second,
//	    MaxRetry: 3,
//	}) // Logs: "Server: localhost:8080"
//
// Example - Composing multiple transformations:
//
//	type Response struct {
//	    StatusCode int
//	    Body       string
//	}
//
//	// Consumer that logs status codes
//	logStatus := func(code int) {
//	    fmt.Printf("Status: %d\n", code)
//	}
//
//	// Extract status code from response
//	getStatus := func(r Response) int {
//	    return r.StatusCode
//	}
//
//	// Adapt to consume responses
//	logResponse := consumer.Local(getStatus)(logStatus)
//	logResponse(Response{StatusCode: 200, Body: "OK"}) // Logs: "Status: 200"
//
// Example - Using with multiple consumers:
//
//	type Event struct {
//	    Type      string
//	    Timestamp time.Time
//	    Data      map[string]any
//	}
//
//	// Consumers for different aspects
//	logType := func(t string) { fmt.Printf("Type: %s\n", t) }
//	logTime := func(t time.Time) { fmt.Printf("Time: %v\n", t) }
//
//	// Adapt them to consume events
//	logEventType := consumer.Local(func(e Event) string { return e.Type })(logType)
//	logEventTime := consumer.Local(func(e Event) time.Time { return e.Timestamp })(logTime)
//
//	event := Event{Type: "UserLogin", Timestamp: time.Now(), Data: nil}
//	logEventType(event) // Logs: "Type: UserLogin"
//	logEventTime(event) // Logs: "Time: ..."
//
// Use Cases:
//   - Type adaptation: Convert between different input types
//   - Field extraction: Extract specific fields from complex structures
//   - Data transformation: Preprocess data before consumption
//   - Interface adaptation: Adapt consumers to work with different interfaces
//   - Logging pipelines: Transform data before logging
//   - Event handling: Extract relevant data from events before processing
//
// Relationship to Reader:
// Consumer is the dual of Reader in category theory:
//   - Reader[R, A] = R -> A (produces output from environment)
//   - Consumer[A] = A -> () (consumes input, produces side effects)
//   - reader.Local transforms the environment before reading
//   - consumer.Local transforms the input before consuming
//   - Both are contravariant functors on their input type
func Local[R1, R2 any](f func(R2) R1) Operator[R1, R2] {
	return func(c Consumer[R1]) Consumer[R2] {
		return func(r2 R2) {
			c(f(r2))
		}
	}
}

// Compose is an alias for Local that emphasizes the composition aspect of consumer transformation.
// It composes a preprocessing function with a consumer, creating a new consumer that applies
// the function before consuming the value.
//
// This function is semantically identical to Local but uses terminology that may be more familiar
// to developers coming from functional programming backgrounds where "compose" is a common operation.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// The name "Compose" highlights that we're composing two operations:
//  1. The transformation function f: R2 -> R1
//  2. The consumer c: R1 -> ()
//
// Result: A composed consumer: R2 -> ()
//
// Type Parameters:
//   - R1: The input type of the original Consumer (what it expects)
//   - R2: The input type of the new Consumer (what you have)
//
// Parameters:
//   - f: A function that converts R2 to R1 (preprocessing function)
//
// Returns:
//   - An Operator that transforms Consumer[R1] into Consumer[R2]
//
// Example - Basic composition:
//
//	// Consumer that logs integers
//	logInt := func(x int) {
//	    fmt.Printf("Value: %d\n", x)
//	}
//
//	// Compose with a string-to-int parser
//	parseToInt := func(s string) int {
//	    n, _ := strconv.Atoi(s)
//	    return n
//	}
//
//	logString := consumer.Compose(parseToInt)(logInt)
//	logString("42") // Logs: "Value: 42"
//
// Example - Composing multiple transformations:
//
//	type Data struct {
//	    Value string
//	}
//
//	type Wrapper struct {
//	    Data Data
//	}
//
//	// Consumer that logs strings
//	logString := func(s string) {
//	    fmt.Println(s)
//	}
//
//	// Compose transformations step by step
//	extractData := func(w Wrapper) Data { return w.Data }
//	extractValue := func(d Data) string { return d.Value }
//
//	logData := consumer.Compose(extractValue)(logString)
//	logWrapper := consumer.Compose(extractData)(logData)
//
//	logWrapper(Wrapper{Data: Data{Value: "Hello"}}) // Logs: "Hello"
//
// Example - Function composition style:
//
//	// Compose is particularly useful when thinking in terms of function composition
//	type Request struct {
//	    Body []byte
//	}
//
//	// Consumer that processes strings
//	processString := func(s string) {
//	    fmt.Printf("Processing: %s\n", s)
//	}
//
//	// Compose byte-to-string conversion with processing
//	bytesToString := func(b []byte) string {
//	    return string(b)
//	}
//	extractBody := func(r Request) []byte {
//	    return r.Body
//	}
//
//	// Chain compositions
//	processBytes := consumer.Compose(bytesToString)(processString)
//	processRequest := consumer.Compose(extractBody)(processBytes)
//
//	processRequest(Request{Body: []byte("test")}) // Logs: "Processing: test"
//
// Relationship to Local:
//   - Compose and Local are identical in implementation
//   - Compose emphasizes the functional composition aspect
//   - Local emphasizes the environment/context transformation aspect
//   - Use Compose when thinking about function composition
//   - Use Local when thinking about adapting to different contexts
//
// Use Cases:
//   - Building processing pipelines with clear composition semantics
//   - Adapting consumers in a functional programming style
//   - Creating reusable consumer transformations
//   - Chaining multiple preprocessing steps
func Compose[R1, R2 any](f func(R2) R1) Operator[R1, R2] {
	return Local(f)
}

// Contramap is the categorical name for the contravariant functor operation on Consumers.
// It transforms a Consumer by preprocessing its input, making it the dual of the covariant
// functor's map operation.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#contravariant
//
// In category theory, a contravariant functor reverses the direction of morphisms.
// While a covariant functor maps f: A -> B to map(f): F[A] -> F[B],
// a contravariant functor maps f: A -> B to contramap(f): F[B] -> F[A].
//
// For Consumers:
//   - Consumer[A] is contravariant in A
//   - Given f: R2 -> R1, contramap(f) transforms Consumer[R1] to Consumer[R2]
//   - The direction is reversed: we go from Consumer[R1] to Consumer[R2]
//
// This is semantically identical to Local and Compose, but uses the standard
// categorical terminology that emphasizes the contravariant nature of the transformation.
//
// Type Parameters:
//   - R1: The input type of the original Consumer (what it expects)
//   - R2: The input type of the new Consumer (what you have)
//
// Parameters:
//   - f: A function that converts R2 to R1 (preprocessing function)
//
// Returns:
//   - An Operator that transforms Consumer[R1] into Consumer[R2]
//
// Example - Basic contravariant mapping:
//
//	// Consumer that logs integers
//	logInt := func(x int) {
//	    fmt.Printf("Value: %d\n", x)
//	}
//
//	// Contramap with a string-to-int parser
//	parseToInt := func(s string) int {
//	    n, _ := strconv.Atoi(s)
//	    return n
//	}
//
//	logString := consumer.Contramap(parseToInt)(logInt)
//	logString("42") // Logs: "Value: 42"
//
// Example - Demonstrating contravariance:
//
//	// In covariant functors (like Option, Array), map goes "forward":
//	// map: (A -> B) -> F[A] -> F[B]
//	//
//	// In contravariant functors (like Consumer), contramap goes "backward":
//	// contramap: (B -> A) -> F[A] -> F[B]
//
//	type Animal struct{ Name string }
//	type Dog struct{ Animal Animal; Breed string }
//
//	// Consumer for animals
//	consumeAnimal := func(a Animal) {
//	    fmt.Printf("Animal: %s\n", a.Name)
//	}
//
//	// Function from Dog to Animal (B -> A)
//	dogToAnimal := func(d Dog) Animal {
//	    return d.Animal
//	}
//
//	// Contramap creates Consumer[Dog] from Consumer[Animal]
//	// Direction is reversed: Consumer[Animal] -> Consumer[Dog]
//	consumeDog := consumer.Contramap(dogToAnimal)(consumeAnimal)
//
//	consumeDog(Dog{
//	    Animal: Animal{Name: "Buddy"},
//	    Breed:  "Golden Retriever",
//	}) // Logs: "Animal: Buddy"
//
// Example - Contravariant functor laws:
//
//	// Law 1: Identity
//	// contramap(identity) = identity
//	identity := func(x int) int { return x }
//	consumer1 := consumer.Contramap(identity)(consumeInt)
//	// consumer1 behaves identically to consumeInt
//
//	// Law 2: Composition
//	// contramap(f . g) = contramap(g) . contramap(f)
//	// Note: composition order is reversed compared to covariant map
//	f := func(s string) int { n, _ := strconv.Atoi(s); return n }
//	g := func(b bool) string { if b { return "1" } else { return "0" } }
//
//	// These two are equivalent:
//	consumer2 := consumer.Contramap(func(b bool) int { return f(g(b)) })(consumeInt)
//	consumer3 := consumer.Contramap(g)(consumer.Contramap(f)(consumeInt))
//
// Example - Practical use with type hierarchies:
//
//	type Logger interface {
//	    Log(string)
//	}
//
//	type Message struct {
//	    Text      string
//	    Timestamp time.Time
//	}
//
//	// Consumer that logs strings
//	logString := func(s string) {
//	    fmt.Println(s)
//	}
//
//	// Contramap to handle Message types
//	extractText := func(m Message) string {
//	    return fmt.Sprintf("[%s] %s", m.Timestamp.Format(time.RFC3339), m.Text)
//	}
//
//	logMessage := consumer.Contramap(extractText)(logString)
//	logMessage(Message{
//	    Text:      "Hello",
//	    Timestamp: time.Now(),
//	}) // Logs: "[2024-01-20T10:00:00Z] Hello"
//
// Relationship to Local and Compose:
//   - Contramap, Local, and Compose are identical in implementation
//   - Contramap emphasizes the categorical/theoretical aspect
//   - Local emphasizes the context transformation aspect
//   - Compose emphasizes the function composition aspect
//   - Use Contramap when working with category theory concepts
//   - Use Local when adapting to different contexts
//   - Use Compose when building functional pipelines
//
// Category Theory Background:
//   - Consumer[A] forms a contravariant functor
//   - The contravariant functor laws must hold:
//     1. contramap(id) = id
//     2. contramap(f ∘ g) = contramap(g) ∘ contramap(f)
//   - This is dual to the covariant functor (map) operation
//   - Consumers are contravariant because they consume rather than produce values
//
// Use Cases:
//   - Working with contravariant functors in a categorical style
//   - Adapting consumers to work with more specific types
//   - Building type-safe consumer transformations
//   - Implementing profunctor patterns (Consumer is a profunctor)
func Contramap[R1, R2 any](f func(R2) R1) Operator[R1, R2] {
	return Local(f)
}
