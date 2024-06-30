package controls

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
	"github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/controls/validate"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
)

// NewElementTemplateModuleOptions - создаёт объект elementtemplate.Options.
func NewElementTemplateModuleOptions(_ context.Context, opts app.Options) (elementtemplate.Options, error) {
	return elementtemplate.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: elementtemplate.RequestParsers{
			// Parser:       opts.RequestParsers.Parser,
			// ExtendParser: opts.RequestParsers.ExtendParser,
			ModuleParser: validate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.FileJson,
				pkgvalidate.NewDetailingParser(),
			),
		},
		ResponseSender: opts.ResponseSenders.FileSender,

		UnitElementTemplate: elementtemplate.UnitElementTemplateOptions{},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// RegisterElementTemplateErrors - comment func.
func RegisterElementTemplateErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.ElementTemplateErrors()))
}
