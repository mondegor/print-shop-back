package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      availability.PaperFactureStorage
		errorWrapper mrerr.UseCaseErrorWrapper
		trace        mrtrace.Tracer
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(
	storage availability.PaperFactureStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
	trace mrtrace.Tracer,
) *PaperFacture {
	return &PaperFacture{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, api.PaperFactureAvailabilityName),
		trace:        trace,
	}
}

// CheckAvailability - comment method.
func (uc *PaperFacture) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", mrargs.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPaperFactureRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return api.ErrPaperFactureNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrPaperFactureNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperFacture) traceCmd(ctx context.Context, command string, data mrargs.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.PaperFactureAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
