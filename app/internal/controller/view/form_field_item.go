package view

import "github.com/mondegor/go-storage/mrentity"

type (
    CreateFormFieldItemRequest struct {
        TemplateId mrentity.KeyInt32 `json:"templateId" validate:"required"`
        ParamName  string `json:"paramName" validate:"omitempty,min=4,max=32,variable"`
        Caption    string `json:"caption" validate:"omitempty,max=128"`
        Required   bool   `json:"fieldRequired"`
    }

    StoreFormFieldItemRequest struct {
        Version  mrentity.Version `json:"version" validate:"required,gte=1"`
        ParamName  string `json:"paramName" validate:"omitempty,min=4,max=32,variable"`
        Caption  string `json:"caption" validate:"omitempty,max=128"`
        Required bool   `json:"fieldRequired"`
    }

    MoveFormFieldItemRequest struct {
        AfterNodeId mrentity.KeyInt32 `json:"afterId"`
    }
)
