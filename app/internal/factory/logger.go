package factory

import (
    "print-shop-back/config"

    "github.com/mondegor/go-webcore/mrcore"
)

func NewLogger(cfg *config.Config) (*mrcore.LoggerAdapter, error) {
    logger, err := mrcore.NewLogger("[" + cfg.AppName + "] ", cfg.Log.Level)

    if err != nil {
        return nil, err
    }

    mrcore.SetDefaultLogger(logger)

    logger.Info("APP %s, version: %s", cfg.AppName, cfg.AppVersion)

    if cfg.Debug {
        logger.Info("DEBUG MODE: ON")
    }

    logger.Info("CONFIG PATH: %s", cfg.ConfigPath)
    logger.Info("APP PATH: %s", cfg.AppPath)
    logger.Info("LOG LEVEL: %s", cfg.Log.Level)

    return logger, nil
}
