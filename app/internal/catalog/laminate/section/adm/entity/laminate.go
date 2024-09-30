package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

const (
	ModelNameLaminate = "admin-api.Catalog.Laminate" // ModelNameLaminate - название сущности
)

type (
	// Laminate - comment struct.
	Laminate struct { // DB: printshop_catalog.laminates
		ID         mrtype.KeyInt32           `json:"id"` // laminate_id
		TagVersion int32                     `json:"tagVersion"`
		Article    string                    `json:"article" sort:"article" upd:"laminate_article"`
		Caption    string                    `json:"caption" sort:"caption,default" upd:"laminate_caption"`
		TypeID     mrtype.KeyInt32           `json:"typeId" upd:"type_id"` // material_types::type_id
		Length     measure.Meter             `json:"length" sort:"length" upd:"laminate_length"`
		Width      measure.Meter             `json:"width" sort:"width" upd:"laminate_width"`
		Thickness  measure.Meter             `json:"thickness" upd:"laminate_thickness"`
		WeightM2   measure.KilogramPerMeter2 `json:"weightM2" sort:"weightM2" upd:"laminate_weight_m2"`
		Status     mrenum.ItemStatus         `json:"status"`
		CreatedAt  time.Time                 `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time                 `json:"updatedAt" sort:"updatedAt"`
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
		Length     measure.RangeMeter
		Width      measure.RangeMeter
		Statuses   []mrenum.ItemStatus
	}
)
