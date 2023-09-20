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

type CatalogLaminate struct {
    storage CatalogLaminateStorage
    storageCatalogLaminateType CatalogLaminateTypeStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewCatalogLaminate(storage CatalogLaminateStorage,
                        storageCatalogLaminateType CatalogLaminateTypeStorage,
                        eventBox mrcore.EventBox,
                        serviceHelper *mrtool.ServiceHelper) *CatalogLaminate {
    return &CatalogLaminate{
        storage: storage,
        storageCatalogLaminateType: storageCatalogLaminateType,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogLaminate) GetList(ctx context.Context, listFilter *entity.CatalogLaminateListFilter) ([]entity.CatalogLaminate, error) {
    items := make([]entity.CatalogLaminate, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogLaminate)
    }

    return items, nil
}

func (uc *CatalogLaminate) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogLaminate, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogLaminate{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminate)
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
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return ErrCatalogLaminateTypeNotFound.Wrap(err, item.TypeId)
        }

        return err
    }

    item.Status = mrcom.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogLaminate)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogLaminate,
        item.Id,
    )

    return nil
}

func (uc *CatalogLaminate) Store(ctx context.Context, item *entity.CatalogLaminate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminate)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogLaminate,
        item.Id,
    )

    return nil
}

func (uc *CatalogLaminate) ChangeStatus(ctx context.Context, item *entity.CatalogLaminate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogLaminate)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogLaminate, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogLaminate)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogLaminate,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogLaminate) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogLaminate)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogLaminate,
        id,
    )

    return nil
}

func (uc *CatalogLaminate) checkArticle(ctx context.Context, item *entity.CatalogLaminate) error {
    id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogLaminate)
    }

    if item.Id == id {
        return nil
    }

    return ErrCatalogLaminateArticleAlreadyExists.New(item.Article)
}
