package insideoutside_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/insideoutside"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func TestMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      rect2d.Format
		out     rect2d.Format
		want    rect2d.Fragments
		wantErr bool
	}{
		{
			name:    "In is not valid",
			in:      rect2d.Format{},
			out:     rect2d.Format{Width: 100, Height: 200},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Out is not valid",
			in:      rect2d.Format{Width: 10, Height: 20},
			out:     rect2d.Format{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Area in",
			in:   rect2d.Format{Width: 100, Height: 100},
			out:  rect2d.Format{Width: 300, Height: 400},
			want: rect2d.Fragments{
				rect2d.Fragment{
					Element: rect2d.Format{Width: 100, Height: 100},
					Layout:  rect2d.Layout{ByWidth: 3, ByHeight: 4},
				},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one) best variant", // [{4, 6}] or [{8, 3}]
			in:   rect2d.Format{Width: 100, Height: 50},
			out:  rect2d.Format{Width: 300, Height: 400},
			want: rect2d.Fragments{
				rect2d.Fragment{
					Element: rect2d.Format{Width: 100, Height: 50},
					Layout:  rect2d.Layout{ByWidth: 3, ByHeight: 8},
				},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one) best variant (invert)", // [{8, 3}] or [{4, 6}]
			in:   rect2d.Format{Width: 50, Height: 100},
			out:  rect2d.Format{Width: 300, Height: 400},
			want: rect2d.Fragments{
				rect2d.Fragment{
					Element: rect2d.Format{Width: 100, Height: 50},
					Layout:  rect2d.Layout{ByWidth: 3, ByHeight: 8},
				},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one and two) best variant", // [{4, 9}, {3, 1}] or [{13, 3}]
			in:   rect2d.Format{Width: 100, Height: 33},
			out:  rect2d.Format{Width: 300, Height: 450},
			want: rect2d.Fragments{
				rect2d.Fragment{
					Element: rect2d.Format{Width: 100, Height: 33},
					Layout:  rect2d.Layout{ByWidth: 3, ByHeight: 13},
				},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one and two) best variant (invert)", // [{13, 3}] or [{4, 9}, {3, 1}]
			in:   rect2d.Format{Width: 33, Height: 100},
			out:  rect2d.Format{Width: 300, Height: 450},
			want: rect2d.Fragments{
				rect2d.Fragment{
					Element: rect2d.Format{Width: 100, Height: 33},
					Layout:  rect2d.Layout{ByWidth: 3, ByHeight: 13},
				},
			},
			wantErr: false,
		},
		{
			name: "Rect in (two fragments)",
			in:   rect2d.Format{Width: 31, Height: 12},
			out:  rect2d.Format{Width: 273, Height: 1231},
			want: rect2d.Fragments{
				rect2d.Fragment{
					Element:  rect2d.Format{Width: 31, Height: 12},
					Layout:   rect2d.Layout{ByWidth: 8, ByHeight: 102},
					Position: enum.PositionTop,
				},
				rect2d.Fragment{
					Element:  rect2d.Format{Width: 12, Height: 31},
					Layout:   rect2d.Layout{ByWidth: 39, ByHeight: 2},
					Position: enum.PositionOnside,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := insideoutside.AlgoMax(tt.in, tt.out)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
