package entity

import (
	"print-shop-back/pkg/modules/controls/enums"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameFormElement = "admin-api.Controls.FormElement"
)

type (
	FormElement struct { // DB: printshop_controls.submit_form_elements
		ID              mrtype.KeyInt32        `json:"id"` // element_id
		TagVersion      int32                  `json:"tagVersion"`
		FormID          uuid.UUID              `json:"-"` // submit_forms::form_id
		ParamName       string                 `json:"paramName" upd:"param_name"`
		Caption         string                 `json:"caption" upd:"element_caption"`
		TemplateID      mrtype.KeyInt32        `json:"templateId"` // element_templates::template_id
		TemplateVersion int32                  `json:"templateVersion"`
		Required        *bool                  `json:"elementRequired" upd:"element_required"`
		Type            enums.ElementType      `json:"elementType"` // element_templates::element_type
		Detailing       enums.ElementDetailing `json:"detailing"`   // element_templates::element_detailing
		Body            []byte                 `json:"-"`           // element_templates::element_body
		CreatedAt       time.Time              `json:"createdAt"`
		UpdatedAt       *time.Time             `json:"updatedAt,omitempty"`
	}
)
