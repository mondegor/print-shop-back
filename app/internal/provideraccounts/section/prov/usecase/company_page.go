package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-sysmess/mrstatus"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"
	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum/publicstatus"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		txManager     mrstorage.DBTxManager
		storage       prov.CompanyPageStorage
		imgBaseURL    mrpath.Builder
		eventEmitter  mrevent.Emitter
		errorWrapper  mrerr.UseCaseErrorWrapper
		statusFlowMap mrstatus.FlowMap[publicstatus.Enum]
	}
)

// NewCompanyPage - создаёт объект CompanyPage.
func NewCompanyPage(
	txManager mrstorage.DBTxManager,
	storage prov.CompanyPageStorage,
	imgBaseURL mrpath.Builder,
	eventEmitter mrevent.Emitter,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *CompanyPage {
	return &CompanyPage{
		txManager:     txManager,
		storage:       storage,
		imgBaseURL:    imgBaseURL,
		eventEmitter:  mrevent.NewSourceEmitter(eventEmitter, entity.ModelNameCompanyPage),
		errorWrapper:  mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNameCompanyPage),
		statusFlowMap: publicstatus.NewFlowMap(),
	}
}

// GetItem - comment method.
func (uc *CompanyPage) GetItem(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error) {
	if accountID == uuid.Nil {
		return entity.CompanyPage{}, mr.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, accountID)
	if err != nil {
		return entity.CompanyPage{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "accountId", accountID)
	}

	item.LogoURL = uc.imgBaseURL.BuildPath(item.LogoURL)

	return item, nil
}

// Store - comment method.
func (uc *CompanyPage) Store(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	item.Status = publicstatus.Draft // only for insert

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		if err := uc.storage.InsertOrUpdate(ctx, item); err != nil {
			return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "accountId", item.AccountID)
		}

		uc.eventEmitter.Emit(ctx, "Store", mrargs.Group{"accountId": item.AccountID})

		return nil
	})
}

// ChangeStatus - comment method.
func (uc *CompanyPage) ChangeStatus(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.AccountID)
	if err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "accountId", item.AccountID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlowMap.IsPossible(currentStatus, item.Status) {
		return mr.ErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	if err = uc.storage.UpdateStatus(ctx, item); err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "accountId", item.AccountID)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", mrargs.Group{"accountId": item.AccountID, "status": item.Status})

	return nil
}

func (uc *CompanyPage) checkItem(ctx context.Context, item *entity.CompanyPage) error {
	return uc.checkRewriteName(ctx, item)
}

func (uc *CompanyPage) checkRewriteName(ctx context.Context, item *entity.CompanyPage) error {
	accountID, err := uc.storage.FetchAccountIDByRewriteName(ctx, item.RewriteName)
	if err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	}

	if item.AccountID != accountID {
		return module.ErrCompanyPageRewriteNameAlreadyExists.New(item.RewriteName)
	}

	return nil
}
