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

func NewFormData(storage FormDataStorage, errorHelper *mrerr.Helper) *FormData {
    return &FormData{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *FormData) GetList(ctx context.Context, listFilter *entity.FormDataListFilter) ([]entity.FormData, error) {
    items := make([]entity.FormData, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "FormData")
    }

    return items, nil
}

func (f *FormData) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormData, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.FormData{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "FormData")
    }

    return item, nil
}

func (f *FormData) CheckAvailability(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.IsExists(ctx, id)

    return f.errorHelper.ReturnErrorIfItemNotFound(err, "FormData")
}

// Create
// modifies: item{Id}
func (f *FormData) Create(ctx context.Context, item *entity.FormData) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "FormData")
    }

    f.logger(ctx).Event("FormData::Created: id=%d", item.Id)

    return nil
}

func (f *FormData) Store(ctx context.Context, item *entity.FormData) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "FormData")
    }

    f.logger(ctx).Event("FormData::Stored: id=%d", item.Id)

    return nil
}

func (f *FormData) ChangeStatus(ctx context.Context, item *entity.FormData) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "FormData")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "FormData", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "FormData")
    }

    f.logger(ctx).Event("FormData::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *FormData) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "FormData")
    }

    f.logger(ctx).Event("FormData::Removed: id=%d", id)

    return nil
}

func (f *FormData) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
