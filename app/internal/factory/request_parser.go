package factory

import (
	"context"
	"print-shop-back/config"
	"print-shop-back/internal"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrjulienrouter"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrview"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"
)

func CreateRequestParsers(ctx context.Context, cfg config.Config) (app.RequestParsers, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init base request parser")

	validator, err := NewValidator(ctx, cfg)

	if err != nil {
		return app.RequestParsers{}, err
	}

	pathFunc := mrjulienrouter.PathParam

	return app.RequestParsers{
		// Bool:       mrparser.NewBool(),
		// DateTime:   mrparser.NewDateTime(),
		Int64:      mrparser.NewInt64(pathFunc),
		ItemStatus: mrparser.NewItemStatus(),
		KeyInt32:   mrparser.NewKeyInt32(pathFunc),
		ListSorter: mrparser.NewListSorter(mrparser.ListSorterOptions{}),
		ListPager: mrparser.NewListPager(
			mrparser.ListPagerOptions{
				PageSizeMax:     cfg.General.PageSizeMax,
				PageSizeDefault: cfg.General.PageSizeDefault,
			},
		),
		String: mrparser.NewString(pathFunc),
		// UUID:       mrparser.NewUUID(pathFunc),
		Validator: mrparser.NewValidator(mrjson.NewDecoder(), validator),
		FileJson: mrparser.NewFile(
			mrparser.FileOptions{
				AllowedExts:             []string{".json"},
				MinSize:                 512,
				MaxSize:                 1024 * 1024,
				MaxFiles:                4,
				CheckRequestContentType: true,
			},
		),
		Image: mrparser.NewImage(mrparser.ImageOptions{}),
	}, nil
}

func NewValidator(ctx context.Context, cfg config.Config) (*mrplayvalidator.ValidatorAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init data validator")

	validator := mrplayvalidator.New()

	// registers custom tags for validation (see mrview.validator_tags.go)

	if err := validator.Register("tag_article", mrview.ValidateAnyNotSpaceSymbol); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_variable", mrview.ValidateVariable); err != nil {
		return nil, err
	}

	return validator, nil
}
