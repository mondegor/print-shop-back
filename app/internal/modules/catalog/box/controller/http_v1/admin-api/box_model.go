package http_v1

import (
	entity "print-shop-back/internal/modules/catalog/box/entity/admin-api"
	"print-shop-back/pkg/libs/measure"
)

type (
	CreateBoxRequest struct {
		Article string             `json:"article" validate:"required,min=3,max=32,tag_article"`
		Caption string             `json:"caption" validate:"required,max=64"`
		Length  measure.Micrometer `json:"length" validate:"required,gte=1,lte=10000000"`
		Width   measure.Micrometer `json:"width" validate:"required,gte=1,lte=10000000"`
		Depth   measure.Micrometer `json:"depth" validate:"required,gte=1,lte=10000000"`
	}

	StoreBoxRequest struct {
		TagVersion int32              `json:"tagVersion" validate:"required,gte=1"`
		Article    string             `json:"article" validate:"omitempty,min=3,max=32,tag_article"`
		Caption    string             `json:"caption" validate:"omitempty,max=64"`
		Length     measure.Micrometer `json:"length" validate:"omitempty,gte=1,lte=10000000"`
		Width      measure.Micrometer `json:"width" validate:"omitempty,gte=1,lte=10000000"`
		Depth      measure.Micrometer `json:"depth" validate:"omitempty,gte=1,lte=10000000"`
	}

	BoxListResponse struct {
		Items []entity.Box `json:"items"`
		Total int64        `json:"total"`
	}
)
