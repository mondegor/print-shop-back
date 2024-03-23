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
	LaminateType struct {
		storage       LaminateTypeStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewLaminateType(
	storage LaminateTypeStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *LaminateType {
	return &LaminateType{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *LaminateType) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return dictionaries.FactoryErrLaminateTypeRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrLaminateTypeNotFound.New(itemID)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, dictionaries.LaminateTypeAPIName)
	} else if status != mrenum.ItemStatusEnabled {
		return dictionaries.FactoryErrLaminateTypeNotAvailable.New(itemID)
	}

	return nil
}

func (uc *LaminateType) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", dictionaries.LaminateTypeAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
