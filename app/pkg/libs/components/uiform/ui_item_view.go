package uiform

import (
	"encoding/json"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_               UIItemView = iota
	UIItemViewBlock            // UIItemViewBlock - comment const
	UIItemViewCheck            // UIItemViewCheck - comment const
	UIItemViewCombo            // UIItemViewCombo - comment const
	UIItemViewList             // UIItemViewList - comment const
	UIItemViewRadio            // UIItemViewRadio - comment const
	UIItemViewRange            // UIItemViewRange - comment const
	UIItemViewText             // UIItemViewText - comment const

	// uiItemViewLast     = uint8(UIItemViewText).
	enumNameUIItemView = "UIItemView"
)

type (
	// UIItemView - comment type.
	UIItemView uint8
)

var (
	uiItemViewName = map[UIItemView]string{ //nolint:gochecknoglobals
		UIItemViewBlock: "BLOCK",
		UIItemViewCheck: "CHECK",
		UIItemViewCombo: "COMBO",
		UIItemViewList:  "LIST",
		UIItemViewRadio: "RADIO",
		UIItemViewRange: "RANGE",
		UIItemViewText:  "TEXT",
	}

	uiItemViewValue = map[string]UIItemView{ //nolint:gochecknoglobals
		"BLOCK": UIItemViewBlock,
		"CHECK": UIItemViewCheck,
		"COMBO": UIItemViewCombo,
		"LIST":  UIItemViewList,
		"RADIO": UIItemViewRadio,
		"RANGE": UIItemViewRange,
		"TEXT":  UIItemViewText,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *UIItemView) ParseAndSet(value string) error {
	if parsedValue, ok := uiItemViewValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameUIItemView)
}

// Set - устанавливает указанное значение, если оно является enum значением.
// func (e *UIItemView) Set(value uint8) error {
// 	if value > 0 && value <= uiItemViewLast {
// 		*e = UIItemView(value)
//
// 		return nil
// 	}
//
// 	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameUIItemView)
// }

// String - возвращает значение в виде строки.
func (e UIItemView) String() string {
	return uiItemViewName[e]
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e UIItemView) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *UIItemView) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}
