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

package tuples

import (
	"bytes"
	"io"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
	R "github.com/IBM/fp-go/v2/result"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

func sampleConvertDocx(r io.Reader) (string, map[string]string, error) {
	content, err := io.ReadAll(r)
	return string(content), map[string]string{}, err
}

func TestSampleConvertDocx1(t *testing.T) {
	// this conversion approach has the disadvantage that it exhausts the reader
	// so we cannot invoke the resulting IOEither multiple times
	convertDocx := func(r io.Reader) IOE.IOEither[error, T.Tuple2[string, map[string]string]] {
		return IOE.TryCatchError(func() (T.Tuple2[string, map[string]string], error) {
			text, meta, err := sampleConvertDocx(r)
			return T.MakeTuple2(text, meta), err
		})
	}

	rdr := strings.NewReader("abc")
	resIOE := convertDocx(rdr)

	resE := resIOE()

	assert.True(t, R.IsRight(resE))
}

func TestSampleConvertDocx2(t *testing.T) {
	// this approach assumes that `sampleConvertDocx` does not have any side effects
	// other than reading from a `Reader`. As a consequence it can be a pure function itself.
	// The disadvantage is that its input has to exist in memory which is probably not a good
	// idea for large inputs
	convertDocx := func(data []byte) Either[T.Tuple2[string, map[string]string]] {
		text, meta, err := sampleConvertDocx(bytes.NewReader(data))
		return R.TryCatchError(T.MakeTuple2(text, meta), err)
	}

	resE := convertDocx([]byte("abc"))

	assert.True(t, R.IsRight(resE))
}

// onClose closes a closeable resource
func onClose[R io.Closer](r R) IOE.IOEither[error, R] {
	return IOE.TryCatchError(func() (R, error) {
		return r, r.Close()
	})
}

// convertDocx3 takes an `acquire` function that creates an instance or a [ReaderCloser] whenever the resulting [IOEither] is invoked. Since
// we return a [Closer] the instance will be closed after use, automatically. This design makes sure that the resulting [IOEither] can be invoked
// as many times as necessary
func convertDocx3[R io.ReadCloser](acquire IOE.IOEither[error, R]) IOE.IOEither[error, T.Tuple2[string, map[string]string]] {
	return IOE.WithResource[T.Tuple2[string, map[string]string]](
		acquire,
		onClose[R])(
		func(r R) IOE.IOEither[error, T.Tuple2[string, map[string]string]] {
			return IOE.TryCatchError(func() (T.Tuple2[string, map[string]string], error) {
				text, meta, err := sampleConvertDocx(r)
				return T.MakeTuple2(text, meta), err
			})
		},
	)
}

// convertDocx4 takes an `acquire` function that creates an instance or a [Reader] whenever the resulting [IOEither] is invoked.
// This design makes sure that the resulting [IOEither] can be invoked
// as many times as necessary
func convertDocx4[R io.Reader](acquire IOE.IOEither[error, R]) IOE.IOEither[error, T.Tuple2[string, map[string]string]] {
	return F.Pipe1(
		acquire,
		IOE.Chain(func(r R) IOE.IOEither[error, T.Tuple2[string, map[string]string]] {
			return IOE.TryCatchError(func() (T.Tuple2[string, map[string]string], error) {
				text, meta, err := sampleConvertDocx(r)
				return T.MakeTuple2(text, meta), err
			})
		}),
	)
}

func TestSampleConvertDocx3(t *testing.T) {
	// IOEither that creates the reader
	acquire := IOEF.Open("./samples/data.txt")

	resIOE := convertDocx3(acquire)
	resE := resIOE()

	assert.True(t, R.IsRight(resE))
}

func TestSampleConvertDocx4(t *testing.T) {
	// IOEither that creates the reader
	acquire := IOE.FromIO[error](func() *strings.Reader {
		return strings.NewReader("abc")
	})

	resIOE := convertDocx4(acquire)
	resE := resIOE()

	assert.True(t, R.IsRight(resE))
}
