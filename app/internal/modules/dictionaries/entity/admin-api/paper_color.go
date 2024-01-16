package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaperColor = "admin-api.Dictionaries.PaperColor"
)

type (
	PaperColor struct { // DB: ps_dictionaries.paper_colors
		ID         mrtype.KeyInt32 `json:"id"`                                   // color_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		Caption string `json:"caption" sort:"caption,default"` // color_caption

		Status mrenum.ItemStatus `json:"status"` // color_status
	}

	PaperColorParams struct {
		Filter PaperColorListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	PaperColorListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
