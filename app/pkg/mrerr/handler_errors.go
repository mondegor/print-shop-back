package mrerr

var (
    ErrHttpRequestParseData = NewFactory(
        "errHttpRequestParseData", ErrorKindUser, "request body is not valid")

    ErrHttpResponseParseData = NewFactory(
        "errHttpResponseParseData", ErrorKindInternal, "response data is not valid")

    ErrHttpResponseSendData = NewFactory(
        "errHttpResponseSendData", ErrorKindInternal, "response data is not send")

    ErrHttpResponseSystemTemporarilyUnableToProcess = NewFactory(
       "errHttpResponseSystemTemporarilyUnableToProcess", ErrorKindUser, "the system is temporarily unable to process your request. please try again")

    ErrHttpResourceNotFound = NewFactory(
        "errHttpResourceNotFound", ErrorKindUser, "resource not found")
)
