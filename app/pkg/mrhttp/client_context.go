package mrhttp

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrerr"
    "print-shop-back/pkg/mrlib"
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

func (c *clientContext) WithContext(ctx context.Context) mrapp.ClientData {
    c.request = c.request.WithContext(ctx)

    return c
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
    err := dec.Decode(&structRequest)

    if err != nil {
        return mrerr.ErrHttpRequestParseData.Wrap(err)
    }

    return nil
}

func (c *clientContext) Validate(structRequest any) error {
    return c.validator.Validate(c.request.Context(), structRequest)
}

func (c *clientContext) ParseAndValidate(structRequest any) error {
    err := c.Parse(structRequest)

    if err != nil {
        return err
    }

    return c.Validate(structRequest)
}

func (c *clientContext) SendResponse(status int, structResponse any) error {
    appError := c.sendResponse(status, structResponse)

    if appError != nil {
        return appError
    }

    return nil
}

func (c *clientContext) sendResponse(status int, structResponse any) *mrerr.AppError {
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

func (c *clientContext) SendResponseNoContent() error {
    c.responseWriter.WriteHeader(http.StatusNoContent)

    return nil
}

func (c *clientContext) sendErrorResponse(err error) {
    var appError *mrerr.AppError

    for { // only for break
        if userErrorList, ok := err.(*mrlib.UserErrorList); ok {
            appError = c.sendUserErrorListResponse(userErrorList)
            break
        }

        if appErrorTmp, ok := err.(*mrerr.AppError); ok {
            if appErrorTmp.Kind() == mrerr.ErrorKindUser {
                appError = c.sendUserErrorResponse(appErrorTmp)
                break
            }

            appError = appErrorTmp
            break
        }

        appError = mrerr.ErrInternal.Wrap(err)
        break
    }

    if appError != nil {
        c.Logger().Error(appError)
        c.sendAppErrorResponse(c.wrapErrorFunc(appError))
    }
}

func (c *clientContext) sendUserErrorListResponse(list *mrlib.UserErrorList) *mrerr.AppError {
    locale := c.Locale()
    errorResponseList := AppErrorListResponse{}

    for _, userError := range *list {
        if userError.Err.Kind() != mrerr.ErrorKindUser {
            c.Logger().Error(userError.Err)
            continue
        }

        errorResponseList.Add(
            userError.Id,
            userError.Err.GetUserMessage(locale).Reason,
        )
    }

    return c.sendResponse(http.StatusBadRequest, errorResponseList)
}

func (c *clientContext) sendUserErrorResponse(appError *mrerr.AppError) *mrerr.AppError {
    return c.sendResponse(
        http.StatusBadRequest,
        AppErrorListResponse{
            AppErrorAttribute{
                Id: AppErrorAttributeNameSystem,
                Value: appError.GetUserMessage(c.Locale()).Reason,
            },
        },
    )
}

func (c *clientContext) sendAppErrorResponse(status int, appError *mrerr.AppError) {
    c.responseWriter.Header().Set("Content-Type", "application/problem+json")
    c.responseWriter.WriteHeader(status)

    errMessage := appError.GetUserMessage(c.Locale())
    errorResponse := AppErrorResponse{
        Title: errMessage.Reason,
        Details: errMessage.DetailsToString(),
        Request: c.request.URL.Path,
        Time: time.Now().Format(time.RFC3339),
        ErrorTraceId: c.getErrorTraceId(appError),
    }

    c.responseWriter.Write(errorResponse.Marshal())
}

func (c *clientContext) getErrorTraceId(err *mrerr.AppError) string {
    errorTraceId := err.EventId()

    if errorTraceId == "" {
        return c.CorrelationId()
    }

    return fmt.Sprintf("%s, %s", c.CorrelationId(), err.EventId())
}

// :TODO: move to package internal
func (c *clientContext) wrapErrorFunc(err *mrerr.AppError) (int, *mrerr.AppError) {
    status := http.StatusInternalServerError

    if mrerr.ErrServiceEntityNotFound.Is(err) {
        status = http.StatusNotFound
        err = mrerr.ErrHttpResourceNotFound.Wrap(err)
    } else if mrerr.ErrServiceEntityTemporarilyUnavailable.Is(err) {
        err = mrerr.ErrHttpResponseSystemTemporarilyUnableToProcess.Wrap(err)
    } else if err.Code() == mrerr.ErrorCodeInternal {
        status = http.StatusTeapot
    }

    return status, err
}
