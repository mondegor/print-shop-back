package dto

import (
	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/warehousing/enum/locationkind"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
)

type (
	// CreateContainer - comment struct.
	CreateContainer struct {
		AccountID  uuid.UUID
		Kind       locationkind.Enum
		Code       string
		Volume     xtype.Volume
		Tags       []string
		Images     []string
		LocationID uint64
		Quantity   int
	}

	// CreateContainerResult - comment struct.
	CreateContainerResult struct {
		ID      uint64
		Code    string
		Marker  uint16
		StockID uint64
	}

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
