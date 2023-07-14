package mrcontext

import (
    "context"
    "fmt"
    "net/http"

    "github.com/google/uuid"
)

// f7479171-83d2-4f64-84ac-892f8c0aaf48
const correlationIdLen= 36

func CorrelationIdFromRequest(r *http.Request) (string, error) {
    value := r.Header.Get("CorrelationID")

    if value == "" {
        return genCorrelationId(), nil
    }

    if len(value) == correlationIdLen {
        return value, nil
    }

    return genCorrelationId(), ErrHttpRequestCorrelationID.New(value)
}

func CorrelationIdNewContext(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, ctxCorrelationIdKey, id)
}

func GetCorrelationId(ctx context.Context) string {
    value, ok := ctx.Value(ctxCorrelationIdKey).(string)

    if ok {
        return value
    }

    return genCorrelationId()
}

func genCorrelationId() string {
    return fmt.Sprintf("TR-%s", uuid.New().String())
}
