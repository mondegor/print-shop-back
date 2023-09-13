package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

type CatalogPaper struct {
    storage CatalogPaperStorage
    storageCatalogPaperColor CatalogPaperColorStorage
    storageCatalogPaperFacture CatalogPaperFactureStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewCatalogPaper(storage CatalogPaperStorage,
                     storageCatalogPaperColor CatalogPaperColorStorage,
                     storageCatalogPaperFacture CatalogPaperFactureStorage,
                     eventBox mrcore.EventBox,
                     serviceHelper *mrtool.ServiceHelper) *CatalogPaper {
    return &CatalogPaper{
        storage: storage,
        storageCatalogPaperColor: storageCatalogPaperColor,
        storageCatalogPaperFacture: storageCatalogPaperFacture,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogPaper) GetList(ctx context.Context, listFilter *entity.CatalogPaperListFilter) ([]entity.CatalogPaper, error) {
    items := make([]entity.CatalogPaper, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogPaper)
    }

    return items, nil
}

func (uc *CatalogPaper) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaper, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogPaper{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaper)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogPaper) Create(ctx context.Context, item *entity.CatalogPaper) error {
    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storageCatalogPaperColor.IsExists(ctx, item.ColorId)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return ErrCatalogPaperColorNotFound.Wrap(err, item.ColorId)
        }

        return err
    }

    err = uc.storageCatalogPaperFacture.IsExists(ctx, item.FactureId)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return ErrCatalogPaperFactureNotFound.Wrap(err, item.FactureId)
        }

        return err
    }

    item.Status = mrcom.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogPaper)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogPaper,
        item.Id,
    )

    return nil
}

func (uc *CatalogPaper) Store(ctx context.Context, item *entity.CatalogPaper) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaper)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogPaper,
        item.Id,
    )

    return nil
}

func (uc *CatalogPaper) ChangeStatus(ctx context.Context, item *entity.CatalogPaper) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaper)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogPaper, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaper)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogPaper,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogPaper) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogPaper)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogPaper,
        id,
    )

    return nil
}

func (uc *CatalogPaper) checkArticle(ctx context.Context, item *entity.CatalogPaper) error {
    id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogPaper)
    }

    if item.Id == id {
        return nil
    }

    return ErrCatalogPaperArticleAlreadyExists.New(item.Article)
}
