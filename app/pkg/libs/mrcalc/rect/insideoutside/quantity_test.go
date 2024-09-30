package insideoutside_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/insideoutside"
)

func TestAlgoQuantity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      rect.Format
		out     rect.Format
		want    base.Fragment
		wantErr bool
	}{
		{
			name:    "test",
			in:      rect.Format{},
			out:     rect.Format{Width: 100, Height: 200},
			want:    base.Fragment{},
			wantErr: true,
		},
		{
			name:    "test",
			in:      rect.Format{Height: 10},
			out:     rect.Format{Width: 100, Height: 200},
			want:    base.Fragment{},
			wantErr: true,
		},
		{
			name:    "test",
			in:      rect.Format{Width: 10},
			out:     rect.Format{Width: 100, Height: 200},
			want:    base.Fragment{},
			wantErr: true,
		},
		{
			name:    "Correct value",
			in:      rect.Format{Width: 10, Height: 5},
			out:     rect.Format{Width: 100, Height: 200},
			want:    base.Fragment{ByWidth: 10, ByHeight: 40},
			wantErr: false,
		},
		{
			name: "test",
			in:   rect.Format{Width: 20, Height: 50},
			out:  rect.Format{Width: 100, Height: 200},
			want: base.Fragment{ByWidth: 5, ByHeight: 4},
		},
		{
			name: "test",
			in:   rect.Format{Width: 11, Height: 33},
			out:  rect.Format{Width: 100, Height: 200},
			want: base.Fragment{ByWidth: 9, ByHeight: 6},
		},
		{
			name: "test",
			in:   rect.Format{Width: 11, Height: 33},
			out:  rect.Format{},
			want: base.Fragment{},
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
