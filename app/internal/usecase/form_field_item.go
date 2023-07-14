package usecase

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrcontext"
    "calc-user-data-back-adm/pkg/mrentity"
    "calc-user-data-back-adm/pkg/mrerr"
    "context"
)

type FormFieldItem struct {
    storage FormFieldItemStorage
    storageFormFieldTemplate FormFieldTemplateStorage
    errorHelper *mrerr.Helper
}

func NewFormFieldItem(storage FormFieldItemStorage, storageFormFieldTemplate FormFieldTemplateStorage, errorHelper *mrerr.Helper) *FormFieldItem {
    return &FormFieldItem{
        storage: storage,
        storageFormFieldTemplate: storageFormFieldTemplate,
        errorHelper: errorHelper,
    }
}

func (f *FormFieldItem) GetList(ctx context.Context, listFilter *entity.FormFieldItemListFilter) ([]entity.FormFieldItem, error) {
    items := make([]entity.FormFieldItem, 0, 4)
    err := f.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "FormFieldItem")
    }

    return items, nil
}

func (f *FormFieldItem) GetItem(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) (*entity.FormFieldItem, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    item := &entity.FormFieldItem{Id: id, FormId: formId}
    err := f.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "FormFieldItem")
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (f *FormFieldItem) Create(ctx context.Context, item *entity.FormFieldItem) error {
    err := f.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    itemTemplate := entity.FormFieldTemplate{Id: item.TemplateId}
    err = f.storageFormFieldTemplate.LoadOne(ctx, &itemTemplate)

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

    err = f.storage.Insert(ctx, item)

    if err != nil {
        return mrerr.ErrServiceEntityNotCreated.Wrap(err, "FormFieldItem")
    }

    f.logger(ctx).Event("FormFieldItem::Created: id=%d", item.Id)

    return nil
}

func (f *FormFieldItem) Store(ctx context.Context, item *entity.FormFieldItem) error {
    if item.Id < 1 || item.Version < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("item.Id=%d; item.Version=%d", item.Id, item.Version)
    }

    err := f.checkParamName(ctx, item)

    if err != nil {
        return err
    }

    err = f.storage.Update(ctx, item)

    if err != nil {
        return f.errorHelper.WrapErrorForUpdate(err, "FormFieldItem")
    }

    f.logger(ctx).Event("FormFieldItem::Stored: id=%d", item.Id)

    return nil
}

func (f *FormFieldItem) Remove(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error {
    if id < 1 {
        return mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    err := f.storage.Delete(ctx, id, formId)

    if err != nil {
        return f.errorHelper.WrapErrorForRemove(err, "FormFieldItem")
    }

    f.logger(ctx).Event("FormFieldItem::Removed: id=%d", id)

    return nil
}

func (f *FormFieldItem) checkParamName(ctx context.Context, item *entity.FormFieldItem) error {
    id, err := f.storage.FetchIdByName(ctx, item)

    if err != nil {
        if mrerr.ErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrerr.ErrServiceEntityTemporarilyUnavailable.Wrap(err, "FormFieldItem")
    }

    if item.Id == id {
        return nil
    }

    return ErrFormFieldItemParamNameAlreadyExists.New(item.ParamName)
}

func (f *FormFieldItem) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
