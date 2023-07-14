package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

type CatalogPaperFacture struct {
    storage CatalogPaperFactureStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogPaperFacture(storage CatalogPaperFactureStorage, errorHelper *mrerr.Helper) *CatalogPaperFacture {
    return &CatalogPaperFacture{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *CatalogPaperFacture) GetList(ctx context.Context, listFilter *entity.CatalogPaperFactureListFilter) ([]entity.CatalogPaperFacture, error) {
    items := make([]entity.CatalogPaperFacture, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "CatalogPaperFacture")
    }

    return items, nil
}

func (f *CatalogPaperFacture) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaperFacture, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.CatalogPaperFacture{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "CatalogPaperFacture")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogPaperFacture) Create(ctx context.Context, item *entity.CatalogPaperFacture) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "CatalogPaperFacture")
    }

    f.logger(ctx).Event("CatalogPaperFacture::Created: id=%d", item.Id)

    return nil
}

func (f *CatalogPaperFacture) Store(ctx context.Context, item *entity.CatalogPaperFacture) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPaperFacture")
    }

    f.logger(ctx).Event("CatalogPaperFacture::Stored: id=%d", item.Id)

    return nil
}

func (f *CatalogPaperFacture) ChangeStatus(ctx context.Context, item *entity.CatalogPaperFacture) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "CatalogPaperFacture")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "CatalogPaperFacture", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPaperFacture")
    }

    f.logger(ctx).Event("CatalogPaperFacture::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *CatalogPaperFacture) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "CatalogPaperFacture")
    }

    f.logger(ctx).Event("CatalogPaperFacture::Removed: id=%d", id)

    return nil
}

func (f *CatalogPaperFacture) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
