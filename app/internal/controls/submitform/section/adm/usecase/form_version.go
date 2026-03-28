package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrlock"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum/activitystatus"
)

type (
	// FormVersion - comment struct.
	// FormVersion - comment struct.
	FormVersion struct {
		storage                     adm.FormVersionStorage
		formComponent               adm.SubmitFormComponent
		formCompiler                adm.FormCompilerComponent
		locker                      mrlock.Locker
		eventEmitter                mrevent.Emitter
		errorWrapper                errors.Wrapper
		errorNotFoundWrapper        errors.Wrapper
		errorVersionConflictWrapper errors.Wrapper
	}
)

// NewFormVersion - создаёт объект FormVersion.
func NewFormVersion(
	storage adm.FormVersionStorage,
	formComponent adm.SubmitFormComponent,
	formCompiler adm.FormCompilerComponent,
	locker mrlock.Locker,
	eventEmitter mrevent.Emitter,
) *FormVersion {
	return &FormVersion{
		storage:                     storage,
		formComponent:               formComponent,
		formCompiler:                formCompiler,
		locker:                      locker,
		eventEmitter:                mrevent.EmitterWithSource(eventEmitter, entity.ModelNameFormVersion),
		errorWrapper:                errors.NewServiceOperationFailedWrapper(),
		errorNotFoundWrapper:        errors.NewServiceRecordNotFoundWrapper(),
		errorVersionConflictWrapper: errors.NewServiceRecordVersionConflictWrapper(),
	}
}

// GetItemJson - comment method.
func (uc *FormVersion) GetItemJson(ctx context.Context, primary entity.PrimaryKey, pretty bool) ([]byte, error) {
	if primary.FormID == uuid.Nil || primary.Version < 0 {
		return nil, errors.ErrRecordNotFound
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
			return nil, uc.errorNotFoundWrapper.Wrap(err, "primary", primary)
		}
	}

	if pretty {
		var prettyJSON bytes.Buffer

		if err = json.Indent(&prettyJSON, item.Body, "", module.JsonPrettyIndent); err != nil {
			return nil, uc.errorWrapper.Wrap(err, "primary", primary)
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

	if lastVersionItem.ActivityStatus == activitystatus.Archived {
		return errors.ErrSwitchStatusRejected.New(
			activitystatus.Archived,
			activitystatus.Testing,
		)
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(formID)); err != nil {
		return uc.errorWrapper.Wrap(err)
	} else {
		defer unlock()
	}

	item, err := uc.createFormVersionForTest(ctx, formID)
	if err != nil {
		return err
	}

	eventName := "Create"
	item.Version = lastVersionItem.Version

	if lastVersionItem.ActivityStatus == activitystatus.Testing {
		if err = uc.storage.Update(ctx, item); err != nil {
			return uc.errorWrapper.Wrap(err, "itemId", item.ID)
		}

		eventName = "Update"
	} else {
		if lastVersionItem.ActivityStatus == activitystatus.Published {
			item.Version++
		}

		if err = uc.storage.Insert(ctx, item); err != nil {
			return uc.errorWrapper.Wrap(err)
		}
	}

	uc.eventEmitter.Emit(ctx, eventName, conv.Group{"formId": item.ID, "version": item.Version})

	return nil
}

// Publish - comment method.
func (uc *FormVersion) Publish(ctx context.Context, formID uuid.UUID) error {
	item, err := uc.getItemLastVersion(ctx, formID)
	if err != nil {
		return err
	}

	if item.ActivityStatus == activitystatus.Published {
		return nil // переключения не требуется
	}

	if item.ActivityStatus != activitystatus.Testing {
		return errors.ErrSwitchStatusRejected.New(
			item.ActivityStatus,
			activitystatus.Published,
		)
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(formID)); err != nil {
		return uc.errorWrapper.Wrap(err)
	} else {
		defer unlock()
	}

	if err = uc.storage.UpdateStatus(ctx, item, activitystatus.Published); err != nil {
		return uc.errorVersionConflictWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Publish", conv.Group{"formId": formID, "version": item.Version})

	return nil
}

func (uc *FormVersion) getLockKey(formID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", entity.ModelNameFormVersion, formID)
}

func (uc *FormVersion) getItemLastVersion(ctx context.Context, formID uuid.UUID) (entity.FormVersionStatus, error) {
	if formID == uuid.Nil {
		return entity.FormVersionStatus{}, errors.ErrRecordNotFound
	}

	formStatus, err := uc.formComponent.GetFormStatus(ctx, formID)
	if err != nil {
		return entity.FormVersionStatus{}, uc.errorNotFoundWrapper.Wrap(err)
	}

	if formStatus != itemstatus.Enabled {
		return entity.FormVersionStatus{}, errors.ErrRecordNotFound // TODO: ErrUseCaseEntityNotAvailable
	}

	item, err := uc.storage.FetchOneLastVersion(ctx, formID)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return entity.FormVersionStatus{
				FormID:         formID,
				Version:        1,
				ActivityStatus: activitystatus.Draft,
			}, nil
		}

		return entity.FormVersionStatus{}, uc.errorWrapper.Wrap(err, "itemId", formID)
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
		ActivityStatus: activitystatus.Testing,
	}, nil
}
