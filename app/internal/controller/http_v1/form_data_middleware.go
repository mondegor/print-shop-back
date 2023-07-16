package http_v1

import (
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrentity"
)

const ctxParentIdKey = mrcontext.CtxParentIdKey

func (f *FormFieldItem) FormDataMiddleware(next mrapp.HttpHandlerFunc) mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))
        err := f.serviceFormData.CheckAvailability(c.Request().Context(), id)

        if err != nil {
            return err
        }

        return next(c.WithContext(mrcontext.IdNewContext(c.Request().Context(), ctxParentIdKey, id)))
    }
}

func (f *FormFieldItem) getFormId(c mrapp.ClientData) mrentity.KeyInt32 {
    return mrcontext.GetId(c.Context(), ctxParentIdKey)
}
