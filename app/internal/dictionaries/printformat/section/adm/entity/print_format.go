package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

const (
	ModelNamePrintFormat = "admin-api.Dictionaries.PrintFormat" // ModelNamePrintFormat - название сущности
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct { // DB: printshop_dictionaries.print_format
		ID         uint64            `json:"id"` // format_id
		TagVersion uint32            `json:"tagVersion"`
		Caption    string            `json:"caption" sort:"caption,default"`
		Width      measure.Meter     `json:"width" sort:"width"`
		Height     measure.Meter     `json:"height" sort:"height"`
		Status     mrenum.ItemStatus `json:"status"`
		CreatedAt  time.Time         `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time         `json:"updatedAt" sort:"updatedAt"`
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
		Statuses   []mrenum.ItemStatus
	}
)
