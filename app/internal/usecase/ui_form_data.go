package usecase

import (
    "context"
    "encoding/json"
    "fmt"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
    "strings"
)

type UIFormData struct {
    storage FormDataStorage
    storageFormFieldItem FormFieldItemStorage
    errorHelper *mrerr.Helper
}

func NewUIFormData(storage FormDataStorage, storageFormFieldItem FormFieldItemStorage, errorHelper *mrerr.Helper) *UIFormData {
    return &UIFormData{
        storage: storage,
        storageFormFieldItem: storageFormFieldItem,
        errorHelper: errorHelper,
    }
}

func (f *UIFormData) CompileForm(ctx context.Context, id mrentity.KeyInt32) (*entity.UIForm, error) {
    if id < 1 {
        return nil, mrerr.ErrServiceIncorrectInputData.NewWithData("id=%d", id)
    }

    form := &entity.FormData{Id: id}
    err := f.storage.LoadOne(ctx, form)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "FormData")
    }

    listFilter := entity.FormFieldItemListFilter{FormId: id}
    items := make([]entity.FormFieldItem, 0, 4)
    err = f.storageFormFieldItem.LoadAll(ctx, &listFilter, &items)

    if err != nil {
        return nil, f.errorHelper.WrapErrorForSelect(err, "FormFieldItem")
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

                f.correctField(field.Id, &field)

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

        f.correctField(group.Id, &group)

        uiForm.Fields = append(uiForm.Fields, group)
    }

    return &uiForm, nil
}

func (f *UIFormData) correctField(parentName string, field *entity.UIFieldItem) {
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
        f.correctField(val.Id, val)
    }
}

//func (f *UIFormData) logger(ctx context.Context) mrapp.Logger {
//    return mrcontext.GetLogger(ctx)
//}
