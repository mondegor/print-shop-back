package entity

import (
	"print-shop-back/pkg/libs/measure"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameLaminate = "admin-api.Catalog.Laminate"
)

type (
	Laminate struct { // DB: printshop_catalog.laminates
		ID         mrtype.KeyInt32        `json:"id"` // laminate_id
		TagVersion int32                  `json:"tagVersion"`
		Article    string                 `json:"article" sort:"article" upd:"laminate_article"`
		Caption    string                 `json:"caption" sort:"caption,default" upd:"laminate_caption"`
		TypeID     mrtype.KeyInt32        `json:"typeId" upd:"type_id"`                       // laminate_types::type_id
		Length     measure.Micrometer     `json:"length" sort:"length" upd:"laminate_length"` // (mm)
		Weight     measure.GramsPerMeter2 `json:"weight" sort:"weight" upd:"laminate_weight"` // (g/m2)
		Thickness  measure.Micrometer     `json:"thickness" upd:"laminate_thickness"`         // (mkm)
		Status     mrenum.ItemStatus      `json:"status"`
		CreatedAt  time.Time              `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time             `json:"updatedAt,omitempty" sort:"updatedAt"`
	}

	LaminateParams struct {
		Filter LaminateListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	LaminateListFilter struct {
		SearchText string
		TypeIDs    []mrtype.KeyInt32
		Length     mrtype.RangeInt64
		Weight     mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
