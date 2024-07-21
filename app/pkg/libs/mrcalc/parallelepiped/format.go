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
func (f Format) IsValid() bool {
	return f.Length > 0 && f.Width > 0 && f.Height > 0
}

// BottomFormat - возвращает формат дна параллелепипеда.
func (f Format) BottomFormat() rect.Format {
	return rect.Format{
		Width:  f.Length,
		Height: f.Width,
	}
}

// IsZero - проверяет что формат является нулевым.
func (f Format) IsZero() bool {
	return f.Length == 0 && f.Width == 0 && f.Height == 0
}

// Volume - возвращает объём формата.
func (f Format) Volume() float64 {
	return f.Length * f.Width * f.Height
}

// Diff - возвращает разницу текущего формата и указанного.
// Если какая-то сторона указанного формата больше той же стороны текущего формата,
// то эта сторона будет приравнена 0.
func (f Format) Diff(second Format) Format {
	diff := func(first, second float64) float64 {
		if first > second {
			return first - second
		}

		return 0
	}

	return Format{
		Length: diff(f.Length, second.Length),
		Width:  diff(f.Width, second.Width),
		Height: diff(f.Height, second.Height),
	}
}
