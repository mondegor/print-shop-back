package rect2d

import (
	"errors"
)

// DivideIn - возвращает новый формат разделённый на указанную величину по самой длинной стороне.
func DivideIn(format Format, divisor uint64) (Format, error) {
	if divisor == 0 {
		return Format{}, errors.New("divide by zero")
	}

	if format.Width > format.Height {
		return Format{
			Width:  format.Width / float64(divisor),
			Height: format.Height,
		}, nil
	}

	return Format{
		Width:  format.Width,
		Height: format.Height / float64(divisor),
	}, nil
}

// DivideInHalf - возвращает новый формат разделённый пополам по самой длинной стороне.
func DivideInHalf(format Format) (Format, error) {
	return DivideIn(format, 2)
}
