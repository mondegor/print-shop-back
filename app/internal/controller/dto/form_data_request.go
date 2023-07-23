package dto

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
)

type (
    CreateFormData struct {
        ParamName string `json:"paramName" validate:"required,min=4,max=32,variable"`
        Caption   string `json:"caption" validate:"required,max=128"`
        Detailing entity.ItemDetailing `json:"formDetailing" validate:"required,max=16"`
    }

    StoreFormData struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        ParamName string `json:"paramName" validate:"omitempty,min=4,max=32,variable"`
        Caption   string `json:"caption" validate:"omitempty,max=128"`
        Detailing entity.ItemDetailing `json:"formDetailing" validate:"omitempty,max=16"`
    }
)
