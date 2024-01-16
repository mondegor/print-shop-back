package usecase

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAdapter struct {
		fileAPI       mrstorage.FileProviderAPI
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewFileProviderAdapter(
	fileAPI mrstorage.FileProviderAPI,
	serviceHelper *mrtool.ServiceHelper,
) *FileProviderAdapter {
	return &FileProviderAdapter{
		fileAPI:       fileAPI,
		serviceHelper: serviceHelper,
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
		return mrtype.File{}, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", filePath)
	}

	return file, nil
}
