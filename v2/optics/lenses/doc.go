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

// Package lenses provides pre-built lens and prism implementations for common data structures.
//
// This package offers ready-to-use optics (lenses and prisms) for working with regex match
// structures and URL components in a functional programming style. Lenses enable immutable
// updates to nested data structures, while prisms provide safe optional access to fields.
//
// # Overview
//
// The package includes optics for:
//   - Match structures: For working with regex match results (indexed capture groups)
//   - NamedMatch structures: For working with regex matches with named capture groups
//   - url.Userinfo: For working with URL authentication information
//
// Each structure has three variants of optics:
//   - Value lenses: Work with value types (immutable updates)
//   - Reference lenses: Work with pointer types (mutable updates)
//   - Prisms: Provide optional access treating zero values as None
//
// # Lenses vs Prisms
//
// Lenses provide guaranteed access to a field within a structure:
//   - Get: Extract a field value
//   - Set: Update a field value (returns new structure for values, mutates for pointers)
//
// Prisms provide optional access to fields that may not be present:
//   - GetOption: Try to extract a field, returning Option[T]
//   - ReverseGet: Construct a structure from a field value
//
// # Match Structures
//
// Match represents a regex match with indexed capture groups:
//
//	type Match struct {
//	    Before string    // Text before the match
//	    Groups []string  // Capture groups (index 0 is full match)
//	    After  string    // Text after the match
//	}
//
// NamedMatch represents a regex match with named capture groups:
//
//	type NamedMatch struct {
//	    Before string              // Text before the match
//	    Groups map[string]string   // Named capture groups
//	    Full   string              // Full matched text
//	    After  string              // Text after the match
//	}
//
// # Usage Examples
//
// Working with Match (value-based):
//
//	lenses := MakeMatchLenses()
//	match := Match{
//	    Before: "Hello ",
//	    Groups: []string{"world", "world"},
//	    After:  "!",
//	}
//
//	// Get a field
//	before := lenses.Before.Get(match)  // "Hello "
//
//	// Update a field (returns new Match)
//	updated := lenses.Before.Set(match, "Hi ")
//	// updated.Before is "Hi ", original match unchanged
//
//	// Use optional lens (treats empty string as None)
//	emptyMatch := Match{Before: "", Groups: []string{}, After: ""}
//	beforeOpt := lenses.BeforeO.GetOption(emptyMatch)  // None
//
// Working with Match (reference-based):
//
//	lenses := MakeMatchRefLenses()
//	match := &Match{
//	    Before: "Hello ",
//	    Groups: []string{"world"},
//	    After:  "!",
//	}
//
//	// Get a field
//	before := lenses.Before.Get(match)  // "Hello "
//
//	// Update a field (mutates the pointer)
//	lenses.Before.Set(match, "Hi ")
//	// match.Before is now "Hi "
//
//	// Use prism for optional access
//	afterOpt := lenses.AfterP.GetOption(match)  // Some("!")
//
// Working with NamedMatch:
//
//	lenses := MakeNamedMatchLenses()
//	match := NamedMatch{
//	    Before: "Email: ",
//	    Groups: map[string]string{
//	        "user":   "john",
//	        "domain": "example.com",
//	    },
//	    Full:  "john@example.com",
//	    After: "",
//	}
//
//	// Get field values
//	full := lenses.Full.Get(match)     // "john@example.com"
//	groups := lenses.Groups.Get(match) // map with user and domain
//
//	// Update a field
//	updated := lenses.Before.Set(match, "Contact: ")
//
//	// Use optional lens
//	afterOpt := lenses.AfterO.GetOption(match)  // None (empty string)
//
// Working with url.Userinfo:
//
//	lenses := MakeUserinfoRefLenses()
//	userinfo := url.UserPassword("john", "secret123")
//
//	// Get username
//	username := lenses.Username.Get(userinfo)  // "john"
//
//	// Update password (returns new Userinfo)
//	updated := lenses.Password.Set(userinfo, "newpass")
//
//	// Use optional lens for password
//	pwdOpt := lenses.PasswordO.GetOption(userinfo)  // Some("secret123")
//
//	// Handle userinfo without password
//	userOnly := url.User("alice")
//
// Working with url.URL:
//
//	lenses := MakeURLLenses()
//	u := url.URL{
//	    Scheme: "https",
//	    Host:   "example.com",
//	    Path:   "/api/v1/users",
//	}
//
//	// Get field values
//	scheme := lenses.Scheme.Get(u)  // "https"
//	host := lenses.Host.Get(u)      // "example.com"
//
//	// Update fields (returns new URL)
//	updated := lenses.Path.Set("/api/v2/users")(u)
//	// updated.Path is "/api/v2/users", original u unchanged
//
//	// Use optional lens for query string
//	queryOpt := lenses.RawQueryO.Get(u)  // None (no query string)
//
//	// Set query string
//	withQuery := lenses.RawQuery.Set("page=1&limit=10")(u)
//
// Working with url.Error:
//
//	lenses := MakeErrorLenses()
//	urlErr := url.Error{
//	    Op:  "Get",
//	    URL: "https://example.com",
//	    Err: errors.New("connection timeout"),
//	}
//
//	// Get field values
//	op := lenses.Op.Get(urlErr)       // "Get"
//	urlStr := lenses.URL.Get(urlErr)  // "https://example.com"
//	err := lenses.Err.Get(urlErr)     // error: "connection timeout"
//
//	// Update fields (returns new Error)
//	updated := lenses.Op.Set("Post")(urlErr)
//	// updated.Op is "Post", original urlErr unchanged
//	pwdOpt = lenses.PasswordO.GetOption(userOnly)  // None
//
// # Composing Optics
//
// Lenses and prisms can be composed to access nested structures:
//
//	// Compose lenses to access nested fields
//	outerLens := MakeSomeLens()
//	innerLens := MakeSomeOtherLens()
//	composed := lens.Compose(outerLens, innerLens)
//
//	// Compose prisms for optional nested access
//	outerPrism := MakeSomePrism()
//	innerPrism := MakeSomeOtherPrism()
//	composed := prism.Compose(outerPrism, innerPrism)
//
// # Optional Lenses
//
// Optional lenses (suffixed with 'O') treat zero values as None:
//   - Empty strings become None
//   - Zero values of other types become None
//   - Non-zero values become Some(value)
//
// This is useful for distinguishing between "field not set" and "field set to zero value":
//
//	lenses := MakeMatchLenses()
//	match := Match{Before: "", Groups: []string{"test"}, After: "!"}
//
//	// Regular lens returns empty string
//	before := lenses.Before.Get(match)  // ""
//
//	// Optional lens returns None
//	beforeOpt := lenses.BeforeO.GetOption(match)  // None
//
//	// Setting None clears the field
//	cleared := lenses.BeforeO.Set(match, option.None[string]())
//	// cleared.Before is ""
//
//	// Setting Some updates the field
//	updated := lenses.BeforeO.Set(match, option.Some("prefix "))
//	// updated.Before is "prefix "
//
// # Code Generation
//
// This package uses code generation for creating lens implementations.
// The generate directive at the top of this file triggers the lens generator:
//
//	//go:generate go run ../../main.go lens --dir . --filename gen_lens.go
//
// To regenerate lenses after modifying structures, run:
//
//	go generate ./optics/lenses
//
// # Performance Considerations
//
// Value-based lenses (MatchLenses, NamedMatchLenses):
//   - Create new structures on each Set operation
//   - Safe for concurrent use (immutable)
//   - Suitable for functional programming patterns
//
// Reference-based lenses (MatchRefLenses, NamedMatchRefLenses):
//   - Mutate existing structures
//   - More efficient for repeated updates
//   - Require careful handling in concurrent contexts
//
// # Related Packages
//
//   - github.com/IBM/fp-go/v2/optics/lens: Core lens functionality
//   - github.com/IBM/fp-go/v2/optics/prism: Core prism functionality
//   - github.com/IBM/fp-go/v2/optics/iso: Isomorphisms for type conversions
//   - github.com/IBM/fp-go/v2/option: Option type for optional values
//
// # See Also
//
// For more information on functional optics:
//   - Lens laws: https://github.com/IBM/fp-go/blob/main/optics/lens/README.md
//   - Prism laws: https://github.com/IBM/fp-go/blob/main/optics/prism/README.md
//   - Optics tutorial: https://github.com/IBM/fp-go/blob/main/docs/optics.md
package lenses

//go:generate go run ../../main.go lens --dir . --filename gen_lens.go
