package main

import (
    "flag"
    "print-shop-back/config"
    "print-shop-back/internal/app"
    "print-shop-back/pkg/mrlang"
    "print-shop-back/pkg/mrlib"
)

const appVersion = "v0.4.1"

var configPath string

func init() {
   flag.StringVar(&configPath,"config-path", "./config/config.yaml", "Path to application config file")
}

func main() {
    flag.Parse()

    cfg := config.New(configPath)
    logger := mrlib.NewLogger(cfg.Log.Level, !cfg.Log.NoColor)

    logger.Info("APP VERSION: %s", appVersion)

    if cfg.Debug {
      logger.Info("DEBUG MODE: ON")
    }

    logger.Info("LOG LEVEL: %s", cfg.Log.Level)
    logger.Info("APP PATH: %s", cfg.AppPath)
    logger.Info("CONFIG PATH: %s", configPath)

    translator := mrlib.NewTranslator(
        logger,
        mrlib.TranslatorOptions{
            DirPath: cfg.Translation.DirPath,
            FileType: cfg.Translation.FileType,
            LangCodes: mrlang.CastToLangCodes(cfg.Translation.LangCodes...),
        },
    )

    app.Run(cfg, logger, translator)
}
