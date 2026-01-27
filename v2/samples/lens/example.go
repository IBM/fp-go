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

package lens

import (
	"github.com/IBM/fp-go/v2/optics/lens/option"
)

//go:generate go run ../../main.go lens --dir . --filename gen_lens.go

// fp-go:Lens
type Person struct {
	Name  string
	Age   int
	Email string
	// Optional field with pointer
	Phone *string
}

// fp-go:Lens
type Address struct {
	Street  string
	City    string
	ZipCode string
	Country string
	// Optional field
	State *string `json:"state,omitempty"`
}

// fp-go:Lens
type Company struct {
	Name    string
	Address Address
	CEO     Person
	// Optional field
	Website *string
}

// fp-go:Lens
type CompanyExtended struct {
	Company
	Extended string
}

// fp-go:Lens
type CheckOption struct {
	Name  option.Option[string]
	Value string `json:",omitempty"`
}

// fp-go:Lens
type WithGeneric[T any] struct {
	Name  string
	Value T
}

// fp-go:Lens
type DataBuilder struct {
	name  string
	value string
}
