package entity

const (
	// ModelNamePrintFormat - название сущности.
	ModelNamePrintFormat = "public-api.Dictionaries.PrintFormat"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct { // DB: printshop_dictionaries.print_formats
		ID      uint64 `json:"id"` // format_id
		Caption string `json:"caption"`
	}

	// PrintFormatParams - comment struct.
	PrintFormatParams struct{}
)
