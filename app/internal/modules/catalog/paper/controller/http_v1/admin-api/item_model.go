package http_v1

import (
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ChangeItemStatusRequest struct {
		TagVersion int32             `json:"tagVersion" validate:"required,gte=1"`
		Status     mrenum.ItemStatus `json:"status" validate:"required"`
	}

	MoveItemRequest struct {
		AfterNodeID mrtype.KeyInt32 `json:"afterId"`
	}

	SuccessCreatedItemResponse struct {
		ItemID  string `json:"id"`
		Message string `json:"message,omitempty"`
	}
)
