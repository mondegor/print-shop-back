package base

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseDoubleSize - возвращает результат парсинга строки вида '{first}x{second}' два вещественных числа.
func ParseDoubleSize(str string) (format [2]float64, err error) {
	items := strings.SplitN(str, "x", 2)

	if len(items) != 2 {
		return [2]float64{}, fmt.Errorf("size '%s' must be like {first}x{second}", str)
	}

	format[0], err = parseSizeItem(items[0])
	if err != nil {
		return [2]float64{}, err
	}

	format[1], err = parseSizeItem(items[1])
	if err != nil {
		return [2]float64{}, err
	}

	return format, nil
}

// ParseTripleSize - возвращает результат парсинга строки вида '{first}x{second}x{third}' три вещественных числа.
func ParseTripleSize(str string) (format [3]float64, err error) {
	items := strings.SplitN(str, "x", 3)

	if len(items) != 3 {
		return [3]float64{}, fmt.Errorf("size '%s' must be like {first}x{second}x{third}", str)
	}

	format[0], err = parseSizeItem(items[0])
	if err != nil {
		return [3]float64{}, err
	}

	format[1], err = parseSizeItem(items[1])
	if err != nil {
		return [3]float64{}, err
	}

	format[2], err = parseSizeItem(items[2])
	if err != nil {
		return [3]float64{}, err
	}

	return format, nil
}

func parseSizeItem(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("error parse item '%s': %w", str, err)
	}

	if value < 0 {
		return 0, fmt.Errorf("error parse item '%s': got negative value", str)
	}

	return value, nil
}
