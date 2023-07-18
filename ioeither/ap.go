package ioeither

import (
	G "github.com/IBM/fp-go/ioeither/generic"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[E, A, B any](first IOEither[E, A], second IOEither[E, B]) IOEither[E, A] {
	return G.MonadApFirst[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) A]](first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[E, A, B any](second IOEither[E, B]) func(IOEither[E, A]) IOEither[E, A] {
	return G.ApFirst[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) A]](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[E, A, B any](first IOEither[E, A], second IOEither[E, B]) IOEither[E, B] {
	return G.MonadApSecond[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) B]](first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[E, A, B any](second IOEither[E, B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.ApSecond[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) B]](second)
}
