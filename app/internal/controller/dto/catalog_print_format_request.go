package dto

import (
    "print-shop-back/pkg/mrentity"
)

type (
    CreateCatalogPrintFormat struct {
        Caption   string `json:"caption" validate:"required,max=128"`
        Length  mrentity.Micrometer `json:"length" validate:"required,gte=1,lte=10000"`
        Width   mrentity.Micrometer `json:"width" validate:"required,gte=1,lte=10000"`
    }

    StoreCatalogPrintFormat struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"required,max=128"`
        Length  mrentity.Micrometer `json:"length" validate:"required,gte=1,lte=10000"`
        Width   mrentity.Micrometer `json:"width" validate:"required,gte=1,lte=10000"`
    }
)
