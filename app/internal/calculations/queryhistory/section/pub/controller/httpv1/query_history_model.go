package httpv1

type (
	// CreateQueryHistoryRequest - comment struct.
	CreateQueryHistoryRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
		Params  string  `json:"params" validate:"required,max=16384"`
		Result  string  `json:"result" validate:"required,max=16384"`
	}
)
