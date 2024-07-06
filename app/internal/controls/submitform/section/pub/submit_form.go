package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/entity"
)

type (
	// SubmitFormUseCase - comment interface.
	SubmitFormUseCase interface {
		GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error)
		GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error)
	}

	// SubmitFormStorage - comment interface.
	SubmitFormStorage interface {
		Fetch(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error)
		FetchByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error)
	}
)
