package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/pkg/controls/type/activitystatus"
	"github.com/mondegor/print-shop-back/pkg/controls/type/elementdetailing"
)

const (
	// ModelNameFormVersion - название сущности.
	ModelNameFormVersion = "admin-api.Controls.FormVersion"
)

type (
	// FormVersion - comment struct.
	// FormVersion - comment struct.
	FormVersion struct { // DB: printshop_controls.submit_form_versions
		ID             uuid.UUID             `json:"-"` // form_id
		Version        int32                 `json:"version"`
		RewriteName    string                `json:"rewriteName"`
		Caption        string                `json:"caption"`
		Detailing      elementdetailing.Enum `json:"-"`
		Body           []byte                `json:"-"`
		ActivityStatus activitystatus.Enum   `json:"activityStatus"`
		CreatedAt      time.Time             `json:"createdAt"`
		UpdatedAt      time.Time             `json:"updatedAt"`
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
		ActivityStatus activitystatus.Enum
	}
)
