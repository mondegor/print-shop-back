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
		ID      mrtype.KeyInt32    `json:"id"` // box_id
		Article string             `json:"article"`
		Caption string             `json:"caption"`
		Length  measure.Micrometer `json:"length"` // mkm (mm * 1000)
		Width   measure.Micrometer `json:"width"`  // mkm (mm * 1000)
		Height  measure.Micrometer `json:"height"` // mkm (mm * 1000)
		Weight  measure.Milligram  `json:"weight"` // mg (g * 1000)
	}

	// BoxParams - comment struct.
	BoxParams struct{}
)
