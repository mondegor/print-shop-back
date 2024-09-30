package dictionaries

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/paperfacture/api/availability"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

// NewPaperFactureModuleOptions - создаёт объект paperfacture.Options.
func NewPaperFactureModuleOptions(_ context.Context, opts app.Options) (paperfacture.Options, error) {
	paperFactureDictionary, err := opts.Translator.Dictionary("dictionaries/paper-factures")
	if err != nil {
		return paperfacture.Options{}, err
	}

	return paperfacture.Options{
		EventEmitter:  opts.EventEmitter,
		UseCaseHelper: opts.UseCaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: paperfacture.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		UnitPaperFacture: paperfacture.UnitPaperFactureOptions{
			Dictionary: paperFactureDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewPaperFactureAvailabilityAPI - создаёт объект usecase.PaperFacture.
func NewPaperFactureAvailabilityAPI(ctx context.Context, opts app.Options) (*usecase.PaperFacture, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries paper facture availability API")

	return availability.NewPaperFacture(opts.PostgresConnManager, opts.UseCaseErrorWrapper), nil
}

// RegisterPaperFactureErrors - comment func.
func RegisterPaperFactureErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.PaperFactureErrors()))
}
