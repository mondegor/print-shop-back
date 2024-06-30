package usecase

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// FileProviderAdapter - comment struct.
	FileProviderAdapter struct {
		fileAPI      mrstorage.FileProviderAPI
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewFileProviderAdapter - создаёт объект FileProviderAdapter.
func NewFileProviderAdapter(fileAPI mrstorage.FileProviderAPI, errorWrapper mrcore.UsecaseErrorWrapper) *FileProviderAdapter {
	return &FileProviderAdapter{
		fileAPI:      fileAPI,
		errorWrapper: errorWrapper,
	}
}

// Get - comment method.
// WARNING you don't forget to call item.File.Body.Close().
func (uc *FileProviderAdapter) Get(ctx context.Context, filePath string) (mrtype.File, error) {
	filePath = strings.TrimLeft(filePath, "/")

	if filePath == "" {
		return mrtype.File{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	file, err := uc.fileAPI.Download(ctx, filePath)
	if err != nil {
		return mrtype.File{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", filePath)
	}

	return file, nil
}
