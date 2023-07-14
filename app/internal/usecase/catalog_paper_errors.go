package usecase

import "print-shop-back/pkg/mrerr"

var (
    ErrCatalogPaperColorNotFound = mrerr.NewFactory(
        "errErrCatalogPaperColorNotFound", mrerr.ErrorKindUser, "paper color with ID={{ .id }} not found")

    ErrCatalogPaperFactureNotFound = mrerr.NewFactory(
        "errErrCatalogPaperFactureNotFound", mrerr.ErrorKindUser, "paper facture with ID={{ .id }} not found")

    ErrCatalogPaperArticleAlreadyExists = mrerr.NewFactory(
        "errCatalogPaperArticleAlreadyExists", mrerr.ErrorKindUser, "paper article '{{ .name }}' is already exists")
)
