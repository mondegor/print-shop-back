package httpv1

import (
	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// CreateFormElementRequest - comment struct.
	CreateFormElementRequest struct {
		FormID     uuid.UUID       `json:"formId" validate:"required,min=16,max=64"`
		TemplateID mrtype.KeyInt32 `json:"templateId" validate:"required,gte=1"`
		ParamName  string          `json:"paramName" validate:"omitempty,min=4,max=32,tag_variable"`
		Caption    string          `json:"caption" validate:"omitempty,max=64"`
		Required   bool            `json:"elementRequired"`
	}

	// StoreFormElementRequest - comment struct.
	StoreFormElementRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		ParamName  string `json:"paramName" validate:"omitempty,min=4,max=32,tag_variable"`
		Caption    string `json:"caption" validate:"omitempty,max=64"`
		Required   *bool  `json:"elementRequired" validate:"omitempty"`
	}
)
