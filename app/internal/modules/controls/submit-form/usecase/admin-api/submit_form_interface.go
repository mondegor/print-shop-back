package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	SubmitFormUseCase interface {
		GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, int64, error)
		GetItem(ctx context.Context, itemID uuid.UUID) (entity.SubmitForm, error)
		Create(ctx context.Context, item entity.SubmitForm) (uuid.UUID, error)
		Store(ctx context.Context, item entity.SubmitForm) error
		ChangeStatus(ctx context.Context, item entity.SubmitForm) error
		Remove(ctx context.Context, itemID uuid.UUID) error
	}

	//UISubmitFormUseCase interface {
	//    CompileForm(ctx context.Context, itemID mrtype.KeyInt32) (entity.UIForm, error)
	//}

	SubmitFormStorage interface {
		NewSelectParams(params entity.SubmitFormParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.SubmitForm, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.SubmitForm, error)
		FetchIdByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error)
		FetchIdByParamName(ctx context.Context, paramName string) (uuid.UUID, error)
		FetchStatus(ctx context.Context, row entity.SubmitForm) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, rowID uuid.UUID) error
		Insert(ctx context.Context, row entity.SubmitForm) (uuid.UUID, error)
		Update(ctx context.Context, row entity.SubmitForm) (int32, error)
		UpdateStatus(ctx context.Context, row entity.SubmitForm) (int32, error)
		Delete(ctx context.Context, rowID uuid.UUID) error
	}
)
