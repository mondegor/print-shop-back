package usecase_api

import (
	"context"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PrintFormat struct {
		storage       PrintFormatStorage
		serviceHelper *mrtool.ServiceHelper
	}

	PrintFormatStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)

func NewPrintFormat(
	storage PrintFormatStorage,
	serviceHelper *mrtool.ServiceHelper,
) *PrintFormat {
	return &PrintFormat{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *PrintFormat) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return dictionaries.FactoryErrPrintFormatNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrPrintFormatNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "Dictionaries.PrintFormatAPI")
	}

	return nil
}

func (uc *PrintFormat) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrctx.Logger(ctx).Debug(
		"Dictionaries.PrintFormatAPI: cmd=%s, data=%s",
		command,
		data,
	)
}
