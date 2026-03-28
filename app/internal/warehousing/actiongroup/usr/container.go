package usr

import (
	"context"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
)

type (
	// ContainerService - comment interface.
	ContainerService interface {
		GetList(ctx context.Context, params dto.ContainerParams) (items []entity.Container, hasNext bool, err error)
		SaveTags(ctx context.Context, item entity.UpdateContainerTags) (tagVersion uint32, err error)
	}

	// ContainerStorage - comment interface.
	ContainerStorage interface {
		FetchByCondition(ctx context.Context, params dto.ContainerParams) (rows []entity.Container, hasNext bool, err error)
		FetchOne(ctx context.Context, accountID uuid.UUID, rowID uint64) (entity.Container, error)
		IsExist(ctx context.Context, accountID uuid.UUID, rowID uint64) error
		FetchMaxMarker(ctx context.Context, accountID uuid.UUID, code string) (marker uint16, err error)
		Insert(ctx context.Context, row entity.Container) (rowID uint64, err error)
		UpdateTags(ctx context.Context, row entity.UpdateContainerTags) (tagVersion uint32, err error)
	}
)
