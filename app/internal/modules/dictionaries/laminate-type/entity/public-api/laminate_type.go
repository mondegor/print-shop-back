package entity

import (
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameLaminateType = "public-api.Dictionaries.LaminateType"
)

type (
	LaminateType struct { // DB: printshop_dictionaries.laminate_types
		ID      mrtype.KeyInt32 `json:"id"` // type_id
		Caption string          `json:"caption"`
	}

	LaminateTypeParams struct {
	}
)
