package entity

import (
	"print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameLaminate = "public-api.Dictionaries.Laminate"
)

type (
	Laminate struct { // DB: printshop_catalog.laminates
		ID        mrtype.KeyInt32        `json:"id"` // laminate_id
		Article   string                 `json:"article"`
		Caption   string                 `json:"caption"`
		TypeID    mrtype.KeyInt32        `json:"typeId"`    // laminate_types::type_id
		Length    measure.Micrometer     `json:"length"`    // (mm)
		Weight    measure.GramsPerMeter2 `json:"weight"`    // (g/m2)
		Thickness measure.Micrometer     `json:"thickness"` // (mkm)
	}

	LaminateParams struct {
	}
)
