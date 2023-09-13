package http_v1

import (
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

type ctxFormIdKey struct {}

func (ht *FormFieldItem) FormDataMiddleware(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        id := c.RequestPath().GetInt("fid")
        err := ht.serviceFormData.CheckAvailability(c.Request().Context(), mrentity.KeyInt32(id))

        if err != nil {
            return err
        }

        return next(c.WithContext(mrctx.WithInt64(c.Request().Context(), ctxFormIdKey{}, id)))
    }
}

func (ht *FormFieldItem) getFormId(c mrcore.ClientData) mrentity.KeyInt32 {
    return mrentity.KeyInt32(mrctx.Int64(c.Context(), ctxFormIdKey{}))
}
