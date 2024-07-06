package entity

import (
	"time"

	"github.com/mondegor/print-shop-back/pkg/catalog/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaper = "admin-api.Catalog.Paper" // ModelNamePaper - название сущности
)

type (
	// Paper - comment struct.
	Paper struct { // DB: printshop_catalog.papers
		ID         mrtype.KeyInt32           `json:"id"` // paper_id
		TagVersion int32                     `json:"tagVersion"`
		Article    string                    `json:"article" sort:"article" upd:"paper_article"`
		Caption    string                    `json:"caption" sort:"caption,default" upd:"paper_caption"`
		TypeID     mrtype.KeyInt32           `json:"typeId" upd:"type_id"`       // material_types::type_id
		ColorID    mrtype.KeyInt32           `json:"colorId" upd:"color_id"`     // paper_colors::color_id
		FactureID  mrtype.KeyInt32           `json:"factureId" upd:"facture_id"` // paper_factures::facture_id
		Length     measure.Meter             `json:"length" sort:"length" upd:"paper_length"`
		Height     measure.Meter             `json:"height" sort:"height" upd:"paper_height"`
		Thickness  measure.Meter             `json:"thickness" upd:"paper_thickness"`
		Density    measure.KilogramPerMeter2 `json:"density" sort:"density" upd:"paper_density"`
		Sides      enum.PaperSide            `json:"sides" upd:"paper_sides"`
		Status     mrenum.ItemStatus         `json:"status"`
		CreatedAt  time.Time                 `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time                `json:"updatedAt,omitempty" sort:"updatedAt"`
	}

	// PaperParams - comment struct.
	PaperParams struct {
		Filter PaperListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// PaperListFilter - comment struct.
	PaperListFilter struct {
		SearchText string
		TypeIDs    []mrtype.KeyInt32
		ColorIDs   []mrtype.KeyInt32
		FactureIDs []mrtype.KeyInt32
		Length     mrtype.RangeInt64
		Height     mrtype.RangeInt64
		Density    mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
