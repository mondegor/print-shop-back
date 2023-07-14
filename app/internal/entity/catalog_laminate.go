package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

type (
    CatalogLaminate struct { // DB: catalog_laminates
        Id        mrentity.KeyInt32 `json:"id"` // laminate_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Article   string `json:"article"` // laminate_article
        Caption   string `json:"caption"` // laminate_caption
        TypeId    mrentity.KeyInt32 `json:"typeId"` // catalog_laminate_types::type_id
        Length    mrentity.Micrometer `json:"length"` // laminate_length (mm)
        Weight    mrentity.GramsPerMeter2 `json:"weight"` // laminate_weight (g/m2)
        Thickness mrentity.Micrometer `json:"thickness"` // laminate_thickness (mkm)
        Status    ItemStatus `json:"status"` // laminate_status
    }

    CatalogLaminateListFilter struct {
        Statuses  []ItemStatus
    }
)
