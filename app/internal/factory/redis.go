package factory

import (
    "context"
    "print-shop-back/config"
    "time"

    "github.com/mondegor/go-storage/mrredis"
    "github.com/mondegor/go-webcore/mrcore"
)

func NewRedis(cfg *config.Config, logger mrcore.Logger) (*mrredis.ConnAdapter, error) {
   logger.Info("Create and init redis connection")

   opt := mrredis.Options{
       Host: cfg.Redis.Host,
       Port: cfg.Redis.Port,
       Password: cfg.Redis.Password,
       ConnTimeout: time.Duration(cfg.Redis.Timeout),
   }

   conn := mrredis.New()
   err := conn.Connect(opt)

   if err != nil {
       return nil, err
   }

   err = conn.Ping(context.Background())

   return conn, err
}
