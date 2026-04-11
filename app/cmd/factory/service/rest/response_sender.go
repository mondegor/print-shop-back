package rest

import (
	"net/http"

	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-sysmess/errors/runtime/hint"
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

	statusMapper, err := mrserver.NewHttpErrorStatusMapper(
		int(opts.Cfg.Debugging.UnexpectedHttpStatus),
		mrauth.ErrTokenNotFoundOrExpired.Code(), http.StatusUnauthorized,
		mrauth.ErrTokenInvalid.Code(), http.StatusForbidden,
	)
	if err != nil {
		return nil, err
	}

	return mrresp.NewErrorSender(
		mrjson.NewEncoder(),
		opts.ErrorHandler,
		func(err error) string {
			return hint.Extract(err).ErrorID()
		},
		opts.Logger,
		opts.RequestParsers.Locale,
		statusMapper,
		opts.DebugFunc,
	), nil
}
