package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrLaminateNotFound = NewFactory(
		"errCatalogLaminateNotFound", ErrorKindUser, "laminate with ID={{ .id }} not found")

	FactoryErrLaminateArticleAlreadyExists = NewFactory(
		"errCatalogLaminateArticleAlreadyExists", ErrorKindUser, "laminate article '{{ .name }}' already exist")
)
