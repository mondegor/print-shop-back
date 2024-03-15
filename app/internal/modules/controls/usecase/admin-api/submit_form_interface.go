package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	SubmitFormUseCase interface {
		GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.SubmitForm, error)
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
		Create(ctx context.Context, item entity.SubmitForm) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.SubmitForm) error
		ChangeStatus(ctx context.Context, item entity.SubmitForm) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	//UISubmitFormUseCase interface {
	//    CompileForm(ctx context.Context, itemID mrtype.KeyInt32) (entity.UIForm, error)
	//}

	SubmitFormStorage interface {
		NewFetchParams(params entity.SubmitFormParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.SubmitForm, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.SubmitForm, error)
		FetchIdByName(ctx context.Context, paramName string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, row entity.SubmitForm) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.SubmitForm) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.SubmitForm) (int32, error)
		UpdateStatus(ctx context.Context, row entity.SubmitForm) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
