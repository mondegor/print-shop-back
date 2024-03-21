package entity

import (
	"print-shop-back/pkg/libs/measure"
	"print-shop-back/pkg/modules/catalog/enums"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaper = "public-api.Catalog.Paper"
)

type (
	Paper struct { // DB: printshop_catalog.papers
		ID        mrtype.KeyInt32        `json:"id"` // paper_id
		Article   string                 `json:"article"`
		Caption   string                 `json:"caption"`
		ColorID   mrtype.KeyInt32        `json:"colorId"`   // paper_colors::color_id
		FactureID mrtype.KeyInt32        `json:"factureId"` // paper_factures::facture_id
		Length    measure.Micrometer     `json:"length"`    // (mm)
		Width     measure.Micrometer     `json:"width"`     // (mm)
		Density   measure.GramsPerMeter2 `json:"density"`   // (g/m2)
		Thickness measure.Micrometer     `json:"thickness"` // (mkm)
		Sides     enums.PaperSide        `json:"sides"`
	}

	PaperParams struct {
	}
)
