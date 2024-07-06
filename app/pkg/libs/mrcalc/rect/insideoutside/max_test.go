package insideoutside_test

import (
	"testing"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/insideoutside"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      rect.Format
		out     rect.Format
		want    base.Fragments
		wantErr bool
	}{
		{
			name:    "In is not valid",
			in:      rect.Format{},
			out:     rect.Format{Width: 100, Height: 200},
			want:    base.Fragments{},
			wantErr: true,
		},
		{
			name:    "Out is not valid",
			in:      rect.Format{Width: 10, Height: 20},
			out:     rect.Format{},
			want:    base.Fragments{},
			wantErr: true,
		},
		{
			name: "Area in",
			in:   rect.Format{Width: 100, Height: 100},
			out:  rect.Format{Width: 300, Height: 400},
			want: base.Fragments{
				base.Fragment{ByWidth: 4, ByHeight: 3},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one) best variant", // [{4, 6}] or [{8, 3}]
			in:   rect.Format{Width: 100, Height: 50},
			out:  rect.Format{Width: 300, Height: 400},
			want: base.Fragments{
				base.Fragment{ByWidth: 8, ByHeight: 3},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one) best variant (invert)", // [{8, 3}] or [{4, 6}]
			in:   rect.Format{Width: 50, Height: 100},
			out:  rect.Format{Width: 300, Height: 400},
			want: base.Fragments{
				base.Fragment{ByWidth: 8, ByHeight: 3},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one and two) best variant", // [{4, 9}, {3, 1}] or [{13, 3}]
			in:   rect.Format{Width: 100, Height: 33},
			out:  rect.Format{Width: 300, Height: 450},
			want: base.Fragments{
				base.Fragment{ByWidth: 13, ByHeight: 3},
			},
			wantErr: false,
		},
		{
			name: "Rect in (one and two) best variant (invert)", // [{13, 3}] or [{4, 9}, {3, 1}]
			in:   rect.Format{Width: 33, Height: 100},
			out:  rect.Format{Width: 300, Height: 450},
			want: base.Fragments{
				base.Fragment{ByWidth: 13, ByHeight: 3},
			},
			wantErr: false,
		},
		{
			name: "Rect in best variant (invert)",
			in:   rect.Format{Width: 33, Height: 100},
			out:  rect.Format{Width: 300, Height: 450},
			want: base.Fragments{
				base.Fragment{ByWidth: 13, ByHeight: 3},
			},
			wantErr: false,
		},
		{
			name: "Rect in best variant (invert)",
			in:   rect.Format{Width: 33, Height: 100},
			out:  rect.Format{Width: 300, Height: 450},
			want: base.Fragments{
				base.Fragment{ByWidth: 13, ByHeight: 3},
			},
			wantErr: false,
		},
		{
			name: "Rect in (two fragments)",
			in:   rect.Format{Width: 31, Height: 12},
			out:  rect.Format{Width: 273, Height: 1231},
			want: base.Fragments{
				base.Fragment{ByWidth: 102, ByHeight: 8},
				base.Fragment{ByWidth: 2, ByHeight: 39},
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
