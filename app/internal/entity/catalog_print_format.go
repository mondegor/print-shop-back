package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogPrintFormat = "CatalogPrintFormat"

type (
    CatalogPrintFormat struct { // DB: catalog_print_format
        Id        mrentity.KeyInt32 `json:"id"` // format_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created

        Caption   string `json:"caption" db:"format_caption"`
        Length    Micrometer `json:"length" db:"format_length"` // (mm)
        Width     Micrometer `json:"width" db:"format_width"` // (mm)

        Status    mrcom.ItemStatus `json:"status"` // format_status
    }

    CatalogPrintFormatListFilter struct {
        Statuses  []mrcom.ItemStatus
    }
)
