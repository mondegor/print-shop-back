package mrcalc

import (
	"fmt"
	"strconv"

	"github.com/mondegor/go-webcore/mrlib"
)

const (
	SizeSeparator = "x" // SizeSeparator - разделитель для представления размера в виде строки
)

// ParseStringToPositiveFloat - comment method.
func ParseStringToPositiveFloat(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("error parse float value '%s': %w", str, err)
	}

	if value < 0 {
		return 0, fmt.Errorf("error parse float value '%s': got negative value", str)
	}

	return value, nil
}

// FormatPositiveFloatToString - comment method.
func FormatPositiveFloatToString(value float64) string {
	if value < 0 {
		return "!BADVALUE"
	}

	return strconv.FormatFloat(mrlib.RoundFloat8(value), 'f', -1, 64)
}
