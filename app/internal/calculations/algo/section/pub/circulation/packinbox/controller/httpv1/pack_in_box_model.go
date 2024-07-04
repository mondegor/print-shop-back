package httpv1

type (
	// CalcCirculationPackInBoxRequest - comment struct.
	CalcCirculationPackInBoxRequest struct {
		Format  string `json:"format" validate:"required,max=16,tag_double_size"`
	}
)
