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

package benchmarks

import (
	"math/rand"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	B "github.com/IBM/fp-go/v2/bytes"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randChar   = F.Pipe2(
		len(charset),
		randInt,
		IO.Map(charAt),
	)
	createRandomString = F.Flow3(
		F.Bind2of2(A.Replicate[IO.IO[byte]])(randChar),
		IO.SequenceArray[byte],
		IO.Map(B.ToString),
	)
)

func createRandom[T any](single IO.IO[T]) func(size int) IO.IO[[]T] {
	return F.Flow2(
		F.Bind2of2(A.Replicate[IO.IO[T]])(single),
		IO.SequenceArray[T],
	)
}

func charAt(idx int) byte {
	return charset[idx]
}

func randInt(count int) IO.IO[int] {
	return func() int {
		return seededRand.Intn(count)
	}
}
