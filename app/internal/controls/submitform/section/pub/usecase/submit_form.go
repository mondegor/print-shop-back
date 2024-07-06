package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/entity"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// SubmitForm - comment struct.
	SubmitForm struct {
		storage      pub.SubmitFormStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewSubmitForm - создаёт объект SubmitForm.
func NewSubmitForm(storage pub.SubmitFormStorage, errorWrapper mrcore.UsecaseErrorWrapper) *SubmitForm {
	return &SubmitForm{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *SubmitForm) GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	return items, nil
}

// GetItemByRewriteName - comment method.
func (uc *SubmitForm) GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error) {
	if rewriteName == "" {
		return entity.SubmitForm{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchByRewriteName(ctx, rewriteName)
	if err != nil {
		return entity.SubmitForm{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, rewriteName)
	}

	return item, nil
}
