package factory

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlocale"
	"github.com/mondegor/go-sysmess/mrlocale/provider/gotext"
	"github.com/mondegor/go-sysmess/mrlog"
	"golang.org/x/text/language"

	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/localization/dict/catalog/boxescat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/catalog/laminatescat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/catalog/paperscat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/dictionaries/materialtypescat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/dictionaries/papercolorscat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/dictionaries/paperfacturescat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/dictionaries/printformatscat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/errcat"
	"github.com/mondegor/print-shop-back/internal/localization/dict/msgcat"
)

// LocalePool - создаёт объект mrlang.LocalePool.
func LocalePool(logger mrlog.Logger, cfg config.Config) (*mrlocale.Pool, error) {
	mrlog.Info(logger, "Create and init language translator")

	var (
		localeProvider mrlocale.MessageProvider
		err            error
	)

	bundle, err := mrlocale.NewBundle(
		mrlocale.WithLanguages(cfg.Localization.Languages...),
		mrlocale.WithFormatMessage(gotext.MessageConverter("{", "}")),
		mrlocale.WithFormatError(mr.ErrorToMessage()),
		mrlocale.WithMessageProvider(
			func(languages []language.Tag) (mrlocale.MessageProvider, error) {
				localeProvider, err = gotext.NewProvider(
					gotext.WithLanguages(languages...),
					gotext.WithCatalog(mrlocale.DefaultMessagesDomain, msgcat.NewCatalog()),
					gotext.WithCatalog(mrlocale.DefaultErrorsDomain, errcat.NewCatalog()),
					gotext.WithCatalog("catalog.boxes", boxescat.NewCatalog()),
					gotext.WithCatalog("catalog.laminates", laminatescat.NewCatalog()),
					gotext.WithCatalog("catalog.papers", paperscat.NewCatalog()),
					gotext.WithCatalog("dictionaries.materialtypes", materialtypescat.NewCatalog()),
					gotext.WithCatalog("dictionaries.papercolors", papercolorscat.NewCatalog()),
					gotext.WithCatalog("dictionaries.paperfactures", paperfacturescat.NewCatalog()),
					gotext.WithCatalog("dictionaries.printformats", printformatscat.NewCatalog()),
				)

				return localeProvider, err
			},
		),
	)
	if err != nil {
		return nil, err
	}

	mrlog.DebugFunc(
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
