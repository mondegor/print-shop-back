package httpv1

type (
	// CalcCirculationPackInBoxRequest - comment struct.
	CalcCirculationPackInBoxRequest struct {
		Product Product `json:"product" validate:"required"`
		Box     Box     `json:"box" validate:"required"`
	}

	// Product - comment struct.
	Product struct {
		Format    string `json:"format" validate:"required,max=16,tag_double_size"` // mm x mm
		Thickness uint64 `json:"thickness" validate:"required,gte=1,lte=1000000"`   // mkm
		WeightM2  uint64 `json:"weightM2"`                                          // g/m2
		Quantity  uint64 `json:"quantity" validate:"required,gte=1,lte=1000000000"`
	}

	// Box - comment struct.
	Box struct {
		Format    string `json:"format" validate:"required,max=24,tag_triple_size"`  // mm x mm x mm
		Thickness uint64 `json:"thickness" validate:"required,gte=1,lte=1000000"`    // mkm
		Margins   string `json:"margins" validate:"required,max=16,tag_triple_size"` // mm x mm x mm
		Weight    uint64 `json:"weight" validate:"required,gte=1,lte=1000000"`       // g
		MaxWeight uint64 `json:"maxWeight" validate:"gte=1,lte=1000000"`             // g
	}
)
