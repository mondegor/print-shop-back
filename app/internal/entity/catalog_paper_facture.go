package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogPaperFacture = "CatalogPaperFacture"

type (
    CatalogPaperFacture struct { // DB: catalog_paper_factures
        Id        mrentity.KeyInt32 `json:"id"` // facture_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Caption   string `json:"caption"` // facture_caption
        Status    mrcom.ItemStatus `json:"status"` // facture_status
    }

    CatalogPaperFactureListFilter struct {
        Statuses  []mrcom.ItemStatus
    }
)
