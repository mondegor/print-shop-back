package enum

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_                  PaperSide = iota
	PaperSideSame                // PaperSideSame - comment const
	PaperSideDifferent           // PaperSideDifferent - comment const

	paperSideLast     = uint8(PaperSideDifferent)
	enumNamePaperSide = "PaperSide"
)

type (
	// PaperSide - comment type.
	PaperSide uint8
)

var (
	paperSideName = map[PaperSide]string{ //nolint:gochecknoglobals
		PaperSideSame:      "SAME",
		PaperSideDifferent: "DIFFERENT",
	}

	paperSideValue = map[string]PaperSide{ //nolint:gochecknoglobals
		"SAME":      PaperSideSame,
		"DIFFERENT": PaperSideDifferent,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *PaperSide) ParseAndSet(value string) error {
	if parsedValue, ok := paperSideValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNamePaperSide)
}

// Set - устанавливает указанное значение, если оно является enum значением.
func (e *PaperSide) Set(value uint8) error {
	if value > 0 && value <= paperSideLast {
		*e = PaperSide(value)

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNamePaperSide)
}

// String - comment method.
func (e PaperSide) String() string {
	return paperSideName[e]
}

// Empty - проверяет, что enum значение не установлено.
func (e PaperSide) Empty() bool {
	return e == 0
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e PaperSide) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *PaperSide) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *PaperSide) Scan(value any) error {
	if val, ok := value.(int64); ok {
		return e.Set(uint8(val))
	}

	return mrcore.ErrInternalTypeAssertion.New(enumNamePaperSide, value)
}

// Value implements the driver.Valuer interface.
func (e PaperSide) Value() (driver.Value, error) {
	return uint8(e), nil
}

// ParsePaperSideList - парсит массив строковых значений и
// возвращает соответствующий массив enum значений.
// ParsePaperSideList - comment func.
func ParsePaperSideList(items []string) ([]PaperSide, error) {
	var tmp PaperSide

	parsedItems := make([]PaperSide, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
