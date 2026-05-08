package option

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readeroption"
)

// Modify transforms the focused value of an option lens when the focus is present.
//
// If the lens resolves to Some(a), the function is applied and the updated value is written back.
// If the lens resolves to None, the original structure is returned unchanged.
//
// Example:
//
//	nameLens := FromNillableRef(L.MakeLensRef((*Street).GetName, (*Street).SetName))
//	updated := Modify[*Street](func(name string) string { return name + "!" })(nameLens)(street)
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(LensO[S, A]) Endomorphism[S] {

	orElse := readeroption.GetOrElse(reader.Ask[S]())

	return func(la LensO[S, A]) Endomorphism[S] {

		return F.Pipe2(
			la.Get,
			readeroption.ChainReaderK(F.Flow3(
				f,
				option.Of,
				la.Set,
			)),
			orElse,
		)

	}
}

// Set returns a function that updates the focus of an option lens to Some(a).
//
// This is a convenience helper for assigning a present value through a LensO without
// constructing the option manually.
//
// Example:
//
//	nameLens := FromNillableRef(L.MakeLensRef((*Street).GetName, (*Street).SetName))
//	updated := Set[*Street]("Main")(nameLens)(street)
func Set[S any, A any](a A) func(LensO[S, A]) Endomorphism[S] {

	oa := option.Of(a)

	return func(la LensO[S, A]) Endomorphism[S] {
		return la.Set(oa)
	}
}
