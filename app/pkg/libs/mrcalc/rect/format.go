package rect

import (
	"errors"
	"strconv"

	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
)

type (
	// Format - прямоугольный формат.
	Format struct {
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
)

// String - возвращает текущий формат в виде строки {width}x{height}.
func (f Format) String() string {
	return strconv.FormatFloat(mrlib.RoundFloat8(f.Width), 'f', -1, 64) + "x" +
		strconv.FormatFloat(mrlib.RoundFloat8(f.Height), 'f', -1, 64)
}

// IsValid - проверяет валиден ли текущий формат.
func (f *Format) IsValid() bool {
	return f.Width > 0 && f.Height > 0
}

// IsZero - проверяет что формат является нулевым.
func (f *Format) IsZero() bool {
	return f.Width == 0 && f.Height == 0
}

// Area - возвращает площадь формата.
func (f *Format) Area() float64 {
	return f.Width * f.Height
}

// Max - возвращает самую длинную сторону формата.
func (f *Format) Max() float64 {
	if f.Width > f.Height {
		return f.Width
	}

	return f.Height
}

// Cast - возвращает формат в таком виде, чтобы ширина формата
// была всегда больше или равна высоте формата.
func (f *Format) Cast() Format {
	if f.Width < f.Height {
		return f.Change()
	}

	return *f
}

// Change - возвращает изменение сторон формата: ширину и высоту местами
// (поворачивает формат на 90 градусов).
func (f *Format) Change() Format {
	return Format{
		Width:  f.Height,
		Height: f.Width,
	}
}

// Sum - возвращает сумму текущего формата с указанным.
func (f *Format) Sum(second Format) Format {
	return Format{
		Width:  f.Width + second.Width,
		Height: f.Height + second.Height,
	}
}

// Diff - возвращает разницу текущего формата и указанного.
// Если какая-то сторона указанного формата больше той же стороны текущего формата,
// то эта сторона будет приравнена 0.
func (f *Format) Diff(second Format) Format {
	res := *f

	if res.Width > second.Width {
		res.Width -= second.Width
	} else {
		res.Width = 0
	}

	if res.Height > second.Height {
		res.Height -= second.Height
	} else {
		res.Height = 0
	}

	return res
}

// DivBy - возвращает формат разделённый на указанную величину по широкой стороне.
func (f *Format) DivBy(divisor uint64) (Format, error) {
	if divisor == 0 {
		return Format{}, errors.New("divide by zero")
	}

	res := *f

	if res.Height >= res.Width {
		res.Height /= float64(divisor)
	} else {
		res.Width /= float64(divisor)
	}

	return res, nil
}

// Compare - возвращает результат сравнения формата с указанным.
func (f *Format) Compare(second Format) base.CompareType {
	first := f.Cast()
	second = second.Cast()

	if first.Width == second.Width && first.Height == second.Height {
		return base.CompareTypeEqual
	}

	if first.Width <= second.Width && first.Height <= second.Height {
		return base.CompareTypeFirstInside
	}

	if first.Width >= second.Width && first.Height >= second.Height {
		return base.CompareTypeSecondInside
	}

	return base.CompareTypeNotCompatible
}

// OrientationType - возвращает тип ориентации формата (книжный, альбомный).
func (f *Format) OrientationType() string {
	if f.Width > f.Height {
		return base.OrientationAlbum
	}

	return base.OrientationBook
}
