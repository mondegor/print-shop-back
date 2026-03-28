package dto

import (
	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
)

type (
	// ContainerParams - comment struct.
	ContainerParams struct {
		AccountID uuid.UUID
		Filter    ContainerFilter
		Cursor    xtype.ContainerCursor
	}

	// ContainerFilter - comment struct.
	ContainerFilter struct {
		SearchCode string
		SearchTags []string
	}
)
