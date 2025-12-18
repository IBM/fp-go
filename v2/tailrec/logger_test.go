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

package tailrec

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrampolineLogValue(t *testing.T) {
	t.Run("Bounce state logs correctly", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		bounce := Bounce[int](42)
		logger.Info("test", "trampoline", bounce)

		var logEntry map[string]any
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		require.NoError(t, err)

		// Check the trampoline field
		trampolineField, ok := logEntry["trampoline"].(map[string]any)
		require.True(t, ok, "trampoline field should be a map")

		// When Landed is false, only the "bouncing" field is present with the value
		assert.Equal(t, float64(42), trampolineField["bouncing"]) // JSON numbers are float64
		_, hasLanded := trampolineField["landed"]
		assert.False(t, hasLanded, "landed field should not be present when bouncing")
	})

	t.Run("Land state logs correctly", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		land := Land[int](100)
		logger.Info("test", "trampoline", land)

		var logEntry map[string]any
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		require.NoError(t, err)

		// Check the trampoline field
		trampolineField, ok := logEntry["trampoline"].(map[string]any)
		require.True(t, ok, "trampoline field should be a map")

		// When Landed is true, only the "landed" field is present with the value
		assert.Equal(t, float64(100), trampolineField["landed"]) // JSON numbers are float64
		_, hasValue := trampolineField["value"]
		assert.False(t, hasValue, "value field should not be present when landed")
	})

	t.Run("Complex type in Bounce", func(t *testing.T) {
		type State struct {
			N   int
			Acc int
		}

		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		bounce := Bounce[int](State{N: 5, Acc: 120})
		logger.Info("test", "state", bounce)

		var logEntry map[string]any
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		require.NoError(t, err)

		stateField, ok := logEntry["state"].(map[string]any)
		require.True(t, ok, "state field should be a map")

		// When Landed is false, only the "bouncing" field is present
		bouncing, ok := stateField["bouncing"].(map[string]any)
		require.True(t, ok, "bouncing should be a map")
		assert.Equal(t, float64(5), bouncing["N"])
		assert.Equal(t, float64(120), bouncing["Acc"])

		_, hasLanded := stateField["landed"]
		assert.False(t, hasLanded, "landed field should not be present when bouncing")
	})

	t.Run("String type in Land", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		land := Land[int]("completed")
		logger.Info("test", "result", land)

		var logEntry map[string]any
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		require.NoError(t, err)

		resultField, ok := logEntry["result"].(map[string]any)
		require.True(t, ok, "result field should be a map")

		// When Landed is true, only the "landed" field is present with the value
		assert.Equal(t, "completed", resultField["landed"])
		_, hasValue := resultField["value"]
		assert.False(t, hasValue, "value field should not be present when landed")
	})
}
