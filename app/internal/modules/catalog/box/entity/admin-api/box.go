package entity

import (
	"print-shop-back/pkg/libs/measure"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameBox = "admin-api.Catalog.Box"
)

type (
	Box struct { // DB: printshop_catalog.boxes
		ID         mrtype.KeyInt32 `json:"id"`                                   // box_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // created_at
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // updated_at

		Article string             `json:"article" sort:"article" upd:"box_article"`
		Caption string             `json:"caption" sort:"caption,default" upd:"box_caption"`
		Length  measure.Micrometer `json:"length" sort:"length" upd:"box_length"` // (mm)
		Width   measure.Micrometer `json:"width" sort:"width" upd:"box_width"`    // (mm)
		Depth   measure.Micrometer `json:"depth" sort:"depth" upd:"box_depth"`    // (mm)

		Status mrenum.ItemStatus `json:"status"` // box_status
	}

	BoxParams struct {
		Filter BoxListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	BoxListFilter struct {
		SearchText string
		Length     mrtype.RangeInt64
		Width      mrtype.RangeInt64
		Depth      mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
