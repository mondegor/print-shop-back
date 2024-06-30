package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/catalog/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// CreatePaperRequest - comment struct.
	CreatePaperRequest struct {
		Article   string                 `json:"article" validate:"required,min=3,max=32,tag_article"`
		Caption   string                 `json:"caption" validate:"required,max=64"`
		TypeID    mrtype.KeyInt32        `json:"typeId" validate:"required,gte=1"`
		ColorID   mrtype.KeyInt32        `json:"colorId" validate:"required,gte=1"`
		FactureID mrtype.KeyInt32        `json:"factureId" validate:"required,gte=1"`
		Length    measure.Micrometer     `json:"length" validate:"required,gte=1,lte=10000000"`
		Height    measure.Micrometer     `json:"height" validate:"required,gte=1,lte=10000000"`
		Thickness measure.Micrometer     `json:"thickness" validate:"required,gte=1,lte=1000000"`
		Density   measure.GramsPerMeter2 `json:"density" validate:"required,gte=1,lte=10000"`
		Sides     enum.PaperSide         `json:"sides" validate:"required"`
	}

	// StorePaperRequest - comment struct.
	StorePaperRequest struct {
		TagVersion int32                  `json:"tagVersion" validate:"required,gte=1"`
		Article    string                 `json:"article" validate:"omitempty,min=3,max=32,tag_article"`
		Caption    string                 `json:"caption" validate:"omitempty,max=64"`
		TypeID     mrtype.KeyInt32        `json:"typeId" validate:"omitempty,gte=1"`
		ColorID    mrtype.KeyInt32        `json:"colorId" validate:"omitempty,gte=1"`
		FactureID  mrtype.KeyInt32        `json:"factureId" validate:"omitempty,gte=1"`
		Length     measure.Micrometer     `json:"length" validate:"omitempty,gte=1,lte=10000000"`
		Height     measure.Micrometer     `json:"height" validate:"omitempty,gte=1,lte=10000000"`
		Thickness  measure.Micrometer     `json:"thickness" validate:"omitempty,gte=1,lte=1000000"`
		Density    measure.GramsPerMeter2 `json:"density" validate:"omitempty,gte=1,lte=10000"`
		Sides      enum.PaperSide         `json:"sides" validate:"omitempty"`
	}

	// PaperListResponse - comment struct.
	PaperListResponse struct {
		Items []entity.Paper `json:"items"`
		Total int64          `json:"total"`
	}
)
