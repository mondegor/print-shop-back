package factory

import (
	"github.com/mondegor/go-components/mrauth/bag/contactaddress"
	"github.com/mondegor/go-sysmess/mrlib/extfile"
	"github.com/mondegor/go-sysmess/mrlog"
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
func CreateRequestParsers(opts app.Options) (app.RequestParsers, error) {
	mrlog.Info(opts.Logger, "Create and init base request parsers")

	validator, err := NewValidator(opts.Logger, opts.Cfg)
	if err != nil {
		return app.RequestParsers{}, err
	}

	cfgValidation := opts.Cfg.Validation

	// WARNING: функция использует контекст роутера chi,
	// поэтому её можно менять только при смене самого роутера
	pathFunc := mrchi.URLPathParam

	registeredMimeTypes := extfile.NewMimeTypeList(cfgValidation.MimeTypes)

	jsonMimeTypes, err := registeredMimeTypes.MimeTypesByExts(cfgValidation.Files.Json.Extensions)
	if err != nil {
		return app.RequestParsers{}, err
	}

	logoMimeTypes, err := registeredMimeTypes.MimeTypesByExts(cfgValidation.Images.Logo.File.Extensions)
	if err != nil {
		return app.RequestParsers{}, err
	}

	parsers := app.RequestParsers{
		// Bool:       mrparser.NewBool(),
		// DateTime:   mrparser.NewDateTime(),
		Int64:      mrparser.NewInt64(pathFunc, opts.Logger),
		ItemStatus: mrparser.NewItemStatus(opts.Logger),
		Uint64:     mrparser.NewUint64(pathFunc, opts.Logger),
		ListSorter: mrparser.NewListSorter(opts.Logger, mrparser.ListSorterOptions{}),
		ListPager: mrparser.NewListPager(
			opts.Logger,
			mrparser.ListPagerOptions{
				PageSizeMax:     opts.Cfg.General.PageSizeMax,
				PageSizeDefault: opts.Cfg.General.PageSizeDefault,
			},
		),
		String:    mrparser.NewString(pathFunc, opts.Logger),
		UUID:      mrparser.NewUUID(pathFunc, opts.Logger),
		Validator: mrparser.NewValidator(mrjson.NewDecoder(), validator),
		ClientIP:  mrparser.NewClientIP(opts.Logger),
		User:      mrparser.NewUser(opts.Logger),
		Locale:    mrparser.NewLocale(opts.LocalePool, opts.Logger, opts.Cfg.Localization.LangURLParam),
		FileJson: mrparser.NewFile(
			opts.Logger,
			mrparser.WithFileMinSize(cfgValidation.Files.Json.MinSize),
			mrparser.WithFileMaxSize(cfgValidation.Files.Json.MaxSize),
			mrparser.WithFileMaxFiles(cfgValidation.Files.Json.MaxFiles),
			mrparser.WithFileCheckRequestContentType(cfgValidation.Files.Json.CheckRequestContentType),
			mrparser.WithFileAllowedMimeTypes(jsonMimeTypes),
		),
		ImageLogo: mrparser.NewImage(
			opts.Logger,
			mrparser.WithImageMaxWidth(cfgValidation.Images.Logo.MaxWidth),
			mrparser.WithImageMaxHeight(cfgValidation.Images.Logo.MaxHeight),
			mrparser.WithImageCheckBody(cfgValidation.Images.Logo.CheckBody),
			mrparser.WithImageFileOptions(
				mrparser.WithFileMinSize(cfgValidation.Images.Logo.File.MinSize),
				mrparser.WithFileMaxSize(cfgValidation.Images.Logo.File.MaxSize),
				mrparser.WithFileMaxFiles(cfgValidation.Images.Logo.File.MaxFiles),
				mrparser.WithFileCheckRequestContentType(cfgValidation.Images.Logo.File.CheckRequestContentType),
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
		parsers.ClientIP,
		parsers.User,
		parsers.Locale,
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
func NewValidator(logger mrlog.Logger, _ config.Config) (*mrplayvalidator.ValidatorAdapter, error) {
	mrlog.Info(logger, "Create and init data validator")

	validator := mrplayvalidator.New(logger)

	// registers custom tags for validation (see mrview.validator_tags.go)

	if err := validator.Register("tag_article", mrview.ValidateAnyNotSpaceSymbol); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_email", contactaddress.ValidateEmail); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_phone", contactaddress.ValidatePhone); err != nil {
		return nil, err
	}

	if err := validator.Register(
		"tag_email_phone",
		mrview.ValidateOr(
			contactaddress.ValidateEmail,
			contactaddress.ValidatePhoneWorld,
		)); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_variable", mrview.ValidateVariable); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_name", mrview.ValidateName); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_rewrite_name", mrview.ValidateRewriteName); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_password", mrview.ValidatePassword); err != nil {
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
