package insideoutside_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/insideoutside"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func TestAlgoQuantity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      rect2d.Format
		out     rect2d.Format
		want    rect2d.Layout
		wantErr bool
	}{
		{
			name:    "test",
			in:      rect2d.Format{},
			out:     rect2d.Format{Width: 100, Height: 200},
			want:    rect2d.Layout{},
			wantErr: true,
		},
		{
			name:    "test",
			in:      rect2d.Format{Height: 10},
			out:     rect2d.Format{Width: 100, Height: 200},
			want:    rect2d.Layout{},
			wantErr: true,
		},
		{
			name:    "test",
			in:      rect2d.Format{Width: 10},
			out:     rect2d.Format{Width: 100, Height: 200},
			want:    rect2d.Layout{},
			wantErr: true,
		},
		{
			name:    "Correct value",
			in:      rect2d.Format{Width: 10, Height: 5},
			out:     rect2d.Format{Width: 100, Height: 200},
			want:    rect2d.Layout{ByWidth: 10, ByHeight: 40},
			wantErr: false,
		},
		{
			name: "test",
			in:   rect2d.Format{Width: 20, Height: 50},
			out:  rect2d.Format{Width: 100, Height: 200},
			want: rect2d.Layout{ByWidth: 5, ByHeight: 4},
		},
		{
			name: "test",
			in:   rect2d.Format{Width: 11, Height: 33},
			out:  rect2d.Format{Width: 100, Height: 200},
			want: rect2d.Layout{ByWidth: 9, ByHeight: 6},
		},
		{
			name: "test",
			in:   rect2d.Format{Width: 11, Height: 33},
			out:  rect2d.Format{},
			want: rect2d.Layout{},
		},
		{
			name: "Rect in (3 x 3)",
			in:   rect2d.Format{Width: 0.1, Height: 0.1},
			out:  rect2d.Format{Width: 0.3, Height: 0.3},
			want: rect2d.Layout{ByWidth: 3, ByHeight: 3},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := insideoutside.AlgoQuantity(tt.in, tt.out)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
