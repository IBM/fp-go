package traversable

import (
	"context"
	"fmt"
	"strconv"

	A "github.com/IBM/fp-go/v2/array"
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/effect"
	F "github.com/IBM/fp-go/v2/function"
	trav "github.com/IBM/fp-go/v2/internal/traversable"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Thunk[T any]  = effect.Thunk[T]
	Option[T any] = O.Option[T]
)

func ExampleTraversable() {

	// Traversable[string, Thunk[int], Option[string], Thunk[Option[int]]
	travOption := O.MakeTraversable[string, int](
		thunk.Of,
		thunk.Map,
	)

	// Traversable[string, Thunk[int], Option[string], Thunk[Option[int]]]
	travArray := A.MakeTraversable[Option[string], Option[int]](
		thunk.Of,
		thunk.Map,
		thunk.Ap,
	)

	// Traversable[string, Thunk[int], []Option[string], Thunk[[]Option[int]]
	travArrayOption := trav.ComposeTraversables(
		travArray,
		travOption,
	)

	// func(string) Thunk[int]
	f := F.Flow2(
		result.Eitherize1(strconv.Atoi),
		thunk.FromResult,
	)

	res := F.Pipe1(
		A.From(O.Of("1"), O.Of("2")),
		travArrayOption(f),
	)

	fmt.Println(res(context.Background())())

	// Output:
	// Right[[]option.Option[int]]([Some[int](1) Some[int](2)])
}
