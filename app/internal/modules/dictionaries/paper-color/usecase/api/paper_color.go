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
	PaperColor struct {
		storage       PaperColorStorage
		usecaseHelper *mrcore.UsecaseHelper
	}

	PaperColorStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
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

func (uc *PaperColor) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return dictionaries.FactoryErrPaperColorNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrPaperColorNotFound.New(id)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, dictionaries.PaperColorAPIName)
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
