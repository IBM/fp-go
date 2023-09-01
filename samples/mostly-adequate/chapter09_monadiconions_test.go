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

package mostlyadequate

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

type (
	Street struct {
		Name   string
		Number int
	}

	Address struct {
		Street   Street
		Postcode string
	}

	AddressBook struct {
		Addresses []Address
	}
)

func getAddresses(ab AddressBook) []Address {
	return ab.Addresses
}

func getStreet(s Address) Street {
	return s.Street
}

var FirstAddressStreet = F.Flow3(
	getAddresses,
	A.Head[Address],
	O.Map(getStreet),
)

func Example_street() {
	s := FirstAddressStreet(AddressBook{
		Addresses: A.From(Address{Street: Street{Name: "Mulburry", Number: 8402}, Postcode: "WC2N"}),
	})
	fmt.Println(s)

	// Output:
	// Some[mostlyadequate.Street]({Mulburry 8402})
}
