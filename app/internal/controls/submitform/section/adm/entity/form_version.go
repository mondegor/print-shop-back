package entity

import (
	"time"

	"github.com/mondegor/print-shop-back/pkg/controls/enum"

	"github.com/google/uuid"
)

const (
	ModelNameFormVersion = "admin-api.Controls.FormVersion" // ModelNameFormVersion - название сущности
)

type (
	// FormVersion - comment struct.
	// FormVersion - comment struct.
	FormVersion struct { // DB: printshop_controls.submit_form_versions
		ID             uuid.UUID             `json:"-"` // form_id
		Version        int32                 `json:"version"`
		RewriteName    string                `json:"rewriteName"`
		Caption        string                `json:"caption"`
		Detailing      enum.ElementDetailing `json:"-"`
		Body           []byte                `json:"-"`
		ActivityStatus enum.ActivityStatus   `json:"activityStatus"`
		CreatedAt      time.Time             `json:"createdAt"`
		UpdatedAt      *time.Time            `json:"updatedAt,omitempty"`
	}

	// PrimaryKey - comment struct.
	// PrimaryKey - comment struct.
	PrimaryKey struct {
		FormID  uuid.UUID
		Version int32
	}

	// FormVersionStatus - comment struct.
	// FormVersionStatus - comment struct.
	FormVersionStatus struct {
		FormID         uuid.UUID
		Version        int32
		ActivityStatus enum.ActivityStatus
	}
)
