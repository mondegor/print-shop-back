package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogPaper = "CatalogPaper"

type (
    CatalogPaper struct { // DB: catalog_papers
        Id        mrentity.KeyInt32 `json:"id"` // paper_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created

        Article   string `json:"article" db:"paper_article"`
        Caption   string `json:"caption" db:"paper_caption"`
        Length    Micrometer `json:"length" db:"paper_length"` // (mm)
        Width     Micrometer `json:"width" db:"paper_width"` // (mm)
        Density   GramsPerMeter2 `json:"density" db:"paper_density"` // (g/m2)
        ColorId   mrentity.KeyInt32 `json:"colorId" db:"color_id"` // catalog_paper_colors::color_id
        FactureId mrentity.KeyInt32 `json:"factureId" db:"facture_id"` // catalog_paper_factures::facture_id
        Thickness Micrometer `json:"thickness" db:"paper_thickness"` // (mkm)
        Sides     CatalogPaperSide `json:"sides" db:"paper_sides"`

        Status    mrcom.ItemStatus `json:"status"` // paper_status
    }

    CatalogPaperListFilter struct {
        ColorId mrentity.KeyInt32
        FactureId mrentity.KeyInt32
        Statuses  []mrcom.ItemStatus
    }
)
