package dto

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrentity"
)

type (
    CreateFormData struct {
        Caption   string `json:"caption" validate:"required,max=128"`
        Detailing entity.ItemDetailing `json:"formDetailing" validate:"required,max=16"`
    }

    StoreFormData struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"required,max=128"`
        Detailing entity.ItemDetailing `json:"formDetailing" validate:"required,max=16"`
    }
)
