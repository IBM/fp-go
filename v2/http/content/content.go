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

// Package content provides constants for common HTTP Content-Type header values.
//
// These constants can be used when setting or checking Content-Type headers in HTTP
// requests and responses, ensuring consistency and avoiding typos in content type strings.
//
// Example usage:
//
//	req.Header.Set("Content-Type", content.JSON)
//	if contentType == content.TextPlain {
//	    // handle plain text
//	}
package content

const (
	// TextPlain represents the "text/plain" content type for plain text data.
	// This is commonly used for simple text responses or requests without any
	// specific formatting or structure.
	//
	// Defined in RFC 2046, Section 4.1.3: https://www.rfc-editor.org/rfc/rfc2046.html#section-4.1.3
	TextPlain = "text/plain"

	// JSON represents the "application/json" content type for JSON-encoded data.
	// This is the standard content type for JSON payloads in HTTP requests and responses.
	//
	// Defined in RFC 8259: https://www.rfc-editor.org/rfc/rfc8259.html
	JSON = "application/json"

	// Json is deprecated. Use [JSON] instead.
	//
	// Deprecated: Use JSON for consistency with Go naming conventions.
	Json = JSON

	// FormEncoded represents the "application/x-www-form-urlencoded" content type.
	// This is used for HTML form submissions where form data is encoded as key-value
	// pairs in the request body, with keys and values URL-encoded.
	//
	// Defined in HTML 4.01 Specification, Section 17.13.4:
	// https://www.w3.org/TR/html401/interact/forms.html#h-17.13.4
	// Also referenced in WHATWG HTML Living Standard:
	// https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#application/x-www-form-urlencoded-encoding-algorithm
	FormEncoded = "application/x-www-form-urlencoded"
)
