package entity

import (
	"time"

	"github.com/mondegor/print-shop-back/pkg/controls/api"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameElementTemplate = "admin-api.Controls.ElementTemplate" // ModelNameElementTemplate - название сущности
)

type (
	// ElementTemplate - comment struct.
	ElementTemplate struct { // DB: printshop_controls.element_templates
		ID         mrtype.KeyInt32       `json:"id"` // template_id
		TagVersion int32                 `json:"tagVersion"`
		ParamName  string                `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption    string                `json:"caption" sort:"caption,default" upd:"template_caption"`
		Type       enum.ElementType      `json:"elementType"`
		Detailing  enum.ElementDetailing `json:"detailing"`
		Body       []byte                `json:"-" upd:"element_body"`
		Status     mrenum.ItemStatus     `json:"status"`
		CreatedAt  time.Time             `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time            `json:"updatedAt,omitempty" sort:"updatedAt"`
	}

	// ElementTemplateHead - comment struct.
	ElementTemplateHead struct {
		api.ElementTemplateDTO
		Status mrenum.ItemStatus
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
		Detailing  []enum.ElementDetailing
		Statuses   []mrenum.ItemStatus
	}
)
