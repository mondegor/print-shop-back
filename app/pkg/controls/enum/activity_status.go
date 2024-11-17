package enum

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_                       ActivityStatus = iota
	ActivityStatusDraft                    // ActivityStatusDraft - comment const
	ActivityStatusTesting                  // ActivityStatusTesting - comment const
	ActivityStatusPublished                // ActivityStatusPublished - comment const
	ActivityStatusArchived                 // ActivityStatusArchived - comment const

	activityStatusLast     = uint8(ActivityStatusArchived)
	enumNameActivityStatus = "ActivityStatus"
)

type (
	// ActivityStatus - comment type.
	ActivityStatus uint8
)

var (
	activityStatusName = map[ActivityStatus]string{ //nolint:gochecknoglobals
		ActivityStatusDraft:     "DRAFT",
		ActivityStatusTesting:   "TESTING",
		ActivityStatusPublished: "PUBLISHED",
		ActivityStatusArchived:  "ARCHIVED",
	}

	activityStatusValue = map[string]ActivityStatus{ //nolint:gochecknoglobals
		"DRAFT":     ActivityStatusDraft,
		"TESTING":   ActivityStatusTesting,
		"PUBLISHED": ActivityStatusPublished,
		"ARCHIVED":  ActivityStatusArchived,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *ActivityStatus) ParseAndSet(value string) error {
	if parsedValue, ok := activityStatusValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameActivityStatus)
}

// Set - устанавливает указанное значение, если оно является enum значением.
func (e *ActivityStatus) Set(value uint8) error {
	if value > 0 && value <= activityStatusLast {
		*e = ActivityStatus(value)

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameActivityStatus)
}

// String - возвращает значение в виде строки.
func (e ActivityStatus) String() string {
	return activityStatusName[e]
}

// Empty - проверяет, что enum значение не установлено.
func (e ActivityStatus) Empty() bool {
	return e == 0
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e ActivityStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *ActivityStatus) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ActivityStatus) Scan(value any) error {
	if val, ok := value.(int64); ok {
		return e.Set(uint8(val))
	}

	return mrcore.ErrInternalTypeAssertion.New(enumNameActivityStatus, value)
}

// Value implements the driver.Valuer interface.
func (e ActivityStatus) Value() (driver.Value, error) {
	return uint8(e), nil
}

// ParseActivityStatusList - парсит массив строковых значений и
// возвращает соответствующий массив enum значений.
// ParseActivityStatusList - comment func.
func ParseActivityStatusList(items []string) ([]ActivityStatus, error) {
	var tmp ActivityStatus

	parsedItems := make([]ActivityStatus, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
