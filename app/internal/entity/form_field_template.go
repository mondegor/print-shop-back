package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

const ModelNameFormFieldTemplate = "FormFieldTemplate"

type (
    FormFieldTemplate struct { // DB: form_field_templates
        Id        mrentity.KeyInt32 `json:"id"` // template_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created

        ParamName string `json:"paramName" db:"param_name"`
        Caption   string `json:"caption" db:"template_caption"`
        Type      FormFieldTemplateType `json:"fieldType" db:"field_type"`
        Detailing ItemDetailing `json:"fieldDetailing" db:"field_detailing"`
        Body      string `json:"fieldBody" db:"field_body"`

        Status    ItemStatus `json:"status"` // template_status
    }

    FormFieldTemplateListFilter struct {
        Detailing []ItemDetailing
        Statuses  []ItemStatus
    }
)
