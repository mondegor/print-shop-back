package entity

import (
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaperColor = "public-api.Dictionaries.PaperColor" // ModelNamePaperColor - название сущности
)

type (
	// PaperColor - comment struct.
	PaperColor struct { // DB: printshop_dictionaries.paper_colors
		ID      mrtype.KeyInt32 `json:"id"` // color_id
		Caption string          `json:"caption"`
	}

	// PaperColorParams - comment struct.
	PaperColorParams struct{}
)
