package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage      availability.PrintFormatStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(storage availability.PrintFormatStorage, errorWrapper mrcore.UsecaseErrorWrapper) *PrintFormat {
	return &PrintFormat{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// CheckingAvailability - comment method.
func (uc *PrintFormat) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return api.ErrPrintFormatRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrPrintFormatNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, api.PrintFormatAvailabilityName)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrPrintFormatNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PrintFormat) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", api.PrintFormatAvailabilityName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
