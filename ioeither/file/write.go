package file

import (
	"io"

	IOE "github.com/IBM/fp-go/ioeither"
)

func onWriteAll[W io.Writer](data []byte) func(w W) IOE.IOEither[error, []byte] {
	return func(w W) IOE.IOEither[error, []byte] {
		return IOE.TryCatchError(func() ([]byte, error) {
			_, err := w.Write(data)
			return data, err
		})
	}
}

// WriteAll uses a generator function to create a stream, writes data to it and closes it
func WriteAll[W io.WriteCloser](data []byte) func(acquire IOE.IOEither[error, W]) IOE.IOEither[error, []byte] {
	onWrite := onWriteAll[W](data)
	return func(onCreate IOE.IOEither[error, W]) IOE.IOEither[error, []byte] {
		return IOE.WithResource[error, W, []byte](
			onCreate,
			onClose[W])(
			onWrite,
		)
	}
}

// Write uses a generator function to create a stream, writes data to it and closes it
func Write[W io.WriteCloser, R any](acquire IOE.IOEither[error, W]) func(use func(W) IOE.IOEither[error, R]) IOE.IOEither[error, R] {
	return IOE.WithResource[error, W, R](
		acquire,
		onClose[W])
}
