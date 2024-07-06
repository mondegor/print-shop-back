package rect_test

import (
	"testing"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat_Cast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   rect.Format
	}{
		{
			name:   "width > height",
			width:  1000,
			height: 100,
			want:   rect.Format{Width: 1000, Height: 100},
		},
		{
			name:   "height > width",
			width:  100,
			height: 1000,
			want:   rect.Format{Width: 1000, Height: 100},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.Cast()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Change(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   rect.Format
	}{
		{
			name:   "width > height",
			width:  1000,
			height: 100,
			want:   rect.Format{Width: 100, Height: 1000},
		},
		{
			name:   "height > width",
			width:  100,
			height: 1000,
			want:   rect.Format{Width: 1000, Height: 100},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.Change()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Compare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		first  rect.Format
		second rect.Format
		want   base.CompareType
	}{
		{
			name:   "Equal",
			first:  rect.Format{Width: 1000, Height: 2000},
			second: rect.Format{Width: 1000, Height: 2000},
			want:   base.CompareTypeEqual,
		},
		{
			name:   "Equal 2",
			first:  rect.Format{Width: 2000, Height: 1000},
			second: rect.Format{Width: 1000, Height: 2000},
			want:   base.CompareTypeEqual,
		},
		{
			name:   "Equal 3",
			first:  rect.Format{Width: 1000, Height: 2000},
			second: rect.Format{Width: 2000, Height: 1000},
			want:   base.CompareTypeEqual,
		},
		{
			name:   "First inside",
			first:  rect.Format{Width: 1000, Height: 200},
			second: rect.Format{Width: 1000, Height: 2000},
			want:   base.CompareTypeFirstInside,
		},
		{
			name:   "Second inside",
			first:  rect.Format{Width: 1000, Height: 2000},
			second: rect.Format{Width: 1000, Height: 200},
			want:   base.CompareTypeSecondInside,
		},
		{
			name:   "Not compatible",
			first:  rect.Format{Width: 2000, Height: 100},
			second: rect.Format{Width: 1000, Height: 200},
			want:   base.CompareTypeNotCompatible,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.first.Width,
				Height: tt.first.Height,
			}

			got := f.Compare(tt.second)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Diff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		first  rect.Format
		second rect.Format
		want   rect.Format
	}{
		{
			name:   "first > second",
			first:  rect.Format{Width: 101, Height: 1010},
			second: rect.Format{Width: 52, Height: 520},
			want:   rect.Format{Width: 49, Height: 490},
		},
		{
			name:   "second > first",
			first:  rect.Format{Width: 101, Height: 1010},
			second: rect.Format{Width: 202, Height: 2020},
			want:   rect.Format{},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.first.Width,
				Height: tt.first.Height,
			}

			got := f.Diff(tt.second)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_DivBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		width   float64
		height  float64
		divisor uint64
		want    rect.Format
		wantErr bool
	}{
		{
			name:    "Divide by zero",
			width:   100,
			height:  200,
			divisor: 0,
			want:    rect.Format{},
			wantErr: true,
		},
		{
			name:    "Div height by 2",
			width:   50,
			height:  200,
			divisor: 2,
			want:    rect.Format{Width: 50, Height: 100},
			wantErr: false,
		},
		{
			name:    "Div width by 2",
			width:   200,
			height:  50,
			divisor: 2,
			want:    rect.Format{Width: 100, Height: 50},
			wantErr: false,
		},
		{
			name:    "Div width by 3",
			width:   200,
			height:  50,
			divisor: 3,
			want:    rect.Format{Width: 66.667, Height: 50},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got, err := f.DivBy(tt.divisor)

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

func TestFormat_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   bool
	}{
		{
			name:   "Zero value",
			width:  0,
			height: 0,
			want:   false,
		},
		{
			name:   "Only zero width",
			width:  0,
			height: 100,
			want:   false,
		},
		{
			name:   "Only zero height",
			width:  0,
			height: 200,
			want:   false,
		},
		{
			name:   "Valid value",
			width:  100,
			height: 200,
			want:   true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_IsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   bool
	}{
		{
			name:   "Zero value",
			width:  0,
			height: 0,
			want:   true,
		},
		{
			name:   "Only zero width",
			width:  0,
			height: 100,
			want:   false,
		},
		{
			name:   "Only zero height",
			width:  0,
			height: 200,
			want:   false,
		},
		{
			name:   "Not zero value",
			width:  100,
			height: 200,
			want:   false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.IsZero()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   string
	}{
		{
			name:   "{width}x{height}",
			width:  100,
			height: 200,
			want:   "100x200",
		},
		{
			name:   "{width}x{height} float",
			width:  100.1,
			height: 200.83,
			want:   "100.1x200.83",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Sum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		first  rect.Format
		second rect.Format
		want   rect.Format
	}{
		{
			name:   "first + second",
			first:  rect.Format{Width: 101, Height: 1010},
			second: rect.Format{Width: 102, Height: 1020},
			want:   rect.Format{Width: 203, Height: 2030},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.first.Width,
				Height: tt.first.Height,
			}

			got := f.Sum(tt.second)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Area(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   float64
	}{
		{
			name:   "Zero value",
			width:  0,
			height: 0,
			want:   0,
		},
		{
			name:   "Only zero width",
			width:  0,
			height: 100,
			want:   0,
		},
		{
			name:   "Only zero height",
			width:  0,
			height: 200,
			want:   0,
		},
		{
			name:   "Not zero value",
			width:  100,
			height: 200,
			want:   20000,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.Area()
			assert.InDelta(t, tt.want, got, measure.DeltaThousand)
		})
	}
}
