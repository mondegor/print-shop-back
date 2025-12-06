package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage      availability.PrintFormatStorage
		errorWrapper mrerr.UseCaseErrorWrapper
		trace        mrtrace.Tracer
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(
	storage availability.PrintFormatStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
	trace mrtrace.Tracer,
) *PrintFormat {
	return &PrintFormat{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, api.PrintFormatAvailabilityName),
		trace:        trace,
	}
}

// CheckAvailability - comment method.
func (uc *PrintFormat) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", mrargs.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPrintFormatRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrPrintFormatNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	} else if status != itemstatus.Enabled {
		return api.ErrPrintFormatNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PrintFormat) traceCmd(ctx context.Context, command string, data mrargs.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.PrintFormatAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
