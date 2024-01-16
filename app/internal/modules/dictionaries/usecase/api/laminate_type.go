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
	LaminateType struct {
		storage       LaminateTypeStorage
		serviceHelper *mrtool.ServiceHelper
	}

	LaminateTypeStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)

func NewLaminateType(
	storage LaminateTypeStorage,
	serviceHelper *mrtool.ServiceHelper,
) *LaminateType {
	return &LaminateType{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *LaminateType) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return dictionaries.FactoryErrLaminateTypeNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return dictionaries.FactoryErrLaminateTypeNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "Dictionaries.LaminateTypeAPI")
	}

	return nil
}

func (uc *LaminateType) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrctx.Logger(ctx).Debug(
		"Dictionaries.LaminateTypeAPI: cmd=%s, data=%s",
		command,
		data,
	)
}
