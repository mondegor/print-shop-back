package view

import (
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    ChangeItemStatus struct {
        Version mrentity.Version `json:"version"`
        Status  mrcom.ItemStatus `json:"status" validate:"required,max=16"`
    }
)
