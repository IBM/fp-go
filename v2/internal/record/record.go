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

package record

func Reduce[M ~map[K]V, K comparable, V, R any](r M, f func(R, V) R, initial R) R {
	current := initial
	for _, v := range r {
		current = f(current, v)
	}
	return current
}

func ReduceWithIndex[M ~map[K]V, K comparable, V, R any](r M, f func(K, R, V) R, initial R) R {
	current := initial
	for k, v := range r {
		current = f(k, current, v)
	}
	return current
}

func ReduceRef[M ~map[K]V, K comparable, V, R any](r M, f func(R, *V) R, initial R) R {
	current := initial
	for _, v := range r {
		current = f(current, &v) // #nosec G601
	}
	return current
}

func ReduceRefWithIndex[M ~map[K]V, K comparable, V, R any](r M, f func(K, R, *V) R, initial R) R {
	current := initial
	for k, v := range r {
		current = f(k, current, &v) // #nosec G601
	}
	return current
}
