package http_v1

import (
	"fmt"
	"net/http"
	module "print-shop-back/internal/modules/file-station"
	usecase "print-shop-back/internal/modules/file-station/usecase/public-api"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ImageProxy struct {
		section   mrcore.ClientSection
		service   usecase.FileProviderAdapterService
		imagesURL string
	}
)

func NewImageProxy(
	section mrcore.ClientSection,
	service usecase.FileProviderAdapterService,
	basePath string, // :TODO: to URL
) *ImageProxy {
	return &ImageProxy{
		section:   section,
		service:   service,
		imagesURL: fmt.Sprintf("/%s/*path", strings.Trim(basePath, "/")),
	}
}

func (ht *ImageProxy) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitImageProxyPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(ht.imagesURL), moduleAccessFunc(ht.Get()))
}

func (ht *ImageProxy) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.Get(c.Context(), c.ParamFromPath("path"))

		if err != nil {
			return err
		}

		defer item.Body.Close()

		return c.SendFile(item.FileInfo, "", item.Body)
	}
}
