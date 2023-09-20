package view

import "github.com/mondegor/go-storage/mrentity"

type (
    CreateCatalogPaperColorRequest struct {
        Caption   string `json:"caption" validate:"required,max=64"`
    }

    StoreCatalogPaperColorRequest struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"required,max=64"`
    }
)
