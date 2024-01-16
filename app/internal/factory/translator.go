package factory

import (
	"print-shop-back/config"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewTranslator(cfg *config.Config, logger mrcore.Logger) (*mrlang.Translator, error) {
	logger.Info("Create and init language translator")

	tr, err := mrlang.NewTranslator(
		mrlang.TranslatorOptions{
			DirPath:           cfg.Translation.DirPath,
			LangCodes:         cfg.Translation.LangCodes,
			DictionaryDirPath: cfg.Translation.Dictionaries.DirPath,
			Dictionaries:      cfg.Translation.Dictionaries.List,
		},
	)

	if err != nil {
		return nil, err
	}

	logger.Debug("Locales:\n")

	for _, localeCode := range tr.RegisteredLocales() {
		locale, _ := tr.LocaleByCode(localeCode)
		logger.Debug("- ID=%d;code=%s\n", locale.LangID(), localeCode)
	}

	logger.Debug("Multi language dictionaries:\n")

	for _, dictName := range tr.RegisteredDictionaries() {
		logger.Debug("- %s\n", dictName)
	}

	return tr, nil
}
