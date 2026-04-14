# fp-go/v2 Complete API Reference

> Module: `github.com/IBM/fp-go/v2`
>
> Authoritative reference of every exported type and function.
> Generated from source code. Intended for LLM consumption during code generation.
>
> **Conventions:**
> - Curried form (e.g., `Map(f)(fa)`) is the primary API
> - `Monad*` prefixed = uncurried variants (take all args at once)  
> - `Operator[A, B]` = `func(Type[A]) Type[B]` (composable pipeline stage)
> - `Kleisli[A, B]` = `func(A) Type[B]` (effectful function)
> - `*G` suffix = generic version for custom slice/map types
> - `SequenceTN` / `TraverseTupleN` = N-ary tuple operations (1..15)
> - `EitherizeN` / `UneitherizeN` = convert (R, error) functions (0..15)
> - `OptimizeN` / `UnoptionizeN` = convert (R, bool) functions (0..10)
> - Do/Bind/Let/ApS = do-notation simulation for building structs

---

# Table of Contents

1. [Core Monads](#core-monads): option, either, result, io, iooption, ioeither, ioresult
2. [Reader Stack](#reader-stack): reader, readeroption, readereither, readerio, readeriooption, readerioeither, readerioresult, readerresult
3. [Context Specializations](#context-specializations): context/readerioresult, context/readerresult
4. [Effect System](#effect-system): effect
5. [State Monads](#state-monads): state, stateio, statereaderioeither
6. [Optics](#optics): lens, prism, iso, optional, traversal, codec
7. [Utilities](#utilities): function, array, record, pair, tuple, predicate, endomorphism
8. [Algebraic Structures](#algebraic-structures): eq, ord, semigroup, monoid
9. [Primitives](#primitives): number, string, boolean, bytes
10. [Other](#other): identity, lazy, constant, json, di, builder, retry, circuitbreaker, tailrec, ioref, consumer, erasure, iterator/stateless
11. [Idiomatic](#idiomatic): idiomatic/option, idiomatic/result, idiomatic/ioresult, idiomatic/readerresult, idiomatic/readerioresult, idiomatic/context/readerresult

---


---

# Core Monads

## package `github.com/IBM/fp-go/v2/option`

Import: `import O "github.com/IBM/fp-go/v2/option"`

Option represents an optional value: Some(value) or None. Type-safe alternative to nil pointers.

Key types:
- `Option[A]` -- the core type, a struct wrapping a value or empty
- `Kleisli[A, B] = func(A) Option[B]` -- effectful function returning Option
- `Operator[A, B] = Kleisli[Option[A], B]` -- composable pipeline operator

### Exported API

```go
func AltMonoid[A any]() M.Monoid[Option[A]]
func AlternativeMonoid[A any](m M.Monoid[A]) M.Monoid[Option[A]]
func ApplicativeMonoid[A any](m M.Monoid[A]) M.Monoid[Option[A]]
func ApplySemigroup[A any](s S.Semigroup[A]) S.Semigroup[Option[A]]
func CompactArray[A any](fa []Option[A]) []A
func CompactArrayG[A1 ~[]Option[A], A2 ~[]A, A any](fa A1) A2
func CompactRecord[K comparable, A any](m map[K]Option[A]) map[K]A
func CompactRecordG[M1 ~map[K]Option[A], M2 ~map[K]A, K comparable, A any](m M1) M2
func Eq[A any](a EQ.Eq[A]) EQ.Eq[Option[A]]
func FirstMonoid[A any]() M.Monoid[Option[A]]
func Fold[A, B any](onNone func() B, onSome func(a A) B) func(ma Option[A]) B
func FromEq[A any](pred eq.Eq[A]) func(A) Kleisli[A, A]
func FromStrictCompare[A C.Ordered]() ord.Ord[Option[A]]
func FromStrictEq[A comparable]() func(A) Kleisli[A, A]
func FromStrictEquals[A comparable]() EQ.Eq[Option[A]]
func Functor[A, B any]() functor.Functor[A, B, Option[A], Option[B]]
func GetOrElse[A any](onNone func() A) func(Option[A]) A
func IsNone[T any](val Option[T]) bool
func IsSome[T any](val Option[T]) bool
func LastMonoid[A any]() M.Monoid[Option[A]]
func Logger[A any](loggers ...*log.Logger) func(string) Kleisli[Option[A], A]
func Monad[A, B any]() monad.Monad[A, B, Option[A], Option[B], Option[func(A) B]]
func MonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B
func MonadGetOrElse[A any](fa Option[A], onNone func() A) A
func Monoid[A any]() func(S.Semigroup[A]) M.Monoid[Option[A]]
func Optionize0[F ~func() (R, bool), R any](f F) func() Option[R]
func Optionize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, bool), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) Option[R]
func Optionize2[F ~func(T0, T1) (R, bool), T0, T1, R any](f F) func(T0, T1) Option[R]
func Optionize3[F ~func(T0, T1, T2) (R, bool), T0, T1, T2, R any](f F) func(T0, T1, T2) Option[R]
func Optionize4[F ~func(T0, T1, T2, T3) (R, bool), T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) Option[R]
func Optionize5[F ~func(T0, T1, T2, T3, T4) (R, bool), T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) Option[R]
func Optionize6[F ~func(T0, T1, T2, T3, T4, T5) (R, bool), T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) Option[R]
func Optionize7[F ~func(T0, T1, T2, T3, T4, T5, T6) (R, bool), T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) Option[R]
func Optionize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) (R, bool), T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) Option[R]
func Optionize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, bool), T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) Option[R]
func Ord[A any](a ord.Ord[A]) ord.Ord[Option[A]]
func Pointed[A any]() pointed.Pointed[A, Option[A]]
func Reduce[A, B any](f func(B, A) B, initial B) func(Option[A]) B
func Semigroup[A any]() func(S.Semigroup[A]) S.Semigroup[Option[A]]
func Sequence[A, HKTA, HKTOA any](
func Sequence2[T1, T2, R any](f func(T1, T2) Option[R]) func(Option[T1], Option[T2]) Option[R]
func Traverse[A, B, HKTB, HKTOB any](
func TraverseTuple1[F1 ~Kleisli[A1, T1], A1, T1 any](f1 F1) func(T.Tuple1[A1]) Option[T.Tuple1[T1]]
func TraverseTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(T.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) Option[T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], A1, T1, A2, T2 any](f1 F1, f2 F2) func(T.Tuple2[A1, A2]) Option[T.Tuple2[T1, T2]]
func TraverseTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], A1, T1, A2, T2, A3, T3 any](f1 F1, f2 F2, f3 F3) func(T.Tuple3[A1, A2, A3]) Option[T.Tuple3[T1, T2, T3]]
func TraverseTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], A1, T1, A2, T2, A3, T3, A4, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T.Tuple4[A1, A2, A3, A4]) Option[T.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T.Tuple5[A1, A2, A3, A4, A5]) Option[T.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(T.Tuple6[A1, A2, A3, A4, A5, A6]) Option[T.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(T.Tuple7[A1, A2, A3, A4, A5, A6, A7]) Option[T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(T.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) Option[T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(T.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) Option[T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Unoptionize0[F ~func() Option[R], R any](f F) func() (R, bool)
func Unoptionize1[F ~Kleisli[T0, R], T0, R any](f F) func(T0) (R, bool)
func Unoptionize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) Option[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, bool)
func Unoptionize2[F ~func(T0, T1) Option[R], T0, T1, R any](f F) func(T0, T1) (R, bool)
func Unoptionize3[F ~func(T0, T1, T2) Option[R], T0, T1, T2, R any](f F) func(T0, T1, T2) (R, bool)
func Unoptionize4[F ~func(T0, T1, T2, T3) Option[R], T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) (R, bool)
func Unoptionize5[F ~func(T0, T1, T2, T3, T4) Option[R], T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) (R, bool)
func Unoptionize6[F ~func(T0, T1, T2, T3, T4, T5) Option[R], T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) (R, bool)
func Unoptionize7[F ~func(T0, T1, T2, T3, T4, T5, T6) Option[R], T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) (R, bool)
func Unoptionize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) Option[R], T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) (R, bool)
func Unoptionize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) Option[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, bool)
func Unwrap[A any](ma Option[A]) (A, bool)
type Endomorphism[T any] = endomorphism.Endomorphism[T]
type Kleisli[A, B any] = func(A) Option[B]
func FromNonZero[A comparable]() Kleisli[A, A]
func FromPredicate[A any](pred func(A) bool) Kleisli[A, A]
func FromValidation[A, B any](f func(A) (B, bool)) Kleisli[A, B]
func FromZero[A comparable]() Kleisli[A, A]
func Optionize1[F ~func(T0) (R, bool), T0, R any](f F) Kleisli[T0, R]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseArrayWithIndex[A, B any](f func(int, A) Option[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) Option[B]) Kleisli[GA, GB]
func TraverseIter[A, B any](f Kleisli[A, B]) Kleisli[Seq[A], Seq[B]]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Option[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) Option[B]) Kleisli[GA, GB]
type Operator[A, B any] = Kleisli[Option[A], B]
func Alt[A any](that func() Option[A]) Operator[A, A]
func Ap[B, A any](fa Option[A]) Operator[func(A) B, B]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func Bind[S1, S2, A any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainTo[A, B any](mb Option[B]) Operator[A, B]
func Filter[A any](pred func(A) bool) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, B any](
func LetL[S, T any](
func LetTo[S1, S2, B any](
func LetToL[S, T any](
func Map[A, B any](f func(a A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
type Option[A any] struct {
func Do[S any](
func Flatten[A any](mma Option[Option[A]]) Option[A]
func FromNillable[A any](a *A) Option[*A]
func InstanceOf[T any](src any) Option[T]
func MonadAlt[A any](fa Option[A], that func() Option[A]) Option[A]
func MonadAp[B, A any](fab Option[func(A) B], fa Option[A]) Option[B]
func MonadChain[A, B any](fa Option[A], f Kleisli[A, B]) Option[B]
func MonadChainFirst[A, B any](ma Option[A], f Kleisli[A, B]) Option[A]
func MonadChainTo[A, B any](ma Option[A], mb Option[B]) Option[B]
func MonadFlap[B, A any](fab Option[func(A) B], a A) Option[B]
func MonadMap[A, B any](fa Option[A], f func(A) B) Option[B]
func MonadMapTo[A, B any](fa Option[A], b B) Option[B]
func MonadSequence2[T1, T2, R any](o1 Option[T1], o2 Option[T2], f func(T1, T2) Option[R]) Option[R]
func None[T any]() Option[T]
func Of[T any](value T) Option[T]
func SequenceArray[A any](ma []Option[A]) Option[[]A]
func SequenceArrayG[GA ~[]A, GOA ~[]Option[A], A any](ma GOA) Option[GA]
func SequenceIter[A any](as Seq[Option[A]]) Option[Seq[A]]
func SequencePair[T1, T2 any](t P.Pair[Option[T1], Option[T2]]) Option[P.Pair[T1, T2]]
func SequenceRecord[K comparable, A any](ma map[K]Option[A]) Option[map[K]A]
func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Option[A], K comparable, A any](ma GOA) Option[GA]
func SequenceT1[T1 any](t1 Option[T1]) Option[T.Tuple1[T1]]
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t1 Option[T1], t2 Option[T2], t3 Option[T3], t4 Option[T4], t5 Option[T5], t6 Option[T6], t7 Option[T7], t8 Option[T8], t9 Option[T9], t10 Option[T10]) Option[T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceT2[T1, T2 any](t1 Option[T1], t2 Option[T2]) Option[T.Tuple2[T1, T2]]
func SequenceT3[T1, T2, T3 any](t1 Option[T1], t2 Option[T2], t3 Option[T3]) Option[T.Tuple3[T1, T2, T3]]
func SequenceT4[T1, T2, T3, T4 any](t1 Option[T1], t2 Option[T2], t3 Option[T3], t4 Option[T4]) Option[T.Tuple4[T1, T2, T3, T4]]
func SequenceT5[T1, T2, T3, T4, T5 any](t1 Option[T1], t2 Option[T2], t3 Option[T3], t4 Option[T4], t5 Option[T5]) Option[T.Tuple5[T1, T2, T3, T4, T5]]
func SequenceT6[T1, T2, T3, T4, T5, T6 any](t1 Option[T1], t2 Option[T2], t3 Option[T3], t4 Option[T4], t5 Option[T5], t6 Option[T6]) Option[T.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](t1 Option[T1], t2 Option[T2], t3 Option[T3], t4 Option[T4], t5 Option[T5], t6 Option[T6], t7 Option[T7]) Option[T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](t1 Option[T1], t2 Option[T2], t3 Option[T3], t4 Option[T4], t5 Option[T5], t6 Option[T6], t7 Option[T7], t8 Option[T8]) Option[T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t1 Option[T1], t2 Option[T2], t3 Option[T3], t4 Option[T4], t5 Option[T5], t6 Option[T6], t7 Option[T7], t8 Option[T8], t9 Option[T9]) Option[T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceTuple1[T1 any](t T.Tuple1[Option[T1]]) Option[T.Tuple1[T1]]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t T.Tuple10[Option[T1], Option[T2], Option[T3], Option[T4], Option[T5], Option[T6], Option[T7], Option[T8], Option[T9], Option[T10]]) Option[T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple2[T1, T2 any](t T.Tuple2[Option[T1], Option[T2]]) Option[T.Tuple2[T1, T2]]
func SequenceTuple3[T1, T2, T3 any](t T.Tuple3[Option[T1], Option[T2], Option[T3]]) Option[T.Tuple3[T1, T2, T3]]
func SequenceTuple4[T1, T2, T3, T4 any](t T.Tuple4[Option[T1], Option[T2], Option[T3], Option[T4]]) Option[T.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t T.Tuple5[Option[T1], Option[T2], Option[T3], Option[T4], Option[T5]]) Option[T.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t T.Tuple6[Option[T1], Option[T2], Option[T3], Option[T4], Option[T5], Option[T6]]) Option[T.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t T.Tuple7[Option[T1], Option[T2], Option[T3], Option[T4], Option[T5], Option[T6], Option[T7]]) Option[T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t T.Tuple8[Option[T1], Option[T2], Option[T3], Option[T4], Option[T5], Option[T6], Option[T7], Option[T8]]) Option[T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t T.Tuple9[Option[T1], Option[T2], Option[T3], Option[T4], Option[T5], Option[T6], Option[T7], Option[T8], Option[T9]]) Option[T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Some[T any](value T) Option[T]
func ToAny[T any](src T) Option[any]
func TryCatch[A any](f func() (A, error)) Option[A]
func Zero[A any]() Option[A]
func (s Option[A]) Format(f fmt.State, c rune)
func (s Option[A]) GoString() string
func (s Option[A]) LogValue() slog.Value
func (s Option[A]) MarshalJSON() ([]byte, error)
func (s Option[A]) String() string
func (s *Option[A]) UnmarshalJSON(data []byte) error
type Seq[T any] = iter.Seq[T]
```

## package `github.com/IBM/fp-go/v2/either`

Import: `import E "github.com/IBM/fp-go/v2/either"`

Either represents a value of one of two types: Left (error) or Right (success).

Key types:
- `Either[E, A]` -- discriminated union: Left[E] or Right[A]
- `Kleisli[E, A, B] = func(A) Either[E, B]` -- effectful function
- `Operator[E, A, B] = Kleisli[E, Either[E, A], B]` -- pipeline operator

### Exported API

```go
func AltSemigroup[E, A any]() S.Semigroup[Either[E, A]]
func ApV[B, A, E any](sg S.Semigroup[E]) func(Either[E, A]) Operator[E, func(A) B, B]
func Applicative[E, A, B any]() applicative.Applicative[A, B, Either[E, A], Either[E, B], Either[E, func(A) B]]
func ApplicativeMonoid[E, A any](m M.Monoid[A]) M.Monoid[Either[E, A]]
func ApplicativeV[E, A, B any](sg S.Semigroup[E]) applicative.Applicative[A, B, Either[E, A], Either[E, B], Either[E, func(A) B]]
func ApplySemigroup[E, A any](s S.Semigroup[A]) S.Semigroup[Either[E, A]]
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(a A) B) func(Either[E1, A]) Either[E2, B]
func ChainOptionK[A, B, E any](onNone func() E) func(func(A) Option[B]) Operator[E, A, B]
func CompactArray[E, A any](fa []Either[E, A]) []A
func CompactArrayG[A1 ~[]Either[E, A], A2 ~[]A, E, A any](fa A1) A2
func CompactRecord[K comparable, E, A any](m map[K]Either[E, A]) map[K]A
func CompactRecordG[M1 ~map[K]Either[E, A], M2 ~map[K]A, K comparable, E, A any](m M1) M2
func Curry0[R any](f func() (R, error)) func() Either[error, R]
func Curry1[T1, R any](f func(T1) (R, error)) func(T1) Either[error, R]
func Curry2[T1, T2, R any](f func(T1, T2) (R, error)) func(T1) func(T2) Either[error, R]
func Curry3[T1, T2, T3, R any](f func(T1, T2, T3) (R, error)) func(T1) func(T2) func(T3) Either[error, R]
func Curry4[T1, T2, T3, T4, R any](f func(T1, T2, T3, T4) (R, error)) func(T1) func(T2) func(T3) func(T4) Either[error, R]
func Eitherize0[F ~func() (R, error), R any](f F) func() Either[error, R]
func Eitherize1[F ~func(T0) (R, error), T0, R any](f F) func(T0) Either[error, R]
func Eitherize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) Either[error, R]
func Eitherize11[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) Either[error, R]
func Eitherize12[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) Either[error, R]
func Eitherize13[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) Either[error, R]
func Eitherize14[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) Either[error, R]
func Eitherize15[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) Either[error, R]
func Eitherize2[F ~func(T0, T1) (R, error), T0, T1, R any](f F) func(T0, T1) Either[error, R]
func Eitherize3[F ~func(T0, T1, T2) (R, error), T0, T1, T2, R any](f F) func(T0, T1, T2) Either[error, R]
func Eitherize4[F ~func(T0, T1, T2, T3) (R, error), T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) Either[error, R]
func Eitherize5[F ~func(T0, T1, T2, T3, T4) (R, error), T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) Either[error, R]
func Eitherize6[F ~func(T0, T1, T2, T3, T4, T5) (R, error), T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) Either[error, R]
func Eitherize7[F ~func(T0, T1, T2, T3, T4, T5, T6) (R, error), T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) Either[error, R]
func Eitherize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) Either[error, R]
func Eitherize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) Either[error, R]
func Eq[E, A any](e EQ.Eq[E], a EQ.Eq[A]) EQ.Eq[Either[E, A]]
func FirstMonoid[E, A any](zero Lazy[Either[E, A]]) M.Monoid[Either[E, A]]
func Fold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B
func FromError[A any](f func(a A) error) func(A) Either[error, A]
func FromOption[A, E any](onNone func() E) func(Option[A]) Either[E, A]
func FromStrictEquals[E, A comparable]() EQ.Eq[Either[E, A]]
func Functor[E, A, B any]() functor.Functor[A, B, Either[E, A], Either[E, B]]
func GetOrElse[E, A any](onLeft func(E) A) func(Either[E, A]) A
func IsLeft[E, A any](val Either[E, A]) bool
func IsRight[E, A any](val Either[E, A]) bool
func LastMonoid[E, A any](zero Lazy[Either[E, A]]) M.Monoid[Either[E, A]]
func Logger[E, A any](loggers ...*log.Logger) func(string) Operator[E, A, A]
func MapLeft[A, E1, E2 any](f func(E1) E2) func(fa Either[E1, A]) Either[E2, A]
func Monad[E, A, B any]() monad.Monad[A, B, Either[E, A], Either[E, B], Either[E, func(A) B]]
func MonadApV[B, A, E any](sg S.Semigroup[E]) func(fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B]
func MonadFold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B
func Partition[E, A any](p Predicate[A], empty E) func(Either[E, A]) Pair[Either[E, A], Either[E, A]]
func PartitionMap[E, A, B, C any](f Kleisli[B, A, C], empty E) func(Either[E, A]) Pair[Either[E, B], Either[E, C]]
func Pointed[E, A any]() pointed.Pointed[A, Either[E, A]]
func Reduce[E, A, B any](f func(B, A) B, initial B) func(Either[E, A]) B
func Sequence[E, A, HKTA, HKTRA any](
func Sequence2[E, T1, T2, R any](f func(T1, T2) Either[E, R]) func(Either[E, T1], Either[E, T2]) Either[E, R]
func Sequence3[E, T1, T2, T3, R any](f func(T1, T2, T3) Either[E, R]) func(Either[E, T1], Either[E, T2], Either[E, T3]) Either[E, R]
func ToError[A any](e Either[error, A]) error
func ToSLogAttr[E, A any]() func(Either[E, A]) slog.Attr
func ToType[A, E any](onError func(any) E) func(any) Either[E, A]
func Traverse[A, E, B, HKTB, HKTRB any](
func TraverseTuple1[F1 ~func(A1) Either[E, T1], E, A1, T1 any](f1 F1) func(T.Tuple1[A1]) Either[E, T.Tuple1[T1]]
func TraverseTuple10[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], F9 ~func(A9) Either[E, T9], F10 ~func(A10) Either[E, T10], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(T.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) Either[E, T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple11[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], F9 ~func(A9) Either[E, T9], F10 ~func(A10) Either[E, T10], F11 ~func(A11) Either[E, T11], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) func(T.Tuple11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11]) Either[E, T.Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func TraverseTuple12[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], F9 ~func(A9) Either[E, T9], F10 ~func(A10) Either[E, T10], F11 ~func(A11) Either[E, T11], F12 ~func(A12) Either[E, T12], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) func(T.Tuple12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12]) Either[E, T.Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func TraverseTuple13[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], F9 ~func(A9) Either[E, T9], F10 ~func(A10) Either[E, T10], F11 ~func(A11) Either[E, T11], F12 ~func(A12) Either[E, T12], F13 ~func(A13) Either[E, T13], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12, A13, T13 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) func(T.Tuple13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13]) Either[E, T.Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func TraverseTuple14[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], F9 ~func(A9) Either[E, T9], F10 ~func(A10) Either[E, T10], F11 ~func(A11) Either[E, T11], F12 ~func(A12) Either[E, T12], F13 ~func(A13) Either[E, T13], F14 ~func(A14) Either[E, T14], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12, A13, T13, A14, T14 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) func(T.Tuple14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14]) Either[E, T.Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func TraverseTuple15[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], F9 ~func(A9) Either[E, T9], F10 ~func(A10) Either[E, T10], F11 ~func(A11) Either[E, T11], F12 ~func(A12) Either[E, T12], F13 ~func(A13) Either[E, T13], F14 ~func(A14) Either[E, T14], F15 ~func(A15) Either[E, T15], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12, A13, T13, A14, T14, A15, T15 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) func(T.Tuple15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15]) Either[E, T.Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func TraverseTuple2[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], E, A1, T1, A2, T2 any](f1 F1, f2 F2) func(T.Tuple2[A1, A2]) Either[E, T.Tuple2[T1, T2]]
func TraverseTuple3[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], E, A1, T1, A2, T2, A3, T3 any](f1 F1, f2 F2, f3 F3) func(T.Tuple3[A1, A2, A3]) Either[E, T.Tuple3[T1, T2, T3]]
func TraverseTuple4[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], E, A1, T1, A2, T2, A3, T3, A4, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T.Tuple4[A1, A2, A3, A4]) Either[E, T.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T.Tuple5[A1, A2, A3, A4, A5]) Either[E, T.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(T.Tuple6[A1, A2, A3, A4, A5, A6]) Either[E, T.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(T.Tuple7[A1, A2, A3, A4, A5, A6, A7]) Either[E, T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(T.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) Either[E, T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[F1 ~func(A1) Either[E, T1], F2 ~func(A2) Either[E, T2], F3 ~func(A3) Either[E, T3], F4 ~func(A4) Either[E, T4], F5 ~func(A5) Either[E, T5], F6 ~func(A6) Either[E, T6], F7 ~func(A7) Either[E, T7], F8 ~func(A8) Either[E, T8], F9 ~func(A9) Either[E, T9], E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(T.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) Either[E, T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Uncurry0[R any](f func() Either[error, R]) func() (R, error)
func Uncurry1[T1, R any](f func(T1) Either[error, R]) func(T1) (R, error)
func Uncurry2[T1, T2, R any](f func(T1) func(T2) Either[error, R]) func(T1, T2) (R, error)
func Uncurry3[T1, T2, T3, R any](f func(T1) func(T2) func(T3) Either[error, R]) func(T1, T2, T3) (R, error)
func Uncurry4[T1, T2, T3, T4, R any](f func(T1) func(T2) func(T3) func(T4) Either[error, R]) func(T1, T2, T3, T4) (R, error)
func Uneitherize0[F ~func() Either[error, R], R any](f F) func() (R, error)
func Uneitherize1[F ~func(T0) Either[error, R], T0, R any](f F) func(T0) (R, error)
func Uneitherize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error)
func Uneitherize11[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error)
func Uneitherize12[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) (R, error)
func Uneitherize13[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) (R, error)
func Uneitherize14[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) (R, error)
func Uneitherize15[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) (R, error)
func Uneitherize2[F ~func(T0, T1) Either[error, R], T0, T1, R any](f F) func(T0, T1) (R, error)
func Uneitherize3[F ~func(T0, T1, T2) Either[error, R], T0, T1, T2, R any](f F) func(T0, T1, T2) (R, error)
func Uneitherize4[F ~func(T0, T1, T2, T3) Either[error, R], T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) (R, error)
func Uneitherize5[F ~func(T0, T1, T2, T3, T4) Either[error, R], T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) (R, error)
func Uneitherize6[F ~func(T0, T1, T2, T3, T4, T5) Either[error, R], T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) (R, error)
func Uneitherize7[F ~func(T0, T1, T2, T3, T4, T5, T6) Either[error, R], T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) (R, error)
func Uneitherize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) (R, error)
func Uneitherize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error)
func Unvariadic0[V, R any](f func(...V) (R, error)) func([]V) Either[error, R]
func Unvariadic1[T1, V, R any](f func(T1, ...V) (R, error)) func(T1, []V) Either[error, R]
func Unvariadic2[T1, T2, V, R any](f func(T1, T2, ...V) (R, error)) func(T1, T2, []V) Either[error, R]
func Unvariadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, ...V) (R, error)) func(T1, T2, T3, []V) Either[error, R]
func Unvariadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, ...V) (R, error)) func(T1, T2, T3, T4, []V) Either[error, R]
func Unwrap[E, A any](ma Either[E, A]) (A, E)
func UnwrapError[A any](ma Either[error, A]) (A, error)
func Variadic0[V, R any](f func([]V) (R, error)) func(...V) Either[error, R]
func Variadic1[T1, V, R any](f func(T1, []V) (R, error)) func(T1, ...V) Either[error, R]
func Variadic2[T1, T2, V, R any](f func(T1, T2, []V) (R, error)) func(T1, T2, ...V) Either[error, R]
func Variadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, []V) (R, error)) func(T1, T2, T3, ...V) Either[error, R]
func Variadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, []V) (R, error)) func(T1, T2, T3, T4, ...V) Either[error, R]
type Either[E, A any] struct {
func Do[E, S any](
func Flatten[E, A any](mma Either[E, Either[E, A]]) Either[E, A]
func FromIO[E any, IO ~func() A, A any](f IO) Either[E, A]
func Left[A, E any](value E) Either[E, A]
func Memoize[E, A any](val Either[E, A]) Either[E, A]
func MonadAlt[E, A any](fa Either[E, A], that Lazy[Either[E, A]]) Either[E, A]
func MonadAp[B, E, A any](fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B]
func MonadBiMap[E1, E2, A, B any](fa Either[E1, A], f func(E1) E2, g func(a A) B) Either[E2, B]
func MonadChain[E, A, B any](fa Either[E, A], f Kleisli[E, A, B]) Either[E, B]
func MonadChainFirst[E, A, B any](ma Either[E, A], f Kleisli[E, A, B]) Either[E, A]
func MonadChainLeft[EA, EB, A any](fa Either[EA, A], f Kleisli[EB, EA, A]) Either[EB, A]
func MonadChainOptionK[A, B, E any](onNone func() E, ma Either[E, A], f func(A) Option[B]) Either[E, B]
func MonadChainTo[A, E, B any](_ Either[E, A], mb Either[E, B]) Either[E, B]
func MonadExtend[E, A, B any](fa Either[E, A], f func(Either[E, A]) B) Either[E, B]
func MonadFlap[E, B, A any](fab Either[E, func(A) B], a A) Either[E, B]
func MonadMap[E, A, B any](fa Either[E, A], f func(a A) B) Either[E, B]
func MonadMapLeft[E1, A, E2 any](fa Either[E1, A], f func(E1) E2) Either[E2, A]
func MonadMapTo[E, A, B any](fa Either[E, A], b B) Either[E, B]
func MonadSequence2[E, T1, T2, R any](e1 Either[E, T1], e2 Either[E, T2], f func(T1, T2) Either[E, R]) Either[E, R]
func MonadSequence3[E, T1, T2, T3, R any](e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3], f func(T1, T2, T3) Either[E, R]) Either[E, R]
func Of[E, A any](value A) Either[E, A]
func Right[E, A any](value A) Either[E, A]
func SequenceArray[E, A any](ma []Either[E, A]) Either[E, []A]
func SequenceArrayG[GA ~[]A, GOA ~[]Either[E, A], E, A any](ma GOA) Either[E, GA]
func SequenceRecord[K comparable, E, A any](ma map[K]Either[E, A]) Either[E, map[K]A]
func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Either[E, A], K comparable, E, A any](ma GOA) Either[E, GA]
func SequenceSeq[E, A any](ma iter.Seq[Either[E, A]]) Either[E, iter.Seq[A]]
func SequenceT1[E, T1 any](t1 Either[E, T1]) Either[E, T.Tuple1[T1]]
func SequenceT10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8], t9 Either[E, T9], t10 Either[E, T10]) Either[E, T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceT11[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8], t9 Either[E, T9], t10 Either[E, T10], t11 Either[E, T11]) Either[E, T.Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func SequenceT12[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8], t9 Either[E, T9], t10 Either[E, T10], t11 Either[E, T11], t12 Either[E, T12]) Either[E, T.Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func SequenceT13[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8], t9 Either[E, T9], t10 Either[E, T10], t11 Either[E, T11], t12 Either[E, T12], t13 Either[E, T13]) Either[E, T.Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func SequenceT14[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8], t9 Either[E, T9], t10 Either[E, T10], t11 Either[E, T11], t12 Either[E, T12], t13 Either[E, T13], t14 Either[E, T14]) Either[E, T.Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func SequenceT15[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8], t9 Either[E, T9], t10 Either[E, T10], t11 Either[E, T11], t12 Either[E, T12], t13 Either[E, T13], t14 Either[E, T14], t15 Either[E, T15]) Either[E, T.Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func SequenceT2[E, T1, T2 any](t1 Either[E, T1], t2 Either[E, T2]) Either[E, T.Tuple2[T1, T2]]
func SequenceT3[E, T1, T2, T3 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3]) Either[E, T.Tuple3[T1, T2, T3]]
func SequenceT4[E, T1, T2, T3, T4 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4]) Either[E, T.Tuple4[T1, T2, T3, T4]]
func SequenceT5[E, T1, T2, T3, T4, T5 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5]) Either[E, T.Tuple5[T1, T2, T3, T4, T5]]
func SequenceT6[E, T1, T2, T3, T4, T5, T6 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6]) Either[E, T.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceT7[E, T1, T2, T3, T4, T5, T6, T7 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7]) Either[E, T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceT8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8]) Either[E, T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceT9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t1 Either[E, T1], t2 Either[E, T2], t3 Either[E, T3], t4 Either[E, T4], t5 Either[E, T5], t6 Either[E, T6], t7 Either[E, T7], t8 Either[E, T8], t9 Either[E, T9]) Either[E, T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceTuple1[E, T1 any](t T.Tuple1[Either[E, T1]]) Either[E, T.Tuple1[T1]]
func SequenceTuple10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t T.Tuple10[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8], Either[E, T9], Either[E, T10]]) Either[E, T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple11[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](t T.Tuple11[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8], Either[E, T9], Either[E, T10], Either[E, T11]]) Either[E, T.Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func SequenceTuple12[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](t T.Tuple12[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8], Either[E, T9], Either[E, T10], Either[E, T11], Either[E, T12]]) Either[E, T.Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func SequenceTuple13[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](t T.Tuple13[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8], Either[E, T9], Either[E, T10], Either[E, T11], Either[E, T12], Either[E, T13]]) Either[E, T.Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func SequenceTuple14[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](t T.Tuple14[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8], Either[E, T9], Either[E, T10], Either[E, T11], Either[E, T12], Either[E, T13], Either[E, T14]]) Either[E, T.Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func SequenceTuple15[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](t T.Tuple15[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8], Either[E, T9], Either[E, T10], Either[E, T11], Either[E, T12], Either[E, T13], Either[E, T14], Either[E, T15]]) Either[E, T.Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func SequenceTuple2[E, T1, T2 any](t T.Tuple2[Either[E, T1], Either[E, T2]]) Either[E, T.Tuple2[T1, T2]]
func SequenceTuple3[E, T1, T2, T3 any](t T.Tuple3[Either[E, T1], Either[E, T2], Either[E, T3]]) Either[E, T.Tuple3[T1, T2, T3]]
func SequenceTuple4[E, T1, T2, T3, T4 any](t T.Tuple4[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4]]) Either[E, T.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[E, T1, T2, T3, T4, T5 any](t T.Tuple5[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5]]) Either[E, T.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[E, T1, T2, T3, T4, T5, T6 any](t T.Tuple6[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6]]) Either[E, T.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[E, T1, T2, T3, T4, T5, T6, T7 any](t T.Tuple7[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7]]) Either[E, T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](t T.Tuple8[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8]]) Either[E, T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t T.Tuple9[Either[E, T1], Either[E, T2], Either[E, T3], Either[E, T4], Either[E, T5], Either[E, T6], Either[E, T7], Either[E, T8], Either[E, T9]]) Either[E, T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Swap[E, A any](val Either[E, A]) Either[A, E]
func TryCatch[FE func(error) E, E, A any](val A, err error, onThrow FE) Either[E, A]
func TryCatchError[A any](val A, err error) Either[error, A]
func Zero[E, A any]() Either[E, A]
func (s Either[E, A]) Format(f fmt.State, c rune)
func (s Either[E, A]) GoString() string
func (s Either[E, A]) LogValue() slog.Value
func (s Either[E, A]) String() string
type Endomorphism[T any] = endomorphism.Endomorphism[T]
func ApSL[E, S, T any](
func BindL[E, S, T any](
func LetL[E, S, T any](
func LetToL[E, S, T any](
type Kleisli[E, A, B any] = reader.Reader[A, Either[E, B]]
func AltW[E, E1, A any](that Lazy[Either[E1, A]]) Kleisli[E1, Either[E, A], A]
func ChainLeft[EA, EB, A any](f Kleisli[EB, EA, A]) Kleisli[EB, Either[EA, A], A]
func FromNillable[A, E any](e E) Kleisli[E, *A, *A]
func FromPredicate[E, A any](pred Predicate[A], onFalse func(A) E) Kleisli[E, A, A]
func OrElse[E1, E2, A any](onLeft Kleisli[E2, E1, A]) Kleisli[E2, Either[E1, A], A]
func TailRec[E, A, B any](f Kleisli[E, A, tailrec.Trampoline[A, B]]) Kleisli[E, A, B]
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B]
func TraverseArrayG[GA ~[]A, GB ~[]B, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, GA, GB]
func TraverseArrayWithIndex[E, A, B any](f func(int, A) Either[E, B]) Kleisli[E, []A, []B]
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, E, A, B any](f func(int, A) Either[E, B]) Kleisli[E, GA, GB]
func TraverseRecord[K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, map[K]A, map[K]B]
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, GA, GB]
func TraverseRecordWithIndex[K comparable, E, A, B any](f func(K, A) Either[E, B]) Kleisli[E, map[K]A, map[K]B]
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, E, A, B any](f func(K, A) Either[E, B]) Kleisli[E, GA, GB]
func TraverseSeq[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, iter.Seq[A], iter.Seq[B]]
func WithResource[A, E, R, ANY any](
type Lazy[T any] = lazy.Lazy[T]
type Lens[S, T any] = lens.Lens[S, T]
type Monoid[E, A any] = monoid.Monoid[Either[E, A]]
func AltMonoid[E, A any](zero Lazy[Either[E, A]]) Monoid[E, A]
func AlternativeMonoid[E, A any](m M.Monoid[A]) Monoid[E, A]
type Operator[E, A, B any] = Kleisli[E, Either[E, A], B]
func Alt[E, A any](that Lazy[Either[E, A]]) Operator[E, A, A]
func Ap[B, E, A any](fa Either[E, A]) Operator[E, func(A) B, B]
func ApS[E, S1, S2, T any](
func Bind[E, S1, S2, T any](
func BindTo[E, S1, T any](
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B]
func ChainFirst[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, A]
func ChainTo[A, E, B any](mb Either[E, B]) Operator[E, A, B]
func Extend[E, A, B any](f func(Either[E, A]) B) Operator[E, A, B]
func Filter[E, A any](p Predicate[A], empty E) Operator[E, A, A]
func FilterMap[E, A, B any](f option.Kleisli[A, B], empty E) Operator[E, A, B]
func FilterOrElse[E, A any](pred Predicate[A], onFalse func(A) E) Operator[E, A, A]
func Flap[E, B, A any](a A) Operator[E, func(A) B, B]
func Let[E, S1, S2, T any](
func LetTo[E, S1, S2, T any](
func Map[E, A, B any](f func(a A) B) Operator[E, A, B]
func MapTo[E, A, B any](b B) Operator[E, A, B]
type Option[A any] = option.Option[A]
func ToOption[E, A any](ma Either[E, A]) Option[A]
type Pair[L, R any] = pair.Pair[L, R]
type Predicate[A any] = predicate.Predicate[A]
```

## package `github.com/IBM/fp-go/v2/result`

Import: `import R "github.com/IBM/fp-go/v2/result"`

Result is Either specialized with error as Left type.

Key types:
- `Result[A] = Either[error, A]` -- the core type
- `Kleisli[A, B] = func(A) Result[B]` -- effectful function
- `Operator[A, B] = Kleisli[Result[A], B]` -- pipeline operator

### Exported API

```go
func AltSemigroup[A any]() S.Semigroup[Result[A]]
func AltW[E1, A any](that Lazy[Either[E1, A]]) func(Result[A]) Either[E1, A]
func ApplicativeMonoid[A any](m M.Monoid[A]) M.Monoid[Result[A]]
func ApplySemigroup[A any](s S.Semigroup[A]) S.Semigroup[Result[A]]
func BiMap[E, A, B any](f func(error) E, g func(a A) B) func(Result[A]) Either[E, B]
func ChainOptionK[A, B any](onNone func() error) func(option.Kleisli[A, B]) Operator[A, B]
func CompactArray[A any](fa []Result[A]) []A
func CompactArrayG[A1 ~[]Result[A], A2 ~[]A, A any](fa A1) A2
func CompactRecord[K comparable, A any](m map[K]Result[A]) map[K]A
func CompactRecordG[M1 ~map[K]Result[A], M2 ~map[K]A, K comparable, A any](m M1) M2
func Curry0[R any](f func() (R, error)) func() Result[R]
func Curry1[T1, R any](f func(T1) (R, error)) func(T1) Result[R]
func Curry2[T1, T2, R any](f func(T1, T2) (R, error)) func(T1) func(T2) Result[R]
func Curry3[T1, T2, T3, R any](f func(T1, T2, T3) (R, error)) func(T1) func(T2) func(T3) Result[R]
func Curry4[T1, T2, T3, T4, R any](f func(T1, T2, T3, T4) (R, error)) func(T1) func(T2) func(T3) func(T4) Result[R]
func Eitherize0[F ~func() (R, error), R any](f F) func() Result[R]
func Eitherize1[F ~func(T0) (R, error), T0, R any](f F) func(T0) Result[R]
func Eitherize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) Result[R]
func Eitherize11[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) Result[R]
func Eitherize12[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) Result[R]
func Eitherize13[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) Result[R]
func Eitherize14[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) Result[R]
func Eitherize15[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) Result[R]
func Eitherize2[F ~func(T0, T1) (R, error), T0, T1, R any](f F) func(T0, T1) Result[R]
func Eitherize3[F ~func(T0, T1, T2) (R, error), T0, T1, T2, R any](f F) func(T0, T1, T2) Result[R]
func Eitherize4[F ~func(T0, T1, T2, T3) (R, error), T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) Result[R]
func Eitherize5[F ~func(T0, T1, T2, T3, T4) (R, error), T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) Result[R]
func Eitherize6[F ~func(T0, T1, T2, T3, T4, T5) (R, error), T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) Result[R]
func Eitherize7[F ~func(T0, T1, T2, T3, T4, T5, T6) (R, error), T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) Result[R]
func Eitherize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) Result[R]
func Eitherize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) Result[R]
func Eq[A any](a eq.Eq[A]) eq.Eq[Result[A]]
func FirstMonoid[A any](zero Lazy[Result[A]]) M.Monoid[Result[A]]
func Fold[A, B any](onLeft func(error) B, onRight func(A) B) func(Result[A]) B
func FromNillable[A any](e error) func(*A) Result[*A]
func FromOption[A any](onNone func() error) func(Option[A]) Result[A]
func FromStrictEquals[A comparable]() eq.Eq[Result[A]]
func Functor[A, B any]() functor.Functor[A, B, Result[A], Result[B]]
func GetOrElse[A any](onLeft func(error) A) func(Result[A]) A
func IsLeft[A any](val Result[A]) bool
func IsRight[A any](val Result[A]) bool
func LastMonoid[A any](zero Lazy[Result[A]]) M.Monoid[Result[A]]
func Logger[A any](loggers ...*log.Logger) func(string) Operator[A, A]
func MapLeft[A, E any](f func(error) E) func(fa Result[A]) Either[E, A]
func Monad[A, B any]() monad.Monad[A, B, Result[A], Result[B], Result[func(A) B]]
func MonadFold[A, B any](ma Result[A], onLeft func(e error) B, onRight func(a A) B) B
func Partition[A any](p Predicate[A], empty error) func(Result[A]) Pair[Result[A], Result[A]]
func PartitionMap[A, B, C any](f either.Kleisli[B, A, C], empty error) func(Result[A]) Pair[Result[B], Result[C]]
func Pointed[A any]() pointed.Pointed[A, Result[A]]
func Reduce[A, B any](f func(B, A) B, initial B) func(Result[A]) B
func Sequence[A, HKTA, HKTRA any](
func Sequence2[T1, T2, R any](f func(T1, T2) Result[R]) func(Result[T1], Result[T2]) Result[R]
func Sequence3[T1, T2, T3, R any](f func(T1, T2, T3) Result[R]) func(Result[T1], Result[T2], Result[T3]) Result[R]
func ToError[A any](e Result[A]) error
func ToSLogAttr[A any]() func(Result[A]) slog.Attr
func Traverse[A, B, HKTB, HKTRB any](
func TraverseTuple1[F1 ~func(A1) Result[T1], A1, T1 any](f1 F1) func(T.Tuple1[A1]) Result[T.Tuple1[T1]]
func TraverseTuple10[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], F9 ~func(A9) Result[T9], F10 ~func(A10) Result[T10], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(T.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) Result[T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple11[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], F9 ~func(A9) Result[T9], F10 ~func(A10) Result[T10], F11 ~func(A11) Result[T11], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) func(T.Tuple11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11]) Result[T.Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func TraverseTuple12[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], F9 ~func(A9) Result[T9], F10 ~func(A10) Result[T10], F11 ~func(A11) Result[T11], F12 ~func(A12) Result[T12], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) func(T.Tuple12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12]) Result[T.Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func TraverseTuple13[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], F9 ~func(A9) Result[T9], F10 ~func(A10) Result[T10], F11 ~func(A11) Result[T11], F12 ~func(A12) Result[T12], F13 ~func(A13) Result[T13], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12, A13, T13 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) func(T.Tuple13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13]) Result[T.Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func TraverseTuple14[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], F9 ~func(A9) Result[T9], F10 ~func(A10) Result[T10], F11 ~func(A11) Result[T11], F12 ~func(A12) Result[T12], F13 ~func(A13) Result[T13], F14 ~func(A14) Result[T14], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12, A13, T13, A14, T14 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) func(T.Tuple14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14]) Result[T.Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func TraverseTuple15[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], F9 ~func(A9) Result[T9], F10 ~func(A10) Result[T10], F11 ~func(A11) Result[T11], F12 ~func(A12) Result[T12], F13 ~func(A13) Result[T13], F14 ~func(A14) Result[T14], F15 ~func(A15) Result[T15], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10, A11, T11, A12, T12, A13, T13, A14, T14, A15, T15 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) func(T.Tuple15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15]) Result[T.Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func TraverseTuple2[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], A1, T1, A2, T2 any](f1 F1, f2 F2) func(T.Tuple2[A1, A2]) Result[T.Tuple2[T1, T2]]
func TraverseTuple3[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], A1, T1, A2, T2, A3, T3 any](f1 F1, f2 F2, f3 F3) func(T.Tuple3[A1, A2, A3]) Result[T.Tuple3[T1, T2, T3]]
func TraverseTuple4[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], A1, T1, A2, T2, A3, T3, A4, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T.Tuple4[A1, A2, A3, A4]) Result[T.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T.Tuple5[A1, A2, A3, A4, A5]) Result[T.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(T.Tuple6[A1, A2, A3, A4, A5, A6]) Result[T.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(T.Tuple7[A1, A2, A3, A4, A5, A6, A7]) Result[T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(T.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) Result[T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[F1 ~func(A1) Result[T1], F2 ~func(A2) Result[T2], F3 ~func(A3) Result[T3], F4 ~func(A4) Result[T4], F5 ~func(A5) Result[T5], F6 ~func(A6) Result[T6], F7 ~func(A7) Result[T7], F8 ~func(A8) Result[T8], F9 ~func(A9) Result[T9], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(T.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) Result[T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Uncurry0[R any](f func() Result[R]) func() (R, error)
func Uncurry1[T1, R any](f func(T1) Result[R]) func(T1) (R, error)
func Uncurry2[T1, T2, R any](f func(T1) func(T2) Result[R]) func(T1, T2) (R, error)
func Uncurry3[T1, T2, T3, R any](f func(T1) func(T2) func(T3) Result[R]) func(T1, T2, T3) (R, error)
func Uncurry4[T1, T2, T3, T4, R any](f func(T1) func(T2) func(T3) func(T4) Result[R]) func(T1, T2, T3, T4) (R, error)
func Uneitherize0[F ~func() Result[R], R any](f F) func() (R, error)
func Uneitherize1[F ~func(T0) Result[R], T0, R any](f F) func(T0) (R, error)
func Uneitherize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error)
func Uneitherize11[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error)
func Uneitherize12[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) (R, error)
func Uneitherize13[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) (R, error)
func Uneitherize14[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) (R, error)
func Uneitherize15[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) (R, error)
func Uneitherize2[F ~func(T0, T1) Result[R], T0, T1, R any](f F) func(T0, T1) (R, error)
func Uneitherize3[F ~func(T0, T1, T2) Result[R], T0, T1, T2, R any](f F) func(T0, T1, T2) (R, error)
func Uneitherize4[F ~func(T0, T1, T2, T3) Result[R], T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) (R, error)
func Uneitherize5[F ~func(T0, T1, T2, T3, T4) Result[R], T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) (R, error)
func Uneitherize6[F ~func(T0, T1, T2, T3, T4, T5) Result[R], T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) (R, error)
func Uneitherize7[F ~func(T0, T1, T2, T3, T4, T5, T6) Result[R], T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) (R, error)
func Uneitherize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) (R, error)
func Uneitherize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) Result[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error)
func Unvariadic0[V, R any](f func(...V) (R, error)) func([]V) Result[R]
func Unvariadic1[T1, V, R any](f func(T1, ...V) (R, error)) func(T1, []V) Result[R]
func Unvariadic2[T1, T2, V, R any](f func(T1, T2, ...V) (R, error)) func(T1, T2, []V) Result[R]
func Unvariadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, ...V) (R, error)) func(T1, T2, T3, []V) Result[R]
func Unvariadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, ...V) (R, error)) func(T1, T2, T3, T4, []V) Result[R]
func Unwrap[A any](ma Result[A]) (A, error)
func UnwrapError[A any](ma Result[A]) (A, error)
func Variadic0[V, R any](f func([]V) (R, error)) func(...V) Result[R]
func Variadic1[T1, V, R any](f func(T1, []V) (R, error)) func(T1, ...V) Result[R]
func Variadic2[T1, T2, V, R any](f func(T1, T2, []V) (R, error)) func(T1, T2, ...V) Result[R]
func Variadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, []V) (R, error)) func(T1, T2, T3, ...V) Result[R]
func Variadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, []V) (R, error)) func(T1, T2, T3, T4, ...V) Result[R]
type Either[E, T any] = either.Either[E, T]
func MonadBiMap[E, A, B any](fa Result[A], f func(error) E, g func(a A) B) Either[E, B]
func MonadMapLeft[A, E any](fa Result[A], f func(error) E) Either[E, A]
func Swap[A any](val Result[A]) Either[A, error]
type Endomorphism[T any] = endomorphism.Endomorphism[T]
type Kleisli[A, B any] = reader.Reader[A, Result[B]]
func FromError[A any](f func(a A) error) Kleisli[A, A]
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A]
func TailRec[A, B any](f Kleisli[A, tailrec.Trampoline[A, B]]) Kleisli[A, B]
func ToType[A any](onError func(any) error) Kleisli[any, A]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseArrayWithIndex[A, B any](f func(int, A) Result[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) Result[B]) Kleisli[GA, GB]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Result[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) Result[B]) Kleisli[GA, GB]
func TraverseSeq[A, B any](f Kleisli[A, B]) Kleisli[iter.Seq[A], iter.Seq[B]]
func WithResource[A, R, ANY any](
type Lazy[T any] = lazy.Lazy[T]
type Lens[S, T any] = lens.Lens[S, T]
type Monoid[A any] = monoid.Monoid[Result[A]]
func AltMonoid[A any](zero Lazy[Result[A]]) Monoid[A]
func AlternativeMonoid[A any](m M.Monoid[A]) Monoid[A]
type Operator[A, B any] = Kleisli[Result[A], B]
func Alt[A any](that Lazy[Result[A]]) Operator[A, A]
func Ap[B, A any](fa Result[A]) Operator[func(A) B, B]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func Bind[S1, S2, T any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A]
func ChainTo[A, B any](mb Result[B]) Operator[A, B]
func Filter[A any](p Predicate[A], empty error) Operator[A, A]
func FilterMap[A, B any](f option.Kleisli[A, B], empty error) Operator[A, B]
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func Map[A, B any](f func(a A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A]
type Option[A any] = option.Option[A]
func ToOption[A any](ma Result[A]) Option[A]
type Pair[L, R any] = pair.Pair[L, R]
type Predicate[A any] = predicate.Predicate[A]
type Result[T any] = Either[error, T]
func Do[S any](
func Flatten[A any](mma Result[Result[A]]) Result[A]
func FromIO[IO ~func() A, A any](f IO) Result[A]
func InstanceOf[A any](a any) Result[A]
func Left[A any](value error) Result[A]
func Memoize[A any](val Result[A]) Result[A]
func MonadAlt[A any](fa Result[A], that Lazy[Result[A]]) Result[A]
func MonadAp[B, A any](fab Result[func(a A) B], fa Result[A]) Result[B]
func MonadChain[A, B any](fa Result[A], f Kleisli[A, B]) Result[B]
func MonadChainFirst[A, B any](ma Result[A], f Kleisli[A, B]) Result[A]
func MonadChainLeft[A any](fa Result[A], f Kleisli[error, A]) Result[A]
func MonadChainOptionK[A, B any](onNone func() error, ma Result[A], f option.Kleisli[A, B]) Result[B]
func MonadChainTo[A, B any](ma Result[A], mb Result[B]) Result[B]
func MonadFlap[B, A any](fab Result[func(A) B], a A) Result[B]
func MonadMap[A, B any](fa Result[A], f func(a A) B) Result[B]
func MonadMapTo[A, B any](fa Result[A], b B) Result[B]
func MonadSequence2[T1, T2, R any](e1 Result[T1], e2 Result[T2], f func(T1, T2) Result[R]) Result[R]
func MonadSequence3[T1, T2, T3, R any](e1 Result[T1], e2 Result[T2], e3 Result[T3], f func(T1, T2, T3) Result[R]) Result[R]
func Of[A any](value A) Result[A]
func Right[A any](value A) Result[A]
func SequenceArray[A any](ma []Result[A]) Result[[]A]
func SequenceArrayG[GA ~[]A, GOA ~[]Result[A], A any](ma GOA) Result[GA]
func SequenceRecord[K comparable, A any](ma map[K]Result[A]) Result[map[K]A]
func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Result[A], K comparable, A any](ma GOA) Result[GA]
func SequenceSeq[A any](ma iter.Seq[Result[A]]) Result[iter.Seq[A]]
func SequenceT1[T1 any](t1 Result[T1]) Result[T.Tuple1[T1]]
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8], t9 Result[T9], t10 Result[T10]) Result[T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceT11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8], t9 Result[T9], t10 Result[T10], t11 Result[T11]) Result[T.Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func SequenceT12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8], t9 Result[T9], t10 Result[T10], t11 Result[T11], t12 Result[T12]) Result[T.Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func SequenceT13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8], t9 Result[T9], t10 Result[T10], t11 Result[T11], t12 Result[T12], t13 Result[T13]) Result[T.Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func SequenceT14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8], t9 Result[T9], t10 Result[T10], t11 Result[T11], t12 Result[T12], t13 Result[T13], t14 Result[T14]) Result[T.Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func SequenceT15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8], t9 Result[T9], t10 Result[T10], t11 Result[T11], t12 Result[T12], t13 Result[T13], t14 Result[T14], t15 Result[T15]) Result[T.Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func SequenceT2[T1, T2 any](t1 Result[T1], t2 Result[T2]) Result[T.Tuple2[T1, T2]]
func SequenceT3[T1, T2, T3 any](t1 Result[T1], t2 Result[T2], t3 Result[T3]) Result[T.Tuple3[T1, T2, T3]]
func SequenceT4[T1, T2, T3, T4 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4]) Result[T.Tuple4[T1, T2, T3, T4]]
func SequenceT5[T1, T2, T3, T4, T5 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5]) Result[T.Tuple5[T1, T2, T3, T4, T5]]
func SequenceT6[T1, T2, T3, T4, T5, T6 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6]) Result[T.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7]) Result[T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8]) Result[T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t1 Result[T1], t2 Result[T2], t3 Result[T3], t4 Result[T4], t5 Result[T5], t6 Result[T6], t7 Result[T7], t8 Result[T8], t9 Result[T9]) Result[T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceTuple1[T1 any](t T.Tuple1[Result[T1]]) Result[T.Tuple1[T1]]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t T.Tuple10[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8], Result[T9], Result[T10]]) Result[T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](t T.Tuple11[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8], Result[T9], Result[T10], Result[T11]]) Result[T.Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func SequenceTuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](t T.Tuple12[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8], Result[T9], Result[T10], Result[T11], Result[T12]]) Result[T.Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func SequenceTuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](t T.Tuple13[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8], Result[T9], Result[T10], Result[T11], Result[T12], Result[T13]]) Result[T.Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func SequenceTuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](t T.Tuple14[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8], Result[T9], Result[T10], Result[T11], Result[T12], Result[T13], Result[T14]]) Result[T.Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func SequenceTuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](t T.Tuple15[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8], Result[T9], Result[T10], Result[T11], Result[T12], Result[T13], Result[T14], Result[T15]]) Result[T.Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func SequenceTuple2[T1, T2 any](t T.Tuple2[Result[T1], Result[T2]]) Result[T.Tuple2[T1, T2]]
func SequenceTuple3[T1, T2, T3 any](t T.Tuple3[Result[T1], Result[T2], Result[T3]]) Result[T.Tuple3[T1, T2, T3]]
func SequenceTuple4[T1, T2, T3, T4 any](t T.Tuple4[Result[T1], Result[T2], Result[T3], Result[T4]]) Result[T.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t T.Tuple5[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5]]) Result[T.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t T.Tuple6[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6]]) Result[T.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t T.Tuple7[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7]]) Result[T.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t T.Tuple8[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8]]) Result[T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t T.Tuple9[Result[T1], Result[T2], Result[T3], Result[T4], Result[T5], Result[T6], Result[T7], Result[T8], Result[T9]]) Result[T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TryCatch[FE Endomorphism[error], A any](val A, err error, onThrow FE) Result[A]
func TryCatchError[A any](val A, err error) Result[A]
func Zero[A any]() Result[A]
```

## package `github.com/IBM/fp-go/v2/io`

Import: `import "github.com/IBM/fp-go/v2/io"`

IO represents a synchronous side-effectful computation.

Key types:
- `IO[A] = func() A` -- lazy computation
- `Kleisli[A, B] = func(A) IO[B]` -- effectful function

### Exported API

```go
func Eq[A any](e EQ.Eq[A]) EQ.Eq[IO[A]]
func FromStrictEquals[A comparable]() EQ.Eq[IO[A]]
func Logger[A any](loggers ...*log.Logger) func(string) Kleisli[A, A]
func Run[A any](fa IO[A]) A
func TraverseParTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IO[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseParTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IO[tuple.Tuple2[T1, T2]]
func TraverseParTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IO[tuple.Tuple3[T1, T2, T3]]
func TraverseParTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IO[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseParTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IO[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseParTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IO[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseParTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IO[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseParTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IO[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseParTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IO[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseSeqTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IO[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseSeqTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IO[tuple.Tuple2[T1, T2]]
func TraverseSeqTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IO[tuple.Tuple3[T1, T2, T3]]
func TraverseSeqTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IO[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseSeqTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IO[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseSeqTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IO[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseSeqTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IO[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseSeqTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IO[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseSeqTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IO[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IO[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IO[tuple.Tuple2[T1, T2]]
func TraverseTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IO[tuple.Tuple3[T1, T2, T3]]
func TraverseTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IO[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IO[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IO[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IO[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IO[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IO[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
type Consumer[A any] = consumer.Consumer[A]
type IO[A any] = func() A
var Now IO[time.Time] = time.Now
func Bracket[A, B, ANY any](
func Defer[A any](gen func() IO[A]) IO[A]
func Do[S any](
func Flatten[A any](mma IO[IO[A]]) IO[A]
func FromIO[A any](a IO[A]) IO[A]
func FromImpure[ANY ~func()](f ANY) IO[Void]
func Memoize[A any](ma IO[A]) IO[A]
func MonadAp[A, B any](mab IO[func(A) B], ma IO[A]) IO[B]
func MonadApFirst[A, B any](first IO[A], second IO[B]) IO[A]
func MonadApPar[A, B any](mab IO[func(A) B], ma IO[A]) IO[B]
func MonadApSecond[A, B any](first IO[A], second IO[B]) IO[B]
func MonadApSeq[A, B any](mab IO[func(A) B], ma IO[A]) IO[B]
func MonadChain[A, B any](fa IO[A], f Kleisli[A, B]) IO[B]
func MonadChainFirst[A, B any](fa IO[A], f Kleisli[A, B]) IO[A]
func MonadChainTo[A, B any](fa IO[A], fb IO[B]) IO[B]
func MonadFlap[B, A any](fab IO[func(A) B], a A) IO[B]
func MonadMap[A, B any](fa IO[A], f func(A) B) IO[B]
func MonadMapTo[A, B any](fa IO[A], b B) IO[B]
func MonadOf[A any](a A) IO[A]
func MonadTraverseArray[A, B any](tas []A, f Kleisli[A, B]) IO[[]B]
func MonadTraverseArraySeq[A, B any](tas []A, f Kleisli[A, B]) IO[[]B]
func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f Kleisli[A, B]) IO[map[K]B]
func MonadTraverseRecordSeq[K comparable, A, B any](tas map[K]A, f Kleisli[A, B]) IO[map[K]B]
func Of[A any](a A) IO[A]
func Retrying[A any](
func SequenceArray[A any](tas []IO[A]) IO[[]A]
func SequenceArraySeq[A any](tas []IO[A]) IO[[]A]
func SequenceIter[A any](as Seq[IO[A]]) IO[Seq[A]]
func SequenceParT1[T1 any](
func SequenceParT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceParT2[T1, T2 any](
func SequenceParT3[T1, T2, T3 any](
func SequenceParT4[T1, T2, T3, T4 any](
func SequenceParT5[T1, T2, T3, T4, T5 any](
func SequenceParT6[T1, T2, T3, T4, T5, T6 any](
func SequenceParT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceParT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceParT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceParTuple1[T1 any](t tuple.Tuple1[IO[T1]]) IO[tuple.Tuple1[T1]]
func SequenceParTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8], IO[T9], IO[T10]]) IO[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceParTuple2[T1, T2 any](t tuple.Tuple2[IO[T1], IO[T2]]) IO[tuple.Tuple2[T1, T2]]
func SequenceParTuple3[T1, T2, T3 any](t tuple.Tuple3[IO[T1], IO[T2], IO[T3]]) IO[tuple.Tuple3[T1, T2, T3]]
func SequenceParTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IO[T1], IO[T2], IO[T3], IO[T4]]) IO[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceParTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5]]) IO[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceParTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6]]) IO[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceParTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7]]) IO[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceParTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8]]) IO[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceParTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8], IO[T9]]) IO[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceRecord[K comparable, A any](tas map[K]IO[A]) IO[map[K]A]
func SequenceRecordSeq[K comparable, A any](tas map[K]IO[A]) IO[map[K]A]
func SequenceSeqT1[T1 any](
func SequenceSeqT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceSeqT2[T1, T2 any](
func SequenceSeqT3[T1, T2, T3 any](
func SequenceSeqT4[T1, T2, T3, T4 any](
func SequenceSeqT5[T1, T2, T3, T4, T5 any](
func SequenceSeqT6[T1, T2, T3, T4, T5, T6 any](
func SequenceSeqT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceSeqT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceSeqT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceSeqTuple1[T1 any](t tuple.Tuple1[IO[T1]]) IO[tuple.Tuple1[T1]]
func SequenceSeqTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8], IO[T9], IO[T10]]) IO[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceSeqTuple2[T1, T2 any](t tuple.Tuple2[IO[T1], IO[T2]]) IO[tuple.Tuple2[T1, T2]]
func SequenceSeqTuple3[T1, T2, T3 any](t tuple.Tuple3[IO[T1], IO[T2], IO[T3]]) IO[tuple.Tuple3[T1, T2, T3]]
func SequenceSeqTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IO[T1], IO[T2], IO[T3], IO[T4]]) IO[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceSeqTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5]]) IO[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceSeqTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6]]) IO[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceSeqTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7]]) IO[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceSeqTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8]]) IO[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceSeqTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8], IO[T9]]) IO[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceT1[T1 any](
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceT2[T1, T2 any](
func SequenceT3[T1, T2, T3 any](
func SequenceT4[T1, T2, T3, T4 any](
func SequenceT5[T1, T2, T3, T4, T5 any](
func SequenceT6[T1, T2, T3, T4, T5, T6 any](
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceTuple1[T1 any](t tuple.Tuple1[IO[T1]]) IO[tuple.Tuple1[T1]]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8], IO[T9], IO[T10]]) IO[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple2[T1, T2 any](t tuple.Tuple2[IO[T1], IO[T2]]) IO[tuple.Tuple2[T1, T2]]
func SequenceTuple3[T1, T2, T3 any](t tuple.Tuple3[IO[T1], IO[T2], IO[T3]]) IO[tuple.Tuple3[T1, T2, T3]]
func SequenceTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IO[T1], IO[T2], IO[T3], IO[T4]]) IO[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5]]) IO[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6]]) IO[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7]]) IO[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8]]) IO[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IO[T1], IO[T2], IO[T3], IO[T4], IO[T5], IO[T6], IO[T7], IO[T8], IO[T9]]) IO[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func WithDuration[A any](a IO[A]) IO[Pair[time.Duration, A]]
func WithTime[A any](a IO[A]) IO[Pair[Pair[time.Time, time.Time], A]]
type IOApplicative[A, B any] = applicative.Applicative[A, B, IO[A], IO[B], IO[func(A) B]]
func Applicative[A, B any]() IOApplicative[A, B]
type IOFunctor[A, B any] = functor.Functor[A, B, IO[A], IO[B]]
func Functor[A, B any]() IOFunctor[A, B]
type IOMonad[A, B any] = monad.Monad[A, B, IO[A], IO[B], IO[func(A) B]]
func Monad[A, B any]() IOMonad[A, B]
type IOPointed[A any] = pointed.Pointed[A, IO[A]]
func Pointed[A any]() IOPointed[A]
type Kleisli[A, B any] = reader.Reader[A, IO[B]]
func FromConsumer[A any](c Consumer[A]) Kleisli[A, Void]
func LogGo[A any](prefix string) Kleisli[A, A]
func Logf[A any](prefix string) Kleisli[A, A]
func PrintGo[A any](prefix string) Kleisli[A, A]
func Printf[A any](prefix string) Kleisli[A, A]
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) IO[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) IO[B]) Kleisli[[]A, []B]
func TraverseIter[A, B any](f Kleisli[A, B]) Kleisli[Seq[A], Seq[B]]
func TraverseParTuple1[F1 ~Kleisli[A1, T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndeSeq[K comparable, A, B any](f func(K, A) IO[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IO[B]) Kleisli[map[K]A, map[K]B]
func TraverseSeqTuple1[F1 ~Kleisli[A1, T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func TraverseTuple1[F1 ~Kleisli[A1, T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func WithResource[
type Monoid[A any] = M.Monoid[IO[A]]
func ApplicativeMonoid[A any](m M.Monoid[A]) Monoid[A]
type Operator[A, B any] = Kleisli[IO[A], B]
func After[A any](timestamp time.Time) Operator[A, A]
func Ap[B, A any](ma IO[A]) Operator[func(A) B, B]
func ApFirst[A, B any](second IO[B]) Operator[A, A]
func ApPar[B, A any](ma IO[A]) Operator[func(A) B, B]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func ApSecond[A, B any](second IO[B]) Operator[A, B]
func ApSeq[B, A any](ma IO[A]) Operator[func(A) B, B]
func Bind[S1, S2, T any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainConsumer[A any](c Consumer[A]) Operator[A, Void]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainTo[A, B any](fb IO[B]) Operator[A, B]
func Delay[A any](delay time.Duration) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func Map[A, B any](f func(A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func WithLock[A any](lock IO[context.CancelFunc]) Operator[A, A]
type Pair[L, R any] = pair.Pair[L, R]
type Predicate[A any] = predicate.Predicate[A]
type RetryStatus = IO[R.RetryStatus]
type Semigroup[A any] = S.Semigroup[IO[A]]
func ApplySemigroup[A any](s S.Semigroup[A]) Semigroup[A]
type Seq[T any] = iter.Seq[T]
type Trampoline[B, L any] = tailrec.Trampoline[B, L]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/iooption`

Import: `import "github.com/IBM/fp-go/v2/iooption"`

IOOption combines IO and Option: `IOOption[A] = IO[Option[A]]` = `func() Option[A]`.

### Exported API

```go
func Eq[A any](eq EQ.Eq[A]) EQ.Eq[IOOption[A]]
func Fold[A, B any](onNone IO[B], onSome io.Kleisli[A, B]) func(IOOption[A]) IO[B]
func FromStrictEquals[A comparable]() EQ.Eq[IOOption[A]]
func Optionize2[T1, T2, A any](f func(t1 T1, t2 T2) (A, bool)) func(T1, T2) IOOption[A]
func Optionize3[T1, T2, T3, A any](f func(t1 T1, t2 T2, t3 T3) (A, bool)) func(T1, T2, T3) IOOption[A]
func Optionize4[T1, T2, T3, T4, A any](f func(t1 T1, t2 T2, t3 T3, t4 T4) (A, bool)) func(T1, T2, T3, T4) IOOption[A]
func TraverseParTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOOption[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseParTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOOption[tuple.Tuple2[T1, T2]]
func TraverseParTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOOption[tuple.Tuple3[T1, T2, T3]]
func TraverseParTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOOption[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseParTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOOption[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseParTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOOption[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseParTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOOption[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseParTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOOption[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseParTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOOption[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseSeqTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOOption[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseSeqTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOOption[tuple.Tuple2[T1, T2]]
func TraverseSeqTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOOption[tuple.Tuple3[T1, T2, T3]]
func TraverseSeqTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOOption[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseSeqTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOOption[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseSeqTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOOption[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseSeqTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOOption[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseSeqTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOOption[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseSeqTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOOption[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOOption[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOOption[tuple.Tuple2[T1, T2]]
func TraverseTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOOption[tuple.Tuple3[T1, T2, T3]]
func TraverseTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOOption[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOOption[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOOption[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOOption[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOOption[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOOption[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func WithLock[E, A any](lock IO[context.CancelFunc]) func(fa IOOption[A]) IOOption[A]
type Consumer[A any] = consumer.Consumer[A]
type Either[E, A any] = either.Either[E, A]
type IO[A any] = io.IO[A]
type IOOption[A any] = io.IO[Option[A]]
func Bracket[A, B, ANY any](
func Defer[A any](gen func() IOOption[A]) IOOption[A]
func Do[S any](
func Flatten[A any](mma IOOption[IOOption[A]]) IOOption[A]
func FromEither[E, A any](e Either[E, A]) IOOption[A]
func FromIO[A any](mr IO[A]) IOOption[A]
func FromOption[A any](o Option[A]) IOOption[A]
func Memoize[A any](ma IOOption[A]) IOOption[A]
func MonadAlt[A any](first, second IOOption[A]) IOOption[A]
func MonadAp[B, A any](mab IOOption[func(A) B], ma IOOption[A]) IOOption[B]
func MonadChain[A, B any](fa IOOption[A], f Kleisli[A, B]) IOOption[B]
func MonadChainFirst[A, B any](ma IOOption[A], f Kleisli[A, B]) IOOption[A]
func MonadChainFirstIOK[A, B any](first IOOption[A], f io.Kleisli[A, B]) IOOption[A]
func MonadChainIOK[A, B any](ma IOOption[A], f io.Kleisli[A, B]) IOOption[B]
func MonadMap[A, B any](fa IOOption[A], f func(A) B) IOOption[B]
func MonadOf[A any](r A) IOOption[A]
func None[A any]() IOOption[A]
func Of[A any](r A) IOOption[A]
func Retrying[A any](
func SequenceArray[A any](ma []IOOption[A]) IOOption[[]A]
func SequenceParT1[T1 any](
func SequenceParT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceParT2[T1, T2 any](
func SequenceParT3[T1, T2, T3 any](
func SequenceParT4[T1, T2, T3, T4 any](
func SequenceParT5[T1, T2, T3, T4, T5 any](
func SequenceParT6[T1, T2, T3, T4, T5, T6 any](
func SequenceParT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceParT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceParT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceParTuple1[T1 any](t tuple.Tuple1[IOOption[T1]]) IOOption[tuple.Tuple1[T1]]
func SequenceParTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8], IOOption[T9], IOOption[T10]]) IOOption[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceParTuple2[T1, T2 any](t tuple.Tuple2[IOOption[T1], IOOption[T2]]) IOOption[tuple.Tuple2[T1, T2]]
func SequenceParTuple3[T1, T2, T3 any](t tuple.Tuple3[IOOption[T1], IOOption[T2], IOOption[T3]]) IOOption[tuple.Tuple3[T1, T2, T3]]
func SequenceParTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4]]) IOOption[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceParTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5]]) IOOption[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceParTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6]]) IOOption[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceParTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7]]) IOOption[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceParTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8]]) IOOption[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceParTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8], IOOption[T9]]) IOOption[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceSeqT1[T1 any](
func SequenceSeqT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceSeqT2[T1, T2 any](
func SequenceSeqT3[T1, T2, T3 any](
func SequenceSeqT4[T1, T2, T3, T4 any](
func SequenceSeqT5[T1, T2, T3, T4, T5 any](
func SequenceSeqT6[T1, T2, T3, T4, T5, T6 any](
func SequenceSeqT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceSeqT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceSeqT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceSeqTuple1[T1 any](t tuple.Tuple1[IOOption[T1]]) IOOption[tuple.Tuple1[T1]]
func SequenceSeqTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8], IOOption[T9], IOOption[T10]]) IOOption[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceSeqTuple2[T1, T2 any](t tuple.Tuple2[IOOption[T1], IOOption[T2]]) IOOption[tuple.Tuple2[T1, T2]]
func SequenceSeqTuple3[T1, T2, T3 any](t tuple.Tuple3[IOOption[T1], IOOption[T2], IOOption[T3]]) IOOption[tuple.Tuple3[T1, T2, T3]]
func SequenceSeqTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4]]) IOOption[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceSeqTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5]]) IOOption[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceSeqTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6]]) IOOption[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceSeqTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7]]) IOOption[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceSeqTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8]]) IOOption[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceSeqTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8], IOOption[T9]]) IOOption[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceT1[T1 any](
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceT2[T1, T2 any](
func SequenceT3[T1, T2, T3 any](
func SequenceT4[T1, T2, T3, T4 any](
func SequenceT5[T1, T2, T3, T4, T5 any](
func SequenceT6[T1, T2, T3, T4, T5, T6 any](
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceTuple1[T1 any](t tuple.Tuple1[IOOption[T1]]) IOOption[tuple.Tuple1[T1]]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8], IOOption[T9], IOOption[T10]]) IOOption[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple2[T1, T2 any](t tuple.Tuple2[IOOption[T1], IOOption[T2]]) IOOption[tuple.Tuple2[T1, T2]]
func SequenceTuple3[T1, T2, T3 any](t tuple.Tuple3[IOOption[T1], IOOption[T2], IOOption[T3]]) IOOption[tuple.Tuple3[T1, T2, T3]]
func SequenceTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4]]) IOOption[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5]]) IOOption[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6]]) IOOption[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7]]) IOOption[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8]]) IOOption[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOOption[T1], IOOption[T2], IOOption[T3], IOOption[T4], IOOption[T5], IOOption[T6], IOOption[T7], IOOption[T8], IOOption[T9]]) IOOption[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Some[A any](r A) IOOption[A]
type Kleisli[A, B any] = reader.Reader[A, IOOption[B]]
func Optionize1[T1, A any](f func(t1 T1) (A, bool)) Kleisli[T1, A]
func TailRec[A, B any](f Kleisli[A, tailrec.Trampoline[A, B]]) Kleisli[A, B]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) IOOption[B]) Kleisli[[]A, []B]
func TraverseParTuple1[F1 ~Kleisli[A1, T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func TraverseSeqTuple1[F1 ~Kleisli[A1, T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func TraverseTuple1[F1 ~Kleisli[A1, T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func WithResource[
type Lazy[A any] = lazy.Lazy[A]
func Optionize0[A any](f func() (A, bool)) Lazy[IOOption[A]]
type Lens[S, T any] = lens.Lens[S, T]
type Operator[A, B any] = Kleisli[IOOption[A], B]
func After[A any](timestamp time.Time) Operator[A, A]
func Alt[A any](second IOOption[A]) Operator[A, A]
func Ap[B, A any](ma IOOption[A]) Operator[func(A) B, B]
func ApPar[B, A any](ma IOOption[A]) Operator[func(A) B, B]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func ApSeq[B, A any](ma IOOption[A]) Operator[func(A) B, B]
func Bind[S1, S2, T any](
func BindL[S, T any](
func BindTo[S1, T any](
func BindToP[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainConsumer[A any](c Consumer[A]) Operator[A, struct{}]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainFirstConsumer[A any](c Consumer[A]) Operator[A, A]
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B]
func ChainOptionK[A, B any](f func(A) Option[B]) Operator[A, B]
func Delay[A any](delay time.Duration) Operator[A, A]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func Map[A, B any](f func(A) B) Operator[A, B]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Prism[S, T any] = prism.Prism[S, T]
type Trampoline[B, L any] = tailrec.Trampoline[B, L]
```

## package `github.com/IBM/fp-go/v2/ioeither`

Import: `import "github.com/IBM/fp-go/v2/ioeither"`

IOEither combines IO and Either: `IOEither[E, A] = IO[Either[E, A]]` = `func() Either[E, A]`.

### Exported API

```go
func ApSeq[B, E, A any](ma IOEither[E, A]) func(IOEither[E, func(A) B]) IOEither[E, B]
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(A) B) func(IOEither[E1, A]) IOEither[E2, B]
func ChainLeft[EA, EB, A any](f Kleisli[EB, EA, A]) func(IOEither[EA, A]) IOEither[EB, A]
func ChainOptionK[A, B, E any](onNone func() E) func(func(A) O.Option[B]) Operator[E, A, B]
func Eitherize0[F ~func() (R, error), R any](f F) func() IOEither[error, R]
func Eitherize1[F ~func(T1) (R, error), T1, R any](f F) func(T1) IOEither[error, R]
func Eitherize10[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) IOEither[error, R]
func Eitherize2[F ~func(T1, T2) (R, error), T1, T2, R any](f F) func(T1, T2) IOEither[error, R]
func Eitherize3[F ~func(T1, T2, T3) (R, error), T1, T2, T3, R any](f F) func(T1, T2, T3) IOEither[error, R]
func Eitherize4[F ~func(T1, T2, T3, T4) (R, error), T1, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) IOEither[error, R]
func Eitherize5[F ~func(T1, T2, T3, T4, T5) (R, error), T1, T2, T3, T4, T5, R any](f F) func(T1, T2, T3, T4, T5) IOEither[error, R]
func Eitherize6[F ~func(T1, T2, T3, T4, T5, T6) (R, error), T1, T2, T3, T4, T5, T6, R any](f F) func(T1, T2, T3, T4, T5, T6) IOEither[error, R]
func Eitherize7[F ~func(T1, T2, T3, T4, T5, T6, T7) (R, error), T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T1, T2, T3, T4, T5, T6, T7) IOEither[error, R]
func Eitherize8[F ~func(T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8) IOEither[error, R]
func Eitherize9[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9) IOEither[error, R]
func Eq[E, A any](eq EQ.Eq[Either[E, A]]) EQ.Eq[IOEither[E, A]]
func Fold[E, A, B any](onLeft func(E) IO[B], onRight io.Kleisli[A, B]) func(IOEither[E, A]) IO[B]
func FromIOOption[A, E any](onNone func() E) func(o IOO.IOOption[A]) IOEither[E, A]
func FromOption[A, E any](onNone func() E) func(o O.Option[A]) IOEither[E, A]
func FromStrictEquals[E, A comparable]() EQ.Eq[IOEither[E, A]]
func Functor[E, A, B any]() functor.Functor[A, B, IOEither[E, A], IOEither[E, B]]
func GetOrElse[E, A any](onLeft func(E) IO[A]) func(IOEither[E, A]) IO[A]
func GetOrElseOf[E, A any](onLeft func(E) A) func(IOEither[E, A]) IO[A]
func MapLeft[A, E1, E2 any](f func(E1) E2) func(IOEither[E1, A]) IOEither[E2, A]
func Monad[E, A, B any]() monad.Monad[A, B, IOEither[E, A], IOEither[E, B], IOEither[E, func(A) B]]
func Pointed[E, A any]() pointed.Pointed[A, IOEither[E, A]]
func ToIOOption[E, A any](ioe IOEither[E, A]) IOO.IOOption[A]
func TraverseParTuple1[E error, F1 ~func(A1) IOEither[E, T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOEither[E, tuple.Tuple1[T1]]
func TraverseParTuple10[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], F9 ~func(A9) IOEither[E, T9], F10 ~func(A10) IOEither[E, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOEither[E, tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseParTuple2[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOEither[E, tuple.Tuple2[T1, T2]]
func TraverseParTuple3[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOEither[E, tuple.Tuple3[T1, T2, T3]]
func TraverseParTuple4[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOEither[E, tuple.Tuple4[T1, T2, T3, T4]]
func TraverseParTuple5[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOEither[E, tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseParTuple6[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOEither[E, tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseParTuple7[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOEither[E, tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseParTuple8[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOEither[E, tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseParTuple9[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], F9 ~func(A9) IOEither[E, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOEither[E, tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseSeqTuple1[E error, F1 ~func(A1) IOEither[E, T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOEither[E, tuple.Tuple1[T1]]
func TraverseSeqTuple10[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], F9 ~func(A9) IOEither[E, T9], F10 ~func(A10) IOEither[E, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOEither[E, tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseSeqTuple2[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOEither[E, tuple.Tuple2[T1, T2]]
func TraverseSeqTuple3[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOEither[E, tuple.Tuple3[T1, T2, T3]]
func TraverseSeqTuple4[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOEither[E, tuple.Tuple4[T1, T2, T3, T4]]
func TraverseSeqTuple5[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOEither[E, tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseSeqTuple6[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOEither[E, tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseSeqTuple7[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOEither[E, tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseSeqTuple8[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOEither[E, tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseSeqTuple9[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], F9 ~func(A9) IOEither[E, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOEither[E, tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseTuple1[E error, F1 ~func(A1) IOEither[E, T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOEither[E, tuple.Tuple1[T1]]
func TraverseTuple10[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], F9 ~func(A9) IOEither[E, T9], F10 ~func(A10) IOEither[E, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOEither[E, tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple2[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOEither[E, tuple.Tuple2[T1, T2]]
func TraverseTuple3[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOEither[E, tuple.Tuple3[T1, T2, T3]]
func TraverseTuple4[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOEither[E, tuple.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOEither[E, tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOEither[E, tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOEither[E, tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOEither[E, tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[E error, F1 ~func(A1) IOEither[E, T1], F2 ~func(A2) IOEither[E, T2], F3 ~func(A3) IOEither[E, T3], F4 ~func(A4) IOEither[E, T4], F5 ~func(A5) IOEither[E, T5], F6 ~func(A6) IOEither[E, T6], F7 ~func(A7) IOEither[E, T7], F8 ~func(A8) IOEither[E, T8], F9 ~func(A9) IOEither[E, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOEither[E, tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Uneitherize0[F ~func() IOEither[error, R], R any](f F) func() (R, error)
func Uneitherize1[F ~func(T1) IOEither[error, R], T1, R any](f F) func(T1) (R, error)
func Uneitherize10[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) IOEither[error, R], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error)
func Uneitherize2[F ~func(T1, T2) IOEither[error, R], T1, T2, R any](f F) func(T1, T2) (R, error)
func Uneitherize3[F ~func(T1, T2, T3) IOEither[error, R], T1, T2, T3, R any](f F) func(T1, T2, T3) (R, error)
func Uneitherize4[F ~func(T1, T2, T3, T4) IOEither[error, R], T1, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) (R, error)
func Uneitherize5[F ~func(T1, T2, T3, T4, T5) IOEither[error, R], T1, T2, T3, T4, T5, R any](f F) func(T1, T2, T3, T4, T5) (R, error)
func Uneitherize6[F ~func(T1, T2, T3, T4, T5, T6) IOEither[error, R], T1, T2, T3, T4, T5, T6, R any](f F) func(T1, T2, T3, T4, T5, T6) (R, error)
func Uneitherize7[F ~func(T1, T2, T3, T4, T5, T6, T7) IOEither[error, R], T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T1, T2, T3, T4, T5, T6, T7) (R, error)
func Uneitherize8[F ~func(T1, T2, T3, T4, T5, T6, T7, T8) IOEither[error, R], T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8) (R, error)
func Uneitherize9[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9) IOEither[error, R], T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error)
type Consumer[A any] = consumer.Consumer[A]
type Either[E, A any] = either.Either[E, A]
type IO[A any] = io.IO[A]
func MonadFold[E, A, B any](ma IOEither[E, A], onLeft func(E) IO[B], onRight io.Kleisli[A, B]) IO[B]
type IOEither[E, A any] = IO[Either[E, A]]
func Bracket[E, A, B, ANY any](
func Defer[E, A any](gen lazy.Lazy[IOEither[E, A]]) IOEither[E, A]
func Do[E, S any](
func Flatten[E, A any](mma IOEither[E, IOEither[E, A]]) IOEither[E, A]
func FromEither[E, A any](e Either[E, A]) IOEither[E, A]
func FromIO[E, A any](mr IO[A]) IOEither[E, A]
func FromImpure[E any](f func()) IOEither[E, Void]
func FromLazy[E, A any](mr lazy.Lazy[A]) IOEither[E, A]
func Left[A, E any](l E) IOEither[E, A]
func LeftIO[A, E any](ml IO[E]) IOEither[E, A]
func Memoize[E, A any](ma IOEither[E, A]) IOEither[E, A]
func MonadAlt[E, A any](first IOEither[E, A], second lazy.Lazy[IOEither[E, A]]) IOEither[E, A]
func MonadAp[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B]
func MonadApFirst[A, E, B any](first IOEither[E, A], second IOEither[E, B]) IOEither[E, A]
func MonadApPar[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B]
func MonadApSecond[A, E, B any](first IOEither[E, A], second IOEither[E, B]) IOEither[E, B]
func MonadApSeq[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B]
func MonadBiMap[E1, E2, A, B any](fa IOEither[E1, A], f func(E1) E2, g func(A) B) IOEither[E2, B]
func MonadChain[E, A, B any](fa IOEither[E, A], f Kleisli[E, A, B]) IOEither[E, B]
func MonadChainEitherK[E, A, B any](ma IOEither[E, A], f either.Kleisli[E, A, B]) IOEither[E, B]
func MonadChainFirst[E, A, B any](ma IOEither[E, A], f Kleisli[E, A, B]) IOEither[E, A]
func MonadChainFirstEitherK[A, E, B any](ma IOEither[E, A], f either.Kleisli[E, A, B]) IOEither[E, A]
func MonadChainFirstIOK[E, A, B any](ma IOEither[E, A], f io.Kleisli[A, B]) IOEither[E, A]
func MonadChainFirstLeft[A, EA, EB, B any](ma IOEither[EA, A], f Kleisli[EB, EA, B]) IOEither[EA, A]
func MonadChainIOK[E, A, B any](ma IOEither[E, A], f io.Kleisli[A, B]) IOEither[E, B]
func MonadChainLeft[EA, EB, A any](fa IOEither[EA, A], f Kleisli[EB, EA, A]) IOEither[EB, A]
func MonadChainTo[A, E, B any](fa IOEither[E, A], fb IOEither[E, B]) IOEither[E, B]
func MonadChainToIO[E, A, B any](fa IOEither[E, A], fb IO[B]) IOEither[E, B]
func MonadFlap[E, B, A any](fab IOEither[E, func(A) B], a A) IOEither[E, B]
func MonadMap[E, A, B any](fa IOEither[E, A], f func(A) B) IOEither[E, B]
func MonadMapLeft[A, E1, E2 any](fa IOEither[E1, A], f func(E1) E2) IOEither[E2, A]
func MonadMapTo[E, A, B any](fa IOEither[E, A], b B) IOEither[E, B]
func MonadOf[E, A any](r A) IOEither[E, A]
func MonadTap[E, A, B any](ma IOEither[E, A], f Kleisli[E, A, B]) IOEither[E, A]
func MonadTapEitherK[A, E, B any](ma IOEither[E, A], f either.Kleisli[E, A, B]) IOEither[E, A]
func MonadTapIOK[E, A, B any](ma IOEither[E, A], f io.Kleisli[A, B]) IOEither[E, A]
func MonadTapLeft[A, EA, EB, B any](ma IOEither[EA, A], f Kleisli[EB, EA, B]) IOEither[EA, A]
func Of[E, A any](r A) IOEither[E, A]
func Retrying[E, A any](
func Right[E, A any](r A) IOEither[E, A]
func RightIO[E, A any](mr IO[A]) IOEither[E, A]
func SequenceArray[E, A any](ma []IOEither[E, A]) IOEither[E, []A]
func SequenceArrayPar[E, A any](ma []IOEither[E, A]) IOEither[E, []A]
func SequenceArraySeq[E, A any](ma []IOEither[E, A]) IOEither[E, []A]
func SequenceParT1[E, T1 any](
func SequenceParT10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceParT2[E, T1, T2 any](
func SequenceParT3[E, T1, T2, T3 any](
func SequenceParT4[E, T1, T2, T3, T4 any](
func SequenceParT5[E, T1, T2, T3, T4, T5 any](
func SequenceParT6[E, T1, T2, T3, T4, T5, T6 any](
func SequenceParT7[E, T1, T2, T3, T4, T5, T6, T7 any](
func SequenceParT8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceParT9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceParTuple1[E, T1 any](t tuple.Tuple1[IOEither[E, T1]]) IOEither[E, tuple.Tuple1[T1]]
func SequenceParTuple10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8], IOEither[E, T9], IOEither[E, T10]]) IOEither[E, tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceParTuple2[E, T1, T2 any](t tuple.Tuple2[IOEither[E, T1], IOEither[E, T2]]) IOEither[E, tuple.Tuple2[T1, T2]]
func SequenceParTuple3[E, T1, T2, T3 any](t tuple.Tuple3[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3]]) IOEither[E, tuple.Tuple3[T1, T2, T3]]
func SequenceParTuple4[E, T1, T2, T3, T4 any](t tuple.Tuple4[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4]]) IOEither[E, tuple.Tuple4[T1, T2, T3, T4]]
func SequenceParTuple5[E, T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5]]) IOEither[E, tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceParTuple6[E, T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6]]) IOEither[E, tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceParTuple7[E, T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7]]) IOEither[E, tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceParTuple8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8]]) IOEither[E, tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceParTuple9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8], IOEither[E, T9]]) IOEither[E, tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceRecord[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A]
func SequenceRecordPar[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A]
func SequenceRecordSeq[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A]
func SequenceSeqT1[E, T1 any](
func SequenceSeqT10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceSeqT2[E, T1, T2 any](
func SequenceSeqT3[E, T1, T2, T3 any](
func SequenceSeqT4[E, T1, T2, T3, T4 any](
func SequenceSeqT5[E, T1, T2, T3, T4, T5 any](
func SequenceSeqT6[E, T1, T2, T3, T4, T5, T6 any](
func SequenceSeqT7[E, T1, T2, T3, T4, T5, T6, T7 any](
func SequenceSeqT8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceSeqT9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceSeqTuple1[E, T1 any](t tuple.Tuple1[IOEither[E, T1]]) IOEither[E, tuple.Tuple1[T1]]
func SequenceSeqTuple10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8], IOEither[E, T9], IOEither[E, T10]]) IOEither[E, tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceSeqTuple2[E, T1, T2 any](t tuple.Tuple2[IOEither[E, T1], IOEither[E, T2]]) IOEither[E, tuple.Tuple2[T1, T2]]
func SequenceSeqTuple3[E, T1, T2, T3 any](t tuple.Tuple3[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3]]) IOEither[E, tuple.Tuple3[T1, T2, T3]]
func SequenceSeqTuple4[E, T1, T2, T3, T4 any](t tuple.Tuple4[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4]]) IOEither[E, tuple.Tuple4[T1, T2, T3, T4]]
func SequenceSeqTuple5[E, T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5]]) IOEither[E, tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceSeqTuple6[E, T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6]]) IOEither[E, tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceSeqTuple7[E, T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7]]) IOEither[E, tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceSeqTuple8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8]]) IOEither[E, tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceSeqTuple9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8], IOEither[E, T9]]) IOEither[E, tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceT1[E, T1 any](
func SequenceT10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceT2[E, T1, T2 any](
func SequenceT3[E, T1, T2, T3 any](
func SequenceT4[E, T1, T2, T3, T4 any](
func SequenceT5[E, T1, T2, T3, T4, T5 any](
func SequenceT6[E, T1, T2, T3, T4, T5, T6 any](
func SequenceT7[E, T1, T2, T3, T4, T5, T6, T7 any](
func SequenceT8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceT9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceTuple1[E, T1 any](t tuple.Tuple1[IOEither[E, T1]]) IOEither[E, tuple.Tuple1[T1]]
func SequenceTuple10[E, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8], IOEither[E, T9], IOEither[E, T10]]) IOEither[E, tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple2[E, T1, T2 any](t tuple.Tuple2[IOEither[E, T1], IOEither[E, T2]]) IOEither[E, tuple.Tuple2[T1, T2]]
func SequenceTuple3[E, T1, T2, T3 any](t tuple.Tuple3[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3]]) IOEither[E, tuple.Tuple3[T1, T2, T3]]
func SequenceTuple4[E, T1, T2, T3, T4 any](t tuple.Tuple4[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4]]) IOEither[E, tuple.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[E, T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5]]) IOEither[E, tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[E, T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6]]) IOEither[E, tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[E, T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7]]) IOEither[E, tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[E, T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8]]) IOEither[E, tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[E, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOEither[E, T1], IOEither[E, T2], IOEither[E, T3], IOEither[E, T4], IOEither[E, T5], IOEither[E, T6], IOEither[E, T7], IOEither[E, T8], IOEither[E, T9]]) IOEither[E, tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Swap[E, A any](val IOEither[E, A]) IOEither[A, E]
func TryCatch[E, A any](f func() (A, error), onThrow func(error) E) IOEither[E, A]
func TryCatchError[A any](f func() (A, error)) IOEither[error, A]
type Kleisli[E, A, B any] = R.Reader[A, IOEither[E, B]]
func LogJSON[A any](prefix string) Kleisli[error, A, string]
func OrElse[E1, E2, A any](onLeft Kleisli[E2, E1, A]) Kleisli[E2, IOEither[E1, A], A]
func TailRec[E, A, B any](f Kleisli[E, A, tailrec.Trampoline[A, B]]) Kleisli[E, A, B]
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B]
func TraverseArrayPar[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B]
func TraverseArraySeq[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B]
func TraverseArrayWithIndex[E, A, B any](f func(int, A) IOEither[E, B]) Kleisli[E, []A, []B]
func TraverseArrayWithIndexPar[E, A, B any](f func(int, A) IOEither[E, B]) Kleisli[E, []A, []B]
func TraverseArrayWithIndexSeq[E, A, B any](f func(int, A) IOEither[E, B]) Kleisli[E, []A, []B]
func TraverseRecord[K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, map[K]A, map[K]B]
func TraverseRecordPar[K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, map[K]A, map[K]B]
func TraverseRecordSeq[K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, map[K]A, map[K]B]
func TraverseRecordWithIndex[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) Kleisli[E, map[K]A, map[K]B]
func TraverseRecordWithIndexPar[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) Kleisli[E, map[K]A, map[K]B]
func TraverseRecordWithIndexSeq[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) Kleisli[E, map[K]A, map[K]B]
func WithResource[A, E, R, ANY any](onCreate IOEither[E, R], onRelease Kleisli[E, R, ANY]) Kleisli[E, Kleisli[E, R, A], A]
type Monoid[E, A any] = monoid.Monoid[IOEither[E, A]]
func ApplicativeMonoid[E, A any](
func ApplicativeMonoidPar[E, A any](
func ApplicativeMonoidSeq[E, A any](
type Operator[E, A, B any] = Kleisli[E, IOEither[E, A], B]
func After[E, A any](timestamp time.Time) Operator[E, A, A]
func Alt[E, A any](second lazy.Lazy[IOEither[E, A]]) Operator[E, A, A]
func Ap[B, E, A any](ma IOEither[E, A]) Operator[E, func(A) B, B]
func ApFirst[A, E, B any](second IOEither[E, B]) Operator[E, A, A]
func ApPar[B, E, A any](ma IOEither[E, A]) Operator[E, func(A) B, B]
func ApS[E, S1, S2, T any](
func ApSL[E, S, T any](
func ApSecond[A, E, B any](second IOEither[E, B]) Operator[E, A, B]
func Bind[E, S1, S2, T any](
func BindL[E, S, T any](
func BindTo[E, S1, T any](
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B]
func ChainConsumer[E, A any](c Consumer[A]) Operator[E, A, struct{}]
func ChainEitherK[E, A, B any](f either.Kleisli[E, A, B]) Operator[E, A, B]
func ChainFirst[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, A]
func ChainFirstConsumer[E, A any](c Consumer[A]) Operator[E, A, A]
func ChainFirstEitherK[A, E, B any](f either.Kleisli[E, A, B]) Operator[E, A, A]
func ChainFirstIOK[E, A, B any](f io.Kleisli[A, B]) Operator[E, A, A]
func ChainFirstLeft[A, EA, EB, B any](f Kleisli[EB, EA, B]) Operator[EA, A, A]
func ChainIOK[E, A, B any](f io.Kleisli[A, B]) Operator[E, A, B]
func ChainLazyK[E, A, B any](f func(A) lazy.Lazy[B]) Operator[E, A, B]
func ChainTo[A, E, B any](fb IOEither[E, B]) Operator[E, A, B]
func ChainToIO[E, A, B any](fb IO[B]) Operator[E, A, B]
func Delay[E, A any](delay time.Duration) Operator[E, A, A]
func FilterOrElse[E, A any](pred Predicate[A], onFalse func(A) E) Operator[E, A, A]
func Flap[E, B, A any](a A) Operator[E, func(A) B, B]
func Let[E, S1, S2, T any](
func LetL[E, S, T any](
func LetTo[E, S1, S2, T any](
func LetToL[E, S, T any](
func LogEntryExit[E, A any](name string) Operator[E, A, A]
func LogEntryExitF[E, A, STARTTOKEN, ANY any](
func Map[E, A, B any](f func(A) B) Operator[E, A, B]
func MapTo[E, A, B any](b B) Operator[E, A, B]
func Tap[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, A]
func TapEitherK[A, E, B any](f either.Kleisli[E, A, B]) Operator[E, A, A]
func TapIOK[E, A, B any](f io.Kleisli[A, B]) Operator[E, A, A]
func TapLeft[A, EA, EB, B any](f Kleisli[EB, EA, B]) Operator[EA, A, A]
func WithLock[E, A any](lock IO[context.CancelFunc]) Operator[E, A, A]
type Predicate[A any] = predicate.Predicate[A]
type Semigroup[E, A any] = semigroup.Semigroup[IOEither[E, A]]
func AltSemigroup[E, A any]() Semigroup[E, A]
type Trampoline[B, L any] = tailrec.Trampoline[B, L]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/ioresult`

Import: `import "github.com/IBM/fp-go/v2/ioresult"`

IOResult combines IO and Result: `IOResult[A] = IO[Result[A]]` = `func() Result[A]`.

### Exported API

```go
func ApSeq[B, A any](ma IOResult[A]) func(IOResult[func(A) B]) IOResult[B]
func BiMap[E, A, B any](f func(error) E, g func(A) B) func(IOResult[A]) ioeither.IOEither[E, B]
func ChainOptionK[A, B any](onNone func() error) func(O.Kleisli[A, B]) Operator[A, B]
func Eitherize0[F ~func() (R, error), R any](f F) func() IOResult[R]
func Eitherize1[F ~func(T1) (R, error), T1, R any](f F) func(T1) IOResult[R]
func Eitherize10[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) IOResult[R]
func Eitherize2[F ~func(T1, T2) (R, error), T1, T2, R any](f F) func(T1, T2) IOResult[R]
func Eitherize3[F ~func(T1, T2, T3) (R, error), T1, T2, T3, R any](f F) func(T1, T2, T3) IOResult[R]
func Eitherize4[F ~func(T1, T2, T3, T4) (R, error), T1, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) IOResult[R]
func Eitherize5[F ~func(T1, T2, T3, T4, T5) (R, error), T1, T2, T3, T4, T5, R any](f F) func(T1, T2, T3, T4, T5) IOResult[R]
func Eitherize6[F ~func(T1, T2, T3, T4, T5, T6) (R, error), T1, T2, T3, T4, T5, T6, R any](f F) func(T1, T2, T3, T4, T5, T6) IOResult[R]
func Eitherize7[F ~func(T1, T2, T3, T4, T5, T6, T7) (R, error), T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T1, T2, T3, T4, T5, T6, T7) IOResult[R]
func Eitherize8[F ~func(T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8) IOResult[R]
func Eitherize9[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9) IOResult[R]
func Eq[A any](eq EQ.Eq[Result[A]]) EQ.Eq[IOResult[A]]
func Fold[A, B any](onLeft func(error) IO[B], onRight io.Kleisli[A, B]) func(IOResult[A]) IO[B]
func FromIOOption[A any](onNone func() error) func(o IOO.IOOption[A]) IOResult[A]
func FromOption[A any](onNone func() error) func(o O.Option[A]) IOResult[A]
func FromStrictEquals[A comparable]() EQ.Eq[IOResult[A]]
func Functor[A, B any]() functor.Functor[A, B, IOResult[A], IOResult[B]]
func GetOrElse[A any](onLeft func(error) IO[A]) func(IOResult[A]) IO[A]
func GetOrElseOf[A any](onLeft func(error) A) func(IOResult[A]) IO[A]
func MapLeft[A, E any](f func(error) E) func(IOResult[A]) ioeither.IOEither[E, A]
func Monad[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]]
func MonadBiMap[E, A, B any](fa IOResult[A], f func(error) E, g func(A) B) ioeither.IOEither[E, B]
func MonadMapLeft[A, E any](fa IOResult[A], f func(error) E) ioeither.IOEither[E, A]
func Pointed[A any]() pointed.Pointed[A, IOResult[A]]
func Swap[A any](val IOResult[A]) ioeither.IOEither[A, error]
func ToIOOption[A any](ioe IOResult[A]) IOO.IOOption[A]
func TraverseParTuple1[F1 ~func(A1) IOResult[T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOResult[tuple.Tuple1[T1]]
func TraverseParTuple10[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], F10 ~func(A10) IOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseParTuple2[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOResult[tuple.Tuple2[T1, T2]]
func TraverseParTuple3[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseParTuple4[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseParTuple5[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseParTuple6[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseParTuple7[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseParTuple8[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseParTuple9[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseSeqTuple1[F1 ~func(A1) IOResult[T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOResult[tuple.Tuple1[T1]]
func TraverseSeqTuple10[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], F10 ~func(A10) IOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseSeqTuple2[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOResult[tuple.Tuple2[T1, T2]]
func TraverseSeqTuple3[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseSeqTuple4[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseSeqTuple5[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseSeqTuple6[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseSeqTuple7[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseSeqTuple8[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseSeqTuple9[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseTuple1[F1 ~func(A1) IOResult[T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOResult[tuple.Tuple1[T1]]
func TraverseTuple10[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], F10 ~func(A10) IOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple2[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOResult[tuple.Tuple2[T1, T2]]
func TraverseTuple3[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseTuple4[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Uneitherize0[F ~func() IOResult[R], R any](f F) func() (R, error)
func Uneitherize1[F ~func(T1) IOResult[R], T1, R any](f F) func(T1) (R, error)
func Uneitherize10[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) IOResult[R], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error)
func Uneitherize2[F ~func(T1, T2) IOResult[R], T1, T2, R any](f F) func(T1, T2) (R, error)
func Uneitherize3[F ~func(T1, T2, T3) IOResult[R], T1, T2, T3, R any](f F) func(T1, T2, T3) (R, error)
func Uneitherize4[F ~func(T1, T2, T3, T4) IOResult[R], T1, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) (R, error)
func Uneitherize5[F ~func(T1, T2, T3, T4, T5) IOResult[R], T1, T2, T3, T4, T5, R any](f F) func(T1, T2, T3, T4, T5) (R, error)
func Uneitherize6[F ~func(T1, T2, T3, T4, T5, T6) IOResult[R], T1, T2, T3, T4, T5, T6, R any](f F) func(T1, T2, T3, T4, T5, T6) (R, error)
func Uneitherize7[F ~func(T1, T2, T3, T4, T5, T6, T7) IOResult[R], T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T1, T2, T3, T4, T5, T6, T7) (R, error)
func Uneitherize8[F ~func(T1, T2, T3, T4, T5, T6, T7, T8) IOResult[R], T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8) (R, error)
func Uneitherize9[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9) IOResult[R], T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error)
type Consumer[A any] = consumer.Consumer[A]
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[A any] = io.IO[A]
func MonadFold[A, B any](ma IOResult[A], onLeft func(error) IO[B], onRight io.Kleisli[A, B]) IO[B]
type IOResult[A any] = IO[Result[A]]
func Bracket[A, B, ANY any](
func Defer[A any](gen Lazy[IOResult[A]]) IOResult[A]
func Do[S any](
func Flatten[A any](mma IOResult[IOResult[A]]) IOResult[A]
func FromEither[A any](e Result[A]) IOResult[A]
func FromEitherI[A any](a A, err error) IOResult[A]
func FromIO[A any](mr IO[A]) IOResult[A]
func FromImpure[E any](f func()) IOResult[Void]
func FromLazy[A any](mr Lazy[A]) IOResult[A]
func FromResult[A any](e Result[A]) IOResult[A]
func FromResultI[A any](a A, err error) IOResult[A]
func Left[A any](l error) IOResult[A]
func LeftIO[A any](ml IO[error]) IOResult[A]
func Memoize[A any](ma IOResult[A]) IOResult[A]
func MonadAlt[A any](first IOResult[A], second Lazy[IOResult[A]]) IOResult[A]
func MonadAp[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B]
func MonadApFirst[A, B any](first IOResult[A], second IOResult[B]) IOResult[A]
func MonadApPar[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B]
func MonadApSecond[A, B any](first IOResult[A], second IOResult[B]) IOResult[B]
func MonadApSeq[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B]
func MonadChain[A, B any](fa IOResult[A], f Kleisli[A, B]) IOResult[B]
func MonadChainEitherK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[B]
func MonadChainFirst[A, B any](ma IOResult[A], f Kleisli[A, B]) IOResult[A]
func MonadChainFirstEitherK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A]
func MonadChainFirstIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[A]
func MonadChainFirstLeft[A, B any](fa IOResult[A], f Kleisli[error, B]) IOResult[A]
func MonadChainFirstResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A]
func MonadChainI[A, B any](fa IOResult[A], f IOI.Kleisli[A, B]) IOResult[B]
func MonadChainIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[B]
func MonadChainLeft[A any](fa IOResult[A], f Kleisli[error, A]) IOResult[A]
func MonadChainResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[B]
func MonadChainTo[A, B any](fa IOResult[A], fb IOResult[B]) IOResult[B]
func MonadFlap[B, A any](fab IOResult[func(A) B], a A) IOResult[B]
func MonadMap[A, B any](fa IOResult[A], f func(A) B) IOResult[B]
func MonadMapTo[A, B any](fa IOResult[A], b B) IOResult[B]
func MonadOf[A any](r A) IOResult[A]
func MonadTap[A, B any](ma IOResult[A], f Kleisli[A, B]) IOResult[A]
func MonadTapEitherK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A]
func MonadTapIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[A]
func MonadTapLeft[A, B any](fa IOResult[A], f Kleisli[error, B]) IOResult[A]
func MonadTapResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A]
func Of[A any](r A) IOResult[A]
func Retrying[A any](
func Right[A any](r A) IOResult[A]
func RightIO[A any](mr IO[A]) IOResult[A]
func SequenceArray[A any](ma []IOResult[A]) IOResult[[]A]
func SequenceArrayPar[A any](ma []IOResult[A]) IOResult[[]A]
func SequenceArraySeq[A any](ma []IOResult[A]) IOResult[[]A]
func SequenceParT1[T1 any](
func SequenceParT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceParT2[T1, T2 any](
func SequenceParT3[T1, T2, T3 any](
func SequenceParT4[T1, T2, T3, T4 any](
func SequenceParT5[T1, T2, T3, T4, T5 any](
func SequenceParT6[T1, T2, T3, T4, T5, T6 any](
func SequenceParT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceParT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceParT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceParTuple1[T1 any](t tuple.Tuple1[IOResult[T1]]) IOResult[tuple.Tuple1[T1]]
func SequenceParTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9], IOResult[T10]]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceParTuple2[T1, T2 any](t tuple.Tuple2[IOResult[T1], IOResult[T2]]) IOResult[tuple.Tuple2[T1, T2]]
func SequenceParTuple3[T1, T2, T3 any](t tuple.Tuple3[IOResult[T1], IOResult[T2], IOResult[T3]]) IOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceParTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4]]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceParTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5]]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceParTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6]]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceParTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7]]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceParTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8]]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceParTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9]]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceRecord[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A]
func SequenceRecordPar[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A]
func SequenceRecordSeq[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A]
func SequenceSeqT1[T1 any](
func SequenceSeqT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceSeqT2[T1, T2 any](
func SequenceSeqT3[T1, T2, T3 any](
func SequenceSeqT4[T1, T2, T3, T4 any](
func SequenceSeqT5[T1, T2, T3, T4, T5 any](
func SequenceSeqT6[T1, T2, T3, T4, T5, T6 any](
func SequenceSeqT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceSeqT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceSeqT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceSeqTuple1[T1 any](t tuple.Tuple1[IOResult[T1]]) IOResult[tuple.Tuple1[T1]]
func SequenceSeqTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9], IOResult[T10]]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceSeqTuple2[T1, T2 any](t tuple.Tuple2[IOResult[T1], IOResult[T2]]) IOResult[tuple.Tuple2[T1, T2]]
func SequenceSeqTuple3[T1, T2, T3 any](t tuple.Tuple3[IOResult[T1], IOResult[T2], IOResult[T3]]) IOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceSeqTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4]]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceSeqTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5]]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceSeqTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6]]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceSeqTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7]]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceSeqTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8]]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceSeqTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9]]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceT1[T1 any](
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceT2[T1, T2 any](
func SequenceT3[T1, T2, T3 any](
func SequenceT4[T1, T2, T3, T4 any](
func SequenceT5[T1, T2, T3, T4, T5 any](
func SequenceT6[T1, T2, T3, T4, T5, T6 any](
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceTuple1[T1 any](t tuple.Tuple1[IOResult[T1]]) IOResult[tuple.Tuple1[T1]]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9], IOResult[T10]]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple2[T1, T2 any](t tuple.Tuple2[IOResult[T1], IOResult[T2]]) IOResult[tuple.Tuple2[T1, T2]]
func SequenceTuple3[T1, T2, T3 any](t tuple.Tuple3[IOResult[T1], IOResult[T2], IOResult[T3]]) IOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4]]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5]]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6]]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7]]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8]]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9]]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TryCatch[A any](f func() (A, error), onThrow Endomorphism[error]) IOResult[A]
func TryCatchError[A any](f func() (A, error)) IOResult[A]
type Kleisli[A, B any] = reader.Reader[A, IOResult[B]]
func LogJSON[A any](prefix string) Kleisli[A, string]
func TailRec[A, B any](f Kleisli[A, tailrec.Trampoline[A, B]]) Kleisli[A, B]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayPar[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexPar[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordPar[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexPar[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexSeq[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B]
func WithResource[A, R, ANY any](onCreate IOResult[R], onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A]
type Lazy[A any] = lazy.Lazy[A]
type Monoid[A any] = monoid.Monoid[IOResult[A]]
func ApplicativeMonoid[A any](
func ApplicativeMonoidPar[A any](
func ApplicativeMonoidSeq[A any](
type Operator[A, B any] = Kleisli[IOResult[A], B]
func After[A any](timestamp time.Time) Operator[A, A]
func Alt[A any](second Lazy[IOResult[A]]) Operator[A, A]
func Ap[B, A any](ma IOResult[A]) Operator[func(A) B, B]
func ApFirst[A, B any](second IOResult[B]) Operator[A, A]
func ApPar[B, A any](ma IOResult[A]) Operator[func(A) B, B]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func ApSecond[A, B any](second IOResult[B]) Operator[A, B]
func Bind[S1, S2, T any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainConsumer[A any](c Consumer[A]) Operator[A, struct{}]
func ChainEitherK[A, B any](f result.Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainFirstConsumer[A any](c Consumer[A]) Operator[A, A]
func ChainFirstEitherK[A, B any](f result.Kleisli[A, B]) Operator[A, A]
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A]
func ChainI[A, B any](f IOI.Kleisli[A, B]) Operator[A, B]
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B]
func ChainLazyK[A, B any](f func(A) Lazy[B]) Operator[A, B]
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A]
func ChainResultK[A, B any](f result.Kleisli[A, B]) Operator[A, B]
func ChainTo[A, B any](fb IOResult[B]) Operator[A, B]
func Delay[A any](delay time.Duration) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func LogEntryExit[A any](name string) Operator[A, A]
func LogEntryExitF[A, STARTTOKEN, ANY any](
func Map[A, B any](f func(A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A]
func Tap[A, B any](f Kleisli[A, B]) Operator[A, A]
func TapEitherK[A, B any](f result.Kleisli[A, B]) Operator[A, A]
func TapIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func TapLeft[A, B any](f Kleisli[error, B]) Operator[A, A]
func WithLock[A any](lock IO[context.CancelFunc]) Operator[A, A]
type Predicate[A any] = predicate.Predicate[A]
type Result[A any] = result.Result[A]
type Semigroup[A any] = semigroup.Semigroup[IOResult[A]]
func AltSemigroup[A any]() Semigroup[A]
type Void = function.Void
```

---

# Reader Stack

## package `github.com/IBM/fp-go/v2/reader`

Import: `import "github.com/IBM/fp-go/v2/reader"`

Reader represents a computation depending on environment R.

Key types:
- `Reader[R, A] = func(R) A`
- `Kleisli[R, A, B] = func(A) Reader[R, B]`

### Exported API

```go
func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[Reader[R, A]]
func Curry2[R, T1, T2, A any](f func(R, T1, T2) A) func(T1) func(T2) Reader[R, A]
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) A) func(T1) func(T2) func(T3) Reader[R, A]
func Curry4[R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) A) func(T1) func(T2) func(T3) func(T4) Reader[R, A]
func From0[F ~func(C) R, C, R any](f F) func() Reader[C, R]
func From1[F ~func(C, T0) R, T0, C, R any](f F) func(T0) Reader[C, R]
func From10[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) R, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) Reader[C, R]
func From2[F ~func(C, T0, T1) R, T0, T1, C, R any](f F) func(T0, T1) Reader[C, R]
func From3[F ~func(C, T0, T1, T2) R, T0, T1, T2, C, R any](f F) func(T0, T1, T2) Reader[C, R]
func From4[F ~func(C, T0, T1, T2, T3) R, T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) Reader[C, R]
func From5[F ~func(C, T0, T1, T2, T3, T4) R, T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) Reader[C, R]
func From6[F ~func(C, T0, T1, T2, T3, T4, T5) R, T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) Reader[C, R]
func From7[F ~func(C, T0, T1, T2, T3, T4, T5, T6) R, T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) Reader[C, R]
func From8[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) R, T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) Reader[C, R]
func From9[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) R, T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) Reader[C, R]
func Read[A, E any](e E) func(Reader[E, A]) A
func Traverse[R2, R1, A, B any](
func TraverseRecord[K comparable, R, A, B any](f Kleisli[R, A, B]) func(map[K]A) Reader[R, map[K]B]
func TraverseRecordWithIndex[K comparable, R, A, B any](f func(K, A) Reader[R, B]) func(map[K]A) Reader[R, map[K]B]
func Uncurry0[R, A any](f Reader[R, A]) func(R) A
func Uncurry1[R, T1, A any](f Kleisli[R, T1, A]) func(R, T1) A
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) Reader[R, A]) func(R, T1, T2) A
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) Reader[R, A]) func(R, T1, T2, T3) A
func Uncurry4[R, T1, T2, T3, T4, A any](f func(T1) func(T2) func(T3) func(T4) Reader[R, A]) func(R, T1, T2, T3, T4) A
type Kleisli[R, A, B any] = func(A) Reader[R, B]
func Compose[C, R, B any](ab Reader[R, B]) Kleisli[R, Reader[B, C], C]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, Reader[R1, A], A]
func Curry[R, T1, A any](f func(R, T1) A) Kleisli[R, T1, A]
func Curry1[R, T1, A any](f func(R, T1) A) Kleisli[R, T1, A]
func Local[A, R1, R2 any](f func(R2) R1) Kleisli[R2, Reader[R1, A], A]
func Promap[E, A, D, B any](f func(D) E, g func(A) B) Kleisli[D, Reader[E, A], B]
func ReduceArray[R, A, B any](reduce func(B, A) B, initial B) Kleisli[R, []Reader[R, A], B]
func ReduceArrayM[R, A any](m monoid.Monoid[A]) Kleisli[R, []Reader[R, A], A]
func Sequence[R1, R2, A any](ma Kleisli[R1, R2, A]) Kleisli[R2, R1, A]
func TailRec[R, A, B any](f Kleisli[R, A, Trampoline[A, B]]) Kleisli[R, A, B]
func TraverseArray[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, []A, []B]
func TraverseArrayWithIndex[R, A, B any](f func(int, A) Reader[R, B]) Kleisli[R, []A, []B]
func TraverseIter[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, Seq[A], Seq[B]]
func TraverseReduceArray[R, A, B, C any](trfrm Kleisli[R, A, B], reduce func(C, B) C, initial C) Kleisli[R, []A, C]
func TraverseReduceArrayM[R, A, B any](trfrm Kleisli[R, A, B], m monoid.Monoid[B]) Kleisli[R, []A, B]
type Lazy[A any] = func() A
type Operator[R, A, B any] = Kleisli[R, Reader[R, A], B]
func Ap[B, R, A any](fa Reader[R, A]) Operator[R, func(A) B, B]
func ApS[R, S1, S2, T any](
func ApSL[R, S, T any](
func Bind[R, S1, S2, T any](
func BindL[R, S, T any](
func BindTo[R, S1, T any](
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainTo[A, R, B any](b Reader[R, B]) Operator[R, A, B]
func Flap[R, B, A any](a A) Operator[R, func(A) B, B]
func Let[R, S1, S2, T any](
func LetL[R, S, T any](
func LetTo[R, S1, S2, T any](
func LetToL[R, S, T any](
func Map[E, A, B any](f func(A) B) Operator[E, A, B]
func MapTo[E, A, B any](b B) Operator[E, A, B]
type Reader[R, A any] = func(R) A
func Ask[R any]() Reader[R, R]
func Asks[R, A any](f Reader[R, A]) Reader[R, A]
func AsksReader[R, A any](f Kleisli[R, R, A]) Reader[R, A]
func Bracket[
func Curry0[R, A any](f func(R) A) Reader[R, A]
func Do[R, S any](
func First[A, B, C any](pab Reader[A, B]) Reader[T.Tuple2[A, C], T.Tuple2[B, C]]
func Flatten[R, A any](mma Reader[R, Reader[R, A]]) Reader[R, A]
func MonadAp[B, R, A any](fab Reader[R, func(A) B], fa Reader[R, A]) Reader[R, B]
func MonadChain[R, A, B any](ma Reader[R, A], f Kleisli[R, A, B]) Reader[R, B]
func MonadChainTo[A, R, B any](_ Reader[R, A], b Reader[R, B]) Reader[R, B]
func MonadFlap[R, B, A any](fab Reader[R, func(A) B], a A) Reader[R, B]
func MonadMap[E, A, B any](fa Reader[E, A], f func(A) B) Reader[E, B]
func MonadMapTo[E, A, B any](_ Reader[E, A], b B) Reader[E, B]
func MonadReduceArray[R, A, B any](as []Reader[R, A], reduce func(B, A) B, initial B) Reader[R, B]
func MonadReduceArrayM[R, A any](as []Reader[R, A], m monoid.Monoid[A]) Reader[R, A]
func MonadTraverseArray[R, A, B any](ma []A, f Kleisli[R, A, B]) Reader[R, []B]
func MonadTraverseRecord[K comparable, R, A, B any](ma map[K]A, f Kleisli[R, A, B]) Reader[R, map[K]B]
func MonadTraverseRecordWithIndex[K comparable, R, A, B any](ma map[K]A, f func(K, A) Reader[R, B]) Reader[R, map[K]B]
func MonadTraverseReduceArray[R, A, B, C any](as []A, trfrm Kleisli[R, A, B], reduce func(C, B) C, initial C) Reader[R, C]
func MonadTraverseReduceArrayM[R, A, B any](as []A, trfrm Kleisli[R, A, B], m monoid.Monoid[B]) Reader[R, B]
func Of[R, A any](a A) Reader[R, A]
func OfLazy[R, A any](fa Lazy[A]) Reader[R, A]
func Second[A, B, C any](pbc Reader[B, C]) Reader[T.Tuple2[A, B], T.Tuple2[A, C]]
func SequenceArray[R, A any](ma []Reader[R, A]) Reader[R, []A]
func SequenceIter[R, A any](as Seq[Reader[R, A]]) Reader[R, Seq[A]]
func SequenceRecord[K comparable, R, A any](ma map[K]Reader[R, A]) Reader[R, map[K]A]
func SequenceT1[R, A any](a Reader[R, A]) Reader[R, T.Tuple1[A]]
func SequenceT2[R, A, B any](a Reader[R, A], b Reader[R, B]) Reader[R, T.Tuple2[A, B]]
func SequenceT3[R, A, B, C any](a Reader[R, A], b Reader[R, B], c Reader[R, C]) Reader[R, T.Tuple3[A, B, C]]
func SequenceT4[R, A, B, C, D any](a Reader[R, A], b Reader[R, B], c Reader[R, C], d Reader[R, D]) Reader[R, T.Tuple4[A, B, C, D]]
func WithLocal[A, R1, R2 any](fa Reader[R1, A], f func(R2) R1) Reader[R2, A]
type Seq[T any] = iter.Seq[T]
type Trampoline[B, L any] = tailrec.Trampoline[B, L]
```

## package `github.com/IBM/fp-go/v2/readeroption`

Import: `import "github.com/IBM/fp-go/v2/readeroption"`

ReaderOption: `ReaderOption[R, A] = func(R) Option[A]`.

### Exported API

```go
func AlternativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderOption[R, A]]
func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderOption[R, A]]
func Curry2[R, T1, T2, A any](f func(R, T1, T2) (A, bool)) func(T1) func(T2) ReaderOption[R, A]
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, bool)) func(T1) func(T2) func(T3) ReaderOption[R, A]
func Fold[E, A, B any](onNone Reader[E, B], onRight reader.Kleisli[E, A, B]) reader.Operator[E, Option[A], B]
func From0[R, A any](f func(R) (A, bool)) func() ReaderOption[R, A]
func From2[R, T1, T2, A any](f func(R, T1, T2) (A, bool)) func(T1, T2) ReaderOption[R, A]
func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, bool)) func(T1, T2, T3) ReaderOption[R, A]
func GetOrElse[E, A any](onNone Reader[E, A]) reader.Operator[E, Option[A], A]
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderOption[R1, A]) ReaderOption[R2, A]
func Read[A, E any](e E) func(ReaderOption[E, A]) Option[A]
func ReadOption[A, E any](e Option[E]) func(ReaderOption[E, A]) Option[A]
func Sequence[R1, R2, A any](ma ReaderOption[R2, ReaderOption[R1, A]]) reader.Kleisli[R2, R1, Option[A]]
func SequenceReader[R1, R2, A any](ma ReaderOption[R2, Reader[R1, A]]) reader.Kleisli[R2, R1, Option[A]]
func Traverse[R2, R1, A, B any](
func TraverseArrayWithIndex[E, A, B any](f func(int, A) ReaderOption[E, B]) func([]A) ReaderOption[E, []B]
func TraverseReader[R2, R1, A, B any](
func Uncurry1[R, T1, A any](f func(T1) ReaderOption[R, A]) func(R, T1) (A, bool)
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) ReaderOption[R, A]) func(R, T1, T2) (A, bool)
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderOption[R, A]) func(R, T1, T2, T3) (A, bool)
type Either[E, A any] = either.Either[E, A]
type Kleisli[R, A, B any] = Reader[A, ReaderOption[R, B]]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderOption[R1, A], A]
func Curry1[R, T1, A any](f func(R, T1) (A, bool)) Kleisli[R, T1, A]
func From1[R, T1, A any](f func(R, T1) (A, bool)) Kleisli[R, T1, A]
func FromPredicate[E, A any](pred Predicate[A]) Kleisli[E, A, A]
func Promap[R, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, ReaderOption[R, A], B]
func TailRec[R, A, B any](f Kleisli[R, A, tailrec.Trampoline[A, B]]) Kleisli[R, A, B]
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B]
type Lazy[A any] = lazy.Lazy[A]
type Operator[R, A, B any] = Reader[ReaderOption[R, A], ReaderOption[R, B]]
func Alt[E, A any](second Lazy[ReaderOption[E, A]]) Operator[E, A, A]
func Ap[B, E, A any](fa ReaderOption[E, A]) Operator[E, func(A) B, B]
func ApS[R, S1, S2, T any](
func ApSL[R, S, T any](
func Bind[R, S1, S2, T any](
func BindL[R, S, T any](
func BindTo[R, S1, T any](
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B]
func ChainOptionK[E, A, B any](f O.Kleisli[A, B]) Operator[E, A, B]
func Flap[E, B, A any](a A) Operator[E, func(A) B, B]
func Let[R, S1, S2, T any](
func LetL[R, S, T any](
func LetTo[R, S1, S2, T any](
func LetToL[R, S, T any](
func Map[E, A, B any](f func(A) B) Operator[E, A, B]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
func MonadFold[E, A, B any](fa ReaderOption[E, A], onNone Reader[E, B], onRight reader.Kleisli[E, A, B]) Reader[E, B]
type ReaderOption[R, A any] = Reader[R, Option[A]]
func Ask[E any]() ReaderOption[E, E]
func Asks[E, A any](r Reader[E, A]) ReaderOption[E, A]
func Curry0[R, A any](f func(R) (A, bool)) ReaderOption[R, A]
func Do[R, S any](
func Flatten[E, A any](mma ReaderOption[E, ReaderOption[E, A]]) ReaderOption[E, A]
func FromOption[E, A any](e Option[A]) ReaderOption[E, A]
func FromReader[E, A any](r Reader[E, A]) ReaderOption[E, A]
func MonadAlt[E, A any](first ReaderOption[E, A], second Lazy[ReaderOption[E, A]]) ReaderOption[E, A]
func MonadAp[E, A, B any](fab ReaderOption[E, func(A) B], fa ReaderOption[E, A]) ReaderOption[E, B]
func MonadChain[E, A, B any](ma ReaderOption[E, A], f Kleisli[E, A, B]) ReaderOption[E, B]
func MonadChainOptionK[E, A, B any](ma ReaderOption[E, A], f O.Kleisli[A, B]) ReaderOption[E, B]
func MonadFlap[E, A, B any](fab ReaderOption[E, func(A) B], a A) ReaderOption[E, B]
func MonadMap[E, A, B any](fa ReaderOption[E, A], f func(A) B) ReaderOption[E, B]
func None[E, A any]() ReaderOption[E, A]
func Of[E, A any](a A) ReaderOption[E, A]
func SequenceArray[E, A any](ma []ReaderOption[E, A]) ReaderOption[E, []A]
func SequenceT1[E, A any](a ReaderOption[E, A]) ReaderOption[E, T.Tuple1[A]]
func SequenceT2[E, A, B any](
func SequenceT3[E, A, B, C any](
func SequenceT4[E, A, B, C, D any](
func Some[E, A any](r A) ReaderOption[E, A]
func SomeReader[E, A any](r Reader[E, A]) ReaderOption[E, A]
```

## package `github.com/IBM/fp-go/v2/readereither`

Import: `import "github.com/IBM/fp-go/v2/readereither"`

ReaderEither: `ReaderEither[R, E, A] = func(R) Either[E, A]`.

### Exported API

```go
func AltMonoid[R, E, A any](zero lazy.Lazy[ReaderEither[R, E, A]]) monoid.Monoid[ReaderEither[R, E, A]]
func AlternativeMonoid[R, E, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderEither[R, E, A]]
func Ap[B, E, L, A any](fa ReaderEither[E, L, A]) func(ReaderEither[E, L, func(A) B]) ReaderEither[E, L, B]
func ApS[R, E, S1, S2, T any](
func ApSL[R, E, S, T any](
func ApplicativeMonoid[R, E, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderEither[R, E, A]]
func BiMap[E, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderEither[E, E1, A]) ReaderEither[E, E2, B]
func Bind[R, E, S1, S2, T any](
func BindEitherK[R, E, S1, S2, T any](
func BindL[R, E, S, T any](
func BindReaderK[R, E, S1, S2, T any](
func BindTo[R, E, S1, T any](
func BindToEither[
func BindToReader[
func Chain[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B]
func ChainEitherK[E, L, A, B any](f func(A) Either[L, B]) func(ma ReaderEither[E, L, A]) ReaderEither[E, L, B]
func ChainLeft[R, EA, EB, A any](f Kleisli[R, EB, EA, A]) func(ReaderEither[R, EA, A]) ReaderEither[R, EB, A]
func ChainOptionK[E, A, B, L any](onNone func() L) func(func(A) Option[B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B]
func ChainReaderK[L, E, A, B any](f reader.Kleisli[E, A, B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B]
func Curry1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderEither[R, error, A]
func Curry2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1) func(T2) ReaderEither[R, error, A]
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderEither[R, error, A]
func Flap[L, E, B, A any](a A) func(ReaderEither[L, E, func(A) B]) ReaderEither[L, E, B]
func Fold[E, L, A, B any](onLeft func(L) Reader[E, B], onRight func(A) Reader[E, B]) func(ReaderEither[E, L, A]) Reader[E, B]
func From0[R, A any](f func(R) (A, error)) func() ReaderEither[R, error, A]
func From1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderEither[R, error, A]
func From2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) ReaderEither[R, error, A]
func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderEither[R, error, A]
func FromPredicate[E, L, A any](pred func(A) bool, onFalse func(A) L) func(A) ReaderEither[E, L, A]
func GetOrElse[E, L, A any](onLeft func(L) Reader[E, A]) func(ReaderEither[E, L, A]) Reader[E, A]
func Let[R, E, S1, S2, T any](
func LetL[R, E, S, T any](
func LetTo[R, E, S1, S2, T any](
func LetToL[R, E, S, T any](
func Local[E, A, R1, R2 any](f func(R2) R1) func(ReaderEither[R1, E, A]) ReaderEither[R2, E, A]
func Map[E, L, A, B any](f func(A) B) func(ReaderEither[E, L, A]) ReaderEither[E, L, B]
func MapLeft[C, E1, E2, A any](f func(E1) E2) func(ReaderEither[C, E1, A]) ReaderEither[C, E2, A]
func OrLeft[A, L1, E, L2 any](onLeft func(L1) Reader[E, L2]) func(ReaderEither[E, L1, A]) ReaderEither[E, L2, A]
func Read[E1, A, E any](e E) func(ReaderEither[E, E1, A]) Either[E1, A]
func ReadEither[E1, A, E any](e Either[E1, E]) func(ReaderEither[E, E1, A]) Either[E1, A]
func Traverse[R2, R1, E, A, B any](
func TraverseArray[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func([]A) ReaderEither[E, L, []B]
func TraverseArrayWithIndex[E, L, A, B any](f func(int, A) ReaderEither[E, L, B]) func([]A) ReaderEither[E, L, []B]
func TraverseReader[R2, R1, E, A, B any](
func Uncurry1[R, T1, A any](f func(T1) ReaderEither[R, error, A]) func(R, T1) (A, error)
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) ReaderEither[R, error, A]) func(R, T1, T2) (A, error)
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderEither[R, error, A]) func(R, T1, T2, T3) (A, error)
type Either[E, A any] = either.Either[E, A]
type Kleisli[R, E, A, B any] = Reader[A, ReaderEither[R, E, B]]
func Contramap[E, A, R1, R2 any](f func(R2) R1) Kleisli[R2, E, ReaderEither[R1, E, A], A]
func OrElse[R, E1, E2, A any](onLeft Kleisli[R, E2, E1, A]) Kleisli[R, E2, ReaderEither[R, E1, A], A]
func Promap[R, E, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, E, ReaderEither[R, E, A], B]
func Sequence[R1, R2, E, A any](ma ReaderEither[R2, E, ReaderEither[R1, E, A]]) Kleisli[R2, E, R1, A]
func SequenceReader[R1, R2, E, A any](ma ReaderEither[R2, E, Reader[R1, A]]) Kleisli[R2, E, R1, A]
func TailRec[R, E, A, B any](f Kleisli[R, E, A, tailrec.Trampoline[A, B]]) Kleisli[R, E, A, B]
type Lazy[A any] = lazy.Lazy[A]
type Operator[R, E, A, B any] = Kleisli[R, E, ReaderEither[R, E, A], B]
func Alt[R, E, A any](second Lazy[ReaderEither[R, E, A]]) Operator[R, E, A, A]
func ApEitherS[
func ApReaderS[
func ChainFirstLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A]
func TapLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A]
type Option[A any] = option.Option[A]
type Reader[R, A any] = reader.Reader[R, A]
func MonadFold[E, L, A, B any](ma ReaderEither[E, L, A], onLeft func(L) Reader[E, B], onRight func(A) Reader[E, B]) Reader[E, B]
type ReaderEither[R, E, A any] = Reader[R, Either[E, A]]
func Ask[E, L any]() ReaderEither[E, L, E]
func Asks[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A]
func Curry0[R, A any](f func(R) (A, error)) ReaderEither[R, error, A]
func Do[R, E, S any](
func Flatten[E, L, A any](mma ReaderEither[E, L, ReaderEither[E, L, A]]) ReaderEither[E, L, A]
func FromEither[E, L, A any](e Either[L, A]) ReaderEither[E, L, A]
func FromReader[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A]
func Left[E, A, L any](l L) ReaderEither[E, L, A]
func LeftReader[A, E, L any](l Reader[E, L]) ReaderEither[E, L, A]
func MonadAlt[R, E, A any](first ReaderEither[R, E, A], second Lazy[ReaderEither[R, E, A]]) ReaderEither[R, E, A]
func MonadAp[B, E, L, A any](fab ReaderEither[E, L, func(A) B], fa ReaderEither[E, L, A]) ReaderEither[E, L, B]
func MonadBiMap[E, E1, E2, A, B any](fa ReaderEither[E, E1, A], f func(E1) E2, g func(A) B) ReaderEither[E, E2, B]
func MonadChain[E, L, A, B any](ma ReaderEither[E, L, A], f func(A) ReaderEither[E, L, B]) ReaderEither[E, L, B]
func MonadChainEitherK[E, L, A, B any](ma ReaderEither[E, L, A], f func(A) Either[L, B]) ReaderEither[E, L, B]
func MonadChainFirstLeft[A, R, EA, EB, B any](ma ReaderEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderEither[R, EA, A]
func MonadChainLeft[R, EA, EB, A any](fa ReaderEither[R, EA, A], f Kleisli[R, EB, EA, A]) ReaderEither[R, EB, A]
func MonadChainReaderK[L, E, A, B any](ma ReaderEither[E, L, A], f reader.Kleisli[E, A, B]) ReaderEither[E, L, B]
func MonadFlap[L, E, A, B any](fab ReaderEither[L, E, func(A) B], a A) ReaderEither[L, E, B]
func MonadMap[E, L, A, B any](fa ReaderEither[E, L, A], f func(A) B) ReaderEither[E, L, B]
func MonadMapLeft[C, E1, E2, A any](fa ReaderEither[C, E1, A], f func(E1) E2) ReaderEither[C, E2, A]
func MonadTapLeft[A, R, EA, EB, B any](ma ReaderEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderEither[R, EA, A]
func Of[E, L, A any](a A) ReaderEither[E, L, A]
func OfLazy[E, L, A any](r Lazy[A]) ReaderEither[E, L, A]
func Right[E, L, A any](r A) ReaderEither[E, L, A]
func RightReader[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A]
func SequenceArray[E, L, A any](ma []ReaderEither[E, L, A]) ReaderEither[E, L, []A]
func SequenceT1[L, E, A any](a ReaderEither[E, L, A]) ReaderEither[E, L, T.Tuple1[A]]
func SequenceT2[L, E, A, B any](
func SequenceT3[L, E, A, B, C any](
func SequenceT4[L, E, A, B, C, D any](
```

## package `github.com/IBM/fp-go/v2/readerio`

Import: `import "github.com/IBM/fp-go/v2/readerio"`

ReaderIO: `ReaderIO[R, A] = func(R) func() A`.

### Exported API

```go
func ApFirst[A, R, B any](second ReaderIO[R, B]) func(ReaderIO[R, A]) ReaderIO[R, A]
func ApS[R, S1, S2, T any](
func ApSL[R, S, T any](
func ApSecond[A, R, B any](second ReaderIO[R, B]) func(ReaderIO[R, A]) ReaderIO[R, B]
func Bind[R, S1, S2, T any](
func BindL[R, S, T any](
func BindTo[R, S1, T any](
func Eq[R, A any](e EQ.Eq[A]) func(r R) EQ.Eq[ReaderIO[R, A]]
func From0[F ~func(R) IO[A], R, A any](f func(R) IO[A]) func() ReaderIO[R, A]
func From1[F ~func(R, T1) IO[A], R, T1, A any](f func(R, T1) IO[A]) func(T1) ReaderIO[R, A]
func From2[F ~func(R, T1, T2) IO[A], R, T1, T2, A any](f func(R, T1, T2) IO[A]) func(T1, T2) ReaderIO[R, A]
func From3[F ~func(R, T1, T2, T3) IO[A], R, T1, T2, T3, A any](f func(R, T1, T2, T3) IO[A]) func(T1, T2, T3) ReaderIO[R, A]
func Let[R, S1, S2, T any](
func LetL[R, S, T any](
func LetTo[R, S1, S2, T any](
func LetToL[R, S, T any](
func Read[A, R any](r R) func(ReaderIO[R, A]) IO[A]
func ReadIO[A, R any](r IO[R]) func(ReaderIO[R, A]) IO[A]
func Traverse[R2, R1, A, B any](
func TraverseArray[R, A, B any](f func(A) ReaderIO[R, B]) func([]A) ReaderIO[R, []B]
func TraverseArrayWithIndex[R, A, B any](f func(int, A) ReaderIO[R, B]) func([]A) ReaderIO[R, []B]
func TraverseReader[R2, R1, A, B any](
type Consumer[A any] = consumer.Consumer[A]
type Either[E, A any] = either.Either[E, A]
type IO[A any] = io.IO[A]
type Kleisli[R, A, B any] = Reader[A, ReaderIO[R, B]]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIO[R1, A], A]
func Local[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIO[R1, A], A]
func LocalIOK[A, R1, R2 any](f io.Kleisli[R2, R1]) Kleisli[R2, ReaderIO[R1, A], A]
func LogGo[R, A any](prefix string) Kleisli[R, A, A]
func Logf[R, A any](prefix string) Kleisli[R, A, A]
func PrintGo[R, A any](prefix string) Kleisli[R, A, A]
func Printf[R, A any](prefix string) Kleisli[R, A, A]
func Promap[E, A, D, B any](f func(D) E, g func(A) B) Kleisli[D, ReaderIO[E, A], B]
func Sequence[R1, R2, A any](ma ReaderIO[R2, ReaderIO[R1, A]]) Kleisli[R2, R1, A]
func SequenceReader[R1, R2, A any](ma ReaderIO[R2, Reader[R1, A]]) Kleisli[R2, R1, A]
func TailRec[R, A, B any](f Kleisli[R, A, Trampoline[A, B]]) Kleisli[R, A, B]
func WithResource[R, A, B, ANY any](
type Operator[R, A, B any] = Kleisli[R, ReaderIO[R, A], B]
func After[R, A any](timestamp time.Time) Operator[R, A, A]
func Ap[B, R, A any](fa ReaderIO[R, A]) Operator[R, func(A) B, B]
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainConsumer[R, A any](c Consumer[A]) Operator[R, A, Void]
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A]
func ChainFirstConsumer[R, A any](c Consumer[A]) Operator[R, A, A]
func ChainFirstIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A]
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A]
func ChainIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, B]
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B]
func Delay[R, A any](delay time.Duration) Operator[R, A, A]
func Flap[R, B, A any](a A) Operator[R, func(A) B, B]
func Map[R, A, B any](f func(A) B) Operator[R, A, B]
func MapTo[R, A, B any](b B) Operator[R, A, B]
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A]
func TapIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A]
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A]
func WithLock[R, A any](lock func() context.CancelFunc) Operator[R, A, A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderIO[R, A any] = Reader[R, IO[A]]
func Ask[R any]() ReaderIO[R, R]
func Asks[R, A any](r Reader[R, A]) ReaderIO[R, A]
func Bracket[
func Defer[R, A any](gen func() ReaderIO[R, A]) ReaderIO[R, A]
func Do[R, S any](
func Flatten[R, A any](mma ReaderIO[R, ReaderIO[R, A]]) ReaderIO[R, A]
func FromIO[R, A any](t IO[A]) ReaderIO[R, A]
func FromReader[R, A any](r Reader[R, A]) ReaderIO[R, A]
func Memoize[R, A any](rdr ReaderIO[R, A]) ReaderIO[R, A]
func MonadAp[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B]
func MonadApFirst[A, R, B any](first ReaderIO[R, A], second ReaderIO[R, B]) ReaderIO[R, A]
func MonadApPar[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B]
func MonadApSecond[A, R, B any](first ReaderIO[R, A], second ReaderIO[R, B]) ReaderIO[R, B]
func MonadApSeq[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B]
func MonadChain[R, A, B any](ma ReaderIO[R, A], f Kleisli[R, A, B]) ReaderIO[R, B]
func MonadChainFirst[R, A, B any](ma ReaderIO[R, A], f Kleisli[R, A, B]) ReaderIO[R, A]
func MonadChainFirstIOK[R, A, B any](ma ReaderIO[R, A], f io.Kleisli[A, B]) ReaderIO[R, A]
func MonadChainFirstReaderK[R, A, B any](ma ReaderIO[R, A], f reader.Kleisli[R, A, B]) ReaderIO[R, A]
func MonadChainIOK[R, A, B any](ma ReaderIO[R, A], f io.Kleisli[A, B]) ReaderIO[R, B]
func MonadChainReaderK[R, A, B any](ma ReaderIO[R, A], f reader.Kleisli[R, A, B]) ReaderIO[R, B]
func MonadFlap[R, B, A any](fab ReaderIO[R, func(A) B], a A) ReaderIO[R, B]
func MonadMap[R, A, B any](fa ReaderIO[R, A], f func(A) B) ReaderIO[R, B]
func MonadMapTo[R, A, B any](fa ReaderIO[R, A], b B) ReaderIO[R, B]
func MonadTap[R, A, B any](ma ReaderIO[R, A], f Kleisli[R, A, B]) ReaderIO[R, A]
func MonadTapIOK[R, A, B any](ma ReaderIO[R, A], f io.Kleisli[A, B]) ReaderIO[R, A]
func MonadTapReaderK[R, A, B any](ma ReaderIO[R, A], f reader.Kleisli[R, A, B]) ReaderIO[R, A]
func Of[R, A any](a A) ReaderIO[R, A]
func Retrying[R, A any](
func SequenceArray[R, A any](ma []ReaderIO[R, A]) ReaderIO[R, []A]
func SequenceT1[R, A any](a ReaderIO[R, A]) ReaderIO[R, T.Tuple1[A]]
func SequenceT2[R, A, B any](a ReaderIO[R, A], b ReaderIO[R, B]) ReaderIO[R, T.Tuple2[A, B]]
func SequenceT3[R, A, B, C any](a ReaderIO[R, A], b ReaderIO[R, B], c ReaderIO[R, C]) ReaderIO[R, T.Tuple3[A, B, C]]
func SequenceT4[R, A, B, C, D any](a ReaderIO[R, A], b ReaderIO[R, B], c ReaderIO[R, C], d ReaderIO[R, D]) ReaderIO[R, T.Tuple4[A, B, C, D]]
type Trampoline[B, L any] = tailrec.Trampoline[B, L]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/readeriooption`

Import: `import "github.com/IBM/fp-go/v2/readeriooption"`

ReaderIOOption: `ReaderIOOption[R, A] = func(R) func() Option[A]`.

### Exported API

```go
func AlternativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderIOOption[R, A]]
func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderIOOption[R, A]]
func Fold[R, A, B any](onNone Reader[R, B], onRight reader.Kleisli[R, A, B]) reader.Operator[R, Option[A], B]
func GetOrElse[R, A any](onNone Reader[R, A]) reader.Operator[R, Option[A], A]
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderIOOption[R1, A]) ReaderIOOption[R2, A]
func Read[A, R any](e R) func(ReaderIOOption[R, A]) IOOption[A]
func TraverseArrayWithIndex[E, A, B any](f func(int, A) ReaderIOOption[E, B]) func([]A) ReaderIOOption[E, []B]
type Either[E, A any] = either.Either[E, A]
type IOOption[A any] = iooption.IOOption[A]
type Kleisli[R, A, B any] = Reader[A, ReaderIOOption[R, B]]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIOOption[R1, A], A]
func FromPredicate[R, A any](pred Predicate[A]) Kleisli[R, A, A]
func Promap[R, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, ReaderIOOption[R, A], B]
func TailRec[R, A, B any](f Kleisli[R, A, tailrec.Trampoline[A, B]]) Kleisli[R, A, B]
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B]
type Lazy[A any] = lazy.Lazy[A]
type Operator[R, A, B any] = Reader[ReaderIOOption[R, A], ReaderIOOption[R, B]]
func Alt[R, A any](second Lazy[ReaderIOOption[R, A]]) Operator[R, A, A]
func Ap[B, R, A any](fa ReaderIOOption[R, A]) Operator[R, func(A) B, B]
func ApS[R, S1, S2, T any](
func ApSL[R, S, T any](
func Bind[R, S1, S2, T any](
func BindL[R, S, T any](
func BindTo[R, S1, T any](
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainOptionK[R, A, B any](f O.Kleisli[A, B]) Operator[R, A, B]
func Flap[R, B, A any](a A) Operator[R, func(A) B, B]
func Let[R, S1, S2, T any](
func LetL[R, S, T any](
func LetTo[R, S1, S2, T any](
func LetToL[R, S, T any](
func Map[R, A, B any](f func(A) B) Operator[R, A, B]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderIO[R, A any] = readerio.ReaderIO[R, A]
func MonadFold[R, A, B any](fa ReaderIOOption[R, A], onNone ReaderIO[R, B], onRight readerio.Kleisli[R, A, B]) ReaderIO[R, B]
type ReaderIOOption[R, A any] = Reader[R, IOOption[A]]
func Ask[R any]() ReaderIOOption[R, R]
func Asks[R, A any](r Reader[R, A]) ReaderIOOption[R, A]
func Do[R, S any](
func Flatten[R, A any](mma ReaderIOOption[R, ReaderIOOption[R, A]]) ReaderIOOption[R, A]
func FromOption[R, A any](t Option[A]) ReaderIOOption[R, A]
func FromReader[R, A any](r Reader[R, A]) ReaderIOOption[R, A]
func MonadAlt[R, A any](first ReaderIOOption[R, A], second Lazy[ReaderIOOption[R, A]]) ReaderIOOption[R, A]
func MonadAp[R, A, B any](fab ReaderIOOption[R, func(A) B], fa ReaderIOOption[R, A]) ReaderIOOption[R, B]
func MonadChain[R, A, B any](ma ReaderIOOption[R, A], f Kleisli[R, A, B]) ReaderIOOption[R, B]
func MonadChainOptionK[R, A, B any](ma ReaderIOOption[R, A], f O.Kleisli[A, B]) ReaderIOOption[R, B]
func MonadFlap[R, A, B any](fab ReaderIOOption[R, func(A) B], a A) ReaderIOOption[R, B]
func MonadMap[R, A, B any](fa ReaderIOOption[R, A], f func(A) B) ReaderIOOption[R, B]
func None[R, A any]() ReaderIOOption[R, A]
func Of[R, A any](a A) ReaderIOOption[R, A]
func SequenceT1[R, A any](a ReaderIOOption[R, A]) ReaderIOOption[R, T.Tuple1[A]]
func SequenceT2[R, A, B any](
func SequenceT3[R, A, B, C any](
func SequenceT4[R, A, B, C, D any](
func Some[R, A any](r A) ReaderIOOption[R, A]
func SomeReader[R, A any](r Reader[R, A]) ReaderIOOption[R, A]
```

## package `github.com/IBM/fp-go/v2/readerioeither`

Import: `import "github.com/IBM/fp-go/v2/readerioeither"`

ReaderIOEither: `ReaderIOEither[R, E, A] = func(R) func() Either[E, A]`.

### Exported API

```go
func Ap[B, R, E, A any](fa ReaderIOEither[R, E, A]) func(fab ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B]
func BiMap[R, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, B]
func ChainFirstReaderOptionK[R, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, E, A, A]
func ChainLeft[R, EA, EB, A any](f Kleisli[R, EB, EA, A]) func(ReaderIOEither[R, EA, A]) ReaderIOEither[R, EB, A]
func ChainOptionK[R, A, B, E any](onNone Lazy[E]) func(func(A) Option[B]) Operator[R, E, A, B]
func ChainReaderOptionK[R, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, E, A, B]
func Eitherize0[F ~func(C) (R, error), C, R any](f F) func() ReaderIOEither[C, error, R]
func Eitherize1[F ~func(C, T0) (R, error), T0, C, R any](f F) func(T0) ReaderIOEither[C, error, R]
func Eitherize10[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOEither[C, error, R]
func Eitherize2[F ~func(C, T0, T1) (R, error), T0, T1, C, R any](f F) func(T0, T1) ReaderIOEither[C, error, R]
func Eitherize3[F ~func(C, T0, T1, T2) (R, error), T0, T1, T2, C, R any](f F) func(T0, T1, T2) ReaderIOEither[C, error, R]
func Eitherize4[F ~func(C, T0, T1, T2, T3) (R, error), T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) ReaderIOEither[C, error, R]
func Eitherize5[F ~func(C, T0, T1, T2, T3, T4) (R, error), T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) ReaderIOEither[C, error, R]
func Eitherize6[F ~func(C, T0, T1, T2, T3, T4, T5) (R, error), T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) ReaderIOEither[C, error, R]
func Eitherize7[F ~func(C, T0, T1, T2, T3, T4, T5, T6) (R, error), T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) ReaderIOEither[C, error, R]
func Eitherize8[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOEither[C, error, R]
func Eitherize9[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOEither[C, error, R]
func Eq[R, E, A any](eq EQ.Eq[either.Either[E, A]]) func(R) EQ.Eq[ReaderIOEither[R, E, A]]
func Flap[R, E, B, A any](a A) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B]
func Fold[R, E, A, B any](onLeft readerio.Kleisli[R, E, B], onRight func(A) ReaderIO[R, B]) func(ReaderIOEither[R, E, A]) ReaderIO[R, B]
func From0[F ~func(C) func() (R, error), C, R any](f F) func() ReaderIOEither[C, error, R]
func From1[F ~func(C, T0) func() (R, error), T0, C, R any](f F) func(T0) ReaderIOEither[C, error, R]
func From10[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOEither[C, error, R]
func From2[F ~func(C, T0, T1) func() (R, error), T0, T1, C, R any](f F) func(T0, T1) ReaderIOEither[C, error, R]
func From3[F ~func(C, T0, T1, T2) func() (R, error), T0, T1, T2, C, R any](f F) func(T0, T1, T2) ReaderIOEither[C, error, R]
func From4[F ~func(C, T0, T1, T2, T3) func() (R, error), T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) ReaderIOEither[C, error, R]
func From5[F ~func(C, T0, T1, T2, T3, T4) func() (R, error), T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) ReaderIOEither[C, error, R]
func From6[F ~func(C, T0, T1, T2, T3, T4, T5) func() (R, error), T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) ReaderIOEither[C, error, R]
func From7[F ~func(C, T0, T1, T2, T3, T4, T5, T6) func() (R, error), T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) ReaderIOEither[C, error, R]
func From8[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOEither[C, error, R]
func From9[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOEither[C, error, R]
func FromOption[R, A, E any](onNone Lazy[E]) func(Option[A]) ReaderIOEither[R, E, A]
func FromPredicate[R, E, A any](pred func(A) bool, onFalse func(A) E) func(A) ReaderIOEither[R, E, A]
func FromStrictEquals[R any, E, A comparable]() func(R) EQ.Eq[ReaderIOEither[R, E, A]]
func Functor[R, E, A, B any]() functor.Functor[A, B, ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]]
func GetOrElse[R, E, A any](onLeft readerio.Kleisli[R, E, A]) func(ReaderIOEither[R, E, A]) ReaderIO[R, A]
func Local[E, A, R1, R2 any](f func(R2) R1) func(ReaderIOEither[R1, E, A]) ReaderIOEither[R2, E, A]
func MapLeft[R, A, E1, E2 any](f func(E1) E2) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A]
func Monad[R, E, A, B any]() monad.Monad[A, B, ReaderIOEither[R, E, A], ReaderIOEither[R, E, B], ReaderIOEither[R, E, func(A) B]]
func OrLeft[A, E1, R, E2 any](onLeft func(E1) ReaderIO[R, E2]) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A]
func Pointed[R, E, A any]() pointed.Pointed[A, ReaderIOEither[R, E, A]]
func Read[E, A, R any](r R) func(ReaderIOEither[R, E, A]) IOEither[E, A]
func ReadIO[E, A, R any](r IO[R]) func(ReaderIOEither[R, E, A]) IOEither[E, A]
func ReadIOEither[A, R, E any](r IOEither[E, R]) func(ReaderIOEither[R, E, A]) IOEither[E, A]
func TapReaderOptionK[R, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, E, A, A]
func Traverse[R2, R1, E, A, B any](
func TraverseArrayWithIndex[R, E, A, B any](f func(int, A) ReaderIOEither[R, E, B]) func([]A) ReaderIOEither[R, E, []B]
func TraverseReader[R2, R1, E, A, B any](
func TraverseRecord[K comparable, R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(map[K]A) ReaderIOEither[R, E, map[K]B]
func TraverseRecordWithIndex[K comparable, R, E, A, B any](f func(K, A) ReaderIOEither[R, E, B]) func(map[K]A) ReaderIOEither[R, E, map[K]B]
func Uneitherize0[F ~func() ReaderIOEither[C, error, R], C, R any](f F) func(C) (R, error)
func Uneitherize1[F ~func(T0) ReaderIOEither[C, error, R], T0, C, R any](f F) func(C, T0) (R, error)
func Uneitherize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOEither[C, error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error)
func Uneitherize2[F ~func(T0, T1) ReaderIOEither[C, error, R], T0, T1, C, R any](f F) func(C, T0, T1) (R, error)
func Uneitherize3[F ~func(T0, T1, T2) ReaderIOEither[C, error, R], T0, T1, T2, C, R any](f F) func(C, T0, T1, T2) (R, error)
func Uneitherize4[F ~func(T0, T1, T2, T3) ReaderIOEither[C, error, R], T0, T1, T2, T3, C, R any](f F) func(C, T0, T1, T2, T3) (R, error)
func Uneitherize5[F ~func(T0, T1, T2, T3, T4) ReaderIOEither[C, error, R], T0, T1, T2, T3, T4, C, R any](f F) func(C, T0, T1, T2, T3, T4) (R, error)
func Uneitherize6[F ~func(T0, T1, T2, T3, T4, T5) ReaderIOEither[C, error, R], T0, T1, T2, T3, T4, T5, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5) (R, error)
func Uneitherize7[F ~func(T0, T1, T2, T3, T4, T5, T6) ReaderIOEither[C, error, R], T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6) (R, error)
func Uneitherize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOEither[C, error, R], T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6, T7) (R, error)
func Uneitherize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOEither[C, error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error)
type Consumer[A any] = consumer.Consumer[A]
type Either[E, A any] = either.Either[E, A]
type IO[T any] = io.IO[T]
type IOEither[E, A any] = ioeither.IOEither[E, A]
type Kleisli[R, E, A, B any] = reader.Reader[A, ReaderIOEither[R, E, B]]
func Contramap[E, A, R1, R2 any](f func(R2) R1) Kleisli[R2, E, ReaderIOEither[R1, E, A], A]
func FromReaderOption[R, A, E any](onNone Lazy[E]) Kleisli[R, E, ReaderOption[R, A], A]
func LocalIOEitherK[A, R1, R2, E any](f ioeither.Kleisli[E, R2, R1]) Kleisli[R2, E, ReaderIOEither[R1, E, A], A]
func LocalIOK[E, A, R1, R2 any](f io.Kleisli[R2, R1]) Kleisli[R2, E, ReaderIOEither[R1, E, A], A]
func OrElse[R, E1, E2, A any](onLeft Kleisli[R, E2, E1, A]) Kleisli[R, E2, ReaderIOEither[R, E1, A], A]
func Promap[R, E, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, E, ReaderIOEither[R, E, A], B]
func ReduceArray[R, E, A, B any](reduce func(B, A) B, initial B) Kleisli[R, E, []ReaderIOEither[R, E, A], B]
func ReduceArrayM[R, E, A any](m monoid.Monoid[A]) Kleisli[R, E, []ReaderIOEither[R, E, A], A]
func Sequence[R1, R2, E, A any](ma ReaderIOEither[R2, E, ReaderIOEither[R1, E, A]]) Kleisli[R2, E, R1, A]
func SequenceReader[R1, R2, E, A any](ma ReaderIOEither[R2, E, Reader[R1, A]]) Kleisli[R2, E, R1, A]
func SequenceReaderEither[R1, R2, E, A any](ma ReaderIOEither[R2, E, ReaderEither[R1, E, A]]) Kleisli[R2, E, R1, A]
func SequenceReaderIO[R1, R2, E, A any](ma ReaderIOEither[R2, E, ReaderIO[R1, A]]) Kleisli[R2, E, R1, A]
func TailRec[R, E, A, B any](f Kleisli[R, E, A, tailrec.Trampoline[A, B]]) Kleisli[R, E, A, B]
func TraverseArray[R, E, A, B any](f Kleisli[R, E, A, B]) Kleisli[R, E, []A, []B]
func TraverseReduceArray[R, E, A, B, C any](trfrm Kleisli[R, E, A, B], reduce func(C, B) C, initial C) Kleisli[R, E, []A, C]
func TraverseReduceArrayM[R, E, A, B any](trfrm Kleisli[R, E, A, B], m monoid.Monoid[B]) Kleisli[R, E, []A, B]
func WithResource[A, L, E, R, ANY any](onCreate ReaderIOEither[L, E, R], onRelease Kleisli[L, E, R, ANY]) Kleisli[L, E, Kleisli[L, E, R, A], A]
type Lazy[A any] = lazy.Lazy[A]
type Monoid[R, E, A any] = monoid.Monoid[ReaderIOEither[R, E, A]]
func AltMonoid[R, E, A any](zero lazy.Lazy[ReaderIOEither[R, E, A]]) Monoid[R, E, A]
func AlternativeMonoid[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A]
func ApplicativeMonoid[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A]
func ApplicativeMonoidPar[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A]
func ApplicativeMonoidSeq[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A]
type Operator[R, E, A, B any] = Kleisli[R, E, ReaderIOEither[R, E, A], B]
func After[R, E, A any](timestamp time.Time) Operator[R, E, A, A]
func Alt[R, E, A any](second L.Lazy[ReaderIOEither[R, E, A]]) Operator[R, E, A, A]
func ApEitherS[R, E, S1, S2, T any](
func ApEitherSL[R, E, S, T any](
func ApIOEitherS[R, E, S1, S2, T any](
func ApIOEitherSL[R, E, S, T any](
func ApIOS[R, E, S1, S2, T any](
func ApIOSL[R, E, S, T any](
func ApReaderIOS[R, E, S1, S2, T any](
func ApReaderIOSL[R, E, S, T any](
func ApReaderS[R, E, S1, S2, T any](
func ApReaderSL[R, E, S, T any](
func ApS[R, E, S1, S2, T any](
func ApSL[R, E, S, T any](
func Bind[R, E, S1, S2, T any](
func BindEitherK[R, E, S1, S2, T any](
func BindIOEitherK[R, E, S1, S2, T any](
func BindIOEitherKL[R, E, S, T any](
func BindIOK[R, E, S1, S2, T any](
func BindIOKL[R, E, S, T any](
func BindL[R, E, S, T any](
func BindReaderIOK[E, R, S1, S2, T any](
func BindReaderIOKL[E, R, S, T any](
func BindReaderK[E, R, S1, S2, T any](
func BindReaderKL[E, R, S, T any](
func BindTo[R, E, S1, T any](
func Chain[R, E, A, B any](f Kleisli[R, E, A, B]) Operator[R, E, A, B]
func ChainConsumer[R, E, A any](c Consumer[A]) Operator[R, E, A, struct{}]
func ChainEitherK[R, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, E, A, B]
func ChainFirst[R, E, A, B any](f Kleisli[R, E, A, B]) Operator[R, E, A, A]
func ChainFirstConsumer[R, E, A any](c Consumer[A]) Operator[R, E, A, A]
func ChainFirstEitherK[R, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, E, A, A]
func ChainFirstIOEitherK[R, E, A, B any](f IOE.Kleisli[E, A, B]) Operator[R, E, A, A]
func ChainFirstIOK[R, E, A, B any](f io.Kleisli[A, B]) Operator[R, E, A, A]
func ChainFirstLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A]
func ChainFirstLeftIOK[A, R, EA, B any](f io.Kleisli[EA, B]) Operator[R, EA, A, A]
func ChainFirstReaderEitherK[E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, E, A, A]
func ChainFirstReaderIOK[E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, E, A, A]
func ChainFirstReaderK[E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, E, A, A]
func ChainIOEitherK[R, E, A, B any](f IOE.Kleisli[E, A, B]) Operator[R, E, A, B]
func ChainIOK[R, E, A, B any](f io.Kleisli[A, B]) Operator[R, E, A, B]
func ChainReaderEitherK[E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, E, A, B]
func ChainReaderIOK[E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, E, A, B]
func ChainReaderK[E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, E, A, B]
func Delay[R, E, A any](delay time.Duration) Operator[R, E, A, A]
func FilterOrElse[R, E, A any](pred Predicate[A], onFalse func(A) E) Operator[R, E, A, A]
func Let[R, E, S1, S2, T any](
func LetL[R, E, S, T any](
func LetTo[R, E, S1, S2, T any](
func LetToL[R, E, S, T any](
func Map[R, E, A, B any](f func(A) B) Operator[R, E, A, B]
func MapTo[R, E, A, B any](b B) Operator[R, E, A, B]
func Tap[R, E, A, B any](f Kleisli[R, E, A, B]) Operator[R, E, A, A]
func TapEitherK[R, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, E, A, A]
func TapIOEitherK[R, E, A, B any](f IOE.Kleisli[E, A, B]) Operator[R, E, A, A]
func TapIOK[R, E, A, B any](f io.Kleisli[A, B]) Operator[R, E, A, A]
func TapLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A]
func TapLeftIOK[A, R, EA, B any](f io.Kleisli[EA, B]) Operator[R, EA, A, A]
func TapReaderEitherK[E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, E, A, A]
func TapReaderIOK[E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, E, A, A]
func TapReaderK[E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, E, A, A]
func WithLock[R, E, A any](lock func() context.CancelFunc) Operator[R, E, A, A]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderEither[R, E, A any] = readereither.ReaderEither[R, E, A]
type ReaderIO[R, A any] = readerio.ReaderIO[R, A]
func MonadFold[R, E, A, B any](ma ReaderIOEither[R, E, A], onLeft readerio.Kleisli[R, E, B], onRight func(A) ReaderIO[R, B]) ReaderIO[R, B]
type ReaderIOEither[R, E, A any] = Reader[R, IOEither[E, A]]
func Ask[R, E any]() ReaderIOEither[R, E, R]
func Asks[E, R, A any](r Reader[R, A]) ReaderIOEither[R, E, A]
func Bracket[
func Defer[R, E, A any](gen L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A]
func Do[R, E, S any](
func Flatten[R, E, A any](mma ReaderIOEither[R, E, ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A]
func FromEither[R, E, A any](t either.Either[E, A]) ReaderIOEither[R, E, A]
func FromIO[R, E, A any](ma IO[A]) ReaderIOEither[R, E, A]
func FromIOEither[R, E, A any](ma IOEither[E, A]) ReaderIOEither[R, E, A]
func FromReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A]
func FromReaderEither[R, E, A any](ma RE.ReaderEither[R, E, A]) ReaderIOEither[R, E, A]
func FromReaderIO[E, R, A any](ma ReaderIO[R, A]) ReaderIOEither[R, E, A]
func Left[R, A, E any](e E) ReaderIOEither[R, E, A]
func LeftIO[R, A, E any](ma IO[E]) ReaderIOEither[R, E, A]
func LeftReader[A, R, E any](ma Reader[R, E]) ReaderIOEither[R, E, A]
func LeftReaderIO[A, R, E any](me ReaderIO[R, E]) ReaderIOEither[R, E, A]
func Memoize[
func MonadAlt[R, E, A any](first ReaderIOEither[R, E, A], second L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A]
func MonadAp[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B]
func MonadApPar[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B]
func MonadApSeq[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B]
func MonadBiMap[R, E1, E2, A, B any](fa ReaderIOEither[R, E1, A], f func(E1) E2, g func(A) B) ReaderIOEither[R, E2, B]
func MonadChain[R, E, A, B any](fa ReaderIOEither[R, E, A], f Kleisli[R, E, A, B]) ReaderIOEither[R, E, B]
func MonadChainEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f either.Kleisli[E, A, B]) ReaderIOEither[R, E, B]
func MonadChainFirst[R, E, A, B any](fa ReaderIOEither[R, E, A], f Kleisli[R, E, A, B]) ReaderIOEither[R, E, A]
func MonadChainFirstEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f either.Kleisli[E, A, B]) ReaderIOEither[R, E, A]
func MonadChainFirstIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f io.Kleisli[A, B]) ReaderIOEither[R, E, A]
func MonadChainFirstLeft[A, R, EA, EB, B any](ma ReaderIOEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderIOEither[R, EA, A]
func MonadChainFirstReaderEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f RE.Kleisli[R, E, A, B]) ReaderIOEither[R, E, A]
func MonadChainFirstReaderIOK[E, R, A, B any](ma ReaderIOEither[R, E, A], f readerio.Kleisli[R, A, B]) ReaderIOEither[R, E, A]
func MonadChainFirstReaderK[E, R, A, B any](ma ReaderIOEither[R, E, A], f reader.Kleisli[R, A, B]) ReaderIOEither[R, E, A]
func MonadChainIOEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f IOE.Kleisli[E, A, B]) ReaderIOEither[R, E, B]
func MonadChainIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f io.Kleisli[A, B]) ReaderIOEither[R, E, B]
func MonadChainLeft[R, EA, EB, A any](fa ReaderIOEither[R, EA, A], f Kleisli[R, EB, EA, A]) ReaderIOEither[R, EB, A]
func MonadChainReaderEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f RE.Kleisli[R, E, A, B]) ReaderIOEither[R, E, B]
func MonadChainReaderIOK[E, R, A, B any](ma ReaderIOEither[R, E, A], f readerio.Kleisli[R, A, B]) ReaderIOEither[R, E, B]
func MonadChainReaderK[E, R, A, B any](ma ReaderIOEither[R, E, A], f reader.Kleisli[R, A, B]) ReaderIOEither[R, E, B]
func MonadFlap[R, E, B, A any](fab ReaderIOEither[R, E, func(A) B], a A) ReaderIOEither[R, E, B]
func MonadMap[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) B) ReaderIOEither[R, E, B]
func MonadMapLeft[R, E1, E2, A any](fa ReaderIOEither[R, E1, A], f func(E1) E2) ReaderIOEither[R, E2, A]
func MonadMapTo[R, E, A, B any](fa ReaderIOEither[R, E, A], b B) ReaderIOEither[R, E, B]
func MonadReduceArray[R, E, A, B any](as []ReaderIOEither[R, E, A], reduce func(B, A) B, initial B) ReaderIOEither[R, E, B]
func MonadReduceArrayM[R, E, A any](as []ReaderIOEither[R, E, A], m monoid.Monoid[A]) ReaderIOEither[R, E, A]
func MonadTap[R, E, A, B any](fa ReaderIOEither[R, E, A], f Kleisli[R, E, A, B]) ReaderIOEither[R, E, A]
func MonadTapEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f either.Kleisli[E, A, B]) ReaderIOEither[R, E, A]
func MonadTapIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f io.Kleisli[A, B]) ReaderIOEither[R, E, A]
func MonadTapLeft[A, R, EA, EB, B any](ma ReaderIOEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderIOEither[R, EA, A]
func MonadTapReaderEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f RE.Kleisli[R, E, A, B]) ReaderIOEither[R, E, A]
func MonadTapReaderIOK[E, R, A, B any](ma ReaderIOEither[R, E, A], f readerio.Kleisli[R, A, B]) ReaderIOEither[R, E, A]
func MonadTapReaderK[E, R, A, B any](ma ReaderIOEither[R, E, A], f reader.Kleisli[R, A, B]) ReaderIOEither[R, E, A]
func MonadTraverseReduceArray[R, E, A, B, C any](as []A, trfrm Kleisli[R, E, A, B], reduce func(C, B) C, initial C) ReaderIOEither[R, E, C]
func MonadTraverseReduceArrayM[R, E, A, B any](as []A, trfrm Kleisli[R, E, A, B], m monoid.Monoid[B]) ReaderIOEither[R, E, B]
func Of[R, E, A any](a A) ReaderIOEither[R, E, A]
func Retrying[R, E, A any](
func Right[R, E, A any](a A) ReaderIOEither[R, E, A]
func RightIO[R, E, A any](ma IO[A]) ReaderIOEither[R, E, A]
func RightReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A]
func RightReaderIO[E, R, A any](ma ReaderIO[R, A]) ReaderIOEither[R, E, A]
func SequenceArray[R, E, A any](ma []ReaderIOEither[R, E, A]) ReaderIOEither[R, E, []A]
func SequenceRecord[K comparable, R, E, A any](ma map[K]ReaderIOEither[R, E, A]) ReaderIOEither[R, E, map[K]A]
func SequenceT1[R, E, A any](a ReaderIOEither[R, E, A]) ReaderIOEither[R, E, T.Tuple1[A]]
func SequenceT2[R, E, A, B any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B]) ReaderIOEither[R, E, T.Tuple2[A, B]]
func SequenceT3[R, E, A, B, C any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B], c ReaderIOEither[R, E, C]) ReaderIOEither[R, E, T.Tuple3[A, B, C]]
func SequenceT4[R, E, A, B, C, D any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B], c ReaderIOEither[R, E, C], d ReaderIOEither[R, E, D]) ReaderIOEither[R, E, T.Tuple4[A, B, C, D]]
func Swap[R, E, A any](val ReaderIOEither[R, E, A]) ReaderIOEither[R, A, E]
func ThrowError[R, A, E any](e E) ReaderIOEither[R, E, A]
func TryCatch[R, E, A any](f func(R) func() (A, error), onThrow func(error) E) ReaderIOEither[R, E, A]
type ReaderOption[R, A any] = readeroption.ReaderOption[R, A]
```

## package `github.com/IBM/fp-go/v2/readerioresult`

Import: `import "github.com/IBM/fp-go/v2/readerioresult"`

ReaderIOResult: `ReaderIOResult[A] = func(context.Context) func() Result[A]`.

This is the primary effect type for context-aware IO with error handling.

### Exported API

```go
func BiMap[R, E, A, B any](f func(error) E, g func(A) B) func(ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, E, B]
func ChainFirstReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A]
func ChainLeft[R, A any](f Kleisli[R, error, A]) func(ReaderIOResult[R, A]) ReaderIOResult[R, A]
func ChainOptionK[R, A, B any](onNone Lazy[error]) func(func(A) Option[B]) Operator[R, A, B]
func ChainReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, B]
func Eitherize0[F ~func(C) (R, error), C, R any](f F) func() ReaderIOResult[C, R]
func Eitherize1[F ~func(C, T0) (R, error), T0, C, R any](f F) func(T0) ReaderIOResult[C, R]
func Eitherize10[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOResult[C, R]
func Eitherize2[F ~func(C, T0, T1) (R, error), T0, T1, C, R any](f F) func(T0, T1) ReaderIOResult[C, R]
func Eitherize3[F ~func(C, T0, T1, T2) (R, error), T0, T1, T2, C, R any](f F) func(T0, T1, T2) ReaderIOResult[C, R]
func Eitherize4[F ~func(C, T0, T1, T2, T3) (R, error), T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) ReaderIOResult[C, R]
func Eitherize5[F ~func(C, T0, T1, T2, T3, T4) (R, error), T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) ReaderIOResult[C, R]
func Eitherize6[F ~func(C, T0, T1, T2, T3, T4, T5) (R, error), T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) ReaderIOResult[C, R]
func Eitherize7[F ~func(C, T0, T1, T2, T3, T4, T5, T6) (R, error), T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) ReaderIOResult[C, R]
func Eitherize8[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOResult[C, R]
func Eitherize9[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOResult[C, R]
func Eq[R, A any](eq eq.Eq[Result[A]]) func(R) eq.Eq[ReaderIOResult[R, A]]
func Fold[R, A, B any](onLeft readerio.Kleisli[R, error, B], onRight func(A) ReaderIO[R, B]) func(ReaderIOResult[R, A]) ReaderIO[R, B]
func From0[F ~func(C) func() (R, error), C, R any](f F) func() ReaderIOResult[C, R]
func From1[F ~func(C, T0) func() (R, error), T0, C, R any](f F) func(T0) ReaderIOResult[C, R]
func From10[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOResult[C, R]
func From2[F ~func(C, T0, T1) func() (R, error), T0, T1, C, R any](f F) func(T0, T1) ReaderIOResult[C, R]
func From3[F ~func(C, T0, T1, T2) func() (R, error), T0, T1, T2, C, R any](f F) func(T0, T1, T2) ReaderIOResult[C, R]
func From4[F ~func(C, T0, T1, T2, T3) func() (R, error), T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) ReaderIOResult[C, R]
func From5[F ~func(C, T0, T1, T2, T3, T4) func() (R, error), T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) ReaderIOResult[C, R]
func From6[F ~func(C, T0, T1, T2, T3, T4, T5) func() (R, error), T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) ReaderIOResult[C, R]
func From7[F ~func(C, T0, T1, T2, T3, T4, T5, T6) func() (R, error), T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) ReaderIOResult[C, R]
func From8[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOResult[C, R]
func From9[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOResult[C, R]
func FromStrictEquals[R any, A comparable]() func(R) eq.Eq[ReaderIOResult[R, A]]
func Functor[R, A, B any]() functor.Functor[A, B, ReaderIOResult[R, A], ReaderIOResult[R, B]]
func GetOrElse[R, A any](onLeft readerio.Kleisli[R, error, A]) func(ReaderIOResult[R, A]) ReaderIO[R, A]
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, A]
func MapLeft[R, A, E any](f func(error) E) func(ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, E, A]
func Monad[R, A, B any]() monad.Monad[A, B, ReaderIOResult[R, A], ReaderIOResult[R, B], ReaderIOResult[R, func(A) B]]
func MonadBiMap[R, E, A, B any](fa ReaderIOResult[R, A], f func(error) E, g func(A) B) RIOE.ReaderIOEither[R, E, B]
func MonadMapLeft[R, E, A any](fa ReaderIOResult[R, A], f func(error) E) RIOE.ReaderIOEither[R, E, A]
func OrLeft[A, R, E any](onLeft readerio.Kleisli[R, error, E]) func(ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, E, A]
func Pointed[R, A any]() pointed.Pointed[A, ReaderIOResult[R, A]]
func Read[A, R any](r R) func(ReaderIOResult[R, A]) IOResult[A]
func ReadIO[A, R any](r IO[R]) func(ReaderIOResult[R, A]) IOResult[A]
func ReadIOEither[A, R any](r IOResult[R]) func(ReaderIOResult[R, A]) IOResult[A]
func ReadIOResult[A, R any](r IOResult[R]) func(ReaderIOResult[R, A]) IOResult[A]
func Swap[R, A any](val ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, A, error]
func TapReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A]
func Traverse[R2, R1, A, B any](
func TraverseReader[R2, R1, A, B any](
func Uneitherize0[F ~func() ReaderIOResult[C, R], C, R any](f F) func(C) (R, error)
func Uneitherize1[F ~func(T0) ReaderIOResult[C, R], T0, C, R any](f F) func(C, T0) (R, error)
func Uneitherize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOResult[C, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error)
func Uneitherize2[F ~func(T0, T1) ReaderIOResult[C, R], T0, T1, C, R any](f F) func(C, T0, T1) (R, error)
func Uneitherize3[F ~func(T0, T1, T2) ReaderIOResult[C, R], T0, T1, T2, C, R any](f F) func(C, T0, T1, T2) (R, error)
func Uneitherize4[F ~func(T0, T1, T2, T3) ReaderIOResult[C, R], T0, T1, T2, T3, C, R any](f F) func(C, T0, T1, T2, T3) (R, error)
func Uneitherize5[F ~func(T0, T1, T2, T3, T4) ReaderIOResult[C, R], T0, T1, T2, T3, T4, C, R any](f F) func(C, T0, T1, T2, T3, T4) (R, error)
func Uneitherize6[F ~func(T0, T1, T2, T3, T4, T5) ReaderIOResult[C, R], T0, T1, T2, T3, T4, T5, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5) (R, error)
func Uneitherize7[F ~func(T0, T1, T2, T3, T4, T5, T6) ReaderIOResult[C, R], T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6) (R, error)
func Uneitherize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOResult[C, R], T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6, T7) (R, error)
func Uneitherize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOResult[C, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error)
type Consumer[A any] = consumer.Consumer[A]
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[A any] = io.IO[A]
type IOEither[E, A any] = ioeither.IOEither[E, A]
type IOResult[A any] = ioresult.IOResult[A]
type Kleisli[R, A, B any] = reader.Reader[A, ReaderIOResult[R, B]]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIOResult[R1, A], A]
func FromOption[R, A any](onNone Lazy[error]) Kleisli[R, Option[A], A]
func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) error) Kleisli[R, A, A]
func FromReaderOption[R, A any](onNone Lazy[error]) Kleisli[R, ReaderOption[R, A], A]
func LocalIOEitherK[A, R1, R2 any](f ioeither.Kleisli[error, R2, R1]) Kleisli[R2, ReaderIOResult[R1, A], A]
func LocalIOK[A, R1, R2 any](f io.Kleisli[R2, R1]) Kleisli[R2, ReaderIOResult[R1, A], A]
func LocalIOResultK[A, R1, R2 any](f ioresult.Kleisli[R2, R1]) Kleisli[R2, ReaderIOResult[R1, A], A]
func Promap[R, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, ReaderIOResult[R, A], B]
func ReduceArray[R, A, B any](reduce func(B, A) B, initial B) Kleisli[R, []ReaderIOResult[R, A], B]
func ReduceArrayM[R, A any](m monoid.Monoid[A]) Kleisli[R, []ReaderIOResult[R, A], A]
func Sequence[R1, R2, A any](ma ReaderIOResult[R2, ReaderIOResult[R1, A]]) Kleisli[R2, R1, A]
func SequenceReader[R1, R2, A any](ma ReaderIOResult[R2, Reader[R1, A]]) Kleisli[R2, R1, A]
func SequenceReaderEither[R1, R2, A any](ma ReaderIOResult[R2, ReaderResult[R1, A]]) Kleisli[R2, R1, A]
func SequenceReaderIO[R1, R2, A any](ma ReaderIOResult[R2, ReaderIO[R1, A]]) Kleisli[R2, R1, A]
func SequenceReaderResult[R1, R2, A any](ma ReaderIOResult[R2, ReaderResult[R1, A]]) Kleisli[R2, R1, A]
func TailRec[R, A, B any](f Kleisli[R, A, tailrec.Trampoline[A, B]]) Kleisli[R, A, B]
func TraverseArray[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, []A, []B]
func TraverseArrayWithIndex[R, A, B any](f func(int, A) ReaderIOResult[R, B]) Kleisli[R, []A, []B]
func TraverseRecord[K comparable, R, A, B any](f Kleisli[R, A, B]) Kleisli[R, map[K]A, map[K]B]
func TraverseRecordWithIndex[K comparable, R, A, B any](f func(K, A) ReaderIOResult[R, B]) Kleisli[R, map[K]A, map[K]B]
func TraverseReduceArray[R, A, B, C any](trfrm Kleisli[R, A, B], reduce func(C, B) C, initial C) Kleisli[R, []A, C]
func TraverseReduceArrayM[R, A, B any](trfrm Kleisli[R, A, B], m monoid.Monoid[B]) Kleisli[R, []A, B]
func WithResource[A, L, R, ANY any](onCreate ReaderIOResult[L, R], onRelease Kleisli[L, R, ANY]) Kleisli[L, Kleisli[L, R, A], A]
type Lazy[A any] = lazy.Lazy[A]
type Monoid[R, A any] = monoid.Monoid[ReaderIOResult[R, A]]
func AltMonoid[R, A any](zero lazy.Lazy[ReaderIOResult[R, A]]) Monoid[R, A]
func AlternativeMonoid[R, A any](m monoid.Monoid[A]) Monoid[R, A]
func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) Monoid[R, A]
func ApplicativeMonoidPar[R, A any](m monoid.Monoid[A]) Monoid[R, A]
func ApplicativeMonoidSeq[R, A any](m monoid.Monoid[A]) Monoid[R, A]
type Operator[R, A, B any] = Kleisli[R, ReaderIOResult[R, A], B]
func After[R, A any](timestamp time.Time) Operator[R, A, A]
func Alt[R, A any](second Lazy[ReaderIOResult[R, A]]) Operator[R, A, A]
func Ap[B, R, A any](fa ReaderIOResult[R, A]) Operator[R, func(A) B, B]
func ApEitherS[R, S1, S2, T any](
func ApEitherSL[R, S, T any](
func ApIOEitherS[R, S1, S2, T any](
func ApIOEitherSL[R, S, T any](
func ApIOResultS[R, S1, S2, T any](
func ApIOResultSL[R, S, T any](
func ApIOS[R, S1, S2, T any](
func ApIOSL[R, S, T any](
func ApReaderIOS[R, S1, S2, T any](
func ApReaderIOSL[R, S, T any](
func ApReaderS[R, S1, S2, T any](
func ApReaderSL[R, S, T any](
func ApResultS[R, S1, S2, T any](
func ApResultSL[R, S, T any](
func ApS[R, S1, S2, T any](
func ApSL[R, S, T any](
func Bind[R, S1, S2, T any](
func BindEitherK[R, S1, S2, T any](
func BindIOEitherK[R, S1, S2, T any](
func BindIOEitherKL[R, S, T any](
func BindIOK[R, S1, S2, T any](
func BindIOKL[R, S, T any](
func BindIOResultK[R, S1, S2, T any](
func BindIOResultKL[R, S, T any](
func BindL[R, S, T any](
func BindReaderIOK[R, S1, S2, T any](
func BindReaderIOKL[R, S, T any](
func BindReaderK[R, S1, S2, T any](
func BindReaderKL[R, S, T any](
func BindResultK[R, S1, S2, T any](
func BindTo[R, S1, T any](
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainConsumer[R, A any](c Consumer[A]) Operator[R, A, Void]
func ChainEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B]
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A]
func ChainFirstConsumer[R, A any](c Consumer[A]) Operator[R, A, A]
func ChainFirstEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A]
func ChainFirstIOEitherK[R, A, B any](f ioresult.Kleisli[A, B]) Operator[R, A, A]
func ChainFirstIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, A]
func ChainFirstIOResultK[R, A, B any](f ioresult.Kleisli[A, B]) Operator[R, A, A]
func ChainFirstLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A]
func ChainFirstLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A]
func ChainFirstReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A]
func ChainFirstReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A]
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A]
func ChainFirstReaderResultK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A]
func ChainFirstResultK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A]
func ChainIOEitherK[R, A, B any](f func(A) IOResult[B]) Operator[R, A, B]
func ChainIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, B]
func ChainIOResultK[R, A, B any](f func(A) IOResult[B]) Operator[R, A, B]
func ChainReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, B]
func ChainReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, B]
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B]
func ChainReaderResultK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, B]
func ChainResultK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B]
func Delay[R, A any](delay time.Duration) Operator[R, A, A]
func FilterOrElse[R, A any](pred Predicate[A], onFalse func(A) error) Operator[R, A, A]
func Flap[R, B, A any](a A) Operator[R, func(A) B, B]
func Let[R, S1, S2, T any](
func LetL[R, S, T any](
func LetTo[R, S1, S2, T any](
func LetToL[R, S, T any](
func Map[R, A, B any](f func(A) B) Operator[R, A, B]
func MapTo[R, A, B any](b B) Operator[R, A, B]
func OrElse[R, A any](onLeft Kleisli[R, error, A]) Operator[R, A, A]
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A]
func TapEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A]
func TapIOEitherK[R, A, B any](f ioresult.Kleisli[A, B]) Operator[R, A, A]
func TapIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, A]
func TapIOResultK[R, A, B any](f ioresult.Kleisli[A, B]) Operator[R, A, A]
func TapLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A]
func TapLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A]
func TapReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A]
func TapReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A]
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A]
func TapReaderResultK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A]
func TapResultK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A]
func WithLock[R, A any](lock func() context.CancelFunc) Operator[R, A, A]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderIO[R, A any] = readerio.ReaderIO[R, A]
type ReaderIOResult[R, A any] = Reader[R, IOResult[A]]
func Ask[R any]() ReaderIOResult[R, R]
func Asks[R, A any](r Reader[R, A]) ReaderIOResult[R, A]
func Bracket[
func Defer[R, A any](gen Lazy[ReaderIOResult[R, A]]) ReaderIOResult[R, A]
func Do[R, S any](
func Flatten[R, A any](mma ReaderIOResult[R, ReaderIOResult[R, A]]) ReaderIOResult[R, A]
func FromEither[R, A any](t Result[A]) ReaderIOResult[R, A]
func FromIO[R, A any](ma IO[A]) ReaderIOResult[R, A]
func FromIOEither[R, A any](ma IOResult[A]) ReaderIOResult[R, A]
func FromIOResult[R, A any](ma IOResult[A]) ReaderIOResult[R, A]
func FromReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A]
func FromReaderEither[R, A any](ma RE.ReaderEither[R, error, A]) ReaderIOResult[R, A]
func FromReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A]
func FromResult[R, A any](t Result[A]) ReaderIOResult[R, A]
func Left[R, A any](e error) ReaderIOResult[R, A]
func LeftIO[R, A any](ma IO[error]) ReaderIOResult[R, A]
func LeftReader[A, R any](ma Reader[R, error]) ReaderIOResult[R, A]
func LeftReaderIO[A, R any](me ReaderIO[R, error]) ReaderIOResult[R, A]
func Memoize[R, A any](rdr ReaderIOResult[R, A]) ReaderIOResult[R, A]
func MonadAlt[R, A any](first ReaderIOResult[R, A], second Lazy[ReaderIOResult[R, A]]) ReaderIOResult[R, A]
func MonadAp[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B]
func MonadApPar[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B]
func MonadApSeq[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B]
func MonadChain[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, B]
func MonadChainEitherK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, B]
func MonadChainFirst[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadChainFirstEitherK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A]
func MonadChainFirstIOK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IO[B]) ReaderIOResult[R, A]
func MonadChainFirstLeft[A, R, B any](ma ReaderIOResult[R, A], f Kleisli[R, error, B]) ReaderIOResult[R, A]
func MonadChainFirstReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A]
func MonadChainFirstReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadChainFirstReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadChainFirstReaderResultK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A]
func MonadChainFirstResultK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A]
func MonadChainIOEitherK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IOResult[B]) ReaderIOResult[R, B]
func MonadChainIOK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IO[B]) ReaderIOResult[R, B]
func MonadChainIOResultK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IOResult[B]) ReaderIOResult[R, B]
func MonadChainLeft[R, A any](fa ReaderIOResult[R, A], f Kleisli[R, error, A]) ReaderIOResult[R, A]
func MonadChainReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, B]
func MonadChainReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, B]
func MonadChainReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, B]
func MonadChainReaderResultK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, B]
func MonadChainResultK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, B]
func MonadFlap[R, B, A any](fab ReaderIOResult[R, func(A) B], a A) ReaderIOResult[R, B]
func MonadMap[R, A, B any](fa ReaderIOResult[R, A], f func(A) B) ReaderIOResult[R, B]
func MonadMapTo[R, A, B any](fa ReaderIOResult[R, A], b B) ReaderIOResult[R, B]
func MonadReduceArray[R, A, B any](as []ReaderIOResult[R, A], reduce func(B, A) B, initial B) ReaderIOResult[R, B]
func MonadReduceArrayM[R, A any](as []ReaderIOResult[R, A], m monoid.Monoid[A]) ReaderIOResult[R, A]
func MonadTap[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadTapEitherK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A]
func MonadTapIOK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IO[B]) ReaderIOResult[R, A]
func MonadTapLeft[A, R, B any](ma ReaderIOResult[R, A], f Kleisli[R, error, B]) ReaderIOResult[R, A]
func MonadTapReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A]
func MonadTapReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadTapReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadTapReaderResultK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A]
func MonadTapResultK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A]
func MonadTraverseReduceArray[R, A, B, C any](as []A, trfrm Kleisli[R, A, B], reduce func(C, B) C, initial C) ReaderIOResult[R, C]
func MonadTraverseReduceArrayM[R, A, B any](as []A, trfrm Kleisli[R, A, B], m monoid.Monoid[B]) ReaderIOResult[R, B]
func Of[R, A any](a A) ReaderIOResult[R, A]
func Retrying[R, A any](
func Right[R, A any](a A) ReaderIOResult[R, A]
func RightIO[R, A any](ma IO[A]) ReaderIOResult[R, A]
func RightReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A]
func RightReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A]
func SequenceArray[R, A any](ma []ReaderIOResult[R, A]) ReaderIOResult[R, []A]
func SequenceRecord[K comparable, R, A any](ma map[K]ReaderIOResult[R, A]) ReaderIOResult[R, map[K]A]
func SequenceT1[R, A any](a ReaderIOResult[R, A]) ReaderIOResult[R, T.Tuple1[A]]
func SequenceT2[R, A, B any](a ReaderIOResult[R, A], b ReaderIOResult[R, B]) ReaderIOResult[R, T.Tuple2[A, B]]
func SequenceT3[R, A, B, C any](a ReaderIOResult[R, A], b ReaderIOResult[R, B], c ReaderIOResult[R, C]) ReaderIOResult[R, T.Tuple3[A, B, C]]
func SequenceT4[R, A, B, C, D any](a ReaderIOResult[R, A], b ReaderIOResult[R, B], c ReaderIOResult[R, C], d ReaderIOResult[R, D]) ReaderIOResult[R, T.Tuple4[A, B, C, D]]
func ThrowError[R, A any](e error) ReaderIOResult[R, A]
func TryCatch[R, A any](f func(R) func() (A, error), onThrow Endomorphism[error]) ReaderIOResult[R, A]
type ReaderOption[R, A any] = readeroption.ReaderOption[R, A]
type ReaderResult[R, A any] = readerresult.ReaderResult[R, A]
type Result[A any] = result.Result[A]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/readerresult`

Import: `import "github.com/IBM/fp-go/v2/readerresult"`

ReaderResult: `ReaderResult[A] = func(context.Context) Result[A]`.

### Exported API

```go
func ApEitherIS[
func ApResultIS[
func BindToEither[
func BindToEitherI[
func BindToReader[
func BindToResult[
func BindToResultI[
func ChainOptionIK[R, A, B any](onNone Lazy[error]) func(OI.Kleisli[A, B]) Operator[R, A, B]
func ChainOptionK[R, A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[R, A, B]
func Curry1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderResult[R, A]
func Curry2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1) func(T2) ReaderResult[R, A]
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderResult[R, A]
func Fold[R, A, B any](onLeft reader.Kleisli[R, error, B], onRight reader.Kleisli[R, A, B]) func(ReaderResult[R, A]) Reader[R, B]
func From0[R, A any](f func(R) (A, error)) func() ReaderResult[R, A]
func From1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderResult[R, A]
func From2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) ReaderResult[R, A]
func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderResult[R, A]
func GetOrElse[R, A any](onLeft reader.Kleisli[R, error, A]) func(ReaderResult[R, A]) Reader[R, A]
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderResult[R1, A]) ReaderResult[R2, A]
func Read[A, R any](r R) func(ReaderResult[R, A]) Result[A]
func Sequence[R1, R2, A any](ma ReaderResult[R2, ReaderResult[R1, A]]) reader.Kleisli[R2, R1, Result[A]]
func SequenceReader[R1, R2, A any](ma ReaderResult[R2, Reader[R1, A]]) reader.Kleisli[R2, R1, Result[A]]
func Traverse[R2, R1, A, B any](
func TraverseReader[R2, R1, A, B any](
func Uncurry1[R, T1, A any](f func(T1) ReaderResult[R, A]) func(R, T1) (A, error)
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) ReaderResult[R, A]) func(R, T1, T2) (A, error)
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderResult[R, A]) func(R, T1, T2, T3) (A, error)
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type Kleisli[R, A, B any] = Reader[A, ReaderResult[R, B]]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderResult[R1, A], A]
func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) error) Kleisli[R, A, A]
func Promap[R, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, ReaderResult[R, A], B]
func TailRec[R, A, B any](f Kleisli[R, A, tailrec.Trampoline[A, B]]) Kleisli[R, A, B]
func TraverseArray[L, A, B any](f Kleisli[L, A, B]) Kleisli[L, []A, []B]
func TraverseArrayWithIndex[L, A, B any](f func(int, A) ReaderResult[L, B]) Kleisli[L, []A, []B]
type Lazy[A any] = lazy.Lazy[A]
type Monoid[R, A any] = monoid.Monoid[ReaderResult[R, A]]
func AltMonoid[R, A any](zero Lazy[ReaderResult[R, A]]) Monoid[R, A]
func AlternativeMonoid[R, A any](m M.Monoid[A]) Monoid[R, A]
func ApplicativeMonoid[R, A any](m M.Monoid[A]) Monoid[R, A]
type Operator[R, A, B any] = Kleisli[R, ReaderResult[R, A], B]
func Alt[R, A any](second Lazy[ReaderResult[R, A]]) Operator[R, A, A]
func AltI[R, A any](second Lazy[RRI.ReaderResult[R, A]]) Operator[R, A, A]
func Ap[B, R, A any](fa ReaderResult[R, A]) Operator[R, func(A) B, B]
func ApEitherS[
func ApI[B, R, A any](fa RRI.ReaderResult[R, A]) Operator[R, func(A) B, B]
func ApIS[R, S1, S2, T any](
func ApISL[R, S, T any](
func ApReader[B, R, A any](fa Reader[R, A]) Operator[R, func(A) B, B]
func ApReaderS[
func ApResult[B, R, A any](fa Result[A]) Operator[R, func(A) B, B]
func ApResultI[B, R, A any](a A, err error) Operator[R, func(A) B, B]
func ApResultS[
func ApS[R, S1, S2, T any](
func ApSL[R, S, T any](
func BiMap[R, A, B any](f Endomorphism[error], g func(A) B) Operator[R, A, B]
func Bind[R, S1, S2, T any](
func BindEitherIK[R, S1, S2, T any](
func BindEitherK[R, S1, S2, T any](
func BindI[R, S1, S2, T any](
func BindIL[R, S, T any](
func BindL[R, S, T any](
func BindReaderK[R, S1, S2, T any](
func BindResultIK[R, S1, S2, T any](
func BindResultK[R, S1, S2, T any](
func BindTo[R, S1, T any](
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainEitherIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, B]
func ChainEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B]
func ChainI[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, B]
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B]
func FilterOrElse[R, A any](pred Predicate[A], onFalse func(A) error) Operator[R, A, A]
func Flap[R, B, A any](a A) Operator[R, func(A) B, B]
func Let[R, S1, S2, T any](
func LetL[R, S, T any](
func LetTo[R, S1, S2, T any](
func LetToL[R, S, T any](
func Map[R, A, B any](f func(A) B) Operator[R, A, B]
func MapLeft[R, A any](f Endomorphism[error]) Operator[R, A, A]
func OrElse[R, A any](onLeft Kleisli[R, error, A]) Operator[R, A, A]
func OrElseI[R, A any](onLeft RRI.Kleisli[R, error, A]) Operator[R, A, A]
func OrLeft[R, A any](onLeft reader.Kleisli[R, error, error]) Operator[R, A, A]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderResult[R, A any] = Reader[R, Result[A]]
func Ask[R any]() ReaderResult[R, R]
func Asks[R, A any](r Reader[R, A]) ReaderResult[R, A]
func Curry0[R, A any](f func(R) (A, error)) ReaderResult[R, A]
func Do[R, S any](
func Flatten[R, A any](mma ReaderResult[R, ReaderResult[R, A]]) ReaderResult[R, A]
func FlattenI[R, A any](mma ReaderResult[R, RRI.ReaderResult[R, A]]) ReaderResult[R, A]
func FromEither[R, A any](e Result[A]) ReaderResult[R, A]
func FromReader[R, A any](r Reader[R, A]) ReaderResult[R, A]
func FromReaderResultI[R, A any](rr RRI.ReaderResult[R, A]) ReaderResult[R, A]
func FromResult[R, A any](e Result[A]) ReaderResult[R, A]
func FromResultI[R, A any](a A, err error) ReaderResult[R, A]
func Left[R, A any](l error) ReaderResult[R, A]
func LeftReader[A, R any](l Reader[R, error]) ReaderResult[R, A]
func MonadAlt[R, A any](first ReaderResult[R, A], second Lazy[ReaderResult[R, A]]) ReaderResult[R, A]
func MonadAltI[R, A any](first ReaderResult[R, A], second Lazy[RRI.ReaderResult[R, A]]) ReaderResult[R, A]
func MonadAp[B, R, A any](fab ReaderResult[R, func(A) B], fa ReaderResult[R, A]) ReaderResult[R, B]
func MonadApI[B, R, A any](fab ReaderResult[R, func(A) B], fa RRI.ReaderResult[R, A]) ReaderResult[R, B]
func MonadApReader[B, R, A any](fab ReaderResult[R, func(A) B], fa Reader[R, A]) ReaderResult[R, B]
func MonadApResult[B, R, A any](fab ReaderResult[R, func(A) B], fa result.Result[A]) ReaderResult[R, B]
func MonadBiMap[R, A, B any](fa ReaderResult[R, A], f Endomorphism[error], g func(A) B) ReaderResult[R, B]
func MonadChain[R, A, B any](ma ReaderResult[R, A], f Kleisli[R, A, B]) ReaderResult[R, B]
func MonadChainEitherIK[R, A, B any](ma ReaderResult[R, A], f RI.Kleisli[A, B]) ReaderResult[R, B]
func MonadChainEitherK[R, A, B any](ma ReaderResult[R, A], f result.Kleisli[A, B]) ReaderResult[R, B]
func MonadChainI[R, A, B any](ma ReaderResult[R, A], f RRI.Kleisli[R, A, B]) ReaderResult[R, B]
func MonadChainReaderK[R, A, B any](ma ReaderResult[R, A], f reader.Kleisli[R, A, B]) ReaderResult[R, B]
func MonadFlap[R, A, B any](fab ReaderResult[R, func(A) B], a A) ReaderResult[R, B]
func MonadMap[R, A, B any](fa ReaderResult[R, A], f func(A) B) ReaderResult[R, B]
func MonadMapLeft[R, A any](fa ReaderResult[R, A], f Endomorphism[error]) ReaderResult[R, A]
func Of[R, A any](a A) ReaderResult[R, A]
func OfLazy[R, A any](r Lazy[A]) ReaderResult[R, A]
func Right[R, A any](r A) ReaderResult[R, A]
func RightReader[R, A any](r Reader[R, A]) ReaderResult[R, A]
func SequenceArray[L, A any](ma []ReaderResult[L, A]) ReaderResult[L, []A]
func SequenceT1[L, A any](a ReaderResult[L, A]) ReaderResult[L, T.Tuple1[A]]
func SequenceT2[L, A, B any](
func SequenceT3[L, A, B, C any](
func SequenceT4[L, A, B, C, D any](
type Result[A any] = result.Result[A]
```

---

# Context Specializations

## package `github.com/IBM/fp-go/v2/context/readerioresult`

Import: `import "github.com/IBM/fp-go/v2/context/readerioresult"`

Context-specialized ReaderIOResult with context.Context as Reader environment.

Provides Eitherize functions that handle context.Context as first parameter.

### Exported API

```go
func ChainFirstReaderOptionK[A, B any](onNone Lazy[error]) func(readeroption.Kleisli[context.Context, A, B]) Operator[A, A]
func ChainOptionK[A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[A, B]
func ChainReaderOptionK[A, B any](onNone Lazy[error]) func(readeroption.Kleisli[context.Context, A, B]) Operator[A, B]
func Contramap[A, R any](f pair.Kleisli[context.CancelFunc, R, context.Context]) RIOR.Kleisli[R, ReaderIOResult[A], A]
func Eitherize0[F ~func(context.Context) (R, error), R any](f F) func() ReaderIOResult[R]
func Eitherize10[F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOResult[R]
func Eitherize2[F ~func(context.Context, T0, T1) (R, error), T0, T1, R any](f F) func(T0, T1) ReaderIOResult[R]
func Eitherize3[F ~func(context.Context, T0, T1, T2) (R, error), T0, T1, T2, R any](f F) func(T0, T1, T2) ReaderIOResult[R]
func Eitherize4[F ~func(context.Context, T0, T1, T2, T3) (R, error), T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) ReaderIOResult[R]
func Eitherize5[F ~func(context.Context, T0, T1, T2, T3, T4) (R, error), T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) ReaderIOResult[R]
func Eitherize6[F ~func(context.Context, T0, T1, T2, T3, T4, T5) (R, error), T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) ReaderIOResult[R]
func Eitherize7[F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6) (R, error), T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) ReaderIOResult[R]
func Eitherize8[F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOResult[R]
func Eitherize9[F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOResult[R]
func Eq[A any](eq eq.Eq[Either[A]]) func(context.Context) eq.Eq[ReaderIOResult[A]]
func Filter[HKTA, A any](
func FilterMap[HKTA, HKTB, A, B any](
func GetOrElse[A any](onLeft readerio.Kleisli[error, A]) func(ReaderIOResult[A]) ReaderIO[A]
func Local[A, R any](f pair.Kleisli[context.CancelFunc, R, context.Context]) RIOR.Kleisli[R, ReaderIOResult[A], A]
func Promap[R, A, B any](f pair.Kleisli[context.CancelFunc, R, context.Context], g func(A) B) RIOR.Kleisli[R, ReaderIOResult[A], B]
func Read[A any](r context.Context) func(ReaderIOResult[A]) IOResult[A]
func ReadIO[A any](r IO[context.Context]) func(ReaderIOResult[A]) IOResult[A]
func ReadIOEither[A any](r IOResult[context.Context]) func(ReaderIOResult[A]) IOResult[A]
func ReadIOResult[A any](r IOResult[context.Context]) func(ReaderIOResult[A]) IOResult[A]
func TapReaderOptionK[A, B any](onNone Lazy[error]) func(readeroption.Kleisli[context.Context, A, B]) Operator[A, A]
func TraverseParTuple10[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], F9 ~func(A9) ReaderIOResult[T9], F10 ~func(A10) ReaderIOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) ReaderIOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseParTuple2[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) ReaderIOResult[tuple.Tuple2[T1, T2]]
func TraverseParTuple3[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) ReaderIOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseParTuple4[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) ReaderIOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseParTuple5[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) ReaderIOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseParTuple6[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) ReaderIOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseParTuple7[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) ReaderIOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseParTuple8[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) ReaderIOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseParTuple9[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], F9 ~func(A9) ReaderIOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) ReaderIOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseReader[R, A, B any](
func TraverseSeqTuple10[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], F9 ~func(A9) ReaderIOResult[T9], F10 ~func(A10) ReaderIOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) ReaderIOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseSeqTuple2[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) ReaderIOResult[tuple.Tuple2[T1, T2]]
func TraverseSeqTuple3[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) ReaderIOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseSeqTuple4[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) ReaderIOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseSeqTuple5[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) ReaderIOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseSeqTuple6[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) ReaderIOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseSeqTuple7[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) ReaderIOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseSeqTuple8[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) ReaderIOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseSeqTuple9[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], F9 ~func(A9) ReaderIOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) ReaderIOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseTuple10[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], F9 ~func(A9) ReaderIOResult[T9], F10 ~func(A10) ReaderIOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) ReaderIOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple2[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) ReaderIOResult[tuple.Tuple2[T1, T2]]
func TraverseTuple3[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) ReaderIOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseTuple4[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) ReaderIOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) ReaderIOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) ReaderIOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) ReaderIOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) ReaderIOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[F1 ~func(A1) ReaderIOResult[T1], F2 ~func(A2) ReaderIOResult[T2], F3 ~func(A3) ReaderIOResult[T3], F4 ~func(A4) ReaderIOResult[T4], F5 ~func(A5) ReaderIOResult[T5], F6 ~func(A6) ReaderIOResult[T6], F7 ~func(A7) ReaderIOResult[T7], F8 ~func(A8) ReaderIOResult[T8], F9 ~func(A9) ReaderIOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) ReaderIOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Uneitherize0[F ~func() ReaderIOResult[R], R any](f F) func(context.Context) (R, error)
func Uneitherize1[F ~func(T0) ReaderIOResult[R], T0, R any](f F) func(context.Context, T0) (R, error)
func Uneitherize10[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOResult[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error)
func Uneitherize2[F ~func(T0, T1) ReaderIOResult[R], T0, T1, R any](f F) func(context.Context, T0, T1) (R, error)
func Uneitherize3[F ~func(T0, T1, T2) ReaderIOResult[R], T0, T1, T2, R any](f F) func(context.Context, T0, T1, T2) (R, error)
func Uneitherize4[F ~func(T0, T1, T2, T3) ReaderIOResult[R], T0, T1, T2, T3, R any](f F) func(context.Context, T0, T1, T2, T3) (R, error)
func Uneitherize5[F ~func(T0, T1, T2, T3, T4) ReaderIOResult[R], T0, T1, T2, T3, T4, R any](f F) func(context.Context, T0, T1, T2, T3, T4) (R, error)
func Uneitherize6[F ~func(T0, T1, T2, T3, T4, T5) ReaderIOResult[R], T0, T1, T2, T3, T4, T5, R any](f F) func(context.Context, T0, T1, T2, T3, T4, T5) (R, error)
func Uneitherize7[F ~func(T0, T1, T2, T3, T4, T5, T6) ReaderIOResult[R], T0, T1, T2, T3, T4, T5, T6, R any](f F) func(context.Context, T0, T1, T2, T3, T4, T5, T6) (R, error)
func Uneitherize8[F ~func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOResult[R], T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7) (R, error)
func Uneitherize9[F ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOResult[R], T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error)
type CircuitBreaker[T any] = State[Env[T], ReaderIOResult[T]]
func MakeCircuitBreaker[T any](
type ClosedState = circuitbreaker.ClosedState
type Consumer[A any] = consumer.Consumer[A]
type ContextCancel = Pair[context.CancelFunc, context.Context]
type Either[A any] = either.Either[error, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type Env[T any] = Pair[IORef[circuitbreaker.BreakerState], ReaderIOResult[T]]
type IO[A any] = io.IO[A]
type IOEither[A any] = ioeither.IOEither[error, A]
type IORef[A any] = ioref.IORef[A]
type IOResult[A any] = ioresult.IOResult[A]
type Kleisli[A, B any] = reader.Reader[A, ReaderIOResult[B]]
func Eitherize1[F ~func(context.Context, T0) (R, error), T0, R any](f F) Kleisli[T0, R]
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A]
func FromReaderOption[A any](onNone Lazy[error]) Kleisli[ReaderOption[context.Context, A], A]
func SLog[A any](message string) Kleisli[Result[A], A]
func SLogWithCallback[A any](
func SequenceReader[R, A any](ma ReaderIOResult[Reader[R, A]]) Kleisli[R, A]
func SequenceReaderIO[R, A any](ma ReaderIOResult[RIO.ReaderIO[R, A]]) Kleisli[R, A]
func SequenceReaderResult[R, A any](ma ReaderIOResult[RR.ReaderResult[R, A]]) Kleisli[R, A]
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayPar[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) ReaderIOResult[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexPar[A, B any](f func(int, A) ReaderIOResult[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) ReaderIOResult[B]) Kleisli[[]A, []B]
func TraverseParTuple1[F1 ~func(A1) ReaderIOResult[T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordPar[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) ReaderIOResult[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexPar[K comparable, A, B any](f func(K, A) ReaderIOResult[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexSeq[K comparable, A, B any](f func(K, A) ReaderIOResult[B]) Kleisli[map[K]A, map[K]B]
func TraverseSeqTuple1[F1 ~func(A1) ReaderIOResult[T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func TraverseTuple1[F1 ~func(A1) ReaderIOResult[T1], T1, A1 any](f1 F1) Kleisli[tuple.Tuple1[A1], tuple.Tuple1[T1]]
func WithCloser[B any, A io.Closer](onCreate ReaderIOResult[A]) Kleisli[Kleisli[A, B], B]
func WithContextK[A, B any](f Kleisli[A, B]) Kleisli[A, B]
func WithResource[A, R, ANY any](onCreate ReaderIOResult[R], onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A]
type Lazy[A any] = lazy.Lazy[A]
type Lens[S, T any] = lens.Lens[S, T]
type LoggingID uint64
type Monoid[A any] = monoid.Monoid[ReaderIOResult[A]]
func AltMonoid[A any](zero Lazy[ReaderIOResult[A]]) Monoid[A]
func AlternativeMonoid[A any](m monoid.Monoid[A]) Monoid[A]
func ApplicativeMonoid[A any](m monoid.Monoid[A]) Monoid[A]
func ApplicativeMonoidPar[A any](m monoid.Monoid[A]) Monoid[A]
func ApplicativeMonoidSeq[A any](m monoid.Monoid[A]) Monoid[A]
type Operator[A, B any] = Kleisli[ReaderIOResult[A], B]
func Alt[A any](second Lazy[ReaderIOResult[A]]) Operator[A, A]
func Ap[B, A any](fa ReaderIOResult[A]) Operator[func(A) B, B]
func ApEitherS[S1, S2, T any](
func ApEitherSL[S, T any](
func ApIOEitherS[S1, S2, T any](
func ApIOEitherSL[S, T any](
func ApIOResultS[S1, S2, T any](
func ApIOResultSL[S, T any](
func ApIOS[S1, S2, T any](
func ApIOSL[S, T any](
func ApPar[B, A any](fa ReaderIOResult[A]) Operator[func(A) B, B]
func ApReaderIOS[S1, S2, T any](
func ApReaderIOSL[S, T any](
func ApReaderS[S1, S2, T any](
func ApReaderSL[S, T any](
func ApResultS[S1, S2, T any](
func ApResultSL[S, T any](
func ApS[S1, S2, T any](
func ApSL[S, T any](
func ApSeq[B, A any](fa ReaderIOResult[A]) Operator[func(A) B, B]
func Bind[S1, S2, T any](
func BindEitherK[S1, S2, T any](
func BindIOEitherK[S1, S2, T any](
func BindIOEitherKL[S, T any](
func BindIOK[S1, S2, T any](
func BindIOKL[S, T any](
func BindIOResultK[S1, S2, T any](
func BindIOResultKL[S, T any](
func BindL[S, T any](
func BindReaderIOK[S1, S2, T any](
func BindReaderIOKL[S, T any](
func BindReaderK[S1, S2, T any](
func BindReaderKL[S, T any](
func BindResultK[S1, S2, T any](
func BindTo[S1, T any](
func BindToP[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainConsumer[A any](c Consumer[A]) Operator[A, struct{}]
func ChainEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainFirstConsumer[A any](c Consumer[A]) Operator[A, A]
func ChainFirstEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A]
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A]
func ChainFirstLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A]
func ChainFirstReaderIOK[A, B any](f readerio.Kleisli[A, B]) Operator[A, A]
func ChainFirstReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, A]
func ChainFirstReaderResultK[A, B any](f readerresult.Kleisli[A, B]) Operator[A, A]
func ChainIOEitherK[A, B any](f ioresult.Kleisli[A, B]) Operator[A, B]
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B]
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A]
func ChainReaderIOK[A, B any](f readerio.Kleisli[A, B]) Operator[A, B]
func ChainReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, B]
func ChainReaderResultK[A, B any](f readerresult.Kleisli[A, B]) Operator[A, B]
func ChainResultK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B]
func ContramapIOK[A any](f io.Kleisli[context.Context, ContextCancel]) Operator[A, A]
func Delay[A any](delay time.Duration) Operator[A, A]
func FilterArray[A any](p Predicate[A]) Operator[[]A, []A]
func FilterIter[A any](p Predicate[A]) Operator[Seq[A], Seq[A]]
func FilterMapArray[A, B any](p option.Kleisli[A, B]) Operator[[]A, []B]
func FilterMapIter[A, B any](p option.Kleisli[A, B]) Operator[Seq[A], Seq[B]]
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Fold[A, B any](onLeft Kleisli[error, B], onRight Kleisli[A, B]) Operator[A, B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func LocalIOK[A any](f io.Kleisli[context.Context, ContextCancel]) Operator[A, A]
func LocalIOResultK[A any](f ioresult.Kleisli[context.Context, ContextCancel]) Operator[A, A]
func LogEntryExit[A any](name string) Operator[A, A]
func LogEntryExitWithCallback[A any](
func MakeSingletonBreaker[T any](
func Map[A, B any](f func(A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A]
func OrLeft[A any](onLeft func(error) ReaderIO[error]) Operator[A, A]
func Tap[A, B any](f Kleisli[A, B]) Operator[A, A]
func TapEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A]
func TapIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func TapLeft[A, B any](f Kleisli[error, B]) Operator[A, A]
func TapLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A]
func TapReaderIOK[A, B any](f readerio.Kleisli[A, B]) Operator[A, A]
func TapReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, A]
func TapReaderResultK[A, B any](f readerresult.Kleisli[A, B]) Operator[A, A]
func TapSLog[A any](message string) Operator[A, A]
func WithDeadline[A any](deadline time.Time) Operator[A, A]
func WithLock[A any](lock ReaderIOResult[context.CancelFunc]) Operator[A, A]
func WithTimeout[A any](timeout time.Duration) Operator[A, A]
type Option[A any] = option.Option[A]
type Pair[A, B any] = pair.Pair[A, B]
type Predicate[A any] = predicate.Predicate[A]
type Prism[S, T any] = prism.Prism[S, T]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderEither[R, E, A any] = readereither.ReaderEither[R, E, A]
type ReaderIO[A any] = readerio.ReaderIO[context.Context, A]
type ReaderIOResult[A any] = RIOR.ReaderIOResult[context.Context, A]
func Ask() ReaderIOResult[context.Context]
func Bracket[
func Defer[A any](gen Lazy[ReaderIOResult[A]]) ReaderIOResult[A]
func Do[S any](
func Flatten[A any](rdr ReaderIOResult[ReaderIOResult[A]]) ReaderIOResult[A]
func FromEither[A any](e Either[A]) ReaderIOResult[A]
func FromIO[A any](t IO[A]) ReaderIOResult[A]
func FromIOEither[A any](t IOResult[A]) ReaderIOResult[A]
func FromIOResult[A any](t IOResult[A]) ReaderIOResult[A]
func FromLazy[A any](t Lazy[A]) ReaderIOResult[A]
func FromReader[A any](t Reader[context.Context, A]) ReaderIOResult[A]
func FromReaderEither[A any](ma ReaderEither[context.Context, error, A]) ReaderIOResult[A]
func FromReaderIO[A any](t ReaderIO[A]) ReaderIOResult[A]
func FromReaderResult[A any](ma ReaderResult[A]) ReaderIOResult[A]
func FromResult[A any](e Result[A]) ReaderIOResult[A]
func Left[A any](l error) ReaderIOResult[A]
func Memoize[A any](rdr ReaderIOResult[A]) ReaderIOResult[A]
func MonadAlt[A any](first ReaderIOResult[A], second Lazy[ReaderIOResult[A]]) ReaderIOResult[A]
func MonadAp[B, A any](fab ReaderIOResult[func(A) B], fa ReaderIOResult[A]) ReaderIOResult[B]
func MonadApPar[B, A any](fab ReaderIOResult[func(A) B], fa ReaderIOResult[A]) ReaderIOResult[B]
func MonadApSeq[B, A any](fab ReaderIOResult[func(A) B], fa ReaderIOResult[A]) ReaderIOResult[B]
func MonadChain[A, B any](ma ReaderIOResult[A], f Kleisli[A, B]) ReaderIOResult[B]
func MonadChainEitherK[A, B any](ma ReaderIOResult[A], f either.Kleisli[error, A, B]) ReaderIOResult[B]
func MonadChainFirst[A, B any](ma ReaderIOResult[A], f Kleisli[A, B]) ReaderIOResult[A]
func MonadChainFirstEitherK[A, B any](ma ReaderIOResult[A], f either.Kleisli[error, A, B]) ReaderIOResult[A]
func MonadChainFirstIOK[A, B any](ma ReaderIOResult[A], f io.Kleisli[A, B]) ReaderIOResult[A]
func MonadChainFirstLeft[A, B any](ma ReaderIOResult[A], f Kleisli[error, B]) ReaderIOResult[A]
func MonadChainFirstReaderIOK[A, B any](ma ReaderIOResult[A], f readerio.Kleisli[A, B]) ReaderIOResult[A]
func MonadChainFirstReaderK[A, B any](ma ReaderIOResult[A], f reader.Kleisli[context.Context, A, B]) ReaderIOResult[A]
func MonadChainFirstReaderResultK[A, B any](ma ReaderIOResult[A], f readerresult.Kleisli[A, B]) ReaderIOResult[A]
func MonadChainIOK[A, B any](ma ReaderIOResult[A], f io.Kleisli[A, B]) ReaderIOResult[B]
func MonadChainLeft[A any](fa ReaderIOResult[A], f Kleisli[error, A]) ReaderIOResult[A]
func MonadChainReaderIOK[A, B any](ma ReaderIOResult[A], f readerio.Kleisli[A, B]) ReaderIOResult[B]
func MonadChainReaderK[A, B any](ma ReaderIOResult[A], f reader.Kleisli[context.Context, A, B]) ReaderIOResult[B]
func MonadChainReaderResultK[A, B any](ma ReaderIOResult[A], f readerresult.Kleisli[A, B]) ReaderIOResult[B]
func MonadFlap[B, A any](fab ReaderIOResult[func(A) B], a A) ReaderIOResult[B]
func MonadMap[A, B any](fa ReaderIOResult[A], f func(A) B) ReaderIOResult[B]
func MonadMapTo[A, B any](fa ReaderIOResult[A], b B) ReaderIOResult[B]
func MonadTap[A, B any](ma ReaderIOResult[A], f Kleisli[A, B]) ReaderIOResult[A]
func MonadTapEitherK[A, B any](ma ReaderIOResult[A], f either.Kleisli[error, A, B]) ReaderIOResult[A]
func MonadTapIOK[A, B any](ma ReaderIOResult[A], f io.Kleisli[A, B]) ReaderIOResult[A]
func MonadTapLeft[A, B any](ma ReaderIOResult[A], f Kleisli[error, B]) ReaderIOResult[A]
func MonadTapReaderIOK[A, B any](ma ReaderIOResult[A], f readerio.Kleisli[A, B]) ReaderIOResult[A]
func MonadTapReaderK[A, B any](ma ReaderIOResult[A], f reader.Kleisli[context.Context, A, B]) ReaderIOResult[A]
func MonadTapReaderResultK[A, B any](ma ReaderIOResult[A], f readerresult.Kleisli[A, B]) ReaderIOResult[A]
func MonadTraverseArrayPar[A, B any](as []A, f Kleisli[A, B]) ReaderIOResult[[]B]
func MonadTraverseArraySeq[A, B any](as []A, f Kleisli[A, B]) ReaderIOResult[[]B]
func MonadTraverseRecordPar[K comparable, A, B any](as map[K]A, f Kleisli[A, B]) ReaderIOResult[map[K]B]
func MonadTraverseRecordSeq[K comparable, A, B any](as map[K]A, f Kleisli[A, B]) ReaderIOResult[map[K]B]
func Never[A any]() ReaderIOResult[A]
func Of[A any](a A) ReaderIOResult[A]
func Retrying[A any](
func Right[A any](r A) ReaderIOResult[A]
func SequenceArray[A any](ma []ReaderIOResult[A]) ReaderIOResult[[]A]
func SequenceArrayPar[A any](ma []ReaderIOResult[A]) ReaderIOResult[[]A]
func SequenceArraySeq[A any](ma []ReaderIOResult[A]) ReaderIOResult[[]A]
func SequenceParT1[T1 any](
func SequenceParT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceParT2[T1, T2 any](
func SequenceParT3[T1, T2, T3 any](
func SequenceParT4[T1, T2, T3, T4 any](
func SequenceParT5[T1, T2, T3, T4, T5 any](
func SequenceParT6[T1, T2, T3, T4, T5, T6 any](
func SequenceParT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceParT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceParT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceParTuple1[T1 any](t tuple.Tuple1[ReaderIOResult[T1]]) ReaderIOResult[tuple.Tuple1[T1]]
func SequenceParTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8], ReaderIOResult[T9], ReaderIOResult[T10]]) ReaderIOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceParTuple2[T1, T2 any](t tuple.Tuple2[ReaderIOResult[T1], ReaderIOResult[T2]]) ReaderIOResult[tuple.Tuple2[T1, T2]]
func SequenceParTuple3[T1, T2, T3 any](t tuple.Tuple3[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3]]) ReaderIOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceParTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4]]) ReaderIOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceParTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5]]) ReaderIOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceParTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6]]) ReaderIOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceParTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7]]) ReaderIOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceParTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8]]) ReaderIOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceParTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8], ReaderIOResult[T9]]) ReaderIOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceRecord[K comparable, A any](ma map[K]ReaderIOResult[A]) ReaderIOResult[map[K]A]
func SequenceRecordPar[K comparable, A any](ma map[K]ReaderIOResult[A]) ReaderIOResult[map[K]A]
func SequenceRecordSeq[K comparable, A any](ma map[K]ReaderIOResult[A]) ReaderIOResult[map[K]A]
func SequenceSeqT1[T1 any](
func SequenceSeqT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceSeqT2[T1, T2 any](
func SequenceSeqT3[T1, T2, T3 any](
func SequenceSeqT4[T1, T2, T3, T4 any](
func SequenceSeqT5[T1, T2, T3, T4, T5 any](
func SequenceSeqT6[T1, T2, T3, T4, T5, T6 any](
func SequenceSeqT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceSeqT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceSeqT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceSeqTuple1[T1 any](t tuple.Tuple1[ReaderIOResult[T1]]) ReaderIOResult[tuple.Tuple1[T1]]
func SequenceSeqTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8], ReaderIOResult[T9], ReaderIOResult[T10]]) ReaderIOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceSeqTuple2[T1, T2 any](t tuple.Tuple2[ReaderIOResult[T1], ReaderIOResult[T2]]) ReaderIOResult[tuple.Tuple2[T1, T2]]
func SequenceSeqTuple3[T1, T2, T3 any](t tuple.Tuple3[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3]]) ReaderIOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceSeqTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4]]) ReaderIOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceSeqTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5]]) ReaderIOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceSeqTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6]]) ReaderIOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceSeqTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7]]) ReaderIOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceSeqTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8]]) ReaderIOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceSeqTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8], ReaderIOResult[T9]]) ReaderIOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceT1[T1 any](
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceT2[T1, T2 any](
func SequenceT3[T1, T2, T3 any](
func SequenceT4[T1, T2, T3, T4 any](
func SequenceT5[T1, T2, T3, T4, T5 any](
func SequenceT6[T1, T2, T3, T4, T5, T6 any](
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceTuple1[T1 any](t tuple.Tuple1[ReaderIOResult[T1]]) ReaderIOResult[tuple.Tuple1[T1]]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8], ReaderIOResult[T9], ReaderIOResult[T10]]) ReaderIOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple2[T1, T2 any](t tuple.Tuple2[ReaderIOResult[T1], ReaderIOResult[T2]]) ReaderIOResult[tuple.Tuple2[T1, T2]]
func SequenceTuple3[T1, T2, T3 any](t tuple.Tuple3[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3]]) ReaderIOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4]]) ReaderIOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5]]) ReaderIOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6]]) ReaderIOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7]]) ReaderIOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8]]) ReaderIOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[ReaderIOResult[T1], ReaderIOResult[T2], ReaderIOResult[T3], ReaderIOResult[T4], ReaderIOResult[T5], ReaderIOResult[T6], ReaderIOResult[T7], ReaderIOResult[T8], ReaderIOResult[T9]]) ReaderIOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Timer(delay time.Duration) ReaderIOResult[time.Time]
func TryCatch[A any](f func(context.Context) func() (A, error)) ReaderIOResult[A]
func WithContext[A any](ma ReaderIOResult[A]) ReaderIOResult[A]
type ReaderOption[R, A any] = readeroption.ReaderOption[R, A]
type ReaderResult[A any] = readerresult.ReaderResult[A]
type Result[A any] = result.Result[A]
type Semigroup[A any] = semigroup.Semigroup[ReaderIOResult[A]]
func AltSemigroup[A any]() Semigroup[A]
type Seq[A any] = iter.Seq[A]
type State[S, A any] = state.State[S, A]
type Trampoline[B, L any] = tailrec.Trampoline[B, L]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/context/readerresult`

Import: `import "github.com/IBM/fp-go/v2/context/readerresult"`

Context-specialized ReaderResult with context.Context as Reader environment.

### Exported API

```go
func Ap[A, B any](fa ReaderResult[A]) func(ReaderResult[func(A) B]) ReaderResult[B]
func ChainEitherK[A, B any](f func(A) Either[B]) func(ma ReaderResult[A]) ReaderResult[B]
func ChainOptionK[A, B any](onNone func() error) func(option.Kleisli[A, B]) Operator[A, B]
func Contramap[A, R any](f pair.Kleisli[context.CancelFunc, R, context.Context]) RR.Kleisli[R, ReaderResult[A], A]
func Curry2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1) Kleisli[T2, A]
func Curry3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1) func(T2) Kleisli[T3, A]
func From0[A any](f func(context.Context) (A, error)) func() ReaderResult[A]
func From2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1, T2) ReaderResult[A]
func From3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderResult[A]
func Local[A, R any](f pair.Kleisli[context.CancelFunc, R, context.Context]) RR.Kleisli[R, ReaderResult[A], A]
func Promap[R, A, B any](f pair.Kleisli[context.CancelFunc, R, context.Context], g func(A) B) RR.Kleisli[R, ReaderResult[A], B]
func Read[A any](r context.Context) func(ReaderResult[A]) Result[A]
func ReadEither[A any](r Result[context.Context]) func(ReaderResult[A]) Result[A]
func ReadIO[A any](r IO[context.Context]) func(ReaderResult[A]) IOResult[A]
func ReadIOEither[A any](r IOResult[context.Context]) func(ReaderResult[A]) IOResult[A]
func ReadIOResult[A any](r IOResult[context.Context]) func(ReaderResult[A]) IOResult[A]
func ReadResult[A any](r Result[context.Context]) func(ReaderResult[A]) Result[A]
func SequenceReader[R, A any](ma ReaderResult[Reader[R, A]]) reader.Kleisli[context.Context, R, Result[A]]
func TraverseReader[R, A, B any](
func Uncurry1[T1, A any](f Kleisli[T1, A]) func(context.Context, T1) (A, error)
func Uncurry2[T1, T2, A any](f func(T1) Kleisli[T2, A]) func(context.Context, T1, T2) (A, error)
func Uncurry3[T1, T2, T3, A any](f func(T1) func(T2) Kleisli[T3, A]) func(context.Context, T1, T2, T3) (A, error)
type Either[A any] = either.Either[error, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[A any] = io.IO[A]
type IOResult[A any] = ioresult.IOResult[A]
type Kleisli[A, B any] = reader.Reader[A, ReaderResult[B]]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func Bind[S1, S2, T any](
func BindL[S, T any](
func Curry1[T1, A any](f func(context.Context, T1) (A, error)) Kleisli[T1, A]
func From1[T1, A any](f func(context.Context, T1) (A, error)) Kleisli[T1, A]
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func OrElse[A any](onLeft Kleisli[error, A]) Kleisli[ReaderResult[A], A]
func SLog[A any](message string) Kleisli[Result[A], A]
func SLogWithCallback[A any](
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) ReaderResult[B]) Kleisli[[]A, []B]
func WithContextK[A, B any](f Kleisli[A, B]) Kleisli[A, B]
type Lens[S, T any] = lens.Lens[S, T]
type Operator[A, B any] = Kleisli[ReaderResult[A], B]
func BindTo[S1, T any](
func BindToP[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A]
func ChainFirstLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A]
func ChainIOEitherK[A, B any](f ioresult.Kleisli[A, B]) Operator[A, B]
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B]
func ChainIOResultK[A, B any](f ioresult.Kleisli[A, B]) Operator[A, B]
func ChainTo[A, B any](b ReaderResult[B]) Operator[A, B]
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Map[A, B any](f func(A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func TapIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func TapLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A]
func TapSLog[A any](message string) Operator[A, A]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Prism[S, T any] = prism.Prism[S, T]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderResult[A any] = readereither.ReaderEither[context.Context, error, A]
func Ask() ReaderResult[context.Context]
func Curry0[A any](f func(context.Context) (A, error)) ReaderResult[A]
func Do[S any](
func FromEither[A any](e Either[A]) ReaderResult[A]
func FromIO[A any](t io.IO[A]) ReaderResult[A]
func FromIOResult[A any](t ioresult.IOResult[A]) ReaderResult[A]
func FromReader[A any](r Reader[context.Context, A]) ReaderResult[A]
func Left[A any](l error) ReaderResult[A]
func MonadAp[A, B any](fab ReaderResult[func(A) B], fa ReaderResult[A]) ReaderResult[B]
func MonadChain[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[B]
func MonadChainEitherK[A, B any](ma ReaderResult[A], f func(A) Either[B]) ReaderResult[B]
func MonadChainFirst[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[A]
func MonadChainFirstIOK[A, B any](ma ReaderResult[A], f io.Kleisli[A, B]) ReaderResult[A]
func MonadChainIOK[A, B any](ma ReaderResult[A], f io.Kleisli[A, B]) ReaderResult[B]
func MonadChainTo[A, B any](ma ReaderResult[A], b ReaderResult[B]) ReaderResult[B]
func MonadFlap[B, A any](fab ReaderResult[func(A) B], a A) ReaderResult[B]
func MonadMap[A, B any](fa ReaderResult[A], f func(A) B) ReaderResult[B]
func MonadMapTo[A, B any](ma ReaderResult[A], b B) ReaderResult[B]
func MonadTapIOK[A, B any](ma ReaderResult[A], f io.Kleisli[A, B]) ReaderResult[A]
func Of[A any](a A) ReaderResult[A]
func Retrying[A any](
func Right[A any](r A) ReaderResult[A]
func SequenceArray[A any](ma []ReaderResult[A]) ReaderResult[[]A]
func SequenceT1[A any](a ReaderResult[A]) ReaderResult[tuple.Tuple1[A]]
func SequenceT2[A, B any](a ReaderResult[A], b ReaderResult[B]) ReaderResult[tuple.Tuple2[A, B]]
func SequenceT3[A, B, C any](a ReaderResult[A], b ReaderResult[B], c ReaderResult[C]) ReaderResult[tuple.Tuple3[A, B, C]]
func SequenceT4[A, B, C, D any](a ReaderResult[A], b ReaderResult[B], c ReaderResult[C], d ReaderResult[D]) ReaderResult[tuple.Tuple4[A, B, C, D]]
func WithContext[A any](ma ReaderResult[A]) ReaderResult[A]
type Result[A any] = result.Result[A]
type Trampoline[A, B any] = tailrec.Trampoline[A, B]
```

---

# Effect System

## package `github.com/IBM/fp-go/v2/effect`

Import: `import "github.com/IBM/fp-go/v2/effect"`

The Effect system: dependency-injection-aware effect type.

Key types:
- `Effect[C, A] = ReaderReaderIOResult[C, A]` -- reads config C, uses context.Context, produces Result[A]
- Use `Provide` to supply config, `Read` to extract thunks

### Exported API

```go
func Filter[C, HKTA, A any](
func FilterMap[C, HKTA, HKTB, A, B any](
func LocalEffectK[A, C1, C2 any](f Kleisli[C2, C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalIOK[A, C1, C2 any](f io.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalIOResultK[A, C1, C2 any](f ioresult.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalReaderK[A, C1, C2 any](f reader.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalResultK[A, C1, C2 any](f result.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalThunkK[A, C1, C2 any](f thunk.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func Provide[A, C any](c C) func(Effect[C, A]) ReaderIOResult[A]
func Read[A, C any](c C) func(Effect[C, A]) Thunk[A]
func RunSync[A any](fa ReaderIOResult[A]) readerresult.ReaderResult[A]
type Effect[C, A any] = readerreaderioresult.ReaderReaderIOResult[C, A]
func Ask[C any]() Effect[C, C]
func Asks[C, A any](r Reader[C, A]) Effect[C, A]
func Do[C, S any](
func Eitherize[C, T any](f func(C, context.Context) (T, error)) Effect[C, T]
func Fail[C, A any](err error) Effect[C, A]
func FromResult[C, A any](r Result[A]) Effect[C, A]
func FromThunk[C, A any](f Thunk[A]) Effect[C, A]
func Of[C, A any](a A) Effect[C, A]
func Retrying[C, A any](
func Succeed[C, A any](a A) Effect[C, A]
func Suspend[C, A any](fa Lazy[Effect[C, A]]) Effect[C, A]
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[A any] = io.IO[A]
type IOEither[E, A any] = ioeither.IOEither[E, A]
type IOResult[A any] = ioresult.IOResult[A]
type Kleisli[C, A, B any] = readerreaderioresult.Kleisli[C, A, B]
func Contramap[A, C1, C2 any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A]
func Eitherize1[C, A, T any](f func(C, context.Context, A) (T, error)) Kleisli[C, A, T]
func Local[A, C1, C2 any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A]
func Promap[E, A, D, B any](f Reader[D, E], g Reader[A, B]) Kleisli[D, Effect[E, A], B]
func Ternary[C, A, B any](pred Predicate[A], onTrue, onFalse Kleisli[C, A, B]) Kleisli[C, A, B]
func TraverseArray[C, A, B any](f Kleisli[C, A, B]) Kleisli[C, []A, []B]
type Lazy[A any] = lazy.Lazy[A]
type Lens[S, T any] = lens.Lens[S, T]
type Monoid[A any] = monoid.Monoid[A]
func AlternativeMonoid[C, A any](m monoid.Monoid[A]) Monoid[Effect[C, A]]
func ApplicativeMonoid[C, A any](m monoid.Monoid[A]) Monoid[Effect[C, A]]
type Operator[C, A, B any] = readerreaderioresult.Operator[C, A, B]
func Ap[B, C, A any](fa Effect[C, A]) Operator[C, func(A) B, B]
func ApEitherS[C, S1, S2, T any](
func ApEitherSL[C, S, T any](
func ApIOEitherS[C, S1, S2, T any](
func ApIOEitherSL[C, S, T any](
func ApIOS[C, S1, S2, T any](
func ApIOSL[C, S, T any](
func ApReaderIOS[C, S1, S2, T any](
func ApReaderIOSL[C, S, T any](
func ApReaderS[C, S1, S2, T any](
func ApReaderSL[C, S, T any](
func ApS[C, S1, S2, T any](
func ApSL[C, S, T any](
func Bind[C, S1, S2, T any](
func BindEitherK[C, S1, S2, T any](
func BindIOEitherK[C, S1, S2, T any](
func BindIOEitherKL[C, S, T any](
func BindIOK[C, S1, S2, T any](
func BindIOKL[C, S, T any](
func BindIOResultK[C, S1, S2, T any](
func BindL[C, S, T any](
func BindReaderIOK[C, S1, S2, T any](
func BindReaderIOKL[C, S, T any](
func BindReaderK[C, S1, S2, T any](
func BindReaderKL[C, S, T any](
func BindTo[C, S1, T any](
func Chain[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, B]
func ChainFirst[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, A]
func ChainFirstIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, A]
func ChainFirstThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, A]
func ChainIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, B]
func ChainReaderIOK[C, A, B any](f readerio.Kleisli[C, A, B]) Operator[C, A, B]
func ChainReaderK[C, A, B any](f reader.Kleisli[C, A, B]) Operator[C, A, B]
func ChainResultK[C, A, B any](f result.Kleisli[A, B]) Operator[C, A, B]
func ChainThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, B]
func FilterArray[C, A any](p Predicate[A]) Operator[C, []A, []A]
func FilterIter[C, A any](p Predicate[A]) Operator[C, Seq[A], Seq[A]]
func FilterMapArray[C, A, B any](p option.Kleisli[A, B]) Operator[C, []A, []B]
func FilterMapIter[C, A, B any](p option.Kleisli[A, B]) Operator[C, Seq[A], Seq[B]]
func Let[C, S1, S2, T any](
func LetL[C, S, T any](
func LetTo[C, S1, S2, T any](
func LetToL[C, S, T any](
func Map[C, A, B any](f func(A) B) Operator[C, A, B]
func Tap[C, A, ANY any](f Kleisli[C, A, ANY]) Operator[C, A, A]
func TapIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, A]
func TapThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderIO[R, A any] = readerio.ReaderIO[R, A]
type ReaderIOResult[A any] = readerioresult.ReaderIOResult[A]
type Result[A any] = result.Result[A]
type Seq[A any] = iter.Seq[A]
type Thunk[A any] = ReaderIOResult[A]
```

---

# State Monads

## package `github.com/IBM/fp-go/v2/state`

Import: `import "github.com/IBM/fp-go/v2/state"`

State monad: `State[S, A] = func(S) Pair[S, A]`.

Encapsulates stateful computations with Get/Put/Modify.

### Exported API

```go
func Applicative[S, A, B any]() applicative.Applicative[A, B, State[S, A], State[S, B], State[S, func(A) B]]
func ApplicativeMonoid[S, A any](m M.Monoid[A]) M.Monoid[State[S, A]]
func Eq[S, A any](w eq.Eq[S], a eq.Eq[A]) func(S) eq.Eq[State[S, A]]
func Evaluate[A, S any](s S) func(State[S, A]) A
func Execute[A, S any](s S) func(State[S, A]) S
func FromStrictEquals[S, A comparable]() func(S) eq.Eq[State[S, A]]
func Functor[S, A, B any]() functor.Functor[A, B, State[S, A], State[S, B]]
func Monad[S, A, B any]() monad.Monad[A, B, State[S, A], State[S, B], State[S, func(A) B]]
func Pointed[S, A any]() pointed.Pointed[A, State[S, A]]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
func ApSL[ST, S, T any](
func BindL[ST, S, T any](
func LetL[ST, S, T any](
func LetToL[ST, S, T any](
type Kleisli[S, A, B any] = Reader[A, State[S, B]]
func IMap[A, S2, S1, B any](f iso.Iso[S2, S1], g func(A) B) Kleisli[S2, State[S1, A], B]
func MapState[A, S2, S1 any](f iso.Iso[S2, S1]) Kleisli[S2, State[S1, A], A]
type Lens[S, A any] = lens.Lens[S, A]
type Operator[S, A, B any] = Kleisli[S, State[S, A], B]
func Ap[B, S, A any](ga State[S, A]) Operator[S, func(A) B, B]
func ApS[ST, S1, S2, T any](
func Bind[ST, S1, S2, T any](
func BindTo[ST, S1, T any](
func Chain[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, B]
func ChainFirst[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, A]
func Flap[S, A, B any](a A) Operator[S, func(A) B, B]
func Let[ST, S1, S2, T any](
func LetTo[ST, S1, S2, T any](
func Map[S any, FCT ~func(A) B, A, B any](f FCT) Operator[S, A, B]
type Pair[L, R any] = pair.Pair[L, R]
type Reader[R, A any] = reader.Reader[R, A]
type State[S, A any] = Reader[S, pair.Pair[S, A]]
func Do[ST, A any](
func Flatten[S, A any](mma State[S, State[S, A]]) State[S, A]
func Get[S any]() State[S, S]
func Gets[FCT ~func(S) A, A, S any](f FCT) State[S, A]
func Modify[FCT ~func(S) S, S any](f FCT) State[S, Void]
func MonadAp[B, S, A any](fab State[S, func(A) B], fa State[S, A]) State[S, B]
func MonadChain[S any, FCT ~func(A) State[S, B], A, B any](fa State[S, A], f FCT) State[S, B]
func MonadChainFirst[S any, FCT ~func(A) State[S, B], A, B any](ma State[S, A], f FCT) State[S, A]
func MonadFlap[FAB ~func(A) B, S, A, B any](fab State[S, FAB], a A) State[S, B]
func MonadMap[S any, FCT ~func(A) B, A, B any](fa State[S, A], f FCT) State[S, B]
func Of[S, A any](a A) State[S, A]
func Put[S any]() State[S, Void]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/stateio`

Import: `import "github.com/IBM/fp-go/v2/stateio"`

StateIO: `StateIO[S, A] = func(S) IO[Pair[S, A]]`.

### Exported API

```go
func Applicative[
func ApplicativeMonoid[S, A any](m M.Monoid[A]) M.Monoid[StateIO[S, A]]
func Eq[
func FromStrictEquals[
func Functor[
func Monad[
func Pointed[
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
func ApSL[ST, S, T any](
func BindL[ST, S, T any](
func LetL[ST, S, T any](
func LetToL[ST, S, T any](
type IO[A any] = io.IO[A]
type Kleisli[S, A, B any] = Reader[A, StateIO[S, B]]
func FromIOK[S, A, B any](f func(A) IO[B]) Kleisli[S, A, B]
func WithResource[A, S, RES, ANY any](
type Lens[S, A any] = lens.Lens[S, A]
type Operator[S, A, B any] = Reader[StateIO[S, A], StateIO[S, B]]
func Ap[B, S, A any](fa StateIO[S, A]) Operator[S, func(A) B, B]
func ApS[ST, S1, S2, T any](
func Bind[ST, S1, S2, T any](
func BindTo[ST, S1, T any](
func Chain[S, A, B any](f Kleisli[S, A, B]) Operator[S, A, B]
func Let[ST, S1, S2, T any](
func LetTo[ST, S1, S2, T any](
func Map[S, A, B any](f func(A) B) Operator[S, A, B]
type Pair[L, R any] = pair.Pair[L, R]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type State[S, A any] = state.State[S, A]
type StateIO[S, A any] = Reader[S, IO[Pair[S, A]]]
func Do[ST, A any](
func FromIO[S, A any](fa IO[A]) StateIO[S, A]
func MonadAp[B, S, A any](fab StateIO[S, func(A) B], fa StateIO[S, A]) StateIO[S, B]
func MonadChain[S, A, B any](fa StateIO[S, A], f Kleisli[S, A, B]) StateIO[S, B]
func MonadMap[S, A, B any](fa StateIO[S, A], f func(A) B) StateIO[S, B]
func Of[S, A any](a A) StateIO[S, A]
type StateIOApplicative[
func (o *StateIOApplicative[S, A, B]) Ap(fa StateIO[S, A]) Operator[S, func(A) B, B]
func (o *StateIOApplicative[S, A, B]) Map(f func(A) B) Operator[S, A, B]
func (o *StateIOApplicative[S, A, B]) Of(a A) StateIO[S, A]
type StateIOFunctor[
func (o *StateIOFunctor[S, A, B]) Map(f func(A) B) Operator[S, A, B]
type StateIOMonad[
func (o *StateIOMonad[S, A, B]) Ap(fa StateIO[S, A]) Operator[S, func(A) B, B]
func (o *StateIOMonad[S, A, B]) Chain(f Kleisli[S, A, B]) Operator[S, A, B]
func (o *StateIOMonad[S, A, B]) Map(f func(A) B) Operator[S, A, B]
func (o *StateIOMonad[S, A, B]) Of(a A) StateIO[S, A]
type StateIOPointed[
func (o *StateIOPointed[S, A]) Of(a A) StateIO[S, A]
```

## package `github.com/IBM/fp-go/v2/statereaderioeither`

Import: `import "github.com/IBM/fp-go/v2/statereaderioeither"`

StateReaderIOEither: `StateReaderIOEither[S, R, E, A] = func(S) ReaderIOEither[R, E, Pair[S, A]]`.

### Exported API

```go
func Applicative[
func Eq[
func FromStrictEquals[
func Functor[
func Local[S, E, A, B, R1, R2 any](f func(R2) R1) func(StateReaderIOEither[S, R1, E, A]) StateReaderIOEither[S, R2, E, A]
func Monad[
func Pointed[
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
func ApSL[ST, R, E, S, T any](
func BindL[ST, R, E, S, T any](
func LetL[ST, R, E, S, T any](
func LetToL[ST, R, E, S, T any](
type IO[A any] = io.IO[A]
type IOEither[E, A any] = ioeither.IOEither[E, A]
type Kleisli[S, R, E, A, B any] = Reader[A, StateReaderIOEither[S, R, E, B]]
func FromEitherK[S, R, E, A, B any](f either.Kleisli[E, A, B]) Kleisli[S, R, E, A, B]
func FromIOEitherK[
func FromIOK[S, R, E, A, B any](f func(A) IO[B]) Kleisli[S, R, E, A, B]
func FromReaderIOEitherK[S, R, E, A, B any](f readerioeither.Kleisli[R, E, A, B]) Kleisli[S, R, E, A, B]
func WithResource[A, S, R, E, RES, ANY any](
type Lens[S, A any] = lens.Lens[S, A]
type Operator[S, R, E, A, B any] = Reader[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]]
func Ap[B, S, R, E, A any](fa StateReaderIOEither[S, R, E, A]) Operator[S, R, E, func(A) B, B]
func ApS[ST, R, E, S1, S2, T any](
func Bind[ST, R, E, S1, S2, T any](
func BindTo[ST, R, E, S1, T any](
func Chain[S, R, E, A, B any](f Kleisli[S, R, E, A, B]) Operator[S, R, E, A, B]
func ChainEitherK[S, R, E, A, B any](f either.Kleisli[E, A, B]) Operator[S, R, E, A, B]
func ChainIOEitherK[S, R, E, A, B any](f ioeither.Kleisli[E, A, B]) Operator[S, R, E, A, B]
func ChainReaderIOEitherK[S, R, E, A, B any](f readerioeither.Kleisli[R, E, A, B]) Operator[S, R, E, A, B]
func FilterOrElse[S, R, E, A any](pred Predicate[A], onFalse func(A) E) Operator[S, R, E, A, A]
func Let[ST, R, E, S1, S2, T any](
func LetTo[ST, R, E, S1, S2, T any](
func Map[S, R, E, A, B any](f func(A) B) Operator[S, R, E, A, B]
type Pair[L, R any] = pair.Pair[L, R]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderEither[R, E, A any] = readereither.ReaderEither[R, E, A]
type ReaderIOEither[R, E, A any] = readerioeither.ReaderIOEither[R, E, A]
type State[S, A any] = state.State[S, A]
type StateReaderIOEither[S, R, E, A any] = Reader[S, ReaderIOEither[R, E, Pair[S, A]]]
func Asks[
func Do[ST, R, E, A any](
func FromEither[S, R, E, A any](ma Either[E, A]) StateReaderIOEither[S, R, E, A]
func FromIO[S, R, E, A any](fa IO[A]) StateReaderIOEither[S, R, E, A]
func FromIOEither[S, R, E, A any](fa IOEither[E, A]) StateReaderIOEither[S, R, E, A]
func FromReader[S, E, R, A any](fa Reader[R, A]) StateReaderIOEither[S, R, E, A]
func FromReaderEither[S, R, E, A any](fa ReaderEither[R, E, A]) StateReaderIOEither[S, R, E, A]
func FromReaderIOEither[S, R, E, A any](fa ReaderIOEither[R, E, A]) StateReaderIOEither[S, R, E, A]
func FromState[R, E, S, A any](sa State[S, A]) StateReaderIOEither[S, R, E, A]
func Left[S, R, A, E any](e E) StateReaderIOEither[S, R, E, A]
func MonadAp[B, S, R, E, A any](fab StateReaderIOEither[S, R, E, func(A) B], fa StateReaderIOEither[S, R, E, A]) StateReaderIOEither[S, R, E, B]
func MonadChain[S, R, E, A, B any](fa StateReaderIOEither[S, R, E, A], f Kleisli[S, R, E, A, B]) StateReaderIOEither[S, R, E, B]
func MonadChainEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f either.Kleisli[E, A, B]) StateReaderIOEither[S, R, E, B]
func MonadChainIOEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f ioeither.Kleisli[E, A, B]) StateReaderIOEither[S, R, E, B]
func MonadChainReaderIOEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f readerioeither.Kleisli[R, E, A, B]) StateReaderIOEither[S, R, E, B]
func MonadMap[S, R, E, A, B any](fa StateReaderIOEither[S, R, E, A], f func(A) B) StateReaderIOEither[S, R, E, B]
func Of[S, R, E, A any](a A) StateReaderIOEither[S, R, E, A]
func Right[S, R, E, A any](a A) StateReaderIOEither[S, R, E, A]
```

---

# Optics

## package `github.com/IBM/fp-go/v2/optics/lens`

Import: `import "github.com/IBM/fp-go/v2/optics/lens"`

Lens: composable getters/setters for immutable data.

- `Lens[S, A]` -- focuses on field A within structure S

### Exported API

```go
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(Lens[S, A]) Endomorphism[S]
func ModifyF[S, A, HKTA, HKTS any](
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type Kleisli[S, A, B any] = func(A) Lens[S, B]
type Lens[S, A any] struct {
func Id[S any]() Lens[S, S]
func IdRef[S any]() Lens[*S, *S]
func MakeLens[GET ~func(S) A, SET ~func(S, A) S, S, A any](get GET, set SET) Lens[S, A]
func MakeLensCurried[GET ~func(S) A, SET ~func(A) Endomorphism[S], S, A any](get GET, set SET) Lens[S, A]
func MakeLensCurriedRefWithName[GET ~func(*S) A, SET ~func(A) Endomorphism[*S], S, A any](get GET, set SET, name string) Lens[*S, A]
func MakeLensCurriedWithName[GET ~func(S) A, SET ~func(A) Endomorphism[S], S, A any](get GET, set SET, name string) Lens[S, A]
func MakeLensRef[GET ~func(*S) A, SET func(*S, A) *S, S, A any](get GET, set SET) Lens[*S, A]
func MakeLensRefCurried[S, A any](get func(*S) A, set func(A) Endomorphism[*S]) Lens[*S, A]
func MakeLensRefCurriedWithName[S, A any](get func(*S) A, set func(A) Endomorphism[*S], name string) Lens[*S, A]
func MakeLensRefWithName[GET ~func(*S) A, SET func(*S, A) *S, S, A any](get GET, set SET, name string) Lens[*S, A]
func MakeLensStrict[GET ~func(*S) A, SET func(*S, A) *S, S any, A comparable](get GET, set SET) Lens[*S, A]
func MakeLensStrictWithName[GET ~func(*S) A, SET func(*S, A) *S, S any, A comparable](get GET, set SET, name string) Lens[*S, A]
func MakeLensWithEq[GET ~func(*S) A, SET func(*S, A) *S, S, A any](pred EQ.Eq[A], get GET, set SET) Lens[*S, A]
func MakeLensWithEqWithName[GET ~func(*S) A, SET func(*S, A) *S, S, A any](pred EQ.Eq[A], get GET, set SET, name string) Lens[*S, A]
func MakeLensWithName[GET ~func(S) A, SET ~func(S, A) S, S, A any](get GET, set SET, name string) Lens[S, A]
func (l Lens[S, T]) Format(f fmt.State, c rune)
func (l Lens[S, T]) GoString() string
func (l Lens[S, T]) LogValue() slog.Value
func (l Lens[S, T]) String() string
type Operator[S, A, B any] = Kleisli[S, Lens[S, A], B]
func Compose[S, A, B any](ab Lens[A, B]) Operator[S, A, B]
func ComposeRef[S, A, B any](ab Lens[A, B]) Operator[*S, A, B]
func IMap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](ab AB, ba BA) Operator[S, A, B]
```

## package `github.com/IBM/fp-go/v2/optics/prism`

Import: `import "github.com/IBM/fp-go/v2/optics/prism"`

Prism: access a variant of a sum type.

- `Prism[S, A]` -- focuses on case A of sum type S

### Exported API

```go
func AsTraversal[R ~func(func(A) HKTA) func(S) HKTS, S, A, HKTS, HKTA any](
func Set[S, A any](a A) func(Prism[S, A]) Endomorphism[S]
type Either[E, T any] = either.Either[E, T]
type Endomorphism[T any] = endomorphism.Endomorphism[T]
type ErrorPrisms struct {
func MakeErrorPrisms() ErrorPrisms
type Kleisli[S, A, B any] = func(A) Prism[S, B]
type Lens[S, A any] = lens.Lens[S, A]
type Match struct {
func (m Match) FullMatch() string
func (m Match) Group(n int) string
func (m Match) Reconstruct() string
type NamedMatch struct {
func (nm NamedMatch) Reconstruct() string
type Operator[S, A, B any] = func(Prism[S, A]) Prism[S, B]
func Compose[S, A, B any](ab Prism[A, B]) Operator[S, A, B]
func IMap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](ab AB, ba BA) Operator[S, A, B]
type Option[T any] = O.Option[T]
type Predicate[A any] = predicate.Predicate[A]
type Prism[S, A any] struct {
func Deref[T any]() Prism[*T, *T]
func FromEither[E, T any]() Prism[Either[E, T], T]
func FromEncoding(enc *base64.Encoding) Prism[string, []byte]
func FromNonZero[T comparable]() Prism[T, T]
func FromOption[T any]() Prism[Option[T], T]
func FromPredicate[S any](pred func(S) bool) Prism[S, S]
func FromResult[T any]() Prism[Result[T], T]
func FromZero[T comparable]() Prism[T, T]
func Id[S any]() Prism[S, S]
func InstanceOf[T any]() Prism[any, T]
func MakePrism[S, A any](get O.Kleisli[S, A], rev func(A) S) Prism[S, A]
func MakePrismWithName[S, A any](get O.Kleisli[S, A], rev func(A) S, name string) Prism[S, A]
func NonEmptyString() Prism[string, string]
func ParseBool() Prism[string, bool]
func ParseDate(layout string) Prism[string, time.Time]
func ParseFloat32() Prism[string, float32]
func ParseFloat64() Prism[string, float64]
func ParseInt() Prism[string, int]
func ParseInt64() Prism[string, int64]
func ParseJSON[A any]() Prism[[]byte, A]
func ParseURL() Prism[string, *url.URL]
func RegexMatcher(re *regexp.Regexp) Prism[string, Match]
func RegexNamedMatcher(re *regexp.Regexp) Prism[string, NamedMatch]
func Some[S, A any](soa Prism[S, Option[A]]) Prism[S, A]
func (p Prism[S, T]) Format(f fmt.State, c rune)
func (p Prism[S, T]) GoString() string
func (p Prism[S, T]) LogValue() slog.Value
func (p Prism[S, T]) String() string
type Reader[R, T any] = reader.Reader[R, T]
type Result[T any] = result.Result[T]
type URLPrisms struct {
func MakeURLPrisms() URLPrisms
```

## package `github.com/IBM/fp-go/v2/optics/iso`

Import: `import "github.com/IBM/fp-go/v2/optics/iso"`

Iso: bidirectional lossless conversion.

- `Iso[S, A]` -- S and A are isomorphic

### Exported API

```go
func Compose[S, A, B any](ab Iso[A, B]) func(Iso[S, A]) Iso[S, B]
func From[S, A any](a A) func(Iso[S, A]) S
func IMap[S, A, B any](ab func(A) B, ba func(B) A) func(Iso[S, A]) Iso[S, B]
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(Iso[S, A]) EM.Endomorphism[S]
func To[A, S any](s S) func(Iso[S, A]) A
func Unwrap[A, S any](s S) func(Iso[S, A]) A
func Wrap[S, A any](a A) func(Iso[S, A]) S
type Either[E, A any] = either.Either[E, A]
type Iso[S, A any] struct {
func Add[T Number](n T) Iso[T, T]
func Head[A any]() Iso[A, NonEmptyArray[A]]
func Id[S any]() Iso[S, S]
func Lines() Iso[[]string, string]
func MakeIso[S, A any](get func(S) A, reverse func(A) S) Iso[S, A]
func Reverse[S, A any](sa Iso[S, A]) Iso[A, S]
func ReverseArray[A any]() Iso[[]A, []A]
func Sub[T Number](n T) Iso[T, T]
func SwapEither[E, A any]() Iso[Either[E, A], Either[A, E]]
func SwapPair[A, B any]() Iso[Pair[A, B], Pair[B, A]]
func UTF8String() Iso[[]byte, string]
func UnixMilli() Iso[int64, time.Time]
func (i Iso[S, T]) Format(f fmt.State, c rune)
func (i Iso[S, T]) GoString() string
func (i Iso[S, T]) LogValue() slog.Value
func (i Iso[S, T]) String() string
type NonEmptyArray[A any] = nonempty.NonEmptyArray[A]
type Number = number.Number
type Pair[A, B any] = pair.Pair[A, B]
```

## package `github.com/IBM/fp-go/v2/optics/optional`

Import: `import "github.com/IBM/fp-go/v2/optics/optional"`

Optional: access a value that may not exist.

- `Optional[S, A]` -- like Lens but focus may be absent

### Exported API

```go
func FromPredicate[S, A any](pred func(A) bool) func(func(S) A, func(S, A) S) Optional[S, A]
func FromPredicateRef[S, A any](pred func(A) bool) func(func(*S) A, func(*S, A) *S) Optional[*S, A]
func ModifyOption[S, A any](f func(A) A) func(Optional[S, A]) O.Kleisli[S, S]
func SetOption[S, A any](a A) func(Optional[S, A]) O.Kleisli[S, S]
type Kleisli[S, A, B any] = func(A) Optional[S, B]
type Operator[S, A, B any] = func(Optional[S, A]) Optional[S, B]
func Compose[S, A, B any](ab Optional[A, B]) Operator[S, A, B]
func ComposeRef[S, A, B any](ab Optional[A, B]) Operator[*S, A, B]
func IChain[S, A, B any](ab O.Kleisli[A, B], ba O.Kleisli[B, A]) Operator[S, A, B]
func IChainAny[S, A any]() Operator[S, any, A]
func IMap[S, A, B any](ab func(A) B, ba func(B) A) Operator[S, A, B]
type Optional[S, A any] struct {
func Id[S any]() Optional[S, S]
func IdRef[S any]() Optional[*S, *S]
func MakeOptional[S, A any](get O.Kleisli[S, A], set func(S, A) S) Optional[S, A]
func MakeOptionalCurried[S, A any](get O.Kleisli[S, A], set func(A) func(S) S) Optional[S, A]
func MakeOptionalCurriedWithName[S, A any](get O.Kleisli[S, A], set func(A) func(S) S, name string) Optional[S, A]
func MakeOptionalRef[S, A any](get O.Kleisli[*S, A], set func(*S, A) *S) Optional[*S, A]
func MakeOptionalRefCurriedWithName[S, A any](get O.Kleisli[*S, A], set func(A) func(*S) *S, name string) Optional[*S, A]
func MakeOptionalRefWithName[S, A any](get O.Kleisli[*S, A], set func(*S, A) *S, name string) Optional[*S, A]
func MakeOptionalWithName[S, A any](get O.Kleisli[S, A], set func(S, A) S, name string) Optional[S, A]
func (o Optional[S, T]) Format(f fmt.State, c rune)
func (o Optional[S, T]) GoString() string
func (o Optional[S, T]) LogValue() slog.Value
func (o Optional[S, T]) String() string
```

## package `github.com/IBM/fp-go/v2/optics/traversal`

Import: `import "github.com/IBM/fp-go/v2/optics/traversal"`

Traversal: focuses on zero or more targets within a structure.

### Exported API

```go
func Compose[
func Fold[S, A any](sa G.Traversal[S, A, C.Const[A, S], C.Const[A, A]]) func(S) A
func FoldMap[M, S, A any](f func(A) M) func(sa G.Traversal[S, A, C.Const[M, S], C.Const[M, A]]) func(S) M
func GetAll[S, A any](s S) func(sa G.Traversal[S, A, C.Const[[]A, S], C.Const[[]A, A]]) []A
func Id[S, A any]() G.Traversal[S, S, A, A]
func Modify[S, A any](f func(A) A) func(sa G.Traversal[S, A, S, A]) func(S) S
func Set[S, A any](a A) func(sa G.Traversal[S, A, S, A]) func(S) S
```

## package `github.com/IBM/fp-go/v2/optics/codec`

Import: `import "github.com/IBM/fp-go/v2/optics/codec"`

Codec: combines encoding and decoding.

### Exported API

```go
type Codec[I, O, A any] struct {
type Context = validation.Context
type Decode[I, A any] = decode.Decode[I, A]
type Decoder[I, A any] interface {
type Encode[A, O any] = Reader[A, O]
type Encoder[A, O any] interface {
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type Formattable = formatting.Formattable
type Iso[S, A any] = iso.Iso[S, A]
type Kleisli[A, B, O, I any] = Reader[A, Type[B, O, I]]
type Lazy[A any] = lazy.Lazy[A]
type Lens[S, A any] = lens.Lens[S, A]
type Monoid[A any] = monoid.Monoid[A]
func AltMonoid[A, O, I any](zero Lazy[Type[A, O, I]]) Monoid[Type[A, O, I]]
type Operator[A, B, O, I any] = Kleisli[Type[A, O, I], B, O, I]
func Alt[A, O, I any](second Lazy[Type[A, O, I]]) Operator[A, A, O, I]
func ApSL[S, T, O, I any](
func ApSO[S, T, O, I any](
func Bind[S, T, O, I any](
func Pipe[O, I, A, B any](ab Type[B, A, A]) Operator[A, B, O, I]
type Option[A any] = option.Option[A]
type Optional[S, A any] = optional.Optional[S, A]
type Pair[L, R any] = pair.Pair[L, R]
type Prism[S, A any] = prism.Prism[S, A]
func TypeToPrism[S, A any](t Type[A, S, S]) Prism[S, A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderResult[R, A any] = readerresult.ReaderResult[R, A]
func Is[T any]() ReaderResult[any, T]
type Refinement[A, B any] = Prism[A, B]
type Result[A any] = result.Result[A]
type Semigroup[A any] = semigroup.Semigroup[A]
type Type[A, O, I any] interface {
func Array[T, O any](item Type[T, O, any]) Type[[]T, []O, any]
func Bool() Type[bool, bool, any]
func BoolFromString() Type[bool, string, string]
func Date(layout string) Type[time.Time, string, string]
func Do[I, A, O any](e Lazy[Pair[O, A]]) Type[A, O, I]
func Either[A, B, O, I any](
func Empty[I, A, O any](e Lazy[Pair[O, A]]) Type[A, O, I]
func FromIso[A, I any](iso Iso[I, A]) Type[A, I, I]
func FromRefinement[A, B any](refinement Refinement[A, B]) Type[B, A, A]
func Id[T any]() Type[T, T, T]
func Int() Type[int, int, any]
func Int64FromString() Type[int64, string, string]
func IntFromString() Type[int, string, string]
func MakeSimpleType[A any]() Type[A, A, any]
func MakeType[A, O, I any](
func MarshalJSON[T any](
func MarshalText[T any](
func MonadAlt[A, O, I any](first Type[A, O, I], second Lazy[Type[A, O, I]]) Type[A, O, I]
func Nil[A any]() Type[*A, *A, any]
func Regex(re *regexp.Regexp) Type[prism.Match, string, string]
func RegexNamed(re *regexp.Regexp) Type[prism.NamedMatch, string, string]
func String() Type[string, string, any]
func TranscodeArray[T, O, I any](item Type[T, O, I]) Type[[]T, []O, []I]
func TranscodeEither[L, R, OL, OR, IL, IR any](leftItem Type[L, OL, IL], rightItem Type[R, OR, IR]) Type[either.Either[L, R], either.Either[OL, OR], either.Either[IL, IR]]
func URL() Type[*url.URL, string, string]
type Validate[I, A any] = validate.Validate[I, A]
type Validation[A any] = validation.Validation[A]
type Void = function.Void
```

---

# Utilities

## package `github.com/IBM/fp-go/v2/function`

Import: `import F "github.com/IBM/fp-go/v2/function"`

Core functional utilities.

Key exports:
- `Pipe1..Pipe20` -- left-to-right function application
- `Flow1..Flow20` -- left-to-right function composition  
- `Curry2..Curry5` / `Uncurry2..Uncurry5`
- `Bind1st/Bind2nd` -- partial application
- `Identity / Constant / Flip / Swap`
- `Ref / Deref / IsNil` -- pointer utilities
- `Ternary / Switch` -- conditional logic
- `Memoize` -- caching

### Exported API

```go
var ConstFalse = Constant(false)
var ConstTrue = Constant(true)
func Bind1234of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) func() R
func Bind123of3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T1, T2, T3) func() R
func Bind123of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1, T2, T3) func(T4) R
func Bind124of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1, T2, T4) func(T3) R
func Bind12of2[F ~func(T1, T2) R, T1, T2, R any](f F) func(T1, T2) func() R
func Bind12of3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T1, T2) func(T3) R
func Bind12of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1, T2) func(T3, T4) R
func Bind134of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1, T3, T4) func(T2) R
func Bind13of3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T1, T3) func(T2) R
func Bind13of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1, T3) func(T2, T4) R
func Bind14of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1, T4) func(T2, T3) R
func Bind1of1[F ~func(T1) R, T1, R any](f F) func(T1) func() R
func Bind1of2[F ~func(T1, T2) R, T1, T2, R any](f F) func(T1) func(T2) R
func Bind1of3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T1) func(T2, T3) R
func Bind1of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T1) func(T2, T3, T4) R
func Bind1st[T1, T2, R any](f func(T1, T2) R, t1 T1) func(T2) R
func Bind234of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T2, T3, T4) func(T1) R
func Bind23of3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T2, T3) func(T1) R
func Bind23of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T2, T3) func(T1, T4) R
func Bind24of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T2, T4) func(T1, T3) R
func Bind2nd[T1, T2, R any](f func(T1, T2) R, t2 T2) func(T1) R
func Bind2of2[F ~func(T1, T2) R, T1, T2, R any](f F) func(T2) func(T1) R
func Bind2of3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T2) func(T1, T3) R
func Bind2of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T2) func(T1, T3, T4) R
func Bind34of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T3, T4) func(T1, T2) R
func Bind3of3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T3) func(T1, T2) R
func Bind3of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T3) func(T1, T2, T4) R
func Bind4of4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(T4) func(T1, T2, T3) R
func CacheCallback[
func ConstNil[A any]() *A
func Constant[A any](a A) func() A
func Constant1[B, A any](a A) func(B) A
func Constant2[B, C, A any](a A) func(B, C) A
func ContramapMemoize[T, A any, K comparable](kf func(A) K) func(func(A) T) func(A) T
func Curry1[FCT ~func(T0) T1, T0, T1 any](f FCT) func(T0) T1
func Curry10[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) T10, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) T10
func Curry11[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) T11, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) T11
func Curry12[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) T12, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) T12
func Curry13[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) T13, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) T13
func Curry14[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) T14, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) T14
func Curry15[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) T15, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) T15
func Curry16[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15) T16, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) T16
func Curry17[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16) T17, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) T17
func Curry18[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17) T18, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) func(T17) T18
func Curry19[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18) T19, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) func(T17) func(T18) T19
func Curry2[FCT ~func(T0, T1) T2, T0, T1, T2 any](f FCT) func(T0) func(T1) T2
func Curry20[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19) T20, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) func(T17) func(T18) func(T19) T20
func Curry3[FCT ~func(T0, T1, T2) T3, T0, T1, T2, T3 any](f FCT) func(T0) func(T1) func(T2) T3
func Curry4[FCT ~func(T0, T1, T2, T3) T4, T0, T1, T2, T3, T4 any](f FCT) func(T0) func(T1) func(T2) func(T3) T4
func Curry5[FCT ~func(T0, T1, T2, T3, T4) T5, T0, T1, T2, T3, T4, T5 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) T5
func Curry6[FCT ~func(T0, T1, T2, T3, T4, T5) T6, T0, T1, T2, T3, T4, T5, T6 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) T6
func Curry7[FCT ~func(T0, T1, T2, T3, T4, T5, T6) T7, T0, T1, T2, T3, T4, T5, T6, T7 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) T7
func Curry8[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7) T8, T0, T1, T2, T3, T4, T5, T6, T7, T8 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) T8
func Curry9[FCT ~func(T0, T1, T2, T3, T4, T5, T6, T7, T8) T9, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](f FCT) func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) T9
func Deref[A any](a *A) A
func First[T1, T2 any](t1 T1, _ T2) T1
func Flip[T1, T2, R any](f func(T1) func(T2) R) func(T2) func(T1) R
func Flow1[F1 ~func(T0) T1, T0, T1 any](f1 F1) func(T0) T1
func Flow10[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(T0) T10
func Flow11[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) func(T0) T11
func Flow12[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) func(T0) T12
func Flow13[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) func(T0) T13
func Flow14[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) func(T0) T14
func Flow15[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) func(T0) T15
func Flow16[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16) func(T0) T16
func Flow17[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17) func(T0) T17
func Flow18[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18) func(T0) T18
func Flow19[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, F19 ~func(T18) T19, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18, f19 F19) func(T0) T19
func Flow2[F1 ~func(T0) T1, F2 ~func(T1) T2, T0, T1, T2 any](f1 F1, f2 F2) func(T0) T2
func Flow20[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, F19 ~func(T18) T19, F20 ~func(T19) T20, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18, f19 F19, f20 F20) func(T0) T20
func Flow3[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, T0, T1, T2, T3 any](f1 F1, f2 F2, f3 F3) func(T0) T3
func Flow4[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, T0, T1, T2, T3, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T0) T4
func Flow5[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, T0, T1, T2, T3, T4, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T0) T5
func Flow6[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, T0, T1, T2, T3, T4, T5, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(T0) T6
func Flow7[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, T0, T1, T2, T3, T4, T5, T6, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(T0) T7
func Flow8[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, T0, T1, T2, T3, T4, T5, T6, T7, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(T0) T8
func Flow9[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(T0) T9
func Identity[A any](a A) A
func Ignore1234of4[T1, T2, T3, T4 any, F ~func() R, R any](f F) func(T1, T2, T3, T4) R
func Ignore123of3[T1, T2, T3 any, F ~func() R, R any](f F) func(T1, T2, T3) R
func Ignore123of4[T1, T2, T3 any, F ~func(T4) R, T4, R any](f F) func(T1, T2, T3, T4) R
func Ignore124of4[T1, T2, T4 any, F ~func(T3) R, T3, R any](f F) func(T1, T2, T3, T4) R
func Ignore12of2[T1, T2 any, F ~func() R, R any](f F) func(T1, T2) R
func Ignore12of3[T1, T2 any, F ~func(T3) R, T3, R any](f F) func(T1, T2, T3) R
func Ignore12of4[T1, T2 any, F ~func(T3, T4) R, T3, T4, R any](f F) func(T1, T2, T3, T4) R
func Ignore134of4[T1, T3, T4 any, F ~func(T2) R, T2, R any](f F) func(T1, T2, T3, T4) R
func Ignore13of3[T1, T3 any, F ~func(T2) R, T2, R any](f F) func(T1, T2, T3) R
func Ignore13of4[T1, T3 any, F ~func(T2, T4) R, T2, T4, R any](f F) func(T1, T2, T3, T4) R
func Ignore14of4[T1, T4 any, F ~func(T2, T3) R, T2, T3, R any](f F) func(T1, T2, T3, T4) R
func Ignore1of1[T1 any, F ~func() R, R any](f F) func(T1) R
func Ignore1of2[T1 any, F ~func(T2) R, T2, R any](f F) func(T1, T2) R
func Ignore1of3[T1 any, F ~func(T2, T3) R, T2, T3, R any](f F) func(T1, T2, T3) R
func Ignore1of4[T1 any, F ~func(T2, T3, T4) R, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) R
func Ignore234of4[T2, T3, T4 any, F ~func(T1) R, T1, R any](f F) func(T1, T2, T3, T4) R
func Ignore23of3[T2, T3 any, F ~func(T1) R, T1, R any](f F) func(T1, T2, T3) R
func Ignore23of4[T2, T3 any, F ~func(T1, T4) R, T1, T4, R any](f F) func(T1, T2, T3, T4) R
func Ignore24of4[T2, T4 any, F ~func(T1, T3) R, T1, T3, R any](f F) func(T1, T2, T3, T4) R
func Ignore2of2[T2 any, F ~func(T1) R, T1, R any](f F) func(T1, T2) R
func Ignore2of3[T2 any, F ~func(T1, T3) R, T1, T3, R any](f F) func(T1, T2, T3) R
func Ignore2of4[T2 any, F ~func(T1, T3, T4) R, T1, T3, T4, R any](f F) func(T1, T2, T3, T4) R
func Ignore34of4[T3, T4 any, F ~func(T1, T2) R, T1, T2, R any](f F) func(T1, T2, T3, T4) R
func Ignore3of3[T3 any, F ~func(T1, T2) R, T1, T2, R any](f F) func(T1, T2, T3) R
func Ignore3of4[T3 any, F ~func(T1, T2, T4) R, T1, T2, T4, R any](f F) func(T1, T2, T3, T4) R
func Ignore4of4[T4 any, F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(T1, T2, T3, T4) R
func IsNil[A any](a *A) bool
func IsNonNil[A any](a *A) bool
func Memoize[K comparable, T any](f func(K) T) func(K) T
func Nullary1[F1 ~func() T1, T1 any](f1 F1) func() T1
func Nullary10[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func() T10
func Nullary11[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) func() T11
func Nullary12[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) func() T12
func Nullary13[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) func() T13
func Nullary14[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) func() T14
func Nullary15[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) func() T15
func Nullary16[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16) func() T16
func Nullary17[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17) func() T17
func Nullary18[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18) func() T18
func Nullary19[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, F19 ~func(T18) T19, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18, f19 F19) func() T19
func Nullary2[F1 ~func() T1, F2 ~func(T1) T2, T1, T2 any](f1 F1, f2 F2) func() T2
func Nullary20[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, F19 ~func(T18) T19, F20 ~func(T19) T20, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18, f19 F19, f20 F20) func() T20
func Nullary3[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, T1, T2, T3 any](f1 F1, f2 F2, f3 F3) func() T3
func Nullary4[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, T1, T2, T3, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func() T4
func Nullary5[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, T1, T2, T3, T4, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func() T5
func Nullary6[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, T1, T2, T3, T4, T5, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func() T6
func Nullary7[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, T1, T2, T3, T4, T5, T6, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func() T7
func Nullary8[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, T1, T2, T3, T4, T5, T6, T7, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func() T8
func Nullary9[F1 ~func() T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func() T9
func Pipe0[T0 any](t0 T0) T0
func Pipe1[F1 ~func(T0) T1, T0, T1 any](t0 T0, f1 F1) T1
func Pipe10[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) T10
func Pipe11[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) T11
func Pipe12[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) T12
func Pipe13[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) T13
func Pipe14[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) T14
func Pipe15[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) T15
func Pipe16[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16) T16
func Pipe17[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17) T17
func Pipe18[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18) T18
func Pipe19[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, F19 ~func(T18) T19, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18, f19 F19) T19
func Pipe2[F1 ~func(T0) T1, F2 ~func(T1) T2, T0, T1, T2 any](t0 T0, f1 F1, f2 F2) T2
func Pipe20[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, F10 ~func(T9) T10, F11 ~func(T10) T11, F12 ~func(T11) T12, F13 ~func(T12) T13, F14 ~func(T13) T14, F15 ~func(T14) T15, F16 ~func(T15) T16, F17 ~func(T16) T17, F18 ~func(T17) T18, F19 ~func(T18) T19, F20 ~func(T19) T20, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15, f16 F16, f17 F17, f18 F18, f19 F19, f20 F20) T20
func Pipe3[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, T0, T1, T2, T3 any](t0 T0, f1 F1, f2 F2, f3 F3) T3
func Pipe4[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, T0, T1, T2, T3, T4 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4) T4
func Pipe5[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, T0, T1, T2, T3, T4, T5 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) T5
func Pipe6[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, T0, T1, T2, T3, T4, T5, T6 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) T6
func Pipe7[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, T0, T1, T2, T3, T4, T5, T6, T7 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) T7
func Pipe8[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, T0, T1, T2, T3, T4, T5, T6, T7, T8 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) T8
func Pipe9[F1 ~func(T0) T1, F2 ~func(T1) T2, F3 ~func(T2) T3, F4 ~func(T3) T4, F5 ~func(T4) T5, F6 ~func(T5) T6, F7 ~func(T6) T7, F8 ~func(T7) T8, F9 ~func(T8) T9, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) T9
func Ref[A any](a A) *A
func SK[T1, T2 any](_ T1, t2 T2) T2
func Second[T1, T2 any](_ T1, t2 T2) T2
func SingleElementCache[K comparable, T any]() func(K, func() func() T) func() T
func Swap[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R
func Switch[K comparable, T, R any](kf func(T) K, n map[K]func(T) R, d func(T) R) func(T) R
func Ternary[A, B any](pred func(A) bool, onTrue, onFalse func(A) B) func(A) B
func ToAny[A any](a A) any
func Uncurry1[FCT ~func(T0) T1, T0, T1 any](f FCT) func(T0) T1
func Uncurry10[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) T10, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) T10
func Uncurry11[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) T11, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) T11
func Uncurry12[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) T12, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) T12
func Uncurry13[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) T13, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) T13
func Uncurry14[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) T14, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) T14
func Uncurry15[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) T15, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) T15
func Uncurry16[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) T16, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15) T16
func Uncurry17[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) T17, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16) T17
func Uncurry18[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) func(T17) T18, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17) T18
func Uncurry19[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) func(T17) func(T18) T19, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18) T19
func Uncurry2[FCT ~func(T0) func(T1) T2, T0, T1, T2 any](f FCT) func(T0, T1) T2
func Uncurry20[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) func(T9) func(T10) func(T11) func(T12) func(T13) func(T14) func(T15) func(T16) func(T17) func(T18) func(T19) T20, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19) T20
func Uncurry3[FCT ~func(T0) func(T1) func(T2) T3, T0, T1, T2, T3 any](f FCT) func(T0, T1, T2) T3
func Uncurry4[FCT ~func(T0) func(T1) func(T2) func(T3) T4, T0, T1, T2, T3, T4 any](f FCT) func(T0, T1, T2, T3) T4
func Uncurry5[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) T5, T0, T1, T2, T3, T4, T5 any](f FCT) func(T0, T1, T2, T3, T4) T5
func Uncurry6[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) T6, T0, T1, T2, T3, T4, T5, T6 any](f FCT) func(T0, T1, T2, T3, T4, T5) T6
func Uncurry7[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) T7, T0, T1, T2, T3, T4, T5, T6, T7 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6) T7
func Uncurry8[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) T8, T0, T1, T2, T3, T4, T5, T6, T7, T8 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7) T8
func Uncurry9[FCT ~func(T0) func(T1) func(T2) func(T3) func(T4) func(T5) func(T6) func(T7) func(T8) T9, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](f FCT) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) T9
func Unsliced0[F ~func([]T) R, T, R any](f F) func() R
func Unsliced1[F ~func([]T) R, T, R any](f F) func(T) R
func Unsliced10[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T) R
func Unsliced11[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced12[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced13[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced14[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced15[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced16[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced17[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced18[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced19[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced2[F ~func([]T) R, T, R any](f F) func(T, T) R
func Unsliced20[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) R
func Unsliced3[F ~func([]T) R, T, R any](f F) func(T, T, T) R
func Unsliced4[F ~func([]T) R, T, R any](f F) func(T, T, T, T) R
func Unsliced5[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T) R
func Unsliced6[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T) R
func Unsliced7[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T) R
func Unsliced8[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T) R
func Unsliced9[F ~func([]T) R, T, R any](f F) func(T, T, T, T, T, T, T, T, T) R
func Unvariadic0[V, R any](f func(...V) R) func([]V) R
func Unvariadic1[T1, V, R any](f func(T1, ...V) R) func(T1, []V) R
func Unvariadic10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, []V) R
func Unvariadic11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, []V) R
func Unvariadic12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, []V) R
func Unvariadic13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, []V) R
func Unvariadic14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, []V) R
func Unvariadic15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, []V) R
func Unvariadic16[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, []V) R
func Unvariadic17[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, []V) R
func Unvariadic18[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, []V) R
func Unvariadic19[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, []V) R
func Unvariadic2[T1, T2, V, R any](f func(T1, T2, ...V) R) func(T1, T2, []V) R
func Unvariadic20[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, []V) R
func Unvariadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, ...V) R) func(T1, T2, T3, []V) R
func Unvariadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, ...V) R) func(T1, T2, T3, T4, []V) R
func Unvariadic5[T1, T2, T3, T4, T5, V, R any](f func(T1, T2, T3, T4, T5, ...V) R) func(T1, T2, T3, T4, T5, []V) R
func Unvariadic6[T1, T2, T3, T4, T5, T6, V, R any](f func(T1, T2, T3, T4, T5, T6, ...V) R) func(T1, T2, T3, T4, T5, T6, []V) R
func Unvariadic7[T1, T2, T3, T4, T5, T6, T7, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, []V) R
func Unvariadic8[T1, T2, T3, T4, T5, T6, T7, T8, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, []V) R
func Unvariadic9[T1, T2, T3, T4, T5, T6, T7, T8, T9, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, ...V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, []V) R
func Variadic0[V, R any](f func([]V) R) func(...V) R
func Variadic1[T1, V, R any](f func(T1, []V) R) func(T1, ...V) R
func Variadic10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, ...V) R
func Variadic11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, ...V) R
func Variadic12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, ...V) R
func Variadic13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, ...V) R
func Variadic14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, ...V) R
func Variadic15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, ...V) R
func Variadic16[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, ...V) R
func Variadic17[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, ...V) R
func Variadic18[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, ...V) R
func Variadic19[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, ...V) R
func Variadic2[T1, T2, V, R any](f func(T1, T2, []V) R) func(T1, T2, ...V) R
func Variadic20[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, ...V) R
func Variadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, []V) R) func(T1, T2, T3, ...V) R
func Variadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, []V) R) func(T1, T2, T3, T4, ...V) R
func Variadic5[T1, T2, T3, T4, T5, V, R any](f func(T1, T2, T3, T4, T5, []V) R) func(T1, T2, T3, T4, T5, ...V) R
func Variadic6[T1, T2, T3, T4, T5, T6, V, R any](f func(T1, T2, T3, T4, T5, T6, []V) R) func(T1, T2, T3, T4, T5, T6, ...V) R
func Variadic7[T1, T2, T3, T4, T5, T6, T7, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, []V) R) func(T1, T2, T3, T4, T5, T6, T7, ...V) R
func Variadic8[T1, T2, T3, T4, T5, T6, T7, T8, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, ...V) R
func Variadic9[T1, T2, T3, T4, T5, T6, T7, T8, T9, V, R any](f func(T1, T2, T3, T4, T5, T6, T7, T8, T9, []V) R) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, ...V) R
func Zero[A any]() A
type Void = struct{}
var VOID Void = struct{}{}
```

## package `github.com/IBM/fp-go/v2/array`

Import: `import A "github.com/IBM/fp-go/v2/array"`

Array (slice) operations: Map, Filter, Reduce, Sort, Find, Traverse, etc.

### Exported API

```go
func Any[A any](pred func(A) bool) func([]A) bool
func AnyWithIndex[A any](pred func(int, A) bool) func([]A) bool
func Append[A any](as []A, a A) []A
func ArrayConcatAll[A any](data ...[]A) []A
func ConcatAll[A any](m M.Monoid[A]) func([]A) A
func ConstNil[A any]() []A
func Copy[A any](b []A) []A
func Do[S any](
func Empty[A any]() []A
func Eq[T any](e E.Eq[T]) E.Eq[[]T]
func Extract[A any](as []A) A
func FindFirst[A any](pred func(A) bool) option.Kleisli[[]A, A]
func FindFirstMap[A, B any](sel option.Kleisli[A, B]) option.Kleisli[[]A, B]
func FindFirstMapWithIndex[A, B any](sel func(int, A) Option[B]) option.Kleisli[[]A, B]
func FindFirstWithIndex[A any](pred func(int, A) bool) option.Kleisli[[]A, A]
func FindLast[A any](pred func(A) bool) option.Kleisli[[]A, A]
func FindLastMap[A, B any](sel option.Kleisli[A, B]) option.Kleisli[[]A, B]
func FindLastMapWithIndex[A, B any](sel func(int, A) Option[B]) option.Kleisli[[]A, B]
func FindLastWithIndex[A any](pred func(int, A) bool) option.Kleisli[[]A, A]
func Flatten[A any](mma [][]A) []A
func Fold[A any](m M.Monoid[A]) func([]A) A
func FoldMap[A, B any](m M.Monoid[B]) func(func(A) B) func([]A) B
func FoldMapWithIndex[A, B any](m M.Monoid[B]) func(func(int, A) B) func([]A) B
func From[A any](data ...A) []A
func Intercalate[A any](m M.Monoid[A]) func(A) func([]A) A
func IsEmpty[A any](as []A) bool
func IsNil[A any](as []A) bool
func IsNonEmpty[A any](as []A) bool
func IsNonNil[A any](as []A) bool
func Lookup[A any](idx int) func([]A) Option[A]
func MakeBy[F ~func(int) A, A any](n int, f F) []A
func MakeTraverseType[A, B, HKT_F_B, HKT_F_T_B, HKT_F_B_T_B any]() traversable.TraverseType[A, B, []A, []B, HKT_F_B, HKT_F_T_B, HKT_F_B_T_B]
func Match[A, B any](onEmpty func() B, onNonEmpty func([]A) B) func([]A) B
func MatchLeft[A, B any](onEmpty func() B, onNonEmpty func(A, []A) B) func([]A) B
func Monad[A, B any]() monad.Monad[A, B, []A, []B, []func(A) B]
func MonadAp[B, A any](fab []func(A) B, fa []A) []B
func MonadChain[A, B any](fa []A, f Kleisli[A, B]) []B
func MonadFilterMap[A, B any](fa []A, f option.Kleisli[A, B]) []B
func MonadFilterMapWithIndex[A, B any](fa []A, f func(int, A) Option[B]) []B
func MonadFlap[B, A any](fab []func(A) B, a A) []B
func MonadMap[A, B any](as []A, f func(A) B) []B
func MonadMapRef[A, B any](as []A, f func(*A) B) []B
func MonadPartition[A any](as []A, pred func(A) bool) pair.Pair[[]A, []A]
func MonadReduce[A, B any](fa []A, f func(B, A) B, initial B) B
func MonadReduceWithIndex[A, B any](fa []A, f func(int, B, A) B, initial B) B
func MonadSequence[HKTA, HKTRA any](
func MonadTraverse[A, B, HKTB, HKTAB, HKTRB any](
func MonadTraverseWithIndex[A, B, HKTB, HKTAB, HKTRB any](
func Monoid[T any]() M.Monoid[[]T]
func Of[A any](a A) []A
func Partition[A any](pred func(A) bool) func([]A) pair.Pair[[]A, []A]
func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B
func ReduceRef[A, B any](f func(B, *A) B, initial B) func([]A) B
func ReduceRight[A, B any](f func(A, B) B, initial B) func([]A) B
func ReduceRightWithIndex[A, B any](f func(int, A, B) B, initial B) func([]A) B
func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func([]A) B
func Replicate[A any](n int, a A) []A
func Reverse[A any](as []A) []A
func Semigroup[T any]() S.Semigroup[[]T]
func Sequence[HKTA, HKTRA any](
func Size[A any](as []A) int
func StrictEquals[T comparable]() E.Eq[[]T]
func StrictUniq[A comparable](as []A) []A
func Traverse[A, B, HKTB, HKTAB, HKTRB any](
func TraverseWithIndex[A, B, HKTB, HKTAB, HKTRB any](
func Unzip[A, B any](cs []pair.Pair[A, B]) pair.Pair[[]A, []B]
func Zero[A any]() []A
func Zip[A, B any](fb []B) func([]A) []pair.Pair[A, B]
func ZipWith[FCT ~func(A, B) C, A, B, C any](fa []A, fb []B, f FCT) []C
type Kleisli[A, B any] = func(A) []B
type Operator[A, B any] = Kleisli[[]A, B]
func Ap[B, A any](fa []A) Operator[func(A) B, B]
func ApS[S1, S2, T any](
func Bind[S1, S2, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainOptionK[A, B any](f option.Kleisli[A, []B]) Operator[A, B]
func Clone[A any](f func(A) A) Operator[A, A]
func Concat[A any](suffix []A) Operator[A, A]
func Extend[A, B any](f func([]A) B) Operator[A, B]
func Filter[A any](pred func(A) bool) Operator[A, A]
func FilterMap[A, B any](f option.Kleisli[A, B]) Operator[A, B]
func FilterMapRef[A, B any](pred func(a *A) bool, f func(*A) B) Operator[A, B]
func FilterMapWithIndex[A, B any](f func(int, A) Option[B]) Operator[A, B]
func FilterRef[A any](pred func(*A) bool) Operator[A, A]
func FilterWithIndex[A any](pred func(int, A) bool) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Intersperse[A any](middle A) Operator[A, A]
func Let[S1, S2, T any](
func LetTo[S1, S2, T any](
func Map[A, B any](f func(A) B) Operator[A, B]
func MapRef[A, B any](f func(*A) B) Operator[A, B]
func MapWithIndex[A, B any](f func(int, A) B) Operator[A, B]
func Prepend[A any](head A) Operator[A, A]
func PrependAll[A any](middle A) Operator[A, A]
func Push[A any](a A) Operator[A, A]
func Slice[A any](low, high int) Operator[A, A]
func SliceRight[A any](start int) Operator[A, A]
func Sort[T any](ord O.Ord[T]) Operator[T, T]
func SortBy[T any](ord []O.Ord[T]) Operator[T, T]
func SortByKey[K, T any](ord O.Ord[K], f func(T) K) Operator[T, T]
func Uniq[A any, K comparable](f func(A) K) Operator[A, A]
func UpsertAt[A any](a A) Operator[A, A]
type Option[A any] = option.Option[A]
func ArrayOption[A any](ma []Option[A]) Option[[]A]
func First[A any](as []A) Option[A]
func Head[A any](as []A) Option[A]
func Last[A any](as []A) Option[A]
func Tail[A any](as []A) Option[[]A]
```

## package `github.com/IBM/fp-go/v2/record`

Import: `import "github.com/IBM/fp-go/v2/record"`

Record (map) operations: Map, Filter, Reduce, Traverse, etc.

### Exported API

```go
func Ap[A any, K comparable, B any](m Monoid[Record[K, B]]) func(fa Record[K, A]) Operator[K, func(A) B, B]
func ApS[S1, T any, K comparable, S2 any](m Monoid[Record[K, S2]]) func(
func Bind[S1, T any, K comparable, S2 any](m Monoid[Record[K, S2]]) func(
func Chain[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(Kleisli[K, V1, V2]) Operator[K, V1, V2]
func ChainWithIndex[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(KleisliWithIndex[K, V1, V2]) Operator[K, V1, V2]
func Collect[K comparable, V, R any](f func(K, V) R) func(Record[K, V]) []R
func CollectOrd[V, R any, K comparable](o ord.Ord[K]) func(func(K, V) R) func(Record[K, V]) []R
func Eq[K comparable, V any](e E.Eq[V]) E.Eq[Record[K, V]]
func FilterChain[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(option.Kleisli[V1, Record[K, V2]]) Operator[K, V1, V2]
func FilterChainWithIndex[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(func(K, V1) Option[Record[K, V2]]) Operator[K, V1, V2]
func Flatten[K comparable, V any](m Monoid[Record[K, V]]) func(Record[K, Record[K, V]]) Record[K, V]
func Fold[K comparable, A any](m Monoid[A]) func(Record[K, A]) A
func FoldMap[K comparable, A, B any](m Monoid[B]) func(func(A) B) func(Record[K, A]) B
func FoldMapOrd[A, B any, K comparable](o ord.Ord[K]) func(m Monoid[B]) func(func(A) B) func(Record[K, A]) B
func FoldMapOrdWithIndex[K comparable, A, B any](o ord.Ord[K]) func(m Monoid[B]) func(func(K, A) B) func(Record[K, A]) B
func FoldMapWithIndex[K comparable, A, B any](m Monoid[B]) func(func(K, A) B) func(Record[K, A]) B
func FoldOrd[A any, K comparable](o ord.Ord[K]) func(m Monoid[A]) func(Record[K, A]) A
func FromArrayMap[
func FromFoldableMap[
func FromStrictEquals[K, V comparable]() E.Eq[Record[K, V]]
func Has[K comparable, V any](k K, r Record[K, V]) bool
func IsEmpty[K comparable, V any](r Record[K, V]) bool
func IsNil[K comparable, V any](m Record[K, V]) bool
func IsNonEmpty[K comparable, V any](r Record[K, V]) bool
func IsNonNil[K comparable, V any](m Record[K, V]) bool
func Keys[K comparable, V any](r Record[K, V]) []K
func KeysOrd[V any, K comparable](o ord.Ord[K]) func(r Record[K, V]) []K
func Lookup[V any, K comparable](k K) option.Kleisli[Record[K, V], V]
func Reduce[K comparable, V, R any](f func(R, V) R, initial R) func(Record[K, V]) R
func ReduceOrd[V, R any, K comparable](o ord.Ord[K]) func(func(R, V) R, R) func(Record[K, V]) R
func ReduceOrdWithIndex[V, R any, K comparable](o ord.Ord[K]) func(func(K, R, V) R, R) func(Record[K, V]) R
func ReduceRef[K comparable, V, R any](f func(R, *V) R, initial R) func(Record[K, V]) R
func ReduceRefWithIndex[K comparable, V, R any](f func(K, R, *V) R, initial R) func(Record[K, V]) R
func ReduceWithIndex[K comparable, V, R any](f func(K, R, V) R, initial R) func(Record[K, V]) R
func Sequence[K comparable, A, HKTA, HKTAA, HKTRA any](
func Size[K comparable, V any](r Record[K, V]) int
func Traverse[K comparable, A, B, HKTB, HKTAB, HKTRB any](
func TraverseWithIndex[K comparable, A, B, HKTB, HKTAB, HKTRB any](
func Union[K comparable, V any](m Mg.Magma[V]) func(Record[K, V]) Operator[K, V, V]
func Values[K comparable, V any](r Record[K, V]) []V
func ValuesOrd[V any, K comparable](o ord.Ord[K]) func(r Record[K, V]) []V
type Collector[K comparable, V, R any] = func(K, V) R
type Endomorphism[A any] = endomorphism.Endomorphism[A]
func Clone[K comparable, V any](f Endomorphism[V]) Endomorphism[Record[K, V]]
type Entries[K comparable, V any] = []Entry[K, V]
func ToArray[K comparable, V any](r Record[K, V]) Entries[K, V]
func ToEntries[K comparable, V any](r Record[K, V]) Entries[K, V]
type Entry[K comparable, V any] = pair.Pair[K, V]
type Kleisli[K comparable, V1, V2 any] = func(V1) Record[K, V2]
func FromArray[
func FromFoldable[
type KleisliWithIndex[K comparable, V1, V2 any] = func(K, V1) Record[K, V2]
type Monoid[A any] = monoid.Monoid[A]
func MergeMonoid[K comparable, V any]() Monoid[Record[K, V]]
func UnionFirstMonoid[K comparable, V any]() Monoid[Record[K, V]]
func UnionLastMonoid[K comparable, V any]() Monoid[Record[K, V]]
func UnionMonoid[K comparable, V any](s S.Semigroup[V]) Monoid[Record[K, V]]
type Operator[K comparable, V1, V2 any] = func(Record[K, V1]) Record[K, V2]
func BindTo[S1, T any, K comparable](setter func(T) S1) Operator[K, T, S1]
func DeleteAt[K comparable, V any](k K) Operator[K, V, V]
func Filter[K comparable, V any](f Predicate[K]) Operator[K, V, V]
func FilterMap[K comparable, V1, V2 any](f option.Kleisli[V1, V2]) Operator[K, V1, V2]
func FilterMapWithIndex[K comparable, V1, V2 any](f func(K, V1) Option[V2]) Operator[K, V1, V2]
func FilterWithIndex[K comparable, V any](f PredicateWithIndex[K, V]) Operator[K, V, V]
func Flap[B any, K comparable, A any](a A) Operator[K, func(A) B, B]
func Let[S1, T any, K comparable, S2 any](
func LetTo[S1, T any, K comparable, S2 any](
func Map[K comparable, V, R any](f func(V) R) Operator[K, V, R]
func MapRef[K comparable, V, R any](f func(*V) R) Operator[K, V, R]
func MapRefWithIndex[K comparable, V, R any](f func(K, *V) R) Operator[K, V, R]
func MapWithIndex[K comparable, V, R any](f func(K, V) R) Operator[K, V, R]
func Merge[K comparable, V any](right Record[K, V]) Operator[K, V, V]
func UpsertAt[K comparable, V any](k K, v V) Operator[K, V, V]
type OperatorWithIndex[K comparable, V1, V2 any] = func(func(K, V1) V2) Operator[K, V1, V2]
type Option[A any] = option.Option[A]
func MonadLookup[V any, K comparable](m Record[K, V], k K) Option[V]
type Predicate[K any] = predicate.Predicate[K]
type PredicateWithIndex[K comparable, V any] = func(K, V) bool
type Record[K comparable, V any] = map[K]V
func ConstNil[K comparable, V any]() Record[K, V]
func Copy[K comparable, V any](m Record[K, V]) Record[K, V]
func Do[K comparable, S any]() Record[K, S]
func Empty[K comparable, V any]() Record[K, V]
func FromEntries[K comparable, V any](fa Entries[K, V]) Record[K, V]
func MonadAp[A any, K comparable, B any](m Monoid[Record[K, B]], fab Record[K, func(A) B], fa Record[K, A]) Record[K, B]
func MonadChain[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]], r Record[K, V1], f Kleisli[K, V1, V2]) Record[K, V2]
func MonadChainWithIndex[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]], r Record[K, V1], f KleisliWithIndex[K, V1, V2]) Record[K, V2]
func MonadFlap[B any, K comparable, A any](fab Record[K, func(A) B], a A) Record[K, B]
func MonadMap[K comparable, V, R any](r Record[K, V], f func(V) R) Record[K, R]
func MonadMapRef[K comparable, V, R any](r Record[K, V], f func(*V) R) Record[K, R]
func MonadMapRefWithIndex[K comparable, V, R any](r Record[K, V], f func(K, *V) R) Record[K, R]
func MonadMapWithIndex[K comparable, V, R any](r Record[K, V], f func(K, V) R) Record[K, R]
func Of[K comparable, A any](k K, a A) Record[K, A]
func Singleton[K comparable, V any](k K, v V) Record[K, V]
type Reducer[V, R any] = func(R, V) R
type ReducerWithIndex[K comparable, V, R any] = func(K, R, V) R
type Semigroup[A any] = semigroup.Semigroup[A]
func UnionFirstSemigroup[K comparable, V any]() Semigroup[Record[K, V]]
func UnionLastSemigroup[K comparable, V any]() Semigroup[Record[K, V]]
func UnionSemigroup[K comparable, V any](s Semigroup[V]) Semigroup[Record[K, V]]
```

## package `github.com/IBM/fp-go/v2/pair`

Import: `import P "github.com/IBM/fp-go/v2/pair"`

Pair[L, R]: generic 2-element product type.

### Exported API

```go
func ApHead[B, A, A1 any](sg Semigroup[B], fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B]
func Applicative[B, A, B1 any](m monoid.Monoid[A]) applicative.Applicative[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]
func ApplicativeHead[A, B, A1 any](m monoid.Monoid[B]) applicative.Applicative[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]]
func ApplicativeMonoid[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]]
func ApplicativeMonoidHead[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]]
func ApplicativeMonoidTail[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]]
func ApplicativeTail[B, A, B1 any](m monoid.Monoid[A]) applicative.Applicative[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]
func BiMap[A, B, A1, B1 any](f func(A) A1, g func(B) B1) func(Pair[A, B]) Pair[A1, B1]
func ChainHead[B, A, A1 any](sg Semigroup[B], f func(A) Pair[A1, B]) func(Pair[A, B]) Pair[A1, B]
func Eq[A, B any](a eq.Eq[A], b eq.Eq[B]) eq.Eq[Pair[A, B]]
func First[A, B any](fa Pair[A, B]) A
func FromStrictEquals[A, B comparable]() eq.Eq[Pair[A, B]]
func Functor[B, A, B1 any]() functor.Functor[B, B1, Pair[A, B], Pair[A, B1]]
func FunctorHead[A, B, A1 any]() functor.Functor[A, A1, Pair[A, B], Pair[A1, B]]
func FunctorTail[B, A, B1 any]() functor.Functor[B, B1, Pair[A, B], Pair[A, B1]]
func Head[A, B any](fa Pair[A, B]) A
func MapHead[B, A, A1 any](f func(A) A1) func(Pair[A, B]) Pair[A1, B]
func Merge[F ~func(B) func(A) R, A, B, R any](f F) func(Pair[A, B]) R
func Monad[B, A, B1 any](m monoid.Monoid[A]) monad.Monad[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]
func MonadHead[A, B, A1 any](m monoid.Monoid[B]) monad.Monad[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]]
func MonadSequence[L, A, HKTA, HKTPA any](
func MonadTail[B, A, B1 any](m monoid.Monoid[A]) monad.Monad[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]
func MonadTraverse[L, A, HKTA, HKTPA any](
func Monoid[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]]
func Paired[F ~func(T1, T2) R, T1, T2, R any](f F) func(Pair[T1, T2]) R
func Pointed[B, A any](m monoid.Monoid[A]) pointed.Pointed[B, Pair[A, B]]
func PointedHead[A, B any](m monoid.Monoid[B]) pointed.Pointed[A, Pair[A, B]]
func PointedTail[B, A any](m monoid.Monoid[A]) pointed.Pointed[B, Pair[A, B]]
func Second[A, B any](fa Pair[A, B]) B
func Sequence[L, A, HKTA, HKTPA any](
func Tail[A, B any](fa Pair[A, B]) B
func Traverse[L, A, HKTA, HKTPA any](
func Unpack[L, R any](p Pair[L, R]) (L, R)
func Unpaired[F ~func(Pair[T1, T2]) R, T1, T2, R any](f F) func(T1, T2) R
type Kleisli[L, R1, R2 any] = func(R1) Pair[L, R2]
func FromHead[B, A any](a A) Kleisli[A, B, B]
func FromTail[A, B any](b B) Kleisli[A, A, B]
type Operator[L, R1, R2 any] = func(Pair[L, R1]) Pair[L, R2]
func Ap[A, B, B1 any](sg Semigroup[A], fa Pair[A, B]) Operator[A, func(B) B1, B1]
func ApTail[A, B, B1 any](sg Semigroup[A], fb Pair[A, B]) Operator[A, func(B) B1, B1]
func Chain[A, B, B1 any](sg Semigroup[A], f Kleisli[A, B, B1]) Operator[A, B, B1]
func ChainTail[A, B, B1 any](sg Semigroup[A], f Kleisli[A, B, B1]) Operator[A, B, B1]
func Map[A, B, B1 any](f func(B) B1) Operator[A, B, B1]
func MapTail[A, B, B1 any](f func(B) B1) Operator[A, B, B1]
type Pair[L, R any] struct {
func FromTuple[A, B any](t Tuple2[A, B]) Pair[A, B]
func MakePair[A, B any](a A, b B) Pair[A, B]
func MonadAp[A, B, B1 any](sg Semigroup[A], faa Pair[A, func(B) B1], fa Pair[A, B]) Pair[A, B1]
func MonadApHead[B, A, A1 any](sg Semigroup[B], faa Pair[func(A) A1, B], fa Pair[A, B]) Pair[A1, B]
func MonadApTail[A, B, B1 any](sg Semigroup[A], fbb Pair[A, func(B) B1], fb Pair[A, B]) Pair[A, B1]
func MonadBiMap[A, B, A1, B1 any](fa Pair[A, B], f func(A) A1, g func(B) B1) Pair[A1, B1]
func MonadChain[A, B, B1 any](sg Semigroup[A], fa Pair[A, B], f Kleisli[A, B, B1]) Pair[A, B1]
func MonadChainHead[B, A, A1 any](sg Semigroup[B], fa Pair[A, B], f func(A) Pair[A1, B]) Pair[A1, B]
func MonadChainTail[A, B, B1 any](sg Semigroup[A], fb Pair[A, B], f Kleisli[A, B, B1]) Pair[A, B1]
func MonadMap[A, B, B1 any](fa Pair[A, B], f func(B) B1) Pair[A, B1]
func MonadMapHead[B, A, A1 any](fa Pair[A, B], f func(A) A1) Pair[A1, B]
func MonadMapTail[A, B, B1 any](fa Pair[A, B], f func(B) B1) Pair[A, B1]
func Of[A any](value A) Pair[A, A]
func Swap[A, B any](fa Pair[A, B]) Pair[B, A]
func Zero[L, R any]() Pair[L, R]
func (p Pair[L, R]) Format(f fmt.State, c rune)
func (p Pair[L, R]) GoString() string
func (p Pair[L, R]) LogValue() slog.Value
func (p Pair[L, R]) String() string
type Semigroup[A any] = semigroup.Semigroup[A]
type Tuple2[A, B any] = tuple.Tuple2[A, B]
func ToTuple[A, B any](t Pair[A, B]) Tuple2[A, B]
```

## package `github.com/IBM/fp-go/v2/tuple`

Import: `import T "github.com/IBM/fp-go/v2/tuple"`

Tuple types from Tuple1 to Tuple15. Fixed-size heterogeneous containers.

### Exported API

```go
func BiMap[E, G, A, B any](mapSnd func(E) G, mapFst func(A) B) func(Tuple2[A, E]) Tuple2[B, G]
func First[T1, T2 any](t Tuple2[T1, T2]) T1
func FromArray1[F1 ~func(R) T1, T1, R any](f1 F1) func(r []R) Tuple1[T1]
func FromArray10[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, F9 ~func(R) T9, F10 ~func(R) T10, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(r []R) Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]
func FromArray11[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, F9 ~func(R) T9, F10 ~func(R) T10, F11 ~func(R) T11, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) func(r []R) Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]
func FromArray12[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, F9 ~func(R) T9, F10 ~func(R) T10, F11 ~func(R) T11, F12 ~func(R) T12, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) func(r []R) Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]
func FromArray13[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, F9 ~func(R) T9, F10 ~func(R) T10, F11 ~func(R) T11, F12 ~func(R) T12, F13 ~func(R) T13, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) func(r []R) Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]
func FromArray14[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, F9 ~func(R) T9, F10 ~func(R) T10, F11 ~func(R) T11, F12 ~func(R) T12, F13 ~func(R) T13, F14 ~func(R) T14, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) func(r []R) Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]
func FromArray15[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, F9 ~func(R) T9, F10 ~func(R) T10, F11 ~func(R) T11, F12 ~func(R) T12, F13 ~func(R) T13, F14 ~func(R) T14, F15 ~func(R) T15, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) func(r []R) Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]
func FromArray2[F1 ~func(R) T1, F2 ~func(R) T2, T1, T2, R any](f1 F1, f2 F2) func(r []R) Tuple2[T1, T2]
func FromArray3[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, T1, T2, T3, R any](f1 F1, f2 F2, f3 F3) func(r []R) Tuple3[T1, T2, T3]
func FromArray4[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, T1, T2, T3, T4, R any](f1 F1, f2 F2, f3 F3, f4 F4) func(r []R) Tuple4[T1, T2, T3, T4]
func FromArray5[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, T1, T2, T3, T4, T5, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(r []R) Tuple5[T1, T2, T3, T4, T5]
func FromArray6[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, T1, T2, T3, T4, T5, T6, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(r []R) Tuple6[T1, T2, T3, T4, T5, T6]
func FromArray7[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, T1, T2, T3, T4, T5, T6, T7, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(r []R) Tuple7[T1, T2, T3, T4, T5, T6, T7]
func FromArray8[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, T1, T2, T3, T4, T5, T6, T7, T8, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(r []R) Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
func FromArray9[F1 ~func(R) T1, F2 ~func(R) T2, F3 ~func(R) T3, F4 ~func(R) T4, F5 ~func(R) T5, F6 ~func(R) T6, F7 ~func(R) T7, F8 ~func(R) T8, F9 ~func(R) T9, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(r []R) Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
func Map1[F1 ~func(T1) R1, T1, R1 any](f1 F1) func(Tuple1[T1]) Tuple1[R1]
func Map10[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, F9 ~func(T9) R9, F10 ~func(T10) R10, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8, T9, R9, T10, R10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) Tuple10[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10]
func Map11[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, F9 ~func(T9) R9, F10 ~func(T10) R10, F11 ~func(T11) R11, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8, T9, R9, T10, R10, T11, R11 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) func(Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) Tuple11[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10, R11]
func Map12[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, F9 ~func(T9) R9, F10 ~func(T10) R10, F11 ~func(T11) R11, F12 ~func(T12) R12, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8, T9, R9, T10, R10, T11, R11, T12, R12 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) func(Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) Tuple12[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10, R11, R12]
func Map13[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, F9 ~func(T9) R9, F10 ~func(T10) R10, F11 ~func(T11) R11, F12 ~func(T12) R12, F13 ~func(T13) R13, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8, T9, R9, T10, R10, T11, R11, T12, R12, T13, R13 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) func(Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) Tuple13[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10, R11, R12, R13]
func Map14[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, F9 ~func(T9) R9, F10 ~func(T10) R10, F11 ~func(T11) R11, F12 ~func(T12) R12, F13 ~func(T13) R13, F14 ~func(T14) R14, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8, T9, R9, T10, R10, T11, R11, T12, R12, T13, R13, T14, R14 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) func(Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) Tuple14[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10, R11, R12, R13, R14]
func Map15[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, F9 ~func(T9) R9, F10 ~func(T10) R10, F11 ~func(T11) R11, F12 ~func(T12) R12, F13 ~func(T13) R13, F14 ~func(T14) R14, F15 ~func(T15) R15, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8, T9, R9, T10, R10, T11, R11, T12, R12, T13, R13, T14, R14, T15, R15 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) func(Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]) Tuple15[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10, R11, R12, R13, R14, R15]
func Map2[F1 ~func(T1) R1, F2 ~func(T2) R2, T1, R1, T2, R2 any](f1 F1, f2 F2) func(Tuple2[T1, T2]) Tuple2[R1, R2]
func Map3[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, T1, R1, T2, R2, T3, R3 any](f1 F1, f2 F2, f3 F3) func(Tuple3[T1, T2, T3]) Tuple3[R1, R2, R3]
func Map4[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, T1, R1, T2, R2, T3, R3, T4, R4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(Tuple4[T1, T2, T3, T4]) Tuple4[R1, R2, R3, R4]
func Map5[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(Tuple5[T1, T2, T3, T4, T5]) Tuple5[R1, R2, R3, R4, R5]
func Map6[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(Tuple6[T1, T2, T3, T4, T5, T6]) Tuple6[R1, R2, R3, R4, R5, R6]
func Map7[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(Tuple7[T1, T2, T3, T4, T5, T6, T7]) Tuple7[R1, R2, R3, R4, R5, R6, R7]
func Map8[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Tuple8[R1, R2, R3, R4, R5, R6, R7, R8]
func Map9[F1 ~func(T1) R1, F2 ~func(T2) R2, F3 ~func(T3) R3, F4 ~func(T4) R4, F5 ~func(T5) R5, F6 ~func(T6) R6, F7 ~func(T7) R7, F8 ~func(T8) R8, F9 ~func(T9) R9, T1, R1, T2, R2, T3, R3, T4, R4, T5, R5, T6, R6, T7, R7, T8, R8, T9, R9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Tuple9[R1, R2, R3, R4, R5, R6, R7, R8, R9]
func Monoid1[T1 any](m1 M.Monoid[T1]) M.Monoid[Tuple1[T1]]
func Monoid10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8], m9 M.Monoid[T9], m10 M.Monoid[T10]) M.Monoid[Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func Monoid11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8], m9 M.Monoid[T9], m10 M.Monoid[T10], m11 M.Monoid[T11]) M.Monoid[Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func Monoid12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8], m9 M.Monoid[T9], m10 M.Monoid[T10], m11 M.Monoid[T11], m12 M.Monoid[T12]) M.Monoid[Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func Monoid13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8], m9 M.Monoid[T9], m10 M.Monoid[T10], m11 M.Monoid[T11], m12 M.Monoid[T12], m13 M.Monoid[T13]) M.Monoid[Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func Monoid14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8], m9 M.Monoid[T9], m10 M.Monoid[T10], m11 M.Monoid[T11], m12 M.Monoid[T12], m13 M.Monoid[T13], m14 M.Monoid[T14]) M.Monoid[Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func Monoid15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8], m9 M.Monoid[T9], m10 M.Monoid[T10], m11 M.Monoid[T11], m12 M.Monoid[T12], m13 M.Monoid[T13], m14 M.Monoid[T14], m15 M.Monoid[T15]) M.Monoid[Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func Monoid2[T1, T2 any](m1 M.Monoid[T1], m2 M.Monoid[T2]) M.Monoid[Tuple2[T1, T2]]
func Monoid3[T1, T2, T3 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3]) M.Monoid[Tuple3[T1, T2, T3]]
func Monoid4[T1, T2, T3, T4 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4]) M.Monoid[Tuple4[T1, T2, T3, T4]]
func Monoid5[T1, T2, T3, T4, T5 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5]) M.Monoid[Tuple5[T1, T2, T3, T4, T5]]
func Monoid6[T1, T2, T3, T4, T5, T6 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6]) M.Monoid[Tuple6[T1, T2, T3, T4, T5, T6]]
func Monoid7[T1, T2, T3, T4, T5, T6, T7 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7]) M.Monoid[Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func Monoid8[T1, T2, T3, T4, T5, T6, T7, T8 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8]) M.Monoid[Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func Monoid9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4], m5 M.Monoid[T5], m6 M.Monoid[T6], m7 M.Monoid[T7], m8 M.Monoid[T8], m9 M.Monoid[T9]) M.Monoid[Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Of2[T1, T2 any](e T2) func(T1) Tuple2[T1, T2]
func Ord1[T1 any](o1 O.Ord[T1]) O.Ord[Tuple1[T1]]
func Ord10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8], o9 O.Ord[T9], o10 O.Ord[T10]) O.Ord[Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func Ord11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8], o9 O.Ord[T9], o10 O.Ord[T10], o11 O.Ord[T11]) O.Ord[Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]]
func Ord12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8], o9 O.Ord[T9], o10 O.Ord[T10], o11 O.Ord[T11], o12 O.Ord[T12]) O.Ord[Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]]
func Ord13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8], o9 O.Ord[T9], o10 O.Ord[T10], o11 O.Ord[T11], o12 O.Ord[T12], o13 O.Ord[T13]) O.Ord[Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]]
func Ord14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8], o9 O.Ord[T9], o10 O.Ord[T10], o11 O.Ord[T11], o12 O.Ord[T12], o13 O.Ord[T13], o14 O.Ord[T14]) O.Ord[Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]]
func Ord15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8], o9 O.Ord[T9], o10 O.Ord[T10], o11 O.Ord[T11], o12 O.Ord[T12], o13 O.Ord[T13], o14 O.Ord[T14], o15 O.Ord[T15]) O.Ord[Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]]
func Ord2[T1, T2 any](o1 O.Ord[T1], o2 O.Ord[T2]) O.Ord[Tuple2[T1, T2]]
func Ord3[T1, T2, T3 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3]) O.Ord[Tuple3[T1, T2, T3]]
func Ord4[T1, T2, T3, T4 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4]) O.Ord[Tuple4[T1, T2, T3, T4]]
func Ord5[T1, T2, T3, T4, T5 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5]) O.Ord[Tuple5[T1, T2, T3, T4, T5]]
func Ord6[T1, T2, T3, T4, T5, T6 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6]) O.Ord[Tuple6[T1, T2, T3, T4, T5, T6]]
func Ord7[T1, T2, T3, T4, T5, T6, T7 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7]) O.Ord[Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func Ord8[T1, T2, T3, T4, T5, T6, T7, T8 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8]) O.Ord[Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func Ord9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4], o5 O.Ord[T5], o6 O.Ord[T6], o7 O.Ord[T7], o8 O.Ord[T8], o9 O.Ord[T9]) O.Ord[Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func Push1[T1, T2 any](value T2) func(Tuple1[T1]) Tuple2[T1, T2]
func Push10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](value T11) func(Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]
func Push11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](value T12) func(Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]
func Push12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](value T13) func(Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]
func Push13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](value T14) func(Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]
func Push14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](value T15) func(Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]
func Push2[T1, T2, T3 any](value T3) func(Tuple2[T1, T2]) Tuple3[T1, T2, T3]
func Push3[T1, T2, T3, T4 any](value T4) func(Tuple3[T1, T2, T3]) Tuple4[T1, T2, T3, T4]
func Push4[T1, T2, T3, T4, T5 any](value T5) func(Tuple4[T1, T2, T3, T4]) Tuple5[T1, T2, T3, T4, T5]
func Push5[T1, T2, T3, T4, T5, T6 any](value T6) func(Tuple5[T1, T2, T3, T4, T5]) Tuple6[T1, T2, T3, T4, T5, T6]
func Push6[T1, T2, T3, T4, T5, T6, T7 any](value T7) func(Tuple6[T1, T2, T3, T4, T5, T6]) Tuple7[T1, T2, T3, T4, T5, T6, T7]
func Push7[T1, T2, T3, T4, T5, T6, T7, T8 any](value T8) func(Tuple7[T1, T2, T3, T4, T5, T6, T7]) Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
func Push8[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](value T9) func(Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
func Push9[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](value T10) func(Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]
func Second[T1, T2 any](t Tuple2[T1, T2]) T2
func ToArray1[F1 ~func(T1) R, T1, R any](f1 F1) func(t Tuple1[T1]) []R
func ToArray10[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, F9 ~func(T9) R, F10 ~func(T10) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(t Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) []R
func ToArray11[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, F9 ~func(T9) R, F10 ~func(T10) R, F11 ~func(T11) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11) func(t Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) []R
func ToArray12[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, F9 ~func(T9) R, F10 ~func(T10) R, F11 ~func(T11) R, F12 ~func(T12) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12) func(t Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) []R
func ToArray13[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, F9 ~func(T9) R, F10 ~func(T10) R, F11 ~func(T11) R, F12 ~func(T12) R, F13 ~func(T13) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13) func(t Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) []R
func ToArray14[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, F9 ~func(T9) R, F10 ~func(T10) R, F11 ~func(T11) R, F12 ~func(T12) R, F13 ~func(T13) R, F14 ~func(T14) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14) func(t Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) []R
func ToArray15[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, F9 ~func(T9) R, F10 ~func(T10) R, F11 ~func(T11) R, F12 ~func(T12) R, F13 ~func(T13) R, F14 ~func(T14) R, F15 ~func(T15) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10, f11 F11, f12 F12, f13 F13, f14 F14, f15 F15) func(t Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]) []R
func ToArray2[F1 ~func(T1) R, F2 ~func(T2) R, T1, T2, R any](f1 F1, f2 F2) func(t Tuple2[T1, T2]) []R
func ToArray3[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, T1, T2, T3, R any](f1 F1, f2 F2, f3 F3) func(t Tuple3[T1, T2, T3]) []R
func ToArray4[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, T1, T2, T3, T4, R any](f1 F1, f2 F2, f3 F3, f4 F4) func(t Tuple4[T1, T2, T3, T4]) []R
func ToArray5[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, T1, T2, T3, T4, T5, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(t Tuple5[T1, T2, T3, T4, T5]) []R
func ToArray6[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, T1, T2, T3, T4, T5, T6, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(t Tuple6[T1, T2, T3, T4, T5, T6]) []R
func ToArray7[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, T1, T2, T3, T4, T5, T6, T7, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(t Tuple7[T1, T2, T3, T4, T5, T6, T7]) []R
func ToArray8[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, T1, T2, T3, T4, T5, T6, T7, T8, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(t Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) []R
func ToArray9[F1 ~func(T1) R, F2 ~func(T2) R, F3 ~func(T3) R, F4 ~func(T4) R, F5 ~func(T5) R, F6 ~func(T6) R, F7 ~func(T7) R, F8 ~func(T8) R, F9 ~func(T9) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(t Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) []R
func Tupled1[F ~func(T1) R, T1, R any](f F) func(Tuple1[T1]) R
func Tupled10[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) R
func Tupled11[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f F) func(Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) R
func Tupled12[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f F) func(Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) R
func Tupled13[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f F) func(Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) R
func Tupled14[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f F) func(Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) R
func Tupled15[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, R any](f F) func(Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]) R
func Tupled2[F ~func(T1, T2) R, T1, T2, R any](f F) func(Tuple2[T1, T2]) R
func Tupled3[F ~func(T1, T2, T3) R, T1, T2, T3, R any](f F) func(Tuple3[T1, T2, T3]) R
func Tupled4[F ~func(T1, T2, T3, T4) R, T1, T2, T3, T4, R any](f F) func(Tuple4[T1, T2, T3, T4]) R
func Tupled5[F ~func(T1, T2, T3, T4, T5) R, T1, T2, T3, T4, T5, R any](f F) func(Tuple5[T1, T2, T3, T4, T5]) R
func Tupled6[F ~func(T1, T2, T3, T4, T5, T6) R, T1, T2, T3, T4, T5, T6, R any](f F) func(Tuple6[T1, T2, T3, T4, T5, T6]) R
func Tupled7[F ~func(T1, T2, T3, T4, T5, T6, T7) R, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(Tuple7[T1, T2, T3, T4, T5, T6, T7]) R
func Tupled8[F ~func(T1, T2, T3, T4, T5, T6, T7, T8) R, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) R
func Tupled9[F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) R
func Untupled1[F ~func(Tuple1[T1]) R, T1, R any](f F) func(T1) R
func Untupled10[F ~func(Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) R
func Untupled11[F ~func(Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) R
func Untupled12[F ~func(Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) R
func Untupled13[F ~func(Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) R
func Untupled14[F ~func(Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) R
func Untupled15[F ~func(Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15) R
func Untupled2[F ~func(Tuple2[T1, T2]) R, T1, T2, R any](f F) func(T1, T2) R
func Untupled3[F ~func(Tuple3[T1, T2, T3]) R, T1, T2, T3, R any](f F) func(T1, T2, T3) R
func Untupled4[F ~func(Tuple4[T1, T2, T3, T4]) R, T1, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) R
func Untupled5[F ~func(Tuple5[T1, T2, T3, T4, T5]) R, T1, T2, T3, T4, T5, R any](f F) func(T1, T2, T3, T4, T5) R
func Untupled6[F ~func(Tuple6[T1, T2, T3, T4, T5, T6]) R, T1, T2, T3, T4, T5, T6, R any](f F) func(T1, T2, T3, T4, T5, T6) R
func Untupled7[F ~func(Tuple7[T1, T2, T3, T4, T5, T6, T7]) R, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T1, T2, T3, T4, T5, T6, T7) R
func Untupled8[F ~func(Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) R, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8) R
func Untupled9[F ~func(Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) R, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9) R
type Tuple1[T1 any] struct {
func MakeTuple1[T1 any](t1 T1) Tuple1[T1]
func Of[T1 any](t T1) Tuple1[T1]
func Replicate1[T any](t T) Tuple1[T]
func (t Tuple1[T1]) MarshalJSON() ([]byte, error)
func (t Tuple1[T1]) String() string
func (t *Tuple1[T1]) UnmarshalJSON(data []byte) error
type Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any] struct {
func MakeTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10) Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]
func Replicate10[T any](t T) Tuple10[T, T, T, T, T, T, T, T, T, T]
func (t Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) MarshalJSON() ([]byte, error)
func (t Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) String() string
func (t *Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) UnmarshalJSON(data []byte) error
type Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any] struct {
func MakeTuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10, t11 T11) Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]
func Replicate11[T any](t T) Tuple11[T, T, T, T, T, T, T, T, T, T, T]
func (t Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) MarshalJSON() ([]byte, error)
func (t Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) String() string
func (t *Tuple11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) UnmarshalJSON(data []byte) error
type Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any] struct {
func MakeTuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10, t11 T11, t12 T12) Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]
func Replicate12[T any](t T) Tuple12[T, T, T, T, T, T, T, T, T, T, T, T]
func (t Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) MarshalJSON() ([]byte, error)
func (t Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) String() string
func (t *Tuple12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12]) UnmarshalJSON(data []byte) error
type Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any] struct {
func MakeTuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10, t11 T11, t12 T12, t13 T13) Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]
func Replicate13[T any](t T) Tuple13[T, T, T, T, T, T, T, T, T, T, T, T, T]
func (t Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) MarshalJSON() ([]byte, error)
func (t Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) String() string
func (t *Tuple13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13]) UnmarshalJSON(data []byte) error
type Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any] struct {
func MakeTuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10, t11 T11, t12 T12, t13 T13, t14 T14) Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]
func Replicate14[T any](t T) Tuple14[T, T, T, T, T, T, T, T, T, T, T, T, T, T]
func (t Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) MarshalJSON() ([]byte, error)
func (t Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) String() string
func (t *Tuple14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14]) UnmarshalJSON(data []byte) error
type Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any] struct {
func MakeTuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10, t11 T11, t12 T12, t13 T13, t14 T14, t15 T15) Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]
func Replicate15[T any](t T) Tuple15[T, T, T, T, T, T, T, T, T, T, T, T, T, T, T]
func (t Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]) MarshalJSON() ([]byte, error)
func (t Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]) String() string
func (t *Tuple15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15]) UnmarshalJSON(data []byte) error
type Tuple2[T1, T2 any] struct {
func MakeTuple2[T1, T2 any](t1 T1, t2 T2) Tuple2[T1, T2]
func Replicate2[T any](t T) Tuple2[T, T]
func Swap[T1, T2 any](t Tuple2[T1, T2]) Tuple2[T2, T1]
func (t Tuple2[T1, T2]) MarshalJSON() ([]byte, error)
func (t Tuple2[T1, T2]) String() string
func (t *Tuple2[T1, T2]) UnmarshalJSON(data []byte) error
type Tuple3[T1, T2, T3 any] struct {
func MakeTuple3[T1, T2, T3 any](t1 T1, t2 T2, t3 T3) Tuple3[T1, T2, T3]
func Replicate3[T any](t T) Tuple3[T, T, T]
func (t Tuple3[T1, T2, T3]) MarshalJSON() ([]byte, error)
func (t Tuple3[T1, T2, T3]) String() string
func (t *Tuple3[T1, T2, T3]) UnmarshalJSON(data []byte) error
type Tuple4[T1, T2, T3, T4 any] struct {
func MakeTuple4[T1, T2, T3, T4 any](t1 T1, t2 T2, t3 T3, t4 T4) Tuple4[T1, T2, T3, T4]
func Replicate4[T any](t T) Tuple4[T, T, T, T]
func (t Tuple4[T1, T2, T3, T4]) MarshalJSON() ([]byte, error)
func (t Tuple4[T1, T2, T3, T4]) String() string
func (t *Tuple4[T1, T2, T3, T4]) UnmarshalJSON(data []byte) error
type Tuple5[T1, T2, T3, T4, T5 any] struct {
func MakeTuple5[T1, T2, T3, T4, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) Tuple5[T1, T2, T3, T4, T5]
func Replicate5[T any](t T) Tuple5[T, T, T, T, T]
func (t Tuple5[T1, T2, T3, T4, T5]) MarshalJSON() ([]byte, error)
func (t Tuple5[T1, T2, T3, T4, T5]) String() string
func (t *Tuple5[T1, T2, T3, T4, T5]) UnmarshalJSON(data []byte) error
type Tuple6[T1, T2, T3, T4, T5, T6 any] struct {
func MakeTuple6[T1, T2, T3, T4, T5, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) Tuple6[T1, T2, T3, T4, T5, T6]
func Replicate6[T any](t T) Tuple6[T, T, T, T, T, T]
func (t Tuple6[T1, T2, T3, T4, T5, T6]) MarshalJSON() ([]byte, error)
func (t Tuple6[T1, T2, T3, T4, T5, T6]) String() string
func (t *Tuple6[T1, T2, T3, T4, T5, T6]) UnmarshalJSON(data []byte) error
type Tuple7[T1, T2, T3, T4, T5, T6, T7 any] struct {
func MakeTuple7[T1, T2, T3, T4, T5, T6, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) Tuple7[T1, T2, T3, T4, T5, T6, T7]
func Replicate7[T any](t T) Tuple7[T, T, T, T, T, T, T]
func (t Tuple7[T1, T2, T3, T4, T5, T6, T7]) MarshalJSON() ([]byte, error)
func (t Tuple7[T1, T2, T3, T4, T5, T6, T7]) String() string
func (t *Tuple7[T1, T2, T3, T4, T5, T6, T7]) UnmarshalJSON(data []byte) error
type Tuple8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
func MakeTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
func Replicate8[T any](t T) Tuple8[T, T, T, T, T, T, T, T]
func (t Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) MarshalJSON() ([]byte, error)
func (t Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) String() string
func (t *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) UnmarshalJSON(data []byte) error
type Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
func MakeTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
func Replicate9[T any](t T) Tuple9[T, T, T, T, T, T, T, T, T]
func (t Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) MarshalJSON() ([]byte, error)
func (t Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) String() string
func (t *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) UnmarshalJSON(data []byte) error
```

## package `github.com/IBM/fp-go/v2/predicate`

Import: `import "github.com/IBM/fp-go/v2/predicate"`

Predicate utilities.

- `Predicate[A] = func(A) bool`
- And, Or, Not, Contramap

### Exported API

```go
type Kleisli[A, B any] = func(A) Predicate[B]
func IsEqual[A any](pred eq.Eq[A]) Kleisli[A, A]
func IsStrictEqual[A comparable]() Kleisli[A, A]
type Monoid[A any] = monoid.Monoid[Predicate[A]]
func MonoidAll[A any]() Monoid[A]
func MonoidAny[A any]() Monoid[A]
type Operator[A, B any] = Kleisli[Predicate[A], B]
func And[A any](second Predicate[A]) Operator[A, A]
func ContraMap[A, B any](f func(B) A) Operator[A, B]
func Or[A any](second Predicate[A]) Operator[A, A]
type Predicate[A any] = func(A) bool
func IsNonZero[A comparable]() Predicate[A]
func IsZero[A comparable]() Predicate[A]
func Not[A any](predicate Predicate[A]) Predicate[A]
type Semigroup[A any] = semigroup.Semigroup[Predicate[A]]
func SemigroupAll[A any]() Semigroup[A]
func SemigroupAny[A any]() Semigroup[A]
```

## package `github.com/IBM/fp-go/v2/endomorphism`

Import: `import "github.com/IBM/fp-go/v2/endomorphism"`

Endomorphism[A] = func(A) A. Function from A to A.

### Exported API

```go
type of the endomorphism, preventing type mismatches at compile time.
func Build[A any](e Endomorphism[A]) A
func Curry2[FCT ~func(T0, T1) T1, T0, T1 any](f FCT) func(T0) Endomorphism[T1]
func Curry3[FCT ~func(T0, T1, T2) T2, T0, T1, T2 any](f FCT) func(T0) func(T1) Endomorphism[T2]
func Monoid[A any]() M.Monoid[Endomorphism[A]]
func Read[A any](a A) func(Endomorphism[A]) A
func Reduce[T any](es []Endomorphism[T]) T
func Semigroup[A any]() S.Semigroup[Endomorphism[A]]
func Unwrap[F ~func(A) A, A any](f Endomorphism[A]) F
type Endomorphism[A any] = func(A) A
func ConcatAll[T any](es []Endomorphism[T]) Endomorphism[T]
func Flatten[A any](mma Endomorphism[Endomorphism[A]]) Endomorphism[A]
func Identity[A any]() Endomorphism[A]
func Join[A any](f Kleisli[A]) Endomorphism[A]
func MonadAp[A any](fab, fa Endomorphism[A]) Endomorphism[A]
func MonadChain[A any](ma, f Endomorphism[A]) Endomorphism[A]
func MonadChainFirst[A any](ma, f Endomorphism[A]) Endomorphism[A]
func MonadCompose[A any](f, g Endomorphism[A]) Endomorphism[A]
func MonadMap[A any](f, ma Endomorphism[A]) Endomorphism[A]
func Of[F ~func(A) A, A any](f F) Endomorphism[A]
func Wrap[F ~func(A) A, A any](f F) Endomorphism[A]
type Kleisli[A any] = func(A) Endomorphism[A]
func FromSemigroup[A any](s S.Semigroup[A]) Kleisli[A]
type Operator[A any] = Endomorphism[Endomorphism[A]]
func Ap[A any](fa Endomorphism[A]) Operator[A]
func Chain[A any](f Endomorphism[A]) Operator[A]
func ChainFirst[A any](f Endomorphism[A]) Operator[A]
func Compose[A any](g Endomorphism[A]) Operator[A]
func Map[A any](f Endomorphism[A]) Operator[A]
```

---

# Algebraic Structures

## package `github.com/IBM/fp-go/v2/eq`

Import: `import "github.com/IBM/fp-go/v2/eq"`

Eq typeclass: equality comparison.

- `Eq[A]` -- interface with `Equals(A, A) bool`

### Exported API

```go
func Contramap[A, B any](f func(b B) A) func(Eq[A]) Eq[B]
func Equals[T any](eq Eq[T]) func(T) func(T) bool
func Monoid[A any]() M.Monoid[Eq[A]]
func Semigroup[A any]() S.Semigroup[Eq[A]]
type Eq[T any] interface {
func Empty[T any]() Eq[T]
func FromEquals[T any](c func(x, y T) bool) Eq[T]
func FromStrictEquals[T comparable]() Eq[T]
```

## package `github.com/IBM/fp-go/v2/ord`

Import: `import "github.com/IBM/fp-go/v2/ord"`

Ord typeclass: total ordering.

- `Ord[A]` -- interface with `Compare(A, A) int`

### Exported API

```go
func Between[A any](o Ord[A]) func(A, A) func(A) bool
func Clamp[A any](o Ord[A]) func(A, A) func(A) A
func Geq[A any](o Ord[A]) func(A) func(A) bool
func Gt[A any](o Ord[A]) func(A) func(A) bool
func Leq[A any](o Ord[A]) func(A) func(A) bool
func Lt[A any](o Ord[A]) func(A) func(A) bool
func Max[A any](o Ord[A]) func(A, A) A
func MaxSemigroup[A any](o Ord[A]) S.Semigroup[A]
func Min[A any](o Ord[A]) func(A, A) A
func MinSemigroup[A any](o Ord[A]) S.Semigroup[A]
func Monoid[A any]() M.Monoid[Ord[A]]
func Semigroup[A any]() S.Semigroup[Ord[A]]
func ToEq[T any](o Ord[T]) E.Eq[T]
type Kleisli[A, B any] = func(A) Ord[B]
type Operator[A, B any] = Kleisli[Ord[A], B]
func Contramap[A, B any](f func(B) A) Operator[A, B]
type Ord[T any] interface {
func FromCompare[T any](compare func(T, T) int) Ord[T]
func FromStrictCompare[A C.Ordered]() Ord[A]
func MakeOrd[T any](c func(x, y T) int, e func(x, y T) bool) Ord[T]
func OrdTime() Ord[time.Time]
func Reverse[T any](o Ord[T]) Ord[T]
```

## package `github.com/IBM/fp-go/v2/semigroup`

Import: `import "github.com/IBM/fp-go/v2/semigroup"`

Semigroup typeclass: associative binary operation.

- `Semigroup[A]` -- interface with `Concat(A, A) A`

### Exported API

```go
func AppendTo[A any](s Semigroup[A]) func(A) func(A) A
func ConcatAll[A any](s Semigroup[A]) func(A) func([]A) A
func ConcatWith[A any](s Semigroup[A]) func(A) func(A) A
func GenericConcatAll[GA ~[]A, A any](s Semigroup[A]) func(A) func(GA) A
func GenericMonadConcatAll[GA ~[]A, A any](s Semigroup[A]) func(GA, A) A
func MonadConcatAll[A any](s Semigroup[A]) func([]A, A) A
func ToMagma[A any](s Semigroup[A]) M.Magma[A]
type Semigroup[A any] interface {
func AltSemigroup[HKTA any, LAZYHKTA ~func() HKTA](
func ApplySemigroup[A, HKTA, HKTFA any](
func First[A any]() Semigroup[A]
func FunctionSemigroup[A, B any](s Semigroup[B]) Semigroup[func(A) B]
func Last[A any]() Semigroup[A]
func MakeSemigroup[A any](c func(A, A) A) Semigroup[A]
func Reverse[A any](m Semigroup[A]) Semigroup[A]
```

## package `github.com/IBM/fp-go/v2/monoid`

Import: `import "github.com/IBM/fp-go/v2/monoid"`

Monoid typeclass: semigroup with identity.

- `Monoid[A]` -- interface with `Concat(A, A) A` and `Empty() A`

### Exported API

```go
func ConcatAll[A any](m Monoid[A]) func([]A) A
func Fold[A any](m Monoid[A]) func([]A) A
func GenericConcatAll[GA ~[]A, A any](m Monoid[A]) func(GA) A
func ToSemigroup[A any](m Monoid[A]) S.Semigroup[A]
type Monoid[A any] interface {
func AltMonoid[HKTA any, LAZYHKTA ~func() HKTA](
func AlternativeMonoid[A, HKTA, HKTFA any, LAZYHKTA ~func() HKTA](
func ApplicativeMonoid[A, HKTA, HKTFA any](
func FunctionMonoid[A, B any](m Monoid[B]) Monoid[func(A) B]
func MakeMonoid[A any](c func(A, A) A, e A) Monoid[A]
func Reverse[A any](m Monoid[A]) Monoid[A]
func VoidMonoid() Monoid[Void]
type Void = function.Void
```

---

# Primitives

## package `github.com/IBM/fp-go/v2/number`

Import: `import N "github.com/IBM/fp-go/v2/number"`

Numeric operations: Add, Sub, Mul, Div, comparisons, eq/ord instances.

### Exported API

```go
func Add[T Number](right T) func(T) T
func Div[T Number](right T) func(T) T
func Inc[T Number](value T) T
func LessThan[A C.Ordered](a A) func(A) bool
func MagmaDiv[A Number]() M.Magma[A]
func MagmaSub[A Number]() M.Magma[A]
func Max[A C.Ordered](a, b A) A
func Min[A C.Ordered](a, b A) A
func MonoidProduct[A Number]() M.Monoid[A]
func MonoidSum[A Number]() M.Monoid[A]
func MoreThan[A C.Ordered](a A) func(A) bool
func Mul[T Number](right T) func(T) T
func SemigroupProduct[A Number]() S.Semigroup[A]
func SemigroupSum[A Number]() S.Semigroup[A]
func Sub[T Number](right T) func(T) T
type Number interface {
```

## package `github.com/IBM/fp-go/v2/string`

Import: `import S "github.com/IBM/fp-go/v2/string"`

String operations and instances.

### Exported API

```go
var (
var Monoid = M.MakeMonoid(concat, "")
var Semigroup = S.MakeSemigroup(concat)
func Append(suffix string) func(string) string
func Eq(left, right string) bool
func Format[T any](format string) func(T) string
func Intersperse(middle string) func(string, string) string
func IntersperseMonoid(middle string) M.Monoid[string]
func IntersperseSemigroup(middle string) S.Semigroup[string]
func IsEmpty(s string) bool
func IsNonEmpty(s string) bool
func Prepend(prefix string) func(string) string
func Size(s string) int
func ToBytes(s string) []byte
func ToRunes(s string) []rune
```

## package `github.com/IBM/fp-go/v2/boolean`

Import: `import B "github.com/IBM/fp-go/v2/boolean"`

Boolean operations, monoids (All, Any), and Eq/Ord instances.

### Exported API

```go
var (
type Monoid = monoid.Monoid[bool]
```

## package `github.com/IBM/fp-go/v2/bytes`

Import: `import "github.com/IBM/fp-go/v2/bytes"`

Byte slice operations.

### Exported API

```go
var (
func Empty() []byte
func Size(as []byte) int
func ToString(a []byte) string
```

---

# Other

## package `github.com/IBM/fp-go/v2/identity`

Import: `import "github.com/IBM/fp-go/v2/identity"`

Identity monad: wraps a value with no additional effect.

### Exported API

```go
func ApS[S1, S2, T any](
func Bind[S1, S2, T any](
func BindTo[S1, T any](
func Do[S any](
func Extract[A any](a A) A
func Let[S1, S2, T any](
func LetTo[S1, S2, B any](
func MapTo[A, B any](b B) func(A) B
func Monad[A, B any]() monad.Monad[A, B, A, B, func(A) B]
func MonadAp[B, A any](fab func(A) B, fa A) B
func MonadChain[A, B any](ma A, f Kleisli[A, B]) B
func MonadChainFirst[A, B any](fa A, f Kleisli[A, B]) A
func MonadFlap[B, A any](fab func(A) B, a A) B
func MonadMap[A, B any](fa A, f func(A) B) B
func MonadMapTo[A, B any](_ A, b B) B
func Of[A any](a A) A
func SequenceT1[T1 any](t1 T1) T.Tuple1[T1]
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10) T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]
func SequenceT2[T1, T2 any](t1 T1, t2 T2) T.Tuple2[T1, T2]
func SequenceT3[T1, T2, T3 any](t1 T1, t2 T2, t3 T3) T.Tuple3[T1, T2, T3]
func SequenceT4[T1, T2, T3, T4 any](t1 T1, t2 T2, t3 T3, t4 T4) T.Tuple4[T1, T2, T3, T4]
func SequenceT5[T1, T2, T3, T4, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) T.Tuple5[T1, T2, T3, T4, T5]
func SequenceT6[T1, T2, T3, T4, T5, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) T.Tuple6[T1, T2, T3, T4, T5, T6]
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) T.Tuple7[T1, T2, T3, T4, T5, T6, T7]
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
func SequenceTuple1[T1 any](t T.Tuple1[T1]) T.Tuple1[T1]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]
func SequenceTuple2[T1, T2 any](t T.Tuple2[T1, T2]) T.Tuple2[T1, T2]
func SequenceTuple3[T1, T2, T3 any](t T.Tuple3[T1, T2, T3]) T.Tuple3[T1, T2, T3]
func SequenceTuple4[T1, T2, T3, T4 any](t T.Tuple4[T1, T2, T3, T4]) T.Tuple4[T1, T2, T3, T4]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t T.Tuple5[T1, T2, T3, T4, T5]) T.Tuple5[T1, T2, T3, T4, T5]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t T.Tuple6[T1, T2, T3, T4, T5, T6]) T.Tuple6[T1, T2, T3, T4, T5, T6]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t T.Tuple7[T1, T2, T3, T4, T5, T6, T7]) T.Tuple7[T1, T2, T3, T4, T5, T6, T7]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
func TraverseTuple1[F1 ~func(A1) T1, A1, T1 any](f1 F1) func(T.Tuple1[A1]) T.Tuple1[T1]
func TraverseTuple10[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, F4 ~func(A4) T4, F5 ~func(A5) T5, F6 ~func(A6) T6, F7 ~func(A7) T7, F8 ~func(A8) T8, F9 ~func(A9) T9, F10 ~func(A10) T10, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(T.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) T.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]
func TraverseTuple2[F1 ~func(A1) T1, F2 ~func(A2) T2, A1, T1, A2, T2 any](f1 F1, f2 F2) func(T.Tuple2[A1, A2]) T.Tuple2[T1, T2]
func TraverseTuple3[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, A1, T1, A2, T2, A3, T3 any](f1 F1, f2 F2, f3 F3) func(T.Tuple3[A1, A2, A3]) T.Tuple3[T1, T2, T3]
func TraverseTuple4[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, F4 ~func(A4) T4, A1, T1, A2, T2, A3, T3, A4, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T.Tuple4[A1, A2, A3, A4]) T.Tuple4[T1, T2, T3, T4]
func TraverseTuple5[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, F4 ~func(A4) T4, F5 ~func(A5) T5, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T.Tuple5[A1, A2, A3, A4, A5]) T.Tuple5[T1, T2, T3, T4, T5]
func TraverseTuple6[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, F4 ~func(A4) T4, F5 ~func(A5) T5, F6 ~func(A6) T6, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(T.Tuple6[A1, A2, A3, A4, A5, A6]) T.Tuple6[T1, T2, T3, T4, T5, T6]
func TraverseTuple7[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, F4 ~func(A4) T4, F5 ~func(A5) T5, F6 ~func(A6) T6, F7 ~func(A7) T7, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(T.Tuple7[A1, A2, A3, A4, A5, A6, A7]) T.Tuple7[T1, T2, T3, T4, T5, T6, T7]
func TraverseTuple8[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, F4 ~func(A4) T4, F5 ~func(A5) T5, F6 ~func(A6) T6, F7 ~func(A7) T7, F8 ~func(A8) T8, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(T.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) T.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
func TraverseTuple9[F1 ~func(A1) T1, F2 ~func(A2) T2, F3 ~func(A3) T3, F4 ~func(A4) T4, F5 ~func(A5) T5, F6 ~func(A6) T6, F7 ~func(A7) T7, F8 ~func(A8) T8, F9 ~func(A9) T9, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(T.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) T.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
type Kleisli[A, B any] = func(A) B
type Operator[A, B any] = Kleisli[A, B]
func Ap[B, A any](fa A) Operator[func(A) B, B]
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func Extend[A, B any](f func(A) B) Operator[A, B]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Map[A, B any](f func(A) B) Operator[A, B]
```

## package `github.com/IBM/fp-go/v2/lazy`

Import: `import "github.com/IBM/fp-go/v2/lazy"`

Lazy evaluation: `Lazy[A] = func() A`. Memoization via `Memoize`.

### Exported API

```go
type A:
func Ap[B, A any](ma Lazy[A]) func(Lazy[func(A) B]) Lazy[B]
func ApplicativeMonoid[A any](m M.Monoid[A]) M.Monoid[Lazy[A]]
func ApplySemigroup[A any](s S.Semigroup[A]) S.Semigroup[Lazy[A]]
func Eq[A any](e EQ.Eq[A]) EQ.Eq[Lazy[A]]
func Map[A, B any](f func(A) B) func(fa Lazy[A]) Lazy[B]
type Kleisli[A, B any] = func(A) Lazy[B]
func ApFirst[A, B any](second Lazy[B]) Kleisli[Lazy[A], A]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func ApSecond[A, B any](second Lazy[B]) Kleisli[Lazy[A], B]
func Bind[S1, S2, T any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Kleisli[Lazy[A], B]
func ChainFirst[A, B any](f Kleisli[A, B]) Kleisli[Lazy[A], A]
func ChainTo[A, B any](fb Lazy[B]) Kleisli[Lazy[A], B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func MapTo[A, B any](b B) Kleisli[Lazy[A], B]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) Lazy[B]) Kleisli[[]A, []B]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Lazy[B]) Kleisli[map[K]A, map[K]B]
type Lazy[A any] = func() A
var Now Lazy[time.Time] = io.Now
func Defer[A any](gen func() Lazy[A]) Lazy[A]
func Do[S any](
func Flatten[A any](mma Lazy[Lazy[A]]) Lazy[A]
func FromImpure(f func()) Lazy[Void]
func FromLazy[A any](a Lazy[A]) Lazy[A]
func Memoize[A any](ma Lazy[A]) Lazy[A]
func MonadAp[B, A any](mab Lazy[func(A) B], ma Lazy[A]) Lazy[B]
func MonadApFirst[A, B any](first Lazy[A], second Lazy[B]) Lazy[A]
func MonadApSecond[A, B any](first Lazy[A], second Lazy[B]) Lazy[B]
func MonadChain[A, B any](fa Lazy[A], f Kleisli[A, B]) Lazy[B]
func MonadChainFirst[A, B any](fa Lazy[A], f Kleisli[A, B]) Lazy[A]
func MonadChainTo[A, B any](fa Lazy[A], fb Lazy[B]) Lazy[B]
func MonadMap[A, B any](fa Lazy[A], f func(A) B) Lazy[B]
func MonadMapTo[A, B any](fa Lazy[A], b B) Lazy[B]
func MonadOf[A any](a A) Lazy[A]
func MonadTraverseArray[A, B any](tas []A, f Kleisli[A, B]) Lazy[[]B]
func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f Kleisli[A, B]) Lazy[map[K]B]
func Of[A any](a A) Lazy[A]
func Retrying[A any](
func SequenceArray[A any](tas []Lazy[A]) Lazy[[]A]
func SequenceRecord[K comparable, A any](tas map[K]Lazy[A]) Lazy[map[K]A]
func SequenceT1[A any](a Lazy[A]) Lazy[tuple.Tuple1[A]]
func SequenceT2[A, B any](a Lazy[A], b Lazy[B]) Lazy[tuple.Tuple2[A, B]]
func SequenceT3[A, B, C any](a Lazy[A], b Lazy[B], c Lazy[C]) Lazy[tuple.Tuple3[A, B, C]]
func SequenceT4[A, B, C, D any](a Lazy[A], b Lazy[B], c Lazy[C], d Lazy[D]) Lazy[tuple.Tuple4[A, B, C, D]]
type Operator[A, B any] = Kleisli[Lazy[A], B]
type Predicate[A any] = predicate.Predicate[A]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/constant`

Import: `import "github.com/IBM/fp-go/v2/constant"`

Constant functor.

### Exported API

```go
func Ap[E, A, B any](s S.Semigroup[E]) func(fa Const[E, A]) func(fab Const[E, func(A) B]) Const[E, B]
func Map[E, A, B any](f func(A) B) func(fa Const[E, A]) Const[E, B]
func MonadAp[E, A, B any](s S.Semigroup[E]) func(fab Const[E, func(A) B], fa Const[E, A]) Const[E, B]
func Monoid[A any](a A) M.Monoid[A]
func Of[E, A any](m M.Monoid[E]) func(A) Const[E, A]
func Unwrap[E, A any](c Const[E, A]) E
type Const[E, A any] struct {
func Make[E, A any](e E) Const[E, A]
func MonadMap[E, A, B any](fa Const[E, A], _ func(A) B) Const[E, B]
```

## package `github.com/IBM/fp-go/v2/json`

Import: `import "github.com/IBM/fp-go/v2/json"`

JSON marshal/unmarshal utilities returning Either/Result.

### Exported API

```go
type Either[A any] = E.Either[error, A]
func Marshal[A any](a A) Either[[]byte]
func MarshalIndent[A any](a A) Either[[]byte]
func ToTypeE[A any](src any) Either[A]
func Unmarshal[A any](data []byte) Either[A]
type Option[A any] = option.Option[A]
func ToTypeO[A any](src any) Option[A]
```

## package `github.com/IBM/fp-go/v2/di`

Import: `import "github.com/IBM/fp-go/v2/di"`

Dependency injection utilities.

### Exported API

```go
var (
var RunMain = F.Flow3(
func ConstProvider[R any](token InjectionToken[R], value R) DIE.Provider
func MakeProvider0[R any](
func MakeProvider1[T1 any, R any](
func MakeProvider10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any, R any](
func MakeProvider11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any, R any](
func MakeProvider12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any, R any](
func MakeProvider13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any, R any](
func MakeProvider14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any, R any](
func MakeProvider15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any, R any](
func MakeProvider2[T1, T2 any, R any](
func MakeProvider3[T1, T2, T3 any, R any](
func MakeProvider4[T1, T2, T3, T4 any, R any](
func MakeProvider5[T1, T2, T3, T4, T5 any, R any](
func MakeProvider6[T1, T2, T3, T4, T5, T6 any, R any](
func MakeProvider7[T1, T2, T3, T4, T5, T6, T7 any, R any](
func MakeProvider8[T1, T2, T3, T4, T5, T6, T7, T8 any, R any](
func MakeProvider9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any, R any](
func MakeProviderFactory0[R any](
func MakeProviderFactory1[T1 any, R any](
func MakeProviderFactory10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any, R any](
func MakeProviderFactory11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any, R any](
func MakeProviderFactory12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any, R any](
func MakeProviderFactory13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any, R any](
func MakeProviderFactory14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any, R any](
func MakeProviderFactory15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any, R any](
func MakeProviderFactory2[T1, T2 any, R any](
func MakeProviderFactory3[T1, T2, T3 any, R any](
func MakeProviderFactory4[T1, T2, T3, T4 any, R any](
func MakeProviderFactory5[T1, T2, T3, T4, T5 any, R any](
func MakeProviderFactory6[T1, T2, T3, T4, T5, T6 any, R any](
func MakeProviderFactory7[T1, T2, T3, T4, T5, T6, T7 any, R any](
func MakeProviderFactory8[T1, T2, T3, T4, T5, T6, T7, T8 any, R any](
func MakeProviderFactory9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any, R any](
func Resolve[T any](token InjectionToken[T]) RIOR.ReaderIOResult[DIE.InjectableFactory, T]
type Dependency[T any] interface {
type Entry[K comparable, V any] = record.Entry[K, V]
type IOOption[T any] = iooption.IOOption[T]
type IOResult[T any] = ioresult.IOResult[T]
type InjectionToken[T any] interface {
func MakeToken[T any](name string) InjectionToken[T]
func MakeTokenWithDefault[T any](name string, providerFactory DIE.ProviderFactory) InjectionToken[T]
func MakeTokenWithDefault0[R any](name string, fct IOResult[R]) InjectionToken[R]
func MakeTokenWithDefault1[T1 any, R any](
func MakeTokenWithDefault10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any, R any](
func MakeTokenWithDefault11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any, R any](
func MakeTokenWithDefault12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any, R any](
func MakeTokenWithDefault13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any, R any](
func MakeTokenWithDefault14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any, R any](
func MakeTokenWithDefault15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any, R any](
func MakeTokenWithDefault2[T1, T2 any, R any](
func MakeTokenWithDefault3[T1, T2, T3 any, R any](
func MakeTokenWithDefault4[T1, T2, T3, T4 any, R any](
func MakeTokenWithDefault5[T1, T2, T3, T4, T5 any, R any](
func MakeTokenWithDefault6[T1, T2, T3, T4, T5, T6 any, R any](
func MakeTokenWithDefault7[T1, T2, T3, T4, T5, T6, T7 any, R any](
func MakeTokenWithDefault8[T1, T2, T3, T4, T5, T6, T7, T8 any, R any](
func MakeTokenWithDefault9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any, R any](
type MultiInjectionToken[T any] interface {
func MakeMultiToken[T any](name string) MultiInjectionToken[T]
type Option[T any] = option.Option[T]
type Result[T any] = result.Result[T]
```

## package `github.com/IBM/fp-go/v2/builder`

Import: `import "github.com/IBM/fp-go/v2/builder"`

Builder pattern utilities.

### Exported API

```go
type Builder[T any] interface {
type Option[T any] = option.Option[T]
type Prism[S, A any] = prism.Prism[S, A]
func BuilderPrism[T any, B Builder[T]](creator func(T) B) Prism[B, T]
type Result[T any] = result.Result[T]
```

## package `github.com/IBM/fp-go/v2/retry`

Import: `import "github.com/IBM/fp-go/v2/retry"`

Retry policies and combinators.

### Exported API

```go
var CumulativeDelayLens = L.MakeLensWithName(
var DefaultRetryStatus = RetryStatus{
var IterNumberLens = L.MakeLensWithName(
var Monoid = M.FunctionMonoid[RetryStatus](O.ApplicativeMonoid(M.MakeMonoid(
var PreviousDelayLens = L.MakeLensWithName(
func Always[A any](a A) func(RetryStatus) A
func IterNumber(rs RetryStatus) uint
type Option[A any] = option.Option[A]
type RetryPolicy = func(RetryStatus) Option[time.Duration]
func CapDelay(maxDelay time.Duration, policy RetryPolicy) RetryPolicy
func ConstantDelay(delay time.Duration) RetryPolicy
func ExponentialBackoff(delay time.Duration) RetryPolicy
func LimitRetries(i uint) RetryPolicy
type RetryStatus struct {
func ApplyPolicy(policy RetryPolicy, status RetryStatus) RetryStatus
```

## package `github.com/IBM/fp-go/v2/circuitbreaker`

Import: `import "github.com/IBM/fp-go/v2/circuitbreaker"`

Circuit breaker pattern.

### Exported API

```go
var (
var AnyError = option.FromPredicate(E.IsNonNil)
var InfrastructureError = option.FromPredicate(shouldOpenCircuit)
var MakeCircuitBreakerError = MakeCircuitBreakerErrorWithName("Generic Circuit Breaker")
func MakeCircuitBreakerErrorWithName(name string) func(time.Time) error
func MakeSingletonBreaker[HKTT any](
type BreakerState = Either[openState, ClosedState]
type CircuitBreakerError struct {
func (e *CircuitBreakerError) Error() string
type ClosedState interface {
func MakeClosedStateCounter(maxFailures uint) ClosedState
func MakeClosedStateHistory(
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[T any] = io.IO[T]
type IORef[T any] = ioref.IORef[T]
type Metrics interface {
func MakeMetricsFromLogger(name string, logger *log.Logger) Metrics
func MakeVoidMetrics() Metrics
type Option[A any] = option.Option[A]
type Ord[A any] = ord.Ord[A]
type Pair[L, R any] = pair.Pair[L, R]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderIO[R, A any] = readerio.ReaderIO[R, A]
type State[T, R any] = state.State[T, R]
func MakeCircuitBreaker[E, T, HKTT, HKTOP, HKTHKTT any](
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/tailrec`

Import: `import "github.com/IBM/fp-go/v2/tailrec"`

Tail recursion via trampolining.

- `Trampoline[A, B]` -- either Continue(A) or Done(B)

### Exported API

```go
type Trampoline[B, L any] struct {
func Bounce[L, B any](b B) Trampoline[B, L]
func Land[B, L any](l L) Trampoline[B, L]
func (t Trampoline[B, L]) Format(f fmt.State, verb rune)
func (t Trampoline[B, L]) GoString() string
func (t Trampoline[B, L]) LogValue() slog.Value
func (t Trampoline[B, L]) String() string
```

## package `github.com/IBM/fp-go/v2/ioref`

Import: `import "github.com/IBM/fp-go/v2/ioref"`

Mutable reference in IO context.

### Exported API

```go
func Modify[A any](f Endomorphism[A]) io.Kleisli[IORef[A], A]
func ModifyIOK[A any](f io.Kleisli[A, A]) io.Kleisli[IORef[A], A]
func ModifyIOKWithResult[A, B any](f io.Kleisli[A, Pair[A, B]]) io.Kleisli[IORef[A], B]
func ModifyReaderIOK[R, A any](f readerio.Kleisli[R, A, A]) readerio.Kleisli[R, IORef[A], A]
func ModifyReaderIOKWithResult[R, A, B any](f readerio.Kleisli[R, A, Pair[A, B]]) readerio.Kleisli[R, IORef[A], B]
func ModifyWithResult[A, B any](f func(A) Pair[A, B]) io.Kleisli[IORef[A], B]
func Write[A any](a A) io.Kleisli[IORef[A], A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[A any] = io.IO[A]
func MakeIORef[A any](a A) IO[IORef[A]]
func Read[A any](ref IORef[A]) IO[A]
type IORef[A any] = *ioRef[A]
type Pair[A, B any] = pair.Pair[A, B]
type ReaderIO[R, A any] = readerio.ReaderIO[R, A]
```

## package `github.com/IBM/fp-go/v2/consumer`

Import: `import "github.com/IBM/fp-go/v2/consumer"`

Consumer[A] = func(A). Void-returning function.

### Exported API

```go
type Consumer[A any] = func(A)
type Operator[A, B any] = func(Consumer[A]) Consumer[B]
func Compose[R1, R2 any](f func(R2) R1) Operator[R1, R2]
func Contramap[R1, R2 any](f func(R2) R1) Operator[R1, R2]
func Local[R1, R2 any](f func(R2) R1) Operator[R1, R2]
```

## package `github.com/IBM/fp-go/v2/erasure`

Import: `import "github.com/IBM/fp-go/v2/erasure"`

Type erasure utilities.

### Exported API

```go
func Erase[T any](t T) any
func Erase0[T1 any](f func() T1) func() any
func Erase1[T1, T2 any](f func(T1) T2) func(any) any
func Erase2[T1, T2, T3 any](f func(T1, T2) T3) func(any, any) any
func SafeUnerase[T any](t any) E.Either[error, T]
func Unerase[T any](t any) T
```

## package `github.com/IBM/fp-go/v2/iterator/stateless`

Import: `import "github.com/IBM/fp-go/v2/iterator/stateless"`

Stateless iterator operations using iter.Seq.

### Exported API

```go
func Current[U any](m Pair[Iterator[U], U]) U
func Fold[U any](m M.Monoid[U]) func(Iterator[U]) U
func FoldMap[U, V any](m M.Monoid[V]) func(func(U) V) func(ma Iterator[U]) V
func Monad[A, B any]() monad.Monad[A, B, Iterator[A], Iterator[B], Iterator[func(A) B]]
func Monoid[U any]() M.Monoid[Iterator[U]]
func Reduce[U, V any](f func(V, U) V, initial V) func(Iterator[U]) V
func ToArray[U any](u Iterator[U]) []U
type IO[A any] = io.IO[A]
type Iterator[U any] Lazy[Option[Pair[Iterator[U], U]]]
func Count(start int) Iterator[int]
func Cycle[U any](ma Iterator[U]) Iterator[U]
func Do[S any](
func Empty[U any]() Iterator[U]
func Flatten[U any](ma Iterator[Iterator[U]]) Iterator[U]
func From[U any](data ...U) Iterator[U]
func FromArray[U any](as []U) Iterator[U]
func FromIO[U any](io IO[U]) Iterator[U]
func FromLazy[U any](l Lazy[U]) Iterator[U]
func FromReflect(val R.Value) Iterator[R.Value]
func MakeBy[FCT ~func(int) U, U any](f FCT) Iterator[U]
func MonadAp[V, U any](fab Iterator[func(U) V], ma Iterator[U]) Iterator[V]
func MonadChain[U, V any](ma Iterator[U], f Kleisli[U, V]) Iterator[V]
func MonadChainFirst[U, V any](ma Iterator[U], f Kleisli[U, V]) Iterator[U]
func MonadMap[U, V any](ma Iterator[U], f func(U) V) Iterator[V]
func Next[U any](m Pair[Iterator[U], U]) Iterator[U]
func Of[U any](a U) Iterator[U]
func Repeat[U any](n int, a U) Iterator[U]
func Replicate[U any](a U) Iterator[U]
func StrictUniq[A comparable](as Iterator[A]) Iterator[A]
func ZipWith[FCT ~func(A, B) C, A, B, C any](fa Iterator[A], fb Iterator[B], f FCT) Iterator[C]
type Kleisli[A, B any] = reader.Reader[A, Iterator[B]]
func Chain[U, V any](f Kleisli[U, V]) Kleisli[Iterator[U], V]
type Lazy[A any] = lazy.Lazy[A]
type Operator[A, B any] = Kleisli[Iterator[A], B]
func Ap[V, U any](ma Iterator[U]) Operator[func(U) V, V]
func ApS[S1, S2, T any](
func Bind[S1, S2, T any](
func BindTo[S1, T any](
func ChainFirst[U, V any](f Kleisli[U, V]) Operator[U, U]
func Compress[U any](sel Iterator[bool]) Operator[U, U]
func DropWhile[U any](pred Predicate[U]) Operator[U, U]
func Filter[U any](f Predicate[U]) Operator[U, U]
func FilterChain[U, V any](f func(U) Option[Iterator[V]]) Operator[U, V]
func FilterMap[U, V any](f func(U) Option[V]) Operator[U, V]
func Let[S1, S2, T any](
func LetTo[S1, S2, T any](
func Map[U, V any](f func(U) V) Operator[U, V]
func Scan[FCT ~func(V, U) V, U, V any](f FCT, initial V) Operator[U, V]
func Take[U any](n int) Operator[U, U]
func Uniq[A any, K comparable](f func(A) K) Operator[A, A]
func Zip[A, B any](fb Iterator[B]) Operator[A, Pair[A, B]]
type Option[A any] = option.Option[A]
func First[U any](mu Iterator[U]) Option[U]
func Last[U any](mu Iterator[U]) Option[U]
type Pair[L, R any] = pair.Pair[L, R]
type Predicate[A any] = predicate.Predicate[A]
func Any[U any](pred Predicate[U]) Predicate[Iterator[U]]
type Seq[T any] = iter.Seq[T]
func ToSeq[T any](it Iterator[T]) Seq[T]
type Seq2[K, V any] = iter.Seq2[K, V]
func ToSeq2[K, V any](it Iterator[Pair[K, V]]) Seq2[K, V]
```

---

# Idiomatic

The idiomatic packages provide Go-native APIs using (value, bool) and (value, error) tuples
instead of Option/Either/Result wrapper types.

## package `github.com/IBM/fp-go/v2/idiomatic/option`

Import: `import "github.com/IBM/fp-go/v2/idiomatic/option"`

Idiomatic Option using Go tuples: `(A, bool)` instead of `Option[A]`.

### Exported API

```go
func ApS[S1, S2, T any](
func ApSL[S, T any](
func Do[S any](
func Eq[A any](eq EQ.Eq[A]) func(A, bool) func(A, bool) bool
func Flow1[F1 ~func(T0, bool) (T1, bool), T0, T1 any](f1 F1) func(T0, bool) (T1, bool)
func Flow2[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), T0, T1, T2 any](f1 F1, f2 F2) func(T0, bool) (T2, bool)
func Flow3[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), T0, T1, T2, T3 any](f1 F1, f2 F2, f3 F3) func(T0, bool) (T3, bool)
func Flow4[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), T0, T1, T2, T3, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T0, bool) (T4, bool)
func Flow5[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), F5 ~func(T4, bool) (T5, bool), T0, T1, T2, T3, T4, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T0, bool) (T5, bool)
func Fold[A, B any](onNone func() B, onSome func(A) B) func(A, bool) B
func FromEq[A any](pred eq.Eq[A]) func(A) Kleisli[A, A]
func FromNillable[A any](a *A) (*A, bool)
func FromStrictCompare[A C.Ordered]() func(A, bool) func(A, bool) int
func FromStrictEquals[A comparable]() func(A, bool) func(A, bool) bool
func GetOrElse[A any](onNone func() A) func(A, bool) A
func IsNone[T any](t T, tok bool) bool
func IsSome[T any](t T, tok bool) bool
func Logger[A any](loggers ...*log.Logger) func(string) Operator[A, A]
func None[T any]() (t T, tok bool)
func Of[T any](value T) (T, bool)
func Ord[A any](o ord.Ord[A]) func(A, bool) func(A, bool) int
func Pipe1[F1 ~func(T0) (T1, bool), T0, T1 any](t0 T0, f1 F1) (T1, bool)
func Pipe2[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), T0, T1, T2 any](t0 T0, f1 F1, f2 F2) (T2, bool)
func Pipe3[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), T0, T1, T2, T3 any](t0 T0, f1 F1, f2 F2, f3 F3) (T3, bool)
func Pipe4[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), T0, T1, T2, T3, T4 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4) (T4, bool)
func Pipe5[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), F5 ~func(T4, bool) (T5, bool), T0, T1, T2, T3, T4, T5 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) (T5, bool)
func Reduce[A, B any](f func(B, A) B, initial B) func(A, bool) B
func Some[T any](value T) (T, bool)
func ToAny[T any](src T) (any, bool)
func ToString[T any](t T, tok bool) string
func ToType[T any](src any) (T, bool)
func TraverseTuple1[F1 ~Kleisli[A1, T1], A1, T1 any](f1 F1) func(A1) (T1, bool)
func TraverseTuple10[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], F10 ~Kleisli[A10, T10], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(A1, A2, A3, A4, A5, A6, A7, A8, A9, A10) (T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, bool)
func TraverseTuple2[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], A1, T1, A2, T2 any](f1 F1, f2 F2) func(A1, A2) (T1, T2, bool)
func TraverseTuple3[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], A1, T1, A2, T2, A3, T3 any](f1 F1, f2 F2, f3 F3) func(A1, A2, A3) (T1, T2, T3, bool)
func TraverseTuple4[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], A1, T1, A2, T2, A3, T3, A4, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(A1, A2, A3, A4) (T1, T2, T3, T4, bool)
func TraverseTuple5[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(A1, A2, A3, A4, A5) (T1, T2, T3, T4, T5, bool)
func TraverseTuple6[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(A1, A2, A3, A4, A5, A6) (T1, T2, T3, T4, T5, T6, bool)
func TraverseTuple7[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(A1, A2, A3, A4, A5, A6, A7) (T1, T2, T3, T4, T5, T6, T7, bool)
func TraverseTuple8[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(A1, A2, A3, A4, A5, A6, A7, A8) (T1, T2, T3, T4, T5, T6, T7, T8, bool)
func TraverseTuple9[F1 ~Kleisli[A1, T1], F2 ~Kleisli[A2, T2], F3 ~Kleisli[A3, T3], F4 ~Kleisli[A4, T4], F5 ~Kleisli[A5, T5], F6 ~Kleisli[A6, T6], F7 ~Kleisli[A7, T7], F8 ~Kleisli[A8, T8], F9 ~Kleisli[A9, T9], A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(A1, A2, A3, A4, A5, A6, A7, A8, A9) (T1, T2, T3, T4, T5, T6, T7, T8, T9, bool)
type Endomorphism[T any] = endomorphism.Endomorphism[T]
type Functor[A, B any] interface {
func MakeFunctor[A, B any]() Functor[A, B]
type Kleisli[A, B any] = func(A) (B, bool)
func FromNonZero[A comparable]() Kleisli[A, A]
func FromPredicate[A any](pred func(A) bool) Kleisli[A, A]
func FromZero[A comparable]() Kleisli[A, A]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseArrayWithIndex[A, B any](f func(int, A) (B, bool)) Kleisli[[]A, []B]
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) (B, bool)) Kleisli[GA, GB]
func TraverseIter[A, B any](f Kleisli[A, B]) Kleisli[Seq[A], Seq[B]]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) (B, bool)) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) (B, bool)) Kleisli[GA, GB]
type Operator[A, B any] = func(A, bool) (B, bool)
func Alt[A any](that func() (A, bool)) Operator[A, A]
func Ap[B, A any](fa A, faok bool) Operator[func(A) B, B]
func Bind[S1, S2, A any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainTo[A, B any](b B, bok bool) Operator[A, B]
func Filter[A any](pred func(A) bool) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, B any](
func LetL[S, T any](
func LetTo[S1, S2, B any](
func LetToL[S, T any](
func Map[A, B any](f func(a A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
type Pointed[A any] interface {
func MakePointed[A any]() Pointed[A]
type Seq[T any] = iter.Seq[T]
```

## package `github.com/IBM/fp-go/v2/idiomatic/result`

Import: `import "github.com/IBM/fp-go/v2/idiomatic/result"`

Idiomatic Result using Go tuples: `(A, error)` instead of `Result[A]`.

### Exported API

```go
func ApS[S1, S2, T any](
func ApSL[S, T any](
func ApV[B, A any](sg S.Semigroup[error]) func(A, error) Operator[func(A) B, B]
func ChainOptionK[A, B any](onNone func() error) func(option.Kleisli[A, B]) Operator[A, B]
func Curry0[R any](f func() (R, error)) func() (R, error)
func Curry1[T1, R any](f func(T1) (R, error)) func(T1) (R, error)
func Curry2[T1, T2, R any](f func(T1, T2) (R, error)) func(T1) func(T2) (R, error)
func Curry3[T1, T2, T3, R any](f func(T1, T2, T3) (R, error)) func(T1) func(T2) func(T3) (R, error)
func Curry4[T1, T2, T3, T4, R any](f func(T1, T2, T3, T4) (R, error)) func(T1) func(T2) func(T3) func(T4) (R, error)
func Do[S any](
func Eq[A any](eq EQ.Eq[A]) func(A, error) func(A, error) bool
func Flow1[F1 ~func(T0, error) (T1, error), T0, T1 any](f1 F1) func(T0, error) (T1, error)
func Flow2[F1 ~func(T0, error) (T1, error), F2 ~func(T1, error) (T2, error), T0, T1, T2 any](f1 F1, f2 F2) func(T0, error) (T2, error)
func Flow3[F1 ~func(T0, error) (T1, error), F2 ~func(T1, error) (T2, error), F3 ~func(T2, error) (T3, error), T0, T1, T2, T3 any](f1 F1, f2 F2, f3 F3) func(T0, error) (T3, error)
func Flow4[F1 ~func(T0, error) (T1, error), F2 ~func(T1, error) (T2, error), F3 ~func(T2, error) (T3, error), F4 ~func(T3, error) (T4, error), T0, T1, T2, T3, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T0, error) (T4, error)
func Flow5[F1 ~func(T0, error) (T1, error), F2 ~func(T1, error) (T2, error), F3 ~func(T2, error) (T3, error), F4 ~func(T3, error) (T4, error), F5 ~func(T4, error) (T5, error), T0, T1, T2, T3, T4, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T0, error) (T5, error)
func Fold[A, B any](onLeft func(error) B, onRight func(A) B) func(A, error) B
func FromOption[A any](onNone func() error) func(A, bool) (A, error)
func FromStrictEquals[A comparable]() func(A, error) func(A, error) bool
func GetOrElse[A any](onLeft func(error) A) func(A, error) A
func IsLeft[A any](_ A, err error) bool
func IsRight[A any](_ A, err error) bool
func Left[A any](err error) (A, error)
func Logger[A any](loggers ...*log.Logger) func(string) Operator[A, A]
func Memoize[A any](a A, err error) (A, error)
func Of[A any](a A) (A, error)
func Pipe1[F1 ~func(T0) (T1, error), T0, T1 any](t0 T0, f1 F1) (T1, error)
func Pipe2[F1 ~func(T0) (T1, error), F2 ~func(T1, error) (T2, error), T0, T1, T2 any](t0 T0, f1 F1, f2 F2) (T2, error)
func Pipe3[F1 ~func(T0) (T1, error), F2 ~func(T1, error) (T2, error), F3 ~func(T2, error) (T3, error), T0, T1, T2, T3 any](t0 T0, f1 F1, f2 F2, f3 F3) (T3, error)
func Pipe4[F1 ~func(T0) (T1, error), F2 ~func(T1, error) (T2, error), F3 ~func(T2, error) (T3, error), F4 ~func(T3, error) (T4, error), T0, T1, T2, T3, T4 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4) (T4, error)
func Pipe5[F1 ~func(T0) (T1, error), F2 ~func(T1, error) (T2, error), F3 ~func(T2, error) (T3, error), F4 ~func(T3, error) (T4, error), F5 ~func(T4, error) (T5, error), T0, T1, T2, T3, T4, T5 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) (T5, error)
func Reduce[A, B any](f func(B, A) B, initial B) func(A, error) B
func Right[A any](a A) (A, error)
func Sequence[A, HKTA, HKTRA any](
func ToError[A any](_ A, err error) error
func ToOption[A any](a A, aerr error) (A, bool)
func ToString[A any](a A, err error) string
func Traverse[A, B, HKTB, HKTRB any](
func TraverseTuple1[F1 ~func(A1) (T1, error), E, A1, T1 any](f1 F1) func(A1) (T1, error)
func TraverseTuple10[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), F4 ~func(A4) (T4, error), F5 ~func(A5) (T5, error), F6 ~func(A6) (T6, error), F7 ~func(A7) (T7, error), F8 ~func(A8) (T8, error), F9 ~func(A9) (T9, error), F10 ~func(A10) (T10, error), E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9, A10, T10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(A1, A2, A3, A4, A5, A6, A7, A8, A9, A10) (T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, error)
func TraverseTuple2[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), E, A1, T1, A2, T2 any](f1 F1, f2 F2) func(A1, A2) (T1, T2, error)
func TraverseTuple3[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), E, A1, T1, A2, T2, A3, T3 any](f1 F1, f2 F2, f3 F3) func(A1, A2, A3) (T1, T2, T3, error)
func TraverseTuple4[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), F4 ~func(A4) (T4, error), E, A1, T1, A2, T2, A3, T3, A4, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(A1, A2, A3, A4) (T1, T2, T3, T4, error)
func TraverseTuple5[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), F4 ~func(A4) (T4, error), F5 ~func(A5) (T5, error), E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(A1, A2, A3, A4, A5) (T1, T2, T3, T4, T5, error)
func TraverseTuple6[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), F4 ~func(A4) (T4, error), F5 ~func(A5) (T5, error), F6 ~func(A6) (T6, error), E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(A1, A2, A3, A4, A5, A6) (T1, T2, T3, T4, T5, T6, error)
func TraverseTuple7[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), F4 ~func(A4) (T4, error), F5 ~func(A5) (T5, error), F6 ~func(A6) (T6, error), F7 ~func(A7) (T7, error), E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(A1, A2, A3, A4, A5, A6, A7) (T1, T2, T3, T4, T5, T6, T7, error)
func TraverseTuple8[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), F4 ~func(A4) (T4, error), F5 ~func(A5) (T5, error), F6 ~func(A6) (T6, error), F7 ~func(A7) (T7, error), F8 ~func(A8) (T8, error), E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(A1, A2, A3, A4, A5, A6, A7, A8) (T1, T2, T3, T4, T5, T6, T7, T8, error)
func TraverseTuple9[F1 ~func(A1) (T1, error), F2 ~func(A2) (T2, error), F3 ~func(A3) (T3, error), F4 ~func(A4) (T4, error), F5 ~func(A5) (T5, error), F6 ~func(A6) (T6, error), F7 ~func(A7) (T7, error), F8 ~func(A8) (T8, error), F9 ~func(A9) (T9, error), E, A1, T1, A2, T2, A3, T3, A4, T4, A5, T5, A6, T6, A7, T7, A8, T8, A9, T9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(A1, A2, A3, A4, A5, A6, A7, A8, A9) (T1, T2, T3, T4, T5, T6, T7, T8, T9, error)
func Uncurry0[R any](f func() (R, error)) func() (R, error)
func Uncurry1[T1, R any](f func(T1) (R, error)) func(T1) (R, error)
func Uncurry2[T1, T2, R any](f func(T1) func(T2) (R, error)) func(T1, T2) (R, error)
func Uncurry3[T1, T2, T3, R any](f func(T1) func(T2) func(T3) (R, error)) func(T1, T2, T3) (R, error)
func Uncurry4[T1, T2, T3, T4, R any](f func(T1) func(T2) func(T3) func(T4) (R, error)) func(T1, T2, T3, T4) (R, error)
type Applicative[A, B any] interface {
type Apply[A, B any] interface {
type Chainable[A, B any] interface {
type Endomorphism[T any] = endomorphism.Endomorphism[T]
type Functor[A, B any] interface {
func MakeFunctor[A, B any]() Functor[A, B]
type Kleisli[A, B any] = func(A) (B, error)
func FromError[A any](f func(a A) error) Kleisli[A, A]
func FromNillable[A any](e error) Kleisli[*A, *A]
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A]
func ToType[A any](onError func(any) error) Kleisli[any, A]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseArrayWithIndex[A, B any](f func(int, A) (B, error)) Kleisli[[]A, []B]
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) (B, error)) Kleisli[GA, GB]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) (B, error)) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) (B, error)) Kleisli[GA, GB]
func WithResource[R, A, ANY any](onCreate func() (R, error), onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A]
type Lens[S, T any] = lens.Lens[S, T]
type Monad[A, B any] interface {
func MakeMonad[A, B any]() Monad[A, B]
type Operator[A, B any] = func(A, error) (B, error)
func Alt[A any](that func() (A, error)) Operator[A, A]
func Ap[B, A any](fa A, faerr error) Operator[func(A) B, B]
func BiMap[A, B any](f Endomorphism[error], g func(a A) B) Operator[A, B]
func Bind[S1, S2, T any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainTo[A, B any](b B, berr error) Operator[A, B]
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func Map[A, B any](f func(A) B) Operator[A, B]
func MapLeft[A any](f Endomorphism[error]) Operator[A, A]
func MapTo[A, B any](b B) Operator[A, B]
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A]
type Option[A any] = option.Option[A]
type Pointed[A any] interface {
func MakePointed[A any]() Pointed[A]
type Predicate[A any] = predicate.Predicate[A]
```

## package `github.com/IBM/fp-go/v2/idiomatic/ioresult`

Import: `import "github.com/IBM/fp-go/v2/idiomatic/ioresult"`

Idiomatic IOResult: `func() (A, error)`.

### Exported API

```go
func ApSeq[B, A any](ma IOResult[A]) func(IOResult[func(A) B]) IOResult[B]
func ChainOptionK[A, B any](onNone Lazy[error]) func(func(A) (B, bool)) Operator[A, B]
func Eitherize0[F ~func() (R, error), R any](f F) func() IOResult[R]
func Eitherize1[F ~func(T1) (R, error), T1, R any](f F) func(T1) IOResult[R]
func Eitherize2[F ~func(T1, T2) (R, error), T1, T2, R any](f F) func(T1, T2) IOResult[R]
func Eitherize3[F ~func(T1, T2, T3) (R, error), T1, T2, T3, R any](f F) func(T1, T2, T3) IOResult[R]
func Eq[A any](eq func(A, error) func(A, error) bool) EQ.Eq[IOResult[A]]
func Fold[A, B any](onLeft func(error) IO[B], onRight io.Kleisli[A, B]) func(IOResult[A]) IO[B]
func FromOption[A any](onNone Lazy[error]) func(A, bool) IOResult[A]
func FromStrictEquals[A comparable]() EQ.Eq[IOResult[A]]
func Functor[A, B any]() functor.Functor[A, B, IOResult[A], IOResult[B]]
func GetOrElse[A any](onLeft func(error) IO[A]) func(IOResult[A]) IO[A]
func Monad[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]]
func MonadPar[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]]
func MonadSeq[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]]
func Pointed[A any]() pointed.Pointed[A, IOResult[A]]
func TraverseParTuple1[E error, F1 ~func(A1) IOResult[T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOResult[tuple.Tuple1[T1]]
func TraverseParTuple10[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], F10 ~func(A10) IOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseParTuple2[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOResult[tuple.Tuple2[T1, T2]]
func TraverseParTuple3[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseParTuple4[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseParTuple5[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseParTuple6[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseParTuple7[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseParTuple8[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseParTuple9[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseSeqTuple1[E error, F1 ~func(A1) IOResult[T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOResult[tuple.Tuple1[T1]]
func TraverseSeqTuple10[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], F10 ~func(A10) IOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseSeqTuple2[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOResult[tuple.Tuple2[T1, T2]]
func TraverseSeqTuple3[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseSeqTuple4[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseSeqTuple5[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseSeqTuple6[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseSeqTuple7[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseSeqTuple8[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseSeqTuple9[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func TraverseTuple1[E error, F1 ~func(A1) IOResult[T1], T1, A1 any](f1 F1) func(tuple.Tuple1[A1]) IOResult[tuple.Tuple1[T1]]
func TraverseTuple10[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], F10 ~func(A10) IOResult[T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9, f10 F10) func(tuple.Tuple10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func TraverseTuple2[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], T1, T2, A1, A2 any](f1 F1, f2 F2) func(tuple.Tuple2[A1, A2]) IOResult[tuple.Tuple2[T1, T2]]
func TraverseTuple3[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], T1, T2, T3, A1, A2, A3 any](f1 F1, f2 F2, f3 F3) func(tuple.Tuple3[A1, A2, A3]) IOResult[tuple.Tuple3[T1, T2, T3]]
func TraverseTuple4[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], T1, T2, T3, T4, A1, A2, A3, A4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(tuple.Tuple4[A1, A2, A3, A4]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func TraverseTuple5[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], T1, T2, T3, T4, T5, A1, A2, A3, A4, A5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(tuple.Tuple5[A1, A2, A3, A4, A5]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func TraverseTuple6[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], T1, T2, T3, T4, T5, T6, A1, A2, A3, A4, A5, A6 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6) func(tuple.Tuple6[A1, A2, A3, A4, A5, A6]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func TraverseTuple7[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], T1, T2, T3, T4, T5, T6, T7, A1, A2, A3, A4, A5, A6, A7 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7) func(tuple.Tuple7[A1, A2, A3, A4, A5, A6, A7]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func TraverseTuple8[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], T1, T2, T3, T4, T5, T6, T7, T8, A1, A2, A3, A4, A5, A6, A7, A8 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8) func(tuple.Tuple8[A1, A2, A3, A4, A5, A6, A7, A8]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func TraverseTuple9[E error, F1 ~func(A1) IOResult[T1], F2 ~func(A2) IOResult[T2], F3 ~func(A3) IOResult[T3], F4 ~func(A4) IOResult[T4], F5 ~func(A5) IOResult[T5], F6 ~func(A6) IOResult[T6], F7 ~func(A7) IOResult[T7], F8 ~func(A8) IOResult[T8], F9 ~func(A9) IOResult[T9], T1, T2, T3, T4, T5, T6, T7, T8, T9, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5, f6 F6, f7 F7, f8 F8, f9 F9) func(tuple.Tuple9[A1, A2, A3, A4, A5, A6, A7, A8, A9]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[A any] = io.IO[A]
func MonadFold[A, B any](ma IOResult[A], onLeft func(error) IO[B], onRight io.Kleisli[A, B]) IO[B]
type IOResult[A any] = func() (A, error)
func Bracket[A, B, ANY any](
func Defer[A any](gen Lazy[IOResult[A]]) IOResult[A]
func Do[S any](
func Flatten[A any](mma IOResult[IOResult[A]]) IOResult[A]
func FromEither[A any](e Result[A]) IOResult[A]
func FromIO[A any](mr IO[A]) IOResult[A]
func FromImpure(f func()) IOResult[Void]
func FromLazy[A any](mr Lazy[A]) IOResult[A]
func FromResult[A any](a A, err error) IOResult[A]
func Left[A any](l error) IOResult[A]
func LeftIO[A any](ml IO[error]) IOResult[A]
func Memoize[A any](ma IOResult[A]) IOResult[A]
func MonadAlt[A any](first IOResult[A], second Lazy[IOResult[A]]) IOResult[A]
func MonadAp[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B]
func MonadApFirst[A, B any](first IOResult[A], second IOResult[B]) IOResult[A]
func MonadApPar[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B]
func MonadApSecond[A, B any](first IOResult[A], second IOResult[B]) IOResult[B]
func MonadApSeq[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B]
func MonadBiMap[A, B any](fa IOResult[A], f Endomorphism[error], g func(A) B) IOResult[B]
func MonadChain[A, B any](fa IOResult[A], f Kleisli[A, B]) IOResult[B]
func MonadChainEitherK[A, B any](ma IOResult[A], f either.Kleisli[error, A, B]) IOResult[B]
func MonadChainFirst[A, B any](ma IOResult[A], f Kleisli[A, B]) IOResult[A]
func MonadChainFirstEitherK[A, B any](ma IOResult[A], f either.Kleisli[error, A, B]) IOResult[A]
func MonadChainFirstIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[A]
func MonadChainFirstLeft[A, B any](ma IOResult[A], f Kleisli[error, B]) IOResult[A]
func MonadChainFirstResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A]
func MonadChainIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[B]
func MonadChainLeft[A any](fa IOResult[A], f Kleisli[error, A]) IOResult[A]
func MonadChainResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[B]
func MonadChainTo[A, B any](fa IOResult[A], fb IOResult[B]) IOResult[B]
func MonadFlap[B, A any](fab IOResult[func(A) B], a A) IOResult[B]
func MonadMap[A, B any](fa IOResult[A], f func(A) B) IOResult[B]
func MonadMapLeft[A any](fa IOResult[A], f Endomorphism[error]) IOResult[A]
func MonadMapTo[A, B any](fa IOResult[A], b B) IOResult[B]
func MonadOf[A any](r A) IOResult[A]
func MonadTap[A, B any](ma IOResult[A], f Kleisli[A, B]) IOResult[A]
func MonadTapEitherK[A, B any](ma IOResult[A], f either.Kleisli[error, A, B]) IOResult[A]
func MonadTapIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[A]
func MonadTapLeft[A, B any](ma IOResult[A], f Kleisli[error, B]) IOResult[A]
func MonadTapResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A]
func Of[A any](r A) IOResult[A]
func Retrying[A any](
func Right[A any](r A) IOResult[A]
func RightIO[A any](mr IO[A]) IOResult[A]
func SequenceArray[A any](ma []IOResult[A]) IOResult[[]A]
func SequenceArrayPar[A any](ma []IOResult[A]) IOResult[[]A]
func SequenceArraySeq[A any](ma []IOResult[A]) IOResult[[]A]
func SequenceParT1[T1 any](
func SequenceParT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceParT2[T1, T2 any](
func SequenceParT3[T1, T2, T3 any](
func SequenceParT4[T1, T2, T3, T4 any](
func SequenceParT5[T1, T2, T3, T4, T5 any](
func SequenceParT6[T1, T2, T3, T4, T5, T6 any](
func SequenceParT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceParT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceParT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceParTuple1[T1 any](t tuple.Tuple1[IOResult[T1]]) IOResult[tuple.Tuple1[T1]]
func SequenceParTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9], IOResult[T10]]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceParTuple2[T1, T2 any](t tuple.Tuple2[IOResult[T1], IOResult[T2]]) IOResult[tuple.Tuple2[T1, T2]]
func SequenceParTuple3[T1, T2, T3 any](t tuple.Tuple3[IOResult[T1], IOResult[T2], IOResult[T3]]) IOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceParTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4]]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceParTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5]]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceParTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6]]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceParTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7]]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceParTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8]]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceParTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9]]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceRecord[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A]
func SequenceRecordPar[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A]
func SequenceRecordSeq[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A]
func SequenceSeqT1[T1 any](
func SequenceSeqT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceSeqT2[T1, T2 any](
func SequenceSeqT3[T1, T2, T3 any](
func SequenceSeqT4[T1, T2, T3, T4 any](
func SequenceSeqT5[T1, T2, T3, T4, T5 any](
func SequenceSeqT6[T1, T2, T3, T4, T5, T6 any](
func SequenceSeqT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceSeqT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceSeqT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceSeqTuple1[T1 any](t tuple.Tuple1[IOResult[T1]]) IOResult[tuple.Tuple1[T1]]
func SequenceSeqTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9], IOResult[T10]]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceSeqTuple2[T1, T2 any](t tuple.Tuple2[IOResult[T1], IOResult[T2]]) IOResult[tuple.Tuple2[T1, T2]]
func SequenceSeqTuple3[T1, T2, T3 any](t tuple.Tuple3[IOResult[T1], IOResult[T2], IOResult[T3]]) IOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceSeqTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4]]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceSeqTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5]]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceSeqTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6]]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceSeqTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7]]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceSeqTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8]]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceSeqTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9]]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
func SequenceT1[T1 any](
func SequenceT10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
func SequenceT2[T1, T2 any](
func SequenceT3[T1, T2, T3 any](
func SequenceT4[T1, T2, T3, T4 any](
func SequenceT5[T1, T2, T3, T4, T5 any](
func SequenceT6[T1, T2, T3, T4, T5, T6 any](
func SequenceT7[T1, T2, T3, T4, T5, T6, T7 any](
func SequenceT8[T1, T2, T3, T4, T5, T6, T7, T8 any](
func SequenceT9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
func SequenceTuple1[T1 any](t tuple.Tuple1[IOResult[T1]]) IOResult[tuple.Tuple1[T1]]
func SequenceTuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](t tuple.Tuple10[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9], IOResult[T10]]) IOResult[tuple.Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
func SequenceTuple2[T1, T2 any](t tuple.Tuple2[IOResult[T1], IOResult[T2]]) IOResult[tuple.Tuple2[T1, T2]]
func SequenceTuple3[T1, T2, T3 any](t tuple.Tuple3[IOResult[T1], IOResult[T2], IOResult[T3]]) IOResult[tuple.Tuple3[T1, T2, T3]]
func SequenceTuple4[T1, T2, T3, T4 any](t tuple.Tuple4[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4]]) IOResult[tuple.Tuple4[T1, T2, T3, T4]]
func SequenceTuple5[T1, T2, T3, T4, T5 any](t tuple.Tuple5[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5]]) IOResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
func SequenceTuple6[T1, T2, T3, T4, T5, T6 any](t tuple.Tuple6[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6]]) IOResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
func SequenceTuple7[T1, T2, T3, T4, T5, T6, T7 any](t tuple.Tuple7[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7]]) IOResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
func SequenceTuple8[T1, T2, T3, T4, T5, T6, T7, T8 any](t tuple.Tuple8[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8]]) IOResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
func SequenceTuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](t tuple.Tuple9[IOResult[T1], IOResult[T2], IOResult[T3], IOResult[T4], IOResult[T5], IOResult[T6], IOResult[T7], IOResult[T8], IOResult[T9]]) IOResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
type Kleisli[A, B any] = Reader[A, IOResult[B]]
func LogJSON[A any](prefix string) Kleisli[A, any]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayPar[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexPar[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B]
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B]
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordPar[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexPar[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B]
func TraverseRecordWithIndexSeq[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B]
func WithResource[A, R, ANY any](
type Lazy[A any] = lazy.Lazy[A]
type Monoid[A any] = monoid.Monoid[IOResult[A]]
func ApplicativeMonoid[A any](
func ApplicativeMonoidPar[A any](
func ApplicativeMonoidSeq[A any](
type Operator[A, B any] = Kleisli[IOResult[A], B]
func After[A any](timestamp time.Time) Operator[A, A]
func Alt[A any](second Lazy[IOResult[A]]) Operator[A, A]
func Ap[B, A any](ma IOResult[A]) Operator[func(A) B, B]
func ApFirst[A, B any](second IOResult[B]) Operator[A, A]
func ApPar[B, A any](ma IOResult[A]) Operator[func(A) B, B]
func ApS[S1, S2, T any](
func ApSL[S, T any](
func ApSecond[A, B any](second IOResult[B]) Operator[A, B]
func BiMap[A, B any](f Endomorphism[error], g func(A) B) Operator[A, B]
func Bind[S1, S2, T any](
func BindL[S, T any](
func BindTo[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainFirstEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A]
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A]
func ChainFirstResultK[A, B any](f result.Kleisli[A, B]) Operator[A, A]
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B]
func ChainLazyK[A, B any](f func(A) Lazy[B]) Operator[A, B]
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A]
func ChainResultK[A, B any](f result.Kleisli[A, B]) Operator[A, B]
func ChainTo[A, B any](fb IOResult[B]) Operator[A, B]
func Delay[A any](delay time.Duration) Operator[A, A]
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func Map[A, B any](f func(A) B) Operator[A, B]
func MapLeft[A any](f Endomorphism[error]) Operator[A, A]
func MapTo[A, B any](b B) Operator[A, B]
func Tap[A, B any](f Kleisli[A, B]) Operator[A, A]
func TapEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A]
func TapIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A]
func TapLeft[A, B any](f Kleisli[error, B]) Operator[A, A]
func TapResultK[A, B any](f result.Kleisli[A, B]) Operator[A, A]
func WithLock[A any](lock IO[context.CancelFunc]) Operator[A, A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type Result[A any] = result.Result[A]
type Semigroup[A any] = semigroup.Semigroup[IOResult[A]]
func AltSemigroup[A any]() Semigroup[A]
type Void = function.Void
```

## package `github.com/IBM/fp-go/v2/idiomatic/readerresult`

Import: `import "github.com/IBM/fp-go/v2/idiomatic/readerresult"`

Idiomatic ReaderResult: `func(context.Context) (A, error)`.

### Exported API

```go
func ApResultS[
func BindToEither[
func BindToReader[
func BindToResult[
func ChainOptionK[R, A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[R, A, B]
func Curry1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderResult[R, A]
func Curry2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1) func(T2) ReaderResult[R, A]
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderResult[R, A]
func Fold[R, A, B any](onLeft reader.Kleisli[R, error, B], onRight reader.Kleisli[R, A, B]) func(ReaderResult[R, A]) Reader[R, B]
func From0[R, A any](f func(R) (A, error)) func() ReaderResult[R, A]
func From1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderResult[R, A]
func From2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) ReaderResult[R, A]
func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderResult[R, A]
func GetOrElse[R, A any](onLeft reader.Kleisli[R, error, A]) func(ReaderResult[R, A]) Reader[R, A]
func LetTo[R, S1, S2, T any](
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderResult[R1, A]) ReaderResult[R2, A]
func Read[A, R any](r R) func(ReaderResult[R, A]) (A, error)
func Traverse[R2, R1, A, B any](
func TraverseReader[R2, R1, A, B any](
func Uncurry1[R, T1, A any](f func(T1) ReaderResult[R, A]) func(R, T1) (A, error)
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) ReaderResult[R, A]) func(R, T1, T2) (A, error)
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderResult[R, A]) func(R, T1, T2, T3) (A, error)
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type Kleisli[R, A, B any] = Reader[A, ReaderResult[R, B]]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderResult[R1, A], A]
func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) error) Kleisli[R, A, A]
func Promap[E, A, D, B any](f func(D) E, g func(A) B) Kleisli[D, ReaderResult[E, A], B]
func Sequence[R1, R2, A any](ma ReaderResult[R2, ReaderResult[R1, A]]) Kleisli[R2, R1, A]
func SequenceReader[R1, R2, A any](ma ReaderResult[R2, Reader[R1, A]]) Kleisli[R2, R1, A]
func TraverseArray[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, []A, []B]
func TraverseArrayWithIndex[R, A, B any](f func(int, A) ReaderResult[R, B]) Kleisli[R, []A, []B]
func WithResource[B, R, A, ANY any](
type Lazy[A any] = lazy.Lazy[A]
type Monoid[R, A any] = monoid.Monoid[ReaderResult[R, A]]
func AltMonoid[R, A any](zero Lazy[ReaderResult[R, A]]) Monoid[R, A]
func AlternativeMonoid[R, A any](m M.Monoid[A]) Monoid[R, A]
func ApplicativeMonoid[R, A any](m M.Monoid[A]) Monoid[R, A]
type Operator[R, A, B any] = Kleisli[R, ReaderResult[R, A], B]
func Alt[R, A any](second Lazy[ReaderResult[R, A]]) Operator[R, A, A]
func Ap[B, R, A any](fa ReaderResult[R, A]) Operator[R, func(A) B, B]
func ApEitherS[
func ApReaderS[
func ApS[R, S1, S2, T any](
func ApSL[R, S, T any](
func BiMap[R, A, B any](f Endomorphism[error], g func(A) B) Operator[R, A, B]
func Bind[R, S1, S2, T any](
func BindEitherK[R, S1, S2, T any](
func BindL[R, S, T any](
func BindReaderK[R, S1, S2, T any](
func BindResultK[R, S1, S2, T any](
func BindTo[R, S1, T any](
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainEitherK[R, A, B any](f RES.Kleisli[A, B]) Operator[R, A, B]
func ChainReaderK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B]
func FilterOrElse[R, A any](pred Predicate[A], onFalse func(A) error) Operator[R, A, A]
func Flap[R, B, A any](a A) Operator[R, func(A) B, B]
func Let[R, S1, S2, T any](
func LetL[R, S, T any](
func LetToL[R, S, T any](
func Map[R, A, B any](f func(A) B) Operator[R, A, B]
func MapLeft[R, A any](f Endomorphism[error]) Operator[R, A, A]
func OrElse[R, A any](onLeft Kleisli[R, error, A]) Operator[R, A, A]
func OrLeft[A, R any](onLeft reader.Kleisli[R, error, error]) Operator[R, A, A]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderResult[R, A any] = func(R) (A, error)
func Ask[R any]() ReaderResult[R, R]
func Asks[R, A any](r Reader[R, A]) ReaderResult[R, A]
func Bracket[
func Curry0[R, A any](f func(R) (A, error)) ReaderResult[R, A]
func Do[R, S any](
func Flatten[R, A any](mma ReaderResult[R, ReaderResult[R, A]]) ReaderResult[R, A]
func FromEither[R, A any](e Result[A]) ReaderResult[R, A]
func FromReader[R, A any](r Reader[R, A]) ReaderResult[R, A]
func FromResult[R, A any](a A, err error) ReaderResult[R, A]
func Left[R, A any](err error) ReaderResult[R, A]
func LeftReader[A, R any](l Reader[R, error]) ReaderResult[R, A]
func MonadAlt[R, A any](first ReaderResult[R, A], second Lazy[ReaderResult[R, A]]) ReaderResult[R, A]
func MonadAp[B, R, A any](fab ReaderResult[R, func(A) B], fa ReaderResult[R, A]) ReaderResult[R, B]
func MonadBiMap[R, A, B any](fa ReaderResult[R, A], f Endomorphism[error], g func(A) B) ReaderResult[R, B]
func MonadChain[R, A, B any](ma ReaderResult[R, A], f Kleisli[R, A, B]) ReaderResult[R, B]
func MonadChainEitherK[R, A, B any](ma ReaderResult[R, A], f RES.Kleisli[A, B]) ReaderResult[R, B]
func MonadChainReaderK[R, A, B any](ma ReaderResult[R, A], f result.Kleisli[A, B]) ReaderResult[R, B]
func MonadFlap[R, A, B any](fab ReaderResult[R, func(A) B], a A) ReaderResult[R, B]
func MonadMap[R, A, B any](fa ReaderResult[R, A], f func(A) B) ReaderResult[R, B]
func MonadMapLeft[R, A any](fa ReaderResult[R, A], f Endomorphism[error]) ReaderResult[R, A]
func MonadTraverseArray[R, A, B any](as []A, f Kleisli[R, A, B]) ReaderResult[R, []B]
func Of[R, A any](a A) ReaderResult[R, A]
func Right[R, A any](a A) ReaderResult[R, A]
func RightReader[R, A any](rdr Reader[R, A]) ReaderResult[R, A]
func SequenceArray[R, A any](ma []ReaderResult[R, A]) ReaderResult[R, []A]
func SequenceT1[R, A any](a ReaderResult[R, A]) ReaderResult[R, T.Tuple1[A]]
func SequenceT2[R, A, B any](
func SequenceT3[R, A, B, C any](
func SequenceT4[R, A, B, C, D any](
type Result[A any] = result.Result[A]
```

## package `github.com/IBM/fp-go/v2/idiomatic/readerioresult`

Import: `import "github.com/IBM/fp-go/v2/idiomatic/readerioresult"`

Idiomatic ReaderIOResult: `func(context.Context) func() (A, error)`.

### Exported API

```go
type with no runtime overhead beyond the underlying computation. The IO
func Ap[B, R, A any](fa ReaderIOResult[R, A]) func(fab ReaderIOResult[R, func(A) B]) ReaderIOResult[R, B]
func Flap[R, B, A any](a A) func(ReaderIOResult[R, func(A) B]) ReaderIOResult[R, B]
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, A]
func Read[A, R any](r R) func(ReaderIOResult[R, A]) IOResult[A]
func ReadIO[A, R any](r IO[R]) func(ReaderIOResult[R, A]) IOResult[A]
func ReadIOResult[A, R any](r IOResult[R]) func(ReaderIOResult[R, A]) IOResult[A]
func Sequence[R1, R2, A any](ma ReaderIOResult[R2, ReaderIOResult[R1, A]]) reader.Kleisli[R2, R1, IOResult[A]]
func SequenceReader[R1, R2, A any](ma ReaderIOResult[R2, Reader[R1, A]]) reader.Kleisli[R2, R1, IOResult[A]]
func Traverse[R2, R1, A, B any](
func TraverseReader[R2, R1, A, B any](
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type IO[A any] = io.IO[A]
type IOResult[A any] = ioresult.IOResult[A]
type Kleisli[R, A, B any] = Reader[A, ReaderIOResult[R, B]]
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIOResult[R1, A], A]
func Promap[E, A, D, B any](f func(D) E, g func(A) B) Kleisli[D, ReaderIOResult[E, A], B]
type Lazy[A any] = lazy.Lazy[A]
type Monoid[R, A any] = monoid.Monoid[ReaderIOResult[R, A]]
type Operator[R, A, B any] = Kleisli[R, ReaderIOResult[R, A], B]
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, B]
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A]
func ChainFirstEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A]
func ChainFirstIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A]
func ChainFirstReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A]
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A]
func ChainIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, B]
func ChainReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, B]
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B]
func FilterOrElse[R, A any](pred Predicate[A], onFalse func(A) error) Operator[R, A, A]
func Map[R, A, B any](f func(A) B) Operator[R, A, B]
func MapTo[R, A, B any](b B) Operator[R, A, B]
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A]
func TapEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A]
func TapIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A]
func TapReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A]
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderIO[R, A any] = readerio.ReaderIO[R, A]
type ReaderIOResult[R, A any] = Reader[R, IOResult[A]]
func Ask[R any]() ReaderIOResult[R, R]
func Asks[R, A any](r Reader[R, A]) ReaderIOResult[R, A]
func Flatten[R, A any](mma ReaderIOResult[R, ReaderIOResult[R, A]]) ReaderIOResult[R, A]
func FromEither[R, A any](t either.Either[error, A]) ReaderIOResult[R, A]
func FromIO[R, E, A any](ma IO[A]) ReaderIOResult[R, A]
func FromIOResult[R, A any](ma IOResult[A]) ReaderIOResult[R, A]
func FromReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A]
func FromReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A]
func Left[R, A any](e error) ReaderIOResult[R, A]
func LeftIO[R, A any](ma IO[error]) ReaderIOResult[R, A]
func LeftReader[A, R any](ma Reader[R, error]) ReaderIOResult[R, A]
func MonadAp[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B]
func MonadApPar[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B]
func MonadApSeq[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B]
func MonadChain[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, B]
func MonadChainEitherK[R, A, B any](ma ReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderIOResult[R, B]
func MonadChainFirst[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadChainFirstEitherK[R, A, B any](ma ReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderIOResult[R, A]
func MonadChainFirstIOK[R, A, B any](ma ReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderIOResult[R, A]
func MonadChainFirstReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadChainFirstReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadChainIOK[R, A, B any](ma ReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderIOResult[R, B]
func MonadChainReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, B]
func MonadChainReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, B]
func MonadFlap[R, B, A any](fab ReaderIOResult[R, func(A) B], a A) ReaderIOResult[R, B]
func MonadMap[R, A, B any](fa ReaderIOResult[R, A], f func(A) B) ReaderIOResult[R, B]
func MonadMapTo[R, A, B any](fa ReaderIOResult[R, A], b B) ReaderIOResult[R, B]
func MonadTap[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadTapEitherK[R, A, B any](ma ReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderIOResult[R, A]
func MonadTapIOK[R, A, B any](ma ReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderIOResult[R, A]
func MonadTapReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A]
func MonadTapReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A]
func Of[R, A any](a A) ReaderIOResult[R, A]
func Right[R, A any](a A) ReaderIOResult[R, A]
func RightIO[R, A any](ma IO[A]) ReaderIOResult[R, A]
func RightReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A]
func RightReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A]
type Result[A any] = result.Result[A]
```

## package `github.com/IBM/fp-go/v2/idiomatic/context/readerresult`

Import: `import "github.com/IBM/fp-go/v2/idiomatic/context/readerresult"`

Idiomatic context-specialized ReaderResult.

### Exported API

```go
func ApResultS[
func BindToEither[
func BindToReader[
func BindToResult[
func ChainOptionK[A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[A, B]
func Curry1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderResult[A]
func Curry2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1) func(T2) ReaderResult[A]
func Curry3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderResult[A]
func Fold[A, B any](onLeft reader.Kleisli[context.Context, error, B], onRight reader.Kleisli[context.Context, A, B]) func(ReaderResult[A]) Reader[context.Context, B]
func From0[A any](f func(context.Context) (A, error)) func() ReaderResult[A]
func From1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderResult[A]
func From2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1, T2) ReaderResult[A]
func From3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderResult[A]
func GetOrElse[A any](onLeft reader.Kleisli[context.Context, error, A]) func(ReaderResult[A]) Reader[context.Context, A]
func Read[A any](ctx context.Context) func(ReaderResult[A]) (A, error)
func ToReaderResult[A any](r ReaderResult[A]) RS.ReaderResult[A]
func TraverseReader[R, A, B any](
func Uncurry1[T1, A any](f func(T1) ReaderResult[A]) func(context.Context, T1) (A, error)
func Uncurry2[T1, T2, A any](f func(T1) func(T2) ReaderResult[A]) func(context.Context, T1, T2) (A, error)
func Uncurry3[T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderResult[A]) func(context.Context, T1, T2, T3) (A, error)
type Either[E, A any] = either.Either[E, A]
type Endomorphism[A any] = endomorphism.Endomorphism[A]
type Kleisli[A, B any] = Reader[A, ReaderResult[B]]
func Contramap[A any](f func(context.Context) (context.Context, context.CancelFunc)) Kleisli[ReaderResult[A], A]
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A]
func SequenceReader[R, A any](ma ReaderResult[Reader[R, A]]) Kleisli[R, A]
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B]
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B]
func TraverseArrayWithIndex[A, B any](f func(int, A) ReaderResult[B]) Kleisli[[]A, []B]
func WithCloser[B any, A io.Closer](onCreate Lazy[ReaderResult[A]]) Kleisli[Kleisli[A, B], B]
func WithContextK[A, B any](f Kleisli[A, B]) Kleisli[A, B]
func WithResource[B, A, ANY any](
type Lazy[A any] = lazy.Lazy[A]
type Lens[S, A any] = lens.Lens[S, A]
type Monoid[A any] = monoid.Monoid[ReaderResult[A]]
func AltMonoid[A any](zero Lazy[ReaderResult[A]]) Monoid[A]
func AlternativeMonoid[A any](m M.Monoid[A]) Monoid[A]
func ApplicativeMonoid[A any](m M.Monoid[A]) Monoid[A]
type Operator[A, B any] = Kleisli[ReaderResult[A], B]
func Alt[A any](second Lazy[ReaderResult[A]]) Operator[A, A]
func Ap[B, A any](fa ReaderResult[A]) Operator[func(A) B, B]
func ApEitherS[
func ApReaderS[
func ApS[S1, S2, T any](
func ApSL[S, T any](
func BiMap[A, B any](f Endomorphism[error], g func(A) B) Operator[A, B]
func Bind[S1, S2, T any](
func BindEitherK[S1, S2, T any](
func BindL[S, T any](
func BindReaderK[S1, S2, T any](
func BindResultK[S1, S2, T any](
func BindTo[S1, T any](
func BindToP[S1, T any](
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainEitherK[A, B any](f RES.Kleisli[A, B]) Operator[A, B]
func ChainReaderK[A, B any](f result.Kleisli[A, B]) Operator[A, B]
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A]
func Flap[B, A any](a A) Operator[func(A) B, B]
func Let[S1, S2, T any](
func LetL[S, T any](
func LetTo[S1, S2, T any](
func LetToL[S, T any](
func Local[A any](f func(context.Context) (context.Context, context.CancelFunc)) Operator[A, A]
func Map[A, B any](f func(A) B) Operator[A, B]
func MapLeft[A any](f Endomorphism[error]) Operator[A, A]
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A]
func OrLeft[A any](onLeft reader.Kleisli[context.Context, error, error]) Operator[A, A]
func Promap[A, B any](f func(context.Context) (context.Context, context.CancelFunc), g func(A) B) Operator[A, B]
func WithDeadline[A any](deadline time.Time) Operator[A, A]
func WithTimeout[A any](timeout time.Duration) Operator[A, A]
type Option[A any] = option.Option[A]
type Predicate[A any] = predicate.Predicate[A]
type Prism[S, A any] = prism.Prism[S, A]
type Reader[R, A any] = reader.Reader[R, A]
type ReaderResult[A any] = func(context.Context) (A, error)
func Ask() ReaderResult[context.Context]
func Asks[A any](r Reader[context.Context, A]) ReaderResult[A]
func Bracket[
func Curry0[A any](f func(context.Context) (A, error)) ReaderResult[A]
func Do[S any](
func Flatten[A any](mma ReaderResult[ReaderResult[A]]) ReaderResult[A]
func FromEither[A any](e Result[A]) ReaderResult[A]
func FromReader[A any](r Reader[context.Context, A]) ReaderResult[A]
func FromReaderResult[A any](r RS.ReaderResult[A]) ReaderResult[A]
func FromResult[A any](a A, err error) ReaderResult[A]
func Left[A any](err error) ReaderResult[A]
func LeftReader[A, R any](l Reader[context.Context, error]) ReaderResult[A]
func MonadAlt[A any](first ReaderResult[A], second Lazy[ReaderResult[A]]) ReaderResult[A]
func MonadAp[B, A any](fab ReaderResult[func(A) B], fa ReaderResult[A]) ReaderResult[B]
func MonadBiMap[A, B any](fa ReaderResult[A], f Endomorphism[error], g func(A) B) ReaderResult[B]
func MonadChain[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[B]
func MonadChainEitherK[A, B any](ma ReaderResult[A], f RES.Kleisli[A, B]) ReaderResult[B]
func MonadChainReaderK[A, B any](ma ReaderResult[A], f result.Kleisli[A, B]) ReaderResult[B]
func MonadFlap[A, B any](fab ReaderResult[func(A) B], a A) ReaderResult[B]
func MonadMap[A, B any](fa ReaderResult[A], f func(A) B) ReaderResult[B]
func MonadMapLeft[A any](fa ReaderResult[A], f Endomorphism[error]) ReaderResult[A]
func MonadTraverseArray[A, B any](as []A, f Kleisli[A, B]) ReaderResult[[]B]
func Of[A any](a A) ReaderResult[A]
func Retrying[A any](
func Right[A any](a A) ReaderResult[A]
func RightReader[A any](rdr Reader[context.Context, A]) ReaderResult[A]
func SequenceArray[A any](ma []ReaderResult[A]) ReaderResult[[]A]
func SequenceT1[A any](a ReaderResult[A]) ReaderResult[T.Tuple1[A]]
func SequenceT2[A, B any](
func SequenceT3[A, B, C any](
func SequenceT4[A, B, C, D any](
func WithContext[A any](ma ReaderResult[A]) ReaderResult[A]
type Result[A any] = result.Result[A]
type Trampoline[A, B any] = tailrec.Trampoline[A, B]
```
