package usecase_api

import (
	"context"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperFacture struct {
		storage       PaperFactureStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewPaperFacture(
	storage PaperFactureStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *PaperFacture {
	return &PaperFacture{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *PaperFacture) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return dictionaries.FactoryErrPaperFactureRequired.New()
	}

	if err := uc.storage.IsExists(ctx, itemID); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrPaperFactureNotFound.New(itemID)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, dictionaries.PaperFactureAPIName)
	}

	return nil
}

func (uc *PaperFacture) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", dictionaries.PaperFactureAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
