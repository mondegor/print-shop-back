package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

type (
    CatalogPaper struct { // DB: catalog_papers
        Id        mrentity.KeyInt32 `json:"id"` // paper_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Article   string `json:"article"` // paper_article
        Caption   string `json:"caption"` // paper_caption
        Length    mrentity.Micrometer `json:"length"` // paper_length (mm)
        Width     mrentity.Micrometer `json:"width"` // paper_width (mm)
        Density   mrentity.GramsPerMeter2 `json:"density"` // paper_density (g/m2)
        ColorId   mrentity.KeyInt32 `json:"colorId"` // catalog_paper_colors::color_id
        FactureId mrentity.KeyInt32 `json:"factureId"` // catalog_paper_factures::facture_id
        Thickness mrentity.Micrometer `json:"thickness"` // paper_thickness (mkm)
        Sides     CatalogPaperSide `json:"sides"` // paper_sides
        Status    ItemStatus `json:"status"` // paper_status
    }

    CatalogPaperListFilter struct {
        ColorId mrentity.KeyInt32
        FactureId mrentity.KeyInt32
        Statuses  []ItemStatus
    }
)
