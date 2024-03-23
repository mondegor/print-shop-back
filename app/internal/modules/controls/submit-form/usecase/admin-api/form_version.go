package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	module "print-shop-back/internal/modules/controls/submit-form"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	"print-shop-back/pkg/modules/controls/enums"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	FormVersion struct {
		storage       FormVersionStorage
		formComponent SubmitFormComponent
		formCompiler  FormCompilerComponent
		locker        mrlock.Locker
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewFormVersion(
	storage FormVersionStorage,
	formComponent SubmitFormComponent,
	formCompiler FormCompilerComponent,
	locker mrlock.Locker,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *FormVersion {
	return &FormVersion{
		storage:       storage,
		formComponent: formComponent,
		formCompiler:  formCompiler,
		locker:        locker,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *FormVersion) GetItemJson(ctx context.Context, primary entity.PrimaryKey, pretty bool) ([]byte, error) {
	if primary.FormID == uuid.Nil || primary.Version < 0 {
		return nil, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if _, err := uc.formComponent.GetFormStatus(ctx, primary.FormID); err != nil {
		return nil, err
	}

	var item entity.FormVersion
	var err error

	// при указании формы нулевой версии возвращается копия текущей формы
	if primary.Version == 0 {
		item, err = uc.createFormVersionForTest(ctx, primary.FormID)

		if err != nil {
			return nil, err
		}
	} else {
		// :TODO: можно оптимизировать получая только body
		item, err = uc.storage.FetchOne(ctx, primary)

		if err != nil {
			return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormVersion, primary)
		}
	}

	if pretty {
		var prettyJSON bytes.Buffer

		if err = json.Indent(&prettyJSON, item.Body, "", module.JsonPrettyIndent); err != nil {
			return nil, uc.usecaseHelper.WrapErrorEntityFailed(err, entity.ModelNameFormVersion, primary)
		}

		return prettyJSON.Bytes(), nil
	}

	return item.Body, nil
}

func (uc *FormVersion) PrepareForTest(ctx context.Context, formID uuid.UUID) error {
	lastVersionItem, err := uc.getItemLastVersion(ctx, formID)

	if err != nil {
		return err
	}

	if lastVersionItem.ActivityStatus == enums.ActivityStatusArchived {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(
			enums.ActivityStatusArchived,
			enums.ActivityStatusTesting,
		)
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(formID)); err != nil {
		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormVersion)
	} else {
		defer unlock()
	}

	item, err := uc.createFormVersionForTest(ctx, formID)

	if err != nil {
		return err
	}

	eventName := "Create"
	item.Version = lastVersionItem.Version

	if lastVersionItem.ActivityStatus == enums.ActivityStatusTesting {
		if err = uc.storage.Update(ctx, item); err != nil {
			return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormVersion, item.ID)
		}

		eventName = "Update"
	} else {
		if lastVersionItem.ActivityStatus == enums.ActivityStatusPublished {
			item.Version++
		}

		if err = uc.storage.Insert(ctx, item); err != nil {
			return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormVersion)
		}
	}

	uc.emitEvent(ctx, eventName, mrmsg.Data{"formId": item.ID, "version": item.Version})

	return nil
}

func (uc *FormVersion) Publish(ctx context.Context, formID uuid.UUID) error {
	item, err := uc.getItemLastVersion(ctx, formID)

	if err != nil {
		return err
	}

	if item.ActivityStatus == enums.ActivityStatusPublished {
		return nil // переключения не требуется
	}

	if item.ActivityStatus != enums.ActivityStatusTesting {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(
			item.ActivityStatus,
			enums.ActivityStatusPublished,
		)
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(formID)); err != nil {
		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormVersion)
	} else {
		defer unlock()
	}

	if err = uc.storage.UpdateStatus(ctx, item, enums.ActivityStatusPublished); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormVersion)
	}

	uc.emitEvent(ctx, "Publish", mrmsg.Data{"formId": formID, "version": item.Version})

	return nil
}

func (uc *FormVersion) getLockKey(formID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", entity.ModelNameFormVersion, formID)
}

func (uc *FormVersion) getItemLastVersion(ctx context.Context, formID uuid.UUID) (entity.FormVersionStatus, error) {
	if formID == uuid.Nil {
		return entity.FormVersionStatus{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if formStatus, err := uc.formComponent.GetFormStatus(ctx, formID); err != nil {
		return entity.FormVersionStatus{}, err
	} else if formStatus != mrenum.ItemStatusEnabled {
		return entity.FormVersionStatus{}, mrcore.FactoryErrUseCaseEntityNotAvailable.New()
	}

	item, err := uc.storage.FetchOneLastVersion(ctx, formID)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return entity.FormVersionStatus{
				FormID:         formID,
				Version:        1,
				ActivityStatus: enums.ActivityStatusDraft,
			}, nil
		}

		return entity.FormVersionStatus{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormVersion, formID)
	}

	return item, nil
}

func (uc *FormVersion) createFormVersionForTest(ctx context.Context, formID uuid.UUID) (entity.FormVersion, error) {
	form, err := uc.formComponent.GetFormWithElements(ctx, formID)

	if err != nil {
		return entity.FormVersion{}, err
	}

	body, err := uc.formCompiler.CompileToBytes(ctx, form)

	if err != nil {
		return entity.FormVersion{}, err
	}

	return entity.FormVersion{
		ID:             formID,
		Version:        1,
		RewriteName:    form.RewriteName,
		Caption:        form.Caption,
		Detailing:      form.Detailing,
		Body:           body,
		ActivityStatus: enums.ActivityStatusTesting,
	}, nil
}

func (uc *FormVersion) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameFormVersion,
		data,
	)
}
