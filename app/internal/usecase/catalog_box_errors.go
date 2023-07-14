package usecase

import "calc-user-data-back-adm/pkg/mrerr"

var (
    ErrCatalogBoxArticleAlreadyExists = mrerr.NewFactory(
        "errCatalogBoxArticleAlreadyExists", mrerr.ErrorKindUser, "box article '{{ .name }}' is already exists")
)
