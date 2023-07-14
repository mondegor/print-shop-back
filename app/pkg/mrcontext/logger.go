package mrcontext

import (
    "context"
    "print-shop-back/pkg/mrapp"
)

func LoggerNewContext(ctx context.Context, logger mrapp.Logger) context.Context {
    return context.WithValue(ctx, ctxLoggerKey, logger)
}

func GetLogger(ctx context.Context) mrapp.Logger {
    value, ok := ctx.Value(ctxLoggerKey).(mrapp.Logger)

    if ok {
        return value
    }

    return mrapp.NewLoggerStub()
}
