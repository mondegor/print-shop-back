package usecase_api

import (
	"context"
	"print-shop-back/pkg/modules/controls"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplate struct {
		storage       ElementTemplateStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewElementTemplate(
	storage ElementTemplateStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *ElementTemplate {
	return &ElementTemplate{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *ElementTemplate) GetItemHead(ctx context.Context, itemID mrtype.KeyInt32) (controls.ElementTemplateHead, error) {
	uc.debugCmd(ctx, "GetHead", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return controls.ElementTemplateHead{}, controls.FactoryErrElementTemplateRequired.New()
	}

	item, err := uc.storage.FetchOneHead(ctx, itemID)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return controls.ElementTemplateHead{}, controls.FactoryErrElementTemplateNotFound.New(itemID)
		}

		return controls.ElementTemplateHead{}, uc.usecaseHelper.WrapErrorFailed(err, controls.ElementTemplateAPIName)
	}

	if item.Status == mrenum.ItemStatusDisabled {
		return controls.ElementTemplateHead{}, controls.FactoryErrElementTemplateIsDisabled.New(itemID)
	}

	return item.ElementTemplateHead, nil
}

func (uc *ElementTemplate) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", controls.ElementTemplateAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
