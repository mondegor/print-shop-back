package rect3d

import (
	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s3"
)

type (
	// Format - формат прямоугольного 3d элемента.
	Format struct {
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
)

// String - возвращает текущий формат в виде строки {length}x{width}x{height}.
func (f Format) String() string {
	return s3.Size{f.Length, f.Width, f.Height}.String()
}

// IsValid - проверяет валиден ли текущий формат.
func (f Format) IsValid() bool {
	return f.Length > 0 && f.Width > 0 && f.Height > 0
}

// IsZero - проверяет что формат является нулевым.
func (f Format) IsZero() bool {
	return f.Length == 0 && f.Width == 0 && f.Height == 0
}

// Volume - возвращает объём формата.
func (f Format) Volume() float64 {
	return f.Length * f.Width * f.Height
}

// Add - возвращает сумму текущего формата с указанным.
func (f Format) Add(other Format) Format {
	return Format{
		Length: f.Length + other.Length,
		Width:  f.Width + other.Width,
		Height: f.Height + other.Height,
	}
}

// AddSame - возвращает текущий формат увеличенный на одинаковый отрезок.
func (f Format) AddSame(value float64) Format {
	return f.Add(Format{Length: value, Width: value, Height: value})
}

// Sub - возвращает разницу текущего формата и указанного.
// Если какая-то сторона указанного формата больше той же стороны текущего формата,
// то эта сторона будет приравнена 0.
func (f Format) Sub(other Format) Format {
	sub := func(first, second float64) float64 {
		if first > second {
			return first - second
		}

		return 0
	}

	return Format{
		Length: sub(f.Length, other.Length),
		Width:  sub(f.Width, other.Width),
		Height: sub(f.Height, other.Height),
	}
}

// SubSame - возвращает текущий формат уменьшенный на одинаковый отрезок.
func (f Format) SubSame(value float64) Format {
	return f.Sub(Format{Length: value, Width: value, Height: value})
}

// BottomFormat - возвращает формат дна прямоугольного 3d элемента.
func (f Format) BottomFormat() rect2d.Format {
	return rect2d.Format{
		Width:  f.Length,
		Height: f.Width,
	}
}

// Transform - возвращает формат преобразованный к другой единице измерения.
func (f Format) Transform(coefficient float64) Format {
	return Format{
		Length: f.Length * coefficient,
		Width:  f.Width * coefficient,
		Height: f.Height * coefficient,
	}
}

// Round - возвращает округлённый формат.
func (f Format) Round() Format {
	return Format{
		Length: mrlib.RoundFloat4(f.Length),
		Width:  mrlib.RoundFloat4(f.Width),
		Height: mrlib.RoundFloat4(f.Height),
	}
}
