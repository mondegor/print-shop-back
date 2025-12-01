package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/print-shop-back/pkg/controls/type/elementdetailing"
)

const (
	// ModelNameSubmitForm - название сущности.
	ModelNameSubmitForm = "admin-api.Controls.SubmitForm"
)

type (
	// SubmitForm - comment struct.
	SubmitForm struct { // DB: printshop_controls.submit_forms
		ID          uuid.UUID             `json:"id"` // form_id
		TagVersion  uint32                `json:"tagVersion"`
		RewriteName string                `json:"rewriteName" sort:"rewriteName" upd:"rewrite_name"`
		ParamName   string                `json:"paramName" sort:"paramName" upd:"param_name"`
		Caption     string                `json:"caption" sort:"caption,default" upd:"form_caption"`
		Detailing   elementdetailing.Enum `json:"detailing"`
		Status      itemstatus.Enum       `json:"status"`
		CreatedAt   time.Time             `json:"createdAt" sort:"createdAt"`
		UpdatedAt   time.Time             `json:"updatedAt" sort:"updatedAt"`

		Elements []FormElement `json:"elements,omitempty"`
		Versions []FormVersion `json:"versions,omitempty"`
	}

	// SubmitFormParams - comment struct.
	SubmitFormParams struct {
		Filter SubmitFormListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// SubmitFormListFilter - comment struct.
	SubmitFormListFilter struct {
		SearchText string
		Detailing  []elementdetailing.Enum
		Statuses   []itemstatus.Enum
	}
)
