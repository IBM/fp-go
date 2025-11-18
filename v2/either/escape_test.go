// Copyright (c) 2025 IBM Corp.
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

package either

// Test functions to analyze escape behavior

//go:noinline
func testOf(x int) Either[error, int] {
	return Of[error](x)
}

//go:noinline
func testRight(x int) Either[error, int] {
	return Right[error](x)
}

//go:noinline
func testLeft(x int) Either[int, string] {
	return Left[string](x)
}
