package view

import (
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
)

type (
	CreateFormDataRequest struct {
		ParamName string                         `json:"paramName" validate:"required,min=4,max=32,tag_variable"`
		Caption   string                         `json:"caption" validate:"required,max=128"`
		Detailing entity_shared.ElementDetailing `json:"formDetailing"`
	}

	StoreFormDataRequest struct {
		TagVersion int32                          `json:"version" validate:"required,gte=1"`
		ParamName  string                         `json:"paramName" validate:"omitempty,min=4,max=32,tag_variable"`
		Caption    string                         `json:"caption" validate:"omitempty,max=128"`
		Detailing  entity_shared.ElementDetailing `json:"formDetailing" validate:"omitempty"`
	}

	FormDataListResponse struct {
		Items []entity.FormData `json:"items"`
		Total int64             `json:"total"`
	}
)
