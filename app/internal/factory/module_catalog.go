package factory

import (
	"print-shop-back/internal/modules"
	"print-shop-back/internal/modules/catalog/factory"
)

func NewCatalogOptions(opts *modules.Options) (*factory.Options, error) {
	boxDictionary, err := opts.Translator.Dictionary("catalog/box")

	if err != nil {
		return nil, err
	}

	laminateDictionary, err := opts.Translator.Dictionary("catalog/laminate")

	if err != nil {
		return nil, err
	}

	paperDictionary, err := opts.Translator.Dictionary("catalog/paper")

	if err != nil {
		return nil, err
	}

	return &factory.Options{
		Logger:          opts.Logger,
		EventBox:        opts.EventBox,
		ServiceHelper:   opts.ServiceHelper,
		PostgresAdapter: opts.PostgresAdapter,

		LaminateTypeAPI: opts.DictionariesLaminateTypeAPI,
		PaperColorAPI:   opts.DictionariesPaperColorAPI,
		PaperFactureAPI: opts.DictionariesPaperFactureAPI,

		UnitBox: &factory.UnitBoxOptions{
			Dictionary: boxDictionary,
		},

		UnitLaminate: &factory.UnitLaminateOptions{
			Dictionary: laminateDictionary,
		},

		UnitPaper: &factory.UnitPaperOptions{
			Dictionary: paperDictionary,
		},
	}, nil
}
