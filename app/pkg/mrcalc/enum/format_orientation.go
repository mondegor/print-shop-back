package enum

import (
	"encoding/json"

	"github.com/mondegor/go-sysmess/errors"
)

// Ориентации формата.
const (
	FormatOrientationBook  Orientation = iota + 1 // книжная ориентация
	FormatOrientationAlbum                        // альбомная ориентация
)

const (
	enumNameOrientation = "Orientation"
)

type (
	// Orientation - вид расположения элементов на листе прямоугольного формата.
	Orientation uint8
)

//nolint:gochecknoglobals
var (
	orientationName = map[Orientation]string{
		FormatOrientationBook:  "BOOK",
		FormatOrientationAlbum: "ALBUM",
	}

	orientationValue = map[string]Orientation{
		"BOOK":  FormatOrientationBook,
		"ALBUM": FormatOrientationAlbum,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *Orientation) ParseAndSet(value string) error {
	if parsedValue, ok := orientationValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return errors.ErrInternalKeyNotFoundInSource.New(
		"key", value,
		"source", enumNameOrientation,
	)
}

// String - возвращает значение в виде строки.
func (e Orientation) String() string {
	return orientationName[e]
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e Orientation) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *Orientation) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}
