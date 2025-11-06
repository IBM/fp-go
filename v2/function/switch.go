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

package function

import (
	G "github.com/IBM/fp-go/v2/function/generic"
)

// Switch creates a function that applies different handlers based on a key extracted from the input.
//
// This implements a switch/case-like pattern in a functional style. Given a key extraction function,
// a map of handlers for different cases, and a default handler, it returns a function that:
//  1. Extracts a key from the input using the key function
//  2. Looks up the handler for that key in the map
//  3. Applies the handler if found, or the default handler if not
//
// This is useful for implementing polymorphic behavior, routing, or state machines in a
// functional way.
//
// Type Parameters:
//   - K: The type of the key (must be comparable for map lookup)
//   - T: The input type
//   - R: The return type
//
// Parameters:
//   - kf: A function that extracts a key from the input
//   - n: A map from keys to handler functions
//   - d: The default handler to use when the key is not found in the map
//
// Returns:
//   - A function that applies the appropriate handler based on the extracted key
//
// Example:
//
//	type Animal struct {
//	    Type string
//	    Name string
//	}
//
//	getType := func(a Animal) string { return a.Type }
//
//	handlers := map[string]func(Animal) string{
//	    "dog": func(a Animal) string { return a.Name + " barks" },
//	    "cat": func(a Animal) string { return a.Name + " meows" },
//	}
//
//	defaultHandler := func(a Animal) string {
//	    return a.Name + " makes a sound"
//	}
//
//	makeSound := Switch(getType, handlers, defaultHandler)
//
//	dog := Animal{Type: "dog", Name: "Rex"}
//	cat := Animal{Type: "cat", Name: "Whiskers"}
//	bird := Animal{Type: "bird", Name: "Tweety"}
//
//	result1 := makeSound(dog)   // "Rex barks"
//	result2 := makeSound(cat)   // "Whiskers meows"
//	result3 := makeSound(bird)  // "Tweety makes a sound"
//
// HTTP routing example:
//
//	type Request struct {
//	    Method string
//	    Path   string
//	}
//
//	getMethod := func(r Request) string { return r.Method }
//
//	routes := map[string]func(Request) string{
//	    "GET":    func(r Request) string { return "Handling GET " + r.Path },
//	    "POST":   func(r Request) string { return "Handling POST " + r.Path },
//	    "DELETE": func(r Request) string { return "Handling DELETE " + r.Path },
//	}
//
//	notFound := func(r Request) string {
//	    return "Method not allowed: " + r.Method
//	}
//
//	router := Switch(getMethod, routes, notFound)
func Switch[K comparable, T, R any](kf func(T) K, n map[K]func(T) R, d func(T) R) func(T) R {
	return G.Switch(kf, n, d)
}
