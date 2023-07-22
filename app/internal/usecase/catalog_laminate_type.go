package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
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
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogLaminateType)
    }

    return items, nil
}

func (f *CatalogLaminateType) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogLaminateType, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogLaminateType{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminateType)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogLaminateType) Create(ctx context.Context, item *entity.CatalogLaminateType) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogLaminateType)
    }

    f.logger(ctx).Event("%s::Create: id=%d", entity.ModelNameCatalogLaminateType, item.Id)

    return nil
}

func (f *CatalogLaminateType) Store(ctx context.Context, item *entity.CatalogLaminateType) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminateType)
    }

    f.logger(ctx).Event("%s::Store: id=%d", entity.ModelNameCatalogLaminateType, item.Id)

    return nil
}

func (f *CatalogLaminateType) ChangeStatus(ctx context.Context, item *entity.CatalogLaminateType) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminateType)
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogLaminateType, item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminateType)
    }

    f.logger(ctx).Event("%s::ChangeStatus: id=%d, status=%s", entity.ModelNameCatalogLaminateType, item.Id, item.Status)

    return nil
}

func (f *CatalogLaminateType) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, entity.ModelNameCatalogLaminateType)
    }

    f.logger(ctx).Event("%s::Remove: id=%d", entity.ModelNameCatalogLaminateType, id)

    return nil
}

func (f *CatalogLaminateType) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
