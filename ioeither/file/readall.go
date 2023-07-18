package file

import (
	"io"

	IOE "github.com/ibm/fp-go/ioeither"
)

func onReadAll[R io.Reader](r R) IOE.IOEither[error, []byte] {
	return IOE.TryCatchError(func() ([]byte, error) {
		return io.ReadAll(r)
	})
}

// ReadAll uses a generator function to create a stream, reads it and closes it
func ReadAll[R io.ReadCloser](acquire IOE.IOEither[error, R]) IOE.IOEither[error, []byte] {
	return IOE.WithResource[error, R, []byte](
		acquire,
		onClose[R])(
		onReadAll[R],
	)
}
