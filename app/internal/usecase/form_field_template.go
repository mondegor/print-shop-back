package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
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
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameFormFieldTemplate)
    }

    return items, nil
}

func (f *FormFieldTemplate) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormFieldTemplate, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.FormFieldTemplate{Id: id}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, entity.ModelNameFormFieldTemplate)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *FormFieldTemplate) Create(ctx context.Context, item *entity.FormFieldTemplate) error {
    item.Status = entity.ItemStatusDraft
    err := f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameFormFieldTemplate)
    }

    f.logger(ctx).Event("%s::Create: id=%d", entity.ModelNameFormFieldTemplate, item.Id)

    return nil
}

func (f *FormFieldTemplate) Store(ctx context.Context, item *entity.FormFieldTemplate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, entity.ModelNameFormFieldTemplate)
    }

    f.logger(ctx).Event("%s::Store: id=%d", entity.ModelNameFormFieldTemplate, item.Id)

    return nil
}

func (f *FormFieldTemplate) ChangeStatus(ctx context.Context, item *entity.FormFieldTemplate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := f.storage.FetchStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForSelect(err, entity.ModelNameFormFieldTemplate)
    }

    if !f.statusFlow.Check(currentStatus, item.Status) {
        return mrerr.ErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameFormFieldTemplate, item.Id)
    }

    err = f.storage.UpdateStatus(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, entity.ModelNameFormFieldTemplate)
    }

    f.logger(ctx).Event("%s::ChangeStatus: id=%d, status=%s", entity.ModelNameFormFieldTemplate, item.Id, item.Status)

    return nil
}

func (f *FormFieldTemplate) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := f.storage.Delete(ctx, id)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, entity.ModelNameFormFieldTemplate)
    }

    f.logger(ctx).Event("%s::Remove: id=%d", entity.ModelNameFormFieldTemplate, id)

    return nil
}

func (f *FormFieldTemplate) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
