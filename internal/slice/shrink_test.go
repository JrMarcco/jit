package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_Shrink(t *testing.T) {
	testCases := []struct {
		name      string
		originCap int
		originLen int
		wantCap   int
	}{
		{
			name:      "空切片",
			originCap: 0,
			originLen: 0,
			wantCap:   0,
		}, {
			name:      "超大容量：当比例 >= 2 时缩容到 1.5 倍比例",
			originCap: 8192,
			originLen: 1024,
			wantCap:   1536,
		}, {
			name:      "大容量：当比例 >= 2 时缩容到原容量的 50%",
			originCap: 2048,
			originLen: 256,
			wantCap:   1024,
		}, {
			name:      "中容量：当比例 >= 2.5 时缩容到原容量的 62.5%",
			originCap: 1024,
			originLen: 256,
			wantCap:   640,
		}, {
			name:      "小容量：当比例 >= 3 时缩容到原容量的 50%",
			originCap: 128,
			originLen: 8,
			wantCap:   64,
		}, {
			name:      "下小容量：当比例 < 3 时，不缩容",
			originCap: 128,
			originLen: 64,
			wantCap:   128,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			slice := make([]int, 0, tc.originCap)

			for i := 0; i < tc.originLen; i++ {
				slice = append(slice, i)
			}

			res := Shrink(slice)
			assert.Equal(t, tc.wantCap, cap(res))
		})
	}
}
