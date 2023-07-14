package mrcontext

import (
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrlang"
    "context"
    "net/http"
)

func AcceptLanguageFromRequest(r *http.Request) []mrlang.LangCode {
    acceptLanguage := r.Header.Get("Accept-Language")
    return mrlang.ParseAcceptLanguage(acceptLanguage)
}

func LocaleNewContext(ctx context.Context, locale mrapp.Locale) context.Context {
    return context.WithValue(ctx, ctxLocaleKey, locale)
}

func GetLocale(ctx context.Context) mrapp.Locale {
    value, ok := ctx.Value(ctxLocaleKey).(mrapp.Locale)

    if ok {
        return value
    }

    return mrapp.NewLocaleStub()
}
