package reader

import (
	"context"
	"fmt"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/function"
	T "github.com/IBM/fp-go/tuple"
	"github.com/stretchr/testify/assert"
)

func GoFunction(ctx context.Context, data string) string {
	return strings.ToUpper(data)
}

func GoIntFunction(ctx context.Context, data string, number int) string {
	return fmt.Sprintf("%s: %d", data, number)
}

func TestReaderFrom(t *testing.T) {
	ctx := context.Background()
	f := From1(GoFunction)

	result := f("input")(ctx)

	assert.Equal(t, result, "INPUT")

}

func MyFinalResult(left, right string) string {
	return fmt.Sprintf("%s-%s", left, right)
}

func TestReadersFrom(t *testing.T) {
	ctx := context.Background()

	f1 := From1(GoFunction)
	f2 := From2(GoIntFunction)

	result1 := f1("input")(ctx)
	result2 := f2("input", 10)(ctx)

	result3 := MyFinalResult(result1, result2)

	h := F.Pipe1(
		SequenceT2(f1("input"), f2("input", 10)),
		Map(T.Tupled2(MyFinalResult)),
	)

	composedResult := h(ctx)

	assert.Equal(t, result1, "INPUT")
	assert.Equal(t, result2, "input: 10")
	assert.Equal(t, result3, "INPUT-input: 10")

	assert.Equal(t, composedResult, "INPUT-input: 10")

}
