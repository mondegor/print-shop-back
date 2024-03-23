package entity

import (
	"print-shop-back/pkg/modules/controls/enums"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameSubmitForm = "admin-api.Controls.SubmitForm"
)

type (
	SubmitForm struct { // DB: printshop_controls.submit_forms
		ID          uuid.UUID              `json:"id"` // form_id
		TagVersion  int32                  `json:"tagVersion"`
		RewriteName string                 `json:"rewriteName" sort:"rewriteName" upd:"rewrite_name"`
		ParamName   string                 `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption     string                 `json:"caption" sort:"caption,default" upd:"form_caption"`
		Detailing   enums.ElementDetailing `json:"detailing"`
		Status      mrenum.ItemStatus      `json:"status"`
		CreatedAt   time.Time              `json:"createdAt" sort:"createdAt"`
		UpdatedAt   *time.Time             `json:"updatedAt,omitempty" sort:"updatedAt"`

		Elements []FormElement `json:"elements,omitempty"`
		Versions []FormVersion `json:"versions,omitempty"`
	}

	SubmitFormParams struct {
		Filter SubmitFormListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	SubmitFormListFilter struct {
		SearchText string
		Detailing  []enums.ElementDetailing
		Statuses   []mrenum.ItemStatus
	}
)
