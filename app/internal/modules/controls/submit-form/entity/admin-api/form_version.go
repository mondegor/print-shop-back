package entity

import (
	"print-shop-back/pkg/modules/controls/enums"
	"time"

	"github.com/google/uuid"
)

const (
	ModelNameFormVersion = "admin-api.Controls.FormVersion"
)

type (
	FormVersion struct { // DB: printshop_controls.submit_form_versions
		ID             uuid.UUID              `json:"-"` // form_id
		Version        int32                  `json:"version"`
		RewriteName    string                 `json:"rewriteName"`
		Caption        string                 `json:"caption"`
		Detailing      enums.ElementDetailing `json:"-"`
		Body           []byte                 `json:"-"`
		ActivityStatus enums.ActivityStatus   `json:"activityStatus"`
		CreatedAt      time.Time              `json:"createdAt"`
		UpdatedAt      *time.Time             `json:"updatedAt,omitempty"`
	}

	PrimaryKey struct {
		FormID  uuid.UUID
		Version int32
	}

	FormVersionStatus struct {
		FormID         uuid.UUID
		Version        int32
		ActivityStatus enums.ActivityStatus
	}
)
