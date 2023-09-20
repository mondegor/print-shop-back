package entity

import (
    "encoding/json"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    _ UIItemView = iota
    UIItemViewBlock
    UIItemViewCheck
    UIItemViewCombo
    UIItemViewList
    UIItemViewRadio
    UIItemViewRange
    UIItemViewText
)

type (
    UIItemView uint8
)

var (
    uiItemViewName = map[UIItemView]string{
        UIItemViewBlock: "block",
        UIItemViewCheck: "check",
        UIItemViewCombo: "combo",
        UIItemViewList: "list",
        UIItemViewRadio: "radio",
        UIItemViewRange: "range",
        UIItemViewText: "text",
    }

    uiItemViewValue = map[string]UIItemView{
        "block": UIItemViewBlock,
        "check": UIItemViewCheck,
        "combo": UIItemViewCombo,
        "list": UIItemViewList,
        "radio": UIItemViewRadio,
        "range": UIItemViewRange,
        "text": UIItemViewText,
    }
)

func (e *UIItemView) ParseAndSet(value string) error {
    if parsedValue, ok := uiItemViewValue[value]; ok {
        *e = parsedValue
        return nil
    }

    return mrcore.FactoryErrInternalMapValueNotFound.New(value, "UIItemView")
}

func (e UIItemView) String() string {
    return uiItemViewName[e]
}

func (e UIItemView) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.String())
}

func (e *UIItemView) UnmarshalJSON(data []byte) error {
    var value string
    err := json.Unmarshal(data, &value)

    if err != nil {
        return err
    }

    return e.ParseAndSet(value)
}

//// Scan implements the Scanner interface.
//func (e *UIItemView) Scan(value any) error {
//    if val, ok := value.(string); ok {
//        return e.ParseAndSet(val)
//    }
//
//    return mrcore.FactoryErrInternalTypeAssertion.New("UIItemView", value)
//}
