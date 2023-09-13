package view

import (
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateCatalogPrintFormat struct {
        Caption   string `json:"caption" validate:"required,max=64"`
        Length  entity.Micrometer `json:"length" validate:"required,gte=1,lte=10000000"`
        Width   entity.Micrometer `json:"width" validate:"required,gte=1,lte=10000000"`
    }

    StoreCatalogPrintFormat struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"omitempty,max=64"`
        Length  entity.Micrometer `json:"length" validate:"omitempty,gte=1,lte=10000000"`
        Width   entity.Micrometer `json:"width" validate:"omitempty,gte=1,lte=10000000"`
    }
)
