package mrcontext

type ctxKey uint8

const (
    _ ctxKey = iota
    ctxCorrelationIdKey
    ctxLoggerKey
    ctxLocaleKey
    ctxPlatformKey
    ctxUserIPKey

    CtxParentIdKey
)
