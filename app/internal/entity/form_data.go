package entity

import (
    "print-shop-back/pkg/mrentity"
    "time"
)

type (
    FormData struct { // DB: form_data
        Id        mrentity.KeyInt32 `json:"id"` // form_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Caption   string `json:"caption"` // form_caption
        Detailing ItemDetailing `json:"formDetailing"` // form_detailing
        Body      string `json:"formBody"` // form_body_compiled
        Status    ItemStatus `json:"status"` // form_status
    }

    FormDataListFilter struct {
        Detailing []ItemDetailing
        Statuses  []ItemStatus
    }
)
