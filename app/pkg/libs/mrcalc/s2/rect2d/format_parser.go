package rect2d

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2"
)

// ParseFormat - возвращает результат парсинга строки вида '{width}x{height}' в прямоугольный 2d Format.
func ParseFormat(str string) (Format, error) {
	size, err := s2.ParseR2Size(str)
	if err != nil {
		return Format{}, err
	}

	return Format{
		Width:  size[0],
		Height: size[1],
	}, nil
}
