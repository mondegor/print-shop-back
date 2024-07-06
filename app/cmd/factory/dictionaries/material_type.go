package dictionaries

import (
	"context"

	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/materialtype/api/availability"
)

// NewMaterialTypeModuleOptions - создаёт объект materialtype.Options.
func NewMaterialTypeModuleOptions(_ context.Context, opts app.Options) (materialtype.Options, error) {
	materialTypeDictionary, err := opts.Translator.Dictionary("dictionaries/material-types")
	if err != nil {
		return materialtype.Options{}, err
	}

	return materialtype.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: materialtype.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		UnitMaterialType: materialtype.UnitMaterialTypeOptions{
			Dictionary: materialTypeDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewMaterialTypeAvailabilityAPI - создаёт объект usecase.MaterialType.
func NewMaterialTypeAvailabilityAPI(ctx context.Context, opts app.Options) (*usecase.MaterialType, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries laminate type availability API")

	return availability.NewMaterialType(opts.PostgresConnManager, opts.UsecaseErrorWrapper), nil
}

// RegisterMaterialTypeErrors - comment func.
func RegisterMaterialTypeErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.MaterialTypeErrors()))
}
