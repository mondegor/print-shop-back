package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameLaminate = "public-api.Dictionaries.Laminate" // ModelNameLaminate - название сущности
)

type (
	// Laminate - comment struct.
	Laminate struct { // DB: printshop_catalog.laminates
		ID        mrtype.KeyInt32        `json:"id"` // laminate_id
		Article   string                 `json:"article"`
		Caption   string                 `json:"caption"`
		TypeID    mrtype.KeyInt32        `json:"typeId"`    // material_types::type_id
		Length    measure.Micrometer     `json:"length"`    // mkm (mm * 1000)
		Width     measure.Micrometer     `json:"width"`     // mkm (mm * 1000)
		Thickness measure.Micrometer     `json:"thickness"` // mkm (mm * 1000)
		Weight    measure.GramsPerMeter2 `json:"weight"`    // g/m2
	}

	// LaminateParams - comment struct.
	LaminateParams struct{}
)
