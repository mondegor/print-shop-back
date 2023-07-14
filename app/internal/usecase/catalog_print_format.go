package usecase

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrcontext"
    "calc-user-data-back-adm/pkg/mrentity"
    "calc-user-data-back-adm/pkg/mrerr"
    "context"
)

type CatalogPrintFormat struct {
    storage CatalogPrintFormatStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogPrintFormat(storage CatalogPrintFormatStorage, errorHelper *mrerr.Helper) *CatalogPrintFormat {
    return &CatalogPrintFormat{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *CatalogPrintFormat) GetList(ctx context.Context, listFilter *entity.CatalogPrintFormatListFilter) ([]entity.CatalogPrintFormat, error) {
    items := make([]entity.CatalogPrintFormat, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "CatalogPrintFormat")
    }

    return items, nil
}

func (f *CatalogPrintFormat) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPrintFormat, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.CatalogPrintFormat{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "CatalogPrintFormat")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogPrintFormat) Create(ctx context.Context, item *entity.CatalogPrintFormat) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "CatalogPrintFormat")
    }

    f.logger(ctx).Event("CatalogPrintFormat::Created: id=%d", item.Id)

    return nil
}

func (f *CatalogPrintFormat) Store(ctx context.Context, item *entity.CatalogPrintFormat) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPrintFormat")
    }

    f.logger(ctx).Event("CatalogPrintFormat::Stored: id=%d", item.Id)

    return nil
}

func (f *CatalogPrintFormat) ChangeStatus(ctx context.Context, item *entity.CatalogPrintFormat) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "CatalogPrintFormat")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "CatalogPrintFormat", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPrintFormat")
    }

    f.logger(ctx).Event("CatalogPrintFormat::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *CatalogPrintFormat) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "CatalogPrintFormat")
    }

    f.logger(ctx).Event("CatalogPrintFormat::Removed: id=%d", id)

    return nil
}

func (f *CatalogPrintFormat) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
