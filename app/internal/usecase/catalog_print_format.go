package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

type CatalogPrintFormat struct {
    storage CatalogPrintFormatStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogPrintFormat(storage CatalogPrintFormatStorage,
                           errorHelper *mrerr.Helper) *CatalogPrintFormat {
    return &CatalogPrintFormat{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (uc *CatalogPrintFormat) GetList(ctx context.Context, listFilter *entity.CatalogPrintFormatListFilter) ([]entity.CatalogPrintFormat, error) {
    items := make([]entity.CatalogPrintFormat, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogPrintFormat)
    }

    return items, nil
}

func (uc *CatalogPrintFormat) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPrintFormat, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogPrintFormat{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPrintFormat)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogPrintFormat) Create(ctx context.Context, item *entity.CatalogPrintFormat) error {
    item.Status = entity.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.logger(ctx).Event(
        "%s::Create: id=%d",
        entity.ModelNameCatalogPrintFormat,
        item.Id,
    )

    return nil
}

func (uc *CatalogPrintFormat) Store(ctx context.Context, item *entity.CatalogPrintFormat) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.logger(ctx).Event(
        "%s::Store: id=%d",
        entity.ModelNameCatalogPrintFormat,
        item.Id,
    )

    return nil
}

func (uc *CatalogPrintFormat) ChangeStatus(ctx context.Context, item *entity.CatalogPrintFormat) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPrintFormat)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogPrintFormat, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.logger(ctx).Event(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogPrintFormat,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogPrintFormat) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.errorHelper.WrapErrorForRemove(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.logger(ctx).Event(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogPrintFormat,
        id,
    )

    return nil
}

func (uc *CatalogPrintFormat) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
