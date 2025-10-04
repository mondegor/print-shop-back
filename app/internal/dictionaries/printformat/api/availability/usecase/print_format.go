package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-webcore/mrenum"

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

// CheckingAvailability - comment method.
func (uc *PrintFormat) CheckingAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckingAvailability", mrargs.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPrintFormatRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return api.ErrPrintFormatNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	} else if status != mrenum.ItemStatusEnabled {
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
