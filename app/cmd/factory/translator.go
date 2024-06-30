package factory

import (
	"context"
	"fmt"

	"github.com/mondegor/print-shop-back/config"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewTranslator - создаёт объект mrlang.Translator.
func NewTranslator(ctx context.Context, cfg config.Config) (*mrlang.Translator, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msg("Create and init language translator")

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

	logger.Debug().MsgFunc(
		func() string {
			var buf []byte

			buf = append(buf, "Locales:\n"...)

			for _, localeCode := range tr.RegisteredLocales() {
				locale, _ := tr.LocaleByCode(localeCode)
				buf = append(buf, fmt.Sprintf("- ID=%d;code=%s\n", locale.LangID(), localeCode)...)
			}

			buf = append(buf, "Multi language dictionaries:"...)

			for _, dictName := range tr.RegisteredDictionaries() {
				buf = append(buf, "\n- "+dictName...)
			}

			return string(buf)
		},
	)

	return tr, nil
}
