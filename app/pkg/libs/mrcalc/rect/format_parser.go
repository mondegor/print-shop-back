package rect

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
)

// ParseFormat - возвращает результат парсинга строки вида '{width}x{height}' в Format.
func ParseFormat(str string) (Format, error) {
	width, height, err := base.ParseDoubleSize(str)
	if err != nil {
		return Format{}, err
	}

	return Format{
		Width:  width,
		Height: height,
	}, nil
}
