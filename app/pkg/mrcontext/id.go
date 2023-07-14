package mrcontext

import (
    "context"
    "print-shop-back/pkg/mrentity"
)

func IdNewContext(ctx context.Context, ctxKey any, id mrentity.KeyInt32) context.Context {
    return context.WithValue(ctx, ctxKey, id)
}

func GetId(ctx context.Context, ctxKey any) mrentity.KeyInt32 {
    value, ok := ctx.Value(ctxKey).(mrentity.KeyInt32)

    if ok {
        return value
    }

    return 0
}
