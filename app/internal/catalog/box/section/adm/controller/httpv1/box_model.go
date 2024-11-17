package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// CreateBoxRequest - comment struct.
	CreateBoxRequest struct {
		Article   string             `json:"article" validate:"required,min=3,max=32,tag_article"`
		Caption   string             `json:"caption" validate:"required,max=64"`
		Length    measure.Millimeter `json:"length" validate:"required,gte=1,lte=10000"`
		Width     measure.Millimeter `json:"width" validate:"required,gte=1,lte=10000"`
		Height    measure.Millimeter `json:"height" validate:"required,gte=1,lte=10000"`
		Thickness measure.Micrometer `json:"thickness" validate:"required,gte=1,lte=10000"`
		Weight    measure.Gram       `json:"weight" validate:"required,gte=1,lte=10000"`
	}

	// StoreBoxRequest - comment struct.
	StoreBoxRequest struct {
		TagVersion uint32             `json:"tagVersion" validate:"required,gte=1"`
		Article    string             `json:"article" validate:"omitempty,min=3,max=32,tag_article"`
		Caption    string             `json:"caption" validate:"omitempty,max=64"`
		Length     measure.Millimeter `json:"length" validate:"omitempty,gte=1,lte=10000"`
		Width      measure.Millimeter `json:"width" validate:"omitempty,gte=1,lte=10000"`
		Height     measure.Millimeter `json:"height" validate:"omitempty,gte=1,lte=10000"`
		Thickness  measure.Micrometer `json:"thickness" validate:"omitempty,gte=1,lte=10000"`
		Weight     measure.Gram       `json:"weight" validate:"required,gte=1,lte=10000"`
	}

	// BoxListResponse - comment struct.
	BoxListResponse struct {
		Items []entity.Box `json:"items"`
		Total uint64       `json:"total"`
	}
)
