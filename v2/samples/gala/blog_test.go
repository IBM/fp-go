package gala

import (
	"fmt"
	"strconv"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

func ExampleOption_chaining() {

	result := F.Pipe3(
		O.Some(42),
		O.Map(N.Mul(2)),
		O.Filter(N.MoreThan(50)),
		O.GetOrElse(L.Of(0)),
	)

	fmt.Println(result)

	// Output: 84
}

func ExampleOption_pointfree() {

	endo := F.Flow4(
		O.Some,
		O.Map(N.Mul(2)),
		O.Filter(N.MoreThan(50)),
		O.GetOrElse(L.Of(0)),
	)

	fmt.Println(endo(42))
	fmt.Println(endo(2))

	// Output:
	// 84
	// 0
}

func ExampleEither_chaining() {

	result := F.Pipe2(
		E.Of[string](25),
		E.Filter(N.MoreThan(0), "age cannot be negative"),
		E.Fold(
			S.Format[string]("Error: %s"),
			S.Format[int]("Valid age: %d"),
		),
	)

	fmt.Println(result)

	// Output:
	// Valid age: 25
}

func ExampleEither_trycatch() {

	conv := F.Flow3(
		R.Eitherize1(strconv.Atoi),
		R.Map(N.Mul(2)),
		R.GetOrElse(F.Constant1[error](-1)),
	)

	fmt.Println(conv("42"))
	fmt.Println(conv("abc"))

	// Output:
	// 84
	// -1
}

func isEven(i int) bool {
	return i%2 == 0
}

func sqr(x int) int {
	return x * x
}

func Example_collections() {

	nums := A.From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	sqrEven := F.Flow3(
		A.Filter(isEven),
		A.Map(sqr),
		A.Fold(N.MonoidSum[int]()),
	)

	fmt.Println(sqrEven(nums))

	// Output:
	// 220

}
