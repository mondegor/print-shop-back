package entity

import (
	"print-shop-back/pkg/modules/controls"
	"print-shop-back/pkg/modules/controls/enums"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameElementTemplate = "admin-api.Controls.ElementTemplate"
)

type (
	ElementTemplate struct { // DB: printshop_controls.element_templates
		ID         mrtype.KeyInt32        `json:"id"` // template_id
		TagVersion int32                  `json:"tagVersion"`
		ParamName  string                 `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption    string                 `json:"caption" sort:"caption,default" upd:"template_caption"`
		Type       enums.ElementType      `json:"elementType"`
		Detailing  enums.ElementDetailing `json:"detailing"`
		Body       []byte                 `json:"-" upd:"element_body"`
		Status     mrenum.ItemStatus      `json:"status"`
		CreatedAt  time.Time              `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time             `json:"updatedAt,omitempty" sort:"updatedAt"`
	}

	ElementTemplateHead struct {
		controls.ElementTemplateHead
		Status mrenum.ItemStatus
	}

	ElementTemplateParams struct {
		Filter ElementTemplateListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	ElementTemplateListFilter struct {
		SearchText string
		Detailing  []enums.ElementDetailing
		Statuses   []mrenum.ItemStatus
	}
)
