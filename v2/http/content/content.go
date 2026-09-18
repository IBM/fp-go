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

// Package content provides constants for common HTTP media types (MIME types),
// as used in the Content-Type and Accept headers.
//
// These constants can be used when setting or checking Content-Type or Accept
// headers in HTTP requests and responses, ensuring consistency and avoiding typos
// in media type strings.
//
// The constants are grouped into text, structured data (JSON, XML, YAML, ...),
// binary, form and multipart, image, audio/video and font types.
//
// All constants are bare media types without parameters. Media types are
// registered with IANA (https://www.iana.org/assignments/media-types) and their
// syntax is defined in RFC 9110, Section 8.3.1. Parameters such as the character
// set are appended with a semicolon, for example
// content.TextPlain + "; charset=utf-8".
//
// Example usage:
//
//	req.Header.Set(headers.ContentType, content.JSON)
//	if contentType == content.TextPlain {
//	    // handle plain text
//	}
package content

// Text media types, for human readable content.
const (
	// TextPlain represents the "text/plain" content type for plain text data.
	// This is commonly used for simple text responses or requests without any
	// specific formatting or structure.
	//
	// Defined in RFC 2046, Section 4.1.3: https://www.rfc-editor.org/rfc/rfc2046.html#section-4.1.3
	TextPlain = "text/plain"

	// TextHTML represents the "text/html" content type for HTML documents.
	//
	// Defined in the WHATWG HTML Living Standard: https://html.spec.whatwg.org/multipage/iana.html#text/html
	TextHTML = "text/html"

	// TextCSS represents the "text/css" content type for Cascading Style Sheets.
	//
	// Defined in RFC 2318: https://www.rfc-editor.org/rfc/rfc2318.html
	TextCSS = "text/css"

	// TextCSV represents the "text/csv" content type for comma-separated values,
	// commonly used for tabular data exports.
	//
	// Defined in RFC 4180: https://www.rfc-editor.org/rfc/rfc4180.html
	TextCSV = "text/csv"

	// TextJavaScript represents the "text/javascript" content type for JavaScript
	// source code. It replaces the obsolete "application/javascript".
	//
	// Defined in RFC 9239: https://www.rfc-editor.org/rfc/rfc9239.html
	TextJavaScript = "text/javascript"

	// TextMarkdown represents the "text/markdown" content type for Markdown documents.
	//
	// Defined in RFC 7763: https://www.rfc-editor.org/rfc/rfc7763.html
	TextMarkdown = "text/markdown"

	// TextEventStream represents the "text/event-stream" content type used by
	// Server-Sent Events (SSE) to push a stream of events from server to client.
	//
	// Defined in the WHATWG HTML Living Standard: https://html.spec.whatwg.org/multipage/server-sent-events.html
	TextEventStream = "text/event-stream"
)

// Structured data media types (JSON, XML, YAML and related formats).
const (
	// JSON represents the "application/json" content type for JSON-encoded data.
	// This is the standard content type for JSON payloads in HTTP requests and responses.
	//
	// Defined in RFC 8259: https://www.rfc-editor.org/rfc/rfc8259.html
	JSON = "application/json"

	// Json is deprecated. Use [JSON] instead.
	//
	// Deprecated: Use JSON for consistency with Go naming conventions.
	Json = JSON

	// ProblemJSON represents the "application/problem+json" content type for
	// machine readable error details in HTTP APIs ("Problem Details").
	//
	// Defined in RFC 9457: https://www.rfc-editor.org/rfc/rfc9457.html
	ProblemJSON = "application/problem+json"

	// JSONPatch represents the "application/json-patch+json" content type for a
	// sequence of operations to apply to a JSON document, used with PATCH.
	//
	// Defined in RFC 6902: https://www.rfc-editor.org/rfc/rfc6902.html
	JSONPatch = "application/json-patch+json"

	// MergePatch represents the "application/merge-patch+json" content type for a
	// partial JSON document that is merged into the target, used with PATCH.
	//
	// Defined in RFC 7396: https://www.rfc-editor.org/rfc/rfc7396.html
	MergePatch = "application/merge-patch+json"

	// NDJSON represents the "application/x-ndjson" content type for newline
	// delimited JSON, where each line is a separate JSON value. It is commonly
	// used for streaming and bulk APIs.
	//
	// Defined by the NDJSON specification: https://github.com/ndjson/ndjson-spec
	NDJSON = "application/x-ndjson"

	// JSONLD represents the "application/ld+json" content type for JSON-based
	// Linked Data.
	//
	// Defined in the W3C JSON-LD 1.1 Recommendation: https://www.w3.org/TR/json-ld11/
	JSONLD = "application/ld+json"

	// HALJSON represents the "application/hal+json" content type for JSON
	// documents that follow the Hypertext Application Language conventions.
	//
	// Defined in draft-kelly-json-hal: https://datatracker.ietf.org/doc/html/draft-kelly-json-hal
	HALJSON = "application/hal+json"

	// JSONAPI represents the "application/vnd.api+json" content type for
	// documents that follow the JSON:API specification.
	//
	// Defined by JSON:API: https://jsonapi.org/format/
	JSONAPI = "application/vnd.api+json"

	// XML represents the "application/xml" content type for XML documents. It is
	// preferred over "text/xml" for machine processed XML.
	//
	// Defined in RFC 7303: https://www.rfc-editor.org/rfc/rfc7303.html
	XML = "application/xml"

	// TextXML represents the "text/xml" content type for XML documents that are
	// intended to be readable by casual users. Prefer XML for APIs.
	//
	// Defined in RFC 7303: https://www.rfc-editor.org/rfc/rfc7303.html
	TextXML = "text/xml"

	// ProblemXML represents the "application/problem+xml" content type, the XML
	// variant of ProblemJSON.
	//
	// Defined in RFC 9457: https://www.rfc-editor.org/rfc/rfc9457.html
	ProblemXML = "application/problem+xml"

	// SOAP represents the "application/soap+xml" content type for SOAP 1.2 messages.
	//
	// Defined in RFC 3902: https://www.rfc-editor.org/rfc/rfc3902.html
	SOAP = "application/soap+xml"

	// YAML represents the "application/yaml" content type for YAML documents.
	//
	// Defined in RFC 9512: https://www.rfc-editor.org/rfc/rfc9512.html
	YAML = "application/yaml"

	// CBOR represents the "application/cbor" content type for Concise Binary
	// Object Representation, a compact binary encoding of the JSON data model.
	//
	// Defined in RFC 8949: https://www.rfc-editor.org/rfc/rfc8949.html
	CBOR = "application/cbor"

	// Protobuf represents the "application/x-protobuf" content type for Protocol
	// Buffers encoded messages. This is the de-facto media type used by most
	// HTTP APIs that exchange protobuf payloads.
	//
	// See https://protobuf.dev/
	Protobuf = "application/x-protobuf"

	// GRPC represents the "application/grpc" content type used by gRPC over HTTP/2.
	//
	// Defined in the gRPC over HTTP/2 protocol: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	GRPC = "application/grpc"
)

// Binary and archive media types.
const (
	// OctetStream represents the "application/octet-stream" content type for
	// arbitrary binary data. It is the default for unknown binary content and
	// typically triggers a download in browsers.
	//
	// Defined in RFC 2046, Section 4.5.1: https://www.rfc-editor.org/rfc/rfc2046.html#section-4.5.1
	OctetStream = "application/octet-stream"

	// PDF represents the "application/pdf" content type for PDF documents.
	//
	// Defined in RFC 8118: https://www.rfc-editor.org/rfc/rfc8118.html
	PDF = "application/pdf"

	// ZIP represents the "application/zip" content type for ZIP archives.
	//
	// Registered with IANA: https://www.iana.org/assignments/media-types/application/zip
	ZIP = "application/zip"

	// Gzip represents the "application/gzip" content type for gzip compressed
	// files. Note that transparent compression of a response body is signalled
	// with the Content-Encoding header instead.
	//
	// Defined in RFC 6713: https://www.rfc-editor.org/rfc/rfc6713.html
	Gzip = "application/gzip"

	// WASM represents the "application/wasm" content type for WebAssembly modules.
	//
	// Defined in the W3C WebAssembly Web API: https://www.w3.org/TR/wasm-web-api-2/#media-type-registration
	WASM = "application/wasm"
)

// Form and multipart media types.
const (
	// FormEncoded represents the "application/x-www-form-urlencoded" content type.
	// This is used for HTML form submissions where form data is encoded as key-value
	// pairs in the request body, with keys and values URL-encoded.
	//
	// Defined in HTML 4.01 Specification, Section 17.13.4:
	// https://www.w3.org/TR/html401/interact/forms.html#h-17.13.4
	// Also referenced in WHATWG HTML Living Standard:
	// https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#application/x-www-form-urlencoded-encoding-algorithm
	FormEncoded = "application/x-www-form-urlencoded"

	// MultipartFormData represents the "multipart/form-data" content type for
	// form submissions that include files or binary data. The actual header
	// value carries a boundary parameter, e.g. "multipart/form-data; boundary=xyz".
	//
	// Defined in RFC 7578: https://www.rfc-editor.org/rfc/rfc7578.html
	MultipartFormData = "multipart/form-data"

	// MultipartMixed represents the "multipart/mixed" content type for a body
	// consisting of several independent parts of possibly different types.
	//
	// Defined in RFC 2046, Section 5.1.3: https://www.rfc-editor.org/rfc/rfc2046.html#section-5.1.3
	MultipartMixed = "multipart/mixed"

	// MultipartByteRanges represents the "multipart/byteranges" content type
	// used in 206 Partial Content responses that contain several ranges.
	//
	// Defined in RFC 9110, Section 14.6: https://www.rfc-editor.org/rfc/rfc9110.html#section-14.6
	MultipartByteRanges = "multipart/byteranges"
)

// Image media types.
const (
	// ImagePNG represents the "image/png" content type for PNG images.
	//
	// Defined in RFC 2083: https://www.rfc-editor.org/rfc/rfc2083.html
	ImagePNG = "image/png"

	// ImageJPEG represents the "image/jpeg" content type for JPEG images.
	//
	// Defined in RFC 2046, Section 4.2: https://www.rfc-editor.org/rfc/rfc2046.html#section-4.2
	ImageJPEG = "image/jpeg"

	// ImageGIF represents the "image/gif" content type for GIF images.
	//
	// Defined in RFC 2046, Section 4.2: https://www.rfc-editor.org/rfc/rfc2046.html#section-4.2
	ImageGIF = "image/gif"

	// ImageWebP represents the "image/webp" content type for WebP images.
	//
	// Defined in RFC 9649: https://www.rfc-editor.org/rfc/rfc9649.html
	ImageWebP = "image/webp"

	// ImageAVIF represents the "image/avif" content type for AVIF images.
	//
	// Registered with IANA: https://www.iana.org/assignments/media-types/image/avif
	ImageAVIF = "image/avif"

	// ImageSVG represents the "image/svg+xml" content type for SVG vector graphics.
	//
	// Defined in the W3C SVG 2 Recommendation: https://www.w3.org/TR/SVG2/mimereg.html
	ImageSVG = "image/svg+xml"
)

// Audio and video media types.
const (
	// AudioMPEG represents the "audio/mpeg" content type for MPEG audio such as MP3.
	//
	// Defined in RFC 3003: https://www.rfc-editor.org/rfc/rfc3003.html
	AudioMPEG = "audio/mpeg"

	// AudioOGG represents the "audio/ogg" content type for Ogg audio.
	//
	// Defined in RFC 5334: https://www.rfc-editor.org/rfc/rfc5334.html
	AudioOGG = "audio/ogg"

	// VideoMP4 represents the "video/mp4" content type for MP4 video.
	//
	// Defined in RFC 4337: https://www.rfc-editor.org/rfc/rfc4337.html
	VideoMP4 = "video/mp4"

	// VideoWebM represents the "video/webm" content type for WebM video.
	//
	// Registered with IANA: https://www.iana.org/assignments/media-types/video/webm
	VideoWebM = "video/webm"
)

// Font media types.
const (
	// FontWOFF2 represents the "font/woff2" content type for Web Open Font Format 2 fonts.
	//
	// Defined in RFC 8081: https://www.rfc-editor.org/rfc/rfc8081.html
	FontWOFF2 = "font/woff2"
)
