package view

import (
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateCatalogLaminate struct {
        Article   string `json:"article" validate:"required,min=3,max=32,article"`
        Caption   string `json:"caption" validate:"required,max=64"`
        TypeId    mrentity.KeyInt32 `json:"typeId" validate:"required,gte=1"`
        Length    entity.Micrometer `json:"length" validate:"required,gte=1,lte=1000000000"`
        Weight    entity.GramsPerMeter2 `json:"weight" validate:"required,gte=1,lte=10000"`
        Thickness entity.Micrometer `json:"thickness" validate:"required,gte=1,lte=1000000"`
    }

    StoreCatalogLaminate struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Article   string `json:"article" validate:"omitempty,min=3,max=32,article"`
        Caption   string `json:"caption" validate:"omitempty,max=64"`
        TypeId    mrentity.KeyInt32 `json:"typeId" validate:"omitempty,gte=1"`
        Length    entity.Micrometer `json:"length" validate:"omitempty,gte=1,lte=1000000000"`
        Weight    entity.GramsPerMeter2 `json:"weight" validate:"omitempty,gte=1,lte=10000"`
        Thickness entity.Micrometer `json:"thickness" validate:"omitempty,gte=1,lte=1000000"`
    }
)
