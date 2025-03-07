package rect3d

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s3"
)

// ParseFormat - возвращает результат парсинга строки вида '{length}x{width}x{height}' в Format коробки.
func ParseFormat(str string) (Format, error) {
	size, err := s3.ParseSize(str)
	if err != nil {
		return Format{}, err
	}

	return Format{
		Length: size[0],
		Width:  size[1],
		Height: size[2],
	}, nil
}
