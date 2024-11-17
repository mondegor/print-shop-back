package controls

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
	"github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/controls/validate"
)

// NewElementTemplateModuleOptions - создаёт объект elementtemplate.Options.
func NewElementTemplateModuleOptions(_ context.Context, opts app.Options) (elementtemplate.Options, error) {
	return elementtemplate.Options{
		EventEmitter:        opts.EventEmitter,
		UseCaseErrorWrapper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager:       opts.PostgresConnManager,
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
