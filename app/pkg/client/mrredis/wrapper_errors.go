package mrredis

import "print-shop-back/pkg/mrerr"

func (c *Connection) wrapError(err error) error {
    return mrerr.ErrStorageQueryFailed.Wrap(err)
}
