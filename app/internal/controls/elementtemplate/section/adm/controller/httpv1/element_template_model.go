package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"
)

type (
	// CreateElementTemplateRequest - comment struct.
	CreateElementTemplateRequest struct {
		ParamName string                `json:"paramName" validate:"required,min=4,max=32,tag_variable"`
		Caption   string                `json:"caption" validate:"required,max=64"`
		Type      enum.ElementType      `json:"elementType" validate:"required"`
		Detailing enum.ElementDetailing `json:"detailing" validate:"required"`
	}

	// StoreElementTemplateRequest - comment struct.
	StoreElementTemplateRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		ParamName  string `json:"paramName" validate:"omitempty,min=4,max=32,tag_variable"`
		Caption    string `json:"caption" validate:"omitempty,max=64"`
	}

	// ElementTemplateListResponse - comment struct.
	ElementTemplateListResponse struct {
		Items []entity.ElementTemplate `json:"items"`
		Total int64                    `json:"total"`
	}
)
