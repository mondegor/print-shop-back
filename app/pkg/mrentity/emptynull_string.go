package mrentity

import (
    "database/sql/driver"
    "print-shop-back/pkg/mrerr"
)

type EmptynullString string

// Value implements the driver Valuer interface.
func (n EmptynullString) Value() (driver.Value, error) {
    if n == "" {
        return nil, nil
    }

    return string(n), nil
}

// Scan implements the Scanner interface.
func (n *EmptynullString) Scan(value any) error {
    if value == nil {
        *n = ""
        return nil
    }

    if val, ok := value.(string); ok {
        *n = EmptynullString(val)
        return nil
    }

    return mrerr.ErrInternalTypeAssertion.New("EmptynullString", value)
}
