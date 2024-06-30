package base_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
)

func TestParseDoubleSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		str     string
		want1   int64
		want2   int64
		wantErr bool
	}{
		{
			name:    "Empty string",
			str:     "",
			want1:   0,
			want2:   0,
			wantErr: true,
		},
		{
			name:    "One param",
			str:     "10",
			want1:   0,
			want2:   0,
			wantErr: true,
		},
		{
			name:    "Tree params",
			str:     "10x20x30",
			want1:   0,
			want2:   0,
			wantErr: true,
		},
		{
			name:    "Correct value",
			str:     "10x20",
			want1:   10,
			want2:   20,
			wantErr: false,
		},
		{
			name:    "Correct zero value",
			str:     "0x0",
			want1:   0,
			want2:   0,
			wantErr: false,
		},
		{
			name:    "Negative left value",
			str:     "-10x20",
			want1:   0,
			want2:   0,
			wantErr: true,
		},
		{
			name:    "Negative right value",
			str:     "10x-20",
			want1:   0,
			want2:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got1, got2, err := base.ParseDoubleSize(tt.str)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
		})
	}
}
