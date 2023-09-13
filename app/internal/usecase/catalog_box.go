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

type CatalogBox struct {
    storage CatalogBoxStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewCatalogBox(storage CatalogBoxStorage,
                   eventBox mrcore.EventBox,
                   serviceHelper *mrtool.ServiceHelper) *CatalogBox {
    return &CatalogBox{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogBox) GetList(ctx context.Context, listFilter *entity.CatalogBoxListFilter) ([]entity.CatalogBox, error) {
    items := make([]entity.CatalogBox, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogBox)
    }

    return items, nil
}

func (uc *CatalogBox) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogBox, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogBox{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogBox)
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

    item.Status = mrcom.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogBox)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogBox,
        item.Id,
    )

    return nil
}

func (uc *CatalogBox) Store(ctx context.Context, item *entity.CatalogBox) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogBox)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogBox,
        item.Id,
    )

    return nil
}

func (uc *CatalogBox) ChangeStatus(ctx context.Context, item *entity.CatalogBox) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogBox)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogBox, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogBox)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogBox,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogBox) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogBox)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogBox,
        id,
    )

    return nil
}

func (uc *CatalogBox) checkArticle(ctx context.Context, item *entity.CatalogBox) error {
    id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogBox)
    }

    if item.Id == id {
        return nil
    }

    return ErrCatalogBoxArticleAlreadyExists.New(item.Article)
}
