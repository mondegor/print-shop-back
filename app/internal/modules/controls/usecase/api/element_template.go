package usecase_api

import (
	"context"
	entity_api "print-shop-back/internal/modules/controls/entity/api"
	usecase_shared "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	elementTemplateAPIName = "Controls.ElementTemplateAPI"
)

type (
	ElementTemplate struct {
		storage       ElementTemplateStorage
		usecaseHelper *mrcore.UsecaseHelper
	}

	ElementTemplateAPI interface {
		GetHead(ctx context.Context, itemID mrtype.KeyInt32) (entity_api.ElementTemplateHead, error)
	}

	ElementTemplateStorage interface {
		FetchHead(ctx context.Context, rowID mrtype.KeyInt32) (entity_api.ElementTemplateHead, error)
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

func (uc *ElementTemplate) GetHead(ctx context.Context, itemID mrtype.KeyInt32) (entity_api.ElementTemplateHead, error) {
	uc.debugCmd(ctx, "GetHead", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return entity_api.ElementTemplateHead{}, usecase_shared.FactoryErrElementTemplateNotFound.New(itemID)
	}

	item, err := uc.storage.FetchHead(ctx, itemID)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return entity_api.ElementTemplateHead{}, usecase_shared.FactoryErrElementTemplateNotFound.New(itemID)
		}

		return entity_api.ElementTemplateHead{}, uc.usecaseHelper.WrapErrorFailed(err, elementTemplateAPIName)
	}

	return item, nil
}

func (uc *ElementTemplate) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", elementTemplateAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
