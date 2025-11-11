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

package functor

func flap[FAB ~func(A) B, A, B any](a A) func(FAB) B {
	return func(f FAB) B {
		return f(a)
	}
}

func MonadFlap[FAB ~func(A) B, A, B, HKTFAB, HKTB any](
	fmap func(HKTFAB, func(FAB) B) HKTB,

	fab HKTFAB,
	a A,
) HKTB {
	return fmap(fab, flap[FAB](a))
}

func Flap[FAB ~func(A) B, A, B, HKTFAB, HKTB any](
	fmap func(func(FAB) B) func(HKTFAB) HKTB,
	a A,
) func(HKTFAB) HKTB {
	return fmap(flap[FAB](a))
}
