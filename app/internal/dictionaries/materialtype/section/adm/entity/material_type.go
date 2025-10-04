package entity

import (
	"time"

	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-webcore/mrenum"
)

const (
	// ModelNameMaterialType - название сущности.
	ModelNameMaterialType = "admin-api.Dictionaries.MaterialType"
)

type (
	// MaterialType - comment struct.
	MaterialType struct { // DB: printshop_dictionaries.material_types
		ID         uint64            `json:"id"` // type_id
		TagVersion uint32            `json:"tagVersion"`
		Caption    string            `json:"caption" sort:"caption,default"`
		Status     mrenum.ItemStatus `json:"status"`
		CreatedAt  time.Time         `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time         `json:"updatedAt" sort:"updatedAt"`
	}

	// MaterialTypeParams - comment struct.
	MaterialTypeParams struct {
		Filter MaterialTypeListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// MaterialTypeListFilter - comment struct.
	MaterialTypeListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
