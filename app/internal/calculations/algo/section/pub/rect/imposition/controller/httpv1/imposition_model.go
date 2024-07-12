package httpv1

type (
	// CalcRectImpositionRequest - comment struct.
	CalcRectImpositionRequest struct {
		ItemFormat      string `json:"itemFormat" validate:"required,max=16,tag_double_size"`
		ItemDistance    string `json:"itemDistance" validate:"max=16,tag_double_size"`
		OutFormat       string `json:"outFormat" validate:"required,max=16,tag_double_size"`
		DisableRotation bool   `json:"disableRotation"`
		UseMirror       bool   `json:"useMirror"`
	}
)
