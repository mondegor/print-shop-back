package http_v1

import (
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
)

type (
	CreateSubmitFormRequest struct {
		ParamName string                         `json:"paramName" validate:"required,min=4,max=32,tag_variable"`
		Caption   string                         `json:"caption" validate:"required,max=128"`
		Detailing entity_shared.ElementDetailing `json:"formDetailing"`
	}

	StoreSubmitFormRequest struct {
		TagVersion int32                          `json:"tagVersion" validate:"required,gte=1"`
		ParamName  string                         `json:"paramName" validate:"omitempty,min=4,max=32,tag_variable"`
		Caption    string                         `json:"caption" validate:"omitempty,max=128"`
		Detailing  entity_shared.ElementDetailing `json:"formDetailing" validate:"omitempty"`
	}

	SubmitFormListResponse struct {
		Items []entity.SubmitForm `json:"items"`
		Total int64               `json:"total"`
	}
)
