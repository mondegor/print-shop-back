package entity

import (
    "encoding/json"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    _ ResourceStatus = iota
    ResourceStatusDraft
    ResourceStatusHidden
    ResourceStatusPublished
)

type (
    ResourceStatus uint8
)

var (
    resourceStatusName = map[ResourceStatus]string{
        ResourceStatusDraft: "DRAFT",
        ResourceStatusHidden: "HIDDEN",
        ResourceStatusPublished: "PUBLISHED",
    }

    resourceStatusValue = map[string]ResourceStatus{
        "DRAFT": ResourceStatusDraft,
        "HIDDEN": ResourceStatusHidden,
        "PUBLISHED": ResourceStatusPublished,
    }
)

func (e *ResourceStatus) ParseAndSet(value string) error {
    if parsedValue, ok := resourceStatusValue[value]; ok {
        *e = parsedValue
        return nil
    }

    return mrcore.FactoryErrInternalMapValueNotFound.New(value, "ResourceStatus")
}

func (e ResourceStatus) String() string {
    return resourceStatusName[e]
}

func (e ResourceStatus) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.String())
}

func (e *ResourceStatus) UnmarshalJSON(data []byte) error {
    var value string
    err := json.Unmarshal(data, &value)

    if err != nil {
        return err
    }

    return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ResourceStatus) Scan(value any) error {
    if val, ok := value.(string); ok {
        return e.ParseAndSet(val)
    }

    return mrcore.FactoryErrInternalTypeAssertion.New("ResourceStatus", value)
}
