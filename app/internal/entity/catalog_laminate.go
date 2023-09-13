package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogLaminate = "CatalogLaminate"

type (
    CatalogLaminate struct { // DB: catalog_laminates
        Id        mrentity.KeyInt32 `json:"id"` // laminate_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created

        Article   string `json:"article" db:"laminate_article"`
        Caption   string `json:"caption" db:"laminate_caption"`
        TypeId    mrentity.KeyInt32 `json:"typeId" db:"type_id"` // catalog_laminate_types::type_id
        Length    Micrometer `json:"length" db:"laminate_length"` // (mm)
        Weight    GramsPerMeter2 `json:"weight" db:"laminate_weight"` // (g/m2)
        Thickness Micrometer `json:"thickness" db:"laminate_thickness"` // (mkm)

        Status    mrcom.ItemStatus `json:"status"` // laminate_status
    }

    CatalogLaminateListFilter struct {
        Statuses  []mrcom.ItemStatus
    }
)
