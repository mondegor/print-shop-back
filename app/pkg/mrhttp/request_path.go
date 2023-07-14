package mrhttp

import (
    "net/http"
    "strconv"

    "github.com/julienschmidt/httprouter"
)

type RequestPath struct {
    params httprouter.Params
}

func newRequestPath(r *http.Request) *RequestPath {
    params, ok := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)

    if !ok {
        params = nil
    }

    return &RequestPath{
        params: params,
    }
}

func (r *RequestPath) Get(name string) string {
    if r.params == nil {
        return ""
    }

    return r.params.ByName(name)
}

func (r *RequestPath) GetInt(name string) int64 {
    value, err := strconv.ParseInt(r.Get(name), 10, 64)

    if err != nil {
        return 0
    }

    return value
}
