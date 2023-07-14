package entity

import (
    "calc-user-data-back-adm/pkg/mrentity"
    "time"
)

type (
    FormFieldTemplate struct { // DB: form_field_templates
        Id        mrentity.KeyInt32 `json:"id"` // template_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        ParamName string `json:"paramName"` // param_name
        Caption   string `json:"caption"` // template_caption
        Type      FormFieldTemplateType `json:"fieldType"` // field_type
        Detailing ItemDetailing `json:"fieldDetailing"` // field_detailing
        Body      string `json:"fieldBody"` // field_body
        Status    ItemStatus `json:"status"` // template_status
    }

    FormFieldTemplateListFilter struct {
        Detailing []ItemDetailing
        Statuses  []ItemStatus
    }
)
