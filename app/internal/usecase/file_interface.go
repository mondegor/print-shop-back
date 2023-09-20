package usecase

import (
    "context"

    "github.com/mondegor/go-storage/mrstorage"
)

type (
    FileService interface {
        Get(ctx context.Context, path string) (*mrstorage.File, error)
    }
)
