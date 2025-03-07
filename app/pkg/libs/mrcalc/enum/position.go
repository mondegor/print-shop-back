package enum

import (
	"encoding/json"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	PositionTop    Position = iota // PositionTop -
	PositionOnside                 // PositionOnside -
	PositionBottom                 // PositionBottom -

	enumNamePosition = "Position"
)

type (
	// Position - вид расположения элементов на листе прямоугольного формата.
	Position uint8
)

var (
	positionName = map[Position]string{ //nolint:gochecknoglobals
		PositionTop:    "TOP",
		PositionOnside: "ONSIDE",
		PositionBottom: "BOTTOM",
	}

	positionValue = map[string]Position{ //nolint:gochecknoglobals
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

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNamePosition)
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
