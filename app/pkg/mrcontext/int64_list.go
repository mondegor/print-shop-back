package mrcontext

import (
    "context"
    "net/http"
    "strconv"
    "strings"
)

const maxInt64ListLen = 256

func Int64ListFromRequest(r *http.Request, key string) ([]int64, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return nil, nil
    }

    if len(value) > maxInt64ListLen {
        return nil, ErrHttpRequestParamLen.New(key, maxEnumListLen)
    }

    var items []int64

    for _, item := range strings.Split(value, ",") {
        item = strings.TrimSpace(item)

        i, err := strconv.ParseInt(item, 10, 64)

        if err != nil {
            return nil, ErrHttpRequestParseParam.New("int64", key, value)
        }

        items = append(items, i)
    }

    return items, nil
}

func Int64ListNewContext(ctx context.Context, ctxKey any, items []int64) context.Context {
    return context.WithValue(ctx, ctxKey, items)
}

func GetInt64List(ctx context.Context, ctxKey any) []int64 {
    value, ok := ctx.Value(ctxKey).([]int64)

    if ok {
        return value
    }

    return []int64{}
}
