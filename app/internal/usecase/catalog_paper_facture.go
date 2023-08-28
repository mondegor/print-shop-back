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

func NewCatalogPaperFacture(storage CatalogPaperFactureStorage,
                            errorHelper *mrerr.Helper) *CatalogPaperFacture {
    return &CatalogPaperFacture{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (uc *CatalogPaperFacture) GetList(ctx context.Context, listFilter *entity.CatalogPaperFactureListFilter) ([]entity.CatalogPaperFacture, error) {
    items := make([]entity.CatalogPaperFacture, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogPaperFacture)
    }

    return items, nil
}

func (uc *CatalogPaperFacture) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaperFacture, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogPaperFacture{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaperFacture)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogPaperFacture) Create(ctx context.Context, item *entity.CatalogPaperFacture) error {
    item.Status = entity.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogPaperFacture)
    }

    uc.logger(ctx).Event(
        "%s::Create: id=%d",
        entity.ModelNameCatalogPaperFacture,
        item.Id,
    )

    return nil
}

func (uc *CatalogPaperFacture) Store(ctx context.Context, item *entity.CatalogPaperFacture) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaperFacture)
    }

    uc.logger(ctx).Event(
        "%s::Store: id=%d",
        entity.ModelNameCatalogPaperFacture,
        item.Id,
    )

    return nil
}

func (uc *CatalogPaperFacture) ChangeStatus(ctx context.Context, item *entity.CatalogPaperFacture) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaperFacture)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogPaperFacture, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaperFacture)
    }

    uc.logger(ctx).Event(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogPaperFacture,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogPaperFacture) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.errorHelper.WrapErrorForRemove(err, entity.ModelNameCatalogPaperFacture)
    }

    uc.logger(ctx).Event(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogPaperFacture,
        id,
    )

    return nil
}

func (uc *CatalogPaperFacture) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
