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

package mostlyadequate

import (
	"fmt"
	"path"
	"regexp"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
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

	Chapter09User struct {
		Id      int
		Name    string
		Address Address
	}
)

var (
	albert09 = Chapter09User{
		Id:   1,
		Name: "Albert",
		Address: Address{
			Street: Street{
				Number: 22,
				Name:   "Walnut St",
			},
		},
	}

	gary09 = Chapter09User{
		Id:   2,
		Name: "Gary",
		Address: Address{
			Street: Street{
				Number: 14,
			},
		},
	}

	theresa09 = Chapter09User{
		Id:   3,
		Name: "Theresa",
	}
)

func (ab AddressBook) getAddresses() []Address {
	return ab.Addresses
}

func (s Address) getStreet() Street {
	return s.Street
}

func (s Street) getName() string {
	return s.Name
}

func (u Chapter09User) getAddress() Address {
	return u.Address
}

var (
	FirstAddressStreet = F.Flow3(
		AddressBook.getAddresses,
		A.Head[Address],
		O.Map(Address.getStreet),
	)

	// getFile :: IO String
	getFile = io.Of("/home/mostly-adequate/ch09.md")

	// pureLog :: String -> IO ()
	pureLog = io.Logf[string]("%s")

	// addToMailingList :: Email -> IOEither([Email])
	addToMailingList = F.Flow2(
		A.Of[string],
		ioresult.Of[[]string],
	)

	// validateEmail :: Email -> Either error Email
	validateEmail = result.FromPredicate(Matches(regexp.MustCompile(`\S+@\S+\.\S+`)), errors.OnSome[string]("email %s is invalid"))

	// emailBlast :: [Email] -> IO ()
	emailBlast = F.Flow2(
		A.Intercalate(S.Monoid)(","),
		ioresult.Of[string],
	)
)

func Example_street() {
	s := FirstAddressStreet(AddressBook{
		Addresses: A.From(Address{Street: Street{Name: "Mulburry", Number: 8402}, Postcode: "WC2N"}),
	})
	fmt.Println(s)

	// Output:
	// Some[mostlyadequate.Street]({Mulburry 8402})
}

func Example_solution09A() {
	// // getStreetName :: User -> Maybe String
	getStreetName := F.Flow4(
		Chapter09User.getAddress,
		Address.getStreet,
		Street.getName,
		O.FromPredicate(S.IsNonEmpty),
	)

	fmt.Println(getStreetName(albert09))
	fmt.Println(getStreetName(gary09))
	fmt.Println(getStreetName(theresa09))

	// Output:
	// Some[string](Walnut St)
	// None[string]
	// None[string]

}

func Example_solution09B() {
	logFilename := F.Flow2(
		io.Map(path.Base),
		io.ChainFirst(pureLog),
	)

	fmt.Println(logFilename(getFile)())

	// Output:
	// ch09.md
}

func Example_solution09C() {

	// // joinMailingList :: Email -> Either String (IO ())
	joinMailingList := F.Flow4(
		validateEmail,
		ioresult.FromEither[string],
		ioresult.Chain(addToMailingList),
		ioresult.Chain(emailBlast),
	)

	fmt.Println(joinMailingList("sleepy@grandpa.net")())
	fmt.Println(joinMailingList("notanemail")())

	// Output:
	// Right[string](sleepy@grandpa.net)
	// Left[*errors.errorString](email notanemail is invalid)
}
