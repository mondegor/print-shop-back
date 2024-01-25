package http_v1

import (
	"fmt"
	"net/http"
	usecase "print-shop-back/internal/modules/file-station/usecase/public-api"
	"strings"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	ImageProxy struct {
		parser    mrserver.RequestParserString
		sender    mrserver.FileResponseSender
		service   usecase.FileProviderAdapterService
		imagesURL string
	}
)

func NewImageProxy(
	parser mrserver.RequestParserString,
	sender mrserver.FileResponseSender,
	service usecase.FileProviderAdapterService,
	basePath string, // :TODO: to URL
) *ImageProxy {
	return &ImageProxy{
		parser:    parser,
		sender:    sender,
		service:   service,
		imagesURL: fmt.Sprintf("/%s/*path", strings.Trim(basePath, "/")),
	}
}

func (ht *ImageProxy) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, ht.imagesURL, "", ht.Get},
	}
}

func (ht *ImageProxy) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.Get(r.Context(), ht.parser.PathParamString(r, "path"))

	if err != nil {
		return err
	}

	defer item.Body.Close()

	return ht.sender.SendFile(w, item.FileInfo, "", item.Body)
}
