package s3

import (
	"fmt"
	"strings"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc"
)

type (
	// Size - трёхсторонний размер.
	Size [3]float64
)

// String - возвращает текущий размер в виде строки {x}x{y}x{z}.
func (s Size) String() string {
	return mrcalc.FormatPositiveFloatToString(s[0]) + mrcalc.SizeSeparator +
		mrcalc.FormatPositiveFloatToString(s[1]) + mrcalc.SizeSeparator +
		mrcalc.FormatPositiveFloatToString(s[2])
}

// ParseSize - возвращает результат парсинга строки вида '{first}x{second}x{third}' - три вещественных числа.
func ParseSize(str string) (size Size, err error) {
	items := strings.SplitN(str, mrcalc.SizeSeparator, 3)

	if len(items) != 3 {
		return Size{}, fmt.Errorf("size '%s' must be like {first}x{second}x{third}", str)
	}

	for i := 0; i < 3; i++ {
		size[i], err = mrcalc.ParseStringToPositiveFloat(items[i])
		if err != nil {
			return Size{}, err
		}
	}

	return size, nil
}
