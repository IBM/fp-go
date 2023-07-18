package file

import (
	"context"
	"io"
	"os"

	RIOE "github.com/IBM/fp-go/context/readerioeither"
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/file"
	IOE "github.com/IBM/fp-go/ioeither"
)

var (
	openIOE = IOE.Eitherize1(os.Open)
	// Open opens a file for reading within the given context
	Open = F.Flow3(
		openIOE,
		RIOE.FromIOEither[*os.File],
		RIOE.WithContext[*os.File],
	)
)

// Close closes an object
func Close[C io.Closer](c C) RIOE.ReaderIOEither[any] {
	return RIOE.FromIOEither(func() ET.Either[error, any] {
		return ET.TryCatchError(func() (any, error) {
			return c, c.Close()
		})
	})
}

// ReadFile reads a file in the scope of a context
func ReadFile(path string) RIOE.ReaderIOEither[[]byte] {
	return RIOE.WithResource[*os.File, []byte](Open(path), Close[*os.File])(func(r *os.File) RIOE.ReaderIOEither[[]byte] {
		return func(ctx context.Context) IOE.IOEither[error, []byte] {
			return IOE.MakeIO(func() ET.Either[error, []byte] {
				return file.ReadAll(ctx, r)
			})
		}
	})
}
