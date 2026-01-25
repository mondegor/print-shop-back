package usecase

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrlock"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/util/conv"

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
		eventEmitter mrevent.Emitter
		errorWrapper errors.Wrapper
		logger       mrlog.Logger
	}
)

// NewCompanyPageLogo - создаёт объект CompanyPageLogo.
func NewCompanyPageLogo(
	storage prov.CompanyPageLogoStorage,
	fileAPI mrstorage.FileProviderAPI,
	locker mrlock.Locker,
	eventEmitter mrevent.Emitter,
	logger mrlog.Logger,
) *CompanyPageLogo {
	return &CompanyPageLogo{
		storage:      storage,
		fileAPI:      fileAPI,
		locker:       locker,
		eventEmitter: mrevent.EmitterWithSource(eventEmitter, entity.ModelNameCompanyPageLogo),
		errorWrapper: errors.NewUseCaseWrapper(),
		logger:       logger,
	}
}

// StoreFile - comment method.
func (uc *CompanyPageLogo) StoreFile(ctx context.Context, accountID uuid.UUID, image mrtype.Image) error {
	if accountID == uuid.Nil {
		return errors.ErrUseCaseEntityNotFound
	}

	if image.OriginalName == "" || image.Size == 0 {
		return errors.ErrUseCaseInvalidFile
	}

	newLogoPath, err := uc.getLogoPath(accountID, image.OriginalName)
	if err != nil {
		return err
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(accountID)); err != nil {
		return uc.errorWrapper.Wrap(err)
	} else {
		defer unlock()
	}

	oldLogoMeta, err := uc.storage.FetchMeta(ctx, accountID)
	if err != nil {
		return uc.errorWrapper.Wrap(err, "accountId", accountID)
	}

	image.Path = newLogoPath

	if err = uc.fileAPI.Upload(ctx, image.ToFile()); err != nil {
		return uc.errorWrapper.Wrap(err, "imagePath", image.Path)
	}

	logoMeta := mrentity.ImageMeta{
		Path:   newLogoPath,
		Width:  image.Width,
		Height: image.Height,
		Size:   image.Size,
	}

	if err = uc.storage.UpdateMeta(ctx, accountID, logoMeta); err != nil {
		uc.removeLogoFile(ctx, newLogoPath, oldLogoMeta.Path)

		return uc.errorWrapper.Wrap(err, "accountId", accountID)
	}

	uc.eventEmitter.Emit(ctx, "StoreFile", conv.Group{"accountId": accountID, "path": newLogoPath, "old-path": oldLogoMeta.Path})
	uc.removeLogoFile(ctx, oldLogoMeta.Path, newLogoPath)

	return nil
}

// RemoveFile - comment method.
func (uc *CompanyPageLogo) RemoveFile(ctx context.Context, accountID uuid.UUID) error {
	if accountID == uuid.Nil {
		return errors.ErrUseCaseEntityNotFound
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(accountID)); err != nil {
		return uc.errorWrapper.Wrap(err)
	} else {
		defer unlock()
	}

	logoMeta, err := uc.storage.FetchMeta(ctx, accountID)
	if err != nil {
		return uc.errorWrapper.Wrap(err, "accountId", accountID)
	}

	if err = uc.storage.DeleteMeta(ctx, accountID); err != nil {
		return uc.errorWrapper.Wrap(err, "accountId", accountID)
	}

	uc.eventEmitter.Emit(ctx, "RemoveFile", conv.Group{"accountId": accountID, "meta": logoMeta})
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
		uc.logger.Error(ctx, "fileAPI.Remove", "error", err)
	}
}
