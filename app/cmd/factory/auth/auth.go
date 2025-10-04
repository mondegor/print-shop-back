package auth

import (
	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

// NewAuthModuleOptions - создаёт объект auth.Options.
func NewAuthModuleOptions(opts app.Options) (auth.Options, error) {
	return auth.Options{
		Logger:              opts.Logger,
		EventEmitter:        opts.EventEmitter,
		UsecaseErrorWrapper: opts.UsecaseErrorWrapper,
		StorageErrorWrapper: opts.StorageErrorWrapper,
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: auth.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender:   opts.ResponseSenders.FileSender,
		NotifierAPI:      opts.NotifierAPI,
		Locker:           opts.Locker,
		WithDebugInfo:    opts.Cfg.Debugging.Debug,
		UserRealms:       opts.Cfg.AccessControl.Realms,
		OperationConfirm: opts.Cfg.AccessControl.OperationConfirm,
		JWT: auth.JWT{
			Method: opts.Cfg.AccessControl.JWTMethod,
			Secret: []byte(opts.Cfg.AccessControl.JWTSecret),
		},
	}, nil
}
