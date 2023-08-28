package mrrabbitmq

import (
    "print-shop-back/pkg/mrerr"
)

func (c *Connection) wrapError(err error) error {
    return mrerr.ErrStorageQueryFailed.Caller(2).Wrap(err)
}
