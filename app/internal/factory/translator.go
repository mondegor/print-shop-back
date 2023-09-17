package factory

import (
    "print-shop-back/config"

    "github.com/mondegor/go-sysmess/mrlang"
    "github.com/mondegor/go-webcore/mrcore"
)

func NewTranslator(cfg *config.Config, logger mrcore.Logger) (*mrlang.Translator, error) {
    logger.Info("Create and init language translator")

    return mrlang.NewTranslator(
        mrlang.TranslatorOptions{
            DirPath: cfg.Translation.DirPath,
            FileType: cfg.Translation.FileType,
            LangCodes: cfg.Translation.LangCodes,
        },
    )
}
