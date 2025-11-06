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

package option

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	// jsonNull is the cached representation of the `null` serialization in JSON
	jsonNull = []byte("null")
)

// Option defines a data structure that logically holds a value or not
type Option[A any] struct {
	isSome bool
	value  A
}

// optString prints some debug info for the object
//
//go:noinline
func optString(isSome bool, value any) string {
	if isSome {
		return fmt.Sprintf("Some[%T](%v)", value, value)
	}
	return fmt.Sprintf("None[%T]", value)
}

// optFormat prints some debug info for the object
//
//go:noinline
func optFormat(isSome bool, value any, f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, optString(isSome, value))
	default:
		fmt.Fprint(f, optString(isSome, value))
	}
}

// String prints some debug info for the object
func (s Option[A]) String() string {
	return optString(s.isSome, s.value)
}

// Format prints some debug info for the object
func (s Option[A]) Format(f fmt.State, c rune) {
	optFormat(s.isSome, s.value, f, c)
}

func optMarshalJSON(isSome bool, value any) ([]byte, error) {
	if isSome {
		return json.Marshal(value)
	}
	return jsonNull, nil
}

func (s Option[A]) MarshalJSON() ([]byte, error) {
	return optMarshalJSON(s.isSome, s.value)
}

// optUnmarshalJSON unmarshals the [Option] from a JSON string
//
//go:noinline
func optUnmarshalJSON(isSome *bool, value any, data []byte) error {
	// decode the value
	if bytes.Equal(data, jsonNull) {
		*isSome = false
		reflect.ValueOf(value).Elem().SetZero()
		return nil
	}
	*isSome = true
	return json.Unmarshal(data, value)
}

func (s *Option[A]) UnmarshalJSON(data []byte) error {
	return optUnmarshalJSON(&s.isSome, &s.value, data)
}

func IsNone[T any](val Option[T]) bool {
	return !val.isSome
}

func Some[T any](value T) Option[T] {
	return Option[T]{isSome: true, value: value}
}

func Of[T any](value T) Option[T] {
	return Some(value)
}

func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}

func IsSome[T any](val Option[T]) bool {
	return val.isSome
}

func MonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B {
	if IsSome(ma) {
		return onSome(ma.value)
	}
	return onNone()
}

func Unwrap[A any](ma Option[A]) (A, bool) {
	return ma.value, ma.isSome
}
