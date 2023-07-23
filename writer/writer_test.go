package writer

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/monoid"
	T "github.com/IBM/fp-go/tuple"
)

func doubleAndLog(data int) Writer[[]string, int] {
	return func() T.Tuple2[int, []string] {
		result := data * 2
		return T.MakeTuple2(result, A.Of(fmt.Sprintf("Doubled %d -> %d", data, result)))
	}
}

func ExampleLoggingWriter() {

	m := A.Monoid[string]()
	s := M.ToSemigroup(m)

	res := F.Pipe3(
		10,
		Of[int](m),
		Chain[int, int](s)(doubleAndLog),
		Chain[int, int](s)(doubleAndLog),
	)

	fmt.Println(res())

	// Output: test

}
