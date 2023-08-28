package mrhttp

import (
    "net/http"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
)

func (rt *Router) MiddlewareFirst() mrapp.HttpMiddleware {
    return mrapp.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            correlationId, err := mrcontext.CorrelationIdFromRequest(r)
            logger := rt.logger.WithContext(correlationId)
            logger.Debug("Exec MiddlewareFirst")

            if err != nil {
                logger.Warn(err.Error())
            }

            logger.Info("CorrelationID: %s", correlationId)
            ctx := mrcontext.CorrelationIdNewContext(r.Context(), correlationId)
            ctx = mrcontext.LoggerNewContext(ctx, logger)

            w.Header().Set("Content-Type", "application/json")
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func (rt *Router) MiddlewareAcceptLanguage(translator mrapp.Translator) mrapp.HttpMiddleware {
    return mrapp.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrcontext.GetLogger(r.Context())
            logger.Debug("Exec MiddlewareAcceptLanguage")

            acceptLanguages := mrcontext.AcceptLanguageFromRequest(r)
            locale := translator.GetLocale(acceptLanguages...)
            logger.Info("Accept-Language: %v; Set-Language: %s", acceptLanguages, locale.GetLang())
            ctx := mrcontext.LocaleNewContext(r.Context(), locale)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func (rt *Router) MiddlewarePlatform() mrapp.HttpMiddleware {
    return mrapp.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrcontext.GetLogger(r.Context())
            logger.Debug("Exec MiddlewarePlatform")

            platform, err := mrcontext.PlatformFromRequest(r)

            if err != nil {
                logger.Warn(err.Error())
            }

            logger.Info("Platform: %s", platform)
            ctx := mrcontext.PlatformNewContext(r.Context(), platform)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func (rt *Router) MiddlewareUserIp() mrapp.HttpMiddleware {
    return mrapp.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrcontext.GetLogger(r.Context())
            logger.Debug("Exec MiddlewareUserIp")

            userIp, err := mrcontext.UserIpFromRequest(r)

            if err != nil {
                logger.Warn(err.Error())
            }

            logger.Info("UserIp: %s", userIp.String())
            ctx := mrcontext.UserIpNewContext(r.Context(), userIp)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func (rt *Router) MiddlewareAuthenticateUser() mrapp.HttpMiddleware {
    return mrapp.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrcontext.GetLogger(r.Context())
            logger.Debug("Exec MiddlewareAuthenticateUser")

            next.ServeHTTP(w, r)
        })
    })
}

func (rt *Router) MiddlewareLast(next mrapp.HttpHandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logger := mrcontext.GetLogger(r.Context())
        logger.Debug("Exec MiddlewareLast")

        c := clientContext{
            request: r,
            responseWriter: w,
            validator: rt.validator,
        }

        err := next(&c)

        if err != nil {
            c.sendErrorResponse(err)
        }
    }
}
