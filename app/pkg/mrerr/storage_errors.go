package mrerr

var (
    ErrStorageConnectionIsAlreadyCreated = NewFactory(
        "errStorageConnectionIsAlreadyCreated", ErrorKindInternal, "connection '{{ .name }}' is already created")

    ErrStorageConnectionFailed = NewFactory(
        "errStorageConnectionFailed", ErrorKindSystem, "connection '{{ .name }}' is failed")

    ErrStorageQueryFailed = NewFactory(
        "errStorageQueryFailed", ErrorKindInternal, "query is failed")

    ErrStorageFetchDataFailed = NewFactory(
        "errStorageFetchDataFailed", ErrorKindInternal, "fetching data is failed")

    ErrStorageNoRowFound = NewFactory(
        "errStorageNoRowFound", ErrorKindInternalNotice, "no row found")

    ErrStorageRowsNotAffected = NewFactory(
        "errStorageRowsNotAffected", ErrorKindInternalNotice, "rows not affected")
)
