package mrredis

import "calc-user-data-back-adm/pkg/mrerr"

func (c *Connection) wrapError(err error) error {
    return mrerr.ErrStorageQueryFailed.Wrap(err)
}
