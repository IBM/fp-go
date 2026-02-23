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

package cli

import (
	"context"
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/effect"
	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
	C "github.com/urfave/cli/v3"
)

func TestToAction_Success(t *testing.T) {
	t.Run("converts successful Effect to action", func(t *testing.T) {
		// Arrange
		effect := func(cmd *C.Command) E.Thunk[F.Void] {
			return func(ctx context.Context) E.IOResult[F.Void] {
				return func() R.Result[F.Void] {
					return R.Of(F.Void{})
				}
			}
		}
		action := ToAction(effect)
		cmd := &C.Command{Name: "test"}

		// Act
		err := action(context.Background(), cmd)

		// Assert
		assert.NoError(t, err)
	})
}

func TestToAction_Failure(t *testing.T) {
	t.Run("converts failed Effect to error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("test error")
		effect := func(cmd *C.Command) E.Thunk[F.Void] {
			return func(ctx context.Context) E.IOResult[F.Void] {
				return func() R.Result[F.Void] {
					return R.Left[F.Void](expectedErr)
				}
			}
		}
		action := ToAction(effect)
		cmd := &C.Command{Name: "test"}

		// Act
		err := action(context.Background(), cmd)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestFromAction_Success(t *testing.T) {
	t.Run("converts successful action to Effect", func(t *testing.T) {
		// Arrange
		action := func(ctx context.Context, cmd *C.Command) error {
			return nil
		}
		effect := FromAction(action)
		cmd := &C.Command{Name: "test"}

		// Act
		result := effect(cmd)(context.Background())()

		// Assert
		assert.True(t, R.IsRight(result))
	})
}

func TestFromAction_Failure(t *testing.T) {
	t.Run("converts failed action to Effect", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("test error")
		action := func(ctx context.Context, cmd *C.Command) error {
			return expectedErr
		}
		effect := FromAction(action)
		cmd := &C.Command{Name: "test"}

		// Act
		result := effect(cmd)(context.Background())()

		// Assert
		assert.True(t, R.IsLeft(result))
		err := R.MonadFold(result, F.Identity[error], func(F.Void) error { return nil })
		assert.Equal(t, expectedErr, err)
	})
}

func TestMakeCommand(t *testing.T) {
	t.Run("creates command with Effect-based action", func(t *testing.T) {
		// Arrange
		effect := func(cmd *C.Command) E.Thunk[F.Void] {
			return func(ctx context.Context) E.IOResult[F.Void] {
				return func() R.Result[F.Void] {
					return R.Of(F.Void{})
				}
			}
		}

		// Act
		cmd := MakeCommand(
			"test",
			"Test command",
			[]C.Flag{},
			effect,
		)

		// Assert
		assert.NotNil(t, cmd)
		assert.Equal(t, "test", cmd.Name)
		assert.Equal(t, "Test command", cmd.Usage)
		assert.NotNil(t, cmd.Action)

		// Test the action
		err := cmd.Action(context.Background(), cmd)
		assert.NoError(t, err)
	})
}

func TestMakeCommandWithSubcommands(t *testing.T) {
	t.Run("creates command with subcommands and Effect-based action", func(t *testing.T) {
		// Arrange
		subCmd := &C.Command{Name: "sub"}
		effect := func(cmd *C.Command) E.Thunk[F.Void] {
			return func(ctx context.Context) E.IOResult[F.Void] {
				return func() R.Result[F.Void] {
					return R.Of(F.Void{})
				}
			}
		}

		// Act
		cmd := MakeCommandWithSubcommands(
			"parent",
			"Parent command",
			[]C.Flag{},
			[]*C.Command{subCmd},
			effect,
		)

		// Assert
		assert.NotNil(t, cmd)
		assert.Equal(t, "parent", cmd.Name)
		assert.Equal(t, "Parent command", cmd.Usage)
		assert.Len(t, cmd.Commands, 1)
		assert.Equal(t, "sub", cmd.Commands[0].Name)
		assert.NotNil(t, cmd.Action)
	})
}

func TestToAction_Integration(t *testing.T) {
	t.Run("Effect can access command flags", func(t *testing.T) {
		// Arrange
		var capturedValue string
		effect := func(cmd *C.Command) E.Thunk[F.Void] {
			return func(ctx context.Context) E.IOResult[F.Void] {
				return func() R.Result[F.Void] {
					capturedValue = cmd.String("input")
					return R.Of(F.Void{})
				}
			}
		}

		cmd := &C.Command{
			Name: "test",
			Flags: []C.Flag{
				&C.StringFlag{
					Name:  "input",
					Value: "default-value",
				},
			},
			Action: ToAction(effect),
		}

		// Act
		err := cmd.Action(context.Background(), cmd)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "default-value", capturedValue)
	})
}

// Made with Bob
