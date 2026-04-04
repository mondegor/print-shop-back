package uiform

import (
	"encoding/json"

	"github.com/mondegor/go-sysmess/errors"
)

// Типы данных.
const (
	UIDataTypeBoolean UIDataType = iota + 1
	UIDataTypeGroup
	UIDataTypeEnum
	UIDataTypeNumber
	UIDataTypeString
)

const (
	enumNameUIDataType = "UIDataType"
)

type (
	// UIDataType - comment type.
	UIDataType uint8
)

//nolint:gochecknoglobals
var (
	uiDataTypeName = map[UIDataType]string{
		UIDataTypeBoolean: "BOOLEAN",
		UIDataTypeGroup:   "GROUP",
		UIDataTypeEnum:    "ENUM",
		UIDataTypeNumber:  "NUMBER",
		UIDataTypeString:  "STRING",
	}

	uiDataTypeValue = map[string]UIDataType{
		"BOOLEAN": UIDataTypeBoolean,
		"GROUP":   UIDataTypeGroup,
		"ENUM":    UIDataTypeEnum,
		"NUMBER":  UIDataTypeNumber,
		"STRING":  UIDataTypeString,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *UIDataType) ParseAndSet(value string) error {
	if parsedValue, ok := uiDataTypeValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return errors.ErrInternalKeyNotFoundInSource.New(
		"key", value,
		"source", enumNameUIDataType,
	)
}

// String - возвращает значение в виде строки.
func (e UIDataType) String() string {
	return uiDataTypeName[e]
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e UIDataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e UIDataType) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}
