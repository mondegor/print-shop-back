package view

import entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"

type (
	CreatePaperFactureRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	StorePaperFactureRequest struct {
		Version int32  `json:"version" validate:"required,gte=1"`
		Caption string `json:"caption" validate:"required,max=64"`
	}

	PaperFactureListResponse struct {
		Items []entity.PaperFacture `json:"items"`
		Total int64                 `json:"total"`
	}
)
