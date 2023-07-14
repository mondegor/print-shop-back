package entity

import (
    "print-shop-back/pkg/mrerr"
    "encoding/json"
)

type ItemDetailing uint8

const (
    _ ItemDetailing = iota
    ItemDetailingNormal
    ItemDetailingExtended
)

var (
    itemDetailingName = map[ItemDetailing]string{
        ItemDetailingNormal: "NORMAL",
        ItemDetailingExtended: "EXTENDED",
    }

    itemDetailingValue = map[string]ItemDetailing{
        "NORMAL": ItemDetailingNormal,
        "EXTENDED": ItemDetailingExtended,
    }
)

func (e *ItemDetailing) ParseAndSet(value string) error {
    if parsedValue, ok := itemDetailingValue[value]; ok {
        *e = parsedValue
        return nil
    }

    return mrerr.ErrInternalMapValueNotFound.New(value, "ItemDetailing")
}

func (e ItemDetailing) String() string {
    return itemDetailingName[e]
}

func (e ItemDetailing) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.String())
}

func (e *ItemDetailing) UnmarshalJSON(data []byte) error {
    var value string
    err := json.Unmarshal(data, &value)

    if err != nil {
        return err
    }

    return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ItemDetailing) Scan(value any) error {
    if val, ok := value.(string); ok {
        return e.ParseAndSet(val)
    }

    return mrerr.ErrInternalTypeAssertion.New("ItemDetailing", value)
}
