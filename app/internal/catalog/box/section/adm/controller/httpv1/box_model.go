package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// CreateBoxRequest - comment struct.
	CreateBoxRequest struct {
		Article string             `json:"article" validate:"required,min=3,max=32,tag_article"`
		Caption string             `json:"caption" validate:"required,max=64"`
		Length  measure.Micrometer `json:"length" validate:"required,gte=1,lte=10000000"`
		Width   measure.Micrometer `json:"width" validate:"required,gte=1,lte=10000000"`
		Height  measure.Micrometer `json:"height" validate:"required,gte=1,lte=10000000"`
		Weight  measure.Milligram  `json:"weight" validate:"required,gte=1,lte=10000000"`
	}

	// StoreBoxRequest - comment struct.
	StoreBoxRequest struct {
		TagVersion int32              `json:"tagVersion" validate:"required,gte=1"`
		Article    string             `json:"article" validate:"omitempty,min=3,max=32,tag_article"`
		Caption    string             `json:"caption" validate:"omitempty,max=64"`
		Length     measure.Micrometer `json:"length" validate:"omitempty,gte=1,lte=10000000"`
		Width      measure.Micrometer `json:"width" validate:"omitempty,gte=1,lte=10000000"`
		Height     measure.Micrometer `json:"height" validate:"omitempty,gte=1,lte=10000000"`
		Weight     measure.Milligram  `json:"weight" validate:"required,gte=1,lte=10000000"`
	}

	// BoxListResponse - comment struct.
	BoxListResponse struct {
		Items []entity.Box `json:"items"`
		Total int64        `json:"total"`
	}
)
