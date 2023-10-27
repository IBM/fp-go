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

package functor

func MonadFlap[FAB ~func(A) B, A, B, HKTFAB, HKTB any](
	fmap func(HKTFAB, func(FAB) B) HKTB,

	fab HKTFAB,
	a A,
) HKTB {
	return fmap(fab, func(f FAB) B {
		return f(a)
	})
}

func Flap[FAB ~func(A) B, A, B, HKTFAB, HKTB any](
	fmap func(HKTFAB, func(FAB) B) HKTB,

	a A,
) func(HKTFAB) HKTB {
	return func(fab HKTFAB) HKTB {
		return MonadFlap(fmap, fab, a)
	}
}
