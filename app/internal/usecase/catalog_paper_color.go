package usecase

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrcontext"
    "calc-user-data-back-adm/pkg/mrentity"
    "calc-user-data-back-adm/pkg/mrerr"
    "context"
)

type CatalogPaperColor struct {
    storage CatalogPaperColorStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogPaperColor(storage CatalogPaperColorStorage, errorHelper *mrerr.Helper) *CatalogPaperColor {
    return &CatalogPaperColor{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *CatalogPaperColor) GetList(ctx context.Context, listFilter *entity.CatalogPaperColorListFilter) ([]entity.CatalogPaperColor, error) {
    items := make([]entity.CatalogPaperColor, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "CatalogPaperColor")
    }

    return items, nil
}

func (f *CatalogPaperColor) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaperColor, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.CatalogPaperColor{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "CatalogPaperColor")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogPaperColor) Create(ctx context.Context, item *entity.CatalogPaperColor) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "CatalogPaperColor")
    }

    f.logger(ctx).Event("CatalogPaperColor::Created: id=%d", item.Id)

    return nil
}

func (f *CatalogPaperColor) Store(ctx context.Context, item *entity.CatalogPaperColor) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPaperColor")
    }

    f.logger(ctx).Event("CatalogPaperColor::Stored: id=%d", item.Id)

    return nil
}

func (f *CatalogPaperColor) ChangeStatus(ctx context.Context, item *entity.CatalogPaperColor) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "CatalogPaperColor")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "CatalogPaperColor", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPaperColor")
    }

    f.logger(ctx).Event("CatalogPaperColor::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *CatalogPaperColor) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "CatalogPaperColor")
    }

    f.logger(ctx).Event("CatalogPaperColor::Removed: id=%d", id)

    return nil
}

func (f *CatalogPaperColor) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
