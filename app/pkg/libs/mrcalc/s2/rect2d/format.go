package rect2d

import (
	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2"
)

type (
	// Format - формат прямоугольного 2d элемента.
	Format struct {
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
)

// String - возвращает текущий формат в виде строки {width}x{height}.
func (f Format) String() string {
	return s2.Size{f.Width, f.Height}.String()
}

// IsValid - проверяет валиден ли текущий формат.
func (f Format) IsValid() bool {
	return f.Width > 0 && f.Height > 0
}

// IsZero - проверяет что формат является нулевым.
func (f Format) IsZero() bool {
	return f.Width == 0 && f.Height == 0
}

// Max - возвращает самую длинную сторону формата.
func (f Format) Max() float64 {
	if f.Width > f.Height {
		return f.Width
	}

	return f.Height
}

// Area - возвращает площадь формата.
func (f Format) Area() float64 {
	if !f.IsValid() {
		return 0
	}

	return f.Width * f.Height
}

// Rotate90 - возвращает новый формат меняя у исходного формата ширину и высоту местами
// (поворачивает формат на 90 градусов).
func (f Format) Rotate90() Format {
	return Format{
		Width:  f.Height,
		Height: f.Width,
	}
}

// Cast - возвращает формат в таком виде, чтобы ширина формата
// была всегда больше или равна высоте формата (альбомная ориентация).
func (f Format) Cast() Format {
	if f.Width > f.Height {
		return f
	}

	return f.Rotate90()
}

// Add - возвращает сумму текущего формата с указанным.
func (f Format) Add(other Format) Format {
	return Format{
		Width:  f.Width + other.Width,
		Height: f.Height + other.Height,
	}
}

// AddSame - возвращает текущий формат увеличенный на одинаковый отрезок.
func (f Format) AddSame(value float64) Format {
	return f.Add(Format{Width: value, Height: value})
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
		Width:  sub(f.Width, other.Width),
		Height: sub(f.Height, other.Height),
	}
}

// SubSame - возвращает текущий формат уменьшенный на одинаковый отрезок.
func (f Format) SubSame(value float64) Format {
	return f.Sub(Format{Width: value, Height: value})
}

// Compare - возвращает результат сравнения формата с указанным.
func (f Format) Compare(other Format) enum.CompareType {
	first := f.Cast()
	other = other.Cast()

	if first.Width == other.Width && first.Height == other.Height {
		return enum.CompareTypeEqual
	}

	if first.Width <= other.Width && first.Height <= other.Height {
		return enum.CompareTypeFirstInside
	}

	if first.Width >= other.Width && first.Height >= other.Height {
		return enum.CompareTypeSecondInside
	}

	return enum.CompareTypeNotCompatible
}

// OrientationType - возвращает тип ориентации формата (книжный, альбомный).
func (f Format) OrientationType() enum.Orientation {
	if f.Width > f.Height {
		return enum.FormatOrientationAlbum
	}

	return enum.FormatOrientationBook
}

// Transform - возвращает формат преобразованный к другой единице измерения.
func (f Format) Transform(coefficient float64) Format {
	return Format{
		Width:  f.Width * coefficient,
		Height: f.Height * coefficient,
	}
}

// Round - возвращает округлённый формат.
func (f Format) Round() Format {
	return Format{
		Width:  mrlib.RoundFloat4(f.Width),
		Height: mrlib.RoundFloat4(f.Height),
	}
}
