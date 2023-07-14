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

func NewCatalogBox(storage CatalogBoxStorage, errorHelper *mrerr.Helper) *CatalogBox {
    return &CatalogBox{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *CatalogBox) GetList(ctx context.Context, listFilter *entity.CatalogBoxListFilter) ([]entity.CatalogBox, error) {
    items := make([]entity.CatalogBox, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "CatalogBox")
    }

    return items, nil
}

func (f *CatalogBox) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogBox, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.CatalogBox{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "CatalogBox")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogBox) Create(ctx context.Context, item *entity.CatalogBox) error {
    err := f.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    item.Status = entity.ItemStatusDraft
    err = f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "CatalogBox")
    }

    f.logger(ctx).Event("CatalogBox::Created: id=%d", item.Id)

    return nil
}

func (f *CatalogBox) Store(ctx context.Context, item *entity.CatalogBox) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogBox")
    }

    f.logger(ctx).Event("CatalogBox::Stored: id=%d", item.Id)

    return nil
}

func (f *CatalogBox) ChangeStatus(ctx context.Context, item *entity.CatalogBox) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "CatalogBox")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "CatalogBox", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogBox")
    }

    f.logger(ctx).Event("CatalogBox::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *CatalogBox) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "CatalogBox")
    }

    f.logger(ctx).Event("CatalogBox::Removed: id=%d", id)

    return nil
}

func (f *CatalogBox) checkArticle(ctx context.Context, item *entity.CatalogBox) error {
    id, err := f.storage.FetchIdByArticle(ctx, item)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "CatalogBox")
    }

    if item.Id == id {
        return nil
    }

    return ErrCatalogBoxArticleAlreadyExists.New(item.Article)
}

func (f *CatalogBox) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
