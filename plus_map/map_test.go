package plus_map

import (
	"testing"

	"github.com/JrMarcco/easy_kit/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	tcs := []struct {
		name string
		m    map[int]int
		want []int
	}{
		{
			name: "basic",
			m:    map[int]int{1: 1, 2: 2, 3: 3},
			want: []int{1, 2, 3},
		}, {
			name: "empty",
			m:    map[int]int{},
			want: []int{},
		}, {
			name: "nil",
			m:    nil,
			want: []int{},
		},
	}

	for _, tc := range tcs {
		res := Keys(tc.m)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func TestVals(t *testing.T) {
	tcs := []struct {
		name string
		m    map[int]int
		want []int
	}{
		{
			name: "basic",
			m:    map[int]int{1: 1, 2: 2, 3: 3},
			want: []int{1, 2, 3},
		}, {
			name: "empty",
			m:    map[int]int{},
			want: []int{},
		}, {
			name: "nil",
			m:    nil,
			want: []int{},
		},
	}

	for _, tc := range tcs {
		res := Vals(tc.m)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func TestKeysVals(t *testing.T) {
	tcs := []struct {
		name string
		m    map[int]int
		want []MapKV[int, int]
	}{
		{
			name: "basic",
			m:    map[int]int{1: 1, 2: 2, 3: 3},
			want: []MapKV[int, int]{{Key: 1, Val: 1}, {Key: 2, Val: 2}, {Key: 3, Val: 3}},
		}, {
			name: "empty",
			m:    map[int]int{},
			want: []MapKV[int, int]{},
		}, {
			name: "nil",
			m:    nil,
			want: []MapKV[int, int]{},
		},
	}

	for _, tc := range tcs {
		res := KeysVals(tc.m)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func TestToMap(t *testing.T) {
	tcs := []struct {
		name    string
		keys    []int
		vals    []int
		wantRes map[int]int
		wantErr error
	}{
		{
			name:    "basic",
			keys:    []int{1, 2, 3},
			vals:    []int{1, 2, 3},
			wantRes: map[int]int{1: 1, 2: 2, 3: 3},
			wantErr: nil,
		}, {
			name:    "nil keys",
			keys:    nil,
			vals:    []int{1, 2, 3},
			wantRes: nil,
			wantErr: errs.NilErr("keys or vals"),
		}, {
			name:    "nil vals",
			keys:    []int{1, 2, 3},
			vals:    nil,
			wantRes: nil,
			wantErr: errs.NilErr("keys or vals"),
		}, {
			name:    "different lengths",
			keys:    []int{1, 2, 3},
			vals:    []int{1, 2},
			wantRes: nil,
			wantErr: errs.InvalidKeyValLenErr(),
		}, {
			name:    "empty keys",
			keys:    []int{},
			vals:    []int{1, 2, 3},
			wantRes: map[int]int{},
			wantErr: nil,
		}, {
			name:    "empty vals",
			keys:    []int{1, 2, 3},
			vals:    []int{},
			wantRes: map[int]int{},
			wantErr: nil,
		},
	}

	for _, tc := range tcs {
		res, err := ToMap(tc.keys, tc.vals)
		assert.Equal(t, tc.wantErr, err)

		if err != nil {
			return
		}

		assert.Equal(t, tc.wantRes, res)
	}
}
