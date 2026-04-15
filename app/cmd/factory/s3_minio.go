package factory

import (
	"context"

	"github.com/mondegor/go-storage/mrminio"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/mime"

	"github.com/mondegor/print-shop-back/config"
)

// NewS3Minio - создаёт объект mrminio.ConnAdapter.
func NewS3Minio(ctx context.Context, logger mrlog.Logger, tracer mrtrace.Tracer, cfg config.Config) (*mrminio.ConnAdapter, error) {
	mrlog.Info(logger, "Create and init file provider pool")

	opts := mrminio.Options{
		Host:     cfg.S3.Host,
		Port:     cfg.S3.Port,
		UseSSL:   cfg.S3.UseSSL,
		User:     cfg.S3.Username,
		Password: cfg.S3.Password,
	}

	conn := mrminio.New(
		cfg.S3.CreateBuckets,
		mime.NewTypeList(cfg.MimeTypes), // TODO: можно вынести в общую переменную
		tracer,
	)

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, nil
}

// RegisterS3ImageStorage - comment func.
func RegisterS3ImageStorage(
	logger mrlog.Logger,
	cfg config.Config,
	pool *mrstorage.FileProviderPool,
	conn *mrminio.ConnAdapter,
) error {
	storage, err := newS3MinioFileProvider(
		logger,
		conn,
		cfg.FileProviders.ImageStorage.BucketName,
	)
	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorage.Name, storage)
}

func newS3MinioFileProvider(
	logger mrlog.Logger,
	conn *mrminio.ConnAdapter,
	bucketName string,
) (*mrminio.FileProvider, error) {
	mrlog.Info(logger, "Create and init file provider with bucket '"+bucketName+"'")

	created, err := conn.InitBucket(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if created {
		mrlog.Debug(logger, "Bucket '%s' created", bucketName)
	} else {
		mrlog.Debug(logger, "Bucket '%s' exists, OK", bucketName)
	}

	return mrminio.NewFileProvider(conn, bucketName), nil
}
