package factory

import (
	"context"

	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/pkg/validate"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrview"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"
)

// CreateRequestParsers - создаются и возвращаются парсеры запросов клиента.
func CreateRequestParsers(ctx context.Context, cfg config.Config) (app.RequestParsers, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init base request parser")

	validator, err := NewValidator(ctx, cfg)
	if err != nil {
		return app.RequestParsers{}, err
	}

	// WARNING: функция использует контекст роутера chi,
	// поэтому её можно менять только при смене самого роутера
	pathFunc := mrchi.URLPathParam

	registeredMimeTypes := mrlib.NewMimeTypeList(cfg.MimeTypes)

	jsonMimeTypeList, err := registeredMimeTypes.NewListByExts(".json")
	if err != nil {
		return app.RequestParsers{}, err
	}

	imageMimeTypeList, err := registeredMimeTypes.NewListByExts(".jpeg", ".jpg", ".png")
	if err != nil {
		return app.RequestParsers{}, err
	}

	parsers := app.RequestParsers{
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
		String:    mrparser.NewString(pathFunc),
		UUID:      mrparser.NewUUID(pathFunc),
		Validator: mrparser.NewValidator(mrjson.NewDecoder(), validator),
		FileJson: mrparser.NewFile(
			mrparser.FileOptions{
				AllowedMimeTypes:        jsonMimeTypeList,
				MinSize:                 1,
				MaxSize:                 512 * 1024, // 512Kb
				MaxFiles:                4,
				CheckRequestContentType: true,
			},
		),
		ImageLogo: mrparser.NewImage(
			mrparser.ImageOptions{
				File: mrparser.FileOptions{
					AllowedMimeTypes:        imageMimeTypeList,
					MinSize:                 512,
					MaxSize:                 128 * 1024, // 128Kb
					CheckRequestContentType: true,
				},
				MaxWidth:  1024,
				MaxHeight: 1024,
				CheckBody: true,
			},
		),
	}

	parsers.Parser = validate.NewParser(
		parsers.Int64,
		parsers.KeyInt32,
		parsers.String,
		parsers.UUID,
		parsers.Validator,
	)

	parsers.ExtendParser = validate.NewExtendParser(
		parsers.Parser,
		parsers.ItemStatus,
		parsers.ListSorter,
		parsers.ListPager,
	)

	return parsers, nil
}

// NewValidator - создаёт объект mrplayvalidator.ValidatorAdapter.
func NewValidator(ctx context.Context, _ config.Config) (*mrplayvalidator.ValidatorAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init data validator")

	validator := mrplayvalidator.New()

	// registers custom tags for validation (see mrview.validator_tags.go)

	if err := validator.Register("tag_article", mrview.ValidateAnyNotSpaceSymbol); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_rewrite_name", mrview.ValidateRewriteName); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_variable", mrview.ValidateVariable); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_double_size", mrview.ValidateDoubleSize); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_triple_size", mrview.ValidateTripleSize); err != nil {
		return nil, err
	}

	return validator, nil
}
