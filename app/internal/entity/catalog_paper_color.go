package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogPaperColor = "CatalogPaperColor"

type (
    CatalogPaperColor struct { // DB: catalog_paper_colors
        Id        mrentity.KeyInt32 `json:"id"` // color_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Caption   string `json:"caption"` // color_caption
        Status    mrcom.ItemStatus `json:"status"` // color_status
    }

    CatalogPaperColorListFilter struct {
        Statuses  []mrcom.ItemStatus
    }
)
