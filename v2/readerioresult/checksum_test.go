package readerioresult

import (
	"crypto/sha256"
	"fmt"
	"path"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	FL "github.com/IBM/fp-go/v2/ioeither/file"
	"github.com/IBM/fp-go/v2/ioresult"
	P "github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/result"
)

func TestChecksum(t *testing.T) {

	tmpDir := t.TempDir()

	expectedChecksum := [sha256.Size]byte{246, 64, 183, 1, 47, 191, 7, 44, 139, 70, 69, 165, 89, 118, 27, 12, 249, 255, 186, 152, 68, 252, 168, 221, 157, 39, 156, 242, 31, 29, 65, 164}

	testFile := FL.CopyFile("data/sample.txt")(path.Join(tmpDir, "test.txt"))

	verifyHash := F.Pipe1(
		FL.ReadFile,
		ChainFirstResultK[string](F.Flow2(
			sha256.Sum256,
			result.FromPredicate(P.IsStrictEqual[[sha256.Size]byte]()(expectedChecksum), errors.OnSome[[sha256.Size]byte]("Invalid checksum")),
		)),
	)

	veryfiyAndRemove := F.Pipe1(
		verifyHash,
		OrElse(F.Constant1[error](F.Flow2(
			FL.Remove,
			ioresult.MapTo[string](A.Empty[byte]()),
		))),
	)

	res := F.Pipe2(
		testFile,
		ioresult.Chain(veryfiyAndRemove),
		ioresult.Map(A.Size[byte]),
	)

	fmt.Println(res())
}
