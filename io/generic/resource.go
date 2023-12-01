// Copyright (c) 2023 IBM Corp.
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

package generic

import (
	F "github.com/IBM/fp-go/function"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GA ~func() A,
	GR ~func() R,
	GANY ~func() ANY,
	R, A, ANY any](onCreate GR, onRelease func(R) GANY) func(func(R) GA) GA {
	// simply map to implementation of bracket
	return F.Bind13of3(Bracket[GR, GA, GANY, R, A, ANY])(onCreate, F.Ignore2of2[A](onRelease))
}
