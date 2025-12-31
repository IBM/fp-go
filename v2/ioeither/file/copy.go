package file

import (
	"io"
	"os"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
)

//go:inline
func CopyFile(src, dst string) IOEither[error, string] {
	withSrc := ioeither.WithResource[int64](Open(src), Close)
	withDst := ioeither.WithResource[int64](Create(dst), Close)

	return F.Pipe1(
		withSrc(func(srcFile *os.File) IOEither[error, int64] {
			return withDst(func(dstFile *os.File) IOEither[error, int64] {
				return func() Either[error, int64] {
					return either.TryCatchError(io.Copy(dstFile, srcFile))
				}
			})
		}),
		ioeither.MapTo[error, int64](dst),
	)
}
