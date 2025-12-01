package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		storage      availability.MaterialTypeStorage
		errorWrapper mrerr.UseCaseErrorWrapper
		trace        mrtrace.Tracer
	}
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(
	storage availability.MaterialTypeStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
	trace mrtrace.Tracer,
) *MaterialType {
	return &MaterialType{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, api.MaterialTypeAvailabilityName),
		trace:        trace,
	}
}

// CheckAvailability - comment method.
func (uc *MaterialType) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", mrargs.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrMaterialTypeRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return api.ErrMaterialTypeNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	} else if status != itemstatus.Enabled {
		return api.ErrMaterialTypeNotAvailable.New(itemID)
	}

	return nil
}

func (uc *MaterialType) traceCmd(ctx context.Context, command string, data mrargs.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.MaterialTypeAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
