package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrBoxNotFound = NewFactory(
		"errCatalogBoxNotFound", ErrorKindUser, "box with ID={{ .id }} not found")

	FactoryErrBoxArticleAlreadyExists = NewFactory(
		"errCatalogBoxArticleAlreadyExists", ErrorKindUser, "box article '{{ .name }}' already exist")

	FactoryErrLaminateNotFound = NewFactory(
		"errCatalogLaminateNotFound", ErrorKindUser, "laminate with ID={{ .id }} not found")

	FactoryErrLaminateArticleAlreadyExists = NewFactory(
		"errCatalogLaminateArticleAlreadyExists", ErrorKindUser, "laminate article '{{ .name }}' already exist")

	FactoryErrPaperNotFound = NewFactory(
		"errCatalogPaperNotFound", ErrorKindUser, "paper with ID={{ .id }} not found")

	FactoryErrPaperArticleAlreadyExists = NewFactory(
		"errCatalogPaperArticleAlreadyExists", ErrorKindUser, "paper article '{{ .name }}' already exist")
)
