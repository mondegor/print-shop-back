package rect2d_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func TestDivideIn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		width   float64
		height  float64
		divisor uint64
		want    rect2d.Format
		wantErr bool
	}{
		{
			name:    "Divide by zero",
			width:   100,
			height:  200,
			divisor: 0,
			want:    rect2d.Format{},
			wantErr: true,
		},
		{
			name:    "Div height by 2",
			width:   50,
			height:  200,
			divisor: 2,
			want:    rect2d.Format{Width: 50, Height: 100},
			wantErr: false,
		},
		{
			name:    "Div width by 2",
			width:   200,
			height:  50,
			divisor: 2,
			want:    rect2d.Format{Width: 100, Height: 50},
			wantErr: false,
		},
		{
			name:    "Div width by 3",
			width:   200,
			height:  50,
			divisor: 3,
			want:    rect2d.Format{Width: 66.667, Height: 50},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := rect2d.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got, err := rect2d.DivideIn(f, tt.divisor)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.InDelta(t, tt.want.Width, got.Width, measure.DeltaThousand)
			assert.InDelta(t, tt.want.Height, got.Height, measure.DeltaThousand)
		})
	}
}
