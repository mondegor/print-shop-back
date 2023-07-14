package mrcontext

import "print-shop-back/pkg/mrerr"

var (
    ErrHttpRequestPlatformValue = mrerr.NewFactory(
        "errHttpRequestPlatformValue", mrerr.ErrorKindInternal, "header 'Platform' contains incorrect value '{{ .value }}'")

    ErrHttpRequestCorrelationID = mrerr.NewFactory(
        "errHttpRequestCorrelationID", mrerr.ErrorKindInternalNotice, "header 'CorrelationID' contains incorrect value '{{ .value }}'")

    ErrHttpRequestUserIP = mrerr.NewFactory(
        "errHttpRequestUserIP", mrerr.ErrorKindInternal, "UserIP '{{ .value }}' is not IP:port")

    ErrHttpRequestParseUserIP = mrerr.NewFactory(
        "errHttpRequestParseUserIP", mrerr.ErrorKindInternal, "UserIP contains incorrect value '{{ .value }}'")

    ErrHttpRequestParseEnum = mrerr.NewFactory(
    "errHttpRequestParseEnum", mrerr.ErrorKindUser, "enum with key '{{ .key }}' contains incorrect value '{{ .value }}'")

    ErrHttpRequestEnumLen = mrerr.NewFactory(
        "errHttpRequestEnumLen", mrerr.ErrorKindUser, "enum with key '{{ .key }}' has value length greater then max '{{ .maxLength }}'")
)
