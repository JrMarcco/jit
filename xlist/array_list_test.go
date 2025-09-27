package xlist

import (
	"fmt"
	"testing"

	"github.com/JrMarcco/jit/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestArrayList_Insert(t *testing.T) {
	tcs := []struct {
		name    string
		al      *ArrayList[int]
		index   int
		value   int
		wantRes []int
		wantErr error
		wantLen int
	}{
		{
			name:    "basic",
			al:      NewArrayList[int](4),
			index:   0,
			value:   1,
			wantLen: 1,
			wantRes: []int{1},
		}, {
			name:    "basic out of bounds",
			al:      NewArrayList[int](4),
			index:   1,
			value:   1,
			wantErr: errs.ErrIndexOutOfBounds(0, 1),
			wantLen: 0,
		}, {
			name:    "insert to middle of empty with size 4",
			al:      NewArrayList[int](4),
			index:   2,
			value:   3,
			wantErr: errs.ErrIndexOutOfBounds(0, 2),
			wantLen: 0,
		}, {
			name:    "insert to head",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			value:   0,
			wantRes: []int{0, 1, 2, 3, 4, 5},
			wantLen: 6,
		}, {
			name:    "insert to tail",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			value:   6,
			wantRes: []int{1, 2, 3, 4, 5, 6},
			wantLen: 6,
		}, {
			name:    "insert to middle",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			value:   3,
			wantRes: []int{1, 2, 3, 3, 4, 5},
			wantLen: 6,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.al.Insert(tc.index, tc.value)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, tc.al.Len(), tc.wantLen)

			if err == nil {
				assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
			}
		})
	}
}

func TestArrayList_Append(t *testing.T) {
	tcs := []struct {
		name    string
		al      *ArrayList[int]
		values  []int
		wantRes []int
		wantLen int
	}{
		{
			name:    "basic",
			al:      NewArrayList[int](4),
			values:  []int{1, 2, 3, 4, 5},
			wantRes: []int{1, 2, 3, 4, 5},
			wantLen: 5,
		}, {
			name:    "append to non-empty list",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			values:  []int{6, 7, 8, 9, 10},
			wantRes: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantLen: 10,
		}, {
			name:    "append single item",
			al:      NewArrayList[int](4),
			values:  []int{1},
			wantRes: []int{1},
			wantLen: 1,
		}, {
			name:    "append single item to non-empty list",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			values:  []int{6},
			wantRes: []int{1, 2, 3, 4, 5, 6},
			wantLen: 6,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.al.Append(tc.values...)
			assert.NoError(t, err)
			assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
			assert.Equal(t, tc.al.Len(), tc.wantLen)
		})
	}
}

func TestArrayList_Del(t *testing.T) {
	tcs := []struct {
		name    string
		al      *ArrayList[int]
		index   int
		wantRes []int
		wantErr error
		wantLen int
	}{
		{
			name:    "basic",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			wantRes: []int{1, 2, 4, 5},
			wantLen: 4,
		}, {
			name:    "out of bounds",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			wantErr: errs.ErrIndexOutOfBounds(5, 5),
			wantRes: []int{1, 2, 3, 4, 5},
			wantLen: 5,
		}, {
			name:    "delete from head",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			wantRes: []int{2, 3, 4, 5},
			wantLen: 4,
		}, {
			name:    "delete from tail",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   4,
			wantRes: []int{1, 2, 3, 4},
			wantLen: 4,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.al.Del(tc.index)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
		})
	}
}

func TestArrayList_Set(t *testing.T) {
	tcs := []struct {
		name    string
		al      *ArrayList[int]
		index   int
		value   int
		wantRes []int
		wantErr error
	}{
		{
			name:    "basic",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			value:   6,
			wantRes: []int{1, 2, 6, 4, 5},
		}, {
			name:    "out of bounds",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			value:   6,
			wantErr: errs.ErrIndexOutOfBounds(5, 5),
			wantRes: []int{1, 2, 3, 4, 5},
		}, {
			name:    "set from head",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			value:   6,
			wantRes: []int{6, 2, 3, 4, 5},
		}, {
			name:    "set from tail",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   4,
			value:   6,
			wantRes: []int{1, 2, 3, 4, 6},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.al.Set(tc.index, tc.value)

			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
		})
	}
}

func TestArrayList_Get(t *testing.T) {
	tcs := []struct {
		name    string
		al      *ArrayList[int]
		index   int
		wantRes int
		wantErr error
	}{
		{
			name:    "basic",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			wantRes: 3,
		}, {
			name:    "out of bounds",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			wantErr: errs.ErrIndexOutOfBounds(5, 5),
		}, {
			name:    "empty list",
			al:      NewArrayList[int](4),
			index:   0,
			wantErr: errs.ErrIndexOutOfBounds(0, 0),
		}, {
			name:    "get from head",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			wantRes: 1,
		}, {
			name:    "get from tail",
			al:      ArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   4,
			wantRes: 5,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.al.Get(tc.index)
			assert.Equal(t, err, tc.wantErr)

			if err == nil {
				assert.Equal(t, res, tc.wantRes)
			}
		})
	}
}

func ExampleArrayList_Iter() {
	al := ArrayListOf([]int{1, 2, 3, 4, 5})
	_ = al.Iter(func(idx int, val int) error {
		fmt.Printf("%d: %d\n", idx, 2*val)
		return nil
	})
	// Output:
	// 0: 2
	// 1: 4
	// 2: 6
	// 3: 8
	// 4: 10
}
