package base

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseDoubleSize - возвращает результат парсинга строки вида '{first}x{second}' два целых числа.
func ParseDoubleSize(str string) (first, second int64, err error) {
	items := strings.SplitN(str, "x", 2)

	if len(items) != 2 {
		return 0, 0, fmt.Errorf("size '%s' must be like {first}x{second}", str)
	}

	first, err = parseDoubleSizeItem(items[0])
	if err != nil {
		return 0, 0, err
	}

	second, err = parseDoubleSizeItem(items[1])
	if err != nil {
		return 0, 0, err
	}

	return first, second, nil
}

func parseDoubleSizeItem(str string) (int64, error) {
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parse item '%s': %w", str, err)
	}

	return int64(value), nil
}
