package usecase

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrtype"
)

type (
	// FileProviderAdapter - comment struct.
	FileProviderAdapter struct {
		fileAPI      mrstorage.FileProviderAPI
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewFileProviderAdapter - создаёт объект FileProviderAdapter.
func NewFileProviderAdapter(
	fileAPI mrstorage.FileProviderAPI,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *FileProviderAdapter {
	return &FileProviderAdapter{
		fileAPI:      fileAPI,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, "FileProviderAPI"),
	}
}

// Get - comment method.
// WARNING you don't forget to call item.File.Body.Close().
func (uc *FileProviderAdapter) Get(ctx context.Context, filePath string) (mrtype.File, error) {
	filePath = strings.TrimLeft(filePath, "/")

	if filePath == "" {
		return mrtype.File{}, mr.ErrUseCaseEntityNotFound.New()
	}

	file, err := uc.fileAPI.Download(ctx, filePath)
	if err != nil {
		return mrtype.File{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "filePath", filePath)
	}

	return file, nil
}
