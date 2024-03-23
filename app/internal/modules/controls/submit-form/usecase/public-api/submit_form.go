package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	SubmitForm struct {
		storage       SubmitFormStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewSubmitForm(
	storage SubmitFormStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *SubmitForm {
	return &SubmitForm{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *SubmitForm) GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error) {
	items, err := uc.storage.Fetch(ctx, params)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	return items, nil
}

func (uc *SubmitForm) GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error) {
	if rewriteName == "" {
		return entity.SubmitForm{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchByRewriteName(ctx, rewriteName)

	if err != nil {
		return entity.SubmitForm{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, rewriteName)
	}

	return item, nil
}
