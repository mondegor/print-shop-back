package dto

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrentity"
)

type (
    ChangeItemStatus struct {
        Version mrentity.Version `json:"version"`
        Status  entity.ItemStatus `json:"status" validate:"required,max=16"`
    }
)
