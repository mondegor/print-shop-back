package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    ErrCatalogLaminateTypeNotFound = NewFactory(
        "errErrCatalogLaminateTypeNotFound", ErrorKindUser, "laminate type with ID={{ .id }} not found")

    ErrCatalogLaminateArticleAlreadyExists = NewFactory(
        "errCatalogLaminateArticleAlreadyExists", ErrorKindUser, "laminate article '{{ .name }}' is already exists")
)
