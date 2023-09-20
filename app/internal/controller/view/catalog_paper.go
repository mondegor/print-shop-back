package view

import (
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateCatalogPaperRequest struct {
        Article   string `json:"article" validate:"required,min=3,max=32,article"`
        Caption   string `json:"caption" validate:"required,max=64"`
        Length    entity.Micrometer `json:"length" validate:"required,gte=1,lte=10000000"`
        Width     entity.Micrometer `json:"width" validate:"required,gte=1,lte=10000000"`
        Density   entity.GramsPerMeter2 `json:"density" validate:"required,gte=1,lte=10000"`
        ColorId   mrentity.KeyInt32 `json:"colorId" validate:"required,gte=1"`
        FactureId mrentity.KeyInt32 `json:"factureId" validate:"required,gte=1"`
        Thickness entity.Micrometer `json:"thickness" validate:"required,gte=1,lte=1000000"`
        Sides     entity.CatalogPaperSide `json:"sides" validate:"required,max=32"`
    }

    StoreCatalogPaperRequest struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Article   string `json:"article" validate:"omitempty,min=3,max=32,article"`
        Caption   string `json:"caption" validate:"omitempty,max=64"`
        Length    entity.Micrometer `json:"length" validate:"omitempty,gte=1,lte=10000000"`
        Width     entity.Micrometer `json:"width" validate:"omitempty,gte=1,lte=10000000"`
        Density   entity.GramsPerMeter2 `json:"density" validate:"omitempty,gte=1,lte=10000"`
        ColorId   mrentity.KeyInt32 `json:"colorId" validate:"omitempty,gte=1"`
        FactureId mrentity.KeyInt32 `json:"factureId" validate:"omitempty,gte=1"`
        Thickness entity.Micrometer `json:"thickness" validate:"omitempty,gte=1,lte=1000000"`
        Sides     entity.CatalogPaperSide `json:"sides" validate:"required,max=32"`
    }
)
