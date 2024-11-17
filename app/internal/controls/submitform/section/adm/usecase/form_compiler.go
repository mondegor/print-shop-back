package usecase

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"
	"github.com/mondegor/print-shop-back/pkg/libs/components/uiform"
)

type (
	// FormCompilerJson - comment struct.
	FormCompilerJson struct{}
)

// NewFormCompilerJson - создаёт объект FormCompilerJson.
func NewFormCompilerJson() *FormCompilerJson {
	return &FormCompilerJson{}
}

// Compile - comment method.
func (uc *FormCompilerJson) Compile(_ context.Context, form entity.SubmitForm) (uiform.UIForm, error) {
	uiForm := uiform.UIForm{
		ID:      form.ParamName,
		Caption: form.Caption,
		Fields:  make([]uiform.UIFieldItem, 0, len(form.Elements)),
	}

	for _, item := range form.Elements {
		paramName := uiForm.ID + " " + item.ParamName

		switch item.Type {
		case enum.ElementTypeList:
			var fields []uiform.UIFieldItem

			if err := json.Unmarshal(item.Body, &fields); err != nil {
				return uiform.UIForm{}, err
			}

			for i := range fields {
				field := fields[i]

				field.ID = strings.Replace(field.ID, "%parentId%", paramName, 1)

				if field.IsRequired == nil {
					field.IsRequired = item.Required
				}

				uc.correctField(field.ID, &field)

				uiForm.Fields = append(uiForm.Fields, field)
			}
		case enum.ElementTypeGroup:
			group := uiform.UIFieldItem{
				ID:         paramName,
				Caption:    item.Caption,
				View:       uiform.UIItemViewBlock,
				IsRequired: item.Required,
			}

			if err := json.Unmarshal(item.Body, &group.Values); err != nil {
				return uiform.UIForm{}, err
			}

			uc.correctField(group.ID, &group)

			uiForm.Fields = append(uiForm.Fields, group)
		}
	}

	return uiForm, nil
}

// CompileToBytes - comment method.
func (uc *FormCompilerJson) CompileToBytes(ctx context.Context, form entity.SubmitForm) ([]byte, error) {
	uiForm, err := uc.Compile(ctx, form)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(&uiForm)
	if err != nil {
		return nil, err // TODO: wrap
	}

	return body, nil
}

func (uc *FormCompilerJson) correctField(parentName string, field *uiform.UIFieldItem) {
	field.EnabledValues = nil // TODO: error if not nil

	if field.View == uiform.UIItemViewBlock {
		isNullOrRequired := field.IsRequired != nil && !*field.IsRequired

		if !isNullOrRequired {
			field.EnabledValues = []uiform.UIFieldItem{
				{ID: parentName + "_Disabled", IsChecked: mrtype.CastBoolToPointer(false)},
				{ID: parentName + "_Enabled", IsChecked: mrtype.CastBoolToPointer(true)},
			}
		}
	}

	// TODO: ОТЛАДИТЬ!!!!!!!!!!!!
	for i := range field.Values {
		val := &(field.Values)[i]
		val.ID = strings.Replace(val.ID, "%parentId%", parentName, 1)
		uc.correctField(val.ID, val)
	}
}
