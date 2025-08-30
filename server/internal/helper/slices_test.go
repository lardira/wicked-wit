package helper

import (
	"math/rand"
	"slices"
	"testing"
)

func TestRandomSubset(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		t.Parallel()

		emptySlice := []string{}
		got := RandomSubset(emptySlice, 1000)

		if got != nil {
			t.Fatalf("got %v want %v", got, nil)
		}
	})

	t.Run("1 item", func(t *testing.T) {
		t.Parallel()

		type TestStruct struct {
			A int
			B bool
			C complex128
		}

		var single TestStruct

		s := []TestStruct{single}
		got := RandomSubset(s, 1)
		want := []TestStruct{single}

		if !slices.Equal(got, want) {
			t.Fatalf("got %v want %v", got, nil)
		}
	})

	t.Run("many items", func(t *testing.T) {
		t.Parallel()

		min := 2
		n := min + rand.Intn(10)
		expLen := n / 2

		manyItemsInput := make([]int8, n)
		for i := range n {
			manyItemsInput[i] = int8(i + 1)
		}
		got := RandomSubset(manyItemsInput, expLen)

		if got == nil {
			t.Fatal("got must not be nill")
		}

		if len(got) != expLen {
			t.Fatalf("expected %v got len %v of %v", expLen, len(got), got)
		}
	})
}
