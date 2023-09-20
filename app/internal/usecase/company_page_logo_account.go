package usecase

import (
    "context"
    "fmt"
    "path/filepath"
    "print-shop-back/internal/entity"
    "time"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrtool"
)

const (
    CompanyPageLogoDir = "accounts/companies-pages" // :TODO: в конфиг вытащить
)

type (
    AccountCompanyPageLogo struct {
        storage AccountCompanyPageLogoStorage
        storageFiles mrstorage.FileProvider
        locker mrcore.Locker
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
    }
)

func NewAccountCompanyPageLogo(storage AccountCompanyPageLogoStorage,
                               storageFiles mrstorage.FileProvider,
                               locker mrcore.Locker,
                               eventBox mrcore.EventBox,
                               serviceHelper *mrtool.ServiceHelper) *AccountCompanyPageLogo {
    return &AccountCompanyPageLogo{
        storage: storage,
        storageFiles: storageFiles,
        locker: locker,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
    }
}

func (uc *AccountCompanyPageLogo) Store(ctx context.Context, item *entity.CompanyPageLogoObject) error {
    if item.AccountId == "" {
        return mrcore.FactoryErrServiceEmptyInputData.New("item.AccountId")
    }

    newLogoPath, err := uc.getLogoPath(item)

    if err != nil {
        return err
    }

    unlock, err := uc.locker.Lock(ctx, uc.getLockKey(item.AccountId))

    if err != nil {
        return err
    }

    defer unlock()

    oldLogoPath, err := uc.storage.Fetch(ctx, item.AccountId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCompanyPageLogo)
    }

    file := item.File
    file.Name = newLogoPath

    err = uc.storageFiles.Upload(ctx, &file)

    if err != nil {
        return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, mrstorage.ModelNameFile)
    }

    err = uc.storage.Update(ctx, item.AccountId, newLogoPath)

    if err != nil {
        uc.removeFile(ctx, newLogoPath)
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCompanyPageLogo)
    }

    uc.eventBox.Emit(
        "%s::Upload: cid=%d; path=%s",
        entity.ModelNameCompanyPageLogo,
        item.AccountId,
        newLogoPath,
    )

    uc.removeFile(ctx, oldLogoPath)

    return nil
}

func (uc *AccountCompanyPageLogo) Remove(ctx context.Context, accountId mrentity.KeyString) error {
    if accountId == "" {
        return mrcore.FactoryErrServiceEmptyInputData.New("accountId")
    }

    unlock, err := uc.locker.Lock(ctx, uc.getLockKey(accountId))

    if err != nil {
        return err
    }

    defer unlock()

    oldLogoPath, err := uc.storage.Fetch(ctx, accountId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCompanyPageLogo)
    }

    err = uc.storage.Delete(ctx, accountId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCompanyPageLogo)
    }

    uc.eventBox.Emit(
        "%s::Remove: cid=%d; path=%s",
        entity.ModelNameCompanyPageLogo,
        accountId,
        oldLogoPath,
    )

    uc.removeFile(ctx, oldLogoPath)

    return nil
}

func (uc *AccountCompanyPageLogo) getLockKey(accountId mrentity.KeyString) string {
    return fmt.Sprintf("%s:%d", entity.ModelNameCompanyPageLogo, accountId)
}

func (uc *AccountCompanyPageLogo) getLogoPath(item *entity.CompanyPageLogoObject) (string, error) {
    ext := filepath.Ext(item.File.Name)

    if ext == "" {
        return "", fmt.Errorf("ext is empty")
    }

    return fmt.Sprintf(
        "%s/%s-%x%s",
        CompanyPageLogoDir,
        item.AccountId,
        time.Now().UnixNano() & 0xffff,
        ext,
    ), nil
}

func (uc *AccountCompanyPageLogo) removeFile(ctx context.Context, filePath string) {
    if filePath == "" {
        return
    }

    err := uc.storageFiles.Remove(ctx, filePath)

    if err != nil {
        mrctx.Logger(ctx).Err(err)
    }
}
