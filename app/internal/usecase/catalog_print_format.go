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

type CatalogPrintFormat struct {
    storage CatalogPrintFormatStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewCatalogPrintFormat(storage CatalogPrintFormatStorage,
                           eventBox mrcore.EventBox,
                           serviceHelper *mrtool.ServiceHelper) *CatalogPrintFormat {
    return &CatalogPrintFormat{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogPrintFormat) GetList(ctx context.Context, listFilter *entity.CatalogPrintFormatListFilter) ([]entity.CatalogPrintFormat, error) {
    items := make([]entity.CatalogPrintFormat, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogPrintFormat)
    }

    return items, nil
}

func (uc *CatalogPrintFormat) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPrintFormat, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogPrintFormat{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPrintFormat)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogPrintFormat) Create(ctx context.Context, item *entity.CatalogPrintFormat) error {
    item.Status = mrcom.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogPrintFormat,
        item.Id,
    )

    return nil
}

func (uc *CatalogPrintFormat) Store(ctx context.Context, item *entity.CatalogPrintFormat) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogPrintFormat,
        item.Id,
    )

    return nil
}

func (uc *CatalogPrintFormat) ChangeStatus(ctx context.Context, item *entity.CatalogPrintFormat) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogPrintFormat)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogPrintFormat, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogPrintFormat,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogPrintFormat) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogPrintFormat)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogPrintFormat,
        id,
    )

    return nil
}
