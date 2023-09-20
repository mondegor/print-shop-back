package view

import (
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    ChangeResourceStatusRequest struct {
        Version mrentity.Version `json:"version"`
        Status  entity.ResourceStatus `json:"status" validate:"required,max=16"`
    }
)
