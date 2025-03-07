package rect2d_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func TestLayout_Count(t *testing.T) {
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

			f := &rect2d.Layout{
				ByWidth:  tt.byWidth,
				ByHeight: tt.byHeight,
			}

			got := f.Quantity()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLayout_Min(t *testing.T) {
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
			want:     10,
		},
		{
			name:     "Positive width and height (invert)",
			byWidth:  15,
			byHeight: 10,
			want:     10,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect2d.Layout{
				ByWidth:  tt.byWidth,
				ByHeight: tt.byHeight,
			}

			got := f.Min()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLayout_Max(t *testing.T) {
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

			f := &rect2d.Layout{
				ByWidth:  tt.byWidth,
				ByHeight: tt.byHeight,
			}

			got := f.Max()
			assert.Equal(t, tt.want, got)
		})
	}
}
