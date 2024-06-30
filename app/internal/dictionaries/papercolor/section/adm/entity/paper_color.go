package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaperColor = "admin-api.Dictionaries.PaperColor" // ModelNamePaperColor - название сущности
)

type (
	// PaperColor - comment struct.
	PaperColor struct { // DB: printshop_dictionaries.paper_colors
		ID         mrtype.KeyInt32   `json:"id"` // color_id
		TagVersion int32             `json:"tagVersion"`
		Caption    string            `json:"caption" sort:"caption,default"`
		Status     mrenum.ItemStatus `json:"status"`
		CreatedAt  time.Time         `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time        `json:"updatedAt,omitempty" sort:"updatedAt"`
	}

	// PaperColorParams - comment struct.
	PaperColorParams struct {
		Filter PaperColorListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// PaperColorListFilter - comment struct.
	PaperColorListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
