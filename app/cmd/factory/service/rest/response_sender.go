package rest

import (
	"net/http"

	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
)

// CreateResponseSenders - создаются и возвращаются компоненты для отправки ответа клиенту.
func CreateResponseSenders(logger mrlog.Logger) (app.ResponseSenders, error) {
	mrlog.Info(logger, "Create and init base response senders")

	sender := mrresp.NewSender(mrjson.NewEncoder())

	return app.ResponseSenders{
		Sender:     mrresp.NewSender(mrjson.NewEncoder()),
		FileSender: mrresp.NewFileSender(sender, logger),
	}, nil
}

// NewErrorResponseSender - создаёт объект mrresp.ErrorSender.
func NewErrorResponseSender(opts app.Options) (*mrresp.ErrorSender, error) {
	mrlog.Info(opts.Logger, "Create and init error response sender")

	return mrresp.NewErrorSender(
		mrjson.NewEncoder(),
		opts.ErrorHandler,
		opts.Logger,
		opts.RequestParsers.Locale,
		mrserver.NewHttpErrorStatusMapper(
			opts.Cfg.Debugging.UnexpectedHttpStatus,
			mrauth.ErrTokenNotFoundOrExpired.Code(), http.StatusUnauthorized,
		),
		opts.Cfg.Debugging.Debug,
	), nil
}
