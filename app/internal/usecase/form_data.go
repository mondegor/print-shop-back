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

type FormData struct {
    storage FormDataStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewFormData(storage FormDataStorage,
                 eventBox mrcore.EventBox,
                 serviceHelper *mrtool.ServiceHelper) *FormData {
    return &FormData{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *FormData) GetList(ctx context.Context, listFilter *entity.FormDataListFilter) ([]entity.FormData, error) {
    items := make([]entity.FormData, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameFormData)
    }

    return items, nil
}

func (uc *FormData) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormData, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.FormData{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameFormData)
    }

    return item, nil
}

func (uc *FormData) CheckAvailability(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.IsExists(ctx, id)

    return uc.serviceHelper.ReturnErrorIfItemNotFound(err, entity.ModelNameFormData)
}

// Create
// modifies: item{Id}
func (uc *FormData) Create(ctx context.Context, item *entity.FormData) error {
    err := uc.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    item.Status = mrcom.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameFormData)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameFormData,
        item.Id,
    )

    return nil
}

func (uc *FormData) Store(ctx context.Context, item *entity.FormData) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    err := uc.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameFormData)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameFormData,
        item.Id,
    )

    return nil
}

func (uc *FormData) ChangeStatus(ctx context.Context, item *entity.FormData) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameFormData)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameFormData, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameFormData)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameFormData,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *FormData) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameFormData)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameFormData,
        id,
    )

    return nil
}

func (uc *FormData) checkParamName(ctx context.Context, item *entity.FormData) error {
    id, err := uc.storage.FetchIdByName(ctx, item.ParamName)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameFormData)
    }

    if item.Id == id {
        return nil
    }

    return ErrFormFieldItemParamNameAlreadyExists.New(item.ParamName)
}
