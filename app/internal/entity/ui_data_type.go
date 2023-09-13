package entity

import (
    "encoding/json"

    "github.com/mondegor/go-webcore/mrcore"
)

type UIDataType uint8

const (
    _ UIDataType = iota
    UIDataTypeBoolean
    UIDataTypeGroup
    UIDataTypeEnum
    UIDataTypeNumber
    UIDataTypeString
)

var (
    uiDataTypeName = map[UIDataType]string{
        UIDataTypeBoolean: "boolean",
        UIDataTypeGroup: "group",
        UIDataTypeEnum: "enum",
        UIDataTypeNumber: "number",
        UIDataTypeString: "string",
    }

    uiDataTypeValue = map[string]UIDataType{
        "boolean": UIDataTypeBoolean,
        "group": UIDataTypeGroup,
        "enum": UIDataTypeEnum,
        "number": UIDataTypeNumber,
        "string": UIDataTypeString,
    }
)

func (e *UIDataType) ParseAndSet(value string) error {
    if parsedValue, ok := uiDataTypeValue[value]; ok {
        *e = parsedValue
        return nil
    }

    return mrcore.FactoryErrInternalMapValueNotFound.New(value, "UIDataType")
}

func (e UIDataType) String() string {
    return uiDataTypeName[e]
}

func (e UIDataType) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.String())
}

func (e *UIDataType) UnmarshalJSON(data []byte) error {
    var value string
    err := json.Unmarshal(data, &value)

    if err != nil {
        return err
    }

    return e.ParseAndSet(value)
}

//// Scan implements the Scanner interface.
//func (e *UIDataType) Scan(value any) error {
//    if val, ok := value.(string); ok {
//        return e.ParseAndSet(val)
//    }
//
//    return mrcore.FactoryErrInternalTypeAssertion.New("UIDataType", value)
//}
