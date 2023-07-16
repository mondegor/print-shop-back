package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

type (
    FormFieldItem struct { // DB: form_fields
        Id         mrentity.KeyInt32 `json:"id"` // field_id
        FormId     mrentity.KeyInt32 `json:"formId"` // form_id
        TemplateId mrentity.KeyInt32 `json:"templateId"` // template_id
        Version    mrentity.Version `json:"version"` // tag_version
        CreatedAt  time.Time `json:"createdAt"` // datetime_created
        ParamName  string `json:"paramName"` // param_name
        Caption    string `json:"caption"` // field_caption
        Type      FormFieldTemplateType `json:"fieldType"` // form_field_templates::field_type
        Detailing ItemDetailing `json:"fieldDetailing"` // form_field_templates::field_detailing
        Required   bool   `json:"fieldRequired"` // field_required
        Body       string `json:"fieldBody"` // form_field_templates::field_body
    }

    FormFieldItemListFilter struct {
        FormId mrentity.KeyInt32
        Detailing []ItemDetailing
    }
)
