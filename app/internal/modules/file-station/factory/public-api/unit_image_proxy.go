package factory

import (
	http_v1 "print-shop-back/internal/modules/file-station/controller/http_v1/public-api"
	"print-shop-back/internal/modules/file-station/factory"
	usecase "print-shop-back/internal/modules/file-station/usecase/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitImageProxy(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	service := usecase.NewFileProviderAdapter(opts.UnitImageProxy.FileAPI, opts.ServiceHelper)
	*c = append(*c, http_v1.NewImageProxy(section, service, opts.UnitImageProxy.BaseURL))

	return nil
}
