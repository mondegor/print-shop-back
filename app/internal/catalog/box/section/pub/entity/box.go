package entity

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameBox = "public-api.Catalog.Box" // ModelNameBox - название сущности
)

type (
	// Box - comment struct.
	Box struct { // DB: printshop_catalog.boxes
		ID      mrtype.KeyInt32  `json:"id"` // box_id
		Article string           `json:"article"`
		Caption string           `json:"caption"`
		Length  measure.Meter    `json:"length"`
		Width   measure.Meter    `json:"width"`
		Height  measure.Meter    `json:"height"`
		Weight  measure.Kilogram `json:"weight"`
	}

	// BoxParams - comment struct.
	BoxParams struct{}
)
