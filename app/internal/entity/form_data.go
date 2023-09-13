package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameFormData = "FormData"

type (
    FormData struct { // DB: form_data
        Id        mrentity.KeyInt32 `json:"id"` // form_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created

        ParamName  string `json:"paramName" db:"param_name"`
        Caption   string `json:"caption" db:"form_caption"`
        Detailing ItemDetailing `json:"formDetailing" db:"form_detailing"`

        Body      string `json:"formBody"` // form_body_compiled
        Status    mrcom.ItemStatus `json:"status"` // form_status
    }

    FormDataListFilter struct {
        Detailing []ItemDetailing
        Statuses  []mrcom.ItemStatus
    }
)
