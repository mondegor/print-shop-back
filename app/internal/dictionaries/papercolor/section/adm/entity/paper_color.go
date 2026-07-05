package entity

import (
	"time"

	"github.com/mondegor/go-core/mrtype"

	"print-shop-back/internal/adapter/workflow"
)

const (
	// ModelNamePaperColor - название сущности.
	ModelNamePaperColor = "admin-api.Dictionaries.PaperColor"
)

type (
	// PaperColor - comment struct.
	PaperColor struct { // DB: printshop_dictionaries.paper_colors
		ID         uint64              `json:"id"` // color_id
		TagVersion uint32              `json:"tagVersion"`
		Caption    string              `json:"caption" sort:"caption,default"`
		Status     workflow.ItemStatus `json:"status"`
		CreatedAt  time.Time           `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time           `json:"updatedAt" sort:"updatedAt"`
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
		Statuses   []workflow.ItemStatus
	}
)
