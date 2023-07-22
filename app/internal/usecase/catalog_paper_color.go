package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
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
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogPaperColor)
    }

    return items, nil
}

func (f *CatalogPaperColor) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaperColor, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogPaperColor{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaperColor)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogPaperColor) Create(ctx context.Context, item *entity.CatalogPaperColor) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogPaperColor)
    }

    f.logger(ctx).Event("%s::Create: id=%d", entity.ModelNameCatalogPaperColor, item.Id)

    return nil
}

func (f *CatalogPaperColor) Store(ctx context.Context, item *entity.CatalogPaperColor) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaperColor)
    }

    f.logger(ctx).Event("%s::Store: id=%d", entity.ModelNameCatalogPaperColor, item.Id)

    return nil
}

func (f *CatalogPaperColor) ChangeStatus(ctx context.Context, item *entity.CatalogPaperColor) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaperColor)
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogPaperColor, item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaperColor)
    }

    f.logger(ctx).Event("%s::ChangeStatus: id=%d, status=%s", entity.ModelNameCatalogPaperColor, item.Id, item.Status)

    return nil
}

func (f *CatalogPaperColor) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, entity.ModelNameCatalogPaperColor)
    }

    f.logger(ctx).Event("%s::Remove: id=%d", entity.ModelNameCatalogPaperColor, id)

    return nil
}

func (f *CatalogPaperColor) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
