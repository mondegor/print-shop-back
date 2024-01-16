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
	PaperColor struct {
		storage       PaperColorStorage
		serviceHelper *mrtool.ServiceHelper
	}

	PaperColorStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)

func NewPaperColor(
	storage PaperColorStorage,
	serviceHelper *mrtool.ServiceHelper,
) *PaperColor {
	return &PaperColor{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *PaperColor) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return dictionaries.FactoryErrPaperColorNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrPaperColorNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "Dictionaries.PaperColorAPI")
	}

	return nil
}

func (uc *PaperColor) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrctx.Logger(ctx).Debug(
		"Dictionaries.PaperColorAPI: cmd=%s, data=%s",
		command,
		data,
	)
}
