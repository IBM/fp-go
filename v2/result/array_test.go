package result

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	TST "github.com/IBM/fp-go/v2/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestCompactArray(t *testing.T) {
	ar := []Result[string]{
		Of("ok"),
		Left[string](errors.New("err")),
		Of("ok"),
	}
	assert.Equal(t, 2, len(CompactArray(ar)))
}

func TestSequenceArray(t *testing.T) {
	s := TST.SequenceArrayTest(
		FromStrictEquals[bool](),
		Pointed[string](),
		Pointed[bool](),
		Functor[[]string, bool](),
		SequenceArray[string],
	)
	for i := range 10 {
		t.Run(fmt.Sprintf("TestSequenceArray %d", i), s(i))
	}
}

func TestSequenceArrayError(t *testing.T) {
	s := TST.SequenceArrayErrorTest(
		FromStrictEquals[bool](),
		Left[string],
		Left[bool],
		Pointed[string](),
		Pointed[bool](),
		Functor[[]string, bool](),
		SequenceArray[string],
	)
	s(4)(t)
}

func TestTraverseSeq_Success(t *testing.T) {
	parse := func(s string) Result[int] {
		v, err := strconv.Atoi(s)
		return TryCatchError(v, err)
	}

	collectInts := func(result Result[iter.Seq[int]]) []int {
		return F.Pipe1(result, Fold(
			func(e error) []int { t.Fatal(e); return nil },
			slices.Collect[int],
		))
	}

	t.Run("transforms all elements successfully", func(t *testing.T) {
		result := TraverseSeq(parse)(slices.Values([]string{"1", "2", "3"}))
		assert.Equal(t, []int{1, 2, 3}, collectInts(result))
	})

	t.Run("works with empty iterator", func(t *testing.T) {
		result := TraverseSeq(parse)(slices.Values([]string{}))
		assert.Empty(t, collectInts(result))
	})

	t.Run("works with single element", func(t *testing.T) {
		result := TraverseSeq(parse)(slices.Values([]string{"42"}))
		assert.Equal(t, []int{42}, collectInts(result))
	})

	t.Run("preserves order of elements", func(t *testing.T) {
		result := TraverseSeq(parse)(slices.Values([]string{"10", "20", "30", "40", "50"}))
		assert.Equal(t, []int{10, 20, 30, 40, 50}, collectInts(result))
	})
}

func TestTraverseSeq_Failure(t *testing.T) {
	parse := func(s string) Result[int] {
		v, err := strconv.Atoi(s)
		return TryCatchError(v, err)
	}

	extractErr := func(result Result[iter.Seq[int]]) error {
		return F.Pipe1(result, Fold(
			F.Identity[error],
			func(_ iter.Seq[int]) error { t.Fatal("expected Left but got Right"); return nil },
		))
	}

	t.Run("short-circuits on first Left", func(t *testing.T) {
		err := extractErr(TraverseSeq(parse)(slices.Values([]string{"1", "invalid", "3"})))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("returns first error when multiple failures exist", func(t *testing.T) {
		err := extractErr(TraverseSeq(parse)(slices.Values([]string{"1", "bad1", "bad2"})))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bad1")
	})

	t.Run("handles custom error types", func(t *testing.T) {
		customErr := errors.New("custom validation error")
		validate := func(n int) Result[int] {
			if n == 2 {
				return Left[int](customErr)
			}
			return Of(n * 10)
		}
		err := extractErr(TraverseSeq(validate)(slices.Values([]int{1, 2, 3})))
		assert.Equal(t, customErr, err)
	})
}

func TestTraverseSeq_EdgeCases(t *testing.T) {
	t.Run("handles complex transformations", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		transform := func(id int) Result[User] {
			return Of(User{ID: id, Name: fmt.Sprintf("User%d", id)})
		}

		result := TraverseSeq(transform)(slices.Values([]int{1, 2, 3}))
		collected := F.Pipe1(result, Fold(
			func(e error) []User { t.Fatal(e); return nil },
			slices.Collect[User],
		))

		assert.Equal(t, []User{
			{ID: 1, Name: "User1"},
			{ID: 2, Name: "User2"},
			{ID: 3, Name: "User3"},
		}, collected)
	})

	t.Run("works with identity transformation", func(t *testing.T) {
		input := slices.Values([]Result[int]{Of(1), Of(2), Of(3)})

		result := TraverseSeq(F.Identity[Result[int]])(input)
		collected := F.Pipe1(result, Fold(
			func(e error) []int { t.Fatal(e); return nil },
			slices.Collect[int],
		))

		assert.Equal(t, []int{1, 2, 3}, collected)
	})
}

func TestSequenceSeq_Success(t *testing.T) {
	collectInts := func(result Result[iter.Seq[int]]) []int {
		return F.Pipe1(result, Fold(
			func(e error) []int { t.Fatal(e); return nil },
			slices.Collect[int],
		))
	}

	t.Run("sequences multiple Right values", func(t *testing.T) {
		input := slices.Values([]Result[int]{Of(1), Of(2), Of(3)})
		assert.Equal(t, []int{1, 2, 3}, collectInts(SequenceSeq(input)))
	})

	t.Run("works with empty iterator", func(t *testing.T) {
		input := slices.Values([]Result[string]{})
		result := F.Pipe1(SequenceSeq(input), Fold(
			func(e error) []string { t.Fatal(e); return nil },
			slices.Collect[string],
		))
		assert.Empty(t, result)
	})

	t.Run("works with single Right value", func(t *testing.T) {
		input := slices.Values([]Result[string]{Of("hello")})
		result := F.Pipe1(SequenceSeq(input), Fold(
			func(e error) []string { t.Fatal(e); return nil },
			slices.Collect[string],
		))
		assert.Equal(t, []string{"hello"}, result)
	})

	t.Run("preserves order of results", func(t *testing.T) {
		input := slices.Values([]Result[int]{Of(5), Of(4), Of(3), Of(2), Of(1)})
		assert.Equal(t, []int{5, 4, 3, 2, 1}, collectInts(SequenceSeq(input)))
	})

	t.Run("works with complex types", func(t *testing.T) {
		type Item struct {
			Value int
			Label string
		}

		input := slices.Values([]Result[Item]{
			Of(Item{Value: 1, Label: "first"}),
			Of(Item{Value: 2, Label: "second"}),
			Of(Item{Value: 3, Label: "third"}),
		})

		collected := F.Pipe1(SequenceSeq(input), Fold(
			func(e error) []Item { t.Fatal(e); return nil },
			slices.Collect[Item],
		))

		assert.Equal(t, []Item{
			{Value: 1, Label: "first"},
			{Value: 2, Label: "second"},
			{Value: 3, Label: "third"},
		}, collected)
	})
}

func TestSequenceSeq_Failure(t *testing.T) {
	extractErr := func(result Result[iter.Seq[int]]) error {
		return F.Pipe1(result, Fold(
			F.Identity[error],
			func(_ iter.Seq[int]) error { t.Fatal("expected Left but got Right"); return nil },
		))
	}

	t.Run("short-circuits on first Left", func(t *testing.T) {
		testErr := errors.New("test error")
		input := slices.Values([]Result[int]{Of(1), Left[int](testErr), Of(3)})
		assert.Equal(t, testErr, extractErr(SequenceSeq(input)))
	})

	t.Run("returns first error when multiple Left values exist", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		input := slices.Values([]Result[int]{Of(1), Left[int](err1), Left[int](err2)})
		assert.Equal(t, err1, extractErr(SequenceSeq(input)))
	})

	t.Run("handles Left at the beginning", func(t *testing.T) {
		testErr := errors.New("first error")
		input := slices.Values([]Result[int]{Left[int](testErr), Of(2), Of(3)})
		assert.Equal(t, testErr, extractErr(SequenceSeq(input)))
	})

	t.Run("handles Left at the end", func(t *testing.T) {
		testErr := errors.New("last error")
		input := slices.Values([]Result[int]{Of(1), Of(2), Left[int](testErr)})
		assert.Equal(t, testErr, extractErr(SequenceSeq(input)))
	})
}

func TestSequenceSeq_Integration(t *testing.T) {
	t.Run("integrates with TraverseSeq", func(t *testing.T) {
		parse := func(s string) Result[int] {
			v, err := strconv.Atoi(s)
			return TryCatchError(v, err)
		}
		result := TraverseSeq(parse)(slices.Values([]string{"1", "2", "3"}))
		assert.True(t, IsRight(result))
	})

	t.Run("SequenceSeq is equivalent to TraverseSeq with Identity", func(t *testing.T) {
		mkInput := func() []Result[int] {
			return []Result[int]{Of(10), Of(20), Of(30)}
		}

		collected1 := F.Pipe1(SequenceSeq(slices.Values(mkInput())), Fold(
			func(e error) []int { t.Fatal(e); return nil },
			slices.Collect[int],
		))
		collected2 := F.Pipe1(TraverseSeq(F.Identity[Result[int]])(slices.Values(mkInput())), Fold(
			func(e error) []int { t.Fatal(e); return nil },
			slices.Collect[int],
		))

		assert.Equal(t, collected1, collected2)
	})
}
