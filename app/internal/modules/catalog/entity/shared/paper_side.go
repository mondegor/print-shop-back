package entity_shared

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_ PaperSide = iota
	PaperSideSame
	PaperSideDifferent

	paperSideLast     = uint8(PaperSideDifferent)
	enumNamePaperSide = "PaperSide"
)

type (
	PaperSide uint8
)

var (
	paperSideName = map[PaperSide]string{
		PaperSideSame:      "SAME",
		PaperSideDifferent: "DIFFERENT",
	}

	paperSideValue = map[string]PaperSide{
		"SAME":      PaperSideSame,
		"DIFFERENT": PaperSideDifferent,
	}
)

func (e *PaperSide) ParseAndSet(value string) error {
	if parsedValue, ok := paperSideValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNamePaperSide)
}

func (e *PaperSide) Set(value uint8) error {
	if value > 0 && value <= paperSideLast {
		*e = PaperSide(value)
		return nil
	}

	return fmt.Errorf("number '%d' is not registered in %s", value, enumNamePaperSide)
}

func (e PaperSide) String() string {
	return paperSideName[e]
}

func (e PaperSide) Empty() bool {
	return e == 0
}

func (e PaperSide) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *PaperSide) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *PaperSide) Scan(value any) error {
	if val, ok := value.(int64); ok {
		return e.Set(uint8(val))
	}

	return mrcore.FactoryErrInternalTypeAssertion.New(enumNamePaperSide, value)
}

// Value implements the driver Valuer interface.
func (e PaperSide) Value() (driver.Value, error) {
	return uint8(e), nil
}

func ParsePaperSideList(items []string) ([]PaperSide, error) {
	var tmp PaperSide
	parsedItems := make([]PaperSide, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
