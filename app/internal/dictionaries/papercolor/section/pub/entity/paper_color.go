package entity

const (
	// ModelNamePaperColor - название сущности.
	ModelNamePaperColor = "public-api.Dictionaries.PaperColor"
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
