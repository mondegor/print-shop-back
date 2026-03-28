package entity

import (
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/warehousing/enum/locationkind"
	"github.com/mondegor/print-shop-back/internal/warehousing/module"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
)

type (
	// Container - comment struct.
	Container struct { // DB: warehousing.containers
		ID         uint64            `json:"id"`
		TagVersion uint32            `json:"version"`
		Kind       locationkind.Enum `json:"kind"` // извлекается из ID
		AccountID  uuid.UUID         `json:"-"`
		Code       string            `json:"code"`
		Marker     uint16            `json:"marker"`
		Volume     xtype.Volume      `json:"volume"`
		Tags       []string          `json:"tags"`
		Images     []string          `json:"images"`
		CreatedAt  time.Time         `json:"created_at"`
		UpdatedAt  time.Time         `json:"updated_at"`
	}

	// UpdateContainerTags - comment struct.
	UpdateContainerTags struct { // DB: warehousing.containers
		ID         uint64
		AccountID  uuid.UUID
		TagVersion uint32
		Tags       []string
	}
)

// SequenceName - comment method.
func (s Container) SequenceName() string {
	if s.Kind == locationkind.Group {
		return module.DBTableNameContainers + "_group_id_seq"
	}

	return module.DBTableNameContainers + "_container_id_seq"
}

// CreateContainerCursorValue - comment func.
func CreateContainerCursorValue(items []Container) string {
	if len(items) == 0 {
		return ""
	}

	item := items[len(items)-1]

	return item.Code + "|" + strconv.Itoa(int(item.Marker))
}
