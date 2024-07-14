package entity

import (
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePrintFormat = "public-api.Dictionaries.PrintFormat" // ModelNamePrintFormat - название сущности
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct { // DB: printshop_dictionaries.print_formats
		ID      mrtype.KeyInt32 `json:"id"` // format_id
		Caption string          `json:"caption"`
	}

	// PrintFormatParams - comment struct.
	PrintFormatParams struct{}
)
