package mrerr

var (
    ErrStorageConnectionIsAlreadyCreated = NewFactory(
        "errStorageConnectionIsAlreadyCreated", ErrorKindInternal, "connection '{{ .name }}' is already created")

    ErrStorageConnectionIsNotOpened = NewFactory(
        "errStorageConnectionIsNotOpened", ErrorKindInternal, "connection '{{ .name }}' is not opened")

    ErrStorageConnectionFailed = NewFactory(
        "errStorageConnectionFailed", ErrorKindSystem, "connection '{{ .name }}' is failed")

    ErrStorageQueryFailed = NewFactory(
        "errStorageQueryFailed", ErrorKindInternal, "query is failed")

    ErrStorageFetchDataFailed = NewFactory(
        "errStorageFetchDataFailed", ErrorKindInternal, "fetching data is failed")

    ErrStorageFetchedInvalidData = NewFactory(
        "errStorageFetchedInvalidData", ErrorKindInternal, "fetched data '{{ .value }}' is invalid")

    ErrStorageNoRowFound = NewFactory(
        "errStorageNoRowFound", ErrorKindInternalNotice, "no row found")

    ErrStorageRowsNotAffected = NewFactory(
        "errStorageRowsNotAffected", ErrorKindInternalNotice, "rows not affected")
)
