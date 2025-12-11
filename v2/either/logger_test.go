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

package either

import (
	"errors"
	"log/slog"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {

	l := Logger[error, string]()

	r := Right[error]("test")

	res := F.Pipe1(
		r,
		l("out"),
	)

	assert.Equal(t, r, res)
}

func TestToSLogAttr_Left(t *testing.T) {
	// Test with Left (error) value
	converter := ToSLogAttr[error, int]()
	testErr := errors.New("test error")
	leftValue := Left[int](testErr)

	attr := converter(leftValue)

	// Verify the attribute has the correct key
	assert.Equal(t, "error", attr.Key)
	// Verify the attribute value is the error
	assert.Equal(t, testErr, attr.Value.Any())
}

func TestToSLogAttr_Right(t *testing.T) {
	// Test with Right (success) value
	converter := ToSLogAttr[error, string]()
	rightValue := Right[error]("success value")

	attr := converter(rightValue)

	// Verify the attribute has the correct key
	assert.Equal(t, "value", attr.Key)
	// Verify the attribute value is the success value
	assert.Equal(t, "success value", attr.Value.Any())
}

func TestToSLogAttr_LeftWithCustomType(t *testing.T) {
	// Test with custom error type
	type CustomError struct {
		Code    int
		Message string
	}

	converter := ToSLogAttr[CustomError, string]()
	customErr := CustomError{Code: 404, Message: "not found"}
	leftValue := Left[string](customErr)

	attr := converter(leftValue)

	assert.Equal(t, "error", attr.Key)
	assert.Equal(t, customErr, attr.Value.Any())
}

func TestToSLogAttr_RightWithCustomType(t *testing.T) {
	// Test with custom success type
	type User struct {
		ID   int
		Name string
	}

	converter := ToSLogAttr[error, User]()
	user := User{ID: 123, Name: "Alice"}
	rightValue := Right[error](user)

	attr := converter(rightValue)

	assert.Equal(t, "value", attr.Key)
	assert.Equal(t, user, attr.Value.Any())
}

func TestToSLogAttr_InPipeline(t *testing.T) {
	// Test ToSLogAttr in a functional pipeline
	converter := ToSLogAttr[error, int]()

	// Test with successful pipeline
	successResult := F.Pipe2(
		Right[error](10),
		Map[error](N.Mul(2)),
		converter,
	)

	assert.Equal(t, "value", successResult.Key)
	// slog.Any converts int to int64
	assert.Equal(t, int64(20), successResult.Value.Any())

	// Test with failed pipeline
	testErr := errors.New("computation failed")
	failureResult := F.Pipe2(
		Left[int](testErr),
		Map[error](N.Mul(2)),
		converter,
	)

	assert.Equal(t, "error", failureResult.Key)
	assert.Equal(t, testErr, failureResult.Value.Any())
}

func TestToSLogAttr_WithNilError(t *testing.T) {
	// Test with nil error (edge case)
	converter := ToSLogAttr[error, string]()
	var nilErr error = nil
	leftValue := Left[string](nilErr)

	attr := converter(leftValue)

	assert.Equal(t, "error", attr.Key)
	assert.Nil(t, attr.Value.Any())
}

func TestToSLogAttr_WithZeroValue(t *testing.T) {
	// Test with zero value of success type
	converter := ToSLogAttr[error, int]()
	rightValue := Right[error](0)

	attr := converter(rightValue)

	assert.Equal(t, "value", attr.Key)
	// slog.Any converts int to int64
	assert.Equal(t, int64(0), attr.Value.Any())
}

func TestToSLogAttr_WithEmptyString(t *testing.T) {
	// Test with empty string as success value
	converter := ToSLogAttr[error, string]()
	rightValue := Right[error]("")

	attr := converter(rightValue)

	assert.Equal(t, "value", attr.Key)
	assert.Equal(t, "", attr.Value.Any())
}

func TestToSLogAttr_AttributeKind(t *testing.T) {
	// Verify that the returned attribute has the correct Kind
	converter := ToSLogAttr[error, string]()

	leftAttr := converter(Left[string](errors.New("error")))
	// Errors are stored as KindAny (which has value 0)
	assert.Equal(t, slog.KindAny, leftAttr.Value.Kind())

	rightAttr := converter(Right[error]("value"))
	// Strings have KindString
	assert.Equal(t, slog.KindString, rightAttr.Value.Kind())
}
