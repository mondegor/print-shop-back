package entity

import (
	"time"

	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtype"
)

const (
	// ModelNamePaperColor - название сущности.
	ModelNamePaperColor = "admin-api.Dictionaries.PaperColor"
)

type (
	// PaperColor - comment struct.
	PaperColor struct { // DB: printshop_dictionaries.paper_colors
		ID         uint64          `json:"id"` // color_id
		TagVersion uint32          `json:"tagVersion"`
		Caption    string          `json:"caption" sort:"caption,default"`
		Status     itemstatus.Enum `json:"status"`
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time       `json:"updatedAt" sort:"updatedAt"`
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
		Statuses   []itemstatus.Enum
	}
)
