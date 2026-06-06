package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrpath"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/provideraccounts/module"
	"print-shop-back/internal/provideraccounts/section/prov"
	"print-shop-back/internal/provideraccounts/section/prov/entity"
	"print-shop-back/pkg/provideraccounts/enum/publicstatus"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		txManager     mrstorage.DBTxManager
		storage       prov.CompanyPageStorage
		imgBaseURL    mrpath.Builder
		eventEmitter  mrevent.Emitter
		errorWrapper  errors.Wrapper
		statusFlowMap workflow.FlowMap[publicstatus.Enum]
	}
)

// NewCompanyPage - создаёт объект CompanyPage.
func NewCompanyPage(
	txManager mrstorage.DBTxManager,
	storage prov.CompanyPageStorage,
	imgBaseURL mrpath.Builder,
	eventEmitter mrevent.Emitter,
) *CompanyPage {
	return &CompanyPage{
		txManager:     txManager,
		storage:       storage,
		imgBaseURL:    imgBaseURL,
		eventEmitter:  mrevent.EmitterWithSource(eventEmitter, entity.ModelNameCompanyPage),
		errorWrapper:  errors.NewServiceRecordNotFoundWrapper(),
		statusFlowMap: publicstatus.NewFlowMap(),
	}
}

// GetItem - comment method.
func (uc *CompanyPage) GetItem(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error) {
	if accountID == uuid.Nil {
		return entity.CompanyPage{}, errors.ErrRecordNotFound
	}

	item, err := uc.storage.FetchOne(ctx, accountID)
	if err != nil {
		return entity.CompanyPage{}, uc.errorWrapper.Wrap(err, "accountId", accountID)
	}

	item.LogoURL = uc.imgBaseURL.BuildPath(item.LogoURL)

	return item, nil
}

// Save - comment method.
func (uc *CompanyPage) Save(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return errors.ErrRecordNotFound
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	item.Status = publicstatus.Draft // only for insert

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		if err := uc.storage.InsertOrUpdate(ctx, item); err != nil {
			return uc.errorWrapper.Wrap(err, "accountId", item.AccountID)
		}

		uc.eventEmitter.Emit(ctx, "Store", "accountId", item.AccountID)

		return nil
	})
}

// ChangeStatus - comment method.
func (uc *CompanyPage) ChangeStatus(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return errors.ErrRecordNotFound
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.AccountID)
	if err != nil {
		return uc.errorWrapper.Wrap(err, "accountId", item.AccountID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlowMap.IsPossible(currentStatus, item.Status) {
		return errors.ErrSwitchStatusRejected.New(currentStatus, item.Status)
	}

	if err = uc.storage.UpdateStatus(ctx, item); err != nil {
		return uc.errorWrapper.Wrap(err, "accountId", item.AccountID)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", "accountId", item.AccountID, "status", item.Status)

	return nil
}

func (uc *CompanyPage) checkItem(ctx context.Context, item *entity.CompanyPage) error {
	return uc.checkRewriteName(ctx, item)
}

func (uc *CompanyPage) checkRewriteName(ctx context.Context, item *entity.CompanyPage) error {
	accountID, err := uc.storage.FetchAccountIDByRewriteName(ctx, item.RewriteName)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return nil
		}

		return uc.errorWrapper.Wrap(err)
	}

	if item.AccountID != accountID {
		return module.ErrCompanyPageRewriteNameAlreadyExists.New(item.RewriteName)
	}

	return nil
}
