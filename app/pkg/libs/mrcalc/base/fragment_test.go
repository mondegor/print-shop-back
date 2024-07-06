package base_test

import (
	"testing"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"

	"github.com/stretchr/testify/assert"
)

func TestFragment_Total(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		byWidth  uint64
		byHeight uint64
		want     uint64
	}{
		{
			name:     "Zero value",
			byWidth:  0,
			byHeight: 0,
			want:     0,
		},
		{
			name:     "Zero width",
			byWidth:  0,
			byHeight: 15,
			want:     0,
		},
		{
			name:     "Zero height",
			byWidth:  10,
			byHeight: 0,
			want:     0,
		},
		{
			name:     "Positive width and height",
			byWidth:  10,
			byHeight: 15,
			want:     150,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &base.Fragment{
				ByWidth:  tt.byWidth,
				ByHeight: tt.byHeight,
			}

			got := f.Total()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFragment_Max(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		byWidth  uint64
		byHeight uint64
		want     uint64
	}{
		{
			name:     "Zero value",
			byWidth:  0,
			byHeight: 0,
			want:     0,
		},
		{
			name:     "Zero width",
			byWidth:  0,
			byHeight: 15,
			want:     15,
		},
		{
			name:     "Zero height",
			byWidth:  10,
			byHeight: 0,
			want:     10,
		},
		{
			name:     "Positive width and height",
			byWidth:  10,
			byHeight: 15,
			want:     15,
		},
		{
			name:     "Positive width and height (invert)",
			byWidth:  15,
			byHeight: 10,
			want:     15,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &base.Fragment{
				ByWidth:  tt.byWidth,
				ByHeight: tt.byHeight,
			}

			got := f.Max()
			assert.Equal(t, tt.want, got)
		})
	}
}
