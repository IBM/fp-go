package readerioeither

import (
	"context"
	"time"

	R "github.com/ibm/fp-go/context/reader"
	RIO "github.com/ibm/fp-go/context/readerio"
	ET "github.com/ibm/fp-go/either"
	ER "github.com/ibm/fp-go/errors"
	F "github.com/ibm/fp-go/function"
	IO "github.com/ibm/fp-go/io"
	IOE "github.com/ibm/fp-go/ioeither"
	O "github.com/ibm/fp-go/option"
	RIE "github.com/ibm/fp-go/readerioeither/generic"
)

func FromEither[A any](e ET.Either[error, A]) ReaderIOEither[A] {
	return RIE.FromEither[ReaderIOEither[A]](e)
}

func RightReader[A any](r R.Reader[A]) ReaderIOEither[A] {
	return RIE.RightReader[R.Reader[A], ReaderIOEither[A]](r)
}

func LeftReader[A any](l R.Reader[error]) ReaderIOEither[A] {
	return RIE.LeftReader[R.Reader[error], ReaderIOEither[A]](l)
}

func Left[A any](l error) ReaderIOEither[A] {
	return RIE.Left[ReaderIOEither[A]](l)
}

func Right[A any](r A) ReaderIOEither[A] {
	return RIE.Right[ReaderIOEither[A]](r)
}

func FromReader[A any](r R.Reader[A]) ReaderIOEither[A] {
	return RIE.FromReader[R.Reader[A], ReaderIOEither[A]](r)
}

func MonadMap[A, B any](fa ReaderIOEither[A], f func(A) B) ReaderIOEither[B] {
	return RIE.MonadMap[ReaderIOEither[A], ReaderIOEither[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return RIE.Map[ReaderIOEither[A], ReaderIOEither[B]](f)
}

func MonadChain[A, B any](ma ReaderIOEither[A], f func(A) ReaderIOEither[B]) ReaderIOEither[B] {
	return RIE.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) ReaderIOEither[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return RIE.Chain[ReaderIOEither[A]](f)
}

func Of[A any](a A) ReaderIOEither[A] {
	return RIE.Of[ReaderIOEither[A]](a)
}

// withCancelCauseFunc wraps an IOEither such that in case of an error the cancel function is invoked
func withCancelCauseFunc[A any](cancel context.CancelCauseFunc, ma IOE.IOEither[error, A]) IOE.IOEither[error, A] {
	return F.Pipe3(
		ma,
		IOE.Swap[error, A],
		IOE.ChainFirstIOK[A, error, any](func(err error) IO.IO[any] {
			return IO.MakeIO(func() any {
				cancel(err)
				return nil
			})
		}),
		IOE.Swap[A, error],
	)
}

// MonadAp implements the `Ap` function for a reader with context. It creates a sub-context that will
// be canceled if any of the input operations errors out or
func MonadAp[A, B any](fab ReaderIOEither[func(A) B], fa ReaderIOEither[A]) ReaderIOEither[B] {
	// context sensitive input
	cfab := WithContext(fab)
	cfa := WithContext(fa)

	return func(ctx context.Context) IOE.IOEither[error, B] {
		// quick check for cancellation
		if err := context.Cause(ctx); err != nil {
			return IOE.Left[B](err)
		}

		return func() ET.Either[error, B] {
			// quick check for cancellation
			if err := context.Cause(ctx); err != nil {
				return ET.Left[B](err)
			}

			// create sub-contexts for fa and fab, so they can cancel one other
			ctxSub, cancelSub := context.WithCancelCause(ctx)
			defer cancelSub(nil) // cancel has to be called in all paths

			fabIOE := withCancelCauseFunc(cancelSub, cfab(ctxSub))
			faIOE := withCancelCauseFunc(cancelSub, cfa(ctxSub))

			return IOE.MonadAp(fabIOE, faIOE)()
		}
	}
}

func Ap[A, B any](fa ReaderIOEither[A]) func(ReaderIOEither[func(A) B]) ReaderIOEither[B] {
	return F.Bind2nd(MonadAp[A, B], fa)
}

func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) func(A) ReaderIOEither[A] {
	return RIE.FromPredicate[ReaderIOEither[A]](pred, onFalse)
}

func Fold[A, B any](onLeft func(error) RIO.ReaderIO[B], onRight func(A) RIO.ReaderIO[B]) func(ReaderIOEither[A]) RIO.ReaderIO[B] {
	return RIE.Fold[RIO.ReaderIO[B], ReaderIOEither[A]](onLeft, onRight)
}

func GetOrElse[A any](onLeft func(error) RIO.ReaderIO[A]) func(ReaderIOEither[A]) RIO.ReaderIO[A] {
	return RIE.GetOrElse[RIO.ReaderIO[A], ReaderIOEither[A]](onLeft)
}

func OrElse[A any](onLeft func(error) ReaderIOEither[A]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return RIE.OrElse[ReaderIOEither[A]](onLeft)
}

func OrLeft[A any](onLeft func(error) RIO.ReaderIO[error]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return RIE.OrLeft[ReaderIOEither[A], RIO.ReaderIO[error], ReaderIOEither[A]](onLeft)
}

func Ask() ReaderIOEither[context.Context] {
	return RIE.Ask[ReaderIOEither[context.Context]]()
}

func Asks[A any](r R.Reader[A]) ReaderIOEither[A] {
	return RIE.Asks[R.Reader[A], ReaderIOEither[A]](r)
}

func MonadChainEitherK[A, B any](ma ReaderIOEither[A], f func(A) ET.Either[error, B]) ReaderIOEither[B] {
	return RIE.MonadChainEitherK[ReaderIOEither[A], ReaderIOEither[B]](ma, f)
}

func ChainEitherK[A, B any](f func(A) ET.Either[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return RIE.ChainEitherK[ReaderIOEither[A], ReaderIOEither[B]](f)
}

func ChainOptionK[A, B any](onNone func() error) func(func(A) O.Option[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return RIE.ChainOptionK[ReaderIOEither[A], ReaderIOEither[B]](onNone)
}

func FromIOEither[A any](t IOE.IOEither[error, A]) ReaderIOEither[A] {
	return RIE.FromIOEither[ReaderIOEither[A]](t)
}

func FromIO[A any](t IO.IO[A]) ReaderIOEither[A] {
	return RIE.FromIO[ReaderIOEither[A]](t)
}

// Never returns a 'ReaderIOEither' that never returns, except if its context gets canceled
func Never[A any]() ReaderIOEither[A] {
	return func(ctx context.Context) IOE.IOEither[error, A] {
		return IOE.MakeIO(func() ET.Either[error, A] {
			<-ctx.Done()
			return ET.Left[A](context.Cause(ctx))
		})
	}
}

func MonadChainIOK[A, B any](ma ReaderIOEither[A], f func(A) IO.IO[B]) ReaderIOEither[B] {
	return RIE.MonadChainIOK[ReaderIOEither[A], ReaderIOEither[B]](ma, f)
}

func ChainIOK[A, B any](f func(A) IO.IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return RIE.ChainIOK[ReaderIOEither[A], ReaderIOEither[B]](f)
}

func ChainIOEitherK[A, B any](f func(A) IOE.IOEither[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return RIE.ChainIOEitherK[ReaderIOEither[A], ReaderIOEither[B]](f)
}

// Delay creates an operation that passes in the value after some delay
func Delay[A any](delay time.Duration) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return func(ma ReaderIOEither[A]) ReaderIOEither[A] {
		return func(ctx context.Context) IOE.IOEither[error, A] {
			return IOE.MakeIO(func() ET.Either[error, A] {
				// manage the timeout
				timeoutCtx, cancelTimeout := context.WithTimeout(ctx, delay)
				defer cancelTimeout()
				// whatever comes first
				select {
				case <-timeoutCtx.Done():
					return ma(ctx)()
				case <-ctx.Done():
					return ET.Left[A](context.Cause(ctx))
				}
			})
		}
	}
}

// Timer will return the current time after an initial delay
func Timer(delay time.Duration) ReaderIOEither[time.Time] {
	return F.Pipe2(
		IO.Now,
		FromIO[time.Time],
		Delay[time.Time](delay),
	)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() ReaderIOEither[A]) ReaderIOEither[A] {
	return RIE.Defer[ReaderIOEither[A]](gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[A any](f func(context.Context) func() (A, error)) ReaderIOEither[A] {
	return RIE.TryCatch[ReaderIOEither[A]](f, ER.IdentityError)
}
