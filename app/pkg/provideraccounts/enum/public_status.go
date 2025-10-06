package enum

import (
	"database/sql/driver"
	"encoding/json"
	"math"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

// Статусы публикации страницы.
const (
	PublicStatusDraft           PublicStatus = iota + 1 // черновик
	PublicStatusHidden                                  // скрытый для всех пользователей
	PublicStatusPublished                               // опубликован
	PublicStatusPublishedShared                         // опубликован и присутствует в каталоге
)

const (
	publicStatusLast     = uint8(PublicStatusPublishedShared)
	enumNamePublicStatus = "PublicStatus"
)

type (
	// PublicStatus - comment type.
	PublicStatus uint8
)

var (
	publicStatusName = map[PublicStatus]string{ //nolint:gochecknoglobals
		PublicStatusDraft:           "DRAFT",
		PublicStatusHidden:          "HIDDEN",
		PublicStatusPublished:       "PUBLISHED",
		PublicStatusPublishedShared: "PUBLISHED_SHARED",
	}

	publicStatusValue = map[string]PublicStatus{ //nolint:gochecknoglobals
		"DRAFT":            PublicStatusDraft,
		"HIDDEN":           PublicStatusHidden,
		"PUBLISHED":        PublicStatusPublished,
		"PUBLISHED_SHARED": PublicStatusPublishedShared,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *PublicStatus) ParseAndSet(value string) error {
	if parsedValue, ok := publicStatusValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNamePublicStatus)
}

// Set - устанавливает указанное значение, если оно является enum значением.
func (e *PublicStatus) Set(value uint8) error {
	if value > 0 && value <= publicStatusLast {
		*e = PublicStatus(value)

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNamePublicStatus)
}

// String - возвращает значение в виде строки.
func (e PublicStatus) String() string {
	return publicStatusName[e]
}

// // Empty - сообщает, установлено ли enum значение.
// func (e PublicStatus) Empty() bool {
// 	return e == 0
// }

// MarshalJSON - переводит enum значение в строковое представление.
func (e PublicStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *PublicStatus) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *PublicStatus) Scan(value any) error {
	if val, ok := value.(int64); ok && val >= 0 && val <= math.MaxUint8 {
		return e.Set(uint8(val)) //nolint:gosec
	}

	return mr.ErrInternalTypeAssertion.New(enumNamePublicStatus, value)
}

// Value implements the driver.Valuer interface.
func (e PublicStatus) Value() (driver.Value, error) {
	return uint8(e), nil
}

// ParsePublicStatusList - парсит массив строковых значений и
// возвращает соответствующий массив enum значений.
func ParsePublicStatusList(items []string) ([]PublicStatus, error) {
	var tmp PublicStatus

	parsedItems := make([]PublicStatus, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
