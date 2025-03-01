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

package function

func Bind1st[T1, T2, R any](f func(T1, T2) R, t1 T1) func(T2) R {
	return func(t2 T2) R {
		return f(t1, t2)
	}
}
func Bind2nd[T1, T2, R any](f func(T1, T2) R, t2 T2) func(T1) R {
	return func(t1 T1) R {
		return f(t1, t2)
	}
}

// SK function (SKI combinator calculus).
func SK[T1, T2 any](_ T1, t2 T2) T2 {
	return t2
}
