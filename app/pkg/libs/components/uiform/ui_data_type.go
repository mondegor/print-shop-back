package uiform

import (
	"encoding/json"
	"fmt"
)

const (
	_ UIDataType = iota
	UIDataTypeBoolean
	UIDataTypeGroup
	UIDataTypeEnum
	UIDataTypeNumber
	UIDataTypeString

	// uiDataTypeLast     = uint8(UIDataTypeString)
	enumNameUIDataType = "UIDataType"
)

type (
	UIDataType uint8
)

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

func (e *UIDataType) ParseAndSet(value string) error {
	if parsedValue, ok := uiDataTypeValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNameUIDataType)
}

//func (e *UIDataType) Set(value uint8) error {
//	if value > 0 && value <= uiDataTypeLast {
//		*e = UIDataType(value)
//		return nil
//	}
//
//	return fmt.Errorf("number '%d' is not registered in %s", value, enumNameUIDataType)
//}

func (e UIDataType) String() string {
	return uiDataTypeName[e]
}

func (e UIDataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e UIDataType) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}
