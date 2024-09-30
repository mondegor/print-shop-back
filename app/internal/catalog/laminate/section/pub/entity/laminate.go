package entity

import (
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

const (
	ModelNameLaminate = "public-api.Dictionaries.Laminate" // ModelNameLaminate - название сущности
)

type (
	// Laminate - comment struct.
	Laminate struct { // DB: printshop_catalog.laminates
		ID        mrtype.KeyInt32           `json:"id"` // laminate_id
		Article   string                    `json:"article"`
		Caption   string                    `json:"caption"`
		TypeID    mrtype.KeyInt32           `json:"typeId"` // material_types::type_id
		Length    measure.Meter             `json:"length"`
		Width     measure.Meter             `json:"width"`
		Thickness measure.Meter             `json:"thickness"`
		Weight    measure.KilogramPerMeter2 `json:"weight"`
	}

	// LaminateParams - comment struct.
	LaminateParams struct{}
)
