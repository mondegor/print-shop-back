package usecase

import "calc-user-data-back-adm/pkg/mrerr"

var (
    ErrCatalogLaminateTypeNotFound = mrerr.NewFactory(
        "errErrCatalogLaminateTypeNotFound", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} not found")

    ErrCatalogLaminateArticleAlreadyExists = mrerr.NewFactory(
        "errCatalogLaminateArticleAlreadyExists", mrerr.ErrorKindUser, "laminate article '{{ .name }}' is already exists")
)
