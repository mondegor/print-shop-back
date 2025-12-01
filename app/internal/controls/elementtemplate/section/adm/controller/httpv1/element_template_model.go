package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/type/elementdetailing"
	"github.com/mondegor/print-shop-back/pkg/controls/type/elementtype"
)

type (
	// CreateElementTemplateRequest - comment struct.
	CreateElementTemplateRequest struct {
		ParamName string                `json:"paramName" validate:"required,min=4,max=32,tag_variable"`
		Caption   string                `json:"caption" validate:"required,max=64"`
		Type      elementtype.Enum      `json:"elementType" validate:"required"`
		Detailing elementdetailing.Enum `json:"detailing" validate:"required"`
	}

	// StoreElementTemplateRequest - comment struct.
	StoreElementTemplateRequest struct {
		TagVersion uint32 `json:"tagVersion" validate:"required,gte=1"`
		ParamName  string `json:"paramName" validate:"omitempty,min=4,max=32,tag_variable"`
		Caption    string `json:"caption" validate:"omitempty,max=64"`
	}

	// ElementTemplateListResponse - comment struct.
	ElementTemplateListResponse struct {
		Items []entity.ElementTemplate `json:"items"`
		Total uint64                   `json:"total"`
	}
)
