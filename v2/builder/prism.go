// Copyright (c) 2024 IBM Corp.
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

package builder

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/result"
)

// BuilderPrism creates a [Prism] that converts between a builder and its built type.
//
// A Prism is an optic that focuses on a case of a sum type, providing bidirectional
// conversion with the possibility of failure. This function creates a prism that:
//   - Extracts: Attempts to build the object from the builder (may fail)
//   - Constructs: Creates a builder from a valid object (always succeeds)
//
// The extraction direction (builder -> object) uses the Build method and converts
// the Result to an Option, where errors become None. The construction direction
// (object -> builder) uses the provided creator function.
//
// Type Parameters:
//   - T: The type of the object being built
//   - B: The builder type that implements Builder[T]
//
// Parameters:
//   - creator: A function that creates a builder from a valid object of type T.
//     This function should initialize the builder with all fields from the object.
//
// Returns:
//   - Prism[B, T]: A prism that can convert between the builder and the built type.
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	type PersonBuilder struct {
//	    name string
//	    age  int
//	}
//
//	func (b PersonBuilder) Build() result.Result[Person] {
//	    if b.name == "" {
//	        return result.Error[Person](errors.New("name required"))
//	    }
//	    return result.Of(Person{Name: b.name, Age: b.age})
//	}
//
//	func NewPersonBuilder(p Person) PersonBuilder {
//	    return PersonBuilder{name: p.Name, age: p.Age}
//	}
//
//	// Create a prism for PersonBuilder
//	prism := BuilderPrism(NewPersonBuilder)
//
//	// Use the prism to extract a Person from a valid builder
//	builder := PersonBuilder{name: "Alice", age: 30}
//	person := prism.GetOption(builder) // Some(Person{Name: "Alice", Age: 30})
//
//	// Use the prism to create a builder from a Person
//	p := Person{Name: "Bob", Age: 25}
//	b := prism.ReverseGet(p) // PersonBuilder{name: "Bob", age: 25}
func BuilderPrism[T any, B Builder[T]](creator func(T) B) Prism[B, T] {
	return prism.MakePrismWithName(F.Flow2(B.Build, result.ToOption[T]), creator, "BuilderPrism")
}
