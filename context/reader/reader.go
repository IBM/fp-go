package reader

import (
	"context"

	R "github.com/IBM/fp-go/reader/generic"
)

func MonadMap[A, B any](fa Reader[A], f func(A) B) Reader[B] {
	return R.MonadMap[Reader[A], Reader[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(Reader[A]) Reader[B] {
	return R.Map[Reader[A], Reader[B]](f)
}

func MonadChain[A, B any](ma Reader[A], f func(A) Reader[B]) Reader[B] {
	return R.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) Reader[B]) func(Reader[A]) Reader[B] {
	return R.Chain[Reader[A]](f)
}

func Of[A any](a A) Reader[A] {
	return R.Of[Reader[A]](a)
}

func MonadAp[A, B any](fab Reader[func(A) B], fa Reader[A]) Reader[B] {
	return R.MonadAp[Reader[A], Reader[B]](fab, fa)
}

func Ap[A, B any](fa Reader[A]) func(Reader[func(A) B]) Reader[B] {
	return R.Ap[Reader[A], Reader[B], Reader[func(A) B]](fa)
}

func Ask() Reader[context.Context] {
	return R.Ask[Reader[context.Context]]()
}

func Asks[A any](r Reader[A]) Reader[A] {
	return R.Asks(r)
}
