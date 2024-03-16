package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNamePaperFacture = "admin-api.Dictionaries.PaperFacture"
)

type (
	PaperFacture struct { // DB: printshop_dictionaries.paper_factures
		ID         mrtype.KeyInt32 `json:"id"`                                   // facture_id
		TagVersion int32           `json:"tagVersion"`                           // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // created_at
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // updated_at

		Caption string `json:"caption" sort:"caption,default"` // facture_caption

		Status mrenum.ItemStatus `json:"status"` // facture_status
	}

	PaperFactureParams struct {
		Filter PaperFactureListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	PaperFactureListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
