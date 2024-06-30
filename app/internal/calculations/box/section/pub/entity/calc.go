package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCalcResult = "public-api.Catalog.CalcResult" // ModelNameCalcResult - название сущности
)

type (
	// CalcResult - comment struct.
	CalcResult struct { // DB: printshop_catalog.calc_results
		ID         mrtype.KeyInt32 `json:"id"` // result_id
		TagVersion int32           `json:"tagVersion"`
		Caption    string          `json:"caption" sort:"caption,default" upd:"result_caption"`
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"`
	}
)
