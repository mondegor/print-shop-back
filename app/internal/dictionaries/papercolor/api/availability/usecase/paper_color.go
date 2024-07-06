package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      availability.PaperColorStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(storage availability.PaperColorStorage, errorWrapper mrcore.UsecaseErrorWrapper) *PaperColor {
	return &PaperColor{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// CheckingAvailability - comment method.
func (uc *PaperColor) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return api.ErrPaperColorRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrPaperColorNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, api.PaperColorAvailabilityName)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrPaperColorNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperColor) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", api.PaperColorAvailabilityName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
