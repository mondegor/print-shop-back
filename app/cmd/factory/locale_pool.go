package factory

import (
	"fmt"

	"github.com/mondegor/go-sysmess/errors/helper"
	"github.com/mondegor/go-sysmess/mrlocale"
	"github.com/mondegor/go-sysmess/mrlocale/provider/gotext"
	"golang.org/x/text/language"

	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/localization/dict/catalog/boxescat"
	"print-shop-back/internal/localization/dict/catalog/laminatescat"
	"print-shop-back/internal/localization/dict/catalog/paperscat"
	"print-shop-back/internal/localization/dict/dictionaries/materialtypescat"
	"print-shop-back/internal/localization/dict/dictionaries/papercolorscat"
	"print-shop-back/internal/localization/dict/dictionaries/paperfacturescat"
	"print-shop-back/internal/localization/dict/dictionaries/printformatscat"
	"print-shop-back/internal/localization/dict/errcat"
	"print-shop-back/internal/localization/dict/msgcat"
)

// LocalePool - создаёт объект mrlang.LocalePool.
func LocalePool(logger log.Logger, cfg config.Config) (*mrlocale.Pool, error) {
	log.Info(logger, "Create and init language translator")

	var (
		localeProvider mrlocale.MessageProvider
		err            error
	)

	bundle, err := mrlocale.NewBundle(
		cfg.AppLanguages,
		mrlocale.WithFormatMessage(gotext.MessageConverter("{", "}")),
		mrlocale.WithFormatError(helper.ExtractMessageForLocalization),
		mrlocale.WithMessageProvider(
			func(languages []language.Tag) (mrlocale.MessageProvider, error) {
				localeProvider, err = gotext.NewProvider(
					languages,
					gotext.WithDomainCatalog(mrlocale.DefaultMessagesDomain, msgcat.NewCatalog()),
					gotext.WithDomainCatalog(mrlocale.DefaultErrorsDomain, errcat.NewCatalog()),
					gotext.WithDomainCatalog("catalog.boxes", boxescat.NewCatalog()),
					gotext.WithDomainCatalog("catalog.laminates", laminatescat.NewCatalog()),
					gotext.WithDomainCatalog("catalog.papers", paperscat.NewCatalog()),
					gotext.WithDomainCatalog("dictionaries.materialtypes", materialtypescat.NewCatalog()),
					gotext.WithDomainCatalog("dictionaries.papercolors", papercolorscat.NewCatalog()),
					gotext.WithDomainCatalog("dictionaries.paperfactures", paperfacturescat.NewCatalog()),
					gotext.WithDomainCatalog("dictionaries.printformats", printformatscat.NewCatalog()),
				)

				return localeProvider, err
			},
		),
	)
	if err != nil {
		return nil, err
	}

	log.DebugFunc(
		logger,
		func() string {
			var buf []byte

			buf = append(buf, "Locales:\n"...)

			for _, lang := range localeProvider.Languages() {
				buf = append(buf, fmt.Sprintf("- Language=%s\n", lang)...)
			}

			buf = append(buf, "Locale domains:"...)

			for _, domain := range localeProvider.Domains() {
				buf = append(buf, "\n- "+domain...)
			}

			return string(buf)
		},
	)

	return mrlocale.NewPool(bundle), nil
}
