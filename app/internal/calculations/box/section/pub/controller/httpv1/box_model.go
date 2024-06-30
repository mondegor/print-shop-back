package httpv1

type (
	// CreateCalcResultRequest - comment struct.
	CreateCalcResultRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	// StoreCalcResultRequest - comment struct.
	StoreCalcResultRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"omitempty,max=64"`
	}
)
