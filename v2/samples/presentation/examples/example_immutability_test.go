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

package examples

import (
	"fmt"
	"strings"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
)

type Person struct {
	name string
	age  int
}

func (p Person) GetName() string {
	return p.name
}

func (p Person) GetAge() int {
	return p.age
}

func (p Person) SetName(name string) Person {
	p.name = name
	return p
}

func (p Person) SetAge(age int) Person {
	p.age = age
	return p
}

type Address struct {
	city string
}

func (a Address) GetCity() string {
	return a.city
}

func (a Address) SetCity(city string) Address {
	a.city = city
	return a
}

type Client struct {
	person  Person
	address Address
}

func (c Client) GetPerson() Person {
	return c.person
}

func (c Client) SetPerson(person Person) Client {
	c.person = person
	return c
}

func (c Client) GetAddress() Address {
	return c.address
}

func (c Client) SetAddress(address Address) Client {
	c.address = address
	return c
}

func MakePerson(name string, age int) Person {
	return Person{name, age}
}

func MakeClient(city, name string, age int) Client {
	return Client{person: Person{name, age}, address: Address{city}}
}

func Example_immutability_struct() {
	p1 := MakePerson("Carsten", 53)

	// func(int) func(Person) Person
	setAge := F.Bind2of2(Person.SetAge)

	p2 := F.Pipe1(
		p1,
		setAge(54),
	)

	fmt.Println(p1)
	fmt.Println(p2)

	// Output:
	// {Carsten 53}
	// {Carsten 54}
}

func Example_immutability_optics() {

	// Lens[Person, int]
	ageLens := L.MakeLens(Person.GetAge, Person.SetAge)
	// func(Person) Person
	incAge := L.Modify[Person](N.Inc[int])(ageLens)

	p1 := MakePerson("Carsten", 53)
	p2 := incAge(p1)

	fmt.Println(p1)
	fmt.Println(p2)

	// Output:
	// {Carsten 53}
	// {Carsten 54}
}

func Example_immutability_lenses() {

	// Lens[Person, string]
	nameLens := L.MakeLens(Person.GetName, Person.SetName)
	// Lens[Client, Person]
	personLens := L.MakeLens(Client.GetPerson, Client.SetPerson)

	// Lens[Client, string]
	clientNameLens := F.Pipe1(
		personLens,
		L.Compose[Client](nameLens),
	)
	// func(Client) Client
	upperName := L.Modify[Client](strings.ToUpper)(clientNameLens)

	c1 := MakeClient("Böblingen", "Carsten", 53)

	c2 := upperName(c1)

	fmt.Println(c1)
	fmt.Println(c2)

	// Output:
	// {{Carsten 53} {Böblingen}}
	// {{CARSTEN 53} {Böblingen}}
}
