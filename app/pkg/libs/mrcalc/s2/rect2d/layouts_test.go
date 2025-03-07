package rect2d_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func TestLayouts_Count(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		f    rect2d.Layouts
		want uint64
	}{
		{
			name: "Empty",
			f:    rect2d.Layouts{},
			want: 0,
		},
		{
			name: "One zero",
			f: rect2d.Layouts{
				rect2d.Layout{},
			},
			want: 0,
		},
		{
			name: "One",
			f: rect2d.Layouts{
				rect2d.Layout{ByWidth: 10, ByHeight: 15},
			},
			want: 150,
		},
		{
			name: "Some",
			f: rect2d.Layouts{
				rect2d.Layout{ByWidth: 5, ByHeight: 3},
				rect2d.Layout{ByWidth: 11, ByHeight: 4},
			},
			want: 59,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.f.TotalQuantity()
			assert.Equal(t, tt.want, got)
		})
	}
}
