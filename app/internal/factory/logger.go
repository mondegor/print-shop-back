package factory

import (
	"print-shop-back/config"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewLogger(cfg *config.Config) (*mrcore.LoggerAdapter, error) {
	mrcore.SetDebug(cfg.Debugging.Debug)

	mrerr.SetCallerOptions(
		mrerr.CallerDeep(cfg.Debugging.ErrorCaller.Deep),
		mrerr.CallerUseShortPath(cfg.Debugging.ErrorCaller.UseShortPath),
	)

	prefix := cfg.Log.Prefix

	if prefix != "" {
		prefix = "[" + prefix + "] "
	}

	logger, err := mrcore.NewLogger(
		mrcore.LoggerOptions{
			Prefix: prefix,
			Level:  cfg.Log.Level,
			CallerOptions: []mrerr.CallerOption{
				mrerr.CallerDeep(cfg.Log.LogCaller.Deep),
				mrerr.CallerUseShortPath(cfg.Log.LogCaller.UseShortPath),
			},
			CallerEnabledFunc: func(err error) bool {
				if appErr, ok := err.(*mrerr.AppError); ok {
					return appErr.Kind() == mrerr.ErrorKindUser
				}

				return true
			},
		},
	)

	if err != nil {
		return nil, err
	}

	mrcore.SetDefaultLogger(logger)

	logger.Info("%s, version: %s", cfg.AppName, cfg.AppVersion)

	if cfg.AppInfo != "" {
		logger.Info(cfg.AppInfo)
	}

	if mrcore.Debug() {
		logger.Info("DEBUG MODE: ON")
	}

	logger.Info("LOG LEVEL: %s", cfg.Log.Level)
	logger.Info("CONFIG PATH: %s", cfg.ConfigPath)
	logger.Info("APP PATH: %s", cfg.AppPath)

	return logger, nil
}
