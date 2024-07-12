package rect

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
)

// ParseFormat - возвращает результат парсинга строки вида '{width}x{height}' в Format.
func ParseFormat(str string) (Format, error) {
	format, err := base.ParseDoubleSize(str)
	if err != nil {
		return Format{}, err
	}

	return Format{
		Width:  format[0],
		Height: format[1],
	}, nil
}
