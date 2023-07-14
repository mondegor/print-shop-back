package mrapp

import (
    "context"
    "net/http"
)

type (
    ClientData interface {
        Request() *http.Request
        Context() context.Context
        RequestPath() RequestPath
        Writer() http.ResponseWriter

        CorrelationId() string
        Logger() Logger
        Locale() Locale

        Parse(structRequest any) error
        Validate(structRequest any) error
        ParseAndValidate(structRequest any) error

        SendResponseNoContent() error
        SendResponse(status int, structResponse any) error
        SendResponseWithError(err error)
    }

    RequestPath interface {
        Get(name string) string
        GetInt(name string) int64
    }
)
