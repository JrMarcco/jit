package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tcs := []struct {
		name string
		src  []int
		fn   func(idx int, src int) int
		want []int
	}{
		{
			name: "basic",
			src:  []int{1, 2, 3},
			fn:   func(idx int, src int) int { return src * 2 },
			want: []int{2, 4, 6},
		}, {
			name: "empty",
			src:  []int{},
			fn:   func(idx int, src int) int { return src * 2 },
			want: []int{},
		}, {
			name: "nil",
			src:  nil,
			fn:   func(idx int, src int) int { return src * 2 },
			want: nil,
		},
	}

	for _, tc := range tcs {
		res := Map(tc.src, tc.fn)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func TestFilterMap(t *testing.T) {
	tcs := []struct {
		name string
		src  []int
		fn   func(idx int, src int) (int, bool)
		want []int
	}{
		{
			name: "basic",
			src:  []int{1, 2, 3, 4, 5},
			fn:   func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 },
			want: []int{2, 6, 10},
		}, {
			name: "empty",
			src:  []int{},
			fn:   func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 },
			want: []int{},
		}, {
			name: "nil",
			src:  nil,
			fn:   func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 },
			want: nil,
		}, {
			name: "no match",
			src:  []int{1, 2, 3, 4, 5},
			fn:   func(idx int, src int) (int, bool) { return src * 2, src > 10 },
			want: []int{},
		},
	}

	for _, tc := range tcs {
		res := FilterMap(tc.src, tc.fn)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func ExampleMap() {
	src := []int{1, 2, 3, 4, 5}
	res := Map(src, func(idx int, src int) int { return src * 2 })
	fmt.Println(res)
	// Output: [2 4 6 8 10]
}

func ExampleFilterMap() {
	src := []int{1, 2, 3, 4, 5}
	res := FilterMap(src, func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 })
	fmt.Println(res)
	// Output: [2 6 10]
}
