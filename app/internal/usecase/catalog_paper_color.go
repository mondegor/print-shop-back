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

type CatalogPaperColor struct {
    storage CatalogPaperColorStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewCatalogPaperColor(storage CatalogPaperColorStorage,
                          eventBox mrcore.EventBox,
                          serviceHelper *mrtool.ServiceHelper) *CatalogPaperColor {
    return &CatalogPaperColor{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogPaperColor) GetList(ctx context.Context, listFilter *entity.CatalogPaperColorListFilter) ([]entity.CatalogPaperColor, error) {
    items := make([]entity.CatalogPaperColor, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogPaperColor)
    }

    return items, nil
}

func (uc *CatalogPaperColor) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaperColor, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogPaperColor{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaperColor)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogPaperColor) Create(ctx context.Context, item *entity.CatalogPaperColor) error {
    item.Status = mrcom.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogPaperColor)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogPaperColor,
        item.Id,
    )

    return nil
}

func (uc *CatalogPaperColor) Store(ctx context.Context, item *entity.CatalogPaperColor) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaperColor)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogPaperColor,
        item.Id,
    )

    return nil
}

func (uc *CatalogPaperColor) ChangeStatus(ctx context.Context, item *entity.CatalogPaperColor) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPaperColor)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogPaperColor, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPaperColor)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogPaperColor,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogPaperColor) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogPaperColor)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogPaperColor,
        id,
    )

    return nil
}
