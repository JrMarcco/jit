package map_ext

import (
	"testing"

	"github.com/JrMarcco/easy_kit"
	"github.com/stretchr/testify/assert"
)

func cmp() easy_kit.Comparator[int] {
	return func(a, b int) int { return a - b }
}

func TestNewTreeMapWithMap(t *testing.T) {
	tcs := []struct {
		name     string
		cmp      easy_kit.Comparator[int]
		m        map[int]string
		wantKeys []int
		wantVals []string
		wantErr  error
	}{
		{
			name:    "nil comparator",
			cmp:     nil,
			wantErr: errNilComparator,
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

			keys, vals := treeMap.Kvs()
			assert.Equal(t, tc.wantKeys, keys)
			assert.Equal(t, tc.wantVals, vals)
		})
	}
}
