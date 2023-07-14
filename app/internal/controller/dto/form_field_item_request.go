package dto

import (
    "calc-user-data-back-adm/pkg/mrentity"
)

type (
    CreateFormFieldItem struct {
        TemplateId mrentity.KeyInt32 `json:"templateId" validate:"required"`
        ParamName  string `json:"paramName" validate:"min=4,max=32,variable"`
        Caption    string `json:"caption" validate:"max=128"`
        Required   bool   `json:"fieldRequired"`
    }

    StoreFormFieldItem struct {
        Version  mrentity.Version `json:"version" validate:"required,gte=1"`
        ParamName  string `json:"paramName" validate:"required,min=4,max=32,variable"`
        Caption  string `json:"caption" validate:"required,max=128"`
        Required bool   `json:"fieldRequired"`
    }

    MoveFormFieldItem struct {
        AfterNodeId mrentity.KeyInt32 `json:"afterId"`
    }
)
