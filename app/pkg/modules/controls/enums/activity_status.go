package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
)

const (
	_ ActivityStatus = iota
	ActivityStatusTesting
	ActivityStatusPublished
	ActivityStatusArchived

	activityStatusLast     = uint8(ActivityStatusArchived)
	enumNameActivityStatus = "ActivityStatus"
)

type (
	ActivityStatus uint8
)

var (
	activityStatusName = map[ActivityStatus]string{
		ActivityStatusTesting:   "TESTING",
		ActivityStatusPublished: "PUBLISHED",
		ActivityStatusArchived:  "ARCHIVED",
	}

	activityStatusValue = map[string]ActivityStatus{
		"TESTING":   ActivityStatusTesting,
		"PUBLISHED": ActivityStatusPublished,
		"ARCHIVED":  ActivityStatusArchived,
	}

	ActivityStatusFlow = mrenum.StatusFlow{
		ActivityStatusTesting: {
			ActivityStatusPublished,
		},
		ActivityStatusPublished: {},
		ActivityStatusArchived:  {},
	}
)

func (e *ActivityStatus) ParseAndSet(value string) error {
	if parsedValue, ok := activityStatusValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNameActivityStatus)
}

func (e *ActivityStatus) Set(value uint8) error {
	if value > 0 && value <= activityStatusLast {
		*e = ActivityStatus(value)
		return nil
	}

	return fmt.Errorf("number '%d' is not registered in %s", value, enumNameActivityStatus)
}

func (e ActivityStatus) String() string {
	return activityStatusName[e]
}

func (e ActivityStatus) Empty() bool {
	return e == 0
}

func (e ActivityStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

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

	return mrcore.FactoryErrInternalTypeAssertion.New(enumNameActivityStatus, value)
}

// Value implements the driver Valuer interface.
func (e ActivityStatus) Value() (driver.Value, error) {
	return uint8(e), nil
}

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
