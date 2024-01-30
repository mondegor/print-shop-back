package usecase

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAdapter struct {
		fileAPI       mrstorage.FileProviderAPI
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewFileProviderAdapter(
	fileAPI mrstorage.FileProviderAPI,
	usecaseHelper *mrcore.UsecaseHelper,
) *FileProviderAdapter {
	return &FileProviderAdapter{
		fileAPI:       fileAPI,
		usecaseHelper: usecaseHelper,
	}
}

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *FileProviderAdapter) Get(ctx context.Context, filePath string) (mrtype.File, error) {
	filePath = strings.TrimLeft(filePath, "/")

	if filePath == "" {
		return mrtype.File{}, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	file, err := uc.fileAPI.Download(ctx, filePath)

	if err != nil {
		return mrtype.File{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", filePath)
	}

	return file, nil
}
