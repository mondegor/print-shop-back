package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/pkg/controls/api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		storage      ElementTemplateStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(storage ElementTemplateStorage, errorWrapper mrcore.UsecaseErrorWrapper) *ElementTemplate {
	return &ElementTemplate{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetItemHeader - comment method.
func (uc *ElementTemplate) GetItemHeader(ctx context.Context, itemID mrtype.KeyInt32) (api.ElementTemplateDTO, error) {
	uc.debugCmd(ctx, "GetHead", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return api.ElementTemplateDTO{}, api.ErrElementTemplateRequired.New()
	}

	item, err := uc.storage.FetchOneHead(ctx, itemID)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ElementTemplateDTO{}, api.ErrElementTemplateNotFound.New(itemID)
		}

		return api.ElementTemplateDTO{}, uc.errorWrapper.WrapErrorFailed(err, api.ElementTemplateHeaderName)
	}

	if item.Status == mrenum.ItemStatusDisabled {
		return api.ElementTemplateDTO{}, api.ErrElementTemplateIsDisabled.New(itemID)
	}

	return item.ElementTemplateDTO, nil
}

func (uc *ElementTemplate) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", api.ElementTemplateHeaderName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
