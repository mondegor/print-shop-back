package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

const ModelNameCatalogBox = "CatalogBox"

type (
    CatalogBox struct { // DB: catalog_boxes
        Id        mrentity.KeyInt32 `json:"id"` // box_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created

        Article   string `json:"article" db:"box_article"`
        Caption   string `json:"caption" db:"box_caption"`
        Length    mrentity.Micrometer `json:"length" db:"box_length"` // (mm)
        Width     mrentity.Micrometer `json:"width" db:"box_width"` // (mm)
        Depth     mrentity.Micrometer `json:"depth" db:"box_depth"` // (mm)

        Status    ItemStatus `json:"status"` // box_status
    }

    CatalogBoxListFilter struct {
        Statuses  []ItemStatus
    }
)
