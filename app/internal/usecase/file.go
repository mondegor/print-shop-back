package usecase

import (
    "context"
    "strings"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

type (
    File struct {
        storage mrstorage.FileProvider
    }
)

func NewFile(storage mrstorage.FileProvider) *File {
    return &File{
        storage: storage,
    }
}

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *File) Get(ctx context.Context, path string) (*mrstorage.File, error) {
    path = strings.TrimLeft(path, "/")

    if path == "" {
        return nil, mrcore.FactoryErrServiceEmptyInputData.New("path")
    }

    item := mrstorage.File{
        Name: path,
    }

    err := uc.storage.Download(ctx, &item)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, mrstorage.ModelNameFile)
    }

    return &item, nil
}
