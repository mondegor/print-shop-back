package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/pkg/catalog/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

const (
	ModelNamePaper = "admin-api.Catalog.Paper" // ModelNamePaper - название сущности
)

type (
	// Paper - comment struct.
	Paper struct { // DB: printshop_catalog.papers
		ID         uint64                    `json:"id"` // paper_id
		TagVersion uint32                    `json:"tagVersion"`
		Article    string                    `json:"article" sort:"article" upd:"paper_article"`
		Caption    string                    `json:"caption" sort:"caption,default" upd:"paper_caption"`
		TypeID     uint64                    `json:"typeId" upd:"type_id"`       // material_types::type_id
		ColorID    uint64                    `json:"colorId" upd:"color_id"`     // paper_colors::color_id
		FactureID  uint64                    `json:"factureId" upd:"facture_id"` // paper_factures::facture_id
		Width      measure.Meter             `json:"width" sort:"width" upd:"paper_width"`
		Height     measure.Meter             `json:"height" sort:"height" upd:"paper_height"`
		Thickness  measure.Meter             `json:"thickness" upd:"paper_thickness"`
		Density    measure.KilogramPerMeter2 `json:"density" sort:"density" upd:"paper_density"`
		Sides      enum.PaperSide            `json:"sides" upd:"paper_sides"`
		Status     mrenum.ItemStatus         `json:"status"`
		CreatedAt  time.Time                 `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time                 `json:"updatedAt" sort:"updatedAt"`
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
		TypeIDs    []uint64
		ColorIDs   []uint64
		FactureIDs []uint64
		Width      measure.RangeMeter
		Height     measure.RangeMeter
		Density    measure.RangeKilogramPerMeter2
		Statuses   []mrenum.ItemStatus
	}
)
