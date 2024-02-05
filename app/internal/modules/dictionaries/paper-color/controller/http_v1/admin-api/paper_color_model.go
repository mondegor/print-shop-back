package http_v1

import entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/admin-api"

type (
	CreatePaperColorRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	StorePaperColorRequest struct {
		Version int32  `json:"version" validate:"required,gte=1"`
		Caption string `json:"caption" validate:"required,max=64"`
	}

	PaperColorListResponse struct {
		Items []entity.PaperColor `json:"items"`
		Total int64               `json:"total"`
	}
)
