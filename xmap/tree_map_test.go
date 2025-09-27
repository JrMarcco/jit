package xmap

import (
	"testing"

	"github.com/JrMarcco/jit"
	"github.com/JrMarcco/jit/xslice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cmp() jit.Comparator[int] {
	return func(a, b int) int { return a - b }
}

func TestNewTreeMapWithMap(t *testing.T) {
	tcs := []struct {
		name     string
		cmp      jit.Comparator[int]
		m        map[int]string
		wantKeys []int
		wantVals []string
		wantErr  error
	}{
		{
			name:    "nil comparator",
			cmp:     nil,
			wantErr: ErrNilComparator,
		}, {
			name:     "empty map",
			cmp:      cmp(),
			m:        map[int]string{},
			wantKeys: []int{},
			wantVals: []string{},
			wantErr:  nil,
		}, {
			name: "map single value",
			cmp:  cmp(),
			m: map[int]string{
				1: "1",
			},
			wantKeys: []int{1},
			wantVals: []string{"1"},
			wantErr:  nil,
		}, {
			name: "map multiple values",
			cmp:  cmp(),
			m: map[int]string{
				1: "1",
				2: "2",
				3: "3",
			},
			wantKeys: []int{1, 2, 3},
			wantVals: []string{"1", "2", "3"},
			wantErr:  nil,
		}, {
			name: "map with disordered keys	",
			cmp:  cmp(),
			m: map[int]string{
				3: "3",
				1: "1",
				2: "2",
				6: "6",
				4: "4",
				5: "5",
			},
			wantKeys: []int{1, 2, 3, 4, 5, 6},
			wantVals: []string{"1", "2", "3", "4", "5", "6"},
			wantErr:  nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			treeMap, err := NewTreeMapWithMap(tc.cmp, tc.m)

			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			keys, vals := treeMap.KeyVals()
			assert.Equal(t, tc.wantKeys, keys)
			assert.Equal(t, tc.wantVals, vals)
		})
	}
}

func TestTreeMap_Del(t *testing.T) {

	type kv struct {
		key int
		val string
	}

	type delParam struct {
		key     int
		wantRes bool
	}

	tcs := []struct {
		name     string
		initData []kv
		delKeys  []delParam
		wantSize int
		wantVals []string
	}{
		{
			name: "del non-existent key",
			initData: []kv{
				{key: 1, val: "1"},
			},
			delKeys:  []delParam{{key: 2, wantRes: false}},
			wantSize: 1,
			wantVals: []string{"1"},
		}, {
			name: "del single key",
			initData: []kv{
				{key: 1, val: "1"},
				{key: 2, val: "2"},
				{key: 3, val: "3"},
			},
			delKeys:  []delParam{{key: 1, wantRes: true}},
			wantSize: 2,
			wantVals: []string{"2", "3"},
		}, {
			name: "del multiple keys",
			initData: []kv{
				{key: 1, val: "1"},
				{key: 2, val: "2"},
				{key: 3, val: "3"},
				{key: 4, val: "4"},
				{key: 5, val: "5"},
			},
			delKeys:  []delParam{{key: 2, wantRes: true}, {key: 5, wantRes: true}},
			wantSize: 3,
			wantVals: []string{"1", "3", "4"},
		}, {
			name: "del all keys",
			initData: []kv{
				{key: 1, val: "1"},
				{key: 2, val: "2"},
			},
			delKeys:  []delParam{{key: 1, wantRes: true}, {key: 2, wantRes: true}},
			wantSize: 0,
			wantVals: []string{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			m := xslice.ToMapWithVal(tc.initData, func(elem kv) (int, string) { return elem.key, elem.val })
			treeMap, err := NewTreeMapWithMap(cmp(), m)
			require.NoError(t, err)

			for _, p := range tc.delKeys {
				val, ok := treeMap.Del(p.key)
				assert.Equal(t, p.wantRes, ok)

				if ok {
					assert.Equal(t, val, m[p.key])
				}
			}

			assert.Equal(t, tc.wantVals, treeMap.Vals())

		})
	}
}

func TestTreeMap_Get(t *testing.T) {
	tcs := []struct {
		name    string
		m       map[int]string
		key     int
		wantRes bool
		wantVal string
	}{
		{
			name: "get non-existent key",
			m: map[int]string{
				1: "1",
			},
			key:     2,
			wantRes: false,
		}, {
			name: "get existent key",
			m: map[int]string{
				1: "1",
				2: "2",
				3: "3",
			},
			key:     1,
			wantRes: true,
			wantVal: "1",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			treeMap, err := NewTreeMapWithMap(cmp(), tc.m)
			require.NoError(t, err)

			val, ok := treeMap.Get(tc.key)
			assert.Equal(t, tc.wantRes, ok)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}
