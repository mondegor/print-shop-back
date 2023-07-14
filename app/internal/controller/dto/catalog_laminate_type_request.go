package dto

import (
    "calc-user-data-back-adm/pkg/mrentity"
)

type (
    CreateCatalogLaminateType struct {
        Caption   string `json:"caption" validate:"required,max=128"`
    }

    StoreCatalogLaminateType struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"required,max=128"`
    }
)
