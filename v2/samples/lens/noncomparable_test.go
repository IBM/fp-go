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

// TestDocumentLensNonComparableFields verifies that the generated value lenses
// work for fields whose types are not comparable: slices, named slice types,
// maps and functions.
func TestDocumentLensNonComparableFields(t *testing.T) {
	lenses := MakeDocumentLenses()

	doc := Document{
		Title:    "Report",
		Tags:     Tags{"draft"},
		Lines:    []string{"first"},
		Metadata: map[string]string{"author": "alice"},
	}

	// Get
	assert.Equal(t, Tags{"draft"}, lenses.Tags.Get(doc))
	assert.Equal(t, []string{"first"}, lenses.Lines.Get(doc))
	assert.Equal(t, map[string]string{"author": "alice"}, lenses.Metadata.Get(doc))

	// Set returns a new value and leaves the original untouched
	updated := lenses.Lines.Set([]string{"second", "third"})(doc)
	assert.Equal(t, []string{"second", "third"}, updated.Lines)
	assert.Equal(t, []string{"first"}, doc.Lines)

	// A function valued field round trips
	withRender := lenses.Render.Set(func(s string) string { return s + "!" })(doc)
	assert.Equal(t, "Report!", lenses.Render.Get(withRender)("Report"))
	assert.Nil(t, lenses.Render.Get(doc))
}

// TestDocumentRefLensNonComparableFields verifies the reference lenses. For
// non-comparable fields they are built with lens.MakeLensRef, which copies the
// struct instead of comparing the old and the new value.
func TestDocumentRefLensNonComparableFields(t *testing.T) {
	lenses := MakeDocumentRefLenses()

	doc := &Document{
		Title: "Report",
		Tags:  Tags{"draft"},
		Lines: []string{"first"},
	}

	updated := lenses.Tags.Set(Tags{"final"})(doc)
	assert.Equal(t, Tags{"final"}, updated.Tags)
	assert.Equal(t, Tags{"draft"}, doc.Tags, "the original must not be modified")

	// A payload holding a non-comparable value must not panic: the strict lens
	// would compare the two interface values, MakeLensRef does not.
	withPayload := lenses.Payload.Set([]int{1, 2, 3})(doc)
	assert.Equal(t, []int{1, 2, 3}, withPayload.Payload)
	assert.Equal(t, []int{4}, lenses.Payload.Set([]int{4})(withPayload).Payload)
}

// TestDocumentLensComparableFields verifies that the optional lenses are still
// generated for the comparable fields of a struct that also has non-comparable
// ones.
func TestDocumentLensComparableFields(t *testing.T) {
	lenses := MakeDocumentLenses()

	// The zero value of a comparable field maps to None
	assert.Equal(t, O.None[string](), lenses.TitleO.Get(Document{}))
	assert.Equal(t, O.Of("Report"), lenses.TitleO.Get(Document{Title: "Report"}))

	// Arrays keep their length, so the optional lens focuses on [4]byte
	checksum := [4]byte{1, 2, 3, 4}
	doc := F.Pipe1(Document{}, lenses.ChecksumO.Set(O.Of(checksum)))
	assert.Equal(t, checksum, doc.Checksum)
	assert.Equal(t, O.Of(checksum), lenses.ChecksumO.Get(doc))
	assert.Equal(t, O.None[[4]byte](), lenses.ChecksumO.Get(Document{}))
}

// TestDocumentPrismNonComparableFields verifies that prisms for
// non-comparable fields always succeed, since there is no zero value to
// compare against.
func TestDocumentPrismNonComparableFields(t *testing.T) {
	prisms := MakeDocumentPrisms()

	// A nil slice is still Some: the prism cannot test for the zero value
	assert.Equal(t, O.Of([]string(nil)), prisms.Lines.GetOption(Document{}))
	assert.Equal(t, O.Of([]string{"a"}), prisms.Lines.GetOption(Document{Lines: []string{"a"}}))

	// A comparable field keeps the zero value check
	assert.Equal(t, O.None[string](), prisms.Title.GetOption(Document{}))
	assert.Equal(t, O.Of("Report"), prisms.Title.GetOption(Document{Title: "Report"}))

	// ReverseGet builds a document from the focused value
	assert.Equal(t, []string{"a"}, prisms.Lines.ReverseGet([]string{"a"}).Lines)
}
