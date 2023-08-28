package mrerr

var (
    ErrServiceIncorrectInputData = NewFactory(
        "errServiceIncorrectData", ErrorKindInternal, "data '{{ .data }}' is incorrect")

    ErrServiceEntityTemporarilyUnavailable = NewFactory(
        "errServiceEntityTemporarilyUnavailable", ErrorKindSystem, "entity '{{ .name }}' is temporarily unavailable")

    ErrServiceEntityNotFound = NewFactory(
        "errServiceEntityNotFound", ErrorKindInternalNotice, "entity '{{ .name }}' is not found")

    ErrServiceEntityNotCreated = NewFactory(
        "errServiceEntityNotCreated", ErrorKindSystem, "entity '{{ .name }}' is not created")

    ErrServiceEntityNotUpdated = NewFactory(
        "errServiceEntityNotUpdated", ErrorKindSystem, "entity '{{ .name }}' is not updated")

    ErrServiceEntityNotRemoved = NewFactory(
        "errServiceEntityNotRemoved", ErrorKindSystem, "entity '{{ .name }}' is not removed")

    ErrServiceIncorrectSwitchStatus = NewFactory(
        "errServiceIncorrectSwitchStatus", ErrorKindInternal, "incorrect switch status: '{{ .currentStatus }}' -> '{{ .statusTo }}' for entity '{{ .name }}(ID={{ .id }})'")
)
