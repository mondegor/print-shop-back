package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameMaterialType = "admin-api.Dictionaries.MaterialType" // ModelNameMaterialType - название сущности
)

type (
	// MaterialType - comment struct.
	MaterialType struct { // DB: printshop_dictionaries.material_types
		ID         mrtype.KeyInt32   `json:"id"` // type_id
		TagVersion int32             `json:"tagVersion"`
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
