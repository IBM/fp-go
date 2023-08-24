// Copyright (c) 2023 IBM Corp.
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

package types

import (
	"fmt"
	"reflect"
	"strconv"

	E "github.com/IBM/fp-go/either"
)

func validateStringFromReflect(i reflect.Value, c Context) E.Either[Errors, string] {
	switch i.Kind() {
	case reflect.String:
		return E.Of[Errors](i.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return E.Of[Errors](strconv.FormatInt(i.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return E.Of[Errors](strconv.FormatUint(i.Uint(), 10))
	case reflect.Bool:
		return E.Of[Errors](strconv.FormatBool(i.Bool()))
	case reflect.Pointer:
		return validateStringFromReflect(i.Elem(), c)
	case reflect.Invalid:
		return Failure[string](c, "Invalid value")
	}

	if i.CanInterface() {
		if strg, ok := i.Interface().(fmt.Stringer); ok {
			return E.Of[Errors](strg.String())
		}
	}

	return E.Of[Errors](i.String())
}

// String returns the type validator for a string
func makeString() *Type[string, string, reflect.Value] {
	return FromValidate(validateStringFromReflect)
}

// converts from any type to string
var String = makeString()
