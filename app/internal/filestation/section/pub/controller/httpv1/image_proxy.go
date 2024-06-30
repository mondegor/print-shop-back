package httpv1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mondegor/print-shop-back/internal/filestation/section/pub/usecase"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// ImageProxy - comment struct.
	ImageProxy struct {
		parser    mrserver.RequestParserString
		sender    mrserver.FileResponseSender
		useCase   usecase.FileProviderAdapterUseCase
		imagesURL string
	}
)

// NewImageProxy - создаёт объект ImageProxy.
func NewImageProxy(
	parser mrserver.RequestParserString,
	sender mrserver.FileResponseSender,
	useCase usecase.FileProviderAdapterUseCase,
	basePath string, // TODO: to URL
) *ImageProxy {
	return &ImageProxy{
		parser:    parser,
		sender:    sender,
		useCase:   useCase,
		imagesURL: fmt.Sprintf("/%s/"+mrserver.VarRestOfURL, strings.Trim(basePath, "/")),
	}
}

// Handlers - comment method.
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
