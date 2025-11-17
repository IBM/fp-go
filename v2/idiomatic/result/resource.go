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

package result

// WithResource constructs a function that creates a resource, operates on it, and then releases it.
// This ensures proper resource cleanup even if operations fail.
// The resource is released immediately after the operation completes.
//
// Parameters:
//   - onCreate: Function to create/acquire the resource
//   - onRelease: Function to release/cleanup the resource
//
// Returns a function that takes an operation to perform on the resource.
//
// Example:
//
//	withFile := either.WithResource(
//	    func() either.Either[error, *os.File] {
//	        return either.TryCatchError(os.Open("file.txt"))
//	    },
//	    func(f *os.File) either.Either[error, any] {
//	        return either.TryCatchError(f.Close())
//	    },
//	)
//	result := withFile(func(f *os.File) either.Either[error, string] {
//	    // Use file here
//	    return either.Right[error]("data")
//	})
func WithResource[R, A, ANY any](onCreate func() (R, error), onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A] {

	return func(f func(R) (A, error)) (A, error) {
		r, rerr := onCreate()
		if rerr != nil {
			return Left[A](rerr)
		}
		a, aerr := f(r)
		_, nerr := onRelease(r)
		if aerr != nil {
			return Left[A](aerr)
		}
		if nerr != nil {
			return Left[A](nerr)

		}
		return Of(a)
	}
}
