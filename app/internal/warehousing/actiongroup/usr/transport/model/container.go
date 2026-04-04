package model

import (
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"github.com/mondegor/print-shop-back/internal/warehousing/xtype"
)

type (
	// CreateContainerRequest - comment struct.
	CreateContainerRequest struct {
		Code             string       `json:"code,omitempty"`
		Volume           xtype.Volume `json:"volume,omitempty"`
		Tags             []string     `json:"tags,omitempty"`
		Images           []string     `json:"images,omitempty"`
		LocationID       uint64       `json:"location_id"`
		ExemplarQuantity int          `json:"exemplar_quantity,omitempty"`
	}

	// CreateContainerGroupRequest - comment struct.
	CreateContainerGroupRequest struct {
		Code       string       `json:"code,omitempty"`
		Volume     xtype.Volume `json:"volume"`
		LocationID uint64       `json:"location_id"`
	}

	// SaveContainerTagsRequest - comment struct.
	SaveContainerTagsRequest struct {
		ContainerID uint64   `json:"container_id"`
		TagVersion  uint32   `json:"tag_version" validate:"required,gte=1"`
		Tags        []string `json:"tags"`
	}

	// ContainerListResponse - comment struct.
	ContainerListResponse struct {
		Items   []entity.Container `json:"items"`
		Cursor  string             `json:"cursor"`
		HasNext bool               `json:"has_next"`
	}

	// SuccessCreatedContainerResponse - comment struct.
	SuccessCreatedContainerResponse struct {
		ContainerID uint64 `json:"container_id"`
		Code        string `json:"code"`
		Marker      uint16 `json:"marker"`
		StockID     uint64 `json:"stock_id"`
	}
)
