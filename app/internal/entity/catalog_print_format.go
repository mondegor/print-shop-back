package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

const ModelNameCatalogPrintFormat = "CatalogPrintFormat"

type (
    CatalogPrintFormat struct { // DB: catalog_print_format
        Id        mrentity.KeyInt32 `json:"id"` // format_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Caption   string `json:"caption"` // format_caption
        Length    mrentity.Micrometer `json:"length"` // format_length (mm)
        Width     mrentity.Micrometer `json:"width"` // format_width (mm)
        Status    ItemStatus `json:"status"` // format_status
    }

    CatalogPrintFormatListFilter struct {
        Statuses  []ItemStatus
    }
)
