package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		storage      availability.MaterialTypeStorage
		errorWrapper errors.Wrapper
		trace        mrtrace.Tracer
	}
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(
	storage availability.MaterialTypeStorage,
	trace mrtrace.Tracer,
) *MaterialType {
	return &MaterialType{
		storage:      storage,
		errorWrapper: errors.NewUseCaseWrapper(),
		trace:        trace,
	}
}

// CheckAvailability - comment method.
func (uc *MaterialType) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", conv.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrMaterialTypeRequired
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return api.ErrMaterialTypeNotFound.Wrap(err, itemID)
		}

		return uc.errorWrapper.Wrap(err)
	} else if status != itemstatus.Enabled {
		return api.ErrMaterialTypeNotAvailable.New(itemID)
	}

	return nil
}

func (uc *MaterialType) traceCmd(ctx context.Context, command string, data conv.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.MaterialTypeAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
