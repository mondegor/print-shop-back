package httpv1

type (
	// CalcCuttingQuantityRequest - comment struct.
	CalcCuttingQuantityRequest struct {
		Fragments      []string `json:"fragments" validate:"required,dive,max=16,tag_double_size"`
		DistanceFormat string   `json:"distance" validate:"required,max=16,tag_double_size"`
	}
)
