package uiform

import (
	"encoding/json"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_                 UIDataType = iota
	UIDataTypeBoolean            // UIDataTypeBoolean - comment const
	UIDataTypeGroup              // UIDataTypeGroup - comment const
	UIDataTypeEnum               // UIDataTypeEnum - comment const
	UIDataTypeNumber             // UIDataTypeNumber - comment const
	UIDataTypeString             // UIDataTypeString - comment const

	// uiDataTypeLast = uint8(UIDataTypeString).

	enumNameUIDataType = "UIDataType"
)

type (
	// UIDataType - comment type.
	UIDataType uint8
)

var (
	uiDataTypeName = map[UIDataType]string{ //nolint:gochecknoglobals
		UIDataTypeBoolean: "BOOLEAN",
		UIDataTypeGroup:   "GROUP",
		UIDataTypeEnum:    "ENUM",
		UIDataTypeNumber:  "NUMBER",
		UIDataTypeString:  "STRING",
	}

	uiDataTypeValue = map[string]UIDataType{ //nolint:gochecknoglobals
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

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameUIDataType)
}

// Set - устанавливает указанное значение, если оно является enum значением.
// func (e *UIDataType) Set(value uint8) error {
// 	if value > 0 && value <= uiDataTypeLast {
// 		*e = UIDataType(value)
//
// 		return nil
// 	}
//
// 	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameUIDataType)
// }

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
