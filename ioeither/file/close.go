package file

import (
	"io"

	IOE "github.com/IBM/fp-go/ioeither"
)

func onClose[R io.Closer](r R) IOE.IOEither[error, R] {
	return IOE.TryCatchError(func() (R, error) {
		return r, r.Close()
	})
}
