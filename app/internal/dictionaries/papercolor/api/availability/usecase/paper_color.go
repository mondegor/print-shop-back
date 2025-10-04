package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      availability.PaperColorStorage
		errorWrapper mrerr.UseCaseErrorWrapper
		trace        mrtrace.Tracer
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(
	storage availability.PaperColorStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
	trace mrtrace.Tracer,
) *PaperColor {
	return &PaperColor{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, api.PaperColorAvailabilityName),
		trace:        trace,
	}
}

// CheckingAvailability - comment method.
func (uc *PaperColor) CheckingAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckingAvailability", mrargs.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPaperColorRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return api.ErrPaperColorNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrPaperColorNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperColor) traceCmd(ctx context.Context, command string, data mrargs.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.PaperColorAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
