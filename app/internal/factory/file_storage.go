package factory

import (
	"os"
	"print-shop-back/config"

	"github.com/mondegor/go-storage/mrfilestorage"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewFileSystem(cfg *config.Config, logger mrcore.Logger) *mrfilestorage.FileSystem {
	logger.Info("Create and init file system")

	opt := mrfilestorage.Options{
		DirMode:    os.FileMode(cfg.FileSystem.DirMode),
		CreateDirs: cfg.FileSystem.CreateDirs,
	}

	return mrfilestorage.New(opt)
}

func RegisterFileImageStorage(
	cfg *config.Config,
	pool *mrstorage.FileProviderPool,
	fs *mrfilestorage.FileSystem,
	logger mrcore.Logger,
) error {
	storage, err := newFileProvider(
		fs,
		cfg.FileProviders.ImageStorage.RootDir,
		logger,
	)

	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorage.Name, storage)
}

func newFileProvider(
	fs *mrfilestorage.FileSystem,
	rootDir string,
	logger mrcore.Logger,
) (*mrfilestorage.FileProvider, error) {
	logger.Info("Create and init file provider with root dir '%s'", rootDir)

	created, err := fs.InitRootDir(rootDir)

	if err != nil {
		return nil, err
	}

	if created {
		mrcore.LogInfo("Root dir '%s' created", rootDir)
	} else {
		mrcore.LogInfo("Root dir '%s' exists, OK", rootDir)
	}

	return mrfilestorage.NewFileProvider(fs, rootDir), nil
}
