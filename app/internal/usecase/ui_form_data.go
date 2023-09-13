package usecase

import (
    "context"
    "encoding/json"
    "fmt"
    "print-shop-back/internal/entity"
    "strings"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

type UIFormData struct {
    storage FormDataStorage
    storageFormFieldItem FormFieldItemStorage
    serviceHelper *mrtool.ServiceHelper
}

func NewUIFormData(storage FormDataStorage,
                   storageFormFieldItem FormFieldItemStorage,
                   serviceHelper *mrtool.ServiceHelper) *UIFormData {
    return &UIFormData{
        storage: storage,
        storageFormFieldItem: storageFormFieldItem,
        serviceHelper: serviceHelper,
    }
}

func (uc *UIFormData) CompileForm(ctx context.Context, id mrentity.KeyInt32) (*entity.UIForm, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    form := &entity.FormData{Id: id}
    err := uc.storage.LoadOne(ctx, form)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, "FormData")
    }

    listFilter := entity.FormFieldItemListFilter{FormId: id}
    items := make([]entity.FormFieldItem, 0, 4)
    err = uc.storageFormFieldItem.LoadAll(ctx, &listFilter, &items)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, "FormFieldItem")
    }

    uiForm := entity.UIForm{
        Id: form.ParamName,
        Name: form.Caption,
    }

    for _, item := range items {
        paramName := fmt.Sprintf("%s_%s", uiForm.Id, item.ParamName)

        if item.Type == entity.FormFieldTemplateTypeFields {
            var fields []entity.UIFieldItem
            err = json.Unmarshal([]byte(item.Body), &fields)

            if err != nil {
                return nil, err
            }

            for _, field := range fields {
                field.Id = strings.Replace(field.Id, "%parentId%", paramName, 1)

                if field.IsRequired == nil {
                    field.IsRequired = mrentity.BoolPointer(item.Required)
                }

                uc.correctField(field.Id, &field)

                uiForm.Fields = append(uiForm.Fields, field)
            }

            continue
        }

        // item.Type == entity.FormFieldTemplateTypeGroup
        group := entity.UIFieldItem{
            Id: paramName,
            Name: item.Caption,
            View: entity.UIItemViewBlock,
            IsRequired: mrentity.BoolPointer(item.Required),
        }

        err = json.Unmarshal([]byte(item.Body), &group.Values)

        if err != nil {
            return nil, err
        }

        uc.correctField(group.Id, &group)

        uiForm.Fields = append(uiForm.Fields, group)
    }

    return &uiForm, nil
}

func (uc *UIFormData) correctField(parentName string, field *entity.UIFieldItem) {
    field.EnabledValues = nil // :TODO: error if not nil

    if field.View == entity.UIItemViewBlock {
        isNullOrRequired := field.IsRequired != nil && !*field.IsRequired

        if !isNullOrRequired {
            field.EnabledValues = []entity.UIFieldItem{
                {Id: fmt.Sprintf("%s_Disabled", parentName), IsChecked: mrentity.BoolPointer(false)},
                {Id: fmt.Sprintf("%s_Enabled", parentName), IsChecked: mrentity.BoolPointer(true)},
            }
        }
    }

    for i := range field.Values {
        val := &(field.Values)[i]
        val.Id = strings.Replace(val.Id, "%parentId%", parentName, 1)
        uc.correctField(val.Id, val)
    }
}
