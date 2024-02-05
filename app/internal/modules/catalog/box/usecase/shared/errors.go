package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrBoxNotFound = NewFactory(
		"errCatalogBoxNotFound", ErrorKindUser, "box with ID={{ .id }} not found")

	FactoryErrBoxArticleAlreadyExists = NewFactory(
		"errCatalogBoxArticleAlreadyExists", ErrorKindUser, "box article '{{ .name }}' already exist")
)
