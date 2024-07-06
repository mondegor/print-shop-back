package httpv1

type (
	// CalcRectInsideOutsideRequest - comment struct.
	CalcRectInsideOutsideRequest struct {
		InFormat  string `json:"inFormat" validate:"required,max=16,tag_double_size"`
		OutFormat string `json:"outFormat" validate:"required,max=16,tag_double_size"`
	}
)
