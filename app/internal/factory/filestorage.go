package factory

import (
    "errors"
    "os"
    "print-shop-back/config"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-storage/mrfilestorage"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

func NewFileStorage(cfg *config.Config, logger mrcore.Logger) (mrstorage.FileProvider, error) {
    logger.Info("Init file storage")

    err := createBaseDir(cfg.FileStorage.DownloadDir, 0755)

    if err != nil {
        return nil, err
    }

    err = os.MkdirAll(cfg.FileStorage.DownloadDir + "/" + usecase.CompanyPageLogoDir, 0755)

    if err != nil {
        return nil, err
    }

    return mrfilestorage.New(cfg.FileStorage.DownloadDir), nil
}

func createBaseDir(dirPath string, perms os.FileMode) error {
    _, err := os.Stat(dirPath)

    if errors.Is(err, os.ErrNotExist) {
        err = os.Mkdir(dirPath, perms)
    }

    return err
}
