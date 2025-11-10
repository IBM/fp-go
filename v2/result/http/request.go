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

// Package http provides utilities for creating HTTP requests with Either-based error handling.
package http

import (
	"github.com/IBM/fp-go/v2/either/http"
)

var (
	// PostRequest creates a POST HTTP request with a body.
	// Usage: PostRequest(url)(body) returns Result[*http.Request]
	//
	// Example:
	//
	//	request := http.PostRequest("https://api.example.com/data")([]byte(`{"key":"value"}`))
	PostRequest = http.PostRequest

	// PutRequest creates a PUT HTTP request with a body.
	// Usage: PutRequest(url)(body) returns Result[*http.Request]
	PutRequest = http.PutRequest

	// GetRequest creates a GET HTTP request without a body.
	// Usage: GetRequest(url) returns Result[*http.Request]
	//
	// Example:
	//
	//	request := http.GetRequest("https://api.example.com/data")
	GetRequest = http.GetRequest

	// DeleteRequest creates a DELETE HTTP request without a body.
	// Usage: DeleteRequest(url) returns Result[*http.Request]
	DeleteRequest = http.DeleteRequest

	// OptionsRequest creates an OPTIONS HTTP request without a body.
	// Usage: OptionsRequest(url) returns Result[*http.Request]
	OptionsRequest = http.OptionsRequest

	// HeadRequest creates a HEAD HTTP request without a body.
	// Usage: HeadRequest(url) returns Result[*http.Request]
	HeadRequest = http.HeadRequest
)
