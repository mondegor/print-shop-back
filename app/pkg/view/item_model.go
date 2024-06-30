package view

import (
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ChangeItemStatusRequest - comment struct.
	ChangeItemStatusRequest struct {
		TagVersion int32             `json:"tagVersion" validate:"required,gte=1"`
		Status     mrenum.ItemStatus `json:"status" validate:"required"`
	}

	// MoveItemRequest - comment struct.
	MoveItemRequest struct {
		AfterNodeID mrtype.KeyInt32 `json:"afterId"`
	}

	// SuccessCreatedItemResponse - comment struct.
	SuccessCreatedItemResponse struct {
		ItemID string `json:"id"`
	}

	// SuccessCreatedItemInt32Response - comment struct.
	SuccessCreatedItemInt32Response struct {
		ItemID mrtype.KeyInt32 `json:"id"`
	}
)
