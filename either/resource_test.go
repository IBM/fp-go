package either

import (
	"os"
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestWithResource(t *testing.T) {
	onCreate := func() Either[error, *os.File] {
		return TryCatchError(func() (*os.File, error) {
			return os.CreateTemp("", "*")
		})
	}
	onDelete := F.Flow2(
		func(f *os.File) Either[error, string] {
			return TryCatchError(func() (string, error) {
				return f.Name(), f.Close()
			})
		},
		Chain(func(name string) Either[error, any] {
			return TryCatchError(func() (any, error) {
				return name, os.Remove(name)
			})
		}),
	)

	onHandler := func(f *os.File) Either[error, string] {
		return Of[error](f.Name())
	}

	tempFile := WithResource[error, *os.File, string](onCreate, onDelete)

	resE := tempFile(onHandler)

	assert.True(t, IsRight(resE))
}
