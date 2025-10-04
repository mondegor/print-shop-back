package enum

import (
	"database/sql/driver"
	"encoding/json"
	"math"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

// Типы сторон бумаги.
const (
	PaperSideSame      PaperSide = iota + 1 // стороны одинаковые
	PaperSideDifferent                      // стороны различаются
)

const (
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

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNamePaperSide)
}

// Set - устанавливает указанное значение, если оно является enum значением.
func (e *PaperSide) Set(value uint8) error {
	if value <= paperSideLast {
		*e = PaperSide(value)

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNamePaperSide)
}

// String - возвращает значение в виде строки.
func (e PaperSide) String() string {
	return paperSideName[e]
}

//
// // Empty - сообщает, установлено ли enum значение.
// func (e PaperSide) Empty() bool {
// 	return e == 0
// }

// MarshalJSON - переводит enum значение в строковое представление.
// !!!! при формировании json переводит в текстовый вид.
func (e PaperSide) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
// !!!! из json переводится в числовой вид.
func (e *PaperSide) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
// !!!! не допускает, чтобы из БД пришло неизвестное значение.
func (e *PaperSide) Scan(value any) error {
	if val, ok := value.(int64); ok && val >= 0 && val <= math.MaxUint8 {
		return e.Set(uint8(val))
	}

	return mr.ErrInternalTypeAssertion.New(enumNamePaperSide, value)
}

// Value implements the driver.Valuer interface.
// !!!! можно возвращать nil значение, если 0.
func (e PaperSide) Value() (driver.Value, error) {
	return uint8(e), nil
}
