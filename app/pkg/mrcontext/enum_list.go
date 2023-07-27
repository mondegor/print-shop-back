package mrcontext

import (
    "context"
    "net/http"
    "strings"
)

const maxEnumListLen = 256

func EnumListFromRequest(r *http.Request, key string) ([]string, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return nil, nil
    }

    if len(value) > maxEnumListLen {
        return nil, ErrHttpRequestParamLen.New(key, maxEnumListLen)
    }

    var items []string

    for _, item := range strings.Split(strings.ToUpper(value), ",") {
        item = strings.TrimSpace(item)

        if !regexpEnum.MatchString(item) {
            return nil, ErrHttpRequestParseParam.New("enum", key, value)
        }

        items = append(items, item)
    }

    return items, nil
}

func EnumListNewContext(ctx context.Context, ctxKey any, items []string) context.Context {
    return context.WithValue(ctx, ctxKey, items)
}

func GetEnumList(ctx context.Context, ctxKey any) []string {
    value, ok := ctx.Value(ctxKey).([]string)

    if ok {
        return value
    }

    return []string{}
}
