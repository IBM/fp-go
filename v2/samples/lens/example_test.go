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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func TestPersonLens(t *testing.T) {
	// Create a person
	person := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	// Create lenses
	lenses := MakePersonLenses()

	// Test Get
	assert.Equal(t, "Alice", lenses.Name.Get(person))
	assert.Equal(t, 30, lenses.Age.Get(person))
	assert.Equal(t, "alice@example.com", lenses.Email.Get(person))

	// Test Set
	updated := lenses.Name.Set("Bob")(person)
	assert.Equal(t, "Bob", updated.Name)
	assert.Equal(t, 30, updated.Age)      // Other fields unchanged
	assert.Equal(t, "Alice", person.Name) // Original unchanged

	// Test Modify
	incrementAge := F.Pipe1(
		lenses.Age,
		L.Modify[Person](func(age int) int { return age + 1 }),
	)
	incremented := incrementAge(person)
	assert.Equal(t, 31, incremented.Age)
	assert.Equal(t, 30, person.Age) // Original unchanged
}

func TestCompanyLens(t *testing.T) {
	// Create a company with nested structures
	company := Company{
		Name: "Acme Corp",
		Address: Address{
			Street:  "123 Main St",
			City:    "Springfield",
			ZipCode: "12345",
			Country: "USA",
		},
		CEO: Person{
			Name:  "John Doe",
			Age:   45,
			Email: "john@acme.com",
		},
	}

	// Create lenses
	companyLenses := MakeCompanyLenses()
	addressLenses := MakeAddressLenses()
	personLenses := MakePersonLenses()

	// Test simple field access
	assert.Equal(t, "Acme Corp", companyLenses.Name.Get(company))

	// Test nested field access using composition
	cityLens := F.Pipe1(
		companyLenses.Address,
		L.Compose[Company](addressLenses.City),
	)
	assert.Equal(t, "Springfield", cityLens.Get(company))

	// Test nested field update
	updatedCompany := cityLens.Set("New York")(company)
	assert.Equal(t, "New York", updatedCompany.Address.City)
	assert.Equal(t, "Springfield", company.Address.City) // Original unchanged

	// Test deeply nested field access
	ceoNameLens := F.Pipe1(
		companyLenses.CEO,
		L.Compose[Company](personLenses.Name),
	)
	assert.Equal(t, "John Doe", ceoNameLens.Get(company))

	// Test deeply nested field update
	updatedCompany2 := ceoNameLens.Set("Jane Smith")(company)
	assert.Equal(t, "Jane Smith", updatedCompany2.CEO.Name)
	assert.Equal(t, "John Doe", company.CEO.Name) // Original unchanged
}

func TestLensComposition(t *testing.T) {
	company := Company{
		Name: "Tech Inc",
		Address: Address{
			Street:  "456 Oak Ave",
			City:    "Boston",
			ZipCode: "02101",
			Country: "USA",
		},
		CEO: Person{
			Name:  "Alice Johnson",
			Age:   50,
			Email: "alice@techinc.com",
		},
	}

	companyLenses := MakeCompanyLenses()
	personLenses := MakePersonLenses()

	// Compose lenses to access CEO's email
	ceoEmailLens := F.Pipe1(
		companyLenses.CEO,
		L.Compose[Company](personLenses.Email),
	)

	// Get the CEO's email
	email := ceoEmailLens.Get(company)
	assert.Equal(t, "alice@techinc.com", email)

	// Update the CEO's email
	updated := ceoEmailLens.Set("alice.johnson@techinc.com")(company)
	assert.Equal(t, "alice.johnson@techinc.com", updated.CEO.Email)
	assert.Equal(t, "alice@techinc.com", company.CEO.Email) // Original unchanged

	// Modify the CEO's age
	ceoAgeLens := F.Pipe1(
		companyLenses.CEO,
		L.Compose[Company](personLenses.Age),
	)

	modifyAge := F.Pipe1(
		ceoAgeLens,
		L.Modify[Company](func(age int) int { return age + 5 }),
	)
	olderCEO := modifyAge(company)
	assert.Equal(t, 55, olderCEO.CEO.Age)
	assert.Equal(t, 50, company.CEO.Age) // Original unchanged
}

func TestPersonRefLensesIdempotent(t *testing.T) {
	// Create a person pointer
	person := &Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	// Create ref lenses
	refLenses := MakePersonRefLenses()

	// Test that setting the same value returns the identical pointer (idempotent)
	// This works because Name, Age, and Email use MakeLensStrict which has equality optimization

	// Test Name field - setting same value should return same pointer
	sameName := refLenses.Name.Set("Alice")(person)
	assert.Same(t, person, sameName, "Setting Name to same value should return identical pointer")

	// Test Age field - setting same value should return same pointer
	sameAge := refLenses.Age.Set(30)(person)
	assert.Same(t, person, sameAge, "Setting Age to same value should return identical pointer")

	// Test Email field - setting same value should return same pointer
	sameEmail := refLenses.Email.Set("alice@example.com")(person)
	assert.Same(t, person, sameEmail, "Setting Email to same value should return identical pointer")

	// Test that setting a different value creates a new pointer
	differentName := refLenses.Name.Set("Bob")(person)
	assert.NotSame(t, person, differentName, "Setting Name to different value should return new pointer")
	assert.Equal(t, "Bob", differentName.Name)
	assert.Equal(t, "Alice", person.Name, "Original should be unchanged")

	differentAge := refLenses.Age.Set(31)(person)
	assert.NotSame(t, person, differentAge, "Setting Age to different value should return new pointer")
	assert.Equal(t, 31, differentAge.Age)
	assert.Equal(t, 30, person.Age, "Original should be unchanged")

	differentEmail := refLenses.Email.Set("bob@example.com")(person)
	assert.NotSame(t, person, differentEmail, "Setting Email to different value should return new pointer")
	assert.Equal(t, "bob@example.com", differentEmail.Email)
	assert.Equal(t, "alice@example.com", person.Email, "Original should be unchanged")
}
