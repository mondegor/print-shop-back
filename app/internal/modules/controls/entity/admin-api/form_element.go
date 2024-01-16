package entity

import (
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
	"time"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameFormElement = "admin-api.Controls.FormElement"
)

type (
	FormElement struct { // DB: ps_controls.form_elements
		ID         mrtype.KeyInt32 `json:"id"`                                   // element_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		FormID     mrtype.KeyInt32 `json:"formId"` // form_id
		ParamName  string          `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption    string          `json:"caption" sort:"caption,default" upd:"element_caption"`
		TemplateID mrtype.KeyInt32 `json:"templateId"` // template_id
		Required   bool            `json:"elementRequired"`

		Type      entity_shared.ElementType      `json:"elementType"`      // element_templates::element_type
		Detailing entity_shared.ElementDetailing `json:"elementDetailing"` // element_templates::element_detailing
		Body      string                         `json:"elementBody"`      // element_templates::element_body
	}

	FormElementParams struct {
		FormID mrtype.KeyInt32
		Filter FormElementListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	FormElementListFilter struct {
		SearchText string
		Detailing  []entity_shared.ElementDetailing
	}
)
