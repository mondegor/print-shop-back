package dictionaries

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/api/availability/usecase"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor"
	"github.com/mondegor/print-shop-back/internal/factory/dictionaries/papercolor/api/availability"
)

// NewPaperColorModuleOptions - создаёт объект papercolor.Options.
func NewPaperColorModuleOptions(_ context.Context, opts app.Options) (papercolor.Options, error) {
	paperColorDictionary, err := opts.Translator.Dictionary("dictionaries/paper-colors")
	if err != nil {
		return papercolor.Options{}, err
	}

	return papercolor.Options{
		EventEmitter:        opts.EventEmitter,
		UseCaseErrorWrapper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: papercolor.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		UnitPaperColor: papercolor.UnitPaperColorOptions{
			Dictionary: paperColorDictionary,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewPaperColorAvailabilityAPI - создаёт объект usecase.PaperColor.
func NewPaperColorAvailabilityAPI(ctx context.Context, opts app.Options) (*usecase.PaperColor, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init dictionaries paper color availability API")

	return availability.NewPaperColor(opts.PostgresConnManager), nil
}
