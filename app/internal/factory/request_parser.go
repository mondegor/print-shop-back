package factory

import (
	"print-shop-back/config"
	"print-shop-back/internal/modules"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrjulienrouter"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrview"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"
)

func NewRequestParsers(cfg *config.Config, logger mrcore.Logger) (*modules.RequestParsers, error) {
	logger.Info("Create and init base request parser")

	validator, err := NewValidator(cfg, logger)

	if err != nil {
		return nil, err
	}

	path := mrserver.RequestParserPathFunc(mrjulienrouter.PathParam)

	return &modules.RequestParsers{
		Path:       path,
		Base:       mrparser.NewBase(path),
		ItemStatus: mrparser.NewItemStatus(),
		KeyInt32:   mrparser.NewKeyInt32(path),
		SortPage:   mrparser.NewSortPage(),
		Validator:  mrparser.NewValidator(mrjson.NewDecoder(), validator),
	}, nil
}

func NewValidator(cfg *config.Config, logger mrcore.Logger) (*mrplayvalidator.ValidatorAdapter, error) {
	logger.Info("Create and init data validator")

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
