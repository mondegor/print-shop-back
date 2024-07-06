package base

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseDoubleSize - возвращает результат парсинга строки вида '{first}x{second}' два вещественных числа.
func ParseDoubleSize(str string) (first, second float64, err error) {
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

func parseDoubleSizeItem(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("error parse item '%s': %w", str, err)
	}

	if value < 0 {
		return 0, fmt.Errorf("error parse item '%s': got negative value", str)
	}

	return value, nil
}
