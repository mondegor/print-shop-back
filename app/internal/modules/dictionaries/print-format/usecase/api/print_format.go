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
	PrintFormat struct {
		storage       PrintFormatStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewPrintFormat(
	storage PrintFormatStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *PrintFormat {
	return &PrintFormat{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *PrintFormat) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return dictionaries.FactoryErrPrintFormatRequired.New()
	}

	if err := uc.storage.IsExists(ctx, itemID); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrPrintFormatNotFound.New(itemID)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, dictionaries.PrintFormatAPIName)
	}

	return nil
}

func (uc *PrintFormat) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", dictionaries.PrintFormatAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
