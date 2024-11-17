package view

import (
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// ChangeItemStatusRequest - comment struct.
	ChangeItemStatusRequest struct {
		TagVersion uint32            `json:"tagVersion" validate:"required,gte=1"`
		Status     mrenum.ItemStatus `json:"status" validate:"required"`
	}

	// MoveItemRequest - comment struct.
	MoveItemRequest struct {
		AfterNodeID uint64 `json:"afterId"`
	}

	// SuccessCreatedItemResponse - comment struct.
	SuccessCreatedItemResponse struct {
		ItemID string `json:"id"`
	}

	// SuccessCreatedItemInt32Response - comment struct.
	SuccessCreatedItemInt32Response struct {
		ItemID uint64 `json:"id"`
	}
)
