package model

import (
	"print-shop-back/internal/warehousing/actiongroup/usr/entity"
)

type (
	// StoreListResponse - comment struct.
	StoreListResponse struct {
		Items   []entity.Store `json:"stores"`
		Cursor  string         `json:"cursor"`
		HasNext bool           `json:"has_next"`
	}
)
