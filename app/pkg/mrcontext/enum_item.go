package mrcontext

import (
    "net/http"
    "regexp"
    "strings"
)

const maxEnumLen = 32

var regexpEnum = regexp.MustCompile(`^[A-Z]([A-Z0-9_]+)?[A-Z0-9]$`)

func EnumItemFromRequest(r *http.Request, key string) (string, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return "", nil
    }

    if len(value) > maxEnumLen {
        return "", ErrHttpRequestParamLen.New(key, maxEnumLen)
    }

    value = strings.ToUpper(strings.TrimSpace(value))

    if !regexpEnum.MatchString(value) {
        return "", ErrHttpRequestParseParam.New("enum", key, value)
    }

    return value, nil
}
