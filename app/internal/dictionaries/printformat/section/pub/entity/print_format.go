package entity

const (
	ModelNamePrintFormat = "public-api.Dictionaries.PrintFormat" // ModelNamePrintFormat - название сущности
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct { // DB: printshop_dictionaries.print_formats
		ID      uint64 `json:"id"` // format_id
		Caption string `json:"caption"`
		Width   uint64 `json:"width"`
		Height  uint64 `json:"height"`
	}

	// PrintFormatParams - comment struct.
	PrintFormatParams struct{}
)
