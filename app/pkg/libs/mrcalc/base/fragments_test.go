package base_test

import (
	"testing"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"

	"github.com/stretchr/testify/assert"
)

func TestFragments_Total(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		f    base.Fragments
		want uint64
	}{
		{
			name: "Empty",
			f:    base.Fragments{},
			want: 0,
		},
		{
			name: "One zero",
			f: base.Fragments{
				base.Fragment{},
			},
			want: 0,
		},
		{
			name: "One",
			f: base.Fragments{
				base.Fragment{ByWidth: 10, ByHeight: 15},
			},
			want: 150,
		},
		{
			name: "Some",
			f: base.Fragments{
				base.Fragment{ByWidth: 5, ByHeight: 3},
				base.Fragment{ByWidth: 11, ByHeight: 4},
			},
			want: 59,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.f.Total()
			assert.Equal(t, tt.want, got)
		})
	}
}
