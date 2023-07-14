package mrcontext

import (
    "context"
    "net/http"
)

const (
    PlatformMobile = "MOBILE"
    PlatformWeb = "WEB"
)

func PlatformFromRequest(r *http.Request) (string, error) {
    value := r.Header.Get("Platform")

    if value == "" || value == PlatformWeb {
        return PlatformWeb, nil
    }

    if value == PlatformMobile {
        return PlatformMobile, nil
    }

    return PlatformWeb, ErrHttpRequestPlatformValue.New(value)
}

func PlatformNewContext(ctx context.Context, platform string) context.Context {
    return context.WithValue(ctx, ctxPlatformKey, platform)
}

func GetPlatform(ctx context.Context) string {
    value, ok := ctx.Value(ctxPlatformKey).(string)

    if ok {
        return value
    }

    return PlatformWeb
}
