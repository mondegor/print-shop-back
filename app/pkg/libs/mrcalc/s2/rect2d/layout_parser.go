package rect2d

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2"
)

// ParseLayout - возвращает результат парсинга строки вида '{byWidth}x{byHeight}' в Layout.
func ParseLayout(str string) (Layout, error) {
	by, err := s2.ParseR2Size(str)
	if err != nil {
		return Layout{}, err
	}

	return Layout{
		ByWidth:  uint64(by[0]),
		ByHeight: uint64(by[1]),
	}, nil
}
