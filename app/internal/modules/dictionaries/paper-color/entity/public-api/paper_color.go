package entity

import (
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaperColor = "public-api.Dictionaries.PaperColor"
)

type (
	PaperColor struct { // DB: printshop_dictionaries.paper_colors
		ID      mrtype.KeyInt32 `json:"id"` // color_id
		Caption string          `json:"caption"`
	}

	PaperColorParams struct {
	}
)
