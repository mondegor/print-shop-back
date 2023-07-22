package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
)

type FormFieldItem struct {
    componentOrderer ItemOrdererComponent
    storage FormFieldItemStorage
    storageFormFieldTemplate FormFieldTemplateStorage
    errorHelper *mrerr.Helper
}

func NewFormFieldItem(componentOrderer ItemOrdererComponent, storage FormFieldItemStorage, storageFormFieldTemplate FormFieldTemplateStorage, errorHelper *mrerr.Helper) *FormFieldItem {
    return &FormFieldItem{
        componentOrderer: componentOrderer,
        storage: storage,
        storageFormFieldTemplate: storageFormFieldTemplate,
        errorHelper: errorHelper,
    }
}

func (f *FormFieldItem) GetList(ctx context.Context, listFilter *entity.FormFieldItemListFilter) ([]entity.FormFieldItem, error) {
    items := make([]entity.FormFieldItem, 0, 4)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameFormFieldItem)
    }

    return items, nil
}

func (f *FormFieldItem) GetItem(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) (*entity.FormFieldItem, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.FormFieldItem{Id: id, FormId: formId}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, entity.ModelNameFormFieldItem)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *FormFieldItem) Create(ctx context.Context, item *entity.FormFieldItem) error {
    itemTemplate := entity.FormFieldTemplate{Id: item.TemplateId}
    err := f.storageFormFieldTemplate.LoadOne(ctx, &itemTemplate)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return ErrFormFieldItemTemplateNotFound.Wrap(err, item.TemplateId)
        }

        return err
    }

    if item.ParamName == "" {
        item.ParamName = itemTemplate.ParamName
    }

    if item.Caption == "" {
        item.Caption = itemTemplate.Caption
    }

    err = f.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    err = f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, entity.ModelNameFormFieldItem)
    }

    f.logger(ctx).Event("%s::Create: id=%d", entity.ModelNameFormFieldItem, item.Id)

    meta := f.storage.GetMetaData(item.FormId)
    component := f.componentOrderer.WithMetaData(meta)

    err = component.MoveToLast(
        ctx,
        item.Id,
    )

    if err != nil {
        f.logger(ctx).Error(err)
    }

    return nil
}

func (f *FormFieldItem) Store(ctx context.Context, item *entity.FormFieldItem) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := f.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    err = f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, entity.ModelNameFormFieldItem)
    }

    f.logger(ctx).Event("%s::Store: id=%d", entity.ModelNameFormFieldItem, item.Id)

    return nil
}

func (f *FormFieldItem) Remove(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := f.storage.Delete(ctx, id, formId)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, entity.ModelNameFormFieldItem)
    }

    f.logger(ctx).Event("%s::Remove: id=%d", entity.ModelNameFormFieldItem, id)

    return nil
}

func (f *FormFieldItem) MoveAfterId(ctx context.Context, id mrentity.KeyInt32, afterId mrentity.KeyInt32, formId mrentity.KeyInt32) error {
    if formId < 1 {
        return mrerr.ErrServiceIncorrectInputData.New(mrerr.Arg{"formId": formId})
    }

    meta := f.storage.GetMetaData(formId)
    component := f.componentOrderer.WithMetaData(meta)

    return component.MoveAfterId(ctx, id, afterId)
}

func (f *FormFieldItem) checkParamName(ctx context.Context, item *entity.FormFieldItem) error {
    id, err := f.storage.FetchIdByName(ctx, item)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameFormFieldItem)
    }

    if item.Id == id {
        return nil
    }

    return ErrFormFieldItemParamNameAlreadyExists.New(item.ParamName)
}

func (f *FormFieldItem) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
