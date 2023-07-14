package usecase

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
    "context"
)

type CatalogLaminateType struct {
    storage CatalogLaminateTypeStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogLaminateType(storage CatalogLaminateTypeStorage, errorHelper *mrerr.Helper) *CatalogLaminateType {
    return &CatalogLaminateType{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *CatalogLaminateType) GetList(ctx context.Context, listFilter *entity.CatalogLaminateTypeListFilter) ([]entity.CatalogLaminateType, error) {
    items := make([]entity.CatalogLaminateType, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "CatalogLaminateType")
    }

    return items, nil
}

func (f *CatalogLaminateType) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogLaminateType, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.CatalogLaminateType{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "CatalogLaminateType")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogLaminateType) Create(ctx context.Context, item *entity.CatalogLaminateType) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "CatalogLaminateType")
    }

    f.logger(ctx).Event("CatalogLaminateType::Created: id=%d", item.Id)

    return nil
}

func (f *CatalogLaminateType) Store(ctx context.Context, item *entity.CatalogLaminateType) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogLaminateType")
    }

    f.logger(ctx).Event("CatalogLaminateType::Stored: id=%d", item.Id)

    return nil
}

func (f *CatalogLaminateType) ChangeStatus(ctx context.Context, item *entity.CatalogLaminateType) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "CatalogLaminateType")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "CatalogLaminateType", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogLaminateType")
    }

    f.logger(ctx).Event("CatalogLaminateType::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *CatalogLaminateType) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "CatalogLaminateType")
    }

    f.logger(ctx).Event("CatalogLaminateType::Removed: id=%d", id)

    return nil
}

func (f *CatalogLaminateType) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
