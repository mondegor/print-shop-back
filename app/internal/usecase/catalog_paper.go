package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

type CatalogPaper struct {
    storage CatalogPaperStorage
    storageCatalogPaperColor CatalogPaperColorStorage
    storageCatalogPaperFacture CatalogPaperFactureStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewCatalogPaper(storage CatalogPaperStorage, storageCatalogPaperColor CatalogPaperColorStorage, storageCatalogPaperFacture CatalogPaperFactureStorage, errorHelper *mrerr.Helper) *CatalogPaper {
    return &CatalogPaper{
        storage: storage,
        storageCatalogPaperColor: storageCatalogPaperColor,
        storageCatalogPaperFacture: storageCatalogPaperFacture,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *CatalogPaper) GetList(ctx context.Context, listFilter *entity.CatalogPaperListFilter) ([]entity.CatalogPaper, error) {
    items := make([]entity.CatalogPaper, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "CatalogPaper")
    }

    return items, nil
}

func (f *CatalogPaper) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaper, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.CatalogPaper{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "CatalogPaper")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *CatalogPaper) Create(ctx context.Context, item *entity.CatalogPaper) error {
    err := f.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = f.storageCatalogPaperColor.IsExists(ctx, item.ColorId)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return ErrCatalogPaperColorNotFound.Wrap(err, item.ColorId)
        }

        return err
    }

    err = f.storageCatalogPaperFacture.IsExists(ctx, item.FactureId)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return ErrCatalogPaperFactureNotFound.Wrap(err, item.FactureId)
        }

        return err
    }

    item.Status = entity.ItemStatusDraft
    err = f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "CatalogPaper")
    }

    f.logger(ctx).Event("CatalogPaper::Created: id=%d", item.Id)

    return nil
}

func (f *CatalogPaper) Store(ctx context.Context, item *entity.CatalogPaper) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPaper")
    }

    f.logger(ctx).Event("CatalogPaper::Stored: id=%d", item.Id)

    return nil
}

func (f *CatalogPaper) ChangeStatus(ctx context.Context, item *entity.CatalogPaper) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "CatalogPaper")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "CatalogPaper", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "CatalogPaper")
    }

    f.logger(ctx).Event("CatalogPaper::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *CatalogPaper) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "CatalogPaper")
    }

    f.logger(ctx).Event("CatalogPaper::Removed: id=%d", id)

    return nil
}

func (f *CatalogPaper) checkArticle(ctx context.Context, item *entity.CatalogPaper) error {
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

    return ErrCatalogPaperArticleAlreadyExists.New(item.Article)
}

func (f *CatalogPaper) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
