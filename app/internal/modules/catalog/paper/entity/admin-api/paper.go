package entity

import (
	entity_shared "print-shop-back/internal/modules/catalog/paper/entity/shared"
	"print-shop-back/pkg/libs/measure"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaper = "admin-api.Catalog.Paper"
)

type (
	Paper struct { // DB: printshop_catalog.papers
		ID         mrtype.KeyInt32 `json:"id"`                                   // paper_id
		TagVersion int32           `json:"tagVersion"`                           // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // created_at
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // updated_at

		Article   string                  `json:"article" sort:"article" upd:"paper_article"`
		Caption   string                  `json:"caption" sort:"caption,default" upd:"paper_caption"`
		ColorID   mrtype.KeyInt32         `json:"colorId" upd:"color_id"`                     // printshop_dictionaries.paper_colors::color_id
		FactureID mrtype.KeyInt32         `json:"factureId" upd:"facture_id"`                 // printshop_dictionaries.paper_factures::facture_id
		Length    measure.Micrometer      `json:"length" sort:"length" upd:"paper_length"`    // (mm)
		Width     measure.Micrometer      `json:"width" sort:"width" upd:"paper_width"`       // (mm)
		Density   measure.GramsPerMeter2  `json:"density" sort:"density" upd:"paper_density"` // (g/m2)
		Thickness measure.Micrometer      `json:"thickness" upd:"paper_thickness"`            // (mkm)
		Sides     entity_shared.PaperSide `json:"sides" upd:"paper_sides"`

		Status mrenum.ItemStatus `json:"status"` // paper_status
	}

	PaperParams struct {
		Filter PaperListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	PaperListFilter struct {
		SearchText string
		ColorIDs   []mrtype.KeyInt32
		FactureIDs []mrtype.KeyInt32
		Length     mrtype.RangeInt64
		Width      mrtype.RangeInt64
		Density    mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
