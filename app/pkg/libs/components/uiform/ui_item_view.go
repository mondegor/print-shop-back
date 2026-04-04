package uiform

import (
	"encoding/json"

	"github.com/mondegor/go-sysmess/errors"
)

// Типы элементов интерфейса.
const (
	UIItemViewBlock UIItemView = iota + 1 // блок
	UIItemViewCheck                       // check box
	UIItemViewCombo                       // combo box
	UIItemViewList                        // список
	UIItemViewRadio                       // radio box
	UIItemViewRange                       // диапазон (2 поля)
	UIItemViewText                        // текс
)

const (
	enumNameUIItemView = "UIItemView"
)

type (
	// UIItemView - comment type.
	UIItemView uint8
)

//nolint:gochecknoglobals
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

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *UIItemView) ParseAndSet(value string) error {
	if parsedValue, ok := uiItemViewValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return errors.ErrInternalKeyNotFoundInSource.New(
		"key", value,
		"source", enumNameUIItemView,
	)
}

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
