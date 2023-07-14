package usecase

import "print-shop-back/pkg/mrerr"

var (
    ErrCatalogBoxArticleAlreadyExists = mrerr.NewFactory(
        "errCatalogBoxArticleAlreadyExists", mrerr.ErrorKindUser, "box article '{{ .name }}' is already exists")
)
