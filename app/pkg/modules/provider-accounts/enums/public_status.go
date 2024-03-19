package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
)

const (
	_ PublicStatus = iota
	PublicStatusDraft
	PublicStatusHidden
	PublicStatusPublished
	PublicStatusPublishedShared

	publicStatusLast     = uint8(PublicStatusPublishedShared)
	enumNamePublicStatus = "PublicStatus"
)

type (
	PublicStatus uint8
)

var (
	publicStatusName = map[PublicStatus]string{
		PublicStatusDraft:           "DRAFT",
		PublicStatusHidden:          "HIDDEN",
		PublicStatusPublished:       "PUBLISHED",
		PublicStatusPublishedShared: "PUBLISHED_SHARED",
	}

	publicStatusValue = map[string]PublicStatus{
		"DRAFT":            PublicStatusDraft,
		"HIDDEN":           PublicStatusHidden,
		"PUBLISHED":        PublicStatusPublished,
		"PUBLISHED_SHARED": PublicStatusPublishedShared,
	}

	PublicStatusFlow = mrenum.StatusFlow{
		PublicStatusDraft: {
			PublicStatusPublished,
			PublicStatusPublishedShared,
		},
		PublicStatusHidden: {
			PublicStatusPublished,
			PublicStatusPublishedShared,
		},
		PublicStatusPublished: {
			PublicStatusHidden,
			PublicStatusPublishedShared,
		},
		PublicStatusPublishedShared: {
			PublicStatusHidden,
			PublicStatusPublished,
		},
	}
)

func (e *PublicStatus) ParseAndSet(value string) error {
	if parsedValue, ok := publicStatusValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNamePublicStatus)
}

func (e *PublicStatus) Set(value uint8) error {
	if value > 0 && value <= publicStatusLast {
		*e = PublicStatus(value)
		return nil
	}

	return fmt.Errorf("number '%d' is not registered in %s", value, enumNamePublicStatus)
}

func (e PublicStatus) String() string {
	return publicStatusName[e]
}

func (e PublicStatus) Empty() bool {
	return e == 0
}

func (e PublicStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *PublicStatus) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *PublicStatus) Scan(value any) error {
	if val, ok := value.(int64); ok {
		return e.Set(uint8(val))
	}

	return mrcore.FactoryErrInternalTypeAssertion.New(enumNamePublicStatus, value)
}

// Value implements the driver Valuer interface.
func (e PublicStatus) Value() (driver.Value, error) {
	return uint8(e), nil
}

func ParsePublicStatusList(items []string) ([]PublicStatus, error) {
	var tmp PublicStatus
	parsedItems := make([]PublicStatus, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
