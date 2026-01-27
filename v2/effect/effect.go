package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
)

func Succeed[C, A any](a A) Effect[C, A] {
	return readerreaderioresult.Of[C](a)
}

func Fail[C, A any](err error) Effect[C, A] {
	return readerreaderioresult.Left[C, A](err)
}

func Of[C, A any](a A) Effect[C, A] {
	return readerreaderioresult.Of[C](a)
}

func Map[C, A, B any](f func(A) B) Operator[C, A, B] {
	return readerreaderioresult.Map[C](f)
}

func Chain[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, B] {
	return readerreaderioresult.Chain(f)
}

func Ap[B, C, A any](fa Effect[C, A]) Operator[C, func(A) B, B] {
	return readerreaderioresult.Ap[B](fa)
}

func Suspend[C, A any](fa Lazy[Effect[C, A]]) Effect[C, A] {
	return readerreaderioresult.Defer(fa)
}

func Tap[C, A, ANY any](f Kleisli[C, A, ANY]) Operator[C, A, A] {
	return readerreaderioresult.Tap(f)
}

func Ternary[C, A, B any](pred Predicate[A], onTrue, onFalse Kleisli[C, A, B]) Kleisli[C, A, B] {
	return function.Ternary(pred, onTrue, onFalse)
}

func ChainResultK[C, A, B any](f result.Kleisli[A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainResultK[C](f)
}

func Read[A, C any](c C) func(Effect[C, A]) Thunk[A] {
	return readerreaderioresult.Read[A](c)
}
