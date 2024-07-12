package parallelepiped

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
)

// ParseFormat - возвращает результат парсинга строки вида '{length}x{width}x{height}' в Format.
func ParseFormat(str string) (Format, error) {
	format, err := base.ParseTripleSize(str)
	if err != nil {
		return Format{}, err
	}

	return Format{
		Length: format[0],
		Width:  format[1],
		Height: format[2],
	}, nil
}
