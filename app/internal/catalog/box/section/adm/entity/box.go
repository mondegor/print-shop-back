package entity

import (
	"time"

	"github.com/mondegor/go-sysmess/mrtype"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/pkg/mrcalc/measure"
)

const (
	// ModelNameBox - название сущности.
	ModelNameBox = "admin-api.Catalog.Box"
)

type (
	// Box - comment struct.
	Box struct { // DB: printshop_catalog.boxes
		ID         uint64              `json:"id"` // box_id
		TagVersion uint32              `json:"tagVersion"`
		Article    string              `json:"article" sort:"article" upd:"box_article"`
		Caption    string              `json:"caption" sort:"caption,default" upd:"box_caption"`
		Length     measure.Meter       `json:"length" sort:"length" upd:"box_length"`
		Width      measure.Meter       `json:"width" sort:"width" upd:"box_width"`
		Height     measure.Meter       `json:"height" sort:"height" upd:"box_height"`
		Thickness  measure.Meter       `json:"thickness" sort:"thickness" upd:"box_thickness"`
		Weight     measure.Kilogram    `json:"weight" sort:"weight" upd:"box_weight"`
		Status     workflow.ItemStatus `json:"status"`
		CreatedAt  time.Time           `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time           `json:"updatedAt" sort:"updatedAt"`
	}

	// BoxParams - comment struct.
	BoxParams struct {
		Filter BoxListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// BoxListFilter - comment struct.
	BoxListFilter struct {
		SearchText string
		Length     measure.RangeMeter
		Width      measure.RangeMeter
		Height     measure.RangeMeter
		Weight     measure.RangeKilogram
		Statuses   []workflow.ItemStatus
	}
)
