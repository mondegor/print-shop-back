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
	LaminateType struct {
		storage       LaminateTypeStorage
		usecaseHelper *mrcore.UsecaseHelper
	}

	LaminateTypeStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
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

func (uc *LaminateType) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return dictionaries.FactoryErrLaminateTypeNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrLaminateTypeNotFound.New(id)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, dictionaries.LaminateTypeAPIName)
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
