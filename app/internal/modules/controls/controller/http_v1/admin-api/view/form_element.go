package view

import (
	entity "print-shop-back/internal/modules/controls/entity/admin-api"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CreateFormElementRequest struct {
		FormID     mrtype.KeyInt32 `json:"formId" validate:"required"`
		TemplateID mrtype.KeyInt32 `json:"templateId" validate:"required"`
		ParamName  string          `json:"paramName" validate:"omitempty,min=4,max=32,variable"`
		Caption    string          `json:"caption" validate:"omitempty,max=64"`
		Required   bool            `json:"fieldRequired"`
	}

	StoreFormElementRequest struct {
		TagVersion int32  `json:"version" validate:"required,gte=1"`
		ParamName  string `json:"paramName" validate:"omitempty,min=4,max=32,variable"`
		Caption    string `json:"caption" validate:"omitempty,max=64"`
		Required   bool   `json:"fieldRequired"`
	}

	MoveFormElementRequest struct {
		AfterNodeID mrtype.KeyInt32 `json:"afterId"`
	}

	FormElementListResponse struct {
		Items []entity.FormElement `json:"items"`
		Total int64                `json:"total"`
	}
)
