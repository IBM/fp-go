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

package form

import (
	"net/url"

	A "github.com/IBM/fp-go/array"
	ENDO "github.com/IBM/fp-go/endomorphism"
	F "github.com/IBM/fp-go/function"
	L "github.com/IBM/fp-go/optics/lens"
	LA "github.com/IBM/fp-go/optics/lens/array"
	LRG "github.com/IBM/fp-go/optics/lens/record/generic"
	O "github.com/IBM/fp-go/option"
)

type (
	// FormEndomorphism returns an [ENDO.Endomorphism] that transforms a form
	FormEndomorphism = ENDO.Endomorphism[url.Values]
)

var (
	// Default is the default form field
	Default = make(url.Values)

	noField = O.None[string]()

	// AtValues is a [L.Lens] that focusses on the values of a form field
	AtValues = LRG.AtRecord[url.Values, []string]

	composeHead = F.Pipe1(
		LA.AtHead[string](),
		L.ComposeOptions[url.Values, string](A.Empty[string]()),
	)

	// AtValue is a [L.Lens] that focusses on first value in form fields
	AtValue = F.Flow2(
		AtValues,
		composeHead,
	)
)

// WithValue creates a [FormBuilder] for a certain field
func WithValue(name string) func(value string) FormEndomorphism {
	return F.Flow3(
		O.Of[string],
		AtValue(name).Set,
		ENDO.Of[func(url.Values) url.Values],
	)
}

// WithoutValue creates a [FormBuilder] that removes a field
func WithoutValue(name string) FormEndomorphism {
	return AtValue(name).Set(noField)
}
