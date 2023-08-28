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

func NewCatalogLaminateType(storage CatalogLaminateTypeStorage,
                            errorHelper *mrerr.Helper) *CatalogLaminateType {
    return &CatalogLaminateType{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (uc *CatalogLaminateType) GetList(ctx context.Context, listFilter *entity.CatalogLaminateTypeListFilter) ([]entity.CatalogLaminateType, error) {
    items := make([]entity.CatalogLaminateType, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogLaminateType)
    }

    return items, nil
}

func (uc *CatalogLaminateType) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogLaminateType, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogLaminateType{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminateType)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogLaminateType) Create(ctx context.Context, item *entity.CatalogLaminateType) error {
    item.Status = entity.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogLaminateType)
    }

    uc.logger(ctx).Event(
        "%s::Create: id=%d",
        entity.ModelNameCatalogLaminateType,
        item.Id,
    )

    return nil
}

func (uc *CatalogLaminateType) Store(ctx context.Context, item *entity.CatalogLaminateType) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminateType)
    }

    uc.logger(ctx).Event(
        "%s::Store: id=%d",
        entity.ModelNameCatalogLaminateType,
        item.Id,
    )

    return nil
}

func (uc *CatalogLaminateType) ChangeStatus(ctx context.Context, item *entity.CatalogLaminateType) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminateType)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogLaminateType, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminateType)
    }

    uc.logger(ctx).Event(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogLaminateType,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogLaminateType) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.errorHelper.WrapErrorForRemove(err, entity.ModelNameCatalogLaminateType)
    }

    uc.logger(ctx).Event(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogLaminateType,
        id,
    )

    return nil
}

func (uc *CatalogLaminateType) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
