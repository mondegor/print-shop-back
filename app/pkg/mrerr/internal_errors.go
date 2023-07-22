package mrerr

var (
    ErrInternal = NewFactory(
        ErrorCodeInternal, ErrorKindInternal, "internal server error")

    ErrInternalNilPointer = NewFactory(
        "errInternalNilPointer", ErrorKindInternal, "nil pointer")

    ErrInternalTypeAssertion = NewFactory(
        "errInternalTypeAssertion", ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

    ErrInternalParseData = NewFactory(
        "errInternalParseData", ErrorKindInternal, "data '{{ .name1 }}' parsed to {{ .name2 }} with error")

    ErrInternalMapValueNotFound = NewFactory(
        "errInternalMapValueNotFound", ErrorKindInternal, "'{{ .value }}' is not found in map {{ .name }}")

    ErrInternalMessageNotEnoughArguments = NewFactory(
        "errInternalMessageNotEnoughArguments", ErrorKindInternal, "Not enough arguments in message '{{ .value }}'")

    ErrInternalMessageTooManyArguments = NewFactory(
        "errInternalMessageTooManyArguments", ErrorKindInternal, "Too many arguments in message '{{ .value }}'")
)
