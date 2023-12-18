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

// Curry2 is a duplicate of [F.Curry2] but because of the type system it's not compatible otherwise
func Curry2[GT1 ~func(T1) T1, FCT ~func(T0, T1) T1, T0, T1 any](f FCT) func(T0) GT1 {
	return func(t0 T0) GT1 {
		return func(t1 T1) T1 {
			return f(t0, t1)
		}
	}
}

// Curry2 is a duplicate of [F.Curry2] but because of the type system it's not compatible otherwise
func Curry3[GT2 ~func(T2) T2, FCT ~func(T0, T1, T2) T2, T0, T1, T2 any](f FCT) func(T0) func(T1) GT2 {
	return func(t0 T0) func(T1) GT2 {
		return func(t1 T1) GT2 {
			return func(t2 T2) T2 {
				return f(t0, t1, t2)
			}
		}
	}
}
