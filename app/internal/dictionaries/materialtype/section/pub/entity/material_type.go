package entity

import (
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameMaterialType = "public-api.Dictionaries.MaterialType" // ModelNameMaterialType - название сущности
)

type (
	// MaterialType - comment struct.
	MaterialType struct { // DB: printshop_dictionaries.material_types
		ID      mrtype.KeyInt32 `json:"id"` // type_id
		Caption string          `json:"caption"`
	}

	// MaterialTypeParams - comment struct.
	MaterialTypeParams struct{}
)
