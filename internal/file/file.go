package file

import (
	"bytes"
	"context"
	"io"

	E "github.com/ibm/fp-go/either"
)

type (
	readerWithContext struct {
		ctx      context.Context
		delegate io.Reader
	}
)

func (rdr *readerWithContext) Read(p []byte) (int, error) {
	// check for cancellarion
	if err := rdr.ctx.Err(); err != nil {
		return 0, err
	}
	// simply dispatch
	return rdr.delegate.Read(p)
}

// MakeReader creates a context aware reader
func MakeReader(ctx context.Context, rdr io.Reader) io.Reader {
	return &readerWithContext{ctx, rdr}
}

// ReadAll reads the content of a reader and allows it to be canceled
func ReadAll(ctx context.Context, rdr io.Reader) E.Either[error, []byte] {
	return E.TryCatchError(func() ([]byte, error) {
		var buffer bytes.Buffer
		_, err := io.Copy(&buffer, MakeReader(ctx, rdr))
		return buffer.Bytes(), err
	})
}
