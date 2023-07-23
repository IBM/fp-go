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
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IF "github.com/IBM/fp-go/internal/file"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GA ~func() ET.Either[E, A],
	GR ~func() ET.Either[E, R],
	GANY ~func() ET.Either[E, ANY],
	E, R, A, ANY any](onCreate GR, onRelease func(R) GANY) func(func(R) GA) GA {
	return IF.WithResource(
		MonadChain[GR, GA, E, R, A],
		MonadFold[GA, GA, E, A, ET.Either[E, A]],
		MonadFold[GANY, GA, E, ANY, ET.Either[E, A]],
		MonadMap[GANY, GA, E, ANY, A],
		Left[GA, E, A],
	)(F.Constant(onCreate), onRelease)
}
