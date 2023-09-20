package entity

import (
    "encoding/json"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    _ FormFieldTemplateType = iota
    FormFieldTemplateTypeGroup
    FormFieldTemplateTypeFields
)

type (
    FormFieldTemplateType uint8
)

var (
    formFieldTemplateTypeName = map[FormFieldTemplateType]string{
        FormFieldTemplateTypeGroup: "GROUP",
        FormFieldTemplateTypeFields: "FIELDS",
    }

    formFieldTemplateTypeValue = map[string]FormFieldTemplateType{
        "GROUP": FormFieldTemplateTypeGroup,
        "FIELDS": FormFieldTemplateTypeFields,
    }
)

func (e *FormFieldTemplateType) ParseAndSet(value string) error {
    if parsedValue, ok := formFieldTemplateTypeValue[value]; ok {
        *e = parsedValue
        return nil
    }

    return mrcore.FactoryErrInternalMapValueNotFound.New(value, "FormFieldTemplateType")
}

func (e FormFieldTemplateType) String() string {
    return formFieldTemplateTypeName[e]
}

func (e FormFieldTemplateType) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.String())
}

func (e *FormFieldTemplateType) UnmarshalJSON(data []byte) error {
    var value string
    err := json.Unmarshal(data, &value)

    if err != nil {
        return err
    }

    return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *FormFieldTemplateType) Scan(value any) error {
    if val, ok := value.(string); ok {
        return e.ParseAndSet(val)
    }

    return mrcore.FactoryErrInternalTypeAssertion.New("FormFieldTemplateType", value)
}
