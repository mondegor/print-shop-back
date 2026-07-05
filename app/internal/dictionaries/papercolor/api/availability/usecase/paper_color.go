package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrworkflow/itemstatus"
	"github.com/mondegor/go-sysmess/util/conv"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/dictionaries/papercolor/api/availability"
	"print-shop-back/pkg/dictionaries/api"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      availability.PaperColorStorage
		errorWrapper errors.Wrapper
		tracer       trace.Tracer
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(
	storage availability.PaperColorStorage,
	tracer trace.Tracer,
) *PaperColor {
	return &PaperColor{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
		tracer:       tracer,
	}
}

// CheckAvailability - comment method.
func (uc *PaperColor) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", conv.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPaperColorRequired
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return api.ErrPaperColorNotFound.Wrap(err, itemID)
		}

		return uc.errorWrapper.Wrap(err)
	} else if status != itemstatus.Enabled {
		return api.ErrPaperColorNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperColor) traceCmd(ctx context.Context, command string, data conv.Group) {
	uc.tracer.Trace(
		ctx,
		"storage", api.PaperColorAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
