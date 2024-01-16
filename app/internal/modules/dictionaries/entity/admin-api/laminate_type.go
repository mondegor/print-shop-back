package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameLaminateType = "admin-api.Dictionaries.LaminateType"
)

type (
	LaminateType struct { // DB: ps_dictionaries.laminate_types
		ID         mrtype.KeyInt32 `json:"id"`                                   // type_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		Caption string `json:"caption" sort:"caption,default"` // type_caption

		Status mrenum.ItemStatus `json:"status"` // type_status
	}

	LaminateTypeParams struct {
		Filter LaminateTypeListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	LaminateTypeListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
