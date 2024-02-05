package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrPaperNotFound = NewFactory(
		"errCatalogPaperNotFound", ErrorKindUser, "paper with ID={{ .id }} not found")

	FactoryErrPaperArticleAlreadyExists = NewFactory(
		"errCatalogPaperArticleAlreadyExists", ErrorKindUser, "paper article '{{ .name }}' already exist")
)
