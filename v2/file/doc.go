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

// Package file provides functional programming utilities for working with file paths
// and I/O interfaces in Go.
//
// # Overview
//
// This package offers a collection of utility functions designed to work seamlessly
// with functional programming patterns, particularly with the fp-go library's pipe
// and composition utilities.
//
// # Path Manipulation
//
// The Join function provides a curried approach to path joining, making it easy to
// create reusable path builders:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/file"
//	)
//
//	// Create a reusable path builder
//	addConfig := file.Join("config.json")
//	configPath := addConfig("/etc/myapp")
//	// Result: "/etc/myapp/config.json"
//
//	// Use in a functional pipeline
//	logPath := F.Pipe1("/var/log", file.Join("app.log"))
//	// Result: "/var/log/app.log"
//
//	// Chain multiple joins
//	deepPath := F.Pipe2(
//	    "/root",
//	    file.Join("subdir"),
//	    file.Join("file.txt"),
//	)
//	// Result: "/root/subdir/file.txt"
//
// # I/O Interface Conversions
//
// The package provides generic type conversion functions for common I/O interfaces.
// These are useful for type erasure when you need to work with interface types
// rather than concrete implementations:
//
//	import (
//	    "bytes"
//	    "io"
//	    "github.com/IBM/fp-go/v2/file"
//	)
//
//	// Convert concrete types to interfaces
//	buf := bytes.NewBuffer([]byte("hello"))
//	var reader io.Reader = file.ToReader(buf)
//
//	writer := &bytes.Buffer{}
//	var w io.Writer = file.ToWriter(writer)
//
//	f, _ := os.Open("file.txt")
//	var closer io.Closer = file.ToCloser(f)
//	defer closer.Close()
//
// # Design Philosophy
//
// The functions in this package follow functional programming principles:
//
//   - Currying: Functions like Join return functions, enabling partial application
//   - Type Safety: Generic functions maintain type safety while providing flexibility
//   - Composability: All functions work well with fp-go's pipe and composition utilities
//   - Immutability: Functions don't modify their inputs
//
// # Performance
//
// The type conversion functions (ToReader, ToWriter, ToCloser) have zero overhead
// as they simply return their input cast to the interface type. The Join function
// uses Go's standard filepath.Join internally, ensuring cross-platform compatibility.
package file
