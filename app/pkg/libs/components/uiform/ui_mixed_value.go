package uiform

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	modelNameUIMixedValue = "UIMixedValue"
	quotesByte            = 34
)

type (
	// UIMixedValue - comment struct.
	UIMixedValue struct {
		StringValue string
		FloatValue  float64
		IsString    bool
	}
)

// String - comment method.
func (v UIMixedValue) String() string {
	if v.IsString {
		return v.StringValue
	}

	return strconv.FormatFloat(v.FloatValue, 'f', 4, 64)
}

// MarshalJSON - переводит данные в строковое представление в зависимости от текущего типа данных.
func (v UIMixedValue) MarshalJSON() ([]byte, error) {
	if v.IsString {
		return json.Marshal(v.StringValue)
	}

	return json.Marshal(v.FloatValue)
}

// UnmarshalJSON - переводит строковые данные в соответствующих им тип.
func (v *UIMixedValue) UnmarshalJSON(data []byte) error {
	v.IsString = data[0] == quotesByte

	if v.IsString {
		if err := json.Unmarshal(data, &v.StringValue); err != nil {
			return fmt.Errorf("%s: '%s' is not parsed (%w)", modelNameUIMixedValue, v.StringValue, err)
		}
	} else {
		if err := json.Unmarshal(data, &v.FloatValue); err != nil {
			return fmt.Errorf("%s: '%f' is not parsed (%w)", modelNameUIMixedValue, v.FloatValue, err)
		}
	}

	return nil
}
