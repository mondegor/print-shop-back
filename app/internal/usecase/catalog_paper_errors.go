package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    ErrCatalogPaperColorNotFound = NewFactory(
        "errErrCatalogPaperColorNotFound", ErrorKindUser, "paper color with ID={{ .id }} not found")

    ErrCatalogPaperFactureNotFound = NewFactory(
        "errErrCatalogPaperFactureNotFound", ErrorKindUser, "paper facture with ID={{ .id }} not found")

    ErrCatalogPaperArticleAlreadyExists = NewFactory(
        "errCatalogPaperArticleAlreadyExists", ErrorKindUser, "paper article '{{ .name }}' is already exists")
)
