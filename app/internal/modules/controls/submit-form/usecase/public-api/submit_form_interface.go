package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/public-api"
)

type (
	SubmitFormUseCase interface {
		GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error)
		GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error)
	}

	SubmitFormStorage interface {
		Fetch(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error)
		FetchByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error)
	}
)
