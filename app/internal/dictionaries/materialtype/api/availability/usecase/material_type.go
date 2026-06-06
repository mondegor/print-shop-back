package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrworkflow/itemstatus"
	"github.com/mondegor/go-sysmess/util/conv"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/dictionaries/materialtype/api/availability"
	"print-shop-back/pkg/dictionaries/api"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		storage      availability.MaterialTypeStorage
		errorWrapper errors.Wrapper
		tracer       trace.Tracer
	}
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(
	storage availability.MaterialTypeStorage,
	tracer trace.Tracer,
) *MaterialType {
	return &MaterialType{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
		tracer:       tracer,
	}
}

// CheckAvailability - comment method.
func (uc *MaterialType) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", conv.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrMaterialTypeRequired
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return api.ErrMaterialTypeNotFound.Wrap(err, itemID)
		}

		return uc.errorWrapper.Wrap(err)
	} else if status != itemstatus.Enabled {
		return api.ErrMaterialTypeNotAvailable.New(itemID)
	}

	return nil
}

func (uc *MaterialType) traceCmd(ctx context.Context, command string, data conv.Group) {
	uc.tracer.Trace(
		ctx,
		"storage", api.MaterialTypeAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
