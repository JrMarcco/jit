package list

import (
	"testing"

	easykit "github.com/JrMarcco/easy-kit"
	"github.com/stretchr/testify/assert"
)

var testCmp = func() easykit.Comparator[int] {
	return func(a, b int) int { return a - b }
}()

func TestNewSkipList(t *testing.T) {
	sl := NewSkipList[int](testCmp)

	assert.Equal(t, sl.size, 0)
	assert.Equal(t, sl.currLevel, 1)
	assert.Equal(t, sl.head, &skipListNode[int]{
		next: make([]*skipListNode[int], MaxLevel),
		span: make([]int, MaxLevel),
	})
}

func TestSkipListFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sl := SkipListFromSlice[int](testCmp, slice)

	assert.Equal(t, sl.size, 10)
	assert.Equal(t, sl.ToSlice(), slice)
}

func TestSkipList_Insert(t *testing.T) {
	tcs := []struct {
		name      string
		list      *SkipList[int]
		val       int
		wantSlice []int
		wantSize  int
	}{
		{
			name:      "basic",
			list:      NewSkipList[int](testCmp),
			val:       1,
			wantSlice: []int{1},
			wantSize:  1,
		}, {
			name:      "insert exists value",
			list:      SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			val:       2,
			wantSlice: []int{1, 2, 2, 3},
			wantSize:  4,
		}, {
			name: "insert to head",
			list: SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			val:  0,
			wantSlice: []int{
				0, 1, 2, 3,
			},
			wantSize: 4,
		}, {
			name:      "insert to tail",
			list:      SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			val:       4,
			wantSlice: []int{1, 2, 3, 4},
			wantSize:  4,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.list.Insert(tc.val)

			assert.Equal(t, tc.list.size, tc.wantSize)
			assert.Equal(t, tc.list.ToSlice(), tc.wantSlice)
		})
	}
}

func TestSkipList_Delete(t *testing.T) {
	tcs := []struct {
		name      string
		list      *SkipList[int]
		val       int
		wantSlice []int
		wantSize  int
		wantRes   bool
	}{
		{
			name:      "basic",
			list:      SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			val:       2,
			wantSlice: []int{1, 3},
			wantSize:  2,
			wantRes:   true,
		}, {
			name:      "delete non-exist value",
			list:      SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			val:       4,
			wantSlice: []int{1, 2, 3},
			wantSize:  3,
			wantRes:   false,
		}, {
			name:      "delete head",
			list:      SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			val:       1,
			wantSlice: []int{2, 3},
			wantSize:  2,
			wantRes:   true,
		}, {
			name:      "delete tail",
			list:      SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			val:       3,
			wantSlice: []int{1, 2},
			wantSize:  2,
			wantRes:   true,
		}, {
			name:      "delete from empty list",
			list:      NewSkipList[int](testCmp),
			val:       1,
			wantSlice: []int{},
			wantSize:  0,
			wantRes:   false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.list.Delete(tc.val)

			assert.Equal(t, tc.list.size, tc.wantSize)
			assert.Equal(t, tc.list.ToSlice(), tc.wantSlice)
			assert.Equal(t, res, tc.wantRes)
		})
	}
}

func TestSkipList_Exist(t *testing.T) {
	tcs := []struct {
		name    string
		list    *SkipList[int]
		target  int
		wantRes bool
	}{
		{
			name:    "basic",
			list:    SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			target:  2,
			wantRes: true,
		}, {
			name:    "not exist",
			list:    SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			target:  4,
			wantRes: false,
		}, {
			name:    "empty list",
			list:    NewSkipList[int](testCmp),
			target:  1,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			exists := tc.list.Exists(tc.target)
			assert.Equal(t, exists, tc.wantRes)
		})
	}
}

func TestSkipList_GetByIndex(t *testing.T) {
	tcs := []struct {
		name    string
		list    *SkipList[int]
		idx     int
		wantVal int
		wantRes bool
	}{
		{
			name:    "basic",
			list:    SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			idx:     1,
			wantVal: 2,
			wantRes: true,
		}, {
			name:    "not exist",
			list:    SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			idx:     4,
			wantVal: 0,
			wantRes: false,
		}, {
			name:    "empty list",
			list:    NewSkipList[int](testCmp),
			idx:     0,
			wantVal: 0,
			wantRes: false,
		}, {
			name:    "head",
			list:    SkipListFromSlice[int](testCmp, []int{1, 1, 1, 2, 3, 4, 4, 9, 9, 10}),
			idx:     0,
			wantVal: 1,
			wantRes: true,
		}, {
			name:    "tail",
			list:    SkipListFromSlice[int](testCmp, []int{1, 1, 1, 2, 3, 4, 4, 9, 9, 10}),
			idx:     9,
			wantVal: 10,
			wantRes: true,
		}, {
			name:    "out of range",
			list:    SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			idx:     6,
			wantRes: false,
		}, {
			name:    "negative index",
			list:    SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			idx:     -1,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := tc.list.GetByIndex(tc.idx)
			assert.Equal(t, ok, tc.wantRes)

			if ok {
				assert.Equal(t, val, tc.wantVal)
			}
		})
	}
}

func TestSkipList_Peek(t *testing.T) {
	tcs := []struct {
		name    string
		list    *SkipList[int]
		wantVal int
		wantRes bool
	}{
		{
			name:    "basic",
			list:    SkipListFromSlice[int](testCmp, []int{1, 2, 3}),
			wantVal: 1,
			wantRes: true,
		}, {
			name:    "empty list",
			list:    NewSkipList[int](testCmp),
			wantVal: 0,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := tc.list.Peek()
			assert.Equal(t, ok, tc.wantRes)

			if ok {
				assert.Equal(t, val, tc.wantVal)
			}
		})
	}
}
