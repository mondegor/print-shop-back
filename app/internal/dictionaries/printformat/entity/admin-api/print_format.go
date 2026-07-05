package entity

import (
	"time"

	"github.com/mondegor/go-core/mrtype"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/pkg/mrcalc/measure"
)

const (
	// ModelNamePrintFormat - название сущности.
	ModelNamePrintFormat = "admin-api.Dictionaries.PrintFormat"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct { // DB: printshop_dictionaries.print_format
		ID         uint64              `json:"id"` // format_id
		TagVersion uint32              `json:"tagVersion"`
		Caption    string              `json:"caption" sort:"caption,default"`
		Width      measure.Meter       `json:"width" sort:"width"`
		Height     measure.Meter       `json:"height" sort:"height"`
		Status     workflow.ItemStatus `json:"status"`
		CreatedAt  time.Time           `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time           `json:"updatedAt" sort:"updatedAt"`
	}

	// PrintFormatParams - comment struct.
	PrintFormatParams struct {
		Filter PrintFormatListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// PrintFormatListFilter - comment struct.
	PrintFormatListFilter struct {
		SearchText string
		Width      measure.RangeMeter
		Height     measure.RangeMeter
		Statuses   []workflow.ItemStatus
	}
)
