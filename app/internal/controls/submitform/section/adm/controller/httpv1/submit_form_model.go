package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"
)

type (
	// CreateSubmitFormRequest - comment struct.
	CreateSubmitFormRequest struct {
		RewriteName string                `json:"rewriteName" validate:"required,min=4,max=32,tag_rewrite_name"`
		ParamName   string                `json:"paramName" validate:"required,min=4,max=32,tag_variable"`
		Caption     string                `json:"caption" validate:"required,max=128"`
		Detailing   enum.ElementDetailing `json:"detailing" validate:"required"`
	}

	// StoreSubmitFormRequest - comment struct.
	StoreSubmitFormRequest struct {
		TagVersion  uint32 `json:"tagVersion" validate:"required,gte=1"`
		RewriteName string `json:"rewriteName" validate:"omitempty,min=4,max=32,tag_rewrite_name"`
		ParamName   string `json:"paramName" validate:"omitempty,min=4,max=32,tag_variable"`
		Caption     string `json:"caption" validate:"omitempty,max=128"`
	}

	// SubmitFormListResponse - comment struct.
	SubmitFormListResponse struct {
		Items []entity.SubmitForm `json:"items"`
		Total uint64              `json:"total"`
	}
)
