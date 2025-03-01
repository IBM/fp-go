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

package either

func Variadic0[V, R any](f func([]V) (R, error)) func(...V) Either[error, R] {
	return func(v ...V) Either[error, R] {
		return TryCatchError(f(v))
	}
}

func Variadic1[T1, V, R any](f func(T1, []V) (R, error)) func(T1, ...V) Either[error, R] {
	return func(t1 T1, v ...V) Either[error, R] {
		return TryCatchError(f(t1, v))
	}
}

func Variadic2[T1, T2, V, R any](f func(T1, T2, []V) (R, error)) func(T1, T2, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, v ...V) Either[error, R] {
		return TryCatchError(f(t1, t2, v))
	}
}

func Variadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, []V) (R, error)) func(T1, T2, T3, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, v ...V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, v))
	}
}

func Variadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, []V) (R, error)) func(T1, T2, T3, T4, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, v ...V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, t4, v))
	}
}

func Unvariadic0[V, R any](f func(...V) (R, error)) func([]V) Either[error, R] {
	return func(v []V) Either[error, R] {
		return TryCatchError(f(v...))
	}
}

func Unvariadic1[T1, V, R any](f func(T1, ...V) (R, error)) func(T1, []V) Either[error, R] {
	return func(t1 T1, v []V) Either[error, R] {
		return TryCatchError(f(t1, v...))
	}
}

func Unvariadic2[T1, T2, V, R any](f func(T1, T2, ...V) (R, error)) func(T1, T2, []V) Either[error, R] {
	return func(t1 T1, t2 T2, v []V) Either[error, R] {
		return TryCatchError(f(t1, t2, v...))
	}
}

func Unvariadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, ...V) (R, error)) func(T1, T2, T3, []V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, v []V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, v...))
	}
}

func Unvariadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, ...V) (R, error)) func(T1, T2, T3, T4, []V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, v []V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, t4, v...))
	}
}
