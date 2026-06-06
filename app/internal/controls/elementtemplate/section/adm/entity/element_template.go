package entity

import (
	"time"

	"github.com/mondegor/go-sysmess/mrtype"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/pkg/controls/api"
	"print-shop-back/pkg/controls/enum/elementdetailing"
	"print-shop-back/pkg/controls/enum/elementtype"
)

const (
	// ModelNameElementTemplate - название сущности.
	ModelNameElementTemplate = "admin-api.Controls.ElementTemplate"
)

type (
	// ElementTemplate - comment struct.
	ElementTemplate struct { // DB: printshop_controls.element_templates
		ID         uint64                `json:"id"` // template_id
		TagVersion uint32                `json:"tagVersion"`
		ParamName  string                `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption    string                `json:"caption" sort:"caption,default" upd:"template_caption"`
		Type       elementtype.Enum      `json:"elementType"`
		Detailing  elementdetailing.Enum `json:"detailing"`
		Body       []byte                `json:"-" upd:"element_body"`
		Status     workflow.ItemStatus   `json:"status"`
		CreatedAt  time.Time             `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time             `json:"updatedAt" sort:"updatedAt"`
	}

	// ElementTemplateHead - comment struct.
	ElementTemplateHead struct {
		api.ElementTemplateDTO
		Status workflow.ItemStatus
	}

	// ElementTemplateParams - comment struct.
	ElementTemplateParams struct {
		Filter ElementTemplateListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// ElementTemplateListFilter - comment struct.
	ElementTemplateListFilter struct {
		SearchText string
		Detailing  []elementdetailing.Enum
		Statuses   []workflow.ItemStatus
	}
)
