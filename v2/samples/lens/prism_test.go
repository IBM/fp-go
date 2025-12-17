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
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestCompanyExtendedPrismWebsite tests that we can create a CompanyExtended
// with only the Website field set using prisms
func TestCompanyExtendedPrismWebsite(t *testing.T) {
	prisms := MakeCompanyExtendedPrisms()

	website := "https://example.com"

	// Use the Website prism to create a CompanyExtended with only Website set
	result := prisms.Website.ReverseGet(&website)

	// Verify the Website field is set
	assert.NotNil(t, result.Website)
	assert.Equal(t, website, *result.Website)

	// Verify other fields are zero values
	assert.Equal(t, "", result.Name)
	assert.Equal(t, "", result.Extended)
	assert.Equal(t, Address{}, result.Address)
	assert.Equal(t, Person{}, result.CEO)
}

// TestCompanyExtendedPrismName tests that we can create a CompanyExtended
// with only the Name field set (from embedded Company struct)
func TestCompanyExtendedPrismName(t *testing.T) {
	prisms := MakeCompanyExtendedPrisms()

	name := "Acme Corp"

	// Use the Name prism to create a CompanyExtended with only Name set
	result := prisms.Name.ReverseGet(name)

	// Verify the Name field is set (from embedded Company)
	assert.Equal(t, name, result.Name)

	// Verify other fields are zero values
	assert.Nil(t, result.Website)
	assert.Equal(t, "", result.Extended)
	assert.Equal(t, Address{}, result.Address)
	assert.Equal(t, Person{}, result.CEO)
}

// TestCompanyExtendedPrismExtended tests that we can create a CompanyExtended
// with only the Extended field set
func TestCompanyExtendedPrismExtended(t *testing.T) {
	prisms := MakeCompanyExtendedPrisms()

	extended := "Extra Info"

	// Use the Extended prism to create a CompanyExtended with only Extended set
	result := prisms.Extended.ReverseGet(extended)

	// Verify the Extended field is set
	assert.Equal(t, extended, result.Extended)

	// Verify other fields are zero values
	assert.Nil(t, result.Website)
	assert.Equal(t, "", result.Name)
	assert.Equal(t, Address{}, result.Address)
	assert.Equal(t, Person{}, result.CEO)
}

// TestCompanyExtendedPrismGetOption tests that GetOption works correctly
func TestCompanyExtendedPrismGetOption(t *testing.T) {
	prisms := MakeCompanyExtendedPrisms()

	t.Run("GetOption returns Some for non-zero Name", func(t *testing.T) {
		company := CompanyExtended{
			Company: Company{
				Name: "Test Corp",
			},
			Extended: "Info",
		}

		result := prisms.Name.GetOption(company)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "Test Corp", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("GetOption returns None for zero Name", func(t *testing.T) {
		company := CompanyExtended{
			Extended: "Info",
		}

		result := prisms.Name.GetOption(company)
		assert.True(t, O.IsNone(result))
	})

	t.Run("GetOption returns Some for non-nil Website", func(t *testing.T) {
		website := "https://example.com"
		company := CompanyExtended{
			Company: Company{
				Website: &website,
			},
		}

		result := prisms.Website.GetOption(company)
		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant[*string](nil))(result)
		assert.NotNil(t, extracted)
		assert.Equal(t, website, *extracted)
	})

	t.Run("GetOption returns None for nil Website", func(t *testing.T) {
		company := CompanyExtended{}

		result := prisms.Website.GetOption(company)
		assert.True(t, O.IsNone(result))
	})
}

// TestCompanyPrismWebsite tests that Company prisms work correctly
func TestCompanyPrismWebsite(t *testing.T) {
	prisms := MakeCompanyPrisms()

	website := "https://company.com"

	// Use the Website prism to create a Company with only Website set
	result := prisms.Website.ReverseGet(&website)

	// Verify the Website field is set
	assert.NotNil(t, result.Website)
	assert.Equal(t, website, *result.Website)

	// Verify other fields are zero values
	assert.Equal(t, "", result.Name)
	assert.Equal(t, Address{}, result.Address)
	assert.Equal(t, Person{}, result.CEO)
}

// TestPersonPrismName tests that Person prisms work correctly
func TestPersonPrismName(t *testing.T) {
	prisms := MakePersonPrisms()

	name := "John Doe"

	// Use the Name prism to create a Person with only Name set
	result := prisms.Name.ReverseGet(name)

	// Verify the Name field is set
	assert.Equal(t, name, result.Name)

	// Verify other fields are zero values
	assert.Equal(t, 0, result.Age)
	assert.Equal(t, "", result.Email)
	assert.Nil(t, result.Phone)
}
