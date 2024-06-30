package factory

import (
	"context"
	"os"

	"github.com/mondegor/print-shop-back/config"

	"github.com/mondegor/go-storage/mrfilestorage"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewFileSystem - создаёт объект mrfilestorage.FileSystem.
func NewFileSystem(ctx context.Context, cfg config.Config) *mrfilestorage.FileSystem {
	mrlog.Ctx(ctx).Info().Msg("Create and init file system")

	return mrfilestorage.New(
		os.FileMode(cfg.FileSystem.DirMode),
		cfg.FileSystem.CreateDirs,
		mrlib.NewMimeTypeList(cfg.MimeTypes), // TODO: можно вынести в общую переменную
	)
}

// RegisterFileImageStorage - comment func.
func RegisterFileImageStorage(ctx context.Context, cfg config.Config, pool *mrstorage.FileProviderPool, fs *mrfilestorage.FileSystem) error {
	storage, err := newFileStorageProvider(
		ctx,
		fs,
		cfg.FileProviders.ImageStorage.RootDir,
	)
	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorage.Name, storage)
}

func newFileStorageProvider(
	ctx context.Context,
	fs *mrfilestorage.FileSystem,
	rootDir string,
) (*mrfilestorage.FileProvider, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msgf("Create and init file provider with root dir '%s'", rootDir)

	created, err := fs.InitRootDir(rootDir)
	if err != nil {
		return nil, err
	}

	if created {
		logger.Debug().Msgf("Root dir '%s' created", rootDir)
	} else {
		logger.Debug().Msgf("Root dir '%s' exists, OK", rootDir)
	}

	return mrfilestorage.NewFileProvider(fs, rootDir), nil
}
