package view

import entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"

type (
	CreateLaminateTypeRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	StoreLaminateTypeRequest struct {
		Version int32  `json:"version" validate:"required,gte=1"`
		Caption string `json:"caption" validate:"required,max=64"`
	}

	LaminateTypeListResponse struct {
		Items []entity.LaminateType `json:"items"`
		Total int64                 `json:"total"`
	}
)
