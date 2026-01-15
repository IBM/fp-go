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

package ord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Kleisli type
func TestKleisli(t *testing.T) {
	// Create a Kleisli that produces different orderings based on input
	var orderingFactory Kleisli[string, int] = func(mode string) Ord[int] {
		if mode == "ascending" {
			return FromStrictCompare[int]()
		}
		return Reverse(FromStrictCompare[int]())
	}

	// Test ascending order
	ascOrd := orderingFactory("ascending")
	assert.Equal(t, -1, ascOrd.Compare(3, 5), "ascending: 3 < 5")
	assert.Equal(t, 1, ascOrd.Compare(5, 3), "ascending: 5 > 3")
	assert.Equal(t, 0, ascOrd.Compare(5, 5), "ascending: 5 == 5")

	// Test descending order
	descOrd := orderingFactory("descending")
	assert.Equal(t, 1, descOrd.Compare(3, 5), "descending: 3 > 5")
	assert.Equal(t, -1, descOrd.Compare(5, 3), "descending: 5 < 3")
	assert.Equal(t, 0, descOrd.Compare(5, 5), "descending: 5 == 5")
}

// Test Kleisli with complex types
func TestKleisli_ComplexType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Kleisli that creates orderings based on a field selector
	var personOrderingFactory Kleisli[string, Person] = func(field string) Ord[Person] {
		stringOrd := FromStrictCompare[string]()
		intOrd := FromStrictCompare[int]()

		switch field {
		case "name":
			return Contramap(func(p Person) string { return p.Name })(stringOrd)
		case "age":
			return Contramap(func(p Person) int { return p.Age })(intOrd)
		default:
			// Default to name ordering
			return Contramap(func(p Person) string { return p.Name })(stringOrd)
		}
	}

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}

	// Order by name
	nameOrd := personOrderingFactory("name")
	assert.Equal(t, -1, nameOrd.Compare(p1, p2), "Alice < Bob by name")

	// Order by age
	ageOrd := personOrderingFactory("age")
	assert.Equal(t, 1, ageOrd.Compare(p1, p2), "30 > 25 by age")
}

// Test Operator type
func TestOperator(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Operator that transforms Ord[int] to Ord[Person] by age
	var ageOperator Operator[int, Person] = Contramap(func(p Person) int {
		return p.Age
	})

	intOrd := FromStrictCompare[int]()
	personOrd := ageOperator(intOrd)

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}
	p3 := Person{Name: "Charlie", Age: 30}

	assert.Equal(t, 1, personOrd.Compare(p1, p2), "30 > 25")
	assert.Equal(t, -1, personOrd.Compare(p2, p1), "25 < 30")
	assert.Equal(t, 0, personOrd.Compare(p1, p3), "30 == 30")
	assert.True(t, personOrd.Equals(p1, p3), "same age")
	assert.False(t, personOrd.Equals(p1, p2), "different age")
}

// Test Operator composition
func TestOperator_Composition(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Address Address
	}

	// Create operators for different transformations
	stringOrd := FromStrictCompare[string]()

	// Operator to order Person by city
	var cityOperator Operator[string, Person] = Contramap(func(p Person) string {
		return p.Address.City
	})

	personOrd := cityOperator(stringOrd)

	p1 := Person{Name: "Alice", Address: Address{Street: "Main St", City: "Boston"}}
	p2 := Person{Name: "Bob", Address: Address{Street: "Oak Ave", City: "Austin"}}

	assert.Equal(t, 1, personOrd.Compare(p1, p2), "Boston > Austin")
	assert.Equal(t, -1, personOrd.Compare(p2, p1), "Austin < Boston")
}

// Test Operator with multiple transformations
func TestOperator_MultipleTransformations(t *testing.T) {
	type Product struct {
		Name  string
		Price float64
	}

	floatOrd := FromStrictCompare[float64]()

	// Operator to order by price
	var priceOperator Operator[float64, Product] = Contramap(func(p Product) float64 {
		return p.Price
	})

	// Operator to reverse the ordering
	var reverseOperator Operator[float64, Product] = func(o Ord[float64]) Ord[Product] {
		return priceOperator(Reverse(o))
	}

	// Order by price descending
	productOrd := reverseOperator(floatOrd)

	prod1 := Product{Name: "Widget", Price: 19.99}
	prod2 := Product{Name: "Gadget", Price: 29.99}

	assert.Equal(t, 1, productOrd.Compare(prod1, prod2), "19.99 > 29.99 (reversed)")
	assert.Equal(t, -1, productOrd.Compare(prod2, prod1), "29.99 < 19.99 (reversed)")
}

// Example test for Kleisli
func ExampleKleisli() {
	// Create a Kleisli that produces different orderings based on input
	var orderingFactory Kleisli[string, int] = func(mode string) Ord[int] {
		if mode == "ascending" {
			return FromStrictCompare[int]()
		}
		return Reverse(FromStrictCompare[int]())
	}

	ascOrd := orderingFactory("ascending")
	descOrd := orderingFactory("descending")

	println(ascOrd.Compare(5, 3))  // 1
	println(descOrd.Compare(5, 3)) // -1
}

// Example test for Operator
func ExampleOperator() {
	type Person struct {
		Name string
		Age  int
	}

	// Operator that transforms Ord[int] to Ord[Person] by age
	var ageOperator Operator[int, Person] = Contramap(func(p Person) int {
		return p.Age
	})

	intOrd := FromStrictCompare[int]()
	personOrd := ageOperator(intOrd)

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}

	result := personOrd.Compare(p1, p2)
	println(result) // 1 (30 > 25)
}
