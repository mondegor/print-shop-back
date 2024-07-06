package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/components/uiform"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// SubmitFormUseCase - comment interface.
	SubmitFormUseCase interface {
		GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, int64, error)
		GetItem(ctx context.Context, itemID uuid.UUID) (entity.SubmitForm, error)
		Create(ctx context.Context, item entity.SubmitForm) (uuid.UUID, error)
		Store(ctx context.Context, item entity.SubmitForm) error
		ChangeStatus(ctx context.Context, item entity.SubmitForm) error
		Remove(ctx context.Context, itemID uuid.UUID) error
	}

	// FormVersionUseCase - comment interface.
	FormVersionUseCase interface {
		GetItemJson(ctx context.Context, primary entity.PrimaryKey, pretty bool) ([]byte, error)
		PrepareForTest(ctx context.Context, formID uuid.UUID) error
		Publish(ctx context.Context, formID uuid.UUID) error
	}

	// SubmitFormComponent - comment interface.
	SubmitFormComponent interface {
		GetFormStatus(ctx context.Context, formID uuid.UUID) (mrenum.ItemStatus, error)
		GetFormWithElements(ctx context.Context, formID uuid.UUID) (entity.SubmitForm, error)
	}

	// FormCompilerComponent - comment interface.
	FormCompilerComponent interface {
		Compile(ctx context.Context, form entity.SubmitForm) (uiform.UIForm, error)
		CompileToBytes(ctx context.Context, form entity.SubmitForm) ([]byte, error)
	}

	// SubmitFormStorage - comment interface.
	SubmitFormStorage interface {
		NewSelectParams(params entity.SubmitFormParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.SubmitForm, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.SubmitForm, error)
		FetchIDByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error)
		FetchIDByParamName(ctx context.Context, paramName string) (uuid.UUID, error)
		FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.SubmitForm) (uuid.UUID, error)
		Update(ctx context.Context, row entity.SubmitForm) (int32, error)
		UpdateStatus(ctx context.Context, row entity.SubmitForm) (int32, error)
		Delete(ctx context.Context, rowID uuid.UUID) error
	}

	// FormVersionStorage - comment interface.
	FormVersionStorage interface {
		Fetch(ctx context.Context, formID uuid.UUID) ([]entity.FormVersion, error)
		FetchOne(ctx context.Context, primary entity.PrimaryKey) (entity.FormVersion, error)
		FetchOneLastVersion(ctx context.Context, formID uuid.UUID) (entity.FormVersionStatus, error)
		Insert(ctx context.Context, row entity.FormVersion) error
		Update(ctx context.Context, row entity.FormVersion) error
		UpdateStatus(ctx context.Context, row entity.FormVersionStatus, toStatus enum.ActivityStatus) error
	}
)
