package entity

import (
    "calc-user-data-back-adm/pkg/mrerr"
    "encoding/json"
)

type ItemStatus uint8

const (
    _ ItemStatus = iota
    ItemStatusDraft
    ItemStatusEnabled
    ItemStatusDisabled
    ItemStatusRemoved
)

var (
    itemStatusName = map[ItemStatus]string{
        ItemStatusDraft: "DRAFT",
        ItemStatusEnabled: "ENABLED",
        ItemStatusDisabled: "DISABLED",
        ItemStatusRemoved: "REMOVED",
    }

    itemStatusValue = map[string]ItemStatus{
        "DRAFT": ItemStatusDraft,
        "ENABLED": ItemStatusEnabled,
        "DISABLED": ItemStatusDisabled,
        "REMOVED": ItemStatusRemoved,
    }
)

func (e *ItemStatus) ParseAndSet(value string) error {
    if parsedValue, ok := itemStatusValue[value]; ok {
        *e = parsedValue
        return nil
    }

    return mrerr.ErrInternalMapValueNotFound.New(value, "ItemStatus")
}

func (e ItemStatus) String() string {
    return itemStatusName[e]
}

func (e ItemStatus) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.String())
}

func (e *ItemStatus) UnmarshalJSON(data []byte) error {
    var value string
    err := json.Unmarshal(data, &value)

    if err != nil {
        return err
    }

    return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ItemStatus) Scan(value any) error {
    if val, ok := value.(string); ok {
        return e.ParseAndSet(val)
    }

    return mrerr.ErrInternalTypeAssertion.New("ItemStatus", value)
}
