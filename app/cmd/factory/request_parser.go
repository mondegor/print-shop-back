package factory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrview"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"

	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/pkg/validate"
	"github.com/mondegor/print-shop-back/pkg/view"
)

// CreateRequestParsers - создаются и возвращаются парсеры запросов клиента.
func CreateRequestParsers(ctx context.Context, cfg config.Config) (app.RequestParsers, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msg("Create and init base request parsers")

	validator, err := NewValidator(ctx, cfg)
	if err != nil {
		return app.RequestParsers{}, err
	}

	// WARNING: функция использует контекст роутера chi,
	// поэтому её можно менять только при смене самого роутера
	pathFunc := mrchi.URLPathParam

	registeredMimeTypes := mrlib.NewMimeTypeList(logger, cfg.Validation.MimeTypes)

	jsonMimeTypes, err := registeredMimeTypes.MimeTypesByExts(cfg.Validation.Files.Json.Extensions)
	if err != nil {
		return app.RequestParsers{}, err
	}

	logoMimeTypes, err := registeredMimeTypes.MimeTypesByExts(cfg.Validation.Images.Logo.File.Extensions)
	if err != nil {
		return app.RequestParsers{}, err
	}

	parsers := app.RequestParsers{
		// Bool:       mrparser.NewBool(),
		// DateTime:   mrparser.NewDateTime(),
		Int64:      mrparser.NewInt64(pathFunc),
		ItemStatus: mrparser.NewItemStatus(),
		Uint64:     mrparser.NewUint64(pathFunc),
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
			logger,
			mrparser.WithFileMinSize(cfg.Validation.Files.Json.MinSize),
			mrparser.WithFileMaxSize(cfg.Validation.Files.Json.MaxSize),
			mrparser.WithFileMaxFiles(cfg.Validation.Files.Json.MaxFiles),
			mrparser.WithFileCheckRequestContentType(cfg.Validation.Files.Json.CheckRequestContentType),
			mrparser.WithFileAllowedMimeTypes(jsonMimeTypes),
		),
		ImageLogo: mrparser.NewImage(
			logger,
			mrparser.WithImageMaxWidth(cfg.Validation.Images.Logo.MaxWidth),
			mrparser.WithImageMaxHeight(cfg.Validation.Images.Logo.MaxHeight),
			mrparser.WithImageCheckBody(cfg.Validation.Images.Logo.CheckBody),
			mrparser.WithImageFileOptions(
				mrparser.WithFileMinSize(cfg.Validation.Images.Logo.File.MinSize),
				mrparser.WithFileMaxSize(cfg.Validation.Images.Logo.File.MaxSize),
				mrparser.WithFileMaxFiles(cfg.Validation.Images.Logo.File.MaxFiles),
				mrparser.WithFileCheckRequestContentType(cfg.Validation.Images.Logo.File.CheckRequestContentType),
				mrparser.WithFileAllowedMimeTypes(logoMimeTypes),
			),
		),
	}

	parsers.Parser = validate.NewParser(
		parsers.Int64,
		parsers.Uint64,
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

	if err := validator.Register("tag_2d_size", view.Validate2dSize); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_3d_size", view.Validate3dSize); err != nil {
		return nil, err
	}

	return validator, nil
}
