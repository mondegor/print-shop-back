package usecase

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
    "context"
)

type FormFieldTemplate struct {
    storage FormFieldTemplateStorage
    errorHelper *mrerr.Helper
    statusFlow entity.ItemStatusFlow
}

func NewFormFieldTemplate(storage FormFieldTemplateStorage, errorHelper *mrerr.Helper) *FormFieldTemplate {
    return &FormFieldTemplate{
        storage: storage,
        errorHelper: errorHelper,
        statusFlow: entity.ItemStatusFlowDefault,
    }
}

func (f *FormFieldTemplate) GetList(ctx context.Context, listFilter *entity.FormFieldTemplateListFilter) ([]entity.FormFieldTemplate, error) {
    items := make([]entity.FormFieldTemplate, 0, 16)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "FormFieldTemplate")
    }

    return items, nil
}

func (f *FormFieldTemplate) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormFieldTemplate, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.FormFieldTemplate{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "FormFieldTemplate")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *FormFieldTemplate) Create(ctx context.Context, item *entity.FormFieldTemplate) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "FormFieldTemplate")
    }

    f.logger(ctx).Event("FormFieldTemplate::Created: id=%d", item.Id)

    return nil
}

func (f *FormFieldTemplate) Store(ctx context.Context, item *entity.FormFieldTemplate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "FormFieldTemplate")
    }

    f.logger(ctx).Event("FormFieldTemplate::Stored: id=%d", item.Id)

    return nil
}

func (f *FormFieldTemplate) ChangeStatus(ctx context.Context, item *entity.FormFieldTemplate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, "FormFieldTemplate")
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, "FormFieldTemplate", item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "FormFieldTemplate")
    }

    f.logger(ctx).Event("FormFieldTemplate::StatusChanged: id=%d, status=%s", item.Id, item.Status)

    return nil
}

func (f *FormFieldTemplate) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "FormFieldTemplate")
    }

    f.logger(ctx).Event("FormFieldTemplate::Removed: id=%d", id)

    return nil
}

func (f *FormFieldTemplate) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
