package s2

import (
	"fmt"
	"strings"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc"
)

type (
	// Size - двухсторонний размер.
	Size [2]float64
)

// String - возвращает текущий размер в виде строки {x}x{y}.
func (s Size) String() string {
	return mrcalc.FormatPositiveFloatToString(s[0]) + mrcalc.SizeSeparator +
		mrcalc.FormatPositiveFloatToString(s[1])
}

// ParseR2Size - возвращает результат парсинга строки вида '{first}x{second}' - два вещественных числа.
func ParseR2Size(str string) (size Size, err error) {
	items := strings.SplitN(str, mrcalc.SizeSeparator, 2)

	if len(items) != 2 {
		return Size{}, fmt.Errorf("size '%s' must be like {first}x{second} size", str)
	}

	size[0], err = mrcalc.ParseStringToPositiveFloat(items[0])
	if err != nil {
		return Size{}, err
	}

	size[1], err = mrcalc.ParseStringToPositiveFloat(items[1])
	if err != nil {
		return Size{}, err
	}

	return size, nil
}
