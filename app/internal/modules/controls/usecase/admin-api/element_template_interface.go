package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplateService interface {
		GetList(ctx context.Context, params entity.ElementTemplateParams) ([]entity.ElementTemplate, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.ElementTemplate, error)
		Create(ctx context.Context, item *entity.ElementTemplate) error
		Store(ctx context.Context, item *entity.ElementTemplate) error
		ChangeStatus(ctx context.Context, item *entity.ElementTemplate) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	ElementTemplateStorage interface {
		NewFetchParams(params entity.ElementTemplateParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.ElementTemplate, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.ElementTemplate) error
		FetchStatus(ctx context.Context, row *entity.ElementTemplate) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.ElementTemplate) error
		Update(ctx context.Context, row *entity.ElementTemplate) (int32, error)
		UpdateStatus(ctx context.Context, row *entity.ElementTemplate) (int32, error)
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
