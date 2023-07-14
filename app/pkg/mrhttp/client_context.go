package mrhttp

import (
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrerr"
    "print-shop-back/pkg/mrlib"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type clientContext struct {
    request *http.Request
    responseWriter http.ResponseWriter
    requestPath *RequestPath
    validator mrapp.Validator
}

func (c *clientContext) Request() *http.Request {
    return c.request
}

func (c *clientContext) Context() context.Context {
    return c.request.Context()
}

func (c *clientContext) RequestPath() mrapp.RequestPath {
    if c.requestPath == nil {
        c.requestPath = newRequestPath(c.request)
    }

    return c.requestPath
}

func (c *clientContext) Writer() http.ResponseWriter {
    return c.responseWriter
}

func (c *clientContext) CorrelationId() string {
    return mrcontext.GetCorrelationId(c.request.Context())
}

func (c *clientContext) Logger() mrapp.Logger {
    return mrcontext.GetLogger(c.request.Context())
}

func (c *clientContext) Locale() mrapp.Locale {
    return mrcontext.GetLocale(c.request.Context())
}

func (c *clientContext) Parse(structRequest any) error {
    dec := json.NewDecoder(c.request.Body)
    dec.DisallowUnknownFields()

    if err := dec.Decode(&structRequest); err != nil {
        return mrerr.ErrHttpRequestParseData.Wrap(err)
    }

    return nil
}

func (c *clientContext) Validate(structRequest any) error {
    return c.validator.Validate(c.request.Context(), structRequest)
}

func (c *clientContext) ParseAndValidate(structRequest any) error {
    if err := c.Parse(structRequest); err != nil {
        return err
    }

    if err := c.Validate(structRequest); err != nil {
        return err
    }

    return nil
}

func (c *clientContext) SendResponseNoContent() error {
    c.responseWriter.WriteHeader(http.StatusNoContent)

    return nil
}

func (c *clientContext) SendResponse(status int, structResponse any) error {
    c.responseWriter.WriteHeader(status)

    bytes, err := json.Marshal(structResponse)

    if err != nil {
        return mrerr.ErrHttpResponseParseData.Wrap(err)
    }

    _, err = c.responseWriter.Write(bytes)

    if err != nil {
        return mrerr.ErrHttpResponseSendData.Wrap(err)
    }

    return nil
}

func (c *clientContext) SendResponseWithError(err error) {
    locale := c.Locale()

    if userErrorList, ok := err.(*mrlib.UserErrorList); ok {
        errorResponseList := AppErrorListResponse{}

        for _, userError := range *userErrorList {
            errorResponseList.Add(userError.Id, userError.Err.GetUserMessage(locale).Reason)
        }

        if err = c.SendResponse(http.StatusBadRequest, errorResponseList); err == nil {
            return
        }
    }

    statusCode, appError := c.getStatusAndAppError(err)

    c.Logger().Error(err)
    c.responseWriter.Header().Set("Content-Type", "application/problem+json")
    c.responseWriter.WriteHeader(statusCode)

    errMessage := appError.GetUserMessage(locale)
    errorResponse := AppErrorResponse{
        Title: errMessage.Reason,
        Details: errMessage.DetailsToString(),
        Request: c.request.URL.Path,
        Time: time.Now().Format(time.RFC3339),
        ErrorTraceId: c.getErrorTraceId(appError),
    }

    c.responseWriter.Write(errorResponse.Marshal())
}

func (c *clientContext) getStatusAndAppError(err error) (int, *mrerr.AppError) {
    statusCode := http.StatusTeapot

    if mrerr.ErrServiceEntityNotFound.Is(err) {
        statusCode = http.StatusNotFound
        err = mrerr.ErrHttpResourceNotFound.Wrap(err)
    } else if mrerr.ErrServiceEntityTemporarilyUnavailable.Is(err) {
        statusCode = http.StatusInternalServerError
        err = mrerr.ErrHttpResponseSystemTemporarilyUnableToProcess.Wrap(err)
    } else {
        if _, ok := err.(*mrerr.AppError); ok {
            statusCode = http.StatusInternalServerError
        } else {
            err = mrerr.ErrInternal.Wrap(err)
        }
    }

    return statusCode, err.(*mrerr.AppError)
}

func (c *clientContext) getErrorTraceId(err *mrerr.AppError) string {
    errorTraceId := err.EventId()

    if errorTraceId == "" {
        return c.CorrelationId()
    }

    return fmt.Sprintf("%s, %s", c.CorrelationId(), err.EventId())
}
