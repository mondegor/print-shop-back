package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"
	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"
	"github.com/mondegor/print-shop-back/pkg/provideraccounts/flow"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		storage      prov.CompanyPageStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
		imgBaseURL   mrpath.PathBuilder
		statusFlow   mrstatus.Flow
	}
)

// NewCompanyPage - создаёт объект CompanyPage.
func NewCompanyPage(
	storage prov.CompanyPageStorage,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UseCaseErrorWrapper,
	imgBaseURL mrpath.PathBuilder,
) *CompanyPage {
	return &CompanyPage{
		storage:      storage,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
		imgBaseURL:   imgBaseURL,
		statusFlow:   flow.PublicStatusFlow(),
	}
}

// GetItem - comment method.
func (uc *CompanyPage) GetItem(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error) {
	if accountID == uuid.Nil {
		return entity.CompanyPage{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, accountID)
	if err != nil {
		return entity.CompanyPage{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, accountID)
	}

	item.LogoURL = uc.imgBaseURL.BuildPath(item.LogoURL)

	return item, nil
}

// Store - comment method.
func (uc *CompanyPage) Store(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	item.Status = enum.PublicStatusDraft // only for insert

	if err := uc.storage.InsertOrUpdate(ctx, item); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"accountId": item.AccountID})

	return nil
}

// ChangeStatus - comment method.
func (uc *CompanyPage) ChangeStatus(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.AccountID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.ErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	if err = uc.storage.UpdateStatus(ctx, item); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"accountId": item.AccountID, "status": item.Status})

	return nil
}

func (uc *CompanyPage) checkItem(ctx context.Context, item *entity.CompanyPage) error {
	return uc.checkRewriteName(ctx, item)
}

func (uc *CompanyPage) checkRewriteName(ctx context.Context, item *entity.CompanyPage) error {
	accountID, err := uc.storage.FetchAccountIDByRewriteName(ctx, item.RewriteName)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	if item.AccountID != accountID {
		return module.ErrCompanyPageRewriteNameAlreadyExists.New(item.RewriteName)
	}

	return nil
}

func (uc *CompanyPage) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameCompanyPage,
		data,
	)
}
