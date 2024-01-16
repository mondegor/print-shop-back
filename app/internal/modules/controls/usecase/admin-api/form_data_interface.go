package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormDataService interface {
		GetList(ctx context.Context, params entity.FormDataParams) ([]entity.FormData, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.FormData, error)
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
		Create(ctx context.Context, item *entity.FormData) error
		Store(ctx context.Context, item *entity.FormData) error
		ChangeStatus(ctx context.Context, item *entity.FormData) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	//UIFormDataService interface {
	//    CompileForm(ctx context.Context, id mrtype.KeyInt32) (*entity.UIForm, error)
	//}

	FormDataStorage interface {
		NewFetchParams(params entity.FormDataParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.FormData, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.FormData) error
		FetchIdByName(ctx context.Context, paramName string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, row *entity.FormData) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.FormData) error
		Update(ctx context.Context, row *entity.FormData) (int32, error)
		UpdateStatus(ctx context.Context, row *entity.FormData) (int32, error)
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
