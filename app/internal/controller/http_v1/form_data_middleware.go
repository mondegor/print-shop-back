package http_v1

import (
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrcontext"
    "calc-user-data-back-adm/pkg/mrentity"
)

const ctxParentIdKey = mrcontext.CtxParentIdKey

func (f *FormFieldItem) FormDataMiddleware(next mrapp.HttpHandlerFunc) mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        r := c.Request()

        id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))
        err := f.serviceFormData.CheckAvailability(r.Context(), id)

        if err != nil {
            return err
        }

        mrcontext.IdNewContext(r.Context(), ctxParentIdKey, id)

        return next(c)
    }
}

func (f *FormFieldItem) getFormId(c mrapp.ClientData) mrentity.KeyInt32 {
    return mrcontext.GetId(c.Context(), ctxParentIdKey)
}
