package usecase_api

import (
	"context"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperColor struct {
		storage       PaperColorStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewPaperColor(
	storage PaperColorStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *PaperColor {
	return &PaperColor{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *PaperColor) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return dictionaries.FactoryErrPaperColorRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrPaperColorNotFound.New(itemID)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, dictionaries.PaperColorAPIName)
	} else if status != mrenum.ItemStatusEnabled {
		return dictionaries.FactoryErrPaperColorNotAvailable.New(itemID)
	}

	return nil
}

func (uc *PaperColor) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", dictionaries.PaperColorAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
