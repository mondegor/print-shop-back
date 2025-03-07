package rect2d_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

func TestFormat_Cast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   rect2d.Format
	}{
		{
			name:   "width > height",
			width:  1000,
			height: 100,
			want:   rect2d.Format{Width: 1000, Height: 100},
		},
		{
			name:   "height > width",
			width:  100,
			height: 1000,
			want:   rect2d.Format{Width: 1000, Height: 100},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect2d.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.Cast()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Rotate90(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  float64
		height float64
		want   rect2d.Format
	}{
		{
			name:   "width > height",
			width:  1000,
			height: 100,
			want:   rect2d.Format{Width: 100, Height: 1000},
		},
		{
			name:   "height > width",
			width:  100,
			height: 1000,
			want:   rect2d.Format{Width: 1000, Height: 100},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect2d.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.Rotate90()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Compare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		first  rect2d.Format
		second rect2d.Format
		want   enum.CompareType
	}{
		{
			name:   "Equal",
			first:  rect2d.Format{Width: 1000, Height: 2000},
			second: rect2d.Format{Width: 1000, Height: 2000},
			want:   enum.CompareTypeEqual,
		},
		{
			name:   "Equal 2",
			first:  rect2d.Format{Width: 2000, Height: 1000},
			second: rect2d.Format{Width: 1000, Height: 2000},
			want:   enum.CompareTypeEqual,
		},
		{
			name:   "Equal 3",
			first:  rect2d.Format{Width: 1000, Height: 2000},
			second: rect2d.Format{Width: 2000, Height: 1000},
			want:   enum.CompareTypeEqual,
		},
		{
			name:   "First inside",
			first:  rect2d.Format{Width: 1000, Height: 200},
			second: rect2d.Format{Width: 1000, Height: 2000},
			want:   enum.CompareTypeFirstInside,
		},
		{
			name:   "Second inside",
			first:  rect2d.Format{Width: 1000, Height: 2000},
			second: rect2d.Format{Width: 1000, Height: 200},
			want:   enum.CompareTypeSecondInside,
		},
		{
			name:   "Not compatible",
			first:  rect2d.Format{Width: 2000, Height: 100},
			second: rect2d.Format{Width: 1000, Height: 200},
			want:   enum.CompareTypeNotCompatible,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect2d.Format{
				Width:  tt.first.Width,
				Height: tt.first.Height,
			}

			got := f.Compare(tt.second)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Sub(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		first  rect2d.Format
		second rect2d.Format
		want   rect2d.Format
	}{
		{
			name:   "first > second",
			first:  rect2d.Format{Width: 101, Height: 1010},
			second: rect2d.Format{Width: 52, Height: 520},
			want:   rect2d.Format{Width: 49, Height: 490},
		},
		{
			name:   "second > first",
			first:  rect2d.Format{Width: 101, Height: 1010},
			second: rect2d.Format{Width: 202, Height: 2020},
			want:   rect2d.Format{},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect2d.Format{
				Width:  tt.first.Width,
				Height: tt.first.Height,
			}

			got := f.Sub(tt.second)
			assert.Equal(t, tt.want, got)
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

			f := &rect2d.Format{
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

			f := &rect2d.Format{
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

			f := &rect2d.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		first  rect2d.Format
		second rect2d.Format
		want   rect2d.Format
	}{
		{
			name:   "first + second",
			first:  rect2d.Format{Width: 101, Height: 1010},
			second: rect2d.Format{Width: 102, Height: 1020},
			want:   rect2d.Format{Width: 203, Height: 2030},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &rect2d.Format{
				Width:  tt.first.Width,
				Height: tt.first.Height,
			}

			got := f.Add(tt.second)
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

			f := &rect2d.Format{
				Width:  tt.width,
				Height: tt.height,
			}

			got := f.Area()
			assert.InDelta(t, tt.want, got, measure.DeltaThousand)
		})
	}
}
