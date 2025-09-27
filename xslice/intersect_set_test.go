package xslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersectSetFunc(t *testing.T) {
	tcs := []struct {
		name    string
		src     []int
		dst     []int
		wantRes []int
	}{
		{
			name:    "empty slices",
			src:     []int{},
			dst:     []int{},
			wantRes: []int{},
		}, {
			name:    "no intersection",
			src:     []int{1, 2, 3},
			dst:     []int{4, 5, 6},
			wantRes: []int{},
		}, {
			name:    "intersection",
			src:     []int{1, 2, 3, 4, 5},
			dst:     []int{4, 5, 6, 7, 8},
			wantRes: []int{4, 5},
		}, {
			name:    "intersection with duplicates",
			src:     []int{1, 2, 3, 4, 5, 4, 5},
			dst:     []int{4, 5, 6, 7, 8, 4, 5},
			wantRes: []int{4, 5},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := IntersectSetFunc(tc.src, tc.dst, func(src, dst int) bool { return src == dst })
			assert.ElementsMatch(t, tc.wantRes, res)
		})
	}
}

func TestIntersectSet(t *testing.T) {
	tcs := []struct {
		name    string
		src     []int
		dst     []int
		wantRes []int
	}{
		{
			name:    "empty slices",
			src:     []int{},
			dst:     []int{},
			wantRes: []int{},
		}, {
			name:    "no intersection",
			src:     []int{1, 2, 3},
			dst:     []int{4, 5, 6},
			wantRes: []int{},
		}, {
			name:    "intersection",
			src:     []int{1, 2, 3, 4, 5},
			dst:     []int{4, 5, 6, 7, 8},
			wantRes: []int{4, 5},
		}, {
			name:    "intersection with duplicates",
			src:     []int{1, 2, 3, 4, 5, 4, 5},
			dst:     []int{4, 5, 6, 7, 8, 4, 5},
			wantRes: []int{4, 5},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := IntersectSet(tc.src, tc.dst)
			assert.ElementsMatch(t, tc.wantRes, res)
		})
	}
}

func ExampleIntersectSetFunc() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{4, 5, 6, 7, 8}
	res := IntersectSetFunc(src, dst, func(src, dst int) bool { return src == dst })
	fmt.Println(res)
	// Output: [4 5]
}

func ExampleIntersectSet() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{4, 5, 6, 7, 8}
	res := IntersectSet(src, dst)
	fmt.Println(res)
	// Output: [4 5]
}
