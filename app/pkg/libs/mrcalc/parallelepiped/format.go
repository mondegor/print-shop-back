package parallelepiped

import (
	"strconv"

	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

type (
	// Format - формат параллелепипеда.
	Format struct {
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
)

// String - возвращает текущий формат в виде строки {length}x{width}x{height}.
func (f Format) String() string {
	return strconv.FormatFloat(mrlib.RoundFloat8(f.Length), 'f', -1, 64) + "x" +
		strconv.FormatFloat(mrlib.RoundFloat8(f.Width), 'f', -1, 64) + "x" +
		strconv.FormatFloat(mrlib.RoundFloat8(f.Height), 'f', -1, 64)
}

// IsValid - проверяет валиден ли текущий формат.
func (f *Format) IsValid() bool {
	return f.Length > 0 && f.Width > 0 && f.Height > 0
}

// BottomFormat - возвращает формат дна параллелепипеда.
func (f *Format) BottomFormat() rect.Format {
	return rect.Format{
		Width:  f.Length,
		Height: f.Width,
	}
}

// IsZero - проверяет что формат является нулевым.
func (f *Format) IsZero() bool {
	return f.Length == 0 && f.Width == 0 && f.Height == 0
}

// Volume - возвращает объём формата.
func (f *Format) Volume() float64 {
	return f.Length * f.Width * f.Height
}
