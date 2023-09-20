package entity

import (
    "encoding/json"
    "fmt"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    quotesByte = 34
)

type (
    UIMixedValue struct {
        FloatValue float64
        IsString bool
        StringValue string
    }
)

func (v UIMixedValue) String() string {
    if v.IsString {
        return v.StringValue
    }

    return fmt.Sprintf("%f", v.FloatValue)
}

func (v UIMixedValue) MarshalJSON() ([]byte, error) {
    if v.IsString {
        return json.Marshal(v.StringValue)
    }

    return json.Marshal(v.FloatValue)
}

func (v *UIMixedValue) UnmarshalJSON(data []byte) error {
    var err error

    v.IsString = data[0] == quotesByte

    if v.IsString {
        err = json.Unmarshal(data, &v.StringValue)
    } else {
        err = json.Unmarshal(data, &v.FloatValue)
    }

    if err != nil {
        return mrcore.FactoryErrInternalParseData.Wrap(err, "UIMixedValue", "JSON")
    }

    return nil
}
