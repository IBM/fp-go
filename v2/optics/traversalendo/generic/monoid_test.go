package generic

import (
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/effect"
	F "github.com/IBM/fp-go/v2/function"
	AR "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/monoid"
	TLG "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
	TTG "github.com/IBM/fp-go/v2/optics/traversalendo/traversable/generic"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type (
	Thunk[A any]                = effect.Thunk[A]
	ThunkTraversable[A, GA any] = TTG.Traversable[A, Thunk[A], GA, Thunk[GA]]
	ThunkTraversal[S, A any]    = Traversal[S, A, Thunk[Endomorphism[S]], Thunk[A]]
)

// fp-go:Lens
type Address struct {
	Street   string
	Name     string
	Keywords []string
}

func transform(s string) Thunk[string] {
	return thunk.Of(strings.ToUpper(s))
}

func fromTraversableLens[S, A, GA any](trv ThunkTraversable[A, GA]) func(Lens[S, GA]) ThunkTraversal[S, A] {
	return TTG.FromTraversableLens[A, Thunk[A], S, GA](thunk.Map)(trv)
}

func TestMonoid(t *testing.T) {

	m := MakeMonoid[string, Thunk[string], *Address](
		thunk.Of,
		thunk.Map,
		thunk.Ap,
	)

	arTraversable := AR.Traversable[[]string, []string](
		thunk.Of,
		thunk.Map,
		thunk.Ap,
	)

	fld := monoid.Fold(m)

	fromArray := fromTraversableLens[*Address](arTraversable)

	lenses := MakeAddressRefLenses()

	fromString := TLG.FromLens[*Address, string](
		thunk.Map,
	)

	streetTrav := fromString(lenses.Street)
	nameTrav := fromString(lenses.Name)
	keyTrav := fromArray(lenses.Keywords)

	addrTravEndo := F.Pipe1(
		A.From(streetTrav, nameTrav, keyTrav),
		fld,
	)

	// change the strings to upper
	addrTrav := ToTraversal[string, Thunk[string], *Address](thunk.Map)(addrTravEndo)

	trfrmS := addrTrav(transform)

	addr := Address{Street: "street", Name: "name", Keywords: A.From("a", "b")}

	res := trfrmS(&addr)

	newAddr := res(t.Context())()

	expAddr := Address{Street: "STREET", Name: "NAME", Keywords: A.From("A", "B")}

	assert.Equal(t, result.Of(&expAddr), newAddr)
}
