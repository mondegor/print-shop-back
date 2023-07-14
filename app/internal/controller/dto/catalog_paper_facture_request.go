package dto

import (
    "calc-user-data-back-adm/pkg/mrentity"
)

type (
    CreateCatalogPaperFacture struct {
        Caption   string `json:"caption" validate:"required,max=128"`
    }

    StoreCatalogPaperFacture struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"required,max=128"`
    }
)
