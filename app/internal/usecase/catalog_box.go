package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

type CatalogBox struct {
    storage CatalogBoxStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogBox(storage CatalogBoxStorage,
                   errorHelper *mrerr.Helper) *CatalogBox {
    return &CatalogBox{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (uc *CatalogBox) GetList(ctx context.Context, listFilter *entity.CatalogBoxListFilter) ([]entity.CatalogBox, error) {
    items := make([]entity.CatalogBox, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogBox)
    }

    return items, nil
}

func (uc *CatalogBox) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogBox, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogBox{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogBox)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogBox) Create(ctx context.Context, item *entity.CatalogBox) error {
    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    item.Status = entity.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogBox)
    }

    uc.logger(ctx).Event("%s::Create: id=%d", entity.ModelNameCatalogBox, item.Id)

    return nil
}

func (uc *CatalogBox) Store(ctx context.Context, item *entity.CatalogBox) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogBox)
    }

    uc.logger(ctx).Event("%s::Store: id=%d", entity.ModelNameCatalogBox, item.Id)

    return nil
}

func (uc *CatalogBox) ChangeStatus(ctx context.Context, item *entity.CatalogBox) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameCatalogBox)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogBox, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogBox)
    }

    uc.logger(ctx).Event("%s::ChangeStatus: id=%d, status=%s", entity.ModelNameCatalogBox, item.Id, item.Status)

    return nil
}

func (uc *CatalogBox) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.errorHelper.WrapErrorForRemove(err, entity.ModelNameCatalogBox)
    }

    uc.logger(ctx).Event("%s::Remove: id=%d", entity.ModelNameCatalogBox, id)

    return nil
}

func (uc *CatalogBox) checkArticle(ctx context.Context, item *entity.CatalogBox) error {
    id, err := uc.storage.FetchIdByArticle(ctx, item)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogBox)
    }

    if item.Id == id {
        return nil
    }

    return ErrCatalogBoxArticleAlreadyExists.New(item.Article)
}

func (uc *CatalogBox) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
