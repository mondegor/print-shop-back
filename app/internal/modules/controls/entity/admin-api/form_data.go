package entity

import (
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameFormData = "admin-api.Controls.FormData"
)

type (
	FormData struct { // DB: ps_controls.forms
		ID         mrtype.KeyInt32 `json:"id"`                                   // form_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		ParamName string                         `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption   string                         `json:"caption" sort:"caption,default" upd:"form_caption"`
		Detailing entity_shared.ElementDetailing `json:"formDetailing" upd:"form_detailing"`

		Body   string            `json:"formBody"` // form_body_compiled
		Status mrenum.ItemStatus `json:"status"`   // form_status
	}

	FormDataParams struct {
		Filter FormDataListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	FormDataListFilter struct {
		SearchText string
		Detailing  []entity_shared.ElementDetailing
		Statuses   []mrenum.ItemStatus
	}
)
