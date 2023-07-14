package dto

import (
    "calc-user-data-back-adm/pkg/mrentity"
)

type (
    CreateCatalogPaper struct {
        Article   string `json:"article" validate:"required,min=4,max=32,article"`
        Caption   string `json:"caption" validate:"required,max=128"`
        Length    mrentity.Micrometer `json:"length" validate:"required,gte=1,lte=10000"`
        Width     mrentity.Micrometer `json:"width" validate:"required,gte=1,lte=10000"`
        Density   mrentity.GramsPerMeter2 `json:"density" validate:"required,gte=1,lte=10000"`
        ColorId   mrentity.KeyInt32 `json:"colorId" validate:"required,gte=1"`
        FactureId mrentity.KeyInt32 `json:"factureId" validate:"required,gte=1"`
        Thickness mrentity.Micrometer `json:"thickness" validate:"required,gte=1,lte=10000"`
    }

    StoreCatalogPaper struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Article   string `json:"article" validate:"required,min=4,max=32,article"`
        Caption   string `json:"caption" validate:"required,max=128"`
        Length    mrentity.Micrometer `json:"length" validate:"required,gte=1,lte=10000"`
        Width     mrentity.Micrometer `json:"width" validate:"required,gte=1,lte=10000"`
        Density   mrentity.GramsPerMeter2 `json:"density" validate:"required,gte=1,lte=10000"`
        ColorId   mrentity.KeyInt32 `json:"colorId" validate:"required,gte=1"`
        FactureId mrentity.KeyInt32 `json:"factureId" validate:"required,gte=1"`
        Thickness mrentity.Micrometer `json:"thickness" validate:"required,gte=1,lte=10000"`
    }
)
