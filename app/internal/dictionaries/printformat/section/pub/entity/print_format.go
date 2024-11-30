package entity

import "github.com/mondegor/print-shop-back/pkg/libs/measure"

const (
	ModelNamePrintFormat = "public-api.Dictionaries.PrintFormat" // ModelNamePrintFormat - название сущности
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct { // DB: printshop_dictionaries.print_formats
		ID      uint64        `json:"id"` // format_id
		Caption string        `json:"caption"`
		Width   measure.Meter `json:"width"`
		Height  measure.Meter `json:"height"`
	}

	// PrintFormatParams - comment struct.
	PrintFormatParams struct{}
)
