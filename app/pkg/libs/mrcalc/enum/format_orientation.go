package enum

import (
	"encoding/json"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	FormatOrientationBook  Orientation = iota // FormatOrientationBook - книжная ориентация
	FormatOrientationAlbum                    // FormatOrientationAlbum - альбомная ориентация

	enumNameOrientation = "Orientation"
)

type (
	// Orientation - вид расположения элементов на листе прямоугольного формата.
	Orientation uint8
)

var (
	orientationName = map[Orientation]string{ //nolint:gochecknoglobals
		FormatOrientationBook:  "BOOK",
		FormatOrientationAlbum: "ALBUM",
	}

	orientationValue = map[string]Orientation{ //nolint:gochecknoglobals
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

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameOrientation)
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
