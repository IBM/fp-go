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

package ioeither

import (
	"fmt"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {

	type SomeData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	src := &SomeData{Key: "key", Value: "value"}

	res := F.Pipe1(
		Of[error](src),
		ChainFirst(LogJSON[*SomeData]("Data: \n%s")),
	)

	dst := res()
	assert.Equal(t, E.Of[error](src), dst)
}

func TestLogEntryExit(t *testing.T) {

	t.Run("fast and successful", func(t *testing.T) {

		data := F.Pipe2(
			Of[error]("test"),
			ChainIOK[error](io.Logf[string]("Data: %s")),
			LogEntryExit[error, string]("fast"),
		)

		assert.Equal(t, E.Of[error]("test"), data())
	})

	t.Run("slow and successful", func(t *testing.T) {

		data := F.Pipe3(
			Of[error]("test"),
			Delay[error, string](1*time.Second),
			ChainIOK[error](io.Logf[string]("Data: %s")),
			LogEntryExit[error, string]("slow"),
		)

		assert.Equal(t, E.Of[error]("test"), data())
	})

	t.Run("with error", func(t *testing.T) {

		err := fmt.Errorf("failure")

		data := F.Pipe2(
			Left[string](err),
			ChainIOK[error](io.Logf[string]("Data: %s")),
			LogEntryExit[error, string]("error"),
		)

		assert.Equal(t, E.Left[string](err), data())
	})
}
