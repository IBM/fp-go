package file

import (
	"os"

	IOE "github.com/IBM/fp-go/ioeither"
)

var (
	// Open opens a file for reading
	Open = IOE.Eitherize1(os.Open)
	// ReadFile reads the context of a file
	ReadFile = IOE.Eitherize1(os.ReadFile)
	// WriteFile writes a data blob to a file
	WriteFile = func(dstName string, perm os.FileMode) func([]byte) IOE.IOEither[error, []byte] {
		return func(data []byte) IOE.IOEither[error, []byte] {
			return IOE.TryCatchError(func() ([]byte, error) {
				return data, os.WriteFile(dstName, data, perm)
			})
		}
	}
)
