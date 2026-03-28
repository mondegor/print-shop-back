package usecase

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrmodel"
)

type (
	// FileProviderAdapter - comment struct.
	FileProviderAdapter struct {
		fileAPI      mrstorage.FileProviderAPI
		errorWrapper errors.Wrapper
	}
)

// NewFileProviderAdapter - создаёт объект FileProviderAdapter.
func NewFileProviderAdapter(
	fileAPI mrstorage.FileProviderAPI,
) *FileProviderAdapter {
	return &FileProviderAdapter{
		fileAPI:      fileAPI,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// Get - comment method.
// WARNING you don't forget to call item.File.Body.Close().
func (uc *FileProviderAdapter) Get(ctx context.Context, filePath string) (mrmodel.File, error) {
	filePath = strings.TrimLeft(filePath, "/")

	if filePath == "" {
		return mrmodel.File{}, errors.ErrRecordNotFound
	}

	file, err := uc.fileAPI.Download(ctx, filePath)
	if err != nil {
		return mrmodel.File{}, uc.errorWrapper.Wrap(err, "filePath", filePath)
	}

	return file, nil
}
