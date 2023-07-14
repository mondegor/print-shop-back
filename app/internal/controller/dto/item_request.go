package dto

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
)

type (
    ChangeItemStatus struct {
        Version mrentity.Version `json:"version"`
        Status  entity.ItemStatus `json:"status" validate:"required,max=16"`
    }
)
