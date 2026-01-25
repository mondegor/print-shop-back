package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      availability.PaperFactureStorage
		errorWrapper errors.Wrapper
		trace        mrtrace.Tracer
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(
	storage availability.PaperFactureStorage,
	trace mrtrace.Tracer,
) *PaperFacture {
	return &PaperFacture{
		storage:      storage,
		errorWrapper: errors.NewUseCaseWrapper(),
		trace:        trace,
	}
}

// CheckAvailability - comment method.
func (uc *PaperFacture) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", conv.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPaperFactureRequired
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return api.ErrPaperFactureNotFound.Wrap(err, itemID)
		}

		return uc.errorWrapper.Wrap(err)
	} else if status != itemstatus.Enabled {
		return api.ErrPaperFactureNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperFacture) traceCmd(ctx context.Context, command string, data conv.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.PaperFactureAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
