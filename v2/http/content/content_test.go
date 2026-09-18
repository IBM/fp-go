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

package content

import (
	"fmt"
	"mime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// allMediaTypes lists every media type constant defined by the package
var allMediaTypes = []string{
	TextPlain, TextHTML, TextCSS, TextCSV, TextJavaScript, TextMarkdown, TextEventStream,
	JSON, ProblemJSON, JSONPatch, MergePatch, NDJSON, JSONLD, HALJSON, JSONAPI,
	XML, TextXML, ProblemXML, SOAP, YAML, CBOR, Protobuf, GRPC,
	OctetStream, PDF, ZIP, Gzip, WASM,
	FormEncoded, MultipartFormData, MultipartMixed, MultipartByteRanges,
	ImagePNG, ImageJPEG, ImageGIF, ImageWebP, ImageAVIF, ImageSVG,
	AudioMPEG, AudioOGG, VideoMP4, VideoWebM,
	FontWOFF2,
}

// TestMediaTypesAreValid verifies that every constant is a unique, lower case,
// parameter-free media type that round-trips through mime.ParseMediaType
func TestMediaTypesAreValid(t *testing.T) {
	seen := make(map[string]bool)
	for _, mediaType := range allMediaTypes {
		t.Run(mediaType, func(t *testing.T) {
			assert.Equal(t, strings.ToLower(mediaType), mediaType, "media type must be lower case")
			assert.False(t, seen[mediaType], "media type must be unique")
			seen[mediaType] = true

			parsed, params, err := mime.ParseMediaType(mediaType)
			assert.NoError(t, err)
			assert.Equal(t, mediaType, parsed)
			assert.Empty(t, params, "media type constants must not carry parameters")
		})
	}
}

// TestJsonDeprecatedAlias verifies that the deprecated alias matches JSON
func TestJsonDeprecatedAlias(t *testing.T) {
	assert.Equal(t, JSON, Json)
}

// ExampleJSON demonstrates adding a charset parameter to a media type constant
func ExampleJSON() {
	fmt.Println(mime.FormatMediaType(JSON, map[string]string{"charset": "utf-8"}))
	// Output: application/json; charset=utf-8
}
