package dictionaries

import (
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/api/availability"
)

// NewMaterialTypeModuleOptions - создаёт объект materialtype.Options.
func NewMaterialTypeModuleOptions(opts app.Options) materialtype.Options {
	return materialtype.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: materialtype.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}
}

// NewMaterialTypeAvailabilityAPI - создаёт объект usecase.MaterialType.
func NewMaterialTypeAvailabilityAPI(opts app.Options) (*usecase.MaterialType, error) {
	mrlog.Info(opts.Logger, "Create and init dictionaries laminate type availability API")

	return availability.NewMaterialType(opts.PostgresConnManager, opts.UsecaseErrorWrapper, opts.Tracer), nil
}
