package entity

import (
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameSubmitForm = "admin-api.Controls.SubmitForm"
)

type (
	SubmitForm struct { // DB: printshop_controls.submit_forms
		ID         mrtype.KeyInt32 `json:"id"`                                   // form_id
		TagVersion int32           `json:"tagVersion"`                           // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // created_at
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // updated_at

		ParamName string                         `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption   string                         `json:"caption" sort:"caption,default" upd:"form_caption"`
		Detailing entity_shared.ElementDetailing `json:"formDetailing" upd:"form_detailing"`

		Body   string            `json:"formBody"` // form_body_compiled
		Status mrenum.ItemStatus `json:"status"`   // form_status
	}

	SubmitFormParams struct {
		Filter SubmitFormListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	SubmitFormListFilter struct {
		SearchText string
		Detailing  []entity_shared.ElementDetailing
		Statuses   []mrenum.ItemStatus
	}
)
