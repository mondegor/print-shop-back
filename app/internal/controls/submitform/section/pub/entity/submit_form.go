package entity

const (
	ModelNameSubmitForm = "public-api.Controls.SubmitForm" // ModelNameSubmitForm - название сущности
)

type (
	// SubmitForm - comment struct.
	SubmitForm struct { // DB: printshop_controls.submit_forms
		Version     int32  `json:"version"`
		RewriteName string `json:"rewriteName"`
		Caption     string `json:"caption"`
	}

	// SubmitFormParams - comment struct.
	SubmitFormParams struct{}
)
