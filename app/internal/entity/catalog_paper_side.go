package entity

import (
    "encoding/json"
    "print-shop-back/pkg/mrerr"
)

type CatalogPaperSide uint8

const (
    _ CatalogPaperSide = iota
    CatalogPaperSideSame
    CatalogPaperSideDifferent
)

var (
    catalogPaperSideName = map[CatalogPaperSide]string{
        CatalogPaperSideSame: "SAME",
        CatalogPaperSideDifferent: "DIFFERENT",
    }

    catalogPaperSideValue = map[string]CatalogPaperSide{
        "SAME": CatalogPaperSideSame,
        "DIFFERENT": CatalogPaperSideDifferent,
    }
)

func (e *CatalogPaperSide) ParseAndSet(value string) error {
    if parsedValue, ok := catalogPaperSideValue[value]; ok {
        *e = parsedValue
        return nil
    }

    return mrerr.ErrInternalMapValueNotFound.New(value, "CatalogPaperSide")
}

func (e CatalogPaperSide) String() string {
    return catalogPaperSideName[e]
}

func (e CatalogPaperSide) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.String())
}

func (e *CatalogPaperSide) UnmarshalJSON(data []byte) error {
    var value string
    err := json.Unmarshal(data, &value)

    if err != nil {
        return err
    }

    return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *CatalogPaperSide) Scan(value any) error {
    if val, ok := value.(string); ok {
        return e.ParseAndSet(val)
    }

    return mrerr.ErrInternalTypeAssertion.New("CatalogPaperSide", value)
}
