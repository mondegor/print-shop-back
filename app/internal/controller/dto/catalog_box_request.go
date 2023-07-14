package dto

import "print-shop-back/pkg/mrentity"

type (
    CreateCatalogBox struct {
        Article string `json:"article" validate:"required,min=4,max=32,article"`
        Caption string `json:"caption" validate:"required,max=128"`
        Length  mrentity.Micrometer `json:"length" validate:"required,gte=1,lte=10000"`
        Width   mrentity.Micrometer `json:"width" validate:"required,gte=1,lte=10000"`
        Depth   mrentity.Micrometer `json:"depth" validate:"required,gte=1,lte=10000"`
    }

    StoreCatalogBox struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Article string `json:"article" validate:"required,min=4,max=32,article"`
        Caption string `json:"caption" validate:"required,max=128"`
        Length  mrentity.Micrometer `json:"length" validate:"required,gte=1,lte=10000"`
        Width   mrentity.Micrometer `json:"width" validate:"required,gte=1,lte=10000"`
        Depth   mrentity.Micrometer `json:"depth" validate:"required,gte=1,lte=10000"`
    }
)
