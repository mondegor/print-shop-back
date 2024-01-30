package factory

import (
	"context"
	"print-shop-back/config"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func NewResponseSender(ctx context.Context, cfg config.Config) (*mrresponse.Sender, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init base response sender")

	return mrresponse.NewSender(mrjson.NewEncoder()), nil
}

func NewErrorResponseSender(ctx context.Context, cfg config.Config) (*mrresponse.ErrorSender, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init error response sender")

	return mrresponse.NewErrorSender(mrjson.NewEncoder()), nil
}
