package view

import (
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateFormFieldTemplate struct {
        ParamName string `json:"paramName" validate:"required,min=4,max=32,variable"`
        Caption   string `json:"caption" validate:"required,max=128"`
        Type      entity.FormFieldTemplateType `json:"fieldType" validate:"required,max=16"`
        Detailing entity.ItemDetailing `json:"fieldDetailing" validate:"required,max=16"`
        Body      string `json:"fieldBody" validate:"required,max=65536"`
    }

    StoreFormFieldTemplate struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        ParamName string `json:"paramName" validate:"omitempty,min=4,max=32,variable"`
        Caption   string `json:"caption" validate:"omitempty,max=128"`
        Type      entity.FormFieldTemplateType `json:"fieldType" validate:"omitempty,max=16"`
        Detailing entity.ItemDetailing `json:"fieldDetailing" validate:"omitempty,max=16"`
        Body      string `json:"fieldBody" validate:"omitempty,max=65536"`
    }
)
