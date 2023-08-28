package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

type FormData struct {
    storage FormDataStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewFormData(storage FormDataStorage,
                 errorHelper *mrerr.Helper) *FormData {
    return &FormData{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (uc *FormData) GetList(ctx context.Context, listFilter *entity.FormDataListFilter) ([]entity.FormData, error) {
    items := make([]entity.FormData, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameFormData)
    }

    return items, nil
}

func (uc *FormData) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormData, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.FormData{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameFormData)
    }

    return item, nil
}

func (uc *FormData) CheckAvailability(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.IsExists(ctx, id)

    return uc.errorHelper.ReturnErrorIfItemNotFound(err, entity.ModelNameFormData)
}

// Create
// modifies: item{Id}
func (uc *FormData) Create(ctx context.Context, item *entity.FormData) error {
    err := uc.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    item.Status = entity.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameFormData)
    }

    uc.logger(ctx).Event(
        "%s::Create: id=%d",
        entity.ModelNameFormData,
        item.Id,
    )

    return nil
}

func (uc *FormData) Store(ctx context.Context, item *entity.FormData) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameFormData)
    }

    uc.logger(ctx).Event(
        "%s::Store: id=%d",
        entity.ModelNameFormData,
        item.Id,
    )

    return nil
}

func (uc *FormData) ChangeStatus(ctx context.Context, item *entity.FormData) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForSelect(err, entity.ModelNameFormData)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameFormData, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.errorHelper.WrapErrorForUpdate(err, entity.ModelNameFormData)
    }

    uc.logger(ctx).Event(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameFormData,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *FormData) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.errorHelper.WrapErrorForRemove(err, entity.ModelNameFormData)
    }

    uc.logger(ctx).Event(
        "%s::Remove: id=%d",
        entity.ModelNameFormData,
        id,
    )

    return nil
}

func (uc *FormData) checkParamName(ctx context.Context, item *entity.FormData) error {
    id, err := uc.storage.FetchIdByName(ctx, item.ParamName)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameFormData)
    }

    if item.Id == id {
        return nil
    }

    return ErrFormFieldItemParamNameAlreadyExists.New(item.ParamName)
}

func (uc *FormData) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
