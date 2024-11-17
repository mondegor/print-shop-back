package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"
)

type (
	// FormVersion - comment struct.
	// FormVersion - comment struct.
	FormVersion struct {
		storage       adm.FormVersionStorage
		formComponent adm.SubmitFormComponent
		formCompiler  adm.FormCompilerComponent
		locker        mrlock.Locker
		eventEmitter  mrsender.EventEmitter
		errorWrapper  mrcore.UseCaseErrorWrapper
	}
)

// NewFormVersion - создаёт объект FormVersion.
func NewFormVersion(
	storage adm.FormVersionStorage,
	formComponent adm.SubmitFormComponent,
	formCompiler adm.FormCompilerComponent,
	locker mrlock.Locker,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UseCaseErrorWrapper,
) *FormVersion {
	return &FormVersion{
		storage:       storage,
		formComponent: formComponent,
		formCompiler:  formCompiler,
		locker:        locker,
		eventEmitter:  decorator.NewSourceEmitter(eventEmitter, entity.ModelNameFormVersion),
		errorWrapper:  errorWrapper,
	}
}

// GetItemJson - comment method.
func (uc *FormVersion) GetItemJson(ctx context.Context, primary entity.PrimaryKey, pretty bool) ([]byte, error) {
	if primary.FormID == uuid.Nil || primary.Version < 0 {
		return nil, mrcore.ErrUseCaseEntityNotFound.New()
	}

	if _, err := uc.formComponent.GetFormStatus(ctx, primary.FormID); err != nil {
		return nil, err
	}

	var (
		item entity.FormVersion
		err  error
	)

	// при указании формы нулевой версии возвращается копия текущей формы
	if primary.Version == 0 {
		item, err = uc.createFormVersionForTest(ctx, primary.FormID)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO: можно оптимизировать получая только body
		item, err = uc.storage.FetchOne(ctx, primary)
		if err != nil {
			return nil, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormVersion, primary)
		}
	}

	if pretty {
		var prettyJSON bytes.Buffer

		if err = json.Indent(&prettyJSON, item.Body, "", module.JsonPrettyIndent); err != nil {
			return nil, uc.errorWrapper.WrapErrorEntityFailed(err, entity.ModelNameFormVersion, primary)
		}

		return prettyJSON.Bytes(), nil
	}

	return item.Body, nil
}

// PrepareForTest - comment method.
func (uc *FormVersion) PrepareForTest(ctx context.Context, formID uuid.UUID) error {
	lastVersionItem, err := uc.getItemLastVersion(ctx, formID)
	if err != nil {
		return err
	}

	if lastVersionItem.ActivityStatus == enum.ActivityStatusArchived {
		return mrcore.ErrUseCaseSwitchStatusRejected.New(
			enum.ActivityStatusArchived,
			enum.ActivityStatusTesting,
		)
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(formID)); err != nil {
		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameFormVersion)
	} else {
		defer unlock()
	}

	item, err := uc.createFormVersionForTest(ctx, formID)
	if err != nil {
		return err
	}

	eventName := "Create"
	item.Version = lastVersionItem.Version

	if lastVersionItem.ActivityStatus == enum.ActivityStatusTesting {
		if err = uc.storage.Update(ctx, item); err != nil {
			return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormVersion, item.ID)
		}

		eventName = "Update"
	} else {
		if lastVersionItem.ActivityStatus == enum.ActivityStatusPublished {
			item.Version++
		}

		if err = uc.storage.Insert(ctx, item); err != nil {
			return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameFormVersion)
		}
	}

	uc.eventEmitter.Emit(ctx, eventName, mrmsg.Data{"formId": item.ID, "version": item.Version})

	return nil
}

// Publish - comment method.
func (uc *FormVersion) Publish(ctx context.Context, formID uuid.UUID) error {
	item, err := uc.getItemLastVersion(ctx, formID)
	if err != nil {
		return err
	}

	if item.ActivityStatus == enum.ActivityStatusPublished {
		return nil // переключения не требуется
	}

	if item.ActivityStatus != enum.ActivityStatusTesting {
		return mrcore.ErrUseCaseSwitchStatusRejected.New(
			item.ActivityStatus,
			enum.ActivityStatusPublished,
		)
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(formID)); err != nil {
		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameFormVersion)
	} else {
		defer unlock()
	}

	if err = uc.storage.UpdateStatus(ctx, item, enum.ActivityStatusPublished); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameFormVersion)
	}

	uc.eventEmitter.Emit(ctx, "Publish", mrmsg.Data{"formId": formID, "version": item.Version})

	return nil
}

func (uc *FormVersion) getLockKey(formID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", entity.ModelNameFormVersion, formID)
}

func (uc *FormVersion) getItemLastVersion(ctx context.Context, formID uuid.UUID) (entity.FormVersionStatus, error) {
	if formID == uuid.Nil {
		return entity.FormVersionStatus{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	if formStatus, err := uc.formComponent.GetFormStatus(ctx, formID); err != nil {
		return entity.FormVersionStatus{}, err
	} else if formStatus != mrenum.ItemStatusEnabled {
		return entity.FormVersionStatus{}, mrcore.ErrUseCaseEntityNotAvailable.New()
	}

	item, err := uc.storage.FetchOneLastVersion(ctx, formID)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return entity.FormVersionStatus{
				FormID:         formID,
				Version:        1,
				ActivityStatus: enum.ActivityStatusDraft,
			}, nil
		}

		return entity.FormVersionStatus{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormVersion, formID)
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
		ActivityStatus: enum.ActivityStatusTesting,
	}, nil
}
