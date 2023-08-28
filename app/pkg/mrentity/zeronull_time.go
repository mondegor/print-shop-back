package mrentity

import (
    "print-shop-back/pkg/mrerr"
    "database/sql/driver"
    "time"
)

type ZeronullTime time.Time

// Value implements the driver Valuer interface.
func (n ZeronullTime) Value() (driver.Value, error) {
    if time.Time(n).IsZero() {
        return nil, nil
    }

    return time.Time(n), nil
}

// Scan implements the Scanner interface.
func (n *ZeronullTime) Scan(value any) error {
    if value == nil {
        *n = ZeronullTime{}
        return nil
    }

    if val, ok := value.(time.Time); ok {
        *n = ZeronullTime(val)
        return nil
    }

    return mrerr.ErrInternalTypeAssertion.New("ZeronullTime", value)
}
