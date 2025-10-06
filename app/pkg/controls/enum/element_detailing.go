//nolint:dupl
package enum

import (
	"database/sql/driver"
	"encoding/json"
	"math"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

// Типы детализации элемента.
const (
	ElementDetailingNormal   ElementDetailing = iota + 1 // обычная детализация
	ElementDetailingExtended                             // расширенная детализация
)

const (
	elementDetailingLast     = uint8(ElementDetailingExtended)
	enumNameElementDetailing = "ElementDetailing"
)

type (
	// ElementDetailing - comment type.
	ElementDetailing uint8
)

var (
	elementDetailingName = map[ElementDetailing]string{ //nolint:gochecknoglobals
		ElementDetailingNormal:   "NORMAL",
		ElementDetailingExtended: "EXTENDED",
	}

	elementDetailingValue = map[string]ElementDetailing{ //nolint:gochecknoglobals
		"NORMAL":   ElementDetailingNormal,
		"EXTENDED": ElementDetailingExtended,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *ElementDetailing) ParseAndSet(value string) error {
	if parsedValue, ok := elementDetailingValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNameElementDetailing)
}

// Set - устанавливает указанное значение, если оно является enum значением.
func (e *ElementDetailing) Set(value uint8) error {
	if value > 0 && value <= elementDetailingLast {
		*e = ElementDetailing(value)

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNameElementDetailing)
}

// String - возвращает значение в виде строки.
func (e ElementDetailing) String() string {
	return elementDetailingName[e]
}

// // Empty - сообщает, установлено ли enum значение.
// func (e ElementDetailing) Empty() bool {
// 	return e == 0
// }

// MarshalJSON - переводит enum значение в строковое представление.
func (e ElementDetailing) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *ElementDetailing) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ElementDetailing) Scan(value any) error {
	if val, ok := value.(int64); ok && val >= 0 && val <= math.MaxUint8 {
		return e.Set(uint8(val)) //nolint:gosec
	}

	return mr.ErrInternalTypeAssertion.New(enumNameElementDetailing, value)
}

// Value implements the driver.Valuer interface.
func (e ElementDetailing) Value() (driver.Value, error) {
	return uint8(e), nil
}

// ParseElementDetailingList - парсит массив строковых значений и
// возвращает соответствующий массив enum значений.
func ParseElementDetailingList(items []string) ([]ElementDetailing, error) {
	var tmp ElementDetailing

	parsedItems := make([]ElementDetailing, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
