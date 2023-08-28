package mrredis

import (
    "print-shop-back/pkg/mrerr"

    "github.com/redis/go-redis/v9"
)

func (c *Connection) wrapError(err error) error {
    if err == redis.Nil {
        return mrerr.ErrStorageNoRowFound.Caller(2).Wrap(err)
    }

    return mrerr.ErrStorageQueryFailed.Caller(2).Wrap(err)
}
