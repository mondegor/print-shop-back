package entity

import (
	"github.com/mondegor/print-shop-back/pkg/catalog/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaper = "public-api.Catalog.Paper" // ModelNamePaper - название сущности
)

type (
	// Paper - comment struct.
	Paper struct { // DB: printshop_catalog.papers
		ID        mrtype.KeyInt32        `json:"id"` // paper_id
		Article   string                 `json:"article"`
		Caption   string                 `json:"caption"`
		TypeID    mrtype.KeyInt32        `json:"typeId"`    // material_types::type_id
		ColorID   mrtype.KeyInt32        `json:"colorId"`   // paper_colors::color_id
		FactureID mrtype.KeyInt32        `json:"factureId"` // paper_factures::facture_id
		Length    measure.Micrometer     `json:"length"`    // mkm (mm * 1000)
		Height    measure.Micrometer     `json:"height"`    // mkm (mm * 1000)
		Thickness measure.Micrometer     `json:"thickness"` // mkm (mm * 1000)
		Density   measure.GramsPerMeter2 `json:"density"`   // g/m2
		Sides     enum.PaperSide         `json:"sides"`
	}

	// PaperParams - comment struct.
	PaperParams struct{}
)
