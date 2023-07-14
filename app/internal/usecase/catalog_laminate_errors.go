package usecase

import "print-shop-back/pkg/mrerr"

var (
    ErrCatalogLaminateTypeNotFound = mrerr.NewFactory(
        "errErrCatalogLaminateTypeNotFound", mrerr.ErrorKindUser, "laminate type with ID={{ .id }} not found")

    ErrCatalogLaminateArticleAlreadyExists = mrerr.NewFactory(
        "errCatalogLaminateArticleAlreadyExists", mrerr.ErrorKindUser, "laminate article '{{ .name }}' is already exists")
)
