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

package formatting

import (
	"fmt"
	"reflect"
	"strings"
)

// FmtString implements the fmt.Formatter interface for Formattable types.
// It handles various format verbs to provide consistent string formatting
// across different output contexts.
//
// Supported format verbs:
//   - %v: Uses String() representation (default format)
//   - %+v: Uses String() representation (verbose format)
//   - %#v: Uses GoString() representation (Go-syntax format)
//   - %s: Uses String() representation (string format)
//   - %q: Uses quoted String() representation (quoted string format)
//   - default: Uses String() representation for any other verb
//
// The function delegates to the appropriate method of the Formattable interface
// based on the format verb and flags provided by fmt.State.
//
// Parameters:
//   - stg: The Formattable value to format
//   - f: The fmt.State that provides formatting context and flags
//   - c: The format verb (rune) being used
//
// Example usage:
//
//	type MyType struct {
//		value int
//	}
//
//	func (m MyType) Format(f fmt.State, verb rune) {
//		formatting.FmtString(m, f, verb)
//	}
//
//	func (m MyType) String() string {
//		return fmt.Sprintf("MyType(%d)", m.value)
//	}
//
//	func (m MyType) GoString() string {
//		return fmt.Sprintf("MyType{value: %d}", m.value)
//	}
//
//	// Usage:
//	mt := MyType{value: 42}
//	fmt.Printf("%v\n", mt)   // Output: MyType(42)
//	fmt.Printf("%#v\n", mt)  // Output: MyType{value: 42}
//	fmt.Printf("%s\n", mt)   // Output: MyType(42)
//	fmt.Printf("%q\n", mt)   // Output: "MyType(42)"
func FmtString(stg Formattable, f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('#') {
			// %#v uses GoString representation
			fmt.Fprint(f, stg.GoString())
		} else {
			// %v and %+v use String representation
			fmt.Fprint(f, stg.String())
		}
	case 's':
		fmt.Fprint(f, stg.String())
	case 'q':
		fmt.Fprintf(f, "%q", stg.String())
	default:
		fmt.Fprint(f, stg.String())
	}
}

// TypeInfo returns a string representation of the type of the given value.
// It uses reflection to determine the type and removes the leading asterisk (*)
// from pointer types to provide a cleaner type name.
//
// This function is useful for generating human-readable type information in
// string representations, particularly for generic types where the concrete
// type needs to be displayed.
//
// Parameters:
//   - v: The value whose type information should be extracted
//
// Returns:
//   - A string representing the type name, with pointer prefix removed
//
// Example usage:
//
//	// For non-pointer types
//	TypeInfo(42)                    // Returns: "int"
//	TypeInfo("hello")               // Returns: "string"
//	TypeInfo([]int{1, 2, 3})        // Returns: "[]int"
//
//	// For pointer types (asterisk is removed)
//	var ptr *int
//	TypeInfo(ptr)                   // Returns: "int" (not "*int")
//
//	// For custom types
//	type MyStruct struct{ Name string }
//	TypeInfo(MyStruct{})            // Returns: "formatting.MyStruct"
//	TypeInfo(&MyStruct{})           // Returns: "formatting.MyStruct" (not "*formatting.MyStruct")
//
//	// For interface types
//	var err error = fmt.Errorf("test")
//	TypeInfo(err)                   // Returns: "errors.errorString"
func TypeInfo(v any) string {
	// Remove the leading * from pointer type
	return strings.TrimPrefix(reflect.TypeOf(v).String(), "*")
}
