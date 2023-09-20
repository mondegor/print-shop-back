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

type FormFieldTemplate struct {
    storage FormFieldTemplateStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewFormFieldTemplate(storage FormFieldTemplateStorage,
                          eventBox mrcore.EventBox,
                          serviceHelper *mrtool.ServiceHelper) *FormFieldTemplate {
    return &FormFieldTemplate{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *FormFieldTemplate) GetList(ctx context.Context, listFilter *entity.FormFieldTemplateListFilter) ([]entity.FormFieldTemplate, error) {
    items := make([]entity.FormFieldTemplate, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameFormFieldTemplate)
    }

    return items, nil
}

func (uc *FormFieldTemplate) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormFieldTemplate, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.FormFieldTemplate{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameFormFieldTemplate)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *FormFieldTemplate) Create(ctx context.Context, item *entity.FormFieldTemplate) error {
    item.Status = mrcom.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameFormFieldTemplate)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameFormFieldTemplate,
        item.Id,
    )

    return nil
}

func (uc *FormFieldTemplate) Store(ctx context.Context, item *entity.FormFieldTemplate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameFormFieldTemplate)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameFormFieldTemplate,
        item.Id,
    )

    return nil
}

func (uc *FormFieldTemplate) ChangeStatus(ctx context.Context, item *entity.FormFieldTemplate) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameFormFieldTemplate)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameFormFieldTemplate, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameFormFieldTemplate)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameFormFieldTemplate,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *FormFieldTemplate) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameFormFieldTemplate)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameFormFieldTemplate,
        id,
    )

    return nil
}
