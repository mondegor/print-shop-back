package entity

import (
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameElementTemplate = "admin-api.Controls.ElementTemplate"
)

type (
	ElementTemplate struct { // DB: ps_controls.element_templates
		ID         mrtype.KeyInt32 `json:"id"`                                   // template_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		ParamName string                         `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption   string                         `json:"caption" sort:"caption,default" upd:"template_caption"`
		Type      entity_shared.ElementType      `json:"elementType" upd:"element_type"`
		Detailing entity_shared.ElementDetailing `json:"elementDetailing" upd:"element_detailing"`
		Body      string                         `json:"elementBody" upd:"element_body"`

		Status mrenum.ItemStatus `json:"status"` // template_status
	}

	ElementTemplateParams struct {
		Filter ElementTemplateListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	ElementTemplateListFilter struct {
		SearchText string
		Detailing  []entity_shared.ElementDetailing
		Statuses   []mrenum.ItemStatus
	}
)
