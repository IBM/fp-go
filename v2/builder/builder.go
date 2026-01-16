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

// Package builder provides a generic Builder pattern interface for constructing
// complex objects with validation.
//
// The Builder pattern is useful when:
//   - Object construction requires multiple steps
//   - Construction may fail with validation errors
//   - You want to separate construction logic from the object itself
//
// Example usage:
//
//	type PersonBuilder struct {
//	    name string
//	    age  int
//	}
//
//	func (b PersonBuilder) Build() result.Result[Person] {
//	    if b.name == "" {
//	        return result.Error[Person](errors.New("name is required"))
//	    }
//	    if b.age < 0 {
//	        return result.Error[Person](errors.New("age must be non-negative"))
//	    }
//	    return result.Of(Person{Name: b.name, Age: b.age})
//	}
package builder

type (
	// Builder is a generic interface for the Builder pattern that constructs
	// objects of type T with validation.
	//
	// The Build method returns a Result[T] which can be either:
	//   - Success: containing the constructed object of type T
	//   - Error: containing an error if validation or construction fails
	//
	// This allows builders to perform validation and return meaningful errors
	// during the construction process, making it explicit that object creation
	// may fail.
	//
	// Type Parameters:
	//   - T: The type of object being built
	//
	// Example:
	//
	//	type ConfigBuilder struct {
	//	    host string
	//	    port int
	//	}
	//
	//	func (b ConfigBuilder) Build() result.Result[Config] {
	//	    if b.host == "" {
	//	        return result.Error[Config](errors.New("host is required"))
	//	    }
	//	    if b.port <= 0 || b.port > 65535 {
	//	        return result.Error[Config](errors.New("invalid port"))
	//	    }
	//	    return result.Of(Config{Host: b.host, Port: b.port})
	//	}
	Builder[T any] interface {
		// Build constructs and validates an object of type T.
		//
		// Returns:
		//   - Result[T]: A Result containing either the successfully built object
		//     or an error if validation or construction fails.
		Build() Result[T]
	}
)
