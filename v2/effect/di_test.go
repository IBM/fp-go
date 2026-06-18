package effect

import (
	"context"
	"errors"
	"fmt"

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/record"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

func Example_di() {

	type Dependency1 struct {
		Value int
	}

	type Dependency2 struct {
		Value int
	}

	type AllDependencies struct {
		Dep1 Dependency1
		Dep2 Dependency2
	}

	type DependencyID string

	type Injector = Effect[DependencyID, any]
	type Provider = Effect[Injector, any]

	const ID1 DependencyID = "dep1"
	const ID2 DependencyID = "dep2"
	const ID3 DependencyID = "all"

	// hardcoding just for the example, normally use code generation
	value1Lens := lens.MakeLens(func(d Dependency1) int { return d.Value }, func(d Dependency1, v int) Dependency1 { d.Value = v; return d })
	value2Lens := lens.MakeLens(func(d Dependency2) int { return d.Value }, func(d Dependency2, v int) Dependency2 { d.Value = v; return d })

	dep1Lens := lens.MakeLens(func(d AllDependencies) Dependency1 { return d.Dep1 }, func(d AllDependencies, v Dependency1) AllDependencies { d.Dep1 = v; return d })
	dep2Lens := lens.MakeLens(func(d AllDependencies) Dependency2 { return d.Dep2 }, func(d AllDependencies, v Dependency2) AllDependencies { d.Dep2 = v; return d })

	dep1Prism := prism.MakePrism(func(d Dependency1) Option[int] { return option.FromNonZero[int]()(d.Value) }, func(s int) Dependency1 { return Dependency1{Value: s} })

	// dependency1 has no other dependencies
	providerDep1 := F.Pipe2(
		io.IntN(100),
		io.Map(dep1Prism.ReverseGet),
		FromIO[Injector],
	)

	// dep2 depends on dep1
	readDep1 := F.Pipe2(
		ID1,
		Read[any],
		ChainResultK[Injector](result.InstanceOf[Dependency1]),
	)

	providerDep2 := F.Pipe1(
		Do[Injector](Dependency2{}),
		ApSL(value2Lens, F.Flow2(
			readDep1,
			thunk.Map(F.Flow2(
				value1Lens.Get,
				N.Mul(2),
			)),
		)),
	)

	// final dependency depends on dep1 and dep2
	readDep2 := F.Pipe2(
		ID2,
		Read[any],
		ChainResultK[Injector](result.InstanceOf[Dependency2]),
	)

	providerFinal := F.Pipe2(
		Do[Injector](AllDependencies{}),
		ApSL(dep1Lens, readDep1),
		ApSL(dep2Lens, readDep2),
	)

	// assemble all providers in a map
	providers := map[DependencyID]Provider{
		ID1: F.Flow2(providerDep1, thunk.Map(F.ToAny[Dependency1])),
		ID2: F.Flow2(providerDep2, thunk.Map(F.ToAny[Dependency2])),
		ID3: F.Flow2(providerFinal, thunk.Map(F.ToAny[AllDependencies])),
	}

	// can we use the y-combinator here?
	injector := func(deps map[DependencyID]Provider) Injector {

		lookup := F.Bind1st(record.MonadLookup, providers)

		var inj Injector

		inj = Memoize(func(id DependencyID) Thunk[any] {
			return F.Pipe4(
				id,
				lookup,
				thunk.FromOption[Provider](lazy.Of(errors.New("dependency not found"))),
				thunk.Flap[Thunk[any]](inj),
				thunk.Flatten,
			)
		})

		return inj

	}(providers)

	allDeps := F.Pipe1(
		injector(ID3),
		thunk.ChainResultK(result.InstanceOf[AllDependencies]),
	)

	res, _ := result.UnwrapError(allDeps(context.Background())())

	fmt.Println(res.Dep1.Value*2 == res.Dep2.Value)

	// Output: true

}
