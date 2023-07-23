package mrerr

var (
    ErrInternal = NewFactory(
        ErrorCodeInternal, ErrorKindInternal, "internal server error")

    ErrInternalNilPointer = NewFactory(
        "errInternalNilPointer", ErrorKindInternal, "nil pointer")

    ErrInternalTypeAssertion = NewFactory(
        "errInternalTypeAssertion", ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

    ErrInternalInvalidType = NewFactory(
        "errInternalInvalidType", ErrorKindInternal, "invalid type '{{ .type1 }}', expected: '{{ .type2 }}'")

    ErrInternalInvalidData = NewFactory(
        "errInternalInvalidData", ErrorKindInternal, "invalid data '{{ .value }}'")

    ErrInternalParseData = NewFactory(
        "errInternalParseData", ErrorKindInternal, "data '{{ .name1 }}' parsed to {{ .name2 }} with error")

    ErrInternalMapValueNotFound = NewFactory(
        "errInternalMapValueNotFound", ErrorKindInternal, "'{{ .value }}' is not found in map {{ .name }}")

    ErrInternalMessageNotEnoughArguments = NewFactory(
        "errInternalMessageNotEnoughArguments", ErrorKindInternal, "not enough arguments in message '{{ .value }}'")

    ErrInternalMessageTooManyArguments = NewFactory(
        "errInternalMessageTooManyArguments", ErrorKindInternal, "too many arguments in message '{{ .value }}'")

    ErrDataContainer = NewFactory(
        "errDataContainer", ErrorKindInternalNotice, "data: '{{ .value }}'")
)
