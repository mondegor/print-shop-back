package enum

import (
	"encoding/json"

	"github.com/mondegor/go-core/errors"
)

// Расположения на прямоугольном формате.
const (
	PositionTop    Position = iota + 1 // сверху
	PositionOnside                     // сбоку -
	PositionBottom                     // снизу
)

const (
	enumNamePosition = "Position"
)

type (
	// Position - вид расположения элементов на листе прямоугольного формата.
	Position uint8
)

//nolint:gochecknoglobals
var (
	positionName = map[Position]string{
		PositionTop:    "TOP",
		PositionOnside: "ONSIDE",
		PositionBottom: "BOTTOM",
	}

	positionValue = map[string]Position{
		"TOP":    PositionTop,
		"ONSIDE": PositionOnside,
		"BOTTOM": PositionBottom,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *Position) ParseAndSet(value string) error {
	if parsedValue, ok := positionValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return errors.ErrInternalKeyNotFoundInSource.New(
		"key", value,
		"source", enumNamePosition,
	)
}

// String - возвращает значение в виде строки.
func (e Position) String() string {
	return positionName[e]
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e Position) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *Position) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}
