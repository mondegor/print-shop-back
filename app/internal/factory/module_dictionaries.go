package factory

import (
	"print-shop-back/internal/modules"
	"print-shop-back/internal/modules/dictionaries/factory"
	factory_api "print-shop-back/internal/modules/dictionaries/factory/api"
	usecase_api "print-shop-back/internal/modules/dictionaries/usecase/api"
)

func NewDictionariesOptions(opts *modules.Options) (*factory.Options, error) {
	laminateTypeDictionary, err := opts.Translator.Dictionary("dictionaries/laminate-type")

	if err != nil {
		return nil, err
	}

	paperColorDictionary, err := opts.Translator.Dictionary("dictionaries/paper-color")

	if err != nil {
		return nil, err
	}

	paperFactureDictionary, err := opts.Translator.Dictionary("dictionaries/paper-facture")

	if err != nil {
		return nil, err
	}

	printFormatDictionary, err := opts.Translator.Dictionary("dictionaries/print-format")

	if err != nil {
		return nil, err
	}

	return &factory.Options{
		Logger:          opts.Logger,
		EventBox:        opts.EventBox,
		ServiceHelper:   opts.ServiceHelper,
		PostgresAdapter: opts.PostgresAdapter,

		UnitLaminateType: &factory.UnitLaminateTypeOptions{
			Dictionary: laminateTypeDictionary,
		},

		UnitPaperColor: &factory.UnitPaperColorOptions{
			Dictionary: paperColorDictionary,
		},

		UnitPaperFacture: &factory.UnitPaperFactureOptions{
			Dictionary: paperFactureDictionary,
		},

		UnitPrintFormatFacture: &factory.UnitPrintFormatOptions{
			Dictionary: printFormatDictionary,
		},
	}, nil
}

func NewDictionariesLaminateTypeAPI(opts *modules.Options) (*usecase_api.LaminateType, error) {
	opts.Logger.Info("Create and init dictionaries laminate type API")

	return factory_api.NewLaminateType(opts.PostgresAdapter, opts.ServiceHelper), nil
}

func NewDictionariesPaperColorAPI(opts *modules.Options) (*usecase_api.PaperColor, error) {
	opts.Logger.Info("Create and init dictionaries paper color API")

	return factory_api.NewPaperColor(opts.PostgresAdapter, opts.ServiceHelper), nil
}

func NewDictionariesPaperFactureAPI(opts *modules.Options) (*usecase_api.PaperFacture, error) {
	opts.Logger.Info("Create and init dictionaries paper facture API")

	return factory_api.NewPaperFacture(opts.PostgresAdapter, opts.ServiceHelper), nil
}

func NewDictionariesPrintFormatAPI(opts *modules.Options) (*usecase_api.PrintFormat, error) {
	opts.Logger.Info("Create and init dictionaries print format API")

	return factory_api.NewPrintFormat(opts.PostgresAdapter, opts.ServiceHelper), nil
}
