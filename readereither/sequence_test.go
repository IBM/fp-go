package readereither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/either"
	T "github.com/IBM/fp-go/tuple"
	"github.com/stretchr/testify/assert"
)

var (
	testError = fmt.Errorf("error")
)

func TestSequenceT1(t *testing.T) {

	t1 := Of[MyContext, error]("s1")
	e1 := Left[MyContext, string](testError)

	res1 := SequenceT1(t1)
	assert.Equal(t, E.Of[error](T.MakeTuple1("s1")), res1(defaultContext))

	res2 := SequenceT1(e1)
	assert.Equal(t, E.Left[T.Tuple1[string]](testError), res2(defaultContext))
}

func TestSequenceT2(t *testing.T) {

	t1 := Of[MyContext, error]("s1")
	e1 := Left[MyContext, string](testError)
	t2 := Of[MyContext, error](2)
	e2 := Left[MyContext, int](testError)

	res1 := SequenceT2(t1, t2)
	assert.Equal(t, E.Of[error](T.MakeTuple2("s1", 2)), res1(defaultContext))

	res2 := SequenceT2(e1, t2)
	assert.Equal(t, E.Left[T.Tuple2[string, int]](testError), res2(defaultContext))

	res3 := SequenceT2(t1, e2)
	assert.Equal(t, E.Left[T.Tuple2[string, int]](testError), res3(defaultContext))
}

func TestSequenceT3(t *testing.T) {

	t1 := Of[MyContext, error]("s1")
	e1 := Left[MyContext, string](testError)
	t2 := Of[MyContext, error](2)
	e2 := Left[MyContext, int](testError)
	t3 := Of[MyContext, error](true)
	e3 := Left[MyContext, bool](testError)

	res1 := SequenceT3(t1, t2, t3)
	assert.Equal(t, E.Of[error](T.MakeTuple3("s1", 2, true)), res1(defaultContext))

	res2 := SequenceT3(e1, t2, t3)
	assert.Equal(t, E.Left[T.Tuple3[string, int, bool]](testError), res2(defaultContext))

	res3 := SequenceT3(t1, e2, t3)
	assert.Equal(t, E.Left[T.Tuple3[string, int, bool]](testError), res3(defaultContext))

	res4 := SequenceT3(t1, t2, e3)
	assert.Equal(t, E.Left[T.Tuple3[string, int, bool]](testError), res4(defaultContext))
}

func TestSequenceT4(t *testing.T) {

	t1 := Of[MyContext, error]("s1")
	t2 := Of[MyContext, error](2)
	t3 := Of[MyContext, error](true)
	t4 := Of[MyContext, error](1.0)

	res := SequenceT4(t1, t2, t3, t4)

	assert.Equal(t, E.Of[error](T.MakeTuple4("s1", 2, true, 1.0)), res(defaultContext))
}
