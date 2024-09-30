package factory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/app"
)

// CreateResponseSenders - создаются и возвращаются компоненты для отправки ответа клиенту.
func CreateResponseSenders(ctx context.Context, _ config.Config) (app.ResponseSenders, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init base response senders")

	sender := mrresp.NewSender(mrjson.NewEncoder())

	return app.ResponseSenders{
		Sender:     mrresp.NewSender(mrjson.NewEncoder()),
		FileSender: mrresp.NewFileSender(sender),
	}, nil
}

// NewErrorResponseSender - создаёт объект mrresp.ErrorSender.
func NewErrorResponseSender(ctx context.Context, opts app.Options) (*mrresp.ErrorSender, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init error response sender")

	return mrresp.NewErrorSender(
		mrjson.NewEncoder(),
		opts.ErrorHandler,
		mrserver.NewHttpErrorStatusGetter(
			opts.Cfg.Debugging.UnexpectedHttpStatus,
		),
		opts.Cfg.Debugging.Debug,
	), nil
}
