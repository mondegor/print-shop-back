package pub

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrpath"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request/parser"

	"github.com/mondegor/print-shop-back/internal/filestation/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/filestation/section/pub/usecase"
)

func initImageProxyController(
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	requestParser *parser.String,
	responseFileSender mrserver.FileResponseSender,
	fileAPI mrstorage.FileProviderAPI,
	basePath mrpath.Builder,
) (mrserver.HttpController, error) {
	useCase := usecase.NewFileProviderAdapter(fileAPI, useCaseErrorWrapper)

	controller := httpv1.NewImageProxy(
		requestParser,
		responseFileSender,
		useCase,
		basePath,
	)

	return controller, nil
}
