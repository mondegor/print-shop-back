package entity

import (
	"time"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameLaminate = "admin-api.Catalog.Laminate" // ModelNameLaminate - название сущности
)

type (
	// Laminate - comment struct.
	Laminate struct { // DB: printshop_catalog.laminates
		ID         mrtype.KeyInt32        `json:"id"` // laminate_id
		TagVersion int32                  `json:"tagVersion"`
		Article    string                 `json:"article" sort:"article" upd:"laminate_article"`
		Caption    string                 `json:"caption" sort:"caption,default" upd:"laminate_caption"`
		TypeID     mrtype.KeyInt32        `json:"typeId" upd:"type_id"`                       // material_types::type_id
		Length     measure.Micrometer     `json:"length" sort:"length" upd:"laminate_length"` // mkm (mm * 1000)
		Width      measure.Micrometer     `json:"width" sort:"width" upd:"laminate_width"`    // mkm (mm * 1000)
		Thickness  measure.Micrometer     `json:"thickness" upd:"laminate_thickness"`         // mkm (mm * 1000)
		Weight     measure.GramsPerMeter2 `json:"weight" sort:"weight" upd:"laminate_weight"` // g/m2
		Status     mrenum.ItemStatus      `json:"status"`
		CreatedAt  time.Time              `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time             `json:"updatedAt,omitempty" sort:"updatedAt"`
	}

	// LaminateParams - comment struct.
	LaminateParams struct {
		Filter LaminateListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// LaminateListFilter - comment struct.
	LaminateListFilter struct {
		SearchText string
		TypeIDs    []mrtype.KeyInt32
		Length     mrtype.RangeInt64
		Width      mrtype.RangeInt64
		Weight     mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
