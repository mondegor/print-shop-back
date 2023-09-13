package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogLaminateType = "CatalogLaminateType"

type (
    CatalogLaminateType struct { // DB: catalog_laminate_types
        Id        mrentity.KeyInt32 `json:"id"` // type_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Caption   string `json:"caption"` // type_caption
        Status    mrcom.ItemStatus `json:"status"` // type_status
    }

    CatalogLaminateTypeListFilter struct {
        Statuses  []mrcom.ItemStatus
    }
)
