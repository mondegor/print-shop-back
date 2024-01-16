package usecase_api

import (
	"context"
	entity_api "print-shop-back/internal/modules/controls/entity/api"
	usecase_shared "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	APINameElementTemplate = "Controls.ElementTemplateAPI"
)

type (
	ElementTemplate struct {
		storage       ElementTemplateStorage
		serviceHelper *mrtool.ServiceHelper
	}

	ElementTemplateAPI interface {
		GetHead(ctx context.Context, id mrtype.KeyInt32) (*entity_api.ElementTemplateHead, error)
	}

	ElementTemplateStorage interface {
		FetchHead(ctx context.Context, id mrtype.KeyInt32) (*entity_api.ElementTemplateHead, error)
	}
)

func NewElementTemplate(
	storage ElementTemplateStorage,
	serviceHelper *mrtool.ServiceHelper,
) *ElementTemplate {
	return &ElementTemplate{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *ElementTemplate) GetHead(ctx context.Context, id mrtype.KeyInt32) (*entity_api.ElementTemplateHead, error) {
	uc.debugCmd(ctx, "GetHead", mrmsg.Data{"id": id})

	if id < 1 {
		return nil, usecase_shared.FactoryErrElementTemplateNotFound.New(id)
	}

	item, err := uc.storage.FetchHead(ctx, id)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return nil, usecase_shared.FactoryErrElementTemplateNotFound.New(id)
		}

		return nil, uc.serviceHelper.WrapErrorFailed(err, APINameElementTemplate)
	}

	return item, nil
}

func (uc *ElementTemplate) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrctx.Logger(ctx).Debug(
		"%s: cmd=%s, data=%s",
		APINameElementTemplate,
		command,
		data,
	)
}
