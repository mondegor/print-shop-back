package entity

const (
	ModelNamePaperFacture = "public-api.Dictionaries.PaperFacture" // ModelNamePaperFacture - название сущности
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct { // DB: printshop_dictionaries.paper_factures
		ID      uint64 `json:"id"` // facture_id
		Caption string `json:"caption"`
	}

	// PaperFactureParams - comment struct.
	PaperFactureParams struct{}
)
