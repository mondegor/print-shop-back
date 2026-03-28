package factory

import (
	"github.com/mondegor/go-components/mrauth/model/contactaddress"
	"github.com/mondegor/go-components/wire/mrauth/mapping"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/mime"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
	"github.com/mondegor/go-webcore/mrview"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"

	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/app"
	mrcalcvalidate "github.com/mondegor/print-shop-back/pkg/mrcalc/validate"
	"github.com/mondegor/print-shop-back/pkg/validate"
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

	registeredMimeTypes := mime.NewTypeList(cfgValidation.MimeTypes)

	jsonMimeTypes, err := registeredMimeTypes.TypesByExts(cfgValidation.Files.Json.Extensions)
	if err != nil {
		return app.RequestParsers{}, err
	}

	logoMimeTypes, err := registeredMimeTypes.TypesByExts(cfgValidation.Images.Logo.File.Extensions)
	if err != nil {
		return app.RequestParsers{}, err
	}

	parsers := app.RequestParsers{
		// Bool:       parser.NewBool(),
		// DateTime:   parser.NewDateTime(),
		Int64:      parser.NewInt64(opts.Logger),
		ItemStatus: parser.NewItemStatus(opts.Logger),
		Uint64:     parser.NewUint64(pathFunc, opts.Logger),
		ListSorter: parser.NewListSorter(opts.Logger, parser.ListSorterOptions{}),
		ListCursor: parser.NewListCursor(
			opts.Logger,
			parser.ListCursorOptions{
				LimitMax:     int(opts.Cfg.General.PageSizeMax),
				LimitDefault: int(opts.Cfg.General.PageSizeDefault),
			},
		),
		ListPager: parser.NewListPager(
			opts.Logger,
			parser.ListPagerOptions{
				PageSizeMax:     int(opts.Cfg.General.PageSizeMax),
				PageSizeDefault: int(opts.Cfg.General.PageSizeDefault),
			},
		),
		String:    parser.NewString(pathFunc, opts.Logger),
		UUID:      parser.NewUUID(pathFunc, opts.Logger),
		Validator: parser.NewValidator(mrjson.NewDecoder(), validator),
		ClientIP:  parser.NewClientIP(opts.Logger),
		User:      parser.NewUser(opts.Logger),
		Locale:    parser.NewLocale(opts.LocalePool, opts.Logger, opts.Cfg.Localization.LangURLParam),
		FileJson: parser.NewFile(
			opts.Logger,
			parser.WithFileMinSize(int64(cfgValidation.Files.Json.MinSize)),
			parser.WithFileMaxSize(int64(cfgValidation.Files.Json.MaxSize)),
			parser.WithFileMaxFiles(int(cfgValidation.Files.Json.MaxFiles)),
			parser.WithFileCheckRequestContentType(cfgValidation.Files.Json.CheckRequestContentType),
			parser.WithFileAllowedMimeTypes(jsonMimeTypes),
		),
		ImageLogo: parser.NewImage(
			opts.Logger,
			parser.WithImageMaxWidth(int32(cfgValidation.Images.Logo.MaxWidth)),
			parser.WithImageMaxHeight(int32(cfgValidation.Images.Logo.MaxHeight)),
			parser.WithImageCheckBody(cfgValidation.Images.Logo.CheckBody),
			parser.WithImageFileOpts(
				parser.WithFileMinSize(int64(cfgValidation.Images.Logo.File.MinSize)),
				parser.WithFileMaxSize(int64(cfgValidation.Images.Logo.File.MaxSize)),
				parser.WithFileMaxFiles(int(cfgValidation.Images.Logo.File.MaxFiles)),
				parser.WithFileCheckRequestContentType(cfgValidation.Images.Logo.File.CheckRequestContentType),
				parser.WithFileAllowedMimeTypes(logoMimeTypes),
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
		parsers.ListCursor,
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
func NewValidator(logger mrlog.Logger, cfg config.Config) (*mrplayvalidator.ValidatorAdapter, error) {
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
		mrview.NewValidateOR(
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

	if err := validator.Register("tag_password", mrview.ValidatePassword); err != nil { // TODO: переименовать в validate.Password
		return nil, err
	}

	if err := validator.Register(
		"tag_realm",
		mrview.NewValidateAND(
			mrview.ValidateName,
			mrview.NewValidateInArray(
				mapping.OptionUserRealmsToStringRealms(cfg.AccessControl.Realms),
			),
		),
	); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_2d_size", mrcalcvalidate.Size2d); err != nil {
		return nil, err
	}

	if err := validator.Register("tag_3d_size", mrcalcvalidate.Size3d); err != nil {
		return nil, err
	}

	return validator, nil
}
