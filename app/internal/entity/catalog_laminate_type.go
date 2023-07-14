package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

type (
    CatalogLaminateType struct { // DB: catalog_laminate_types
        Id        mrentity.KeyInt32 `json:"id"` // type_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Caption   string `json:"caption"` // type_caption
        Status    ItemStatus `json:"status"` // type_status
    }

    CatalogLaminateTypeListFilter struct {
        Statuses  []ItemStatus
    }
)
