package factory

import (
	"print-shop-back/config"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func NewResponseSender(cfg *config.Config, logger mrcore.Logger) (*mrresponse.Sender, error) {
	logger.Info("Create and init base response sender")

	return mrresponse.NewSender(mrjson.NewEncoder()), nil
}

func NewErrorResponseSender(cfg *config.Config, logger mrcore.Logger) (*mrresponse.ErrorSender, error) {
	logger.Info("Create and init error response sender")

	return mrresponse.NewErrorSender(mrjson.NewEncoder()), nil
}
