package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	_ ElementType = iota
	ElementTypeGroup
	ElementTypeList

	elementTypeLast     = uint8(ElementTypeList)
	enumNameElementType = "ElementType"
)

type (
	ElementType uint8
)

var (
	elementTypeName = map[ElementType]string{
		ElementTypeGroup: "GROUP",
		ElementTypeList:  "ELEMENT_LIST",
	}

	elementTypeValue = map[string]ElementType{
		"GROUP":        ElementTypeGroup,
		"ELEMENT_LIST": ElementTypeList,
	}
)

func (e *ElementType) ParseAndSet(value string) error {
	if parsedValue, ok := elementTypeValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNameElementType)
}

func (e *ElementType) Set(value uint8) error {
	if value > 0 && value <= elementTypeLast {
		*e = ElementType(value)
		return nil
	}

	return fmt.Errorf("number '%d' is not registered in %s", value, enumNameElementType)
}

func (e ElementType) String() string {
	return elementTypeName[e]
}

func (e ElementType) Empty() bool {
	return e == 0
}

func (e ElementType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *ElementType) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ElementType) Scan(value any) error {
	if val, ok := value.(int64); ok {
		return e.Set(uint8(val))
	}

	return mrcore.FactoryErrInternalTypeAssertion.New(enumNameElementType, value)
}

// Value implements the driver Valuer interface.
func (e ElementType) Value() (driver.Value, error) {
	return uint8(e), nil
}

func ParseElementTypeList(items []string) ([]ElementType, error) {
	var tmp ElementType
	parsedItems := make([]ElementType, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
