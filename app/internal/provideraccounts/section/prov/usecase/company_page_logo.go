package usecase

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"
)

type (
	// CompanyPageLogo - comment struct.
	CompanyPageLogo struct {
		storage      prov.CompanyPageLogoStorage
		fileAPI      mrstorage.FileProviderAPI
		locker       mrlock.Locker
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewCompanyPageLogo - создаёт объект CompanyPageLogo.
func NewCompanyPageLogo(
	storage prov.CompanyPageLogoStorage,
	fileAPI mrstorage.FileProviderAPI,
	locker mrlock.Locker,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UseCaseErrorWrapper,
) *CompanyPageLogo {
	return &CompanyPageLogo{
		storage:      storage,
		fileAPI:      fileAPI,
		locker:       locker,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNameCompanyPageLogo),
		errorWrapper: errorWrapper,
	}
}

// StoreFile - comment method.
func (uc *CompanyPageLogo) StoreFile(ctx context.Context, accountID uuid.UUID, image mrtype.Image) error {
	if accountID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if image.OriginalName == "" || image.Size == 0 {
		return mrcore.ErrUseCaseInvalidFile.New()
	}

	newLogoPath, err := uc.getLogoPath(accountID, image.OriginalName)
	if err != nil {
		return err
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(accountID)); err != nil {
		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCompanyPageLogo)
	} else {
		defer unlock()
	}

	oldLogoMeta, err := uc.storage.FetchMeta(ctx, accountID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	image.Path = newLogoPath

	if err = uc.fileAPI.Upload(ctx, image.ToFile()); err != nil {
		return uc.errorWrapper.WrapErrorEntityFailed(err, "FileProviderAPI", image.Path)
	}

	logoMeta := mrentity.ImageMeta{
		Path:   newLogoPath,
		Width:  image.Width,
		Height: image.Height,
		Size:   image.Size,
	}

	if err = uc.storage.UpdateMeta(ctx, accountID, logoMeta); err != nil {
		uc.removeLogoFile(ctx, newLogoPath, oldLogoMeta.Path)

		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	uc.eventEmitter.Emit(ctx, "StoreFile", mrmsg.Data{"accountId": accountID, "path": newLogoPath, "old-path": oldLogoMeta.Path})
	uc.removeLogoFile(ctx, oldLogoMeta.Path, newLogoPath)

	return nil
}

// RemoveFile - comment method.
func (uc *CompanyPageLogo) RemoveFile(ctx context.Context, accountID uuid.UUID) error {
	if accountID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(accountID)); err != nil {
		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCompanyPageLogo)
	} else {
		defer unlock()
	}

	logoMeta, err := uc.storage.FetchMeta(ctx, accountID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	if err = uc.storage.DeleteMeta(ctx, accountID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	uc.eventEmitter.Emit(ctx, "RemoveFile", mrmsg.Data{"accountId": accountID, "meta": logoMeta})
	uc.removeLogoFile(ctx, logoMeta.Path, "")

	return nil
}

func (uc *CompanyPageLogo) getLockKey(accountID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", entity.ModelNameCompanyPageLogo, accountID)
}

func (uc *CompanyPageLogo) getLogoPath(accountID uuid.UUID, filePath string) (string, error) {
	if ext := path.Ext(filePath); ext != "" {
		return fmt.Sprintf(
			"%s/%s-%x%s",
			module.UnitCompanyPageLogoDir,
			accountID,
			time.Now().UTC().UnixNano()&0xffff,
			ext,
		), nil
	}

	return "", fmt.Errorf("file %s: ext is empty", filePath)
}

func (uc *CompanyPageLogo) removeLogoFile(ctx context.Context, filePath, prevFilePath string) {
	if filePath == "" || filePath == prevFilePath {
		return
	}

	if err := uc.fileAPI.Remove(ctx, filePath); err != nil {
		mrlog.Ctx(ctx).Error().Err(err).Msg("fileAPI.Remove()")
	}
}
