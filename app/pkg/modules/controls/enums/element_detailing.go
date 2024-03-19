package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_ ElementDetailing = iota
	ElementDetailingNormal
	ElementDetailingExtended

	elementDetailingLast     = uint8(ElementDetailingExtended)
	enumNameElementDetailing = "ElementDetailing"
)

type (
	ElementDetailing uint8
)

var (
	elementDetailingName = map[ElementDetailing]string{
		ElementDetailingNormal:   "NORMAL",
		ElementDetailingExtended: "EXTENDED",
	}

	elementDetailingValue = map[string]ElementDetailing{
		"NORMAL":   ElementDetailingNormal,
		"EXTENDED": ElementDetailingExtended,
	}
)

func (e *ElementDetailing) ParseAndSet(value string) error {
	if parsedValue, ok := elementDetailingValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNameElementDetailing)
}

func (e *ElementDetailing) Set(value uint8) error {
	if value > 0 && value <= elementDetailingLast {
		*e = ElementDetailing(value)
		return nil
	}

	return fmt.Errorf("number '%d' is not registered in %s", value, enumNameElementDetailing)
}

func (e ElementDetailing) String() string {
	return elementDetailingName[e]
}

func (e ElementDetailing) Empty() bool {
	return e == 0
}

func (e ElementDetailing) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *ElementDetailing) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ElementDetailing) Scan(value any) error {
	if val, ok := value.(int64); ok {
		return e.Set(uint8(val))
	}

	return mrcore.FactoryErrInternalTypeAssertion.New(enumNameElementDetailing, value)
}

// Value implements the driver Valuer interface.
func (e ElementDetailing) Value() (driver.Value, error) {
	return uint8(e), nil
}

func ParseElementDetailingList(items []string) ([]ElementDetailing, error) {
	var tmp ElementDetailing
	parsedItems := make([]ElementDetailing, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
