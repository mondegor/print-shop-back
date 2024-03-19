package http_v1

import (
	entity "print-shop-back/internal/modules/catalog/paper/entity/admin-api"
	"print-shop-back/pkg/libs/measure"
	"print-shop-back/pkg/modules/catalog/enums"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CreatePaperRequest struct {
		Article   string                 `json:"article" validate:"required,min=3,max=32,tag_article"`
		Caption   string                 `json:"caption" validate:"required,max=64"`
		ColorID   mrtype.KeyInt32        `json:"colorId" validate:"required,gte=1"`
		FactureID mrtype.KeyInt32        `json:"factureId" validate:"required,gte=1"`
		Length    measure.Micrometer     `json:"length" validate:"required,gte=1,lte=10000000"`
		Width     measure.Micrometer     `json:"width" validate:"required,gte=1,lte=10000000"`
		Density   measure.GramsPerMeter2 `json:"density" validate:"required,gte=1,lte=10000"`
		Thickness measure.Micrometer     `json:"thickness" validate:"required,gte=1,lte=1000000"`
		Sides     enums.PaperSide        `json:"sides" validate:"required"`
	}

	StorePaperRequest struct {
		TagVersion int32                  `json:"tagVersion" validate:"required,gte=1"`
		Article    string                 `json:"article" validate:"omitempty,min=3,max=32,tag_article"`
		Caption    string                 `json:"caption" validate:"omitempty,max=64"`
		ColorID    mrtype.KeyInt32        `json:"colorId" validate:"omitempty,gte=1"`
		FactureID  mrtype.KeyInt32        `json:"factureId" validate:"omitempty,gte=1"`
		Length     measure.Micrometer     `json:"length" validate:"omitempty,gte=1,lte=10000000"`
		Width      measure.Micrometer     `json:"width" validate:"omitempty,gte=1,lte=10000000"`
		Density    measure.GramsPerMeter2 `json:"density" validate:"omitempty,gte=1,lte=10000"`
		Thickness  measure.Micrometer     `json:"thickness" validate:"omitempty,gte=1,lte=1000000"`
		Sides      enums.PaperSide        `json:"sides" validate:"omitempty"`
	}

	PaperListResponse struct {
		Items []entity.Paper `json:"items"`
		Total int64          `json:"total"`
	}
)
