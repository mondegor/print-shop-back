package model

type (
	// SheetCuttingQuantityRequest - входные параметры алгоритма SheetCuttingQuantity.
	SheetCuttingQuantityRequest struct {
		Layouts        []string `json:"layouts" validate:"required,dive,max=16,tag_2d_size"` // mm x mm
		DistanceFormat string   `json:"distance" validate:"required,max=16,tag_2d_size"`     // mm x mm
	}

	// SheetCuttingQuantityResult - результат работы алгоритма SheetCuttingQuantity.
	SheetCuttingQuantityResult struct {
		Quantity uint64 `json:"quantity"`
	}
)
