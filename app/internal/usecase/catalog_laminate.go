package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

type CatalogLaminate struct {
    storage CatalogLaminateStorage
    storageCatalogLaminateType CatalogLaminateTypeStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogLaminate(storage CatalogLaminateStorage,
                        storageCatalogLaminateType CatalogLaminateTypeStorage,
                        errorHelper *mrerr.Helper) *CatalogLaminate {
    return &CatalogLaminate{
        storage: storage,
        storageCatalogLaminateType: storageCatalogLaminateType,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (uc *CatalogLaminate) GetList(ctx context.Context, listFilter *entity.CatalogLaminateListFilter) ([]entity.CatalogLaminate, error) {
    items := make([]entity.CatalogLaminate, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogLaminate)
    }

    return items, nil
}

func (uc *CatalogLaminate) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogLaminate, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogLaminate{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminate)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogLaminate) Create(ctx context.Context, item *entity.CatalogLaminate) error {
    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storageCatalogLaminateType.IsExists(ctx, item.TypeId)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return ErrCatalogLaminateTypeNotFound.Wrap(err, item.TypeId)
        }

        return err
    }

    item.Status = entity.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogLaminate)
    }

    uc.logger(ctx).Event("%s::Create: id=%d", entity.ModelNameCatalogLaminate, item.Id)

    return nil
}

func (uc *CatalogLaminate) Store(ctx context.Context, item *entity.CatalogLaminate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminate)
    }

    uc.logger(ctx).Event("%s::Store: id=%d", entity.ModelNameCatalogLaminate, item.Id)

    return nil
}

func (uc *CatalogLaminate) ChangeStatus(ctx context.Context, item *entity.CatalogLaminate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminate)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogLaminate, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminate)
    }

    uc.logger(ctx).Event("%s::ChangeStatus: id=%d, status=%s", entity.ModelNameCatalogLaminate, item.Id, item.Status)

    return nil
}

func (uc *CatalogLaminate) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.errorHelper.WrapErrorForRemove(err, entity.ModelNameCatalogLaminate)
    }

    uc.logger(ctx).Event("%s::Remove: id=%d", entity.ModelNameCatalogLaminate, id)

    return nil
}

func (uc *CatalogLaminate) checkArticle(ctx context.Context, item *entity.CatalogLaminate) error {
    id, err := uc.storage.FetchIdByArticle(ctx, item)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogLaminate)
    }

    if item.Id == id {
        return nil
    }

    return ErrCatalogLaminateArticleAlreadyExists.New(item.Article)
}

func (uc *CatalogLaminate) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
