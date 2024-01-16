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
	PaperFacture struct {
		storage       PaperFactureStorage
		serviceHelper *mrtool.ServiceHelper
	}

	PaperFactureStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)

func NewPaperFacture(
	storage PaperFactureStorage,
	serviceHelper *mrtool.ServiceHelper,
) *PaperFacture {
	return &PaperFacture{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *PaperFacture) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return dictionaries.FactoryErrPaperFactureNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrPaperFactureNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "Dictionaries.PaperFactureAPI")
	}

	return nil
}

func (uc *PaperFacture) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrctx.Logger(ctx).Debug(
		"Dictionaries.PaperFactureAPI: cmd=%s, data=%s",
		command,
		data,
	)
}
