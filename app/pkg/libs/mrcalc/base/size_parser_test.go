package base_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
)

func TestParseDoubleSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		str     string
		want    [2]float64
		wantErr bool
	}{
		{
			name:    "Empty string",
			str:     "",
			want:    [2]float64{},
			wantErr: true,
		},
		{
			name:    "One param",
			str:     "10",
			want:    [2]float64{},
			wantErr: true,
		},
		{
			name:    "Tree params",
			str:     "10x20x30",
			want:    [2]float64{},
			wantErr: true,
		},
		{
			name:    "Correct value",
			str:     "10x20",
			want:    [2]float64{10, 20},
			wantErr: false,
		},
		{
			name:    "Correct zero value",
			str:     "0x0",
			want:    [2]float64{},
			wantErr: false,
		},
		{
			name:    "Negative left value",
			str:     "-10x20",
			want:    [2]float64{},
			wantErr: true,
		},
		{
			name:    "Negative right value",
			str:     "10x-20",
			want:    [2]float64{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := base.ParseDoubleSize(tt.str)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.InDelta(t, tt.want[0], got[0], measure.DeltaThousand)
			assert.InDelta(t, tt.want[1], got[1], measure.DeltaThousand)
		})
	}
}

func TestParseTripleSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		str     string
		want    [3]float64
		wantErr bool
	}{
		{
			name:    "Empty string",
			str:     "",
			want:    [3]float64{},
			wantErr: true,
		},
		{
			name:    "One param",
			str:     "10",
			want:    [3]float64{},
			wantErr: true,
		},
		{
			name:    "Tow params",
			str:     "10x20",
			want:    [3]float64{},
			wantErr: true,
		},
		{
			name:    "Correct value",
			str:     "10x20x30",
			want:    [3]float64{10, 20, 30},
			wantErr: false,
		},
		{
			name:    "Correct zero value",
			str:     "0x0x0",
			want:    [3]float64{},
			wantErr: false,
		},
		{
			name:    "Negative left value",
			str:     "-10x20x30",
			want:    [3]float64{},
			wantErr: true,
		},
		{
			name:    "Negative middle value",
			str:     "10x-20x30",
			want:    [3]float64{},
			wantErr: true,
		},
		{
			name:    "Negative right value",
			str:     "10x20x-30",
			want:    [3]float64{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := base.ParseTripleSize(tt.str)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.InDelta(t, tt.want[0], got[0], measure.DeltaThousand)
			assert.InDelta(t, tt.want[1], got[1], measure.DeltaThousand)
			assert.InDelta(t, tt.want[2], got[2], measure.DeltaThousand)
		})
	}
}
