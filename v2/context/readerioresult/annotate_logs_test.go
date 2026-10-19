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

package readerioresult

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// bufferLoggerContext returns a context whose logger writes text records into buf.
func bufferLoggerContext(ctx context.Context, buf *bytes.Buffer) context.Context {
	return pair.Second(logging.WithLogger(slog.New(slog.NewTextHandler(buf, nil)))(ctx))
}

func TestAnnotateLogs_Success(t *testing.T) {
	t.Run("adds the attributes to logs inside the computation", func(t *testing.T) {
		var buf bytes.Buffer
		ctx := bufferLoggerContext(t.Context(), &buf)

		res := F.Pipe2(
			Of(42),
			TapSLog[int]("inside"),
			AnnotateLogs[int]("requestID", "r-1"),
		)(ctx)()

		assert.Equal(t, result.Of(42), res)
		assert.Contains(t, buf.String(), "msg=inside")
		assert.Contains(t, buf.String(), "requestID=r-1")
	})

	t.Run("nested annotations accumulate", func(t *testing.T) {
		var buf bytes.Buffer
		ctx := bufferLoggerContext(t.Context(), &buf)

		F.Pipe3(
			Of("x"),
			TapSLog[string]("inner"),
			AnnotateLogs[string]("span", "inner"),
			AnnotateLogs[string](slog.String("service", "users")),
		)(ctx)()

		assert.Contains(t, buf.String(), "service=users span=inner")
	})
}

func TestAnnotateLogs_EdgeCases(t *testing.T) {
	t.Run("does not affect logging outside the computation", func(t *testing.T) {
		var buf bytes.Buffer
		ctx := bufferLoggerContext(t.Context(), &buf)

		F.Pipe2(
			Of(1),
			AnnotateLogs[int]("requestID", "r-1"),
			TapSLog[int]("outside"),
		)(ctx)()

		assert.Contains(t, buf.String(), "msg=outside")
		assert.NotContains(t, buf.String(), "requestID")
	})

	t.Run("falls back to the global logger", func(t *testing.T) {
		var buf bytes.Buffer
		old := logging.SetLogger(slog.New(slog.NewTextHandler(&buf, nil)))
		defer logging.SetLogger(old)

		F.Pipe2(
			Of(1),
			TapSLog[int]("global"),
			AnnotateLogs[int]("k", "v"),
		)(context.Background())()

		assert.Contains(t, buf.String(), "k=v")
	})
}
