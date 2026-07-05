package entity

import (
	"time"

	"github.com/mondegor/go-core/mrtype"

	"print-shop-back/internal/adapter/workflow"
)

const (
	// ModelNamePaperFacture - название сущности.
	ModelNamePaperFacture = "admin-api.Dictionaries.PaperFacture"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct { // DB: printshop_dictionaries.paper_factures
		ID         uint64              `json:"id"` // facture_id
		TagVersion uint32              `json:"tagVersion"`
		Caption    string              `json:"caption" sort:"caption,default"`
		Status     workflow.ItemStatus `json:"status"`
		CreatedAt  time.Time           `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time           `json:"updatedAt" sort:"updatedAt"`
	}

	// PaperFactureParams - comment struct.
	PaperFactureParams struct {
		Filter PaperFactureListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// PaperFactureListFilter - comment struct.
	PaperFactureListFilter struct {
		SearchText string
		Statuses   []workflow.ItemStatus
	}
)
