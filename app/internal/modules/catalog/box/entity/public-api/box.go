package entity

import (
	"print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameBox = "public-api.Catalog.Box"
)

type (
	Box struct { // DB: printshop_catalog.boxes
		ID      mrtype.KeyInt32    `json:"id"` // box_id
		Article string             `json:"article"`
		Caption string             `json:"caption"`
		Length  measure.Micrometer `json:"length"` // (mm)
		Width   measure.Micrometer `json:"width"`  // (mm)
		Depth   measure.Micrometer `json:"depth"`  // (mm)
	}

	BoxParams struct {
	}
)
