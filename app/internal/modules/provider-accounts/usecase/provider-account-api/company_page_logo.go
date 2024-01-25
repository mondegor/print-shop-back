package usecase

import (
	"context"
	"fmt"
	"path"
	module "print-shop-back/internal/modules/provider-accounts"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	"time"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CompanyPageLogo struct {
		storage       CompanyPageLogoStorage
		fileAPI       mrstorage.FileProviderAPI
		locker        mrcore.Locker
		eventBox      mrcore.EventBox
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewCompanyPageLogo(
	storage CompanyPageLogoStorage,
	fileAPI mrstorage.FileProviderAPI,
	locker mrcore.Locker,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *CompanyPageLogo {
	return &CompanyPageLogo{
		storage:       storage,
		fileAPI:       fileAPI,
		locker:        locker,
		eventBox:      eventBox,
		serviceHelper: serviceHelper,
	}
}

func (uc *CompanyPageLogo) StoreFile(ctx context.Context, accountID mrtype.KeyString, image mrtype.Image) error {
	if accountID == "" {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if image.OriginalName == "" || image.Size == 0 {
		return mrcore.FactoryErrServiceInvalidFile.New()
	}

	newLogoPath, err := uc.getLogoPath(accountID, image.OriginalName)

	if err != nil {
		return err
	}

	unlock, err := uc.locker.Lock(ctx, uc.getLockKey(accountID))

	if err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCompanyPageLogo)
	}

	defer unlock()

	oldLogoMeta, err := uc.storage.FetchMeta(ctx, accountID)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	image.Path = newLogoPath

	if err = uc.fileAPI.Upload(ctx, image.ToFile()); err != nil {
		return uc.serviceHelper.WrapErrorEntityFailed(err, "FileProviderAPI", image.Path)
	}

	logoMeta := mrentity.ImageMeta{
		Path:   newLogoPath,
		Width:  image.Width,
		Height: image.Height,
		Size:   image.Size,
	}

	if err = uc.storage.UpdateMeta(ctx, accountID, logoMeta); err != nil {
		uc.removeLogoFile(ctx, newLogoPath, oldLogoMeta.Path)
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	uc.eventBoxEmitEntity(ctx, "StoreFile", mrmsg.Data{"accountId": accountID, "path": newLogoPath, "old-path": oldLogoMeta.Path})
	uc.removeLogoFile(ctx, oldLogoMeta.Path, newLogoPath)

	return nil
}

func (uc *CompanyPageLogo) RemoveFile(ctx context.Context, accountID mrtype.KeyString) error {
	if accountID == "" {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	unlock, err := uc.locker.Lock(ctx, uc.getLockKey(accountID))

	if err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCompanyPageLogo)
	}

	defer unlock()

	logoMeta, err := uc.storage.FetchMeta(ctx, accountID)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	if err = uc.storage.DeleteMeta(ctx, accountID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPageLogo, accountID)
	}

	uc.eventBoxEmitEntity(ctx, "RemoveFile", mrmsg.Data{"accountId": accountID, "meta": logoMeta})
	uc.removeLogoFile(ctx, logoMeta.Path, "")

	return nil
}

func (uc *CompanyPageLogo) getLockKey(accountID mrtype.KeyString) string {
	return fmt.Sprintf("%s:%s", entity.ModelNameCompanyPageLogo, accountID)
}

func (uc *CompanyPageLogo) getLogoPath(accountID mrtype.KeyString, filePath string) (string, error) {
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

func (uc *CompanyPageLogo) removeLogoFile(ctx context.Context, filePath string, prevFilePath string) {
	if filePath == "" || filePath == prevFilePath {
		return
	}

	if err := uc.fileAPI.Remove(ctx, filePath); err != nil {
		mrctx.Logger(ctx).Err(err)
	}
}

func (uc *CompanyPageLogo) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNameCompanyPageLogo,
		eventName,
		data,
	)
}
