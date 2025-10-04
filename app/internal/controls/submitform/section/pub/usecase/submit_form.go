package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/entity"
)

type (
	// SubmitForm - comment struct.
	SubmitForm struct {
		storage      pub.SubmitFormStorage
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewSubmitForm - создаёт объект SubmitForm.
func NewSubmitForm(
	storage pub.SubmitFormStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *SubmitForm {
	return &SubmitForm{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNameSubmitForm),
	}
}

// GetList - comment method.
func (uc *SubmitForm) GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err)
	}

	return items, nil
}

// GetItemByRewriteName - comment method.
func (uc *SubmitForm) GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error) {
	if rewriteName == "" {
		return entity.SubmitForm{}, mr.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchByRewriteName(ctx, rewriteName)
	if err != nil {
		return entity.SubmitForm{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "rewriteName", rewriteName)
	}

	return item, nil
}
