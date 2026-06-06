package factory

import (
	"fmt"
	"os"

	"github.com/mondegor/go-storage/mrfilestorage"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/util/mime"

	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/adapter/trace"
)

// NewFileSystem - создаёт объект mrfilestorage.FileSystem.
func NewFileSystem(logger log.Logger, cfg config.Config) *mrfilestorage.FileSystem {
	log.Info(logger, "Create and init file system")

	return mrfilestorage.New(
		os.FileMode(cfg.FSDirMode),
		cfg.FSCreateDirs,
		mime.NewTypeList(cfg.AllowedMimeTypes), // TODO: можно вынести в общую переменную
	)
}

// RegisterFileImageStorage - comment func.
func RegisterFileImageStorage(
	logger log.Logger,
	tracer trace.Tracer,
	cfg config.Config,
	pool *mrstorage.FileProviderPool,
	fs *mrfilestorage.FileSystem,
) error {
	storage, err := newFileStorageProvider(
		logger,
		tracer,
		fs,
		cfg.FileProviders.ImageStorage2RootDir,
	)
	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorageName, storage)
}

func newFileStorageProvider(
	logger log.Logger,
	tracer trace.Tracer,
	fs *mrfilestorage.FileSystem,
	rootDir string,
) (*mrfilestorage.FileProvider, error) {
	log.Info(logger, fmt.Sprintf("Create and init file provider with root dir '%s'", rootDir))

	created, err := fs.InitRootDir(rootDir)
	if err != nil {
		return nil, err
	}

	if created {
		log.Debug(logger, fmt.Sprintf("Root dir '%s' created", rootDir))
	} else {
		log.Debug(logger, fmt.Sprintf("Root dir '%s' exists, OK", rootDir))
	}

	return mrfilestorage.NewFileProvider(fs, tracer, rootDir), nil
}
