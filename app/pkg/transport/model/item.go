package model

import (
	"print-shop-back/internal/adapter/workflow"
)

type (
	// ChangeItemStatusRequest - comment struct.
	ChangeItemStatusRequest struct {
		TagVersion uint32              `json:"tagVersion" validate:"required,gte=1"`
		Status     workflow.ItemStatus `json:"status" validate:"required"`
	}

	// MoveItemRequest - comment struct.
	MoveItemRequest struct {
		AfterNodeID uint64 `json:"afterId"`
	}

	// SuccessCreatedItemResponse - comment struct.
	SuccessCreatedItemResponse struct {
		ItemID     string `json:"id"`
		TagVersion uint32 `json:"tag_version,omitempty"`
	}

	// SuccessCreatedItemUintResponse - comment struct.
	SuccessCreatedItemUintResponse struct {
		ItemID     uint64 `json:"id"`
		TagVersion uint32 `json:"tag_version,omitempty"`
	}

	// SuccessSavedItemResponse - comment struct.
	SuccessSavedItemResponse struct {
		TagVersion uint32 `json:"tag_version"`
	}
)
