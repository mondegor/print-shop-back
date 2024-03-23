package entity

const (
	ModelNameSubmitForm = "public-api.Controls.SubmitForm"
)

type (
	SubmitForm struct { // DB: printshop_controls.submit_forms
		Version     int32  `json:"version"`
		RewriteName string `json:"rewriteName"`
		Caption     string `json:"caption"`
	}

	SubmitFormParams struct {
	}
)
