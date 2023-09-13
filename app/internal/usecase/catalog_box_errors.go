package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    ErrCatalogBoxArticleAlreadyExists = NewFactory(
        "errCatalogBoxArticleAlreadyExists", ErrorKindUser, "box article '{{ .name }}' is already exists")
)
