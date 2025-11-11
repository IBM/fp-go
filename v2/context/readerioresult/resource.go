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

package readerioresult

import (
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource.
// This implements the RAII (Resource Acquisition Is Initialization) pattern, ensuring that resources are
// properly released even if the operation fails or the context is canceled.
//
// The resource is created, used, and released in a safe manner:
//   - onCreate: Creates the resource
//   - The provided function uses the resource
//   - onRelease: Releases the resource (always called, even on error)
//
// Parameters:
//   - onCreate: ReaderIOResult that creates the resource
//   - onRelease: Function to release the resource
//
// Returns a function that takes a resource-using function and returns a ReaderIOResult.
//
// Example:
//
//	file := WithResource(
//	    openFile("data.txt"),
//	    func(f *os.File) ReaderIOResult[any] {
//	        return TryCatch(func(ctx context.Context) func() (any, error) {
//	            return func() (any, error) { return nil, f.Close() }
//	        })
//	    },
//	)
//	result := file(func(f *os.File) ReaderIOResult[string] {
//	    return TryCatch(func(ctx context.Context) func() (string, error) {
//	        return func() (string, error) {
//	            data, err := io.ReadAll(f)
//	            return string(data), err
//	        }
//	    })
//	})
func WithResource[A, R, ANY any](onCreate ReaderIOResult[R], onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A] {
	return RIOR.WithResource[A](onCreate, onRelease)
}
