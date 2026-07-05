package entity

import (
	"print-shop-back/pkg/catalog/type/paperside"
	"print-shop-back/pkg/mrcalc/measure"
)

const (
	// ModelNamePaper - название сущности.
	ModelNamePaper = "public-api.Catalog.Paper"
)

type (
	// Paper - comment struct.
	Paper struct { // DB: printshop_catalog.papers
		ID        uint64                    `json:"id"` // paper_id
		Article   string                    `json:"article"`
		Caption   string                    `json:"caption"`
		TypeID    uint64                    `json:"typeId"`    // material_types::type_id
		ColorID   uint64                    `json:"colorId"`   // paper_colors::color_id
		FactureID uint64                    `json:"factureId"` // paper_factures::facture_id
		Width     measure.Meter             `json:"width"`
		Height    measure.Meter             `json:"height"`
		Thickness measure.Meter             `json:"thickness"`
		Density   measure.KilogramPerMeter2 `json:"density"`
		Sides     paperside.Enum            `json:"sides"`
	}

	// PaperParams - comment struct.
	PaperParams struct{}
)
