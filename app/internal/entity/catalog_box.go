package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogBox = "CatalogBox"

type (
    CatalogBox struct { // DB: catalog_boxes
        Id        mrentity.KeyInt32 `json:"id"` // box_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created

        Article   string `json:"article" db:"box_article"`
        Caption   string `json:"caption" db:"box_caption"`
        Length    Micrometer `json:"length" db:"box_length"` // (mm)
        Width     Micrometer `json:"width" db:"box_width"` // (mm)
        Depth     Micrometer `json:"depth" db:"box_depth"` // (mm)

        Status    mrcom.ItemStatus `json:"status"` // box_status
    }

    CatalogBoxListFilter struct {
        Statuses  []mrcom.ItemStatus
    }
)
