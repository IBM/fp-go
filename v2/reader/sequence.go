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

package reader

import (
	G "github.com/IBM/fp-go/v2/reader/generic"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded type of n strongly typed values,
// represented as a tuple. This is useful for combining multiple independent Reader computations into a single
// Reader that produces a tuple of all results.

// SequenceT1 combines 1 Reader into a Reader of a 1-tuple.
//
// Example:
//
//	type Config struct { Value int }
//	r := reader.Asks(func(c Config) int { return c.Value })
//	result := reader.SequenceT1(r)
//	tuple := result(Config{Value: 42}) // Tuple1{F1: 42}
func SequenceT1[R, A any](a Reader[R, A]) Reader[R, T.Tuple1[A]] {
	return G.SequenceT1[Reader[R, A], Reader[R, T.Tuple1[A]]](a)
}

// SequenceT2 combines 2 Readers into a Reader of a 2-tuple.
// All Readers share the same environment and are evaluated with it.
//
// Example:
//
//	type Config struct { X, Y int }
//	getX := reader.Asks(func(c Config) int { return c.X })
//	getY := reader.Asks(func(c Config) int { return c.Y })
//	result := reader.SequenceT2(getX, getY)
//	tuple := result(Config{X: 10, Y: 20}) // Tuple2{F1: 10, F2: 20}
func SequenceT2[R, A, B any](a Reader[R, A], b Reader[R, B]) Reader[R, T.Tuple2[A, B]] {
	return G.SequenceT2[Reader[R, A], Reader[R, B], Reader[R, T.Tuple2[A, B]]](a, b)
}

// SequenceT3 combines 3 Readers into a Reader of a 3-tuple.
// All Readers share the same environment and are evaluated with it.
//
// Example:
//
//	type Config struct { Host string; Port int; Secure bool }
//	getHost := reader.Asks(func(c Config) string { return c.Host })
//	getPort := reader.Asks(func(c Config) int { return c.Port })
//	getSecure := reader.Asks(func(c Config) bool { return c.Secure })
//	result := reader.SequenceT3(getHost, getPort, getSecure)
//	tuple := result(Config{Host: "localhost", Port: 8080, Secure: true})
//	// Tuple3{F1: "localhost", F2: 8080, F3: true}
func SequenceT3[R, A, B, C any](a Reader[R, A], b Reader[R, B], c Reader[R, C]) Reader[R, T.Tuple3[A, B, C]] {
	return G.SequenceT3[Reader[R, A], Reader[R, B], Reader[R, C], Reader[R, T.Tuple3[A, B, C]]](a, b, c)
}

// SequenceT4 combines 4 Readers into a Reader of a 4-tuple.
// All Readers share the same environment and are evaluated with it.
//
// Example:
//
//	type Config struct { A, B, C, D int }
//	getA := reader.Asks(func(c Config) int { return c.A })
//	getB := reader.Asks(func(c Config) int { return c.B })
//	getC := reader.Asks(func(c Config) int { return c.C })
//	getD := reader.Asks(func(c Config) int { return c.D })
//	result := reader.SequenceT4(getA, getB, getC, getD)
//	tuple := result(Config{A: 1, B: 2, C: 3, D: 4})
//	// Tuple4{F1: 1, F2: 2, F3: 3, F4: 4}
func SequenceT4[R, A, B, C, D any](a Reader[R, A], b Reader[R, B], c Reader[R, C], d Reader[R, D]) Reader[R, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[Reader[R, A], Reader[R, B], Reader[R, C], Reader[R, D], Reader[R, T.Tuple4[A, B, C, D]]](a, b, c, d)
}
