package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/filestation/section/pub"
)

type (
	// ImageProxy - comment struct.
	ImageProxy struct {
		parser    mrserver.RequestParserString
		sender    mrserver.FileResponseSender
		useCase   pub.FileProviderAdapterUseCase
		imagesURL string
	}
)

// NewImageProxy - создаёт контроллер ImageProxy.
func NewImageProxy(
	parser mrserver.RequestParserString,
	sender mrserver.FileResponseSender,
	useCase pub.FileProviderAdapterUseCase,
	basePath mrpath.PathBuilder,
) *ImageProxy {
	return &ImageProxy{
		parser:    parser,
		sender:    sender,
		useCase:   useCase,
		imagesURL: basePath.BuildPath(mrserver.VarRestOfURL),
	}
}

// Handlers - возвращает обработчики контроллера ImageProxy.
func (ht *ImageProxy) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: ht.imagesURL, Func: ht.Get},
	}
}

// Get - comment method.
func (ht *ImageProxy) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.Get(r.Context(), ht.parser.PathParamString(r, mrserver.VarRestOfURL))
	if err != nil {
		return err
	}

	defer item.Body.Close()

	return ht.sender.SendFile(r.Context(), w, item)
}
