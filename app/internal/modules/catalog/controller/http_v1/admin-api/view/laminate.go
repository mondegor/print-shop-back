package view

import (
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	"print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CreateLaminateRequest struct {
		Article   string                 `json:"article" validate:"required,min=3,max=32,article"`
		Caption   string                 `json:"caption" validate:"required,max=64"`
		TypeID    mrtype.KeyInt32        `json:"typeId" validate:"required,gte=1"`
		Length    measure.Micrometer     `json:"length" validate:"required,gte=1,lte=1000000000"`
		Weight    measure.GramsPerMeter2 `json:"weight" validate:"required,gte=1,lte=10000"`
		Thickness measure.Micrometer     `json:"thickness" validate:"required,gte=1,lte=1000000"`
	}

	StoreLaminateRequest struct {
		Version   int32                  `json:"version" validate:"required,gte=1"`
		Article   string                 `json:"article" validate:"omitempty,min=3,max=32,article"`
		Caption   string                 `json:"caption" validate:"omitempty,max=64"`
		TypeID    mrtype.KeyInt32        `json:"typeId" validate:"omitempty,gte=1"`
		Length    measure.Micrometer     `json:"length" validate:"omitempty,gte=1,lte=1000000000"`
		Weight    measure.GramsPerMeter2 `json:"weight" validate:"omitempty,gte=1,lte=10000"`
		Thickness measure.Micrometer     `json:"thickness" validate:"omitempty,gte=1,lte=1000000"`
	}

	LaminateListResponse struct {
		Items []entity.Laminate `json:"items"`
		Total int64             `json:"total"`
	}
)
