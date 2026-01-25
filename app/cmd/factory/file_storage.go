package factory

import (
	"fmt"
	"os"

	"github.com/mondegor/go-storage/mrfilestorage"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/mime"

	"github.com/mondegor/print-shop-back/config"
)

// NewFileSystem - создаёт объект mrfilestorage.FileSystem.
func NewFileSystem(logger mrlog.Logger, cfg config.Config) *mrfilestorage.FileSystem {
	mrlog.Info(logger, "Create and init file system")

	return mrfilestorage.New(
		os.FileMode(cfg.FileSystem.DirMode),
		cfg.FileSystem.CreateDirs,
		mime.NewTypeList(cfg.MimeTypes), // TODO: можно вынести в общую переменную
	)
}

// RegisterFileImageStorage - comment func.
func RegisterFileImageStorage(
	logger mrlog.Logger,
	tracer mrtrace.Tracer,
	cfg config.Config,
	pool *mrstorage.FileProviderPool,
	fs *mrfilestorage.FileSystem,
) error {
	storage, err := newFileStorageProvider(
		logger,
		tracer,
		fs,
		cfg.FileProviders.ImageStorage.RootDir,
	)
	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorage.Name, storage)
}

func newFileStorageProvider(
	logger mrlog.Logger,
	tracer mrtrace.Tracer,
	fs *mrfilestorage.FileSystem,
	rootDir string,
) (*mrfilestorage.FileProvider, error) {
	mrlog.Info(logger, fmt.Sprintf("Create and init file provider with root dir '%s'", rootDir))

	created, err := fs.InitRootDir(rootDir)
	if err != nil {
		return nil, err
	}

	if created {
		mrlog.Debug(logger, fmt.Sprintf("Root dir '%s' created", rootDir))
	} else {
		mrlog.Debug(logger, fmt.Sprintf("Root dir '%s' exists, OK", rootDir))
	}

	return mrfilestorage.NewFileProvider(fs, tracer, rootDir), nil
}
