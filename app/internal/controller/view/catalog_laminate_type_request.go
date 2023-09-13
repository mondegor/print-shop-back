package view

import "github.com/mondegor/go-storage/mrentity"

type (
    CreateCatalogLaminateType struct {
        Caption   string `json:"caption" validate:"required,max=64"`
    }

    StoreCatalogLaminateType struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"required,max=64"`
    }
)
