package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// CreateLaminateRequest - comment struct.
	CreateLaminateRequest struct {
		Article   string                `json:"article" validate:"required,min=3,max=32,tag_article"`
		Caption   string                `json:"caption" validate:"required,max=64"`
		TypeID    uint64                `json:"typeId" validate:"required,gte=1"`
		Length    measure.Meter         `json:"length" validate:"required,gte=1,lte=1000"`
		Width     measure.Millimeter    `json:"width" validate:"required,gte=1,lte=10000"`
		Thickness measure.Micrometer    `json:"thickness" validate:"required,gte=1,lte=10000"`
		WeightM2  measure.GramPerMeter2 `json:"weightM2" validate:"required,gte=1,lte=10000"`
	}

	// StoreLaminateRequest - comment struct.
	StoreLaminateRequest struct {
		TagVersion uint32                `json:"tagVersion" validate:"required,gte=1"`
		Article    string                `json:"article" validate:"omitempty,min=3,max=32,tag_article"`
		Caption    string                `json:"caption" validate:"omitempty,max=64"`
		TypeID     uint64                `json:"typeId" validate:"omitempty,gte=1"`
		Length     measure.Meter         `json:"length" validate:"omitempty,gte=1,lte=1000"`
		Width      measure.Millimeter    `json:"width" validate:"omitempty,gte=1,lte=10000"`
		Thickness  measure.Micrometer    `json:"thickness" validate:"omitempty,gte=1,lte=10000"`
		WeightM2   measure.GramPerMeter2 `json:"weightM2" validate:"omitempty,gte=1,lte=10000"`
	}

	// LaminateListResponse - comment struct.
	LaminateListResponse struct {
		Items []entity.Laminate `json:"items"`
		Total uint64            `json:"total"`
	}
)
