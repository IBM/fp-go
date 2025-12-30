package readerioresult

import (
	"crypto/sha256"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	FL "github.com/IBM/fp-go/v2/ioeither/file"
	P "github.com/IBM/fp-go/v2/predicate"
)

func TestChecksum(t *testing.T) {

	expectedChecksum := A.From[byte](246, 64, 183, 1, 47, 191, 7, 44, 139, 70, 69, 165, 89, 118, 27, 12, 249, 255, 186, 152, 68, 252, 168, 221, 157, 39, 156, 242, 31, 29, 65, 164)

	verifyHash := F.Pipe1(
		FL.ReadFile,
		Map[string](F.Flow3(
			sha256.Sum256,
			func(s [sha256.Size]byte) []byte { return s[:] },
			P.IsEqual(A.StrictEquals[byte]())(expectedChecksum),
		)),
	)

	// removeFile := FL.Remove

	fmt.Println(verifyHash("data/sample.txt")())
}
