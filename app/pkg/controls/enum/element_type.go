//nolint:dupl
package enum

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_                ElementType = iota
	ElementTypeGroup             // ElementTypeGroup - группа (блок) с названием, способная вмещать другие элементы
	ElementTypeList              // ElementTypeList - массив элементов

	elementTypeLast     = uint8(ElementTypeList)
	enumNameElementType = "ElementType"
)

type (
	// ElementType - тип элемента формы.
	ElementType uint8
)

var (
	elementTypeName = map[ElementType]string{ //nolint:gochecknoglobals
		ElementTypeGroup: "GROUP",
		ElementTypeList:  "ELEMENT_LIST",
	}

	elementTypeValue = map[string]ElementType{ //nolint:gochecknoglobals
		"GROUP":        ElementTypeGroup,
		"ELEMENT_LIST": ElementTypeList,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *ElementType) ParseAndSet(value string) error {
	if parsedValue, ok := elementTypeValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameElementType)
}

// Set - устанавливает указанное значение, если оно является enum значением.
func (e *ElementType) Set(value uint8) error {
	if value > 0 && value <= elementTypeLast {
		*e = ElementType(value)

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameElementType)
}

// String - возвращает значение в виде строки.
func (e ElementType) String() string {
	return elementTypeName[e]
}

// Empty - проверяет, что enum значение не установлено.
func (e ElementType) Empty() bool {
	return e == 0
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e ElementType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *ElementType) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ElementType) Scan(value any) error {
	if val, ok := value.(int64); ok {
		return e.Set(uint8(val))
	}

	return mrcore.ErrInternalTypeAssertion.New(enumNameElementType, value)
}

// Value implements the driver.Valuer interface.
func (e ElementType) Value() (driver.Value, error) {
	return uint8(e), nil
}

// ParseElementTypeList - парсит массив строковых значений и
// возвращает соответствующий массив enum значений.
func ParseElementTypeList(items []string) ([]ElementType, error) {
	var tmp ElementType

	parsedItems := make([]ElementType, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
