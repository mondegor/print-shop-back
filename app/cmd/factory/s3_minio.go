package factory

import (
	"context"

	"github.com/mondegor/go-core/mrstorage"
	"github.com/mondegor/go-core/util/mime"
	"github.com/mondegor/go-storage/mrminio"

	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/adapter/trace"
)

// InitS3Minio - создаёт объект mrminio.ConnAdapter.
func InitS3Minio(ctx context.Context, logger log.Logger, tracer trace.Tracer, cfg config.Config) (*mrminio.ConnAdapter, error) {
	log.Info(logger, "Create and init s3 minio connection")

	opts := mrminio.Options{
		Host:     cfg.S3Host,
		Port:     cfg.S3Port,
		UseSSL:   cfg.S3UseSSL,
		User:     cfg.S3Username,
		Password: cfg.S3Password,
	}

	conn := mrminio.New(
		cfg.S3CreateBuckets,
		mime.NewTypeList(cfg.AllowedMimeTypes), // TODO: можно вынести в общую переменную
		tracer,
	)

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, nil
}

// RegisterS3ImageStorage - регистрирует в пуле S3-хранилище изображений
// (создаёт бакет при необходимости).
func RegisterS3ImageStorage(
	logger log.Logger,
	cfg config.Config,
	pool *mrstorage.FileProviderPool,
	conn *mrminio.ConnAdapter,
) error {
	storage, err := newS3MinioFileProvider(
		logger,
		conn,
		cfg.FileProviders.ImageStorageBucketName,
	)
	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorageName, storage)
}

func newS3MinioFileProvider(
	logger log.Logger,
	conn *mrminio.ConnAdapter,
	bucketName string,
) (*mrminio.FileProvider, error) {
	log.Info(logger, "Create and init file provider with bucket '"+bucketName+"'")

	created, err := conn.InitBucket(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if created {
		log.Debug(logger, "Bucket '%s' created", bucketName)
	} else {
		log.Debug(logger, "Bucket '%s' exists, OK", bucketName)
	}

	return mrminio.NewFileProvider(conn, bucketName), nil
}
