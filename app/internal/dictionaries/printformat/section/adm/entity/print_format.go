package entity

import (
	"time"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePrintFormat = "admin-api.Dictionaries.PrintFormat" // ModelNamePrintFormat - название сущности
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct { // DB: printshop_dictionaries.print_format
		ID         mrtype.KeyInt32    `json:"id"` // format_id
		TagVersion int32              `json:"tagVersion"`
		Caption    string             `json:"caption" sort:"caption,default"`
		Width      measure.Micrometer `json:"width" sort:"width"`   // mkm (mm * 1000)
		Height     measure.Micrometer `json:"height" sort:"height"` // mkm (mm * 1000)
		Status     mrenum.ItemStatus  `json:"status"`
		CreatedAt  time.Time          `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time         `json:"updatedAt,omitempty" sort:"updatedAt"`
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
		Width      mrtype.RangeInt64
		Height     mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
