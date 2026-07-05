package usecase

import (
	"context"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrworkflow/itemstatus"
	"github.com/mondegor/go-core/util/conv"

	"print-shop-back/internal/adapter/trace"
	"print-shop-back/internal/dictionaries/printformat/api/availability"
	"print-shop-back/pkg/dictionaries/api"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage      availability.PrintFormatStorage
		errorWrapper errors.Wrapper
		tracer       trace.Tracer
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(
	storage availability.PrintFormatStorage,
	tracer trace.Tracer,
) *PrintFormat {
	return &PrintFormat{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
		tracer:       tracer,
	}
}

// CheckAvailability - comment method.
func (uc *PrintFormat) CheckAvailability(ctx context.Context, itemID uint64) error {
	uc.traceCmd(ctx, "CheckAvailability", conv.Group{"id": itemID})

	if itemID == 0 {
		return api.ErrPrintFormatRequired
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return api.ErrPrintFormatNotFound.Wrap(err, itemID)
		}

		return uc.errorWrapper.Wrap(err)
	} else if status != itemstatus.Enabled {
		return api.ErrPrintFormatNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PrintFormat) traceCmd(ctx context.Context, command string, data conv.Group) {
	uc.tracer.Trace(
		ctx,
		"storage", api.PrintFormatAvailabilityName,
		"cmd", command,
		"data", data,
	)
}
