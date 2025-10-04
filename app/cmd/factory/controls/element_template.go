package controls

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
	"github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/controls/validate"
)

// NewElementTemplateModuleOptions - создаёт объект elementtemplate.Options.
func NewElementTemplateModuleOptions(opts app.Options) elementtemplate.Options {
	return elementtemplate.Options{
		Logger:               opts.Logger,
		EventEmitter:         opts.EventEmitter,
		UsecaseErrorWrapper:  opts.UsecaseErrorWrapper,
		FileUserErrorWrapper: opts.FileUserErrorWrapper,
		DBConnManager:        opts.PostgresConnManager,
		RequestParsers: elementtemplate.RequestParsers{
			// Parser:       opts.RequestParsers.Parser,
			// ExtendParser: opts.RequestParsers.ExtendParser,
			ModuleParser: validate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.FileJson,
				pkgvalidate.NewDetailingParser(opts.Logger),
			),
		},
		ResponseSender: opts.ResponseSenders.FileSender,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}
