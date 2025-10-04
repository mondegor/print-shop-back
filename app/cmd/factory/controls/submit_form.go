package controls

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	"github.com/mondegor/print-shop-back/internal/factory/controls/elementtemplate/api/header"
	"github.com/mondegor/print-shop-back/internal/factory/controls/submitform"
	pkgvalidate "github.com/mondegor/print-shop-back/pkg/controls/validate"
)

// NewSubmitFormModuleOptions - создаёт объект submitform.Options.
func NewSubmitFormModuleOptions(opts app.Options) submitform.Options {
	return submitform.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		StorageErrorWrapper: opts.StorageErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		Locker:              opts.Locker,
		RequestParsers: submitform.RequestParsers{
			// Parser:       opts.RequestParsers.Parser,
			// ExtendParser: opts.RequestParsers.ExtendParser,
			ModuleParser: validate.NewParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.FileJson,
				pkgvalidate.NewDetailingParser(opts.Logger),
			),
		},
		ResponseSender: opts.ResponseSenders.FileSender,

		ElementTemplateAPI: header.NewElementTemplate(opts.PostgresConnManager, opts.UsecaseErrorWrapper, opts.Tracer),

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}
