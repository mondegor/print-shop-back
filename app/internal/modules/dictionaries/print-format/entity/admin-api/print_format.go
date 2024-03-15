package entity

import (
	"print-shop-back/pkg/libs/measure"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePrintFormat = "admin-api.Dictionaries.PrintFormat"
)

type (
	PrintFormat struct { // DB: printshop_dictionaries.print_format
		ID         mrtype.KeyInt32 `json:"id"`                                   // format_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // created_at
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // updated_at

		Caption string             `json:"caption" sort:"caption,default" upd:"format_caption"`
		Length  measure.Micrometer `json:"length" sort:"length" upd:"format_length"` // (mm)
		Width   measure.Micrometer `json:"width" sort:"width" upd:"format_width"`    // (mm)

		Status mrenum.ItemStatus `json:"status"` // format_status
	}

	PrintFormatParams struct {
		Filter PrintFormatListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	PrintFormatListFilter struct {
		SearchText string
		Length     mrtype.RangeInt64
		Width      mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
