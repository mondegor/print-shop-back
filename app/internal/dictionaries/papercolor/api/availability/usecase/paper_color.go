package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      availability.PaperColorStorage
		errorWrapper errors.Wrapper
		trace        mrtrace.Tracer
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(
	storage availability.PaperColorStorage,
	trace mrtrace.Tracer,
) *PaperColor {
	return &PaperColor{
		storage:      storage,
		errorWrapper: errors.NewUseCaseWrapper(),
		trace:        trace,
	}
}

// CheckAvailability - comment method.
func (uc *PaperColor) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", conv.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPaperColorRequired
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return api.ErrPaperColorNotFound.Wrap(err, itemID)
		}

		return uc.errorWrapper.Wrap(err)
	} else if status != itemstatus.Enabled {
		return api.ErrPaperColorNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperColor) traceCmd(ctx context.Context, command string, data conv.Group) {
	uc.trace.Trace(
		ctx,
		"storage", api.PaperColorAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
