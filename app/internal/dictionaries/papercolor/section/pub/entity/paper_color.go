package entity

const (
	ModelNamePaperColor = "public-api.Dictionaries.PaperColor" // ModelNamePaperColor - название сущности
)

type (
	// PaperColor - comment struct.
	PaperColor struct { // DB: printshop_dictionaries.paper_colors
		ID      uint64 `json:"id"` // color_id
		Caption string `json:"caption"`
	}

	// PaperColorParams - comment struct.
	PaperColorParams struct{}
)
