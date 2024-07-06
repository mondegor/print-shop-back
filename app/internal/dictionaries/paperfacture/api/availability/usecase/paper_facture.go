package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      availability.PaperFactureStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(storage availability.PaperFactureStorage, errorWrapper mrcore.UsecaseErrorWrapper) *PaperFacture {
	return &PaperFacture{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// CheckingAvailability - comment method.
func (uc *PaperFacture) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return api.ErrPaperFactureRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrPaperFactureNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, api.PaperFactureAvailabilityName)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrPaperFactureNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperFacture) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", api.PaperFactureAvailabilityName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
