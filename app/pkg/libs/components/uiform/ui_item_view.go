package uiform

import (
	"encoding/json"
	"fmt"
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

	// uiItemViewLast     = uint8(UIItemViewText)
	enumNameUIItemView = "UIItemView"
)

type (
	UIItemView uint8
)

var (
	uiItemViewName = map[UIItemView]string{
		UIItemViewBlock: "BLOCK",
		UIItemViewCheck: "CHECK",
		UIItemViewCombo: "COMBO",
		UIItemViewList:  "LIST",
		UIItemViewRadio: "RADIO",
		UIItemViewRange: "RANGE",
		UIItemViewText:  "TEXT",
	}

	uiItemViewValue = map[string]UIItemView{
		"BLOCK": UIItemViewBlock,
		"CHECK": UIItemViewCheck,
		"COMBO": UIItemViewCombo,
		"LIST":  UIItemViewList,
		"RADIO": UIItemViewRadio,
		"RANGE": UIItemViewRange,
		"TEXT":  UIItemViewText,
	}
)

func (e *UIItemView) ParseAndSet(value string) error {
	if parsedValue, ok := uiItemViewValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNameUIItemView)
}

//func (e *UIItemView) Set(value uint8) error {
//	if value > 0 && value <= uiItemViewLast {
//		*e = UIItemView(value)
//		return nil
//	}
//
//	return fmt.Errorf("number '%d' is not registered in %s", value, enumNameUIItemView)
//}

func (e UIItemView) String() string {
	return uiItemViewName[e]
}

func (e UIItemView) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *UIItemView) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}
